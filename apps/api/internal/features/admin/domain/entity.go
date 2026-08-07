package domain

import (
	"context"
	"errors"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

var (
	ErrNotFound   = errors.New("admin resource not found")
	ErrValidation = errors.New("admin request validation failed")
)

type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	Status      string    `json:"status"`
	Roles       []string  `json:"roles"`
	CreatedAt   time.Time `json:"createdAt"`
}

type AuditEvent struct {
	ID            string    `json:"id"`
	EventType     string    `json:"eventType"`
	ActorUserID   string    `json:"actorUserId,omitempty"`
	TargetType    string    `json:"targetType"`
	TargetID      string    `json:"targetId,omitempty"`
	CorrelationID string    `json:"correlationId,omitempty"`
	Reason        string    `json:"reason,omitempty"`
	PreviousState any       `json:"previousState,omitempty"`
	NewState      any       `json:"newState,omitempty"`
	OccurredAt    time.Time `json:"occurredAt"`
}

type AuditFilter struct {
	ActorUserID, TargetType, TargetID, EventType, CorrelationID string
	From, To                                                    *time.Time
}

type Repository interface {
	ListUsers(context.Context, string, string, int, int) ([]User, int, error)
	SetRoles(context.Context, string, string, []string) (User, error)
	AuditEvents(context.Context, AuditFilter, int, int) ([]AuditEvent, int, error)
}

type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
