package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/notification/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

type recipient struct {
	id    string
	email string
	name  string
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) List(ctx context.Context, userID string, limit, offset int) ([]domain.Notification, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE recipient_user_id=$1`, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count notifications: %w", err)
	}
	rows, err := r.pool.Query(ctx, `SELECT id,event_type,title_key,body_key,payload,read_at,created_at FROM notifications WHERE recipient_user_id=$1 ORDER BY created_at DESC,id DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()
	items := make([]domain.Notification, 0)
	for rows.Next() {
		var item domain.Notification
		if err := rows.Scan(&item.ID, &item.EventType, &item.TitleKey, &item.BodyKey, &item.Payload, &item.ReadAt, &item.CreatedAt); err != nil {
			return nil, 0, err
		}
		item.RecipientUserID = userID
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) UnreadCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM notifications WHERE recipient_user_id=$1 AND read_at IS NULL`, userID).Scan(&count)
	return count, err
}

func (r *PostgresRepository) MarkRead(ctx context.Context, userID, id string) error {
	result, err := r.pool.Exec(ctx, `UPDATE notifications SET read_at=COALESCE(read_at,CURRENT_TIMESTAMP) WHERE id=$1 AND recipient_user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET read_at=CURRENT_TIMESTAMP WHERE recipient_user_id=$1 AND read_at IS NULL`, userID)
	return err
}

func (r *PostgresRepository) ClaimOutbox(ctx context.Context, limit int, lease time.Duration) ([]domain.OutboxEvent, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	rows, err := tx.Query(ctx, `
		SELECT id,event_type,aggregate_type,aggregate_id,payload,attempt_count
		FROM outbox_events
		WHERE processed_at IS NULL
		  AND available_at <= CURRENT_TIMESTAMP
		  AND (locked_at IS NULL OR locked_at < CURRENT_TIMESTAMP - $1::interval)
		ORDER BY created_at,id
		FOR UPDATE SKIP LOCKED
		LIMIT $2`, intervalLiteral(lease), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := make([]domain.OutboxEvent, 0)
	for rows.Next() {
		var event domain.OutboxEvent
		if err := rows.Scan(&event.ID, &event.EventType, &event.AggregateType, &event.AggregateID, &event.Payload, &event.AttemptCount); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for _, event := range events {
		if _, err := tx.Exec(ctx, `UPDATE outbox_events SET locked_at=CURRENT_TIMESTAMP,attempt_count=attempt_count+1 WHERE id=$1`, event.ID); err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	for i := range events {
		events[i].AttemptCount++
	}
	return events, nil
}

func (r *PostgresRepository) DeliverOutbox(ctx context.Context, event domain.OutboxEvent) ([]domain.Notification, error) {
	titleKey, bodyKey, ok := templateFor(event.EventType)
	if !ok {
		return []domain.Notification{}, nil
	}
	recipients, err := r.recipients(ctx, event)
	if err != nil {
		return nil, err
	}
	if len(recipients) == 0 {
		return []domain.Notification{}, nil
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	items := make([]domain.Notification, 0, len(recipients))
	for _, recipient := range recipients {
		var item domain.Notification
		err := tx.QueryRow(ctx, `
			INSERT INTO notifications(recipient_user_id,event_type,title_key,body_key,payload,dedupe_key)
			VALUES($1,$2,$3,$4,$5,$6)
			ON CONFLICT (recipient_user_id,dedupe_key) DO UPDATE SET payload=EXCLUDED.payload
			RETURNING id,event_type,title_key,body_key,payload,read_at,created_at`, recipient.id, event.EventType, titleKey, bodyKey, event.Payload, event.ID).Scan(&item.ID, &item.EventType, &item.TitleKey, &item.BodyKey, &item.Payload, &item.ReadAt, &item.CreatedAt)
		if err == pgx.ErrNoRows {
			continue
		}
		if err != nil {
			return nil, err
		}
		item.RecipientUserID = recipient.id
		item.RecipientEmail = recipient.email
		item.RecipientName = recipient.name
		items = append(items, item)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return items, nil
}

// ClaimEmailDelivery makes email sending idempotent per notification. A stale
// SENDING claim is recoverable after a process crash, while SENT and FAILED
// notifications are never sent again by the outbox worker.
func (r *PostgresRepository) ClaimEmailDelivery(ctx context.Context, notificationID string) (bool, error) {
	var claimedID string
	err := r.pool.QueryRow(ctx, `
		UPDATE notifications
		SET email_delivery_status='SENDING',
		    email_delivery_attempted_at=CURRENT_TIMESTAMP,
		    email_delivery_error=NULL
		WHERE id=$1
		  AND (
				email_delivery_status='PENDING'
				OR (email_delivery_status='SENDING'
				    AND COALESCE(email_delivery_attempted_at, TIMESTAMP 'epoch') < CURRENT_TIMESTAMP - INTERVAL '15 minutes')
		  )
		RETURNING id`, notificationID).Scan(&claimedID)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("claim notification email delivery: %w", err)
	}
	return claimedID != "", nil
}

func (r *PostgresRepository) MarkEmailDeliverySent(ctx context.Context, notificationID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE notifications
		SET email_delivery_status='SENT', email_delivery_error=NULL
		WHERE id=$1 AND email_delivery_status='SENDING'`, notificationID)
	return err
}

func (r *PostgresRepository) MarkEmailDeliveryFailed(ctx context.Context, notificationID string, cause error) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE notifications
		SET email_delivery_status='FAILED', email_delivery_error=$2
		WHERE id=$1 AND email_delivery_status='SENDING'`, notificationID, truncateError(cause))
	return err
}

func (r *PostgresRepository) MarkOutboxProcessed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE outbox_events SET processed_at=CURRENT_TIMESTAMP,locked_at=NULL,last_error=NULL WHERE id=$1`, id)
	return err
}

func (r *PostgresRepository) MarkOutboxFailed(ctx context.Context, id string, availableAt time.Time, cause error) error {
	_, err := r.pool.Exec(ctx, `UPDATE outbox_events SET available_at=$2,locked_at=NULL,last_error=$3 WHERE id=$1`, id, availableAt, truncateError(cause))
	return err
}

func (r *PostgresRepository) recipients(ctx context.Context, event domain.OutboxEvent) ([]recipient, error) {
	items, err := r.aggregateRecipients(ctx, event)
	if err != nil {
		return nil, err
	}
	roles := audienceRoles(event.EventType)
	if len(roles) == 0 {
		return items, nil
	}
	peers, err := r.roleRecipients(ctx, roles...)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(items)+len(peers))
	for _, item := range items {
		seen[item.id] = struct{}{}
	}
	for _, peer := range peers {
		if _, ok := seen[peer.id]; ok {
			continue
		}
		seen[peer.id] = struct{}{}
		items = append(items, peer)
	}
	return items, nil
}

func (r *PostgresRepository) aggregateRecipients(ctx context.Context, event domain.OutboxEvent) ([]recipient, error) {
	var query string
	switch event.AggregateType {
	case "user":
		query = `SELECT id::text,COALESCE(email,''),COALESCE(display_name,'') FROM users WHERE id=$1`
	case "artisan_profile":
		query = `SELECT u.id::text,COALESCE(u.email,''),COALESCE(u.display_name,'') FROM artisan_profiles a JOIN users u ON u.id=a.user_id WHERE a.id=$1`
	case "artisan_membership":
		query = `SELECT u.id::text,COALESCE(u.email,''),COALESCE(u.display_name,'') FROM artisan_memberships a JOIN users u ON u.id=a.user_id WHERE a.id=$1`
	case "workshop":
		query = `SELECT u.id::text,COALESCE(u.email,''),COALESCE(u.display_name,'') FROM workshops w JOIN artisan_profiles a ON a.id=w.artisan_profile_id JOIN users u ON u.id=a.user_id WHERE w.id=$1`
	case "product":
		query = `SELECT u.id::text,COALESCE(u.email,''),COALESCE(u.display_name,'') FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN users u ON u.id=a.user_id WHERE p.id=$1`
	case "order":
		query = `
			SELECT u.id::text,COALESCE(u.email,''),COALESCE(u.display_name,'')
			FROM orders o JOIN users u ON u.id=o.user_id WHERE o.id=$1
			UNION
			SELECT u.id::text,COALESCE(u.email,''),COALESCE(u.display_name,'')
			FROM order_items oi
			JOIN products p ON p.id=oi.product_id
			JOIN artisan_profiles a ON a.id=p.artisan_profile_id
			JOIN users u ON u.id=a.user_id
			WHERE oi.order_id=$1`
	case "warehouse_reception":
		query = `SELECT u.id::text,COALESCE(u.email,''),COALESCE(u.display_name,'') FROM warehouse_receptions r JOIN products p ON p.id=r.product_id JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN users u ON u.id=a.user_id WHERE r.id=$1`
	default:
		return nil, nil
	}
	rows, err := r.pool.Query(ctx, query, event.AggregateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]recipient, 0, 1)
	for rows.Next() {
		var item recipient
		if err := rows.Scan(&item.id, &item.email, &item.name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) roleRecipients(ctx context.Context, roles ...string) ([]recipient, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT u.id::text,COALESCE(u.email,''),COALESCE(u.display_name,'')
		FROM users u
		JOIN user_roles ur ON ur.user_id=u.id
		JOIN roles r ON r.id=ur.role_id
		WHERE u.status='ACTIVE' AND r.code = ANY($1::text[])
		ORDER BY u.id::text`, roles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]recipient, 0)
	for rows.Next() {
		var item recipient
		if err := rows.Scan(&item.id, &item.email, &item.name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func audienceRoles(eventType string) []string {
	switch eventType {
	case "ARTISAN_APPLICATION_SUBMITTED", "PRODUCT_SUBMITTED":
		return []string{"moderator", "administrator"}
	case "WAREHOUSE_RECEPTION_CREATED", "WAREHOUSE_RECEPTION_INSPECTED", "INVENTORY_ACCEPTED":
		return []string{"warehouse_agent", "administrator"}
	case "ORDER_CHECKOUT_CREATED", "ORDER_RETURN_RECORDED":
		return []string{"warehouse_agent", "administrator"}
	case "PAYMENT_CONFIRMED":
		return []string{"warehouse_agent", "administrator"}
	default:
		return nil
	}
}

func templateFor(eventType string) (string, string, bool) {
	if title, body, ok := specificTemplateFor(eventType); ok {
		return title, body, true
	}
	if strings.HasPrefix(eventType, "ARTISAN_APPLICATION_") {
		return "notifications.artisanApplication.title", "notifications.artisanApplication.body", true
	}
	if strings.HasPrefix(eventType, "ARTISAN_VERIFICATION_") {
		return "notifications.verification.title", "notifications.verification.body", true
	}
	if strings.HasPrefix(eventType, "ARTISAN_MEMBERSHIP_") {
		return "notifications.membership.title", "notifications.membership.body", true
	}
	if strings.HasPrefix(eventType, "WORKSHOP_") {
		return "notifications.workshop.title", "notifications.workshop.body", true
	}
	if strings.HasPrefix(eventType, "PRODUCT_") {
		return "notifications.product.title", "notifications.product.body", true
	}
	if strings.HasPrefix(eventType, "ORDER_") {
		return "notifications.order.title", "notifications.order.body", true
	}
	if strings.HasPrefix(eventType, "PAYMENT_") {
		return "notifications.payment.title", "notifications.payment.body", true
	}
	switch eventType {
	case "USER_STATUS_CHANGED", "USER_ROLES_CHANGED":
		return "notifications.account.title", "notifications.account.body", true
	case "INVENTORY_ADJUSTED":
		return "notifications.inventory.title", "notifications.inventory.body", true
	case "WAREHOUSE_RECEPTION_CREATED", "WAREHOUSE_RECEPTION_INSPECTED", "INVENTORY_ACCEPTED", "RECEIVED_PENDING_INSPECTION", "INSPECTED":
		return "notifications.warehouse.title", "notifications.warehouse.body", true
	default:
		return "", "", false
	}
}

func specificTemplateFor(eventType string) (string, string, bool) {
	const (
		artisanApplicationBody = "notifications.artisanApplication.body"
		verificationBody       = "notifications.verification.body"
		membershipBody         = "notifications.membership.body"
		workshopBody           = "notifications.workshop.body"
		productBody            = "notifications.product.body"
		orderBody              = "notifications.order.body"
		paymentBody            = "notifications.payment.body"
		warehouseBody          = "notifications.warehouse.body"
		inventoryBody          = "notifications.inventory.body"
		accountBody            = "notifications.account.body"
	)
	switch eventType {
	case "ARTISAN_APPLICATION_SUBMITTED":
		return "notifications.artisanApplication.submitted.title", artisanApplicationBody, true
	case "ARTISAN_APPLICATION_APPROVED":
		return "notifications.artisanApplication.approved.title", artisanApplicationBody, true
	case "ARTISAN_APPLICATION_CHANGES_REQUESTED":
		return "notifications.artisanApplication.changesRequested.title", artisanApplicationBody, true
	case "ARTISAN_APPLICATION_REJECTED":
		return "notifications.artisanApplication.rejected.title", artisanApplicationBody, true
	case "ARTISAN_VERIFICATION_SUBMITTED":
		return "notifications.verification.submitted.title", verificationBody, true
	case "ARTISAN_VERIFICATION_VERIFIED":
		return "notifications.verification.verified.title", verificationBody, true
	case "ARTISAN_VERIFICATION_CHANGES_REQUESTED":
		return "notifications.verification.changesRequested.title", verificationBody, true
	case "ARTISAN_VERIFICATION_REJECTED":
		return "notifications.verification.rejected.title", verificationBody, true
	case "ARTISAN_MEMBERSHIP_ACTIVE", "ARTISAN_MEMBERSHIP_ACTIVATED":
		return "notifications.membership.active.title", membershipBody, true
	case "ARTISAN_MEMBERSHIP_SUSPENDED":
		return "notifications.membership.suspended.title", membershipBody, true
	case "ARTISAN_MEMBERSHIP_CLOSED":
		return "notifications.membership.closed.title", membershipBody, true
	case "WORKSHOP_CREATED":
		return "notifications.workshop.created.title", workshopBody, true
	case "WORKSHOP_ADMIN_STATUS_CHANGED":
		return "notifications.workshop.statusChanged.title", workshopBody, true
	case "WORKSHOP_DELETED":
		return "notifications.workshop.deleted.title", workshopBody, true
	case "PRODUCT_DRAFT_CREATED":
		return "notifications.product.draftCreated.title", productBody, true
	case "PRODUCT_SUBMITTED":
		return "notifications.product.submitted.title", productBody, true
	case "PRODUCT_MODERATION_DECIDED":
		return "notifications.product.moderation.title", productBody, true
	case "PRODUCT_AUTO_ACTIVATED_AFTER_INSPECTION":
		return "notifications.product.autoActivated.title", productBody, true
	case "PRODUCT_ARCHIVED":
		return "notifications.product.archived.title", productBody, true
	case "ORDER_CHECKOUT_CREATED":
		return "notifications.order.checkout.title", orderBody, true
	case "ORDER_CANCELLED":
		return "notifications.order.cancelled.title", orderBody, true
	case "ORDER_RETURN_RECORDED":
		return "notifications.order.returnRecorded.title", orderBody, true
	case "PAYMENT_CONFIRMED":
		return "notifications.payment.confirmed.title", paymentBody, true
	case "PAYMENT_FAILED":
		return "notifications.payment.failed.title", paymentBody, true
	case "PAYMENT_REFUNDED":
		return "notifications.payment.refunded.title", paymentBody, true
	case "WAREHOUSE_RECEPTION_CREATED", "RECEIVED_PENDING_INSPECTION":
		return "notifications.warehouse.received.title", warehouseBody, true
	case "WAREHOUSE_RECEPTION_INSPECTED", "INSPECTED":
		return "notifications.warehouse.inspected.title", warehouseBody, true
	case "INVENTORY_ACCEPTED":
		return "notifications.warehouse.inventoryAccepted.title", warehouseBody, true
	case "INVENTORY_ADJUSTED":
		return "notifications.inventory.adjusted.title", inventoryBody, true
	case "USER_STATUS_CHANGED":
		return "notifications.account.statusChanged.title", accountBody, true
	case "USER_ROLES_CHANGED":
		return "notifications.account.rolesChanged.title", accountBody, true
	case "USER_REGISTERED":
		return "notifications.account.registered.title", accountBody, true
	case "PASSWORD_RESET_REQUESTED":
		return "notifications.account.passwordReset.title", accountBody, true
	default:
		// Keep the outbox forward-compatible: every committed event with a
		// resolvable recipient still gets an in-app and email notification even
		// before a dedicated localized title is added.
		if strings.TrimSpace(eventType) != "" {
			return "notifications.account.title", "notifications.account.body", true
		}
		return "", "", false
	}
}

func intervalLiteral(value time.Duration) string {
	seconds := int64(value / time.Second)
	if seconds < 1 {
		seconds = 1
	}
	return fmt.Sprintf("%d seconds", seconds)
}

func truncateError(err error) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	if len(message) > 2000 {
		return message[:2000]
	}
	return message
}
