package repositories

import (
	"context"
	"fmt"

	"github.com/aisha-platform/aisha/apps/api/internal/features/cart/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const listQuery = `SELECT ci.product_id,COALESCE(pt.name,p.product_type),COALESCE(a.public_display_name,''),COALESCE(w.name,''),ci.quantity,p.price_minor,p.currency,
 GREATEST(COALESCE((SELECT SUM(quantity_delta) FROM inventory_movements im WHERE im.product_id=p.id AND im.stock_bucket='AVAILABLE'),0),0),
 (p.status='ACTIVE' AND p.published_at IS NOT NULL AND a.status='APPROVED' AND w.status='ACTIVE' AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE')),
 COALESCE((SELECT pm.object_key FROM product_media pm WHERE pm.product_id=p.id AND pm.visibility='PUBLIC' AND pm.media_kind='IMAGE' ORDER BY pm.sort_order,pm.created_at LIMIT 1),''),
 ci.updated_at
 FROM cart_items ci JOIN carts c ON c.id=ci.cart_id JOIN products p ON p.id=ci.product_id
 JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN workshops w ON w.id=p.workshop_id
 LEFT JOIN product_translations pt ON pt.product_id=p.id AND pt.locale='en'
 WHERE c.user_id=$1 ORDER BY ci.updated_at DESC,ci.product_id`

func (r *PostgresRepository) List(ctx context.Context, userID string) ([]domain.Item, error) {
	rows, err := r.pool.Query(ctx, listQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("list cart: %w", err)
	}
	defer rows.Close()
	out := []domain.Item{}
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.ArtisanName, &item.WorkshopName, &item.Quantity, &item.PriceMinor, &item.Currency, &item.Available, &item.Active, &item.Image, &item.UpdatedAt); err != nil {
			return nil, err
		}
		item.Warning = warning(item)
		out = append(out, item)
	}
	return out, rows.Err()
}

func warning(item domain.Item) string {
	if !item.Active {
		return "PRODUCT_UNAVAILABLE"
	}
	if item.Available < int64(item.Quantity) {
		return "INSUFFICIENT_STOCK"
	}
	return ""
}

func (r *PostgresRepository) ensureCart(ctx context.Context, tx pgx.Tx, userID string) (string, error) {
	var id string
	err := tx.QueryRow(ctx, `INSERT INTO carts(user_id) VALUES($1) ON CONFLICT(user_id) DO UPDATE SET updated_at=CURRENT_TIMESTAMP RETURNING id`, userID).Scan(&id)
	return id, err
}

func (r *PostgresRepository) Add(ctx context.Context, userID string, in domain.Input) ([]domain.Item, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	cartID, err := r.ensureCart(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	result, err := tx.Exec(ctx, `INSERT INTO cart_items(cart_id,product_id,quantity)
        SELECT $1,p.id,$3 FROM products p
        JOIN artisan_profiles a ON a.id=p.artisan_profile_id
        JOIN workshops w ON w.id=p.workshop_id
        WHERE p.id=$2 AND p.status='ACTIVE' AND p.published_at IS NOT NULL
          AND a.status='APPROVED' AND w.status='ACTIVE'
          AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE')
        ON CONFLICT(cart_id,product_id) DO UPDATE SET quantity=LEAST(100,cart_items.quantity+EXCLUDED.quantity),updated_at=CURRENT_TIMESTAMP`, cartID, in.ProductID, in.Quantity)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() == 0 {
		return nil, domain.ErrUnavailable
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.List(ctx, userID)
}

func (r *PostgresRepository) Set(ctx context.Context, userID, productID string, quantity int) ([]domain.Item, error) {
	result, err := r.mutate(ctx, userID, `INSERT INTO cart_items(cart_id,product_id,quantity)
        SELECT $1,p.id,$3 FROM products p
        JOIN artisan_profiles a ON a.id=p.artisan_profile_id
        JOIN workshops w ON w.id=p.workshop_id
        WHERE p.id=$2 AND p.status='ACTIVE' AND p.published_at IS NOT NULL
          AND a.status='APPROVED' AND w.status='ACTIVE'
          AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE')
        ON CONFLICT(cart_id,product_id) DO UPDATE SET quantity=EXCLUDED.quantity,updated_at=CURRENT_TIMESTAMP`, productID, quantity)
	return result, err
}

func (r *PostgresRepository) Remove(ctx context.Context, userID, productID string) ([]domain.Item, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	cartID, err := r.ensureCart(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM cart_items WHERE cart_id=$1 AND product_id=$2`, cartID, productID); err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.List(ctx, userID)
}

func (r *PostgresRepository) mutate(ctx context.Context, userID, query, productID string, quantity int) ([]domain.Item, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	cartID, err := r.ensureCart(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	result, err := tx.Exec(ctx, query, cartID, productID, quantity)
	if err != nil {
		return nil, err
	}
	if result.RowsAffected() == 0 {
		return nil, domain.ErrUnavailable
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.List(ctx, userID)
}

func (r *PostgresRepository) Merge(ctx context.Context, userID string, inputs []domain.Input) ([]domain.Item, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	cartID, err := r.ensureCart(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	for _, in := range inputs {
		result, execErr := tx.Exec(ctx, `INSERT INTO cart_items(cart_id,product_id,quantity)
            SELECT $1,p.id,$3 FROM products p
            JOIN artisan_profiles a ON a.id=p.artisan_profile_id
            JOIN workshops w ON w.id=p.workshop_id
            WHERE p.id=$2 AND p.status='ACTIVE' AND p.published_at IS NOT NULL
              AND a.status='APPROVED' AND w.status='ACTIVE'
              AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE')
            ON CONFLICT(cart_id,product_id) DO UPDATE SET quantity=LEAST(100,cart_items.quantity+EXCLUDED.quantity),updated_at=CURRENT_TIMESTAMP`, cartID, in.ProductID, in.Quantity)
		if execErr != nil {
			return nil, execErr
		}
		if result.RowsAffected() == 0 {
			return nil, domain.ErrUnavailable
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.List(ctx, userID)
}
