package response

import (
	"reflect"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/pkg/health"
)

func TestHealthResponseFrom(t *testing.T) {
	services := map[string]string{"postgres": "ok"}
	got := HealthResponseFrom(health.Result{Status: "ok", Services: services})
	if got.Status != "ok" || !reflect.DeepEqual(got.Services, services) {
		t.Fatalf("unexpected health DTO: %#v", got)
	}
}
