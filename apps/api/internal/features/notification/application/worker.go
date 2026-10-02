package application

import (
	"context"
	"encoding/json"
	"html"
	"log/slog"
	"strconv"
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
		var emailErr error
		for _, notification := range notifications {
			if w.publisher != nil {
				w.publisher.Publish(notification.RecipientUserID, notification)
			}
			if w.mailer != nil && strings.TrimSpace(notification.RecipientEmail) != "" {
				if err := w.mailer.Send(ctx, notification.RecipientEmail, emailSubject(notification.EventType), emailBody(notification)); err != nil {
					emailErr = err
					w.logger.Error("send notification email", "notification_id", notification.ID, "event_type", notification.EventType, "error", err)
				}
			}
		}
		if emailErr != nil {
			next := time.Now().UTC().Add(backoff(event.AttemptCount))
			if markErr := w.repository.MarkOutboxFailed(ctx, event.ID, next, emailErr); markErr != nil {
				w.logger.Error("mark notification email retry", "event_id", event.ID, "error", markErr)
			}
			continue
		}
		if err = w.repository.MarkOutboxProcessed(ctx, event.ID); err != nil {
			w.logger.Error("mark notification outbox event processed", "event_id", event.ID, "error", err)
		}
	}
}

func emailSubject(eventType string) string {
	return "AISHA notification — " + humanizeEventType(eventType)
}

func emailBody(notification domain.Notification) string {
	name := html.EscapeString(strings.TrimSpace(notification.RecipientName))
	if name == "" {
		name = "there"
	}
	title := html.EscapeString(humanizeEventType(notification.EventType))
	body := strings.Builder{}
	body.WriteString(`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>AISHA notification</title><style>
body{margin:0;background:#f5f1eb;color:#2b211b;font-family:Arial,Helvetica,sans-serif;line-height:1.6}.shell{max-width:620px;margin:32px auto;padding:0 16px}.card{overflow:hidden;background:#fffdf9;border:1px solid #e1d7ca;border-radius:18px;box-shadow:0 8px 28px rgba(53,37,25,.08)}.brand{padding:22px 28px;background:#2b211b;color:#fffdf9;font-family:Georgia,serif;font-size:22px;letter-spacing:.22em}.content{padding:30px 28px}.eyebrow{margin:0 0 8px;color:#805131;font-size:11px;font-weight:700;letter-spacing:.18em;text-transform:uppercase}.title{margin:0;color:#2b211b;font-family:Georgia,serif;font-size:30px;font-weight:400;line-height:1.2}.greeting{margin:22px 0 8px;font-size:16px}.summary{margin:0 0 24px;color:#6e6258}.details{margin:0 0 24px;border:1px solid #e1d7ca;border-radius:12px}.detail{display:flex;justify-content:space-between;gap:20px;padding:12px 14px;border-bottom:1px solid #eee7df;font-size:14px}.detail:last-child{border-bottom:0}.label{color:#806f61}.value{font-weight:700;text-align:right}.cta{display:inline-block;padding:11px 18px;border-radius:8px;background:#805131;color:#fff;text-decoration:none;font-size:14px;font-weight:700}.footer{padding:18px 28px;border-top:1px solid #e1d7ca;color:#8d8176;font-size:12px}@media(max-width:640px){.shell{margin:12px auto}.content{padding:24px 20px}.brand{padding:18px 20px}.title{font-size:26px}.detail{display:block}.value{display:block;margin-top:3px;text-align:left}}
</style></head><body><main class="shell"><article class="card"><div class="brand">AISHA</div><div class="content"><p class="eyebrow">Workflow update</p><h1 class="title">`)
	body.WriteString(title)
	body.WriteString(`</h1><p class="greeting">Hello `)
	body.WriteString(name)
	body.WriteString(`,</p><p class="summary">`)
	body.WriteString(html.EscapeString(emailSummary(notification.EventType)))
	body.WriteString(`</p>`)
	for _, detail := range emailDetails(notification.Payload) {
		body.WriteString(`<div class="details"><div class="detail"><span class="label">`)
		body.WriteString(html.EscapeString(detail.label))
		body.WriteString(`</span><span class="value">`)
		body.WriteString(html.EscapeString(detail.value))
		body.WriteString(`</span></div></div>`)
	}
	ctaLabel, ctaURL := emailCTA(notification)
	body.WriteString(`<a class="cta" href="`)
	body.WriteString(html.EscapeString(ctaURL))
	body.WriteString(`">`)
	body.WriteString(html.EscapeString(ctaLabel))
	body.WriteString(`</a></div><footer class="footer">This message contains workflow information only. Internal identifiers and private system data are never included.</footer></article></main></body></html>`)
	return body.String()
}

func emailCTA(notification domain.Notification) (string, string) {
	if notification.EventType == "USER_REGISTERED" {
		var fields map[string]any
		if json.Unmarshal(notification.Payload, &fields) == nil {
			if rawURL, ok := fields["activationUrl"].(string); ok && strings.HasPrefix(rawURL, "http") {
				return "Activate your AISHA account", rawURL
			}
		}
	}
	return "Sign in to AISHA", "#"
}

type emailDetail struct {
	label string
	value string
}

func emailSummary(eventType string) string {
	switch eventType {
	case "USER_REGISTERED":
		return "Your AISHA account was created successfully. Activate it using the button below."
	case "PAYMENT_CONFIRMED":
		return "A payment was confirmed for an AISHA order."
	case "PAYMENT_FAILED":
		return "A payment could not be completed. Please review the order status in AISHA."
	case "ORDER_CHECKOUT_CREATED":
		return "Your order was created and is waiting for payment confirmation."
	case "ORDER_CANCELLED":
		return "An AISHA order was cancelled successfully."
	case "WAREHOUSE_RECEPTION_CREATED", "RECEIVED_PENDING_INSPECTION":
		return "A warehouse reception was recorded and is waiting for inspection."
	case "WAREHOUSE_RECEPTION_INSPECTED", "INSPECTED":
		return "A warehouse inspection was completed and the stock status was updated."
	default:
		return "AISHA recorded a new workflow update for your account."
	}
}

func emailDetails(payload json.RawMessage) []emailDetail {
	var fields map[string]any
	if len(payload) == 0 || string(payload) == "null" || string(payload) == "{}" || json.Unmarshal(payload, &fields) != nil {
		return nil
	}
	ordered := []struct {
		key   string
		label string
	}{
		{key: "orderNumber", label: "Order"},
		{key: "status", label: "Status"},
		{key: "action", label: "Action"},
		{key: "provider", label: "Payment method"},
		{key: "reason", label: "Reason"},
		{key: "quantityDelta", label: "Quantity change"},
	}
	details := make([]emailDetail, 0, len(ordered))
	for _, field := range ordered {
		value, ok := fields[field.key]
		if !ok {
			continue
		}
		formatted := safeEmailValue(value)
		if formatted != "" {
			details = append(details, emailDetail{label: field.label, value: formatted})
		}
	}
	return details
}

func safeEmailValue(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case float64:
		if typed == float64(int64(typed)) {
			return strconv.FormatInt(int64(typed), 10)
		}
		return strconv.FormatFloat(typed, 'f', 2, 64)
	default:
		return ""
	}
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
