package repositories

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/product/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresProductWorkshopOwnershipAndMoveRules(t *testing.T) {
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

	var userID, otherUserID, artisanID, otherArtisanID, categoryID, workshopOne, workshopTwo, foreignWorkshop, inactiveWorkshop, productID string
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name)VALUES('product-repository-'||gen_random_uuid()::text,'Product Test')RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name)VALUES('product-repository-other-'||gen_random_uuid()::text,'Other Product Test')RETURNING id`).Scan(&otherUserID); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM inventory_movements WHERE product_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM product_translations WHERE product_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM workshops WHERE artisan_profile_id IN ($1,$2)`, artisanID, otherArtisanID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_verifications WHERE artisan_membership_id IN (SELECT id FROM artisan_memberships WHERE artisan_profile_id IN ($1,$2))`, artisanID, otherArtisanID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_memberships WHERE artisan_profile_id IN ($1,$2)`, artisanID, otherArtisanID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_profiles WHERE id IN ($1,$2)`, artisanID, otherArtisanID)
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE id=$1`, categoryID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1,$2)`, userID, otherUserID)
	}()
	if err = pool.QueryRow(ctx, `INSERT INTO artisan_profiles(user_id,public_display_name,status)VALUES($1,'Product Test Artisan','APPROVED')RETURNING id`, userID).Scan(&artisanID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO artisan_profiles(user_id,public_display_name,status)VALUES($1,'Other Product Artisan','APPROVED')RETURNING id`, otherUserID).Scan(&otherArtisanID); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `INSERT INTO artisan_memberships(user_id,artisan_profile_id,status,activated_at)VALUES($1,$2,'ACTIVE',CURRENT_TIMESTAMP),($3,$4,'ACTIVE',CURRENT_TIMESTAMP)`, userID, artisanID, otherUserID, otherArtisanID); err != nil {
		t.Fatal(err)
	}
	for _, item := range []struct {
		query  string
		args   []any
		target *string
	}{
		{`INSERT INTO workshops(artisan_profile_id,name,status,is_default,is_public)VALUES($1,'One','ACTIVE',true,true)RETURNING id`, []any{artisanID}, &workshopOne},
		{`INSERT INTO workshops(artisan_profile_id,name,status,is_default,is_public)VALUES($1,'Two','ACTIVE',false,true)RETURNING id`, []any{artisanID}, &workshopTwo},
		{`INSERT INTO workshops(artisan_profile_id,name,status,is_default,is_public)VALUES($1,'Foreign','ACTIVE',true,true)RETURNING id`, []any{otherArtisanID}, &foreignWorkshop},
		{`INSERT INTO workshops(artisan_profile_id,name,status,is_default,is_public)VALUES($1,'Inactive','INACTIVE',false,false)RETURNING id`, []any{artisanID}, &inactiveWorkshop},
	} {
		if err = pool.QueryRow(ctx, item.query, item.args...).Scan(item.target); err != nil {
			t.Fatal(err)
		}
	}
	if err = pool.QueryRow(ctx, `INSERT INTO categories(slug,display_name)VALUES('product-repository-'||gen_random_uuid()::text,'Product category')RETURNING id`).Scan(&categoryID); err != nil {
		t.Fatal(err)
	}
	repository := NewPostgresRepository(pool)
	input := domain.Input{WorkshopID: workshopOne, CategoryID: categoryID, ProductType: "ARTISAN_SPECIFIC", PriceMinor: 1000, Currency: "EUR"}
	item, err := repository.Create(ctx, userID, input)
	if err != nil {
		t.Fatal(err)
	}
	productID = item.ID
	if item.WorkshopID != workshopOne || item.Status != "DRAFT" {
		t.Fatalf("created product=%#v", item)
	}
	if _, err = repository.Create(ctx, userID, domain.Input{WorkshopID: foreignWorkshop, CategoryID: categoryID, ProductType: "ARTISAN_SPECIFIC", PriceMinor: 1000, Currency: "EUR"}); !errors.Is(err, domain.ErrWorkshopNotOwned) {
		t.Fatalf("foreign workshop error=%v", err)
	}
	if _, err = repository.Create(ctx, userID, domain.Input{WorkshopID: inactiveWorkshop, CategoryID: categoryID, ProductType: "ARTISAN_SPECIFIC", PriceMinor: 1000, Currency: "EUR"}); !errors.Is(err, domain.ErrWorkshopNotOwned) {
		t.Fatalf("inactive workshop error=%v", err)
	}
	input.WorkshopID = workshopTwo
	if _, err = repository.Update(ctx, userID, productID, input); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `INSERT INTO inventory_movements(product_id,movement_type,quantity_delta,reference_key,reason,actor_user_id)VALUES($1::text::uuid,'ACCEPTED',1,'product-repository:'||$1::text,'Repository test movement',$2::text::uuid)`, productID, userID); err != nil {
		t.Fatal(err)
	}
	input.WorkshopID = workshopOne
	if _, err = repository.Update(ctx, userID, productID, input); !errors.Is(err, domain.ErrWorkshopProtectedHistory) {
		t.Fatalf("protected move error=%v", err)
	}
	archived, err := repository.Archive(ctx, userID, productID)
	if err != nil {
		t.Fatal(err)
	}
	if archived.Status != "ARCHIVED" {
		t.Fatalf("archived product=%#v", archived)
	}
	var publishedAt any
	if err = pool.QueryRow(ctx, `SELECT published_at FROM products WHERE id=$1`, productID).Scan(&publishedAt); err != nil {
		t.Fatal(err)
	}
	if publishedAt != nil {
		t.Fatalf("expected archived product to have no published_at, got %v", publishedAt)
	}
	if _, err = repository.Archive(ctx, userID, productID); !errors.Is(err, domain.ErrInvalidTransition) {
		t.Fatalf("second archive error=%v", err)
	}
	workshops, err := repository.ListOwnedWorkshops(ctx, userID)
	if err != nil || len(workshops) != 3 {
		t.Fatalf("workshops=%#v err=%v", workshops, err)
	}
}
