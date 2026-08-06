package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/aisha-platform/aisha/apps/api/internal/features/user/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile = domain.Profile
type Address = domain.Address
type AddressInput = domain.AddressInput

var ErrNotFound = domain.ErrNotFound

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}
func (r *PostgresRepository) Profile(ctx context.Context, userID string) (Profile, error) {
	var p Profile
	err := r.pool.QueryRow(ctx, `SELECT id,email,COALESCE(display_name,''),COALESCE(phone,'') FROM users WHERE id=$1`, userID).Scan(&p.ID, &p.Email, &p.DisplayName, &p.Phone)
	if errors.Is(err, pgx.ErrNoRows) {
		return Profile{}, ErrNotFound
	}
	return p, err
}
func (r *PostgresRepository) UpdateProfile(ctx context.Context, userID, name, phone string) (Profile, error) {
	var p Profile
	err := r.pool.QueryRow(ctx, `UPDATE users SET display_name=$2,phone=NULLIF($3,''),updated_at=CURRENT_TIMESTAMP WHERE id=$1 RETURNING id,email,display_name,COALESCE(phone,'')`, userID, name, phone).Scan(&p.ID, &p.Email, &p.DisplayName, &p.Phone)
	if err != nil {
		return Profile{}, fmt.Errorf("update customer profile: %w", err)
	}
	return p, nil
}
func (r *PostgresRepository) Addresses(ctx context.Context, userID string) ([]Address, error) {
	rows, err := r.pool.Query(ctx, `SELECT id,full_name,COALESCE(phone,''),line1,COALESCE(line2,''),city,postal_code,country,is_default FROM addresses WHERE user_id=$1 ORDER BY is_default DESC,created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("query customer addresses: %w", err)
	}
	defer rows.Close()
	items := []Address{}
	for rows.Next() {
		var a Address
		if err := rows.Scan(&a.ID, &a.FullName, &a.Phone, &a.Line1, &a.Line2, &a.City, &a.PostalCode, &a.Country, &a.Default); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}
func (r *PostgresRepository) CreateAddress(ctx context.Context, userID string, i AddressInput) (Address, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Address{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if i.Default {
		if _, err = tx.Exec(ctx, `UPDATE addresses SET is_default=false,updated_at=CURRENT_TIMESTAMP WHERE user_id=$1 AND is_default=true`, userID); err != nil {
			return Address{}, err
		}
	}
	var a Address
	err = tx.QueryRow(ctx, `INSERT INTO addresses(user_id,full_name,phone,line1,line2,city,postal_code,country,is_default)VALUES($1,$2,NULLIF($3,''),$4,NULLIF($5,''),$6,$7,$8,$9)RETURNING id,full_name,COALESCE(phone,''),line1,COALESCE(line2,''),city,postal_code,country,is_default`, userID, i.FullName, i.Phone, i.Line1, i.Line2, i.City, i.PostalCode, i.Country, i.Default).Scan(&a.ID, &a.FullName, &a.Phone, &a.Line1, &a.Line2, &a.City, &a.PostalCode, &a.Country, &a.Default)
	if err != nil {
		return Address{}, fmt.Errorf("insert customer address: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return Address{}, err
	}
	return a, nil
}
func (r *PostgresRepository) UpdateAddress(ctx context.Context, userID, id string, i AddressInput) (Address, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Address{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if i.Default {
		if _, err = tx.Exec(ctx, `UPDATE addresses SET is_default=false,updated_at=CURRENT_TIMESTAMP WHERE user_id=$1 AND id<>$2`, userID, id); err != nil {
			return Address{}, err
		}
	}
	var a Address
	err = tx.QueryRow(ctx, `UPDATE addresses SET full_name=$3,phone=NULLIF($4,''),line1=$5,line2=NULLIF($6,''),city=$7,postal_code=$8,country=$9,is_default=$10,updated_at=CURRENT_TIMESTAMP WHERE id=$2 AND user_id=$1 RETURNING id,full_name,COALESCE(phone,''),line1,COALESCE(line2,''),city,postal_code,country,is_default`, userID, id, i.FullName, i.Phone, i.Line1, i.Line2, i.City, i.PostalCode, i.Country, i.Default).Scan(&a.ID, &a.FullName, &a.Phone, &a.Line1, &a.Line2, &a.City, &a.PostalCode, &a.Country, &a.Default)
	if errors.Is(err, pgx.ErrNoRows) {
		return Address{}, ErrNotFound
	}
	if err != nil {
		return Address{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return Address{}, err
	}
	return a, nil
}
func (r *PostgresRepository) DeleteAddress(ctx context.Context, userID, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM addresses WHERE id=$2 AND user_id=$1`, userID, id)
	if err != nil {
		return fmt.Errorf("delete customer address: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
