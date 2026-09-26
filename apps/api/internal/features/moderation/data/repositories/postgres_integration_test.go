package repositories

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/moderation/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresModerationQueueAndActivationGates(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	var actorID, artisanUserID, artisanID, categoryID, workshopID, productID, submissionID string
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM product_moderation_decisions WHERE product_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM audit_events WHERE target_type='product' AND target_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM outbox_events WHERE aggregate_type='product' AND aggregate_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM inventory_movements WHERE product_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM product_media WHERE product_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM product_submissions WHERE product_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM workshops WHERE id=$1`, workshopID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_memberships WHERE artisan_profile_id=$1`, artisanID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_profiles WHERE id=$1`, artisanID)
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE id=$1`, categoryID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1,$2)`, actorID, artisanUserID)
	}()

	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES ('moderation-repository-actor-'||gen_random_uuid()::text,'Moderation Actor') RETURNING id`).Scan(&actorID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES ('moderation-repository-artisan-'||gen_random_uuid()::text,'Moderation Artisan') RETURNING id`).Scan(&artisanUserID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO artisan_profiles(user_id,public_display_name,status) VALUES ($1,'Moderation Artisan','APPROVED') RETURNING id`, artisanUserID).Scan(&artisanID); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `INSERT INTO artisan_memberships(user_id,artisan_profile_id,status,activated_at) VALUES ($1,$2,'ACTIVE',CURRENT_TIMESTAMP)`, artisanUserID, artisanID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO workshops(artisan_profile_id,name,status,is_default,is_public) VALUES ($1,'Moderation Workshop','ACTIVE',true,true) RETURNING id`, artisanID).Scan(&workshopID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO categories(slug,display_name) VALUES ('moderation-repository-'||gen_random_uuid()::text,'Moderation Category') RETURNING id`).Scan(&categoryID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO products(artisan_profile_id,category_id,workshop_id,product_type,status,price_minor,currency) VALUES ($1,$2,$3,'ARTISAN_SPECIFIC','PENDING_REVIEW',2500,'EUR') RETURNING id`, artisanID, categoryID, workshopID).Scan(&productID); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `INSERT INTO product_media(product_id,media_kind,object_key,media_type,size_bytes) VALUES ($1::uuid,'IMAGE','moderation-repository/'||$1::text,'image/jpeg',100)`, productID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO product_submissions(product_id,version,submitted_by_user_id,snapshot) VALUES ($1,1,$2,'{"name":"Moderation product"}'::jsonb) RETURNING id`, productID, artisanUserID).Scan(&submissionID); err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresRepository(pool)
	items, total, err := repository.ListQueue(ctx, "PENDING_REVIEW", 1000, 0)
	if err != nil {
		t.Fatal(err)
	}
	var pendingFound bool
	for _, item := range items {
		if item.ProductID == productID {
			pendingFound = item.SubmissionID == submissionID && item.ProductStatus == "PENDING_REVIEW"
			break
		}
	}
	if !pendingFound || total < 1 {
		t.Fatalf("pending queue total=%d items=%#v", total, items)
	}

	if _, err = repository.Decide(ctx, actorID, domain.DecisionInput{ID: submissionID, Action: "APPROVE"}); err != nil {
		t.Fatal(err)
	}
	items, total, err = repository.ListQueue(ctx, "APPROVED", 1000, 0)
	if err != nil {
		t.Fatal(err)
	}
	var approvedFound bool
	for _, item := range items {
		if item.ProductID == productID {
			approvedFound = item.ProductStatus == "APPROVED"
			break
		}
	}
	if !approvedFound || total < 1 {
		t.Fatalf("approved queue total=%d items=%#v", total, items)
	}
	if _, err = repository.Decide(ctx, actorID, domain.DecisionInput{ID: productID, Action: "ACTIVATE"}); !errors.Is(err, domain.ErrActivationNotReady) {
		t.Fatalf("activation without stock error=%v", err)
	}

	if _, err = pool.Exec(ctx, `INSERT INTO inventory_movements(product_id,movement_type,quantity_delta,reference_key,reason,actor_user_id) VALUES ($1::uuid,'ACCEPTED',1,'moderation-repository:'||$1::text,'Accepted test stock',$2::uuid)`, productID, artisanUserID); err != nil {
		t.Fatal(err)
	}
	if _, err = repository.Decide(ctx, actorID, domain.DecisionInput{ID: productID, Action: "ACTIVATE"}); err != nil {
		t.Fatal(err)
	}
	var status, visibility string
	if err = pool.QueryRow(ctx, `SELECT status FROM products WHERE id=$1`, productID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `SELECT visibility FROM product_media WHERE product_id=$1`, productID).Scan(&visibility); err != nil {
		t.Fatal(err)
	}
	if status != "ACTIVE" || visibility != "PUBLIC" {
		t.Fatalf("status=%q visibility=%q", status, visibility)
	}
	var auditCount, outboxCount int
	if err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM audit_events WHERE target_type='product' AND target_id=$1 AND event_type='PRODUCT_MODERATION_DECIDED'`, productID).Scan(&auditCount); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM outbox_events WHERE aggregate_type='product' AND aggregate_id=$1 AND event_type='PRODUCT_MODERATION_DECIDED'`, productID).Scan(&outboxCount); err != nil {
		t.Fatal(err)
	}
	if auditCount != 2 || outboxCount != 2 {
		t.Fatalf("audit=%d outbox=%d", auditCount, outboxCount)
	}
}
