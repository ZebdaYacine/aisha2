package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aisha-platform/aisha/apps/api/internal/features/order/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(p *pgxpool.Pool) *PostgresRepository { return &PostgresRepository{pool: p} }

type addressRow struct{ fullName, phone, line1, line2, city, postalCode, country string }
type productRow struct {
	id, artisanID, workshopID, status, currency, artisanName, workshopName, productName string
	price                                                                               int64
	made                                                                                bool
	madeToOrder                                                                         bool
	available                                                                           int64
}

func (r *PostgresRepository) Checkout(ctx context.Context, userID, addressID string, items []domain.CartItem, idempotencyKey, requestHash string) (domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var idemID, state, resourceID string
	err = tx.QueryRow(ctx, `INSERT INTO idempotency_keys(actor_user_id,scope,idempotency_key,request_hash,expires_at) VALUES($1,'checkout',$2,$3,CURRENT_TIMESTAMP+interval '30 minutes') ON CONFLICT(actor_user_id,scope,idempotency_key) DO NOTHING RETURNING id,state,COALESCE(resource_id::text,'')`, userID, idempotencyKey, requestHash).Scan(&idemID, &state, &resourceID)
	// The request hash is also the stable key supplied by the service. The header itself is retained in the request hash table key.
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, fmt.Errorf("reserve checkout idempotency key: %w", err)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		var storedHash string
		err = tx.QueryRow(ctx, `SELECT state,request_hash,COALESCE(resource_id::text,'') FROM idempotency_keys WHERE actor_user_id=$1 AND scope='checkout' AND idempotency_key=$2 FOR UPDATE`, userID, idempotencyKey).Scan(&state, &storedHash, &resourceID)
		if err != nil {
			return domain.Order{}, err
		}
		if storedHash != requestHash {
			return domain.Order{}, domain.ErrIdempotencyConflict
		}
		if state == "COMPLETED" && resourceID != "" {
			_ = tx.Rollback(ctx)
			return r.GetMine(ctx, userID, resourceID)
		}
		return domain.Order{}, domain.ErrIdempotencyConflict
	}
	if err = releaseExpiredReservations(ctx, tx); err != nil {
		return domain.Order{}, fmt.Errorf("release expired reservations: %w", err)
	}
	var address addressRow
	err = tx.QueryRow(ctx, `SELECT full_name,COALESCE(phone,''),line1,COALESCE(line2,''),city,postal_code,country FROM addresses WHERE id=$1 AND user_id=$2`, addressID, userID).Scan(&address.fullName, &address.phone, &address.line1, &address.line2, &address.city, &address.postalCode, &address.country)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, domain.ErrValidation
	}
	if err != nil {
		return domain.Order{}, err
	}
	snap, _ := json.Marshal(domain.AddressSnapshot{FullName: address.fullName, Phone: address.phone, Line1: address.line1, Line2: address.line2, City: address.city, PostalCode: address.postalCode, Country: address.country})
	orderID := uuid.NewString()
	orderNumber := "AIS-" + strings.ToUpper(uuid.NewString()[:8])
	var currency string
	var subtotal int64
	prepared := make([]productRow, 0, len(items))
	for _, item := range items {
		var v productRow
		err = tx.QueryRow(ctx, `SELECT p.id,p.artisan_profile_id,p.workshop_id,p.status,p.currency,p.price_minor,p.made_to_order_eligible,COALESCE(a.public_display_name,''),COALESCE(w.name,''),COALESCE((SELECT t.name FROM product_translations t WHERE t.product_id=p.id AND t.locale='en' LIMIT 1),(SELECT t.name FROM product_translations t WHERE t.product_id=p.id ORDER BY t.locale LIMIT 1),p.product_type) FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN workshops w ON w.id=p.workshop_id WHERE p.id=$1 FOR UPDATE OF p`, item.ProductID).Scan(&v.id, &v.artisanID, &v.workshopID, &v.status, &v.currency, &v.price, &v.madeToOrder, &v.artisanName, &v.workshopName, &v.productName)
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Order{}, domain.ErrNotFound
		}
		if err != nil {
			return domain.Order{}, err
		}
		if v.status != "ACTIVE" {
			return domain.Order{}, domain.ErrInvalidTransition
		}
		if err = tx.QueryRow(ctx, `SELECT a.status='APPROVED' AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE') AND w.status='ACTIVE' AND p.published_at IS NOT NULL FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN workshops w ON w.id=p.workshop_id WHERE p.id=$1`, v.id).Scan(&v.made); err != nil {
			return domain.Order{}, err
		}
		if !v.made {
			return domain.Order{}, domain.ErrInvalidTransition
		}
		if !v.madeToOrder {
			if err = tx.QueryRow(ctx, `SELECT COALESCE(SUM(quantity_delta),0) FROM inventory_movements WHERE product_id=$1 AND stock_bucket='AVAILABLE'`, v.id).Scan(&v.available); err != nil {
				return domain.Order{}, err
			}
		}
		if item.ExpectedPriceMinor != nil && *item.ExpectedPriceMinor != v.price {
			return domain.Order{}, domain.ErrPriceChanged
		}
		if item.ExpectedCurrency != "" && item.ExpectedCurrency != strings.TrimSpace(v.currency) {
			return domain.Order{}, domain.ErrPriceChanged
		}
		if !v.madeToOrder && v.available < int64(item.Quantity) {
			return domain.Order{}, domain.ErrOutOfStock
		}
		if currency == "" {
			currency = strings.TrimSpace(v.currency)
		}
		if currency != strings.TrimSpace(v.currency) {
			return domain.Order{}, domain.ErrValidation
		}
		prepared = append(prepared, v)
		subtotal += v.price * int64(item.Quantity)
	}
	if currency == "" {
		return domain.Order{}, domain.ErrValidation
	}
	if _, err = tx.Exec(ctx, `INSERT INTO orders(id,order_number,user_id,status,currency,subtotal_minor,shipping_minor,total_minor,address_snapshot) VALUES($1,$2,$3,'PENDING_PAYMENT',$4,$5,0,$5,$6)`, orderID, orderNumber, userID, currency, subtotal, snap); err != nil {
		return domain.Order{}, err
	}
	for i, v := range prepared {
		item := items[i]
		if _, err = tx.Exec(ctx, `INSERT INTO order_items(order_id,product_id,artisan_profile_id,workshop_id,product_name,artisan_name,workshop_name,unit_price_minor,currency,quantity,subtotal_minor) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`, orderID, v.id, v.artisanID, v.workshopID, v.productName, v.artisanName, v.workshopName, v.price, v.currency, item.Quantity, v.price*int64(item.Quantity)); err != nil {
			return domain.Order{}, err
		}
		if !v.madeToOrder {
			var reservationID string
			if err = tx.QueryRow(ctx, `INSERT INTO stock_reservations(order_id,product_id,quantity,expires_at) VALUES($1,$2,$3,CURRENT_TIMESTAMP+interval '30 minutes') RETURNING id`, orderID, v.id, item.Quantity).Scan(&reservationID); err != nil {
				return domain.Order{}, err
			}
			if _, err = tx.Exec(ctx, `INSERT INTO inventory_movements(product_id,movement_type,quantity_delta,reference_key,reason,actor_user_id,stock_bucket) VALUES($1,'RESERVED',$2,$3,'checkout stock reservation',$4,'AVAILABLE')`, v.id, -item.Quantity, "reservation:"+reservationID, userID); err != nil {
				return domain.Order{}, err
			}
		}
	}
	if _, err = tx.Exec(ctx, `INSERT INTO payment_attempts(order_id,status,amount_minor,currency) VALUES($1,'PENDING',$2,$3)`, orderID, subtotal, currency); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO shipment_events(order_id,status) VALUES($1,'PENDING')`, orderID); err != nil {
		return domain.Order{}, err
	}
	payload, _ := json.Marshal(map[string]any{"orderId": orderID, "orderNumber": orderNumber, "totalMinor": subtotal})
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,new_state) VALUES($1,$2,$3,$4,$5)`, `ORDER_CHECKOUT_CREATED`, userID, "order", orderID, payload); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,$2,$3,$4)`, `ORDER_CHECKOUT_CREATED`, `order`, orderID, payload); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE idempotency_keys SET state='COMPLETED',response_status=201,response_body=$2,resource_type='order',resource_id=$3,completed_at=CURRENT_TIMESTAMP WHERE id=$1`, idemID, payload, orderID); err != nil {
		return domain.Order{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}
	return r.GetMine(ctx, userID, orderID)
}

func releaseExpiredReservations(ctx context.Context, tx pgx.Tx) error {
	rows, err := tx.Query(ctx, `SELECT id,product_id,quantity FROM stock_reservations WHERE status='HELD' AND expires_at <= CURRENT_TIMESTAMP ORDER BY expires_at,id FOR UPDATE SKIP LOCKED`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id, productID string
		var quantity int64
		if err = rows.Scan(&id, &productID, &quantity); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `UPDATE stock_reservations SET status='RELEASED',released_at=CURRENT_TIMESTAMP WHERE id=$1`, id); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO inventory_movements(product_id,movement_type,quantity_delta,reference_key,reason,stock_bucket) VALUES($1,'RELEASED',$2,$3,'reservation expired','AVAILABLE') ON CONFLICT(reference_key) DO NOTHING`, productID, quantity, "release:expiry:"+id); err != nil {
			return err
		}
	}
	return rows.Err()
}

func (r *PostgresRepository) ListMine(ctx context.Context, userID string, limit, offset int) ([]domain.Order, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE user_id=$1`, userID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, `SELECT id FROM orders WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}
	out := make([]domain.Order, 0, len(ids))
	for _, id := range ids {
		o, e := r.GetMine(ctx, userID, id)
		if e != nil {
			return nil, 0, e
		}
		out = append(out, o)
	}
	return out, total, rows.Err()
}
func (r *PostgresRepository) GetMine(ctx context.Context, userID, id string) (domain.Order, error) {
	return r.load(ctx, r.pool, userID, id)
}
func (r *PostgresRepository) load(ctx context.Context, q interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
}, userID, id string) (domain.Order, error) {
	var o domain.Order
	var address []byte
	var created, updated time.Time
	err := q.QueryRow(ctx, `SELECT id,order_number,user_id,status,currency,subtotal_minor,shipping_minor,total_minor,address_snapshot,created_at,updated_at FROM orders WHERE id=$1 AND user_id=$2`, id, userID).Scan(&o.ID, &o.OrderNumber, &o.UserID, &o.Status, &o.Currency, &o.SubtotalMinor, &o.ShippingMinor, &o.TotalMinor, &address, &created, &updated)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	o.CreatedAt = created
	o.UpdatedAt = updated
	_ = json.Unmarshal(address, &o.Address)
	rows, err := q.Query(ctx, `SELECT id,product_id,product_name,artisan_profile_id,artisan_name,workshop_id,workshop_name,unit_price_minor,currency,quantity,subtotal_minor FROM order_items WHERE order_id=$1 ORDER BY created_at`, o.ID)
	if err != nil {
		return domain.Order{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var i domain.OrderItem
		if err = rows.Scan(&i.ID, &i.ProductID, &i.ProductName, &i.ArtisanID, &i.ArtisanName, &i.WorkshopID, &i.WorkshopName, &i.UnitPriceMinor, &i.Currency, &i.Quantity, &i.SubtotalMinor); err != nil {
			return domain.Order{}, err
		}
		o.Items = append(o.Items, i)
	}
	var payment domain.PaymentAttempt
	err = q.QueryRow(ctx, `SELECT id,provider,COALESCE(provider_reference,''),status,amount_minor,currency,COALESCE(failure_reason,'') FROM payment_attempts WHERE order_id=$1 ORDER BY created_at DESC LIMIT 1`, o.ID).Scan(&payment.ID, &payment.Provider, &payment.ProviderReference, &payment.Status, &payment.AmountMinor, &payment.Currency, &payment.FailureReason)
	if err == nil {
		o.Payment = &payment
	}
	shipmentRows, err := q.Query(ctx, `SELECT id,status,COALESCE(tracking_reference,''),occurred_at FROM shipment_events WHERE order_id=$1 ORDER BY occurred_at,id`, o.ID)
	if err != nil {
		return domain.Order{}, err
	}
	defer shipmentRows.Close()
	for shipmentRows.Next() {
		var event domain.ShipmentEvent
		if err = shipmentRows.Scan(&event.ID, &event.Status, &event.TrackingReference, &event.OccurredAt); err != nil {
			return domain.Order{}, err
		}
		o.ShipmentEvents = append(o.ShipmentEvents, event)
	}
	if err = shipmentRows.Err(); err != nil {
		return domain.Order{}, err
	}
	return o, rows.Err()
}

func (r *PostgresRepository) ListSeller(ctx context.Context, userID string, limit, offset int) ([]domain.SellerItem, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM order_items oi JOIN artisan_profiles a ON a.id=oi.artisan_profile_id WHERE a.user_id=$1`, userID).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, `SELECT o.id,o.order_number,o.status,oi.product_id,oi.product_name,oi.quantity,oi.subtotal_minor,oi.currency,o.created_at FROM order_items oi JOIN orders o ON o.id=oi.order_id JOIN artisan_profiles a ON a.id=oi.artisan_profile_id WHERE a.user_id=$1 ORDER BY o.created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []domain.SellerItem{}
	for rows.Next() {
		var v domain.SellerItem
		if err = rows.Scan(&v.OrderID, &v.OrderNumber, &v.OrderStatus, &v.ProductID, &v.ProductName, &v.Quantity, &v.SubtotalMinor, &v.Currency, &v.CreatedAt); err != nil {
			return nil, 0, err
		}
		out = append(out, v)
	}
	return out, total, rows.Err()
}

func (r *PostgresRepository) Cancel(ctx context.Context, userID, id string) (domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var status string
	if err = tx.QueryRow(ctx, `SELECT status FROM orders WHERE id=$1 AND user_id=$2 FOR UPDATE`, id, userID).Scan(&status); errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	if status == "CANCELLED" {
		_ = tx.Rollback(ctx)
		return r.GetMine(ctx, userID, id)
	}
	if status != "PENDING_PAYMENT" && status != "PAID" {
		return domain.Order{}, domain.ErrInvalidTransition
	}
	rows, err := tx.Query(ctx, `SELECT id,product_id,quantity FROM stock_reservations WHERE order_id=$1 AND status='HELD' FOR UPDATE`, id)
	if err != nil {
		return domain.Order{}, err
	}
	type res struct {
		id, product string
		qty         int
	}
	var reservations []res
	for rows.Next() {
		var x res
		if err = rows.Scan(&x.id, &x.product, &x.qty); err != nil {
			rows.Close()
			return domain.Order{}, err
		}
		reservations = append(reservations, x)
	}
	rows.Close()
	for _, x := range reservations {
		if _, err = tx.Exec(ctx, `UPDATE stock_reservations SET status='RELEASED',released_at=CURRENT_TIMESTAMP WHERE id=$1`, x.id); err != nil {
			return domain.Order{}, err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO inventory_movements(product_id,movement_type,quantity_delta,reference_key,reason,actor_user_id,stock_bucket) VALUES($1,'RELEASED',$2,$3,'order cancellation',$4,'AVAILABLE')`, x.product, x.qty, "release:"+x.id, userID); err != nil {
			return domain.Order{}, err
		}
	}
	if _, err = tx.Exec(ctx, `UPDATE orders SET status='CANCELLED',cancelled_at=CURRENT_TIMESTAMP,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE payment_attempts SET status=CASE WHEN status='CONFIRMED' THEN 'REFUND_PENDING' ELSE 'CANCELLED' END WHERE order_id=$1`, id); err != nil {
		return domain.Order{}, err
	}
	payload, _ := json.Marshal(map[string]any{"orderId": id, "status": "CANCELLED"})
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,new_state) VALUES($1,$2,$3,$4,$5)`, `ORDER_CANCELLED`, userID, "order", id, payload); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,$2,$3,$4)`, `ORDER_CANCELLED`, `order`, id, payload); err != nil {
		return domain.Order{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}
	return r.GetMine(ctx, userID, id)
}

func (r *PostgresRepository) RecordReturn(ctx context.Context, actor, id, reason string) (domain.Return, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Return{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var status string
	if err = tx.QueryRow(ctx, `SELECT status FROM orders WHERE id=$1 FOR UPDATE`, id).Scan(&status); errors.Is(err, pgx.ErrNoRows) {
		return domain.Return{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Return{}, err
	}
	if status != "SHIPPED" && status != "DELIVERED" {
		return domain.Return{}, domain.ErrInvalidTransition
	}
	var out domain.Return
	err = tx.QueryRow(ctx, `INSERT INTO order_returns(order_id,actor_user_id,reason) VALUES($1,$2,$3) ON CONFLICT(order_id) DO NOTHING RETURNING id,order_id,status,reason,created_at`, id, actor, reason).Scan(&out.ID, &out.OrderID, &out.Status, &out.Reason, &out.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		if err = tx.QueryRow(ctx, `SELECT id,order_id,status,reason,created_at FROM order_returns WHERE order_id=$1`, id).Scan(&out.ID, &out.OrderID, &out.Status, &out.Reason, &out.CreatedAt); err != nil {
			return domain.Return{}, err
		}
		_ = tx.Rollback(ctx)
		return out, nil
	}
	if err != nil {
		return domain.Return{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO inventory_movements(product_id,movement_type,quantity_delta,reference_key,reason,actor_user_id,stock_bucket)
        SELECT product_id,'RETURN',quantity,'return:' || $1::text || ':' || id::text,'returned stock quarantined',$2,'QUARANTINED'
        FROM order_items WHERE order_id=$3
        ON CONFLICT(reference_key) DO NOTHING`, out.ID, actor, id); err != nil {
		return domain.Return{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE orders SET status='RETURNED',updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id); err != nil {
		return domain.Return{}, err
	}
	payload, _ := json.Marshal(map[string]any{"orderId": id, "returnId": out.ID, "status": out.Status})
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,new_state) VALUES($1,$2,$3,$4,$5)`, `ORDER_RETURN_RECORDED`, actor, "order", id, payload); err != nil {
		return domain.Return{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,$2,$3,$4)`, `ORDER_RETURN_RECORDED`, `order`, id, payload); err != nil {
		return domain.Return{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Return{}, err
	}
	return out, nil
}

// GetPayment returns payment metadata only for an order owned by the caller.
func (r *PostgresRepository) GetPayment(ctx context.Context, userID, paymentID string) (domain.PaymentAttempt, error) {
	var payment domain.PaymentAttempt
	err := r.pool.QueryRow(ctx, `
		SELECT p.id,p.provider,COALESCE(p.provider_reference,''),p.status,p.amount_minor,p.currency,COALESCE(p.failure_reason,'')
		FROM payment_attempts p JOIN orders o ON o.id=p.order_id
		WHERE p.id=$1 AND o.user_id=$2`, paymentID, userID).
		Scan(&payment.ID, &payment.Provider, &payment.ProviderReference, &payment.Status, &payment.AmountMinor, &payment.Currency, &payment.FailureReason)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.PaymentAttempt{}, domain.ErrNotFound
	}
	return payment, err
}

// ConfirmPayment is the internal manual development adapter. It trusts only
// the server-side order total, commits held reservations, and is idempotent
// when the customer clicks the payment action more than once.
func (r *PostgresRepository) ConfirmPayment(ctx context.Context, userID, paymentID, _ string) (domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var orderID, paymentStatus, orderStatus, paymentCurrency, orderCurrency string
	var paymentAmount, orderTotal int64
	err = tx.QueryRow(ctx, `
		SELECT p.order_id,p.status,p.amount_minor,p.currency,o.status,o.total_minor,o.currency
		FROM payment_attempts p JOIN orders o ON o.id=p.order_id
		WHERE p.id=$1 AND o.user_id=$2
		FOR UPDATE OF p,o`, paymentID, userID).
		Scan(&orderID, &paymentStatus, &paymentAmount, &paymentCurrency, &orderStatus, &orderTotal, &orderCurrency)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	if paymentStatus == "CONFIRMED" && orderStatus == "PAID" {
		_ = tx.Rollback(ctx)
		return r.GetMine(ctx, userID, orderID)
	}
	if paymentStatus != "PENDING" || orderStatus != "PENDING_PAYMENT" {
		return domain.Order{}, domain.ErrInvalidTransition
	}
	if paymentAmount != orderTotal || strings.TrimSpace(paymentCurrency) != strings.TrimSpace(orderCurrency) {
		return domain.Order{}, domain.ErrPaymentAmountMismatch
	}
	providerReference := "manual:" + paymentID
	if _, err = tx.Exec(ctx, `UPDATE payment_attempts SET provider='manual',provider_reference=$2,status='CONFIRMED',failure_reason=NULL WHERE id=$1`, paymentID, providerReference); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE stock_reservations SET status='COMMITTED' WHERE order_id=$1 AND status='HELD'`, orderID); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE orders SET status='PAID',updated_at=CURRENT_TIMESTAMP WHERE id=$1`, orderID); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO shipment_events(order_id,status) VALUES($1,'PAID') ON CONFLICT DO NOTHING`, orderID); err != nil {
		return domain.Order{}, err
	}
	payload, _ := json.Marshal(map[string]any{"orderId": orderID, "paymentId": paymentID, "provider": "manual", "status": "CONFIRMED"})
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,new_state) VALUES($1,$2,$3,$4,$5)`, "PAYMENT_CONFIRMED", userID, "payment", paymentID, payload); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,$2,$3,$4)`, "PAYMENT_CONFIRMED", "order", orderID, payload); err != nil {
		return domain.Order{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}
	return r.GetMine(ctx, userID, orderID)
}

func (r *PostgresRepository) FailPayment(ctx context.Context, userID, paymentID, reason string) (domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var orderID, paymentStatus, orderStatus string
	err = tx.QueryRow(ctx, `SELECT p.order_id,p.status,o.status FROM payment_attempts p JOIN orders o ON o.id=p.order_id WHERE p.id=$1 AND o.user_id=$2 FOR UPDATE OF p,o`, paymentID, userID).Scan(&orderID, &paymentStatus, &orderStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	if paymentStatus == "FAILED" && orderStatus == "PAYMENT_FAILED" {
		_ = tx.Rollback(ctx)
		return r.GetMine(ctx, userID, orderID)
	}
	if paymentStatus != "PENDING" || orderStatus != "PENDING_PAYMENT" {
		return domain.Order{}, domain.ErrInvalidTransition
	}
	rows, err := tx.Query(ctx, `SELECT id,product_id,quantity FROM stock_reservations WHERE order_id=$1 AND status='HELD' FOR UPDATE`, orderID)
	if err != nil {
		return domain.Order{}, err
	}
	type reservation struct {
		id, productID string
		quantity      int
	}
	var reservations []reservation
	for rows.Next() {
		var item reservation
		if err = rows.Scan(&item.id, &item.productID, &item.quantity); err != nil {
			rows.Close()
			return domain.Order{}, err
		}
		reservations = append(reservations, item)
	}
	rows.Close()
	for _, item := range reservations {
		if _, err = tx.Exec(ctx, `UPDATE stock_reservations SET status='RELEASED',released_at=CURRENT_TIMESTAMP WHERE id=$1`, item.id); err != nil {
			return domain.Order{}, err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO inventory_movements(product_id,movement_type,quantity_delta,reference_key,reason,actor_user_id,stock_bucket) VALUES($1,'RELEASED',$2,$3,'payment failed',$4,'AVAILABLE') ON CONFLICT(reference_key) DO NOTHING`, item.productID, item.quantity, "release:payment:"+item.id, userID); err != nil {
			return domain.Order{}, err
		}
	}
	if _, err = tx.Exec(ctx, `UPDATE payment_attempts SET status='FAILED',failure_reason=$2 WHERE id=$1`, paymentID, strings.TrimSpace(reason)); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE orders SET status='PAYMENT_FAILED',updated_at=CURRENT_TIMESTAMP WHERE id=$1`, orderID); err != nil {
		return domain.Order{}, err
	}
	payload, _ := json.Marshal(map[string]any{"orderId": orderID, "paymentId": paymentID, "status": "FAILED", "reason": strings.TrimSpace(reason)})
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,new_state) VALUES($1,$2,$3,$4,$5)`, "PAYMENT_FAILED", userID, "payment", paymentID, payload); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,$2,$3,$4)`, "PAYMENT_FAILED", "order", orderID, payload); err != nil {
		return domain.Order{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}
	return r.GetMine(ctx, userID, orderID)
}

func (r *PostgresRepository) loadByID(ctx context.Context, id string) (domain.Order, error) {
	var userID string
	if err := r.pool.QueryRow(ctx, `SELECT user_id FROM orders WHERE id=$1`, id).Scan(&userID); errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, domain.ErrNotFound
	} else if err != nil {
		return domain.Order{}, err
	}
	return r.load(ctx, r.pool, userID, id)
}

func (r *PostgresRepository) ListFulfilment(ctx context.Context, limit, offset int) ([]domain.FulfilmentOrder, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE status IN ('PAID','PREPARING','READY_TO_SHIP')`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.pool.Query(ctx, `SELECT o.id,o.order_number,o.status,o.currency,o.total_minor,COALESCE(u.email,''),o.created_at FROM orders o JOIN users u ON u.id=o.user_id WHERE o.status IN ('PAID','PREPARING','READY_TO_SHIP') ORDER BY o.created_at LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]domain.FulfilmentOrder, 0)
	for rows.Next() {
		var item domain.FulfilmentOrder
		if err = rows.Scan(&item.ID, &item.OrderNumber, &item.Status, &item.Currency, &item.TotalMinor, &item.CustomerEmail, &item.CreatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) Prepare(ctx context.Context, actor, id string) (domain.Order, error) {
	return r.transitionFulfilment(ctx, actor, id, "PAID", "PREPARING", "PREPARING", "")
}

func (r *PostgresRepository) Ship(ctx context.Context, actor, id string, input domain.ShipmentInput) (domain.Order, error) {
	return r.transitionFulfilment(ctx, actor, id, "PREPARING", "SHIPPED", "SHIPPED", strings.TrimSpace(input.Carrier)+" / "+strings.TrimSpace(input.TrackingReference))
}

func (r *PostgresRepository) Deliver(ctx context.Context, actor, id string) (domain.Order, error) {
	return r.transitionFulfilment(ctx, actor, id, "SHIPPED", "DELIVERED", "DELIVERED", "")
}

func (r *PostgresRepository) transitionFulfilment(ctx context.Context, actor, id, from, to, eventStatus, tracking string) (domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var ownerID, status string
	if err = tx.QueryRow(ctx, `SELECT user_id,status FROM orders WHERE id=$1 FOR UPDATE`, id).Scan(&ownerID, &status); errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	if status != from {
		return domain.Order{}, domain.ErrInvalidTransition
	}
	if _, err = tx.Exec(ctx, `UPDATE orders SET status=$2,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id, to); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO shipment_events(order_id,status,tracking_reference,actor_user_id) VALUES($1,$2,NULLIF($3,''),$4) ON CONFLICT DO NOTHING`, id, eventStatus, tracking, actor); err != nil {
		return domain.Order{}, err
	}
	payload, _ := json.Marshal(map[string]any{"orderId": id, "status": to})
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,new_state) VALUES($1,$2,$3,$4,$5)`, "ORDER_"+to, actor, "order", id, payload); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,$2,$3,$4)`, "ORDER_"+to, "order", id, payload); err != nil {
		return domain.Order{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}
	return r.load(ctx, r.pool, ownerID, id)
}

func (r *PostgresRepository) Refund(ctx context.Context, actor, id string, amount int64, key, reason string) (domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Order{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var ownerID, orderStatus, paymentID, paymentStatus, currency string
	var total, paymentAmount int64
	err = tx.QueryRow(ctx, `SELECT o.user_id,o.status,p.id,p.status,p.amount_minor,p.currency,o.total_minor FROM orders o JOIN payment_attempts p ON p.order_id=o.id WHERE o.id=$1 ORDER BY p.created_at DESC LIMIT 1 FOR UPDATE OF o,p`, id).Scan(&ownerID, &orderStatus, &paymentID, &paymentStatus, &paymentAmount, &currency, &total)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	if orderStatus != "PAID" && orderStatus != "PARTIALLY_REFUNDED" || paymentStatus != "CONFIRMED" && paymentStatus != "REFUND_PENDING" {
		return domain.Order{}, domain.ErrInvalidTransition
	}
	var refunded int64
	if err = tx.QueryRow(ctx, `SELECT COALESCE(SUM(amount_minor),0) FROM payment_refunds WHERE payment_attempt_id=$1`, paymentID).Scan(&refunded); err != nil {
		return domain.Order{}, err
	}
	if amount <= 0 || amount > paymentAmount-refunded {
		return domain.Order{}, domain.ErrPaymentAmountMismatch
	}
	var refundID string
	err = tx.QueryRow(ctx, `INSERT INTO payment_refunds(payment_attempt_id,order_id,actor_user_id,idempotency_key,amount_minor,currency,reason) VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT(actor_user_id,idempotency_key) DO NOTHING RETURNING id`, paymentID, id, actor, key, amount, currency, reason).Scan(&refundID)
	if errors.Is(err, pgx.ErrNoRows) {
		_ = tx.Rollback(ctx)
		return r.loadByID(ctx, id)
	}
	if err != nil {
		return domain.Order{}, err
	}
	newStatus := "PARTIALLY_REFUNDED"
	if amount+refunded == paymentAmount {
		newStatus = "REFUNDED"
	}
	if _, err = tx.Exec(ctx, `UPDATE payment_attempts SET status=$2 WHERE id=$1`, paymentID, map[string]string{"PARTIALLY_REFUNDED": "REFUND_PENDING", "REFUNDED": "REFUNDED"}[newStatus]); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `UPDATE orders SET status=$2,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id, newStatus); err != nil {
		return domain.Order{}, err
	}
	payload, _ := json.Marshal(map[string]any{"orderId": id, "paymentId": paymentID, "refundId": refundID, "amountMinor": amount, "status": newStatus})
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,new_state) VALUES($1,$2,$3,$4,$5)`, "PAYMENT_REFUNDED", actor, "payment", paymentID, payload); err != nil {
		return domain.Order{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES($1,$2,$3,$4)`, "PAYMENT_REFUNDED", "order", id, payload); err != nil {
		return domain.Order{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Order{}, err
	}
	return r.load(ctx, r.pool, ownerID, id)
}
