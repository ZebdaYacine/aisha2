package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/aisha-platform/aisha/backend/internal/platform/config"
	"github.com/aisha-platform/aisha/backend/internal/platform/health"
)

type checker struct{ err error }

func (checker) Name() string                  { return "dependency" }
func (c checker) Check(context.Context) error { return c.err }

func TestLiveHealth(t *testing.T) {
	app := New(config.Config{AppName: "AISHA", AllowedOrigins: []string{"http://localhost:3000"}}, health.New())
	response, err := app.Test(httptest.NewRequest("GET", "/health/live", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}
}

func TestReadinessFailure(t *testing.T) {
	app := New(config.Config{AppName: "AISHA", AllowedOrigins: []string{"http://localhost:3000"}}, health.New(checker{err: errors.New("down")}))
	response, err := app.Test(httptest.NewRequest("GET", "/health/ready", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 503 {
		t.Fatalf("expected 503, got %d", response.StatusCode)
	}
}

func TestUnknownRouteUsesStandardErrorDTO(t *testing.T) {
	app := New(config.Config{AppName: "AISHA", AllowedOrigins: []string{"http://localhost:3000"}}, health.New())
	response, err := app.Test(httptest.NewRequest("GET", "/missing", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 404 {
		t.Fatalf("expected 404, got %d", response.StatusCode)
	}
	var body ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != CodeResourceNotFound || body.Error.RequestID == "" {
		t.Fatalf("unexpected response: %#v", body)
	}
}
