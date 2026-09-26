package tests

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var seededTables = []string{
	"roles", "users", "user_roles", "sessions", "artisan_profiles",
	"artisan_profile_translations", "artisan_documents", "artisan_media", "categories",
	"artisan_profile_categories", "workshops", "workshop_translations", "products", "product_translations", "product_media", "inventory_movements",
	"addresses", "password_reset_tokens", "audit_events", "outbox_events", "idempotency_keys",
}

func seedPath(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve seed test path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "db", "seeds", "development.sql"))
}

func TestDevelopmentSeedCoversEveryApplicationTable(t *testing.T) {
	contents, err := os.ReadFile(seedPath(t))
	if err != nil {
		t.Fatal(err)
	}
	seedSQL := strings.ToLower(string(contents))
	for _, table := range seededTables {
		pattern := regexp.MustCompile(`(?m)insert\s+into\s+` + regexp.QuoteMeta(table) + `\s*\(`)
		if !pattern.MatchString(seedSQL) {
			t.Errorf("development seed does not insert into %s", table)
		}
	}
}

func TestDevelopmentSeedIsIdempotent(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is required for the PostgreSQL seed integration test")
	}
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	contents, err := os.ReadFile(seedPath(t))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for attempt := 0; attempt < 2; attempt++ {
		if _, err := db.ExecContext(ctx, string(contents)); err != nil {
			t.Fatalf("apply development seed attempt %d: %v", attempt+1, err)
		}
	}
	for _, table := range seededTables {
		var count int
		if err := db.QueryRowContext(ctx, "SELECT count(*) FROM "+table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if count == 0 {
			t.Errorf("expected seeded rows in %s", table)
		}
	}
}
