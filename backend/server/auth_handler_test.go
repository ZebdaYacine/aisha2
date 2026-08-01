package httpapi

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func TestAuthHandlerReturnsFieldValidationErrors(t *testing.T) {
	handler := NewAuthHandler(nil, NewRequestValidator())
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Use(requestid.New())
	app.Post("/register", handler.Register)
	request := httptest.NewRequest("POST", "/register", strings.NewReader(`{"email":"invalid","password":"short","displayName":"A"}`))
	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.StatusCode)
	}
	var body ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Error.Code != CodeValidationError || body.Error.Fields["email"] != "INVALID_EMAIL" {
		t.Fatalf("unexpected validation response: %#v", body)
	}
}
