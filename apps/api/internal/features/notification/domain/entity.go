package domain

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrNotFound   = errors.New("notification not found")
	ErrValidation = errors.New("invalid notification request")
)

type Notification struct {
	RecipientUserID string          `json:"-"`
	ID              string          `json:"id"`
	EventType       string          `json:"eventType"`
	TitleKey        string          `json:"titleKey"`
	BodyKey         string          `json:"bodyKey"`
	Payload         json.RawMessage `json:"payload"`
	ReadAt          *time.Time      `json:"readAt,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
}

type OutboxEvent struct {
	ID            string
	EventType     string
	AggregateType string
	AggregateID   string
	Payload       json.RawMessage
	AttemptCount  int
}

type Repository interface {
	List(context.Context, string, int, int) ([]Notification, int, error)
	UnreadCount(context.Context, string) (int, error)
	MarkRead(context.Context, string, string) error
	MarkAllRead(context.Context, string) error
	ClaimOutbox(context.Context, int, time.Duration) ([]OutboxEvent, error)
	DeliverOutbox(context.Context, OutboxEvent) ([]Notification, error)
	MarkOutboxProcessed(context.Context, string) error
	MarkOutboxFailed(context.Context, string, time.Time, error) error
}
