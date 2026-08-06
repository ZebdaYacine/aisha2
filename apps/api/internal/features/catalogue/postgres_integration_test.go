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
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM product_media WHERE product_id IN ($1,$2)`, activeID, draftID)
		_, _ = pool.Exec(ctx, `DELETE FROM product_translations WHERE product_id IN ($1,$2)`, activeID, draftID)
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE id IN ($1,$2)`, activeID, draftID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_profile_translations WHERE artisan_profile_id=$1`, artisanID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_profiles WHERE id=$1`, artisanID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM category_translations WHERE category_id=$1`, categoryID)
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE id=$1`, categoryID)
	}()
	if err = pool.QueryRow(ctx, `INSERT INTO artisan_profiles(user_id,public_display_name,status,profile_image_object_key)VALUES($1,'Maker','APPROVED','/images/aisha/test-artisan.jpg')RETURNING id`, userID).Scan(&artisanID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO categories(slug,display_name)VALUES('catalogue-test-'||gen_random_uuid()::text,'Test category')RETURNING id`).Scan(&categoryID); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `INSERT INTO category_translations(category_id,locale,name)VALUES($1,'fr','Catégorie de test')`, categoryID); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `INSERT INTO artisan_profile_translations(artisan_profile_id,locale,display_name,biography)VALUES($1,'en','Maker EN','Story EN'),($1,'fr','Artisan FR','Story FR')`, artisanID); err != nil {
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
	if _, err = pool.Exec(ctx, `INSERT INTO product_translations(product_id,locale,name,description)VALUES($1,'fr','Vase français','Description française')`, activeID); err != nil {
		t.Fatal(err)
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
	artisanPage, err := repository.Artisans(ctx, PageRequest{Locale: "en", Page: 1, PageSize: 100})
	if err != nil {
		t.Fatal(err)
	}
	artisanFound := false
	for _, item := range artisanPage.Items {
		if item.ID == artisanID {
			artisanFound = true
			if len(item.Media) != 1 || item.Media[0] != "/images/aisha/test-artisan.jpg" {
				t.Fatalf("artisan media=%v", item.Media)
			}
		}
	}
	if !artisanFound {
		t.Fatal("approved artisan missing")
	}
	localizedCategories, err := repository.Categories(ctx, PageRequest{Locale: "fr", Page: 1, PageSize: 100})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range localizedCategories.Items {
		if item.ID == categoryID && item.Name != "Catégorie de test" {
			t.Fatalf("category name=%q", item.Name)
		}
	}
	localizedProducts, err := repository.Products(ctx, PageRequest{Locale: "fr", Page: 1, PageSize: 100}, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range localizedProducts.Items {
		if item.ID == activeID {
			if item.Name != "Vase français" || item.ArtisanName != "Artisan FR" {
				t.Fatalf("localized product=%q artisan=%q", item.Name, item.ArtisanName)
			}
		}
	}
}
