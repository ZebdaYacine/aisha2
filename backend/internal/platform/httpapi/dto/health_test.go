package dto

import (
	"reflect"
	"testing"

	"github.com/aisha-platform/aisha/backend/internal/platform/health"
)

func TestHealthResponseFrom(t *testing.T) {
	services := map[string]string{"postgres": "ok"}
	got := HealthResponseFrom(health.Result{Status: "ok", Services: services})
	if got.Status != "ok" || !reflect.DeepEqual(got.Services, services) {
		t.Fatalf("unexpected health DTO: %#v", got)
	}
}
