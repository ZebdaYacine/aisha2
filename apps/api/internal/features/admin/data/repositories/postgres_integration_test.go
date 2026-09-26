package repositories

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresUserStatusRevokesSessionsAndPreservesAudit(t *testing.T) {
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

	var actorID, userID, sessionID string
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES('admin-repository-actor-'||gen_random_uuid()::text,'Admin actor') RETURNING id`).Scan(&actorID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES('admin-repository-user-'||gen_random_uuid()::text,'Managed user') RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM audit_events WHERE actor_user_id=$1 OR target_id=$2`, actorID, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM outbox_events WHERE aggregate_id=$1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM sessions WHERE user_id=$1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1,$2)`, actorID, userID)
	}()
	if err = pool.QueryRow(ctx, `INSERT INTO sessions(user_id,token_hash,family_id,expires_at) VALUES($1,'admin-status-test-'||gen_random_uuid()::text,gen_random_uuid(),CURRENT_TIMESTAMP+INTERVAL '1 hour') RETURNING id`, userID).Scan(&sessionID); err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresRepository(pool)
	item, err := repository.SetUserStatus(ctx, actorID, userID, "SUSPENDED", "policy review")
	if err != nil {
		t.Fatal(err)
	}
	if item.Status != "SUSPENDED" {
		t.Fatalf("status=%q", item.Status)
	}
	var revoked bool
	if err = pool.QueryRow(ctx, `SELECT revoked_at IS NOT NULL FROM sessions WHERE id=$1`, sessionID).Scan(&revoked); err != nil {
		t.Fatal(err)
	}
	if !revoked {
		t.Fatal("expected active sessions to be revoked")
	}
	var count int
	if err = pool.QueryRow(ctx, `SELECT count(*) FROM audit_events WHERE target_id=$1 AND event_type='USER_STATUS_CHANGED'`, userID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("audit count=%d", count)
	}
	if _, err = repository.SetUserStatus(ctx, actorID, userID, "SUSPENDED", "again"); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("same status error=%v", err)
	}
}
