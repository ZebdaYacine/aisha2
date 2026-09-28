package httpapi

import (
	"errors"
	"testing"
)

func TestRequestValidatorReturnsJSONFieldCodes(t *testing.T) {
	err := NewRequestValidator().Validate(RegisterRequest{
		Email:       "not-an-email",
		Password:    "short",
		DisplayName: "A",
	})
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected API error, got %v", err)
	}
	if apiErr.Code != CodeValidationError {
		t.Fatalf("expected validation error, got %s", apiErr.Code)
	}
	if apiErr.Fields["email"] != "INVALID_EMAIL" || apiErr.Fields["password"] != "MIN_LENGTH_12" || apiErr.Fields["displayName"] != "MIN_LENGTH_2" {
		t.Fatalf("unexpected validation fields: %#v", apiErr.Fields)
	}
}

func TestRequestValidatorAcceptsValidAuthenticationDTOs(t *testing.T) {
	validator := NewRequestValidator()
	requests := []any{
		RegisterRequest{Email: "user@example.com", Password: "long-password-value", DisplayName: "Amina"},
		LoginRequest{Identifier: "user@example.com", Password: "password"},
		LoginRequest{Email: "user@example.com", Password: "password"},
		RefreshRequest{RefreshToken: "token"},
		LogoutRequest{RefreshToken: "token"},
		ForgotPasswordRequest{Email: "user@example.com"},
		ResetPasswordRequest{Token: "token", Password: "long-password-value"},
	}
	for _, request := range requests {
		if err := validator.Validate(request); err != nil {
			t.Fatalf("expected valid request %T, got %v", request, err)
		}
	}
}
