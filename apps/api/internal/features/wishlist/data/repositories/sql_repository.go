package repositories

import (
	"context"
	"fmt"
	"github.com/aisha-platform/aisha/apps/api/internal/features/wishlist/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(p *pgxpool.Pool) *PostgresRepository { return &PostgresRepository{pool: p} }

const listQuery = `SELECT wi.product_id,COALESCE(pt.name,p.product_type),p.price_minor,p.currency,(p.status='ACTIVE' AND p.published_at IS NOT NULL AND a.status='APPROVED' AND w.status='ACTIVE' AND EXISTS(SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE')),wi.created_at FROM wishlist_items wi JOIN products p ON p.id=wi.product_id JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN workshops w ON w.id=p.workshop_id LEFT JOIN product_translations pt ON pt.product_id=p.id AND pt.locale='en' WHERE wi.user_id=$1 ORDER BY wi.created_at DESC`

func (r *PostgresRepository) List(ctx context.Context, userID string) ([]domain.Item, error) {
	rows, err := r.pool.Query(ctx, listQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("list wishlist: %w", err)
	}
	defer rows.Close()
	out := []domain.Item{}
	for rows.Next() {
		var v domain.Item
		if err := rows.Scan(&v.ProductID, &v.ProductName, &v.PriceMinor, &v.Currency, &v.Active, &v.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (r *PostgresRepository) Add(ctx context.Context, userID, id string) ([]domain.Item, error) {
	_, err := r.pool.Exec(ctx, `INSERT INTO wishlist_items(user_id,product_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, userID, id)
	if err != nil {
		return nil, err
	}
	return r.List(ctx, userID)
}
func (r *PostgresRepository) Remove(ctx context.Context, userID, id string) ([]domain.Item, error) {
	_, err := r.pool.Exec(ctx, `DELETE FROM wishlist_items WHERE user_id=$1 AND product_id=$2`, userID, id)
	if err != nil {
		return nil, err
	}
	return r.List(ctx, userID)
}
