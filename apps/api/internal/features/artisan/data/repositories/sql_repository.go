package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApplicationInput = domain.ApplicationInput
type Application = domain.Application
type Document = domain.Document

var (
	ErrNotFound          = domain.ErrNotFound
	ErrValidation        = domain.ErrValidation
	ErrInvalidTransition = domain.ErrInvalidTransition
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const applicationSelect = `SELECT a.id,a.user_id,a.public_display_name,COALESCE(a.internal_name,''),COALESCE(a.workshop_name,''),COALESCE(a.wilaya,''),COALESCE(a.location_text,''),COALESCE(a.contact_email,''),COALESCE(a.contact_phone,''),a.contact_visibility,a.status,COALESCE(a.review_reason,''),COALESCE((SELECT json_agg(category_id::text) FROM artisan_profile_categories WHERE artisan_profile_id=a.id),'[]'),COALESCE((SELECT json_agg(json_build_object('Locale',locale,'Biography',COALESCE(biography,'')) ORDER BY locale) FROM artisan_profile_translations WHERE artisan_profile_id=a.id),'[]') FROM artisan_profiles a `

type rowScanner interface{ Scan(...any) error }

func scanApplication(row rowScanner) (Application, error) {
	var a Application
	var categories, translations []byte
	err := row.Scan(&a.ID, &a.UserID, &a.PublicDisplayName, &a.InternalName, &a.WorkshopName, &a.Wilaya, &a.Location, &a.ContactEmail, &a.ContactPhone, &a.ContactVisibility, &a.Status, &a.ReviewReason, &categories, &translations)
	if err != nil {
		return Application{}, err
	}
	if err = json.Unmarshal(categories, &a.CategoryIDs); err != nil {
		return Application{}, fmt.Errorf("decode artisan categories: %w", err)
	}
	if err = json.Unmarshal(translations, &a.Translations); err != nil {
		return Application{}, fmt.Errorf("decode artisan translations: %w", err)
	}
	return a, nil
}

func (r *PostgresRepository) Submit(ctx context.Context, userID string, input ApplicationInput) (Application, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Application{}, fmt.Errorf("begin artisan submission: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var id, status string
	err = tx.QueryRow(ctx, `SELECT id,status FROM artisan_profiles WHERE user_id=$1 FOR UPDATE`, userID).Scan(&id, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		err = tx.QueryRow(ctx, `INSERT INTO artisan_profiles(user_id,public_display_name,internal_name,workshop_name,wilaya,location_text,contact_email,contact_phone,contact_visibility,status,submitted_at) VALUES($1,$2,NULLIF($3,''),NULLIF($4,''),$5,NULLIF($6,''),NULLIF($7,''),NULLIF($8,''),$9,'SUBMITTED',CURRENT_TIMESTAMP) RETURNING id,status`, userID, input.PublicDisplayName, input.InternalName, input.WorkshopName, input.Wilaya, input.Location, input.ContactEmail, input.ContactPhone, input.ContactVisibility).Scan(&id, &status)
	} else if err == nil {
		if status != "DRAFT" && status != "CHANGES_REQUESTED" && status != "REJECTED" {
			return Application{}, ErrInvalidTransition
		}
		_, err = tx.Exec(ctx, `UPDATE artisan_profiles SET public_display_name=$2,internal_name=NULLIF($3,''),workshop_name=NULLIF($4,''),wilaya=$5,location_text=NULLIF($6,''),contact_email=NULLIF($7,''),contact_phone=NULLIF($8,''),contact_visibility=$9,status='SUBMITTED',review_reason=NULL,submitted_at=CURRENT_TIMESTAMP,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id, input.PublicDisplayName, input.InternalName, input.WorkshopName, input.Wilaya, input.Location, input.ContactEmail, input.ContactPhone, input.ContactVisibility)
	}
	if err != nil {
		return Application{}, fmt.Errorf("persist artisan submission: %w", err)
	}
	if err = replaceDetails(ctx, tx, id, input); err != nil {
		return Application{}, err
	}
	if err = event(ctx, tx, "ARTISAN_APPLICATION_SUBMITTED", userID, id, "", status, "SUBMITTED"); err != nil {
		return Application{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Application{}, fmt.Errorf("commit artisan submission: %w", err)
	}
	return r.Mine(ctx, userID)
}

func replaceDetails(ctx context.Context, tx pgx.Tx, id string, input ApplicationInput) error {
	if _, err := tx.Exec(ctx, `DELETE FROM artisan_profile_categories WHERE artisan_profile_id=$1`, id); err != nil {
		return err
	}
	for _, categoryID := range input.CategoryIDs {
		tag, err := tx.Exec(ctx, `INSERT INTO artisan_profile_categories(artisan_profile_id,category_id) SELECT $1,id FROM categories WHERE id=$2 AND is_active=true`, id, categoryID)
		if err != nil || tag.RowsAffected() != 1 {
			return ErrValidation
		}
	}
	if _, err := tx.Exec(ctx, `DELETE FROM artisan_profile_translations WHERE artisan_profile_id=$1`, id); err != nil {
		return err
	}
	for _, t := range input.Translations {
		if _, err := tx.Exec(ctx, `INSERT INTO artisan_profile_translations(artisan_profile_id,locale,biography) VALUES($1,$2,NULLIF($3,''))`, id, t.Locale, t.Biography); err != nil {
			return err
		}
	}
	return nil
}
func event(ctx context.Context, tx pgx.Tx, eventType, actor, id, reason, previous, next string) error {
	if _, err := tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,reason,previous_state,new_state) VALUES($1,$2,'artisan_profile',$3,NULLIF($4,''),jsonb_build_object('status',$5::text),jsonb_build_object('status',$6::text))`, eventType, actor, id, reason, previous, next); err != nil {
		return fmt.Errorf("insert artisan audit event: %w", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,'artisan_profile',$2::uuid,jsonb_build_object('artisanProfileId',($2::uuid)::text,'status',$3::text))`, eventType, id, next); err != nil {
		return fmt.Errorf("insert artisan outbox event: %w", err)
	}
	return nil
}
func (r *PostgresRepository) Mine(ctx context.Context, userID string) (Application, error) {
	a, err := scanApplication(r.pool.QueryRow(ctx, applicationSelect+`WHERE a.user_id=$1`, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return Application{}, ErrNotFound
	}
	return a, err
}
func (r *PostgresRepository) UpdateApproved(ctx context.Context, userID string, input ApplicationInput) (Application, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Application{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var id string
	err = tx.QueryRow(ctx, `UPDATE artisan_profiles SET public_display_name=$2,internal_name=NULLIF($3,''),workshop_name=NULLIF($4,''),wilaya=$5,location_text=NULLIF($6,''),contact_email=NULLIF($7,''),contact_phone=NULLIF($8,''),contact_visibility=$9,updated_at=CURRENT_TIMESTAMP WHERE user_id=$1 AND status='APPROVED' RETURNING id`, userID, input.PublicDisplayName, input.InternalName, input.WorkshopName, input.Wilaya, input.Location, input.ContactEmail, input.ContactPhone, input.ContactVisibility).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Application{}, ErrInvalidTransition
	}
	if err != nil {
		return Application{}, err
	}
	if err = replaceDetails(ctx, tx, id, input); err != nil {
		return Application{}, err
	}
	if err = event(ctx, tx, "ARTISAN_PROFILE_UPDATED", userID, id, "", "APPROVED", "APPROVED"); err != nil {
		return Application{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Application{}, err
	}
	return r.Mine(ctx, userID)
}
func (r *PostgresRepository) List(ctx context.Context, status string, limit, offset int) ([]Application, int, error) {
	rows, err := r.pool.Query(ctx, applicationSelect+`WHERE ($1='' OR a.status=$1) ORDER BY a.submitted_at DESC NULLS LAST,a.created_at DESC LIMIT $2 OFFSET $3`, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []Application{}
	for rows.Next() {
		a, err := scanApplication(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, a)
	}
	var total int
	err = r.pool.QueryRow(ctx, `SELECT count(*) FROM artisan_profiles WHERE ($1='' OR status=$1)`, status).Scan(&total)
	return items, total, err
}
func (r *PostgresRepository) Decide(ctx context.Context, actor, id, decision, reason string) (Application, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Application{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var previous, userID string
	err = tx.QueryRow(ctx, `SELECT status,user_id FROM artisan_profiles WHERE id=$1 FOR UPDATE`, id).Scan(&previous, &userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return Application{}, ErrNotFound
	}
	if err != nil {
		return Application{}, err
	}
	if previous != "SUBMITTED" && previous != "UNDER_REVIEW" {
		return Application{}, ErrInvalidTransition
	}
	_, err = tx.Exec(ctx, `UPDATE artisan_profiles SET status=$2,review_reason=NULLIF($3,''),approved_at=CASE WHEN $2='APPROVED' THEN CURRENT_TIMESTAMP ELSE approved_at END,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id, decision, reason)
	if err != nil {
		return Application{}, err
	}
	if decision == "APPROVED" {
		_, err = tx.Exec(ctx, `INSERT INTO user_roles(user_id,role_id,assigned_by_user_id) SELECT $1,id,$2 FROM roles WHERE code='artisan' ON CONFLICT DO NOTHING`, userID, actor)
		if err != nil {
			return Application{}, err
		}
	}
	if err = event(ctx, tx, "ARTISAN_APPLICATION_"+decision, actor, id, reason, previous, decision); err != nil {
		return Application{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Application{}, err
	}
	return scanApplication(r.pool.QueryRow(ctx, applicationSelect+`WHERE a.id=$1`, id))
}
func (r *PostgresRepository) Documents(ctx context.Context, id string) ([]Document, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,document_type,COALESCE(original_filename,''),media_type,size_bytes FROM artisan_documents WHERE artisan_profile_id=$1 ORDER BY created_at`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Document{}
	for rows.Next() {
		var d Document
		if err := rows.Scan(&d.ID, &d.DocumentType, &d.OriginalFilename, &d.MediaType, &d.SizeBytes); err != nil {
			return nil, err
		}
		items = append(items, d)
	}
	return items, rows.Err()
}
