package error_handler

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/pkg/apperror"
	"github.com/gofiber/fiber/v3"
)

func TestHandlerLogsUnderlyingCauseAndFinalStatus(t *testing.T) {
	var logs bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })

	app := fiber.New(fiber.Config{ErrorHandler: Handler})
	app.Get("/artisans", func(fiber.Ctx) error {
		return WrapAPIError(errors.New("query public artisans: relation does not exist"), apperror.InternalError, "An unexpected error occurred.")
	})

	response, err := app.Test(httptest.NewRequest("GET", "/artisans", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.StatusCode, fiber.StatusInternalServerError)
	}
	output := logs.String()
	if !strings.Contains(output, "query public artisans: relation does not exist") {
		t.Fatalf("log does not contain underlying cause: %s", output)
	}
	if !strings.Contains(output, "status=500") {
		t.Fatalf("log does not contain final status: %s", output)
	}
}
