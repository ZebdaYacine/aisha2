package customer

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
)

func TestPostgresAddressOwnershipAndSingleDefault(t *testing.T) {
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
	var firstUser, secondUser string
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name)VALUES('customer-one@example.test','One')RETURNING id`).Scan(&firstUser); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name)VALUES('customer-two@example.test','Two')RETURNING id`).Scan(&secondUser); err != nil {
		t.Fatal(err)
	}
	defer func() { _, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN($1,$2)`, firstUser, secondUser) }()
	repository := NewPostgresRepository(pool)
	input := AddressInput{FullName: "One", Line1: "Street", City: "Alger", PostalCode: "16000", Country: "Algeria", Default: true}
	first, err := repository.CreateAddress(ctx, firstUser, input)
	if err != nil {
		t.Fatal(err)
	}
	second, err := repository.CreateAddress(ctx, firstUser, input)
	if err != nil {
		t.Fatal(err)
	}
	items, err := repository.Addresses(ctx, firstUser)
	if err != nil {
		t.Fatal(err)
	}
	defaults := 0
	for _, item := range items {
		if item.Default {
			defaults++
		}
	}
	if defaults != 1 || !second.Default {
		t.Fatalf("defaults=%d second=%v", defaults, second.Default)
	}
	if err = repository.DeleteAddress(ctx, secondUser, first.ID); err != ErrNotFound {
		t.Fatalf("cross-owner delete error=%v", err)
	}
}
