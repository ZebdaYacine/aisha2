package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/notification/domain"
)

type workerRepository struct {
	events     []domain.OutboxEvent
	items      []domain.Notification
	deliverErr error
	processed  []string
	failed     []string
}

func (r *workerRepository) List(context.Context, string, int, int) ([]domain.Notification, int, error) {
	return nil, 0, nil
}
func (r *workerRepository) UnreadCount(context.Context, string) (int, error) { return 0, nil }
func (r *workerRepository) MarkRead(context.Context, string, string) error   { return nil }
func (r *workerRepository) MarkAllRead(context.Context, string) error        { return nil }
func (r *workerRepository) ClaimOutbox(context.Context, int, time.Duration) ([]domain.OutboxEvent, error) {
	events := r.events
	r.events = nil
	return events, nil
}
func (r *workerRepository) DeliverOutbox(context.Context, domain.OutboxEvent) ([]domain.Notification, error) {
	return r.items, r.deliverErr
}
func (r *workerRepository) MarkOutboxProcessed(_ context.Context, id string) error {
	r.processed = append(r.processed, id)
	return nil
}
func (r *workerRepository) MarkOutboxFailed(_ context.Context, id string, _ time.Time, _ error) error {
	r.failed = append(r.failed, id)
	return nil
}

type workerPublisher struct {
	userID string
	item   domain.Notification
}

func (p *workerPublisher) Publish(userID string, item domain.Notification) {
	p.userID = userID
	p.item = item
}

func TestWorkerPersistsAndPublishesAfterProcessing(t *testing.T) {
	repository := &workerRepository{
		events: []domain.OutboxEvent{{ID: "event-1", AttemptCount: 1}},
		items:  []domain.Notification{{RecipientUserID: "user-1", ID: "notification-1"}},
	}
	publisher := &workerPublisher{}
	worker := NewWorker(repository, publisher, time.Second, nil)

	worker.tick(context.Background())

	if len(repository.processed) != 1 || repository.processed[0] != "event-1" {
		t.Fatalf("processed = %#v", repository.processed)
	}
	if publisher.userID != "user-1" || publisher.item.ID != "notification-1" {
		t.Fatalf("published = %#v", publisher)
	}
}

func TestWorkerSchedulesRetryWhenDeliveryFails(t *testing.T) {
	repository := &workerRepository{
		events:     []domain.OutboxEvent{{ID: "event-2", AttemptCount: 2}},
		deliverErr: errors.New("temporary database failure"),
	}
	worker := NewWorker(repository, nil, time.Second, nil)

	worker.tick(context.Background())

	if len(repository.failed) != 1 || repository.failed[0] != "event-2" {
		t.Fatalf("failed = %#v", repository.failed)
	}
	if len(repository.processed) != 0 {
		t.Fatalf("unexpected processed events = %#v", repository.processed)
	}
}
