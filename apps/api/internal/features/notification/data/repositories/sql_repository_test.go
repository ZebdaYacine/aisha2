package repositories

import "testing"

func TestWorkflowAudienceRoles(t *testing.T) {
	tests := []struct {
		event string
		want  int
	}{
		{event: "ARTISAN_APPLICATION_SUBMITTED", want: 2},
		{event: "PRODUCT_SUBMITTED", want: 2},
		{event: "WAREHOUSE_RECEPTION_INSPECTED", want: 2},
		{event: "ORDER_CHECKOUT_CREATED", want: 2},
		{event: "PAYMENT_CONFIRMED", want: 2},
		{event: "PRODUCT_MODERATION_DECIDED", want: 0},
		{event: "USER_REGISTERED", want: 0},
	}
	for _, test := range tests {
		if got := len(audienceRoles(test.event)); got != test.want {
			t.Errorf("audienceRoles(%q) = %d, want %d", test.event, got, test.want)
		}
	}
}

func TestWarehouseWorkflowEventsHaveTemplates(t *testing.T) {
	for _, event := range []string{"WAREHOUSE_RECEPTION_CREATED", "WAREHOUSE_RECEPTION_INSPECTED", "INVENTORY_ACCEPTED"} {
		if _, _, ok := templateFor(event); !ok {
			t.Errorf("templateFor(%q) returned no template", event)
		}
	}
}

func TestPaymentWorkflowEventsHaveTemplates(t *testing.T) {
	for _, event := range []string{"PAYMENT_CONFIRMED", "PAYMENT_FAILED", "PAYMENT_REFUNDED"} {
		title, body, ok := templateFor(event)
		if !ok || title == "" || body != "notifications.payment.body" {
			t.Errorf("templateFor(%q) = %q, %q, %v", event, title, body, ok)
		}
	}
}

func TestAccountLifecycleEventsHaveTemplates(t *testing.T) {
	for _, event := range []string{"USER_REGISTERED", "PASSWORD_RESET_REQUESTED"} {
		title, body, ok := templateFor(event)
		if !ok || title == "" || body != "notifications.account.body" {
			t.Errorf("templateFor(%q) = %q, %q, %v", event, title, body, ok)
		}
	}
	for _, event := range []string{"USER_ACCOUNT_UPDATED", "FUTURE_WORKFLOW_EVENT"} {
		if _, _, ok := templateFor(event); !ok {
			t.Fatalf("%s should use the generic notification template", event)
		}
	}
	if _, _, ok := templateFor(""); ok {
		t.Fatal("empty workflow event should not create a notification")
	}
}
