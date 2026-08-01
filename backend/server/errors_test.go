package httpapi

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func TestStatusForCode(t *testing.T) {
	tests := []struct {
		code ErrorCode
		want int
	}{
		{CodeValidationError, fiber.StatusBadRequest},
		{CodeAuthenticationRequired, fiber.StatusUnauthorized},
		{CodeForbidden, fiber.StatusForbidden},
		{CodeResourceNotFound, fiber.StatusNotFound},
		{CodeIdempotencyConflict, fiber.StatusConflict},
		{CodeRateLimited, fiber.StatusTooManyRequests},
		{CodeInternalError, fiber.StatusInternalServerError},
	}
	for _, test := range tests {
		t.Run(string(test.code), func(t *testing.T) {
			if got := statusForCode(test.code); got != test.want {
				t.Fatalf("statusForCode(%q) = %d, want %d", test.code, got, test.want)
			}
		})
	}
}

func TestErrorHandlerWritesValidationEnvelope(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Use(requestid.New())
	app.Get("/failure", func(fiber.Ctx) error {
		return NewAPIError(CodeValidationError, "The request is invalid.", FieldErrors{"email": "INVALID_EMAIL"})
	})
	request := httptest.NewRequest("GET", "/failure", nil)
	request.Header.Set(fiber.HeaderXRequestID, "request-123")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.StatusCode, fiber.StatusBadRequest)
	}
	var body ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != CodeValidationError || body.Error.Fields["email"] != "INVALID_EMAIL" || body.Error.RequestID == "" {
		t.Fatalf("unexpected error response: %#v", body)
	}
}

func TestErrorHandlerDoesNotLeakInternalError(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Get("/failure", func(fiber.Ctx) error { return errors.New("database password leaked") })
	response, err := app.Test(httptest.NewRequest("GET", "/failure", nil))
	if err != nil {
		t.Fatal(err)
	}
	var body ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != CodeInternalError || body.Error.Message != "An unexpected error occurred." {
		t.Fatalf("unexpected internal response: %#v", body)
	}
}
