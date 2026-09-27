package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aisha-platform/aisha/apps/api/internal/features/moderation/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(p *pgxpool.Pool) *PostgresRepository { return &PostgresRepository{pool: p} }
func (r *PostgresRepository) ListQueue(ctx context.Context, status string, limit, offset int) ([]domain.QueueItem, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM products p WHERE ($1='' OR p.status=$1)`, status).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count moderation queue: %w", err)
	}
	rows, err := r.pool.Query(ctx, `SELECT COALESCE(s.id::text,''),p.id,COALESCE(p.product_code,''),COALESCE((SELECT t.name FROM product_translations t WHERE t.product_id=p.id AND t.locale='en' LIMIT 1),(SELECT t.name FROM product_translations t WHERE t.product_id=p.id ORDER BY t.locale LIMIT 1),p.product_type),COALESCE(s.version,0),COALESCE(s.snapshot,'{}'::jsonb),p.status,p.price_minor,p.currency,COALESCE(a.public_display_name,''),COALESCE(w.name,'') FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN workshops w ON w.id=p.workshop_id LEFT JOIN LATERAL (SELECT ps.id,ps.version,ps.snapshot,ps.submitted_at FROM product_submissions ps WHERE ps.product_id=p.id ORDER BY ps.version DESC LIMIT 1) s ON true WHERE ($1='' OR p.status=$1) ORDER BY COALESCE(s.submitted_at,p.updated_at) ASC,p.id LIMIT $2 OFFSET $3`, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []domain.QueueItem{}
	for rows.Next() {
		var i domain.QueueItem
		if err := rows.Scan(&i.SubmissionID, &i.ProductID, &i.ProductCode, &i.ProductName, &i.Version, &i.Snapshot, &i.ProductStatus, &i.PriceMinor, &i.Currency, &i.ArtisanName, &i.WorkshopName); err != nil {
			return nil, 0, err
		}
		i.Media, err = loadQueueMedia(ctx, r.pool, i.ProductID)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, i)
	}
	return items, total, rows.Err()
}

func loadQueueMedia(ctx context.Context, pool *pgxpool.Pool, productID string) ([]domain.Media, error) {
	rows, err := pool.Query(ctx, `SELECT id,media_kind,COALESCE(original_filename,''),media_type,size_bytes,COALESCE(alt_text,''),visibility,object_key FROM product_media WHERE product_id=$1 ORDER BY sort_order,created_at`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.Media{}
	for rows.Next() {
		var item domain.Media
		if err = rows.Scan(&item.ID, &item.MediaKind, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.AltText, &item.Visibility, &item.ObjectKey); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
func (r *PostgresRepository) Decide(ctx context.Context, actor string, in domain.DecisionInput) (domain.QueueItem, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.QueueItem{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var productID, productCode, status, artisanID, workshopID, currency string
	var price int64
	var made bool
	var submissionID *string
	if in.Action == "ACTIVATE" || in.Action == "ARCHIVE" {
		err = tx.QueryRow(ctx, `SELECT id,COALESCE(product_code,''),status,artisan_profile_id,workshop_id,price_minor,currency,made_to_order_eligible FROM products WHERE id=$1 FOR UPDATE`, in.ID).Scan(&productID, &productCode, &status, &artisanID, &workshopID, &price, &currency, &made)
	} else {
		err = tx.QueryRow(ctx, `SELECT p.id,COALESCE(p.product_code,''),p.status,p.artisan_profile_id,p.workshop_id,p.price_minor,p.currency,p.made_to_order_eligible,s.id FROM product_submissions s JOIN products p ON p.id=s.product_id WHERE s.id=$1 FOR UPDATE OF p`, in.ID).Scan(&productID, &productCode, &status, &artisanID, &workshopID, &price, &currency, &made, &submissionID)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.QueueItem{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.QueueItem{}, err
	}
	previous := status
	next := ""
	switch in.Action {
	case "APPROVE":
		if status != "PENDING_REVIEW" {
			return domain.QueueItem{}, domain.ErrInvalidTransition
		}
		next = "APPROVED"
	case "REQUEST_CHANGES":
		if status != "PENDING_REVIEW" {
			return domain.QueueItem{}, domain.ErrInvalidTransition
		}
		next = "CHANGES_REQUESTED"
	case "REJECT":
		if status != "PENDING_REVIEW" {
			return domain.QueueItem{}, domain.ErrInvalidTransition
		}
		next = "ARCHIVED"
	case "SUSPEND":
		if status != "ACTIVE" {
			return domain.QueueItem{}, domain.ErrInvalidTransition
		}
		next = "SUSPENDED"
	case "ACTIVATE":
		if status != "APPROVED" {
			return domain.QueueItem{}, domain.ErrInvalidTransition
		}
		next = "ACTIVE"
	case "ARCHIVE":
		if status == "ARCHIVED" {
			return domain.QueueItem{}, domain.ErrInvalidTransition
		}
		next = "ARCHIVED"
	}
	if next == "APPROVED" || next == "ACTIVE" {
		var ready bool
		if err = tx.QueryRow(ctx, `SELECT a.status='APPROVED' AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE') AND w.status='ACTIVE' AND p.price_minor>0 AND p.currency ~ '^[A-Z]{3}$' AND EXISTS(SELECT 1 FROM product_media m WHERE m.product_id=p.id) AND (p.made_to_order_eligible OR COALESCE((SELECT SUM(quantity_delta) FROM inventory_movements i WHERE i.product_id=p.id AND i.stock_bucket='AVAILABLE'),0)>0) FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN workshops w ON w.id=p.workshop_id WHERE p.id=$1`, productID).Scan(&ready); err != nil {
			return domain.QueueItem{}, err
		}
		if !ready {
			if next == "ACTIVE" {
				return domain.QueueItem{}, domain.ErrActivationNotReady
			}
			next = "APPROVED"
		}
	}
	if err = tx.QueryRow(ctx, `UPDATE products SET status=$2,product_code=CASE WHEN $2='APPROVED' AND product_code IS NULL THEN 'AISHA-' || upper(substr(replace(id::text, '-', ''), 1, 10)) ELSE product_code END,published_at=CASE WHEN $2='ACTIVE' THEN COALESCE(published_at,CURRENT_TIMESTAMP) WHEN $2='ARCHIVED' THEN NULL ELSE published_at END,updated_at=CURRENT_TIMESTAMP WHERE id=$1 RETURNING COALESCE(product_code,'')`, productID, next).Scan(&productCode); err != nil {
		return domain.QueueItem{}, err
	}
	if next == "ACTIVE" {
		if _, err = tx.Exec(ctx, `UPDATE product_media SET visibility='PUBLIC' WHERE product_id=$1`, productID); err != nil {
			return domain.QueueItem{}, err
		}
	}
	if _, err = tx.Exec(ctx, `INSERT INTO product_moderation_decisions(product_id,submission_id,actor_user_id,action,previous_status,new_status,reason) VALUES($1,$2,$3,$4,$5,$6,NULLIF($7,''))`, productID, submissionID, actor, in.Action, previous, next, in.Reason); err != nil {
		return domain.QueueItem{}, err
	}
	payload, _ := json.Marshal(map[string]any{"productId": productID, "action": in.Action, "status": next, "reason": in.Reason})
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,reason,previous_state,new_state) VALUES($1,$2,$3,$4,NULLIF($5,''),jsonb_build_object('status',$6::text),jsonb_build_object('status',$7::text))`, `PRODUCT_MODERATION_DECIDED`, actor, "product", productID, in.Reason, previous, next); err != nil {
		return domain.QueueItem{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,$2,$3,$4)`, `PRODUCT_MODERATION_DECIDED`, `product`, productID, payload); err != nil {
		return domain.QueueItem{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.QueueItem{}, err
	}
	return domain.QueueItem{ProductID: productID, ProductCode: productCode, ProductStatus: next, PriceMinor: price, Currency: strings.TrimSpace(currency)}, nil
}
