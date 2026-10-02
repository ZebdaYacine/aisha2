package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/notification/domain"
)

type Publisher interface {
	Publish(string, domain.Notification)
}

type Mailer interface {
	Send(context.Context, string, string, string) error
}

type Worker struct {
	repository domain.Repository
	publisher  Publisher
	interval   time.Duration
	logger     *slog.Logger
	mailer     Mailer
}

func NewWorker(repository domain.Repository, publisher Publisher, interval time.Duration, logger *slog.Logger, mailers ...Mailer) *Worker {
	if interval <= 0 {
		interval = time.Second
	}
	if logger == nil {
		logger = slog.Default()
	}
	var mailer Mailer
	if len(mailers) > 0 {
		mailer = mailers[0]
	}
	return &Worker{repository: repository, publisher: publisher, interval: interval, logger: logger, mailer: mailer}
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
			if w.mailer != nil && strings.TrimSpace(notification.RecipientEmail) != "" {
				if err := w.mailer.Send(ctx, notification.RecipientEmail, emailSubject(notification.EventType), emailBody(notification)); err != nil {
					w.logger.Error("send notification email", "notification_id", notification.ID, "event_type", notification.EventType, "error", err)
				}
			}
		}
	}
}

func emailSubject(eventType string) string {
	return "AISHA notification — " + humanizeEventType(eventType)
}

func emailBody(notification domain.Notification) string {
	name := strings.TrimSpace(notification.RecipientName)
	if name == "" {
		name = "there"
	}
	body := fmt.Sprintf("Hello %s,\n\nA new AISHA workflow update is available.\n\nTransaction: %s\n", name, humanizeEventType(notification.EventType))
	if len(notification.Payload) > 0 && string(notification.Payload) != "null" && string(notification.Payload) != "{}" {
		var payload any
		if err := json.Unmarshal(notification.Payload, &payload); err == nil {
			if formatted, err := json.MarshalIndent(payload, "", "  "); err == nil {
				body += "\nDetails:\n" + string(formatted) + "\n"
			}
		}
	}
	return body + "\nSign in to AISHA to review the latest status.\n"
}

func humanizeEventType(eventType string) string {
	words := strings.Fields(strings.ReplaceAll(strings.ToLower(strings.TrimSpace(eventType)), "_", " "))
	for index, word := range words {
		if word != "" {
			words[index] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
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
