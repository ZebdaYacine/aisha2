package artisan

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresSubmissionAndApprovalAreAtomic(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	var userID, adminID, categoryID string
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES('artisan-integration@example.test','Applicant') RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES('admin-integration@example.test','Admin') RETURNING id`).Scan(&adminID); err != nil {
		t.Fatal(err)
	}
	defer func() { _, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1,$2)`, userID, adminID) }()
	if err = pool.QueryRow(ctx, `SELECT id FROM categories WHERE is_active=true LIMIT 1`).Scan(&categoryID); err != nil {
		t.Fatal(err)
	}
	repository := NewPostgresRepository(pool)
	application, err := repository.Submit(ctx, userID, ApplicationInput{PublicDisplayName: "Integration Atelier", Wilaya: "Alger", ContactVisibility: "PRIVATE", CategoryIDs: []string{categoryID}, Translations: []Translation{{Locale: "en", Biography: "Workshop story"}}})
	if err != nil {
		t.Fatal(err)
	}
	if application.Status != "SUBMITTED" {
		t.Fatalf("status=%q", application.Status)
	}
	approved, err := repository.Decide(ctx, adminID, application.ID, "APPROVED", "")
	if err != nil {
		t.Fatal(err)
	}
	if approved.Status != "APPROVED" {
		t.Fatalf("status=%q", approved.Status)
	}
	var roles, audits, outbox int
	if err = pool.QueryRow(ctx, `SELECT count(*) FROM user_roles ur JOIN roles r ON r.id=ur.role_id WHERE ur.user_id=$1 AND r.code='artisan'`, userID).Scan(&roles); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `SELECT count(*) FROM audit_events WHERE target_id=$1`, application.ID).Scan(&audits); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `SELECT count(*) FROM outbox_events WHERE aggregate_id=$1`, application.ID).Scan(&outbox); err != nil {
		t.Fatal(err)
	}
	if roles != 1 || audits != 2 || outbox != 2 {
		t.Fatalf("roles=%d audits=%d outbox=%d", roles, audits, outbox)
	}
}
