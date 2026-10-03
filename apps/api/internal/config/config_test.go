package config

import (
	"testing"
	"time"
)

func TestLoadUsesFoundationDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("WEB_BASE_URL", "")
	t.Setenv("DATABASE_MAX_OPEN_CONNS", "12")
	t.Setenv("DATABASE_MAX_IDLE_CONNS", "3")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AppName != "AISHA" || cfg.DatabaseMaxOpen != 12 || cfg.DatabaseMaxIdle != 3 {
		t.Fatalf("unexpected config: %#v", cfg)
	}
	if cfg.MinIOPublicEndpoint != "localhost:9000" {
		t.Fatalf("unexpected public MinIO endpoint: %q", cfg.MinIOPublicEndpoint)
	}
	if cfg.WebBaseURL != "http://localhost:3033" {
		t.Fatalf("unexpected development web base URL: %q", cfg.WebBaseURL)
	}
	wantOrigins := []string{"http://localhost:3033", "http://127.0.0.1:3033", "http://167.86.79.16"}
	if len(cfg.AllowedOrigins) != len(wantOrigins) {
		t.Fatalf("unexpected origins: %#v", cfg.AllowedOrigins)
	}
	for index, origin := range wantOrigins {
		if cfg.AllowedOrigins[index] != origin {
			t.Fatalf("unexpected origins: %#v", cfg.AllowedOrigins)
		}
	}
	if cfg.AuthRateLimitMax != 10 || cfg.AuthRateLimitWindow != 60*time.Second {
		t.Fatalf("unexpected authentication rate limit: max=%d window=%s", cfg.AuthRateLimitMax, cfg.AuthRateLimitWindow)
	}
	if cfg.SMTPPort != 465 || cfg.SMTPHost != "" || cfg.SMTPPassword != "" {
		t.Fatalf("unexpected SMTP defaults: %#v", cfg)
	}
}

func TestLoadUsesProductionWebBaseURLWhenUnset(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("WEB_BASE_URL", "")
	t.Setenv("DATABASE_URL", "postgres://aisha:production-secret@postgres:5432/aisha?sslmode=disable")
	t.Setenv("MINIO_SECRET_KEY", "production-minio-secret")
	t.Setenv("AUTH_SIGNING_KEY", "production-signing-key-with-at-least-32-characters")
	t.Setenv("SMTP_HOST", "smtp.example.test")
	t.Setenv("SMTP_USER", "mailer@example.test")
	t.Setenv("SMTP_PASSWORD", "production-smtp-secret")
	t.Setenv("SMTP_FROM", "AISHA <mailer@example.test>")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.WebBaseURL != "https://aishasouk.com" {
		t.Fatalf("unexpected production web base URL: %q", cfg.WebBaseURL)
	}
}

func TestLoadRejectsProductionWithoutSMTP(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DATABASE_URL", "postgres://aisha:production-secret@postgres:5432/aisha?sslmode=disable")
	t.Setenv("MINIO_SECRET_KEY", "production-minio-secret")
	t.Setenv("AUTH_SIGNING_KEY", "production-signing-key-with-at-least-32-characters")
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SMTP_USER", "")
	t.Setenv("SMTP_PASSWORD", "")
	t.Setenv("SMTP_FROM", "")

	if _, err := Load(); err == nil {
		t.Fatal("expected production SMTP configuration to be required")
	}
}

func TestLoadSupportsSMTPPasswordAliases(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.hostinger.com")
	t.Setenv("SMTP_PORT", "465")
	t.Setenv("SMTP_USER", "support-team@aichasouk.com")
	t.Setenv("SMTP_PASS", "secret")
	t.Setenv("MAIL_FROM", "Aicha Souk <support-team@aichasouk.com>")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SMTPPassword != "secret" || cfg.SMTPFrom != "Aicha Souk <support-team@aichasouk.com>" {
		t.Fatalf("unexpected SMTP aliases: %#v", cfg)
	}
}

func TestLoadSupportsSMTPPasswordAliasesInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DATABASE_URL", "postgres://aisha:production-secret@postgres:5432/aisha?sslmode=disable")
	t.Setenv("MINIO_SECRET_KEY", "production-minio-secret")
	t.Setenv("AUTH_SIGNING_KEY", "production-signing-key-with-at-least-32-characters")
	t.Setenv("SMTP_HOST", "smtp.hostinger.com")
	t.Setenv("SMTP_PORT", "465")
	t.Setenv("SMTP_USER", "support-team@aichasouk.com")
	t.Setenv("SMTP_PASS", "production-smtp-secret")
	t.Setenv("MAIL_FROM", "Aicha Souk <support-team@aichasouk.com>")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SMTPPassword != "production-smtp-secret" || cfg.SMTPFrom != "Aicha Souk <support-team@aichasouk.com>" {
		t.Fatalf("unexpected production SMTP aliases: %#v", cfg)
	}
}

func TestLoadRejectsPartialSMTPConfiguration(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.hostinger.com")
	if _, err := Load(); err == nil {
		t.Fatal("expected partial SMTP configuration to fail")
	}
}

func TestLoadRejectsInvalidAuthenticationRateLimit(t *testing.T) {
	t.Setenv("AUTH_RATE_LIMIT_MAX", "0")
	if _, err := Load(); err == nil {
		t.Fatal("expected invalid authentication rate limit to fail")
	}
}

func TestLoadRejectsInvalidConnectionLimit(t *testing.T) {
	t.Setenv("DATABASE_MAX_OPEN_CONNS", "0")
	if _, err := Load(); err == nil {
		t.Fatal("expected invalid connection limit to fail")
	}
}

func TestLoadRejectsShortAuthSigningKey(t *testing.T) {
	t.Setenv("AUTH_SIGNING_KEY", "short")
	if _, err := Load(); err == nil {
		t.Fatal("expected short authentication signing key to fail")
	}
}
