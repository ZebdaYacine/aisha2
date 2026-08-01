package health

import (
	"context"
	"errors"
	"testing"
)

type fakeChecker struct {
	name string
	err  error
}

func (f fakeChecker) Name() string                { return f.name }
func (f fakeChecker) Check(context.Context) error { return f.err }

func TestReadyReportsAllDependencies(t *testing.T) {
	service := New(fakeChecker{name: "postgres"}, fakeChecker{name: "redis"}, fakeChecker{name: "minio"})
	result, err := service.Ready(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "ok" || len(result.Services) != 3 {
		t.Fatalf("unexpected readiness result: %#v", result)
	}
}

func TestReadyFailsWhenDependencyIsUnavailable(t *testing.T) {
	service := New(fakeChecker{name: "postgres", err: errors.New("down")})
	result, err := service.Ready(context.Background())
	if err == nil || result.Status != "unavailable" || result.Services["postgres"] != "unavailable" {
		t.Fatalf("unexpected readiness result: %#v, %v", result, err)
	}
}
