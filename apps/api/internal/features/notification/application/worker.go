package application

import (
	"context"
	"log/slog"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/notification/domain"
)

type Publisher interface {
	Publish(string, domain.Notification)
}

type Worker struct {
	repository domain.Repository
	publisher  Publisher
	interval   time.Duration
	logger     *slog.Logger
}

func NewWorker(repository domain.Repository, publisher Publisher, interval time.Duration, logger *slog.Logger) *Worker {
	if interval <= 0 {
		interval = time.Second
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Worker{repository: repository, publisher: publisher, interval: interval, logger: logger}
}

func (w *Worker) Run(ctx context.Context) {
	w.tick(ctx)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	events, err := w.repository.ClaimOutbox(ctx, 50, 30*time.Second)
	if err != nil {
		if ctx.Err() == nil {
			w.logger.Error("claim notification outbox events", "error", err)
		}
		return
	}
	for _, event := range events {
		notifications, err := w.repository.DeliverOutbox(ctx, event)
		if err != nil {
			next := time.Now().UTC().Add(backoff(event.AttemptCount))
			if markErr := w.repository.MarkOutboxFailed(ctx, event.ID, next, err); markErr != nil {
				w.logger.Error("mark notification outbox event failed", "event_id", event.ID, "error", markErr)
			}
			continue
		}
		if err = w.repository.MarkOutboxProcessed(ctx, event.ID); err != nil {
			w.logger.Error("mark notification outbox event processed", "event_id", event.ID, "error", err)
			continue
		}
		for _, notification := range notifications {
			if w.publisher != nil {
				w.publisher.Publish(notification.RecipientUserID, notification)
			}
		}
	}
}

func backoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 6 {
		attempt = 6
	}
	return time.Duration(1<<uint(attempt-1)) * time.Second
}
