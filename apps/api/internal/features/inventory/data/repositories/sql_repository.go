package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/inventory/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const balanceQuery = `
SELECT p.id,
       COALESCE(p.product_code,''),
       COALESCE((SELECT t.name FROM product_translations t WHERE t.product_id=p.id AND t.locale='en' LIMIT 1),(SELECT t.name FROM product_translations t WHERE t.product_id=p.id ORDER BY t.locale LIMIT 1),p.product_type),
       w.id,w.name,a.id,COALESCE(a.public_display_name,''),
       GREATEST(COALESCE(SUM(CASE WHEN im.stock_bucket='AVAILABLE' THEN im.quantity_delta ELSE 0 END),0),0)
         + COALESCE((SELECT SUM(sr.quantity) FROM stock_reservations sr WHERE sr.product_id=p.id AND sr.status='HELD' AND sr.expires_at > CURRENT_TIMESTAMP),0),
       GREATEST(COALESCE(SUM(CASE WHEN im.stock_bucket='AVAILABLE' THEN im.quantity_delta ELSE 0 END),0),0),
       COALESCE((SELECT SUM(sr.quantity) FROM stock_reservations sr WHERE sr.product_id=p.id AND sr.status='HELD' AND sr.expires_at > CURRENT_TIMESTAMP),0),
       GREATEST(COALESCE(SUM(CASE WHEN im.stock_bucket='QUARANTINED' THEN im.quantity_delta ELSE 0 END),0),0),
       GREATEST(COALESCE(SUM(CASE WHEN im.stock_bucket='DAMAGED' THEN im.quantity_delta ELSE 0 END),0),0),
       GREATEST(COALESCE(SUM(CASE WHEN im.stock_bucket='REJECTED' THEN im.quantity_delta ELSE 0 END),0),0),
       GREATEST(COALESCE(SUM(CASE WHEN im.stock_bucket='SHIPPED' THEN im.quantity_delta ELSE 0 END),0),0),
       COALESCE(MAX(im.created_at),p.updated_at)
FROM products p
JOIN workshops w ON w.id=p.workshop_id
JOIN artisan_profiles a ON a.id=p.artisan_profile_id
LEFT JOIN inventory_movements im ON im.product_id=p.id
WHERE ($1 OR a.user_id=$2) AND ($3='' OR w.id=NULLIF($3,'')::uuid)
GROUP BY p.id,p.product_type,p.updated_at,w.id,w.name,a.id,a.public_display_name
ORDER BY COALESCE(MAX(im.created_at),p.updated_at) DESC,p.id`

func (r *PostgresRepository) ListBalances(ctx context.Context, actor string, global bool, workshopID string, limit, offset int) ([]domain.Balance, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN workshops w ON w.id=p.workshop_id WHERE ($1 OR a.user_id=$2) AND ($3='' OR w.id=NULLIF($3,'')::uuid)`, global, actor, workshopID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count inventory balances: %w", err)
	}
	rows, err := r.pool.Query(ctx, balanceQuery+` LIMIT $4 OFFSET $5`, global, actor, workshopID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list inventory balances: %w", err)
	}
	defer rows.Close()
	items := make([]domain.Balance, 0)
	for rows.Next() {
		var item domain.Balance
		if err := rows.Scan(&item.ProductID, &item.ProductCode, &item.ProductName, &item.WorkshopID, &item.WorkshopName, &item.ArtisanID, &item.ArtisanName, &item.OnHand, &item.Available, &item.Reserved, &item.Quarantined, &item.Damaged, &item.Rejected, &item.Shipped, &item.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) Adjust(ctx context.Context, actor string, global bool, productID string, input domain.AdjustmentInput) (domain.Movement, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Movement{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var found string
	err = tx.QueryRow(ctx, `SELECT p.id FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id WHERE p.id=$1 AND ($2 OR a.user_id=$3) FOR UPDATE`, productID, global, actor).Scan(&found)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Movement{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Movement{}, err
	}
	var available int64
	if err = tx.QueryRow(ctx, `SELECT COALESCE(SUM(quantity_delta),0) FROM inventory_movements WHERE product_id=$1 AND stock_bucket='AVAILABLE'`, productID).Scan(&available); err != nil {
		return domain.Movement{}, err
	}
	if input.QuantityDelta < 0 && available+input.QuantityDelta < 0 {
		return domain.Movement{}, domain.ErrInsufficientStock
	}
	var item domain.Movement
	err = tx.QueryRow(ctx, `INSERT INTO inventory_movements(product_id,movement_type,quantity_delta,reference_key,reason,actor_user_id,stock_bucket) VALUES($1,'ADJUSTMENT',$2,$3,$4,$5,'AVAILABLE') RETURNING id,product_id,movement_type,quantity_delta,stock_bucket,reference_key,reason,COALESCE(actor_user_id::text,''),created_at`, productID, input.QuantityDelta, input.ReferenceKey, input.Reason, actor).Scan(&item.ID, &item.ProductID, &item.MovementType, &item.QuantityDelta, &item.StockBucket, &item.ReferenceKey, &item.Reason, &item.ActorUserID, &item.CreatedAt)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.Movement{}, domain.ErrDuplicate
	}
	if err != nil {
		return domain.Movement{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,reason,new_state) VALUES('INVENTORY_ADJUSTED',$1,'product',$2,$3,jsonb_build_object('quantityDelta',$4::bigint,'referenceKey',$5::text))`, actor, productID, input.Reason, input.QuantityDelta, input.ReferenceKey); err != nil {
		return domain.Movement{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES('INVENTORY_ADJUSTED','product',$1,jsonb_build_object('productId',$1::text,'quantityDelta',$2::bigint,'referenceKey',$3::text))`, productID, input.QuantityDelta, input.ReferenceKey); err != nil {
		return domain.Movement{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Movement{}, err
	}
	return item, nil
}

func ReleaseExpiredReservations(ctx context.Context, tx pgx.Tx, now time.Time) error {
	rows, err := tx.Query(ctx, `SELECT id,product_id,quantity FROM stock_reservations WHERE status='HELD' AND expires_at <= $1 ORDER BY expires_at,id FOR UPDATE SKIP LOCKED`, now)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id, product string
		var quantity int64
		if err = rows.Scan(&id, &product, &quantity); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `UPDATE stock_reservations SET status='RELEASED',released_at=$2 WHERE id=$1`, id, now); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO inventory_movements(product_id,movement_type,quantity_delta,reference_key,reason,stock_bucket) VALUES($1,'RELEASED',$2,$3,'reservation expired','AVAILABLE') ON CONFLICT(reference_key) DO NOTHING`, product, quantity, "release:expiry:"+id); err != nil {
			return err
		}
	}
	return rows.Err()
}

// ReleaseExpired runs the expiry transition in its own short transaction so
// the background worker and checkout can safely race using row locks.
func (r *PostgresRepository) ReleaseExpired(ctx context.Context) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := ReleaseExpiredReservations(ctx, tx, time.Now().UTC()); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
