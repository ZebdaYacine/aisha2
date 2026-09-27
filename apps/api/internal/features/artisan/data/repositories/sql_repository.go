package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

const applicationSelect = `SELECT a.id,a.user_id,a.public_display_name,COALESCE(a.internal_name,''),COALESCE(a.workshop_name,''),COALESCE(a.wilaya,''),COALESCE(a.location_text,''),COALESCE(a.contact_email,''),COALESCE(a.contact_phone,''),a.contact_visibility,a.status,COALESCE(a.review_reason,''),COALESCE((SELECT status FROM artisan_memberships m WHERE m.artisan_profile_id=a.id),'NOT_STARTED'),COALESCE((SELECT json_agg(category_id::text) FROM artisan_profile_categories WHERE artisan_profile_id=a.id),'[]'),COALESCE((SELECT json_agg(json_build_object('Locale',locale,'Biography',COALESCE(biography,'')) ORDER BY locale) FROM artisan_profile_translations WHERE artisan_profile_id=a.id),'[]') FROM artisan_profiles a `

type rowScanner interface{ Scan(...any) error }

func scanApplication(row rowScanner) (Application, error) {
	var a Application
	var categories, translations []byte
	err := row.Scan(&a.ID, &a.UserID, &a.PublicDisplayName, &a.InternalName, &a.WorkshopName, &a.Wilaya, &a.Location, &a.ContactEmail, &a.ContactPhone, &a.ContactVisibility, &a.Status, &a.ReviewReason, &a.MembershipStatus, &categories, &translations)
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
	if _, err := r.SaveDraft(ctx, userID, input); err != nil {
		return Application{}, err
	}
	return r.FinalizeSubmission(ctx, userID)
}

func (r *PostgresRepository) SaveDraft(ctx context.Context, userID string, input ApplicationInput) (Application, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Application{}, fmt.Errorf("begin artisan draft: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var id, status string
	err = tx.QueryRow(ctx, `SELECT id,status FROM artisan_profiles WHERE user_id=$1 FOR UPDATE`, userID).Scan(&id, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		err = tx.QueryRow(ctx, `INSERT INTO artisan_profiles(user_id,public_display_name,internal_name,workshop_name,wilaya,location_text,contact_email,contact_phone,contact_visibility,status) VALUES($1,$2,NULLIF($3,''),NULLIF($4,''),$5,NULLIF($6,''),NULLIF($7,''),NULLIF($8,''),$9,'DRAFT') RETURNING id,status`, userID, input.PublicDisplayName, input.InternalName, input.WorkshopName, input.Wilaya, input.Location, input.ContactEmail, input.ContactPhone, input.ContactVisibility).Scan(&id, &status)
	} else if err == nil {
		if status != "DRAFT" && status != "CHANGES_REQUESTED" && status != "REJECTED" {
			return Application{}, ErrInvalidTransition
		}
		_, err = tx.Exec(ctx, `UPDATE artisan_profiles SET public_display_name=$2,internal_name=NULLIF($3,''),workshop_name=NULLIF($4,''),wilaya=$5,location_text=NULLIF($6,''),contact_email=NULLIF($7,''),contact_phone=NULLIF($8,''),contact_visibility=$9,status='DRAFT',review_reason=NULL,submitted_at=NULL,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id, input.PublicDisplayName, input.InternalName, input.WorkshopName, input.Wilaya, input.Location, input.ContactEmail, input.ContactPhone, input.ContactVisibility)
	}
	if err != nil {
		return Application{}, fmt.Errorf("persist artisan draft: %w", err)
	}
	if err = replaceDetails(ctx, tx, id, input); err != nil {
		return Application{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Application{}, fmt.Errorf("commit artisan draft: %w", err)
	}
	return r.Mine(ctx, userID)
}

func (r *PostgresRepository) FinalizeSubmission(ctx context.Context, userID string) (Application, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Application{}, fmt.Errorf("begin artisan final submission: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var id, status, workshopName string
	err = tx.QueryRow(ctx, `SELECT id,status,COALESCE(workshop_name,'') FROM artisan_profiles WHERE user_id=$1 FOR UPDATE`, userID).Scan(&id, &status, &workshopName)
	if errors.Is(err, pgx.ErrNoRows) {
		return Application{}, ErrNotFound
	}
	if err != nil {
		return Application{}, err
	}
	if status != "DRAFT" && status != "CHANGES_REQUESTED" && status != "REJECTED" {
		return Application{}, ErrInvalidTransition
	}
	if strings.TrimSpace(workshopName) == "" {
		return Application{}, ErrValidation
	}
	var documents, media int
	if err = tx.QueryRow(ctx, `SELECT count(*) FROM artisan_documents WHERE artisan_profile_id=$1`, id).Scan(&documents); err != nil {
		return Application{}, err
	}
	if err = tx.QueryRow(ctx, `SELECT count(*) FROM artisan_media WHERE artisan_profile_id=$1`, id).Scan(&media); err != nil {
		return Application{}, err
	}
	if documents == 0 || media == 0 {
		return Application{}, ErrValidation
	}
	if _, err = tx.Exec(ctx, `UPDATE artisan_profiles SET status='SUBMITTED',review_reason=NULL,submitted_at=CURRENT_TIMESTAMP,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id); err != nil {
		return Application{}, err
	}
	if err = event(ctx, tx, "ARTISAN_APPLICATION_SUBMITTED", userID, id, "", status, "SUBMITTED"); err != nil {
		return Application{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Application{}, fmt.Errorf("commit artisan final submission: %w", err)
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
	return eventTarget(ctx, tx, eventType, actor, "artisan_profile", "artisan_profile", id, reason, previous, next)
}

func eventTarget(ctx context.Context, tx pgx.Tx, eventType, actor, targetType, aggregateType, id, reason, previous, next string) error {
	if _, err := tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,reason,previous_state,new_state) VALUES($1,$2,$3,$4,NULLIF($5,''),jsonb_build_object('status',$6::text),jsonb_build_object('status',$7::text))`, eventType, actor, targetType, id, reason, previous, next); err != nil {
		return fmt.Errorf("insert artisan audit event: %w", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,$2,$3::uuid,jsonb_build_object('id',($3::uuid)::text,'status',$4::text))`, eventType, aggregateType, id, next); err != nil {
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
		var membershipID string
		if err = tx.QueryRow(ctx, `INSERT INTO artisan_memberships(user_id,artisan_profile_id,status,activated_at,reason,suspended_at,closed_at) VALUES($1,$2,'ACTIVE',CURRENT_TIMESTAMP,NULL,NULL,NULL) ON CONFLICT (artisan_profile_id) DO UPDATE SET status='ACTIVE',activated_at=COALESCE(artisan_memberships.activated_at,CURRENT_TIMESTAMP),reason=NULL,suspended_at=NULL,closed_at=NULL,updated_at=CURRENT_TIMESTAMP RETURNING id`, userID, id).Scan(&membershipID); err != nil {
			return Application{}, err
		}
		if err = ensureDefaultWorkshop(ctx, tx, id); err != nil {
			return Application{}, err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO artisan_verifications(artisan_membership_id) VALUES($1) ON CONFLICT (artisan_membership_id) DO NOTHING`, membershipID); err != nil {
			return Application{}, err
		}
		if err = eventTarget(ctx, tx, "ARTISAN_MEMBERSHIP_ACTIVE", actor, "artisan_membership", "artisan_membership", membershipID, "Approved artisan application", "NOT_STARTED", "ACTIVE"); err != nil {
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
	rows, err := r.pool.Query(ctx, `SELECT id,document_type,object_key,COALESCE(original_filename,''),media_type,size_bytes,created_at FROM artisan_documents WHERE artisan_profile_id=$1 ORDER BY created_at`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Document{}
	for rows.Next() {
		var d Document
		if err := rows.Scan(&d.ID, &d.DocumentType, &d.ObjectKey, &d.OriginalFilename, &d.MediaType, &d.SizeBytes, &d.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, d)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) Media(ctx context.Context, id string) ([]domain.Media, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,media_kind,object_key,COALESCE(original_filename,''),media_type,size_bytes,sort_order,visibility,created_at FROM artisan_media WHERE artisan_profile_id=$1 ORDER BY sort_order,created_at`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.Media{}
	for rows.Next() {
		var item domain.Media
		if err = rows.Scan(&item.ID, &item.MediaKind, &item.ObjectKey, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.SortOrder, &item.Visibility, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) Profile(ctx context.Context, userID string) (string, string, error) {
	var id, status string
	err := r.pool.QueryRow(ctx, `SELECT id,status FROM artisan_profiles WHERE user_id=$1`, userID).Scan(&id, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", domain.ErrNotFound
	}
	return id, status, err
}

func (r *PostgresRepository) AddDocument(ctx context.Context, userID string, input domain.DocumentUploadInput) (Document, error) {
	var item Document
	err := r.pool.QueryRow(ctx, `INSERT INTO artisan_documents(artisan_profile_id,document_type,object_key,original_filename,media_type,size_bytes,checksum_sha256) SELECT id,$2,$3,NULLIF($4,''),$5,$6,$7 FROM artisan_profiles WHERE user_id=$1 AND status <> 'SUSPENDED' RETURNING id,document_type,object_key,COALESCE(original_filename,''),media_type,size_bytes,created_at`, userID, input.DocumentType, input.ObjectKey, input.OriginalFilename, input.MediaType, input.SizeBytes, input.Checksum).Scan(&item.ID, &item.DocumentType, &item.ObjectKey, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Document{}, domain.ErrInvalidTransition
	}
	return item, err
}

func (r *PostgresRepository) OwnDocuments(ctx context.Context, userID string) ([]Document, error) {
	rows, err := r.pool.Query(ctx, `SELECT d.id,d.document_type,d.object_key,COALESCE(d.original_filename,''),d.media_type,d.size_bytes,d.created_at FROM artisan_documents d JOIN artisan_profiles a ON a.id=d.artisan_profile_id WHERE a.user_id=$1 ORDER BY d.created_at`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Document{}
	for rows.Next() {
		var item Document
		if err = rows.Scan(&item.ID, &item.DocumentType, &item.ObjectKey, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) AddMedia(ctx context.Context, userID string, input domain.MediaUploadInput) (domain.Media, error) {
	var item domain.Media
	err := r.pool.QueryRow(ctx, `INSERT INTO artisan_media(artisan_profile_id,media_kind,object_key,original_filename,media_type,size_bytes,checksum_sha256,sort_order,visibility) SELECT id,$2,$3,NULLIF($4,''),$5,$6,$7,$8,'PRIVATE' FROM artisan_profiles WHERE user_id=$1 AND status <> 'SUSPENDED' RETURNING id,media_kind,object_key,COALESCE(original_filename,''),media_type,size_bytes,sort_order,visibility,created_at`, userID, input.MediaKind, input.ObjectKey, input.OriginalFilename, input.MediaType, input.SizeBytes, input.Checksum, input.SortOrder).Scan(&item.ID, &item.MediaKind, &item.ObjectKey, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.SortOrder, &item.Visibility, &item.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Media{}, domain.ErrInvalidTransition
	}
	return item, err
}

func (r *PostgresRepository) OwnMedia(ctx context.Context, userID string) ([]domain.Media, error) {
	rows, err := r.pool.Query(ctx, `SELECT m.id,m.media_kind,m.object_key,COALESCE(m.original_filename,''),m.media_type,m.size_bytes,m.sort_order,m.visibility,m.created_at FROM artisan_media m JOIN artisan_profiles a ON a.id=m.artisan_profile_id WHERE a.user_id=$1 ORDER BY m.sort_order,m.created_at`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.Media{}
	for rows.Next() {
		var item domain.Media
		if err = rows.Scan(&item.ID, &item.MediaKind, &item.ObjectKey, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.SortOrder, &item.Visibility, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) ReplaceMedia(ctx context.Context, userID, mediaID string, input domain.MediaUploadInput) (domain.Media, string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Media{}, "", err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var oldKey string
	if err = tx.QueryRow(ctx, `SELECT m.object_key FROM artisan_media m JOIN artisan_profiles a ON a.id=m.artisan_profile_id WHERE m.id=$1 AND a.user_id=$2 AND a.status <> 'SUSPENDED' FOR UPDATE`, mediaID, userID).Scan(&oldKey); errors.Is(err, pgx.ErrNoRows) {
		return domain.Media{}, "", domain.ErrNotFound
	} else if err != nil {
		return domain.Media{}, "", err
	}
	var item domain.Media
	err = tx.QueryRow(ctx, `UPDATE artisan_media SET media_kind=$3,object_key=$4,original_filename=NULLIF($5,''),media_type=$6,size_bytes=$7,checksum_sha256=$8 WHERE id=$1 AND artisan_profile_id=(SELECT id FROM artisan_profiles WHERE user_id=$2) RETURNING id,media_kind,object_key,COALESCE(original_filename,''),media_type,size_bytes,sort_order,visibility,created_at`, mediaID, userID, input.MediaKind, input.ObjectKey, input.OriginalFilename, input.MediaType, input.SizeBytes, input.Checksum).Scan(&item.ID, &item.MediaKind, &item.ObjectKey, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.SortOrder, &item.Visibility, &item.CreatedAt)
	if err != nil {
		return domain.Media{}, "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Media{}, "", err
	}
	return item, oldKey, nil
}

func (r *PostgresRepository) DeleteMedia(ctx context.Context, userID, mediaID string) (string, error) {
	var objectKey string
	err := r.pool.QueryRow(ctx, `DELETE FROM artisan_media m USING artisan_profiles a WHERE m.id=$1 AND m.artisan_profile_id=a.id AND a.user_id=$2 RETURNING m.object_key`, mediaID, userID).Scan(&objectKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrNotFound
	}
	return objectKey, err
}

func (r *PostgresRepository) ActivateMembership(ctx context.Context, userID string, input domain.WorkshopInput, key, requestHash string) (Application, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Application{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var idemID, state, resourceID string
	err = tx.QueryRow(ctx, `INSERT INTO idempotency_keys(actor_user_id,scope,idempotency_key,request_hash,expires_at) VALUES($1,'artisan-membership-activation',$2,$3,CURRENT_TIMESTAMP+interval '1 day') ON CONFLICT(actor_user_id,scope,idempotency_key) DO NOTHING RETURNING id,state,COALESCE(resource_id::text,'')`, userID, key, requestHash).Scan(&idemID, &state, &resourceID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return Application{}, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		var storedHash string
		err = tx.QueryRow(ctx, `SELECT state,request_hash,COALESCE(resource_id::text,'') FROM idempotency_keys WHERE actor_user_id=$1 AND scope='artisan-membership-activation' AND idempotency_key=$2 FOR UPDATE`, userID, key).Scan(&state, &storedHash, &resourceID)
		if err != nil {
			return Application{}, err
		}
		if strings.TrimSpace(storedHash) != strings.TrimSpace(requestHash) {
			return Application{}, domain.ErrInvalidTransition
		}
		if state == "COMPLETED" && resourceID != "" {
			_ = tx.Rollback(ctx)
			return r.Mine(ctx, userID)
		}
		return Application{}, domain.ErrInvalidTransition
	}
	var profileID, profileStatus string
	if err = tx.QueryRow(ctx, `SELECT id,status FROM artisan_profiles WHERE user_id=$1 FOR UPDATE`, userID).Scan(&profileID, &profileStatus); errors.Is(err, pgx.ErrNoRows) {
		return Application{}, domain.ErrNotFound
	}
	if err != nil {
		return Application{}, err
	}
	if profileStatus != "APPROVED" {
		return Application{}, domain.ErrInvalidTransition
	}
	var membershipID string
	if err = tx.QueryRow(ctx, `SELECT id FROM artisan_memberships WHERE artisan_profile_id=$1 FOR UPDATE`, profileID).Scan(&membershipID); err == nil {
		return Application{}, domain.ErrInvalidTransition
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return Application{}, err
	}
	if err = tx.QueryRow(ctx, `INSERT INTO artisan_memberships(user_id,artisan_profile_id,status,activated_at) VALUES($1,$2,'ACTIVE',CURRENT_TIMESTAMP) RETURNING id`, userID, profileID).Scan(&membershipID); err != nil {
		return Application{}, err
	}
	var workshopID string
	if err = tx.QueryRow(ctx, `INSERT INTO workshops(artisan_profile_id,name,description,wilaya,location_text,status,is_default,is_public) VALUES($1,$2,NULLIF($3,''),$4,NULLIF($5,''),'ACTIVE',true,$6) RETURNING id`, profileID, input.Name, input.Description, input.Wilaya, input.Location, input.IsPublic).Scan(&workshopID); err != nil {
		return Application{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO artisan_verifications(artisan_membership_id) VALUES($1)`, membershipID); err != nil {
		return Application{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO user_roles(user_id,role_id) SELECT $1,id FROM roles WHERE code='artisan' ON CONFLICT DO NOTHING`, userID); err != nil {
		return Application{}, err
	}
	if err = event(ctx, tx, "ARTISAN_MEMBERSHIP_ACTIVATED", userID, profileID, "", "APPROVED", "ACTIVE"); err != nil {
		return Application{}, err
	}
	payload, _ := json.Marshal(map[string]any{"artisanProfileId": profileID, "membershipId": membershipID, "workshopId": workshopID, "status": "ACTIVE"})
	if _, err = tx.Exec(ctx, `UPDATE idempotency_keys SET state='COMPLETED',response_status=201,response_body=$2,resource_type='artisan_membership',resource_id=$3,completed_at=CURRENT_TIMESTAMP WHERE id=$1`, idemID, payload, membershipID); err != nil {
		return Application{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Application{}, err
	}
	return r.Mine(ctx, userID)
}
func (r *PostgresRepository) ListWorkshops(ctx context.Context, userID string) ([]domain.Workshop, error) {
	rows, err := r.pool.Query(ctx, `SELECT w.id,w.name,COALESCE(w.description,''),COALESCE(w.wilaya,''),COALESCE(w.location_text,''),w.status,w.is_default,w.is_public,COUNT(p.id)::int FROM workshops w JOIN artisan_profiles a ON a.id=w.artisan_profile_id JOIN artisan_memberships m ON m.artisan_profile_id=a.id AND m.user_id=$1 LEFT JOIN products p ON p.workshop_id=w.id GROUP BY w.id ORDER BY w.is_default DESC,w.name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Workshop{}
	for rows.Next() {
		var v domain.Workshop
		if err = rows.Scan(&v.ID, &v.Name, &v.Description, &v.Wilaya, &v.Location, &v.Status, &v.IsDefault, &v.IsPublic, &v.ProductCount); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (r *PostgresRepository) CreateWorkshop(ctx context.Context, userID string, input domain.WorkshopInput, key, requestHash string) (domain.Workshop, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Workshop{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var idemID, state, resourceID, storedHash string
	if err = tx.QueryRow(ctx, `INSERT INTO idempotency_keys(actor_user_id,scope,idempotency_key,request_hash,expires_at) VALUES($1,'workshop-create',$2,$3,CURRENT_TIMESTAMP+interval '1 day') ON CONFLICT(actor_user_id,scope,idempotency_key) DO NOTHING RETURNING id`, userID, key, requestHash).Scan(&idemID); errors.Is(err, pgx.ErrNoRows) {
		if err = tx.QueryRow(ctx, `SELECT state,request_hash,COALESCE(resource_id::text,'') FROM idempotency_keys WHERE actor_user_id=$1 AND scope='workshop-create' AND idempotency_key=$2 FOR UPDATE`, userID, key).Scan(&state, &storedHash, &resourceID); err != nil {
			return domain.Workshop{}, err
		}
		if strings.TrimSpace(storedHash) != strings.TrimSpace(requestHash) {
			return domain.Workshop{}, domain.ErrInvalidTransition
		}
		if state == "COMPLETED" && resourceID != "" {
			var v domain.Workshop
			err = tx.QueryRow(ctx, `SELECT w.id,w.name,COALESCE(w.description,''),COALESCE(w.wilaya,''),COALESCE(w.location_text,''),w.status,w.is_default,w.is_public,(SELECT count(*)::int FROM products p WHERE p.workshop_id=w.id) FROM workshops w JOIN artisan_profiles a ON a.id=w.artisan_profile_id JOIN artisan_memberships m ON m.artisan_profile_id=a.id AND m.user_id=$1 WHERE w.id=$2`, userID, resourceID).Scan(&v.ID, &v.Name, &v.Description, &v.Wilaya, &v.Location, &v.Status, &v.IsDefault, &v.IsPublic, &v.ProductCount)
			if err != nil {
				return domain.Workshop{}, err
			}
			_ = tx.Rollback(ctx)
			return v, nil
		}
		return domain.Workshop{}, domain.ErrInvalidTransition
	}
	if err != nil {
		return domain.Workshop{}, err
	}
	var membershipStatus string
	if err = tx.QueryRow(ctx, `SELECT m.status FROM artisan_memberships m WHERE m.user_id=$1 FOR UPDATE`, userID).Scan(&membershipStatus); errors.Is(err, pgx.ErrNoRows) {
		return domain.Workshop{}, domain.ErrInvalidTransition
	}
	if err != nil {
		return domain.Workshop{}, err
	}
	if membershipStatus != "ACTIVE" {
		return domain.Workshop{}, domain.ErrInvalidTransition
	}
	var v domain.Workshop
	err = tx.QueryRow(ctx, `INSERT INTO workshops(artisan_profile_id,name,description,wilaya,location_text,status,is_default,is_public) SELECT artisan_profile_id,$2,NULLIF($3,''),$4,NULLIF($5,''),'ACTIVE',false,$6 FROM artisan_memberships WHERE user_id=$1 RETURNING id,name,COALESCE(description,''),COALESCE(wilaya,''),COALESCE(location_text,''),status,is_default,is_public`, userID, input.Name, input.Description, input.Wilaya, input.Location, input.IsPublic).Scan(&v.ID, &v.Name, &v.Description, &v.Wilaya, &v.Location, &v.Status, &v.IsDefault, &v.IsPublic)
	if err != nil {
		return domain.Workshop{}, err
	}
	if err = eventTarget(ctx, tx, "WORKSHOP_CREATED", userID, "workshop", "workshop", v.ID, "", "", "ACTIVE"); err != nil {
		return domain.Workshop{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE idempotency_keys SET state='COMPLETED',response_status=201,resource_type='workshop',resource_id=$2,completed_at=CURRENT_TIMESTAMP WHERE id=$1`, idemID, v.ID); err != nil {
		return domain.Workshop{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Workshop{}, err
	}
	return v, nil
}
func (r *PostgresRepository) UpdateWorkshop(ctx context.Context, userID, id string, input domain.WorkshopInput) (domain.Workshop, error) {
	var v domain.Workshop
	err := r.pool.QueryRow(ctx, `UPDATE workshops w SET name=$3,description=NULLIF($4,''),wilaya=$5,location_text=NULLIF($6,''),is_public=$7,updated_at=CURRENT_TIMESTAMP FROM artisan_profiles a JOIN artisan_memberships m ON m.artisan_profile_id=a.id AND m.user_id=$1 AND m.status='ACTIVE' WHERE w.id=$2 AND w.artisan_profile_id=a.id RETURNING w.id,w.name,COALESCE(w.description,''),COALESCE(w.wilaya,''),COALESCE(w.location_text,''),w.status,w.is_default,w.is_public,(SELECT count(*)::int FROM products p WHERE p.workshop_id=w.id)`, userID, id, input.Name, input.Description, input.Wilaya, input.Location, input.IsPublic).Scan(&v.ID, &v.Name, &v.Description, &v.Wilaya, &v.Location, &v.Status, &v.IsDefault, &v.IsPublic, &v.ProductCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Workshop{}, domain.ErrNotFound
	}
	return v, err
}
func (r *PostgresRepository) SetWorkshopStatus(ctx context.Context, userID, id, status string) (domain.Workshop, error) {
	var v domain.Workshop
	err := r.pool.QueryRow(ctx, `UPDATE workshops w SET status=$3,is_public=CASE WHEN $3='ACTIVE' THEN is_public ELSE false END,updated_at=CURRENT_TIMESTAMP FROM artisan_profiles a JOIN artisan_memberships m ON m.artisan_profile_id=a.id AND m.user_id=$1 AND m.status='ACTIVE' WHERE w.id=$2 AND w.artisan_profile_id=a.id RETURNING w.id,w.name,COALESCE(w.description,''),COALESCE(w.wilaya,''),COALESCE(w.location_text,''),w.status,w.is_default,w.is_public,(SELECT count(*)::int FROM products p WHERE p.workshop_id=w.id)`, userID, id, status).Scan(&v.ID, &v.Name, &v.Description, &v.Wilaya, &v.Location, &v.Status, &v.IsDefault, &v.IsPublic, &v.ProductCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Workshop{}, domain.ErrNotFound
	}
	return v, err
}
func (r *PostgresRepository) SetWorkshopStatusAdmin(ctx context.Context, actor, id, status, reason string) (domain.Workshop, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Workshop{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var previous string
	if err = tx.QueryRow(ctx, `SELECT status FROM workshops WHERE id=$1 FOR UPDATE`, id).Scan(&previous); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Workshop{}, domain.ErrNotFound
		}
		return domain.Workshop{}, err
	}
	if previous == status {
		return domain.Workshop{}, domain.ErrInvalidTransition
	}
	var v domain.Workshop
	if err = tx.QueryRow(ctx, `UPDATE workshops SET status=$2,is_public=CASE WHEN $2='ACTIVE' THEN is_public ELSE false END,updated_at=CURRENT_TIMESTAMP WHERE id=$1 RETURNING id,name,COALESCE(description,''),COALESCE(wilaya,''),COALESCE(location_text,''),status,is_default,is_public,(SELECT count(*)::int FROM products p WHERE p.workshop_id=workshops.id)`, id, status).Scan(&v.ID, &v.Name, &v.Description, &v.Wilaya, &v.Location, &v.Status, &v.IsDefault, &v.IsPublic, &v.ProductCount); err != nil {
		return domain.Workshop{}, err
	}
	if err = eventTarget(ctx, tx, "WORKSHOP_ADMIN_STATUS_CHANGED", actor, "workshop", "workshop", id, reason, previous, status); err != nil {
		return domain.Workshop{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Workshop{}, err
	}
	return v, nil
}
func (r *PostgresRepository) DeleteWorkshop(ctx context.Context, userID, id string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var profileID string
	var def bool
	if err = tx.QueryRow(ctx, `SELECT w.artisan_profile_id,w.is_default FROM workshops w JOIN artisan_profiles a ON a.id=w.artisan_profile_id JOIN artisan_memberships m ON m.artisan_profile_id=a.id AND m.user_id=$1 AND m.status='ACTIVE' WHERE w.id=$2 FOR UPDATE`, userID, id).Scan(&profileID, &def); errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	if err != nil {
		return err
	}
	if def {
		return domain.ErrInvalidTransition
	}
	var protected int
	if err = tx.QueryRow(ctx, `SELECT (SELECT count(*) FROM products WHERE workshop_id=$1)+(SELECT count(*) FROM inventory_movements i JOIN products p ON p.id=i.product_id WHERE p.workshop_id=$1)+(SELECT count(*) FROM order_items oi WHERE oi.workshop_id=$1)`, id).Scan(&protected); err != nil {
		return err
	}
	if protected > 0 {
		return domain.ErrInvalidTransition
	}
	if _, err = tx.Exec(ctx, `DELETE FROM workshops WHERE id=$1`, id); err != nil {
		return err
	}
	if err = eventTarget(ctx, tx, "WORKSHOP_DELETED", userID, "workshop", "workshop", id, "", "ACTIVE", "DELETED"); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *PostgresRepository) SetMembershipStatus(ctx context.Context, actor, id, status, reason string) (Application, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Application{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var membershipID, profileID, previous string
	if err = tx.QueryRow(ctx, `SELECT m.id,m.artisan_profile_id,m.status FROM artisan_memberships m WHERE m.id=$1 OR m.artisan_profile_id=$1 FOR UPDATE`, id).Scan(&membershipID, &profileID, &previous); errors.Is(err, pgx.ErrNoRows) {
		if status != "ACTIVE" {
			return Application{}, domain.ErrNotFound
		}
		var userID string
		if err = tx.QueryRow(ctx, `SELECT id,user_id FROM artisan_profiles WHERE id=$1 AND status='APPROVED' FOR UPDATE`, id).Scan(&profileID, &userID); errors.Is(err, pgx.ErrNoRows) {
			return Application{}, domain.ErrNotFound
		} else if err != nil {
			return Application{}, err
		}
		if err = tx.QueryRow(ctx, `INSERT INTO artisan_memberships(user_id,artisan_profile_id,status,activated_at) VALUES($1,$2,'ACTIVE',CURRENT_TIMESTAMP) RETURNING id`, userID, profileID).Scan(&membershipID); err != nil {
			return Application{}, err
		}
		if err = ensureDefaultWorkshop(ctx, tx, profileID); err != nil {
			return Application{}, err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO artisan_verifications(artisan_membership_id) VALUES($1) ON CONFLICT (artisan_membership_id) DO NOTHING`, membershipID); err != nil {
			return Application{}, err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO user_roles(user_id,role_id) SELECT $1,id FROM roles WHERE code='artisan' ON CONFLICT DO NOTHING`, userID); err != nil {
			return Application{}, err
		}
		if err = eventTarget(ctx, tx, "ARTISAN_MEMBERSHIP_ACTIVE", actor, "artisan_membership", "artisan_membership", membershipID, reason, "NOT_STARTED", "ACTIVE"); err != nil {
			return Application{}, err
		}
		if err = tx.Commit(ctx); err != nil {
			return Application{}, err
		}
		return scanApplication(r.pool.QueryRow(ctx, applicationSelect+`WHERE a.id=$1`, profileID))
	}
	if err != nil {
		return Application{}, err
	}
	if previous == status {
		return domain.Application{}, domain.ErrInvalidTransition
	}
	if _, err = tx.Exec(ctx, `UPDATE artisan_memberships SET status=$2,reason=NULLIF($3,''),suspended_at=CASE WHEN $2='SUSPENDED' THEN CURRENT_TIMESTAMP ELSE suspended_at END,closed_at=CASE WHEN $2='CLOSED' THEN CURRENT_TIMESTAMP ELSE closed_at END,updated_at=CURRENT_TIMESTAMP WHERE artisan_profile_id=$1`, profileID, status, reason); err != nil {
		return Application{}, err
	}
	if status == "ACTIVE" {
		if err = ensureDefaultWorkshop(ctx, tx, profileID); err != nil {
			return Application{}, err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO artisan_verifications(artisan_membership_id) VALUES($1) ON CONFLICT (artisan_membership_id) DO NOTHING`, membershipID); err != nil {
			return Application{}, err
		}
	}
	if err = eventTarget(ctx, tx, "ARTISAN_MEMBERSHIP_"+status, actor, "artisan_membership", "artisan_membership", membershipID, reason, previous, status); err != nil {
		return Application{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Application{}, err
	}
	return scanApplication(r.pool.QueryRow(ctx, applicationSelect+`WHERE a.id=$1`, profileID))
}

func ensureDefaultWorkshop(ctx context.Context, tx pgx.Tx, profileID string) error {
	if _, err := tx.Exec(ctx, `UPDATE workshops SET status='ACTIVE',is_public=true,updated_at=CURRENT_TIMESTAMP WHERE artisan_profile_id=$1 AND is_default=true`, profileID); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, `INSERT INTO workshops(artisan_profile_id,name,wilaya,location_text,status,is_default,is_public) SELECT id,COALESCE(NULLIF(btrim(workshop_name),''),public_display_name),wilaya,location_text,'ACTIVE',true,true FROM artisan_profiles WHERE id=$1 AND NOT EXISTS (SELECT 1 FROM workshops WHERE artisan_profile_id=$1 AND is_default=true)`, profileID)
	return err
}
func (r *PostgresRepository) MineVerification(ctx context.Context, userID string) (domain.Verification, error) {
	var v domain.Verification
	err := r.pool.QueryRow(ctx, `SELECT v.id,v.artisan_membership_id,v.status,COALESCE(v.reason,''),COALESCE(v.decided_by_user_id::text,''),v.decided_at FROM artisan_verifications v JOIN artisan_memberships m ON m.id=v.artisan_membership_id WHERE m.user_id=$1`, userID).Scan(&v.ID, &v.MembershipID, &v.Status, &v.Reason, &v.DecidedBy, &v.DecidedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Verification{Status: "NOT_SUBMITTED"}, nil
	}
	return v, err
}
func (r *PostgresRepository) ListVerifications(ctx context.Context, status string, limit, offset int) ([]domain.Verification, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT count(*) FROM artisan_verifications WHERE ($1='' OR status=$1)`, status).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, `SELECT v.id,v.artisan_membership_id,v.status,COALESCE(v.reason,''),COALESCE(v.decided_by_user_id::text,''),v.decided_at FROM artisan_verifications v WHERE ($1='' OR v.status=$1) ORDER BY v.updated_at DESC LIMIT $2 OFFSET $3`, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []domain.Verification{}
	for rows.Next() {
		var v domain.Verification
		if err = rows.Scan(&v.ID, &v.MembershipID, &v.Status, &v.Reason, &v.DecidedBy, &v.DecidedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, v)
	}
	return out, total, rows.Err()
}
func (r *PostgresRepository) DecideVerification(ctx context.Context, actor string, input domain.VerificationDecision) (domain.Verification, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Verification{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var v domain.Verification
	if err = tx.QueryRow(ctx, `UPDATE artisan_verifications SET status=$2,reason=NULLIF($3,''),decided_by_user_id=$4,decided_at=CURRENT_TIMESTAMP,updated_at=CURRENT_TIMESTAMP WHERE id=$1 RETURNING id,artisan_membership_id,status,COALESCE(reason,''),COALESCE(decided_by_user_id::text,''),decided_at`, input.ID, input.Status, input.Reason, actor).Scan(&v.ID, &v.MembershipID, &v.Status, &v.Reason, &v.DecidedBy, &v.DecidedAt); errors.Is(err, pgx.ErrNoRows) {
		return domain.Verification{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Verification{}, err
	}
	if err = eventTarget(ctx, tx, "ARTISAN_VERIFICATION_"+input.Status, actor, "artisan_membership", "artisan_membership", v.MembershipID, input.Reason, "PENDING", input.Status); err != nil {
		return domain.Verification{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Verification{}, err
	}
	return v, nil
}
func (r *PostgresRepository) MarkVerificationPending(ctx context.Context, userID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	_, err = tx.Exec(ctx, `
		WITH changed AS (
			UPDATE artisan_verifications v
			SET status='PENDING',updated_at=CURRENT_TIMESTAMP
			FROM artisan_memberships m
			WHERE v.artisan_membership_id=m.id AND m.user_id=$1 AND v.status IN ('NOT_SUBMITTED','CHANGES_REQUESTED')
			RETURNING v.artisan_membership_id
		)
		INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload)
		SELECT 'ARTISAN_VERIFICATION_SUBMITTED','artisan_membership',artisan_membership_id,jsonb_build_object('membershipId',artisan_membership_id::text,'status','PENDING')
		FROM changed`, userID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
