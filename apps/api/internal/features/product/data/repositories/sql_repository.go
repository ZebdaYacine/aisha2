package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aisha-platform/aisha/apps/api/internal/features/product/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

type rowQuerier interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

const productColumns = `p.id,p.artisan_profile_id,p.category_id,p.product_type,p.status,p.price_minor,p.currency,
COALESCE(p.materials,''),COALESCE(p.production_method,''),COALESCE(p.intended_use,''),COALESCE(p.dimensions,''),
p.weight_grams,COALESCE(p.country_of_origin,''),COALESCE(p.region_of_origin,''),p.eco_friendly_verified,
p.fair_trade_verified,p.made_to_order_eligible,
COALESCE((SELECT json_agg(json_build_object('locale',t.locale,'name',t.name,'description',t.description,'story',COALESCE(t.story,''),'culturalContext',COALESCE(t.cultural_context,'')) ORDER BY t.locale) FROM product_translations t WHERE t.product_id=p.id),'[]')`

func (r *PostgresRepository) Create(ctx context.Context, userID string, input domain.Input) (domain.Product, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Product{}, fmt.Errorf("begin product creation: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err = ensureApprovedArtisan(ctx, tx, userID); err != nil {
		return domain.Product{}, err
	}
	if err = ensureCategory(ctx, tx, input.CategoryID); err != nil {
		return domain.Product{}, err
	}
	var id string
	err = tx.QueryRow(ctx, `INSERT INTO products(artisan_profile_id,category_id,product_type,status,price_minor,currency,materials,production_method,intended_use,dimensions,weight_grams,country_of_origin,region_of_origin,eco_friendly_verified,fair_trade_verified,made_to_order_eligible) SELECT a.id,$2,$3,'DRAFT',$4,$5,NULLIF($6,''),NULLIF($7,''),NULLIF($8,''),NULLIF($9,''),$10,NULLIF($11,''),NULLIF($12,''),$13,$14,$15 FROM artisan_profiles a WHERE a.user_id=$1 AND a.status='APPROVED' RETURNING id`, userID, input.CategoryID, input.ProductType, input.PriceMinor, input.Currency, input.Materials, input.ProductionMethod, input.IntendedUse, input.Dimensions, input.WeightGrams, input.CountryOfOrigin, input.RegionOfOrigin, input.EcoFriendlyVerified, input.FairTradeVerified, input.MadeToOrderEligible).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, domain.ErrArtisanNotApproved
	}
	if err != nil {
		return domain.Product{}, fmt.Errorf("insert product: %w", err)
	}
	if err = replaceTranslations(ctx, tx, id, input.Translations); err != nil {
		return domain.Product{}, err
	}
	if err = productEvent(ctx, tx, "PRODUCT_DRAFT_CREATED", userID, id, "", "", "DRAFT"); err != nil {
		return domain.Product{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Product{}, fmt.Errorf("commit product creation: %w", err)
	}
	return r.GetOwned(ctx, userID, id)
}

func (r *PostgresRepository) ListOwned(ctx context.Context, userID, status string, limit, offset int) ([]domain.Product, int, error) {
	rows, err := r.pool.Query(ctx, `SELECT p.id,count(*) OVER() FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id WHERE a.user_id=$1 AND ($2='' OR p.status=$2) ORDER BY p.updated_at DESC LIMIT $3 OFFSET $4`, userID, status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list artisan products: %w", err)
	}
	defer rows.Close()
	var ids []string
	total := 0
	for rows.Next() {
		var id string
		if err = rows.Scan(&id, &total); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	items := make([]domain.Product, 0, len(ids))
	for _, id := range ids {
		item, loadErr := r.GetOwned(ctx, userID, id)
		if loadErr != nil {
			return nil, 0, loadErr
		}
		items = append(items, item)
	}
	return items, total, nil
}

func (r *PostgresRepository) GetOwned(ctx context.Context, userID, id string) (domain.Product, error) {
	item, err := loadProduct(ctx, r.pool, `WHERE p.id=$1 AND a.user_id=$2`, id, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, domain.ErrNotFound
	}
	return item, err
}

func (r *PostgresRepository) Update(ctx context.Context, userID, id string, input domain.Input) (domain.Product, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Product{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var status string
	err = tx.QueryRow(ctx, `SELECT p.status FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id WHERE p.id=$2 AND a.user_id=$1 FOR UPDATE`, userID, id).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Product{}, err
	}
	if status != "DRAFT" && status != "CHANGES_REQUESTED" {
		return domain.Product{}, domain.ErrNotEditable
	}
	if err = ensureCategory(ctx, tx, input.CategoryID); err != nil {
		return domain.Product{}, err
	}
	_, err = tx.Exec(ctx, `UPDATE products SET category_id=$3,product_type=$4,price_minor=$5,currency=$6,materials=NULLIF($7,''),production_method=NULLIF($8,''),intended_use=NULLIF($9,''),dimensions=NULLIF($10,''),weight_grams=$11,country_of_origin=NULLIF($12,''),region_of_origin=NULLIF($13,''),eco_friendly_verified=$14,fair_trade_verified=$15,made_to_order_eligible=$16,updated_at=CURRENT_TIMESTAMP WHERE id=$2 AND artisan_profile_id=(SELECT id FROM artisan_profiles WHERE user_id=$1)`, userID, id, input.CategoryID, input.ProductType, input.PriceMinor, input.Currency, input.Materials, input.ProductionMethod, input.IntendedUse, input.Dimensions, input.WeightGrams, input.CountryOfOrigin, input.RegionOfOrigin, input.EcoFriendlyVerified, input.FairTradeVerified, input.MadeToOrderEligible)
	if err != nil {
		return domain.Product{}, fmt.Errorf("update product: %w", err)
	}
	if err = replaceTranslations(ctx, tx, id, input.Translations); err != nil {
		return domain.Product{}, err
	}
	if err = productEvent(ctx, tx, "PRODUCT_DRAFT_UPDATED", userID, id, "", status, status); err != nil {
		return domain.Product{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Product{}, err
	}
	return r.GetOwned(ctx, userID, id)
}

func (r *PostgresRepository) Submit(ctx context.Context, userID, id string) (domain.Product, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Product{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var status string
	err = tx.QueryRow(ctx, `SELECT p.status FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id WHERE p.id=$2 AND a.user_id=$1 FOR UPDATE`, userID, id).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Product{}, err
	}
	if status != "DRAFT" && status != "CHANGES_REQUESTED" {
		return domain.Product{}, domain.ErrInvalidTransition
	}
	item, err := loadProduct(ctx, tx, `WHERE p.id=$1 AND a.user_id=$2`, id, userID)
	if err != nil {
		return domain.Product{}, err
	}
	snapshot, err := json.Marshal(item)
	if err != nil {
		return domain.Product{}, fmt.Errorf("encode product submission: %w", err)
	}
	var version int
	if err = tx.QueryRow(ctx, `SELECT COALESCE(MAX(version),0)+1 FROM product_submissions WHERE product_id=$1`, id).Scan(&version); err != nil {
		return domain.Product{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE products SET status='PENDING_REVIEW',published_at=NULL,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id); err != nil {
		return domain.Product{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO product_submissions(product_id,version,submitted_by_user_id,snapshot) VALUES($1,$2,$3,$4)`, id, version, userID, snapshot); err != nil {
		return domain.Product{}, fmt.Errorf("insert product submission: %w", err)
	}
	if err = productEvent(ctx, tx, "PRODUCT_SUBMITTED", userID, id, "", status, "PENDING_REVIEW"); err != nil {
		return domain.Product{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Product{}, err
	}
	return r.GetOwned(ctx, userID, id)
}

func (r *PostgresRepository) AddMedia(ctx context.Context, userID, productID string, media domain.Media) (domain.Media, error) {
	var result domain.Media
	err := r.pool.QueryRow(ctx, `INSERT INTO product_media(product_id,media_kind,object_key,original_filename,media_type,size_bytes,checksum_sha256,alt_text,visibility) SELECT p.id,$3,$4,NULLIF($5,''),$6,$7,$8,NULLIF($9,''),'PRIVATE' FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id WHERE p.id=$2 AND a.user_id=$1 AND p.status IN ('DRAFT','CHANGES_REQUESTED') RETURNING id,product_id,media_kind,object_key,COALESCE(original_filename,''),media_type,size_bytes,COALESCE(alt_text,''),sort_order,visibility`, userID, productID, media.MediaKind, media.ObjectKey, media.OriginalFilename, media.MediaType, media.SizeBytes, media.Checksum, media.AltText).Scan(&result.ID, &result.ProductID, &result.MediaKind, &result.ObjectKey, &result.OriginalFilename, &result.MediaType, &result.SizeBytes, &result.AltText, &result.SortOrder, &result.Visibility)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Media{}, domain.ErrNotEditable
	}
	if err != nil {
		return domain.Media{}, fmt.Errorf("insert product media: %w", err)
	}
	return result, nil
}

func (r *PostgresRepository) DeleteMedia(ctx context.Context, userID, productID, mediaID string) (domain.Media, error) {
	var result domain.Media
	err := r.pool.QueryRow(ctx, `DELETE FROM product_media m USING products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id WHERE m.id=$3 AND m.product_id=p.id AND p.id=$2 AND a.user_id=$1 AND p.status IN ('DRAFT','CHANGES_REQUESTED') RETURNING m.id,m.product_id,m.object_key`, userID, productID, mediaID).Scan(&result.ID, &result.ProductID, &result.ObjectKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Media{}, domain.ErrMediaNotFound
	}
	if err != nil {
		return domain.Media{}, fmt.Errorf("delete product media: %w", err)
	}
	return result, nil
}

func loadProduct(ctx context.Context, q rowQuerier, where string, args ...any) (domain.Product, error) {
	var item domain.Product
	var translations []byte
	var weight pgtype.Int4
	err := q.QueryRow(ctx, `SELECT `+productColumns+` FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id `+where, args...).Scan(&item.ID, &item.ArtisanID, &item.CategoryID, &item.ProductType, &item.Status, &item.PriceMinor, &item.Currency, &item.Materials, &item.ProductionMethod, &item.IntendedUse, &item.Dimensions, &weight, &item.CountryOfOrigin, &item.RegionOfOrigin, &item.EcoFriendlyVerified, &item.FairTradeVerified, &item.MadeToOrderEligible, &translations)
	if err != nil {
		return domain.Product{}, err
	}
	if weight.Valid {
		value := int(weight.Int32)
		item.WeightGrams = &value
	}
	if err = json.Unmarshal(translations, &item.Translations); err != nil {
		return domain.Product{}, fmt.Errorf("decode product translations: %w", err)
	}
	item.Media, err = loadMedia(ctx, q, item.ID)
	return item, err
}

func loadMedia(ctx context.Context, q rowQuerier, productID string) ([]domain.Media, error) {
	rows, err := q.Query(ctx, `SELECT id,product_id,media_kind,object_key,COALESCE(original_filename,''),media_type,size_bytes,COALESCE(alt_text,''),sort_order,visibility FROM product_media WHERE product_id=$1 ORDER BY sort_order,created_at`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.Media{}
	for rows.Next() {
		var media domain.Media
		if err = rows.Scan(&media.ID, &media.ProductID, &media.MediaKind, &media.ObjectKey, &media.OriginalFilename, &media.MediaType, &media.SizeBytes, &media.AltText, &media.SortOrder, &media.Visibility); err != nil {
			return nil, err
		}
		items = append(items, media)
	}
	return items, rows.Err()
}

func ensureApprovedArtisan(ctx context.Context, tx pgx.Tx, userID string) error {
	var id string
	err := tx.QueryRow(ctx, `SELECT id FROM artisan_profiles WHERE user_id=$1 AND status='APPROVED'`, userID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrArtisanNotApproved
	}
	return err
}

func ensureCategory(ctx context.Context, tx pgx.Tx, id string) error {
	var exists bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1 AND is_active=true)`, id).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return domain.ErrValidation
	}
	return nil
}

func replaceTranslations(ctx context.Context, tx pgx.Tx, productID string, translations []domain.Translation) error {
	if _, err := tx.Exec(ctx, `DELETE FROM product_translations WHERE product_id=$1`, productID); err != nil {
		return err
	}
	for _, translation := range translations {
		if _, err := tx.Exec(ctx, `INSERT INTO product_translations(product_id,locale,name,description,story,cultural_context) VALUES($1,$2,$3,$4,NULLIF($5,''),NULLIF($6,''))`, productID, translation.Locale, translation.Name, translation.Description, translation.Story, translation.CulturalContext); err != nil {
			return fmt.Errorf("replace product translation: %w", err)
		}
	}
	return nil
}

func productEvent(ctx context.Context, tx pgx.Tx, eventType, actor, productID, reason, previous, next string) error {
	if _, err := tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,reason,previous_state,new_state) VALUES($1,$2,'product',$3,NULLIF($4,''),jsonb_build_object('status',NULLIF($5,'')),jsonb_build_object('status',$6))`, eventType, actor, productID, reason, previous, next); err != nil {
		return fmt.Errorf("insert product audit event: %w", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,'product',$2,jsonb_build_object('productId',$2::text,'status',$3))`, eventType, productID, next); err != nil {
		return fmt.Errorf("insert product outbox event: %w", err)
	}
	return nil
}
