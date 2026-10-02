package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresAuthenticationLifecycle(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is required for the authentication repository integration test")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	email := "auth-integration-" + uuid.NewString() + "@example.com"
	service := NewService(NewPostgresRepository(pool), nil, "integration-test-signing-key-at-least-32-bytes")
	user, tokens, err := service.Register(ctx, email, "initial-password-value", "Integration User")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `UPDATE users SET email_verified_at=CURRENT_TIMESTAMP WHERE id=$1`, user.ID); err != nil {
		t.Fatalf("verify integration user: %v", err)
	}
	t.Cleanup(func() {
		defer pool.Close()
		if _, cleanupErr := pool.Exec(context.Background(), `DELETE FROM users WHERE id=$1`, user.ID); cleanupErr != nil {
			t.Errorf("cleanup authentication user: %v", cleanupErr)
		}
	})

	if _, _, err := service.Login(ctx, email, "initial-password-value"); err != nil {
		t.Fatalf("login persisted user: %v", err)
	}
	rotated, err := service.Refresh(ctx, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("rotate persisted session: %v", err)
	}
	if err := service.Logout(ctx, rotated.RefreshToken); err != nil {
		t.Fatalf("logout persisted session: %v", err)
	}
	if _, err := service.Authenticate(ctx, rotated.AccessToken); err == nil {
		t.Fatal("expected logged-out access token to be rejected")
	}

	resetToken := "known-reset-token"
	repository := NewPostgresRepository(pool)
	if err := repository.StorePasswordReset(ctx, user.ID, hashToken(resetToken), time.Now().UTC().Add(time.Hour)); err != nil {
		t.Fatalf("store persisted reset token: %v", err)
	}
	if err := service.ResetPassword(ctx, resetToken, "replacement-password-value"); err != nil {
		t.Fatalf("reset persisted password: %v", err)
	}
	if _, _, err := service.Login(ctx, email, "replacement-password-value"); err != nil {
		t.Fatalf("login with replacement password: %v", err)
	}
}
