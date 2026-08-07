package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aisha-platform/aisha/apps/api/internal/features/admin/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) ListUsers(ctx context.Context, role, status string, limit, offset int) ([]domain.User, int, error) {
	rows, err := r.pool.Query(ctx, `SELECT u.id,COALESCE(u.email,''),COALESCE(u.display_name,''),u.status,u.created_at,COALESCE(array_agg(r.code) FILTER(WHERE r.code IS NOT NULL),'{}'),count(*) OVER() FROM users u LEFT JOIN user_roles ur ON ur.user_id=u.id LEFT JOIN roles r ON r.id=ur.role_id WHERE ($1='' OR u.status=$1) AND ($2='' OR EXISTS(SELECT 1 FROM user_roles ur2 JOIN roles r2 ON r2.id=ur2.role_id WHERE ur2.user_id=u.id AND r2.code=$2)) GROUP BY u.id ORDER BY u.created_at DESC LIMIT $3 OFFSET $4`, status, role, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []domain.User{}
	total := 0
	for rows.Next() {
		var item domain.User
		if err = rows.Scan(&item.ID, &item.Email, &item.DisplayName, &item.Status, &item.CreatedAt, &item.Roles, &total); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) SetRoles(ctx context.Context, actorID, userID string, roles []string) (domain.User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var exists bool
	if err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, userID).Scan(&exists); err != nil {
		return domain.User{}, err
	}
	if !exists {
		return domain.User{}, domain.ErrNotFound
	}
	var previous []string
	if err = tx.QueryRow(ctx, `SELECT COALESCE(array_agg(r.code ORDER BY r.code),'{}') FROM user_roles ur JOIN roles r ON r.id=ur.role_id WHERE ur.user_id=$1`, userID).Scan(&previous); err != nil {
		return domain.User{}, err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id=$1`, userID); err != nil {
		return domain.User{}, err
	}
	for _, role := range roles {
		result, insertErr := tx.Exec(ctx, `INSERT INTO user_roles(user_id,role_id,assigned_by_user_id) SELECT $1,id,$2 FROM roles WHERE code=$3`, userID, actorID, role)
		if insertErr != nil {
			return domain.User{}, insertErr
		}
		if result.RowsAffected() != 1 {
			return domain.User{}, domain.ErrValidation
		}
	}
	previousJSON, _ := json.Marshal(previous)
	nextJSON, _ := json.Marshal(roles)
	if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,previous_state,new_state) VALUES('USER_ROLES_CHANGED',$1,'user',$2,$3,$4)`, actorID, userID, previousJSON, nextJSON); err != nil {
		return domain.User{}, err
	}
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES('USER_ROLES_CHANGED','user',$1,jsonb_build_object('userId',$1::text))`, userID); err != nil {
		return domain.User{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.User{}, err
	}
	return r.user(ctx, userID)
}

func (r *PostgresRepository) AuditEvents(ctx context.Context, filter domain.AuditFilter, limit, offset int) ([]domain.AuditEvent, int, error) {
	conditions := []string{"1=1"}
	args := []any{}
	add := func(condition string, value any) {
		args = append(args, value)
		conditions = append(conditions, fmt.Sprintf(condition, len(args)))
	}
	if filter.ActorUserID != "" {
		add("actor_user_id=$%d", filter.ActorUserID)
	}
	if filter.TargetType != "" {
		add("target_type=$%d", filter.TargetType)
	}
	if filter.TargetID != "" {
		add("target_id=$%d", filter.TargetID)
	}
	if filter.EventType != "" {
		add("event_type=$%d", filter.EventType)
	}
	if filter.CorrelationID != "" {
		add("correlation_id=$%d", filter.CorrelationID)
	}
	if filter.From != nil {
		add("occurred_at >= $%d", *filter.From)
	}
	if filter.To != nil {
		add("occurred_at < $%d", *filter.To)
	}
	args = append(args, limit, offset)
	query := `SELECT id,COALESCE(actor_user_id::text,''),target_type,COALESCE(target_id::text,''),COALESCE(correlation_id,''),COALESCE(reason,''),COALESCE(previous_state,'null'::jsonb),COALESCE(new_state,'null'::jsonb),occurred_at,count(*) OVER() FROM audit_events WHERE ` + strings.Join(conditions, " AND ") + fmt.Sprintf(" ORDER BY occurred_at DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []domain.AuditEvent{}
	total := 0
	for rows.Next() {
		var item domain.AuditEvent
		var previous, next []byte
		if err = rows.Scan(&item.ID, &item.ActorUserID, &item.TargetType, &item.TargetID, &item.CorrelationID, &item.Reason, &previous, &next, &item.OccurredAt, &total); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(previous, &item.PreviousState)
		_ = json.Unmarshal(next, &item.NewState)
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) user(ctx context.Context, id string) (domain.User, error) {
	var item domain.User
	err := r.pool.QueryRow(ctx, `SELECT u.id,COALESCE(u.email,''),COALESCE(u.display_name,''),u.status,u.created_at,COALESCE(array_agg(r.code) FILTER(WHERE r.code IS NOT NULL),'{}') FROM users u LEFT JOIN user_roles ur ON ur.user_id=u.id LEFT JOIN roles r ON r.id=ur.role_id WHERE u.id=$1 GROUP BY u.id`, id).Scan(&item.ID, &item.Email, &item.DisplayName, &item.Status, &item.CreatedAt, &item.Roles)
	return item, err
}
