package config

import "testing"

func TestLoadUsesFoundationDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("DATABASE_MAX_OPEN_CONNS", "12")
	t.Setenv("DATABASE_MAX_IDLE_CONNS", "3")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AppName != "AISHA" || cfg.DatabaseMaxOpen != 12 || cfg.DatabaseMaxIdle != 3 {
		t.Fatalf("unexpected config: %#v", cfg)
	}
	if len(cfg.AllowedOrigins) != 1 || cfg.AllowedOrigins[0] != "http://localhost:3000" {
		t.Fatalf("unexpected origins: %#v", cfg.AllowedOrigins)
	}
}

func TestLoadRejectsInvalidConnectionLimit(t *testing.T) {
	t.Setenv("DATABASE_MAX_OPEN_CONNS", "0")
	if _, err := Load(); err == nil {
		t.Fatal("expected invalid connection limit to fail")
	}
}
