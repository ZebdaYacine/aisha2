package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) ListUserMedia(ctx context.Context, limit, offset int) ([]domain.UserMedia, int, error) {
	rows, err := r.pool.Query(ctx, `WITH media AS (
SELECT am.id AS media_id,u.id AS user_id,COALESCE(u.display_name,'') AS user_display_name,COALESCE(u.email,'') AS user_email,'PROFILE_MEDIA' AS media_kind,NULLIF(am.media_kind,'') AS document_type,COALESCE(am.original_filename,'') AS original_filename,am.media_type,am.size_bytes,am.object_key,am.created_at
FROM artisan_media am JOIN artisan_profiles ap ON ap.id=am.artisan_profile_id JOIN users u ON u.id=ap.user_id
UNION ALL
SELECT ad.id AS media_id,u.id AS user_id,COALESCE(u.display_name,'') AS user_display_name,COALESCE(u.email,'') AS user_email,'DOCUMENT' AS media_kind,ad.document_type,COALESCE(ad.original_filename,'') AS original_filename,ad.media_type,ad.size_bytes,ad.object_key,ad.created_at
FROM artisan_documents ad JOIN artisan_profiles ap ON ap.id=ad.artisan_profile_id JOIN users u ON u.id=ap.user_id
)
SELECT media_id,user_id,user_display_name,user_email,media_kind,document_type,original_filename,media_type,size_bytes,object_key,created_at,count(*) OVER() FROM media ORDER BY created_at DESC,media_id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []domain.UserMedia{}
	total := 0
	for rows.Next() {
		var item domain.UserMedia
		if err = rows.Scan(&item.ID, &item.UserID, &item.UserDisplayName, &item.UserEmail, &item.MediaKind, &item.DocumentType, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.ObjectKey, &item.CreatedAt, &total); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) ListProductMedia(ctx context.Context, limit, offset int) ([]domain.ProductMedia, int, error) {
	rows, err := r.pool.Query(ctx, `SELECT pm.id,pm.product_id,COALESCE((SELECT t.name FROM product_translations t WHERE t.product_id=p.id AND t.locale='en' LIMIT 1),(SELECT t.name FROM product_translations t WHERE t.product_id=p.id ORDER BY t.locale LIMIT 1),p.product_type),p.status,COALESCE(a.public_display_name,''),pm.media_kind,COALESCE(pm.original_filename,''),pm.media_type,pm.size_bytes,COALESCE(pm.alt_text,''),pm.visibility,pm.object_key,pm.created_at,count(*) OVER() FROM product_media pm JOIN products p ON p.id=pm.product_id JOIN artisan_profiles a ON a.id=p.artisan_profile_id ORDER BY pm.created_at DESC,pm.id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []domain.ProductMedia{}
	total := 0
	for rows.Next() {
		var item domain.ProductMedia
		if err = rows.Scan(&item.ID, &item.ProductID, &item.ProductName, &item.ProductStatus, &item.ArtisanName, &item.MediaKind, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.AltText, &item.Visibility, &item.ObjectKey, &item.CreatedAt, &total); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) DeleteUserMedia(ctx context.Context, actorID, mediaID string) (domain.UserMedia, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.UserMedia{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var item domain.UserMedia
	err = tx.QueryRow(ctx, `DELETE FROM artisan_media WHERE id=$1 RETURNING id,object_key`, mediaID).Scan(&item.ID, &item.ObjectKey)
	if errors.Is(err, pgx.ErrNoRows) {
		err = tx.QueryRow(ctx, `DELETE FROM artisan_documents WHERE id=$1 RETURNING id,object_key`, mediaID).Scan(&item.ID, &item.ObjectKey)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.UserMedia{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.UserMedia{}, fmt.Errorf("delete user media: %w", err)
	}
	if err = recordMediaDeletion(ctx, tx, actorID, "user_media", item.ID); err != nil {
		return domain.UserMedia{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.UserMedia{}, err
	}
	return item, nil
}

func (r *PostgresRepository) DeleteProductMedia(ctx context.Context, actorID, mediaID string) (domain.ProductMedia, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.ProductMedia{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var item domain.ProductMedia
	err = tx.QueryRow(ctx, `DELETE FROM product_media WHERE id=$1 RETURNING id,product_id,object_key`, mediaID).Scan(&item.ID, &item.ProductID, &item.ObjectKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ProductMedia{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.ProductMedia{}, fmt.Errorf("delete product media: %w", err)
	}
	if err = recordMediaDeletion(ctx, tx, actorID, "product_media", item.ID); err != nil {
		return domain.ProductMedia{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.ProductMedia{}, err
	}
	return item, nil
}

func recordMediaDeletion(ctx context.Context, tx pgx.Tx, actorID, targetType, mediaID string) error {
	if _, err := tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,reason) VALUES('MEDIA_DELETED',$1,$2,$3,'Deleted by administrator')`, actorID, targetType, mediaID); err != nil {
		return fmt.Errorf("record media deletion audit: %w", err)
	}
	return nil
}

func (r *PostgresRepository) ListUsers(ctx context.Context, role, status string, limit, offset int) ([]domain.User, int, error) {
	rows, err := r.pool.Query(ctx, `SELECT u.id,COALESCE(u.email,''),COALESCE(u.phone,''),COALESCE(u.display_name,''),u.status,u.created_at,COALESCE(array_agg(r.code) FILTER(WHERE r.code IS NOT NULL),'{}'),count(*) OVER() FROM users u LEFT JOIN user_roles ur ON ur.user_id=u.id LEFT JOIN roles r ON r.id=ur.role_id WHERE ($1='' OR u.status=$1) AND ($2='' OR (($2='worker' AND EXISTS(SELECT 1 FROM user_roles ur2 JOIN roles r2 ON r2.id=ur2.role_id WHERE ur2.user_id=u.id AND r2.code IN ('moderator','warehouse_agent'))) OR ($2='user' AND EXISTS(SELECT 1 FROM user_roles ur3 JOIN roles r3 ON r3.id=ur3.role_id WHERE ur3.user_id=u.id AND r3.code IN ('customer','artisan'))) OR ($2 NOT IN ('worker','user') AND EXISTS(SELECT 1 FROM user_roles ur4 JOIN roles r4 ON r4.id=ur4.role_id WHERE ur4.user_id=u.id AND r4.code=$2)))) GROUP BY u.id ORDER BY u.created_at DESC LIMIT $3 OFFSET $4`, status, role, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []domain.User{}
	total := 0
	for rows.Next() {
		var item domain.User
		if err = rows.Scan(&item.ID, &item.Email, &item.Phone, &item.DisplayName, &item.Status, &item.CreatedAt, &item.Roles, &total); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) CreateUser(ctx context.Context, actor string, input domain.UserInput, assignedBy string) (domain.User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var id string
	if err = tx.QueryRow(ctx, `INSERT INTO users(email,phone,password_hash,display_name) VALUES($1,NULLIF($2,''),$3,$4) RETURNING id`, strings.ToLower(strings.TrimSpace(input.Email)), strings.TrimSpace(input.Phone), input.Password, strings.TrimSpace(input.DisplayName)).Scan(&id); err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, domain.ErrValidation
		}
		return domain.User{}, fmt.Errorf("create admin user: %w", err)
	}
	if err = replaceUserRoles(ctx, tx, id, assignedBy, input.Roles); err != nil {
		return domain.User{}, err
	}
	if err = recordUserAudit(ctx, tx, actor, id, "ADMIN_USER_CREATED"); err != nil {
		return domain.User{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.User{}, err
	}
	return r.adminUser(ctx, id)
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, actor, userID string, input domain.UserInput) (domain.User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `UPDATE users SET email=$2,phone=NULLIF($3,''),display_name=$4,password_hash=CASE WHEN $5='' THEN password_hash ELSE $5 END,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, userID, strings.ToLower(strings.TrimSpace(input.Email)), strings.TrimSpace(input.Phone), strings.TrimSpace(input.DisplayName), input.Password)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, domain.ErrValidation
		}
		return domain.User{}, err
	}
	if result.RowsAffected() != 1 {
		return domain.User{}, domain.ErrNotFound
	}
	if err = replaceUserRoles(ctx, tx, userID, actor, input.Roles); err != nil {
		return domain.User{}, err
	}
	if err = recordUserAudit(ctx, tx, actor, userID, "ADMIN_USER_UPDATED"); err != nil {
		return domain.User{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.User{}, err
	}
	return r.adminUser(ctx, userID)
}

func replaceUserRoles(ctx context.Context, tx pgx.Tx, userID, assignedBy string, roles []string) error {
	if _, err := tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id=$1`, userID); err != nil {
		return err
	}
	for _, role := range roles {
		result, err := tx.Exec(ctx, `INSERT INTO user_roles(user_id,role_id,assigned_by_user_id) SELECT $1,id,$2 FROM roles WHERE code=$3`, userID, assignedBy, strings.TrimSpace(role))
		if err != nil {
			return err
		}
		if result.RowsAffected() != 1 {
			return domain.ErrValidation
		}
	}
	return nil
}

func recordUserAudit(ctx context.Context, tx pgx.Tx, actor, userID, event string) error {
	_, err := tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,reason) VALUES($1,$2,'user',$3,'Administrative user change')`, event, actor, userID)
	return err
}

func (r *PostgresRepository) adminUser(ctx context.Context, id string) (domain.User, error) {
	var item domain.User
	err := r.pool.QueryRow(ctx, `SELECT u.id,COALESCE(u.email,''),COALESCE(u.phone,''),COALESCE(u.display_name,''),u.status,u.created_at,COALESCE(array_agg(r.code) FILTER(WHERE r.code IS NOT NULL),'{}') FROM users u LEFT JOIN user_roles ur ON ur.user_id=u.id LEFT JOIN roles r ON r.id=ur.role_id WHERE u.id=$1 GROUP BY u.id`, id).Scan(&item.ID, &item.Email, &item.Phone, &item.DisplayName, &item.Status, &item.CreatedAt, &item.Roles)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}
	return item, err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (r *PostgresRepository) ListCategories(ctx context.Context, limit, offset int) ([]domain.Category, int, error) {
	rows, err := r.pool.Query(ctx, `SELECT c.id,c.slug,c.display_name,c.benefit_rate_basis_points,c.is_active,c.created_at,c.updated_at,COALESCE(jsonb_object_agg(ct.locale,ct.name) FILTER (WHERE ct.locale IS NOT NULL),'{}'::jsonb),count(*) OVER() FROM categories c LEFT JOIN category_translations ct ON ct.category_id=c.id GROUP BY c.id ORDER BY c.sort_order,c.slug LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list admin categories: %w", err)
	}
	defer rows.Close()
	items := []domain.Category{}
	total := 0
	for rows.Next() {
		var item domain.Category
		var translations []byte
		if err := rows.Scan(&item.ID, &item.Slug, &item.DisplayName, &item.BenefitRateBasisPoints, &item.IsActive, &item.CreatedAt, &item.UpdatedAt, &translations, &total); err != nil {
			return nil, 0, err
		}
		if err := json.Unmarshal(translations, &item.Translations); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) ListOrders(ctx context.Context, limit, offset int) ([]domain.Order, int, error) {
	rows, err := r.pool.Query(ctx, `SELECT o.id,o.order_number,COALESCE(u.display_name,''),COALESCE(u.email,''),o.status,o.currency,o.total_minor,o.created_at,count(*) OVER() FROM orders o JOIN users u ON u.id=o.user_id ORDER BY o.created_at DESC,o.id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list admin orders: %w", err)
	}
	defer rows.Close()
	items := []domain.Order{}
	total := 0
	for rows.Next() {
		var item domain.Order
		if err := rows.Scan(&item.ID, &item.OrderNumber, &item.CustomerName, &item.CustomerEmail, &item.Status, &item.Currency, &item.TotalMinor, &item.CreatedAt, &total); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) CreateCategory(ctx context.Context, actor string, input domain.CategoryInput) (domain.Category, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Category{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var id string
	if err = tx.QueryRow(ctx, `INSERT INTO categories(slug,display_name,benefit_rate_basis_points,is_active) VALUES($1,$2,$3,$4) RETURNING id`, strings.TrimSpace(input.Slug), strings.TrimSpace(input.DisplayName), input.BenefitRateBasisPoints, input.IsActive).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.Category{}, domain.ErrValidation
		}
		return domain.Category{}, err
	}
	if err = saveCategoryTranslations(ctx, tx, id, input.Translations); err != nil {
		return domain.Category{}, err
	}
	if err = recordCategoryAudit(ctx, tx, actor, id, "CATEGORY_CREATED"); err != nil {
		return domain.Category{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Category{}, err
	}
	return r.category(ctx, id)
}

func (r *PostgresRepository) UpdateCategory(ctx context.Context, actor, id string, input domain.CategoryInput) (domain.Category, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Category{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var exists bool
	if err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1)`, id).Scan(&exists); err != nil {
		return domain.Category{}, err
	}
	if !exists {
		return domain.Category{}, domain.ErrNotFound
	}
	if _, err = tx.Exec(ctx, `UPDATE categories SET slug=$2,display_name=$3,benefit_rate_basis_points=$4,is_active=$5,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id, strings.TrimSpace(input.Slug), strings.TrimSpace(input.DisplayName), input.BenefitRateBasisPoints, input.IsActive); err != nil {
		return domain.Category{}, err
	}
	if err = saveCategoryTranslations(ctx, tx, id, input.Translations); err != nil {
		return domain.Category{}, err
	}
	if err = recordCategoryAudit(ctx, tx, actor, id, "CATEGORY_UPDATED"); err != nil {
		return domain.Category{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Category{}, err
	}
	return r.category(ctx, id)
}

func (r *PostgresRepository) DeleteCategory(ctx context.Context, actor, id string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `UPDATE categories SET is_active=false,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return domain.ErrNotFound
	}
	if err = recordCategoryAudit(ctx, tx, actor, id, "CATEGORY_DEACTIVATED"); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func saveCategoryTranslations(ctx context.Context, tx pgx.Tx, id string, translations map[string]string) error {
	for _, locale := range []string{"ar", "en", "fr", "es"} {
		if _, err := tx.Exec(ctx, `INSERT INTO category_translations(category_id,locale,name) VALUES($1,$2,$3) ON CONFLICT(category_id,locale) DO UPDATE SET name=EXCLUDED.name,updated_at=CURRENT_TIMESTAMP`, id, locale, strings.TrimSpace(translations[locale])); err != nil {
			return err
		}
	}
	return nil
}

func recordCategoryAudit(ctx context.Context, tx pgx.Tx, actor, id, event string) error {
	_, err := tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,reason) VALUES($1,$2,'category',$3,'Administrative category change')`, event, actor, id)
	return err
}

func (r *PostgresRepository) category(ctx context.Context, id string) (domain.Category, error) {
	var item domain.Category
	var translations []byte
	err := r.pool.QueryRow(ctx, `SELECT c.id,c.slug,c.display_name,c.benefit_rate_basis_points,c.is_active,c.created_at,c.updated_at,COALESCE((SELECT jsonb_object_agg(ct.locale,ct.name) FROM category_translations ct WHERE ct.category_id=c.id),'{}'::jsonb) FROM categories c WHERE c.id=$1`, id).Scan(&item.ID, &item.Slug, &item.DisplayName, &item.BenefitRateBasisPoints, &item.IsActive, &item.CreatedAt, &item.UpdatedAt, &translations)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Category{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Category{}, err
	}
	if err = json.Unmarshal(translations, &item.Translations); err != nil {
		return domain.Category{}, err
	}
	return item, nil
}

func (r *PostgresRepository) SetRoles(ctx context.Context, actorID, userID string, roles []string) (domain.User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var exists bool
	if err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, userID).Scan(&exists); err != nil {
		return domain.User{}, err
	}
	if !exists {
		return domain.User{}, domain.ErrNotFound
	}
	var previous []string
	if err = tx.QueryRow(ctx, `SELECT COALESCE(array_agg(r.code ORDER BY r.code),'{}') FROM user_roles ur JOIN roles r ON r.id=ur.role_id WHERE ur.user_id=$1`, userID).Scan(&previous); err != nil {
		return domain.User{}, err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id=$1`, userID); err != nil {
		return domain.User{}, err
	}
	for _, role := range roles {
		result, insertErr := tx.Exec(ctx, `INSERT INTO user_roles(user_id,role_id,assigned_by_user_id) SELECT $1,id,$2 FROM roles WHERE code=$3`, userID, actorID, role)
		if insertErr != nil {
			return domain.User{}, insertErr
		}
		if result.RowsAffected() != 1 {
			return domain.User{}, domain.ErrValidation
		}
	}
	previousJSON, _ := json.Marshal(previous)
	nextJSON, _ := json.Marshal(roles)
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,previous_state,new_state) VALUES('USER_ROLES_CHANGED',$1,'user',$2,$3,$4)`, actorID, userID, previousJSON, nextJSON); err != nil {
		return domain.User{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES('USER_ROLES_CHANGED','user',$1,jsonb_build_object('userId',$1::text))`, userID); err != nil {
		return domain.User{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.User{}, err
	}
	return r.user(ctx, userID)
}

func (r *PostgresRepository) SetUserStatus(ctx context.Context, actorID, userID, status, reason string) (domain.User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var previous string
	if err = tx.QueryRow(ctx, `SELECT status FROM users WHERE id=$1 FOR UPDATE`, userID).Scan(&previous); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}
	if previous == status {
		return domain.User{}, domain.ErrValidation
	}
	if _, err = tx.Exec(ctx, `UPDATE users SET status=$2,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, userID, status); err != nil {
		return domain.User{}, err
	}
	if status != "ACTIVE" {
		if _, err = tx.Exec(ctx, `UPDATE sessions SET revoked_at=COALESCE(revoked_at,CURRENT_TIMESTAMP) WHERE user_id=$1`, userID); err != nil {
			return domain.User{}, err
		}
	}
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,reason,previous_state,new_state) VALUES('USER_STATUS_CHANGED',$1,'user',$2,$3,jsonb_build_object('status',$4::text),jsonb_build_object('status',$5::text))`, actorID, userID, reason, previous, status); err != nil {
		return domain.User{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES('USER_STATUS_CHANGED','user',$1,jsonb_build_object('userId',$1::uuid::text,'status',$2::text))`, userID, status); err != nil {
		return domain.User{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.User{}, err
	}
	return r.user(ctx, userID)
}

func (r *PostgresRepository) AuditEvents(ctx context.Context, filter domain.AuditFilter, limit, offset int) ([]domain.AuditEvent, int, error) {
	conditions := []string{"1=1"}
	args := []any{}
	add := func(condition string, value any) {
		args = append(args, value)
		conditions = append(conditions, fmt.Sprintf(condition, len(args)))
	}
	if filter.ActorUserID != "" {
		add("actor_user_id=$%d", filter.ActorUserID)
	}
	if filter.TargetType != "" {
		add("target_type=$%d", filter.TargetType)
	}
	if filter.TargetID != "" {
		add("target_id=$%d", filter.TargetID)
	}
	if filter.EventType != "" {
		add("event_type=$%d", filter.EventType)
	}
	if filter.CorrelationID != "" {
		add("correlation_id=$%d", filter.CorrelationID)
	}
	if filter.From != nil {
		add("occurred_at >= $%d", *filter.From)
	}
	if filter.To != nil {
		add("occurred_at < $%d", *filter.To)
	}
	args = append(args, limit, offset)
	query := `SELECT id,COALESCE(actor_user_id::text,''),target_type,COALESCE(target_id::text,''),COALESCE(correlation_id,''),COALESCE(reason,''),COALESCE(previous_state,'null'::jsonb),COALESCE(new_state,'null'::jsonb),occurred_at,count(*) OVER() FROM audit_events WHERE ` + strings.Join(conditions, " AND ") + fmt.Sprintf(" ORDER BY occurred_at DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []domain.AuditEvent{}
	total := 0
	for rows.Next() {
		var item domain.AuditEvent
		var previous, next []byte
		if err = rows.Scan(&item.ID, &item.ActorUserID, &item.TargetType, &item.TargetID, &item.CorrelationID, &item.Reason, &previous, &next, &item.OccurredAt, &total); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(previous, &item.PreviousState)
		_ = json.Unmarshal(next, &item.NewState)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) user(ctx context.Context, id string) (domain.User, error) {
	var item domain.User
	err := r.pool.QueryRow(ctx, `SELECT u.id,COALESCE(u.email,''),COALESCE(u.display_name,''),u.status,u.created_at,COALESCE(array_agg(r.code) FILTER(WHERE r.code IS NOT NULL),'{}') FROM users u LEFT JOIN user_roles ur ON ur.user_id=u.id LEFT JOIN roles r ON r.id=ur.role_id WHERE u.id=$1 GROUP BY u.id`, id).Scan(&item.ID, &item.Email, &item.DisplayName, &item.Status, &item.CreatedAt, &item.Roles)
	return item, err
}
