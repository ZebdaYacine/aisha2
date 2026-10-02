package application

import (
	"context"
	"errors"
	"strings"
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

type workerMailer struct {
	to      string
	subject string
	body    string
	err     error
}

func (m *workerMailer) Send(_ context.Context, to, subject, body string) error {
	m.to, m.subject, m.body = to, subject, body
	return m.err
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

func TestWorkerSendsTransactionalEmailForDeliveredNotification(t *testing.T) {
	repository := &workerRepository{
		events: []domain.OutboxEvent{{ID: "event-3", AttemptCount: 1}},
		items: []domain.Notification{{
			RecipientUserID: "user-1",
			RecipientEmail:  "nour@example.test",
			RecipientName:   "Nour",
			EventType:       "ORDER_CHECKOUT_CREATED",
			Payload:         []byte(`{"orderNumber":"AIS-123","userId":"internal-user-uuid"}`),
		}},
	}
	mailer := &workerMailer{}
	worker := NewWorker(repository, nil, time.Second, nil, mailer)

	worker.tick(context.Background())

	if mailer.to != "nour@example.test" {
		t.Fatalf("recipient = %q", mailer.to)
	}
	if mailer.subject != "AISHA notification — Order Checkout Created" {
		t.Fatalf("subject = %q", mailer.subject)
	}
	if !strings.Contains(mailer.body, "Hello Nour") || !strings.Contains(mailer.body, "AIS-123") {
		t.Fatalf("body = %q", mailer.body)
	}
	if !strings.Contains(mailer.body, "<style>") || strings.Contains(mailer.body, "userId") || strings.Contains(mailer.body, "internal-user-uuid") {
		t.Fatalf("email contains unsafe or unstyled details: %q", mailer.body)
	}
}

func TestWorkerRetriesOutboxWhenEmailDeliveryFails(t *testing.T) {
	repository := &workerRepository{
		events: []domain.OutboxEvent{{ID: "event-email-failure", AttemptCount: 1}},
		items: []domain.Notification{{
			RecipientUserID: "user-1",
			RecipientEmail:  "nour@example.test",
			EventType:       "PAYMENT_FAILED",
		}},
	}
	mailer := &workerMailer{err: errors.New("smtp unavailable")}
	worker := NewWorker(repository, nil, time.Second, nil, mailer)

	worker.tick(context.Background())

	if len(repository.failed) != 1 || repository.failed[0] != "event-email-failure" {
		t.Fatalf("failed = %#v", repository.failed)
	}
	if len(repository.processed) != 0 {
		t.Fatalf("event was processed despite email failure: %#v", repository.processed)
	}
}

func TestWorkerEmailUsesRegistrationActivationLink(t *testing.T) {
	repository := &workerRepository{
		events: []domain.OutboxEvent{{ID: "event-registration", AttemptCount: 1}},
		items: []domain.Notification{{
			RecipientUserID: "user-1",
			RecipientEmail:  "nour@example.test",
			RecipientName:   "Nour",
			EventType:       "USER_REGISTERED",
			Payload:         []byte(`{"activationUrl":"https://aisha.example/en/activate?token=raw-token"}`),
		}},
	}
	mailer := &workerMailer{}
	NewWorker(repository, nil, time.Second, nil, mailer).tick(context.Background())

	if !strings.Contains(mailer.body, `href="https://aisha.example/en/activate?token=raw-token"`) || !strings.Contains(mailer.body, "Activate your AISHA account") {
		t.Fatalf("activation link missing from email: %q", mailer.body)
	}
}
