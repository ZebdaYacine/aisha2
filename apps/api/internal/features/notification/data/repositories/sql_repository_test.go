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
		{event: "PRODUCT_MODERATION_DECIDED", want: 0},
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
