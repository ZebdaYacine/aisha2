package catalogue

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func TestPostgresCataloguePublicationRules(t *testing.T) {
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
	var userID, artisanID, categoryID, activeID, draftID string
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name)VALUES('catalogue-test@example.test','Maker')RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	defer func() { _, _ = pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, userID) }()
	if err = pool.QueryRow(ctx, `INSERT INTO artisan_profiles(user_id,public_display_name,status)VALUES($1,'Maker','APPROVED')RETURNING id`, userID).Scan(&artisanID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `SELECT id FROM categories WHERE is_active=true LIMIT 1`).Scan(&categoryID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO products(artisan_profile_id,category_id,product_type,status,price_minor,currency,published_at)VALUES($1,$2,'ARTISAN_SPECIFIC','ACTIVE',1000,'EUR',CURRENT_TIMESTAMP)RETURNING id`, artisanID, categoryID).Scan(&activeID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO products(artisan_profile_id,category_id,product_type,status,price_minor,currency)VALUES($1,$2,'ARTISAN_SPECIFIC','DRAFT',1000,'EUR')RETURNING id`, artisanID, categoryID).Scan(&draftID); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{activeID, draftID} {
		if _, err = pool.Exec(ctx, `INSERT INTO product_translations(product_id,locale,name,description)VALUES($1,'en','Vase','Description')`, id); err != nil {
			t.Fatal(err)
		}
	}
	if _, err = pool.Exec(ctx, `INSERT INTO product_media(product_id,media_kind,object_key,media_type,size_bytes,visibility)VALUES($1,'IMAGE','catalogue/public.jpg','image/jpeg',10,'PUBLIC'),($1,'IMAGE','catalogue/private.jpg','image/jpeg',10,'PRIVATE')`, activeID); err != nil {
		t.Fatal(err)
	}
	repository := NewPostgresRepository(pool)
	page, err := repository.Products(ctx, PageRequest{Locale: "en", Page: 1, PageSize: 100}, "")
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, item := range page.Items {
		if item.ID == draftID {
			t.Fatal("draft product leaked")
		}
		if item.ID == activeID {
			found = true
			if len(item.Media) != 1 || item.Media[0] != "catalogue/public.jpg" {
				t.Fatalf("media=%v", item.Media)
			}
		}
	}
	if !found {
		t.Fatal("active product missing")
	}
}
