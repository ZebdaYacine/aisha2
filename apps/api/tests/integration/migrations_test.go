package tests

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const expectedMigrationCount = 21

func migrationsPath(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve migration test path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", "db", "migrations"))
}

func TestMigrationFilesArePairedAndDeclareDeletionBehaviour(t *testing.T) {
	path := migrationsPath(t)
	upFiles, err := filepath.Glob(filepath.Join(path, "*.up.sql"))
	if err != nil {
		t.Fatal(err)
	}
	downFiles, err := filepath.Glob(filepath.Join(path, "*.down.sql"))
	if err != nil {
		t.Fatal(err)
	}
	if len(upFiles) != expectedMigrationCount || len(downFiles) != expectedMigrationCount {
		t.Fatalf("expected %d up/down pairs, got %d up and %d down", expectedMigrationCount, len(upFiles), len(downFiles))
	}

	for _, upFile := range upFiles {
		contents, err := os.ReadFile(upFile)
		if err != nil {
			t.Fatal(err)
		}
		sqlText := strings.ToUpper(string(contents))
		if strings.Contains(sqlText, "REFERENCES ") && !strings.Contains(sqlText, "ON DELETE ") {
			t.Fatalf("%s contains a foreign key without explicit deletion behaviour", filepath.Base(upFile))
		}
	}
}

func TestMigrationsRoundTrip(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is required for the PostgreSQL migration integration test")
	}

	migrator, err := migrate.New("file://"+filepath.ToSlash(migrationsPath(t)), databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		sourceErr, databaseErr := migrator.Close()
		if sourceErr != nil || databaseErr != nil {
			t.Errorf("close migrator: source=%v database=%v", sourceErr, databaseErr)
		}
	}()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("initial migration up: %v", err)
	}
	assertFoundation(t, databaseURL)
	if err := migrator.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("migration down: %v", err)
	}
	assertTableAbsent(t, databaseURL, "users")
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("second migration up: %v", err)
	}
	assertFoundation(t, databaseURL)
}

func assertFoundation(t *testing.T, databaseURL string) {
	t.Helper()
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	var roleCount, categoryCount, modelTableCount int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM roles").Scan(&roleCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM categories").Scan(&categoryCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, `
		SELECT count(*) FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = ANY($1)
	`, []string{
		"users", "sessions", "artisan_profiles", "artisan_profile_translations",
		"artisan_documents", "artisan_media", "categories", "artisan_profile_categories",
		"products", "product_translations",
		"product_media", "audit_events", "outbox_events", "idempotency_keys",
		"password_reset_tokens",
		"addresses",
		"category_translations",
		"product_submissions",
		"workshops", "workshop_translations", "inventory_movements",
		"product_moderation_decisions", "orders", "order_items", "stock_reservations",
		"payment_attempts", "order_returns", "shipment_events",
		"carts", "cart_items", "wishlist_items",
		"artisan_memberships", "artisan_verifications",
		"warehouse_receptions", "warehouse_inspections", "warehouse_evidence",
	}).Scan(&modelTableCount); err != nil {
		t.Fatal(err)
	}
	if roleCount != 6 || categoryCount != 12 || modelTableCount != 36 {
		t.Fatalf("unexpected foundation counts: roles=%d categories=%d tables=%d", roleCount, categoryCount, modelTableCount)
	}

	var auditID string
	if err := db.QueryRowContext(ctx, `
		INSERT INTO audit_events (event_type, target_type)
		VALUES ('MIGRATION_TEST', 'migration')
		RETURNING id
	`).Scan(&auditID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "UPDATE audit_events SET reason = 'mutated' WHERE id = $1", auditID); err == nil {
		t.Fatal("expected audit event update to be rejected")
	}
}

func assertTableAbsent(t *testing.T, databaseURL, tableName string) {
	t.Helper()
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var regclass *string
	if err := db.QueryRow("SELECT to_regclass($1)", "public."+tableName).Scan(&regclass); err != nil {
		t.Fatal(err)
	}
	if regclass != nil {
		t.Fatalf("expected %s to be removed by down migrations", tableName)
	}
}
