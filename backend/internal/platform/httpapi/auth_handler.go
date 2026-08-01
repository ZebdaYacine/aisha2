package httpapi

import (
	"errors"
	"strings"

	"github.com/aisha-platform/aisha/backend/internal/auth"
	"github.com/gofiber/fiber/v3"
)

const principalLocal = "authenticated_principal"

type AuthHandler struct{ service *auth.Service }

func NewAuthHandler(service *auth.Service) *AuthHandler { return &AuthHandler{service: service} }
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	user, tokens, err := h.service.Register(c.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		return authAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(AuthenticationResponse{User: userResponseFrom(user), Tokens: tokenResponseFrom(tokens)})
}
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	user, tokens, err := h.service.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return authAPIError(err)
	}
	return c.JSON(AuthenticationResponse{User: userResponseFrom(user), Tokens: tokenResponseFrom(tokens)})
}
func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req RefreshRequest
	if err := c.Bind().Body(&req); err != nil || req.RefreshToken == "" {
		return NewAPIError(CodeValidationError, "The request is invalid.", FieldErrors{"refreshToken": "REQUIRED"})
	}
	tokens, err := h.service.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return authAPIError(err)
	}
	return c.JSON(tokenResponseFrom(tokens))
}
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	var req LogoutRequest
	if err := c.Bind().Body(&req); err != nil || req.RefreshToken == "" {
		return NewAPIError(CodeValidationError, "The request is invalid.", FieldErrors{"refreshToken": "REQUIRED"})
	}
	if err := h.service.Logout(c.Context(), req.RefreshToken); err != nil {
		return authAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *AuthHandler) ForgotPassword(c fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.Bind().Body(&req); err != nil || strings.TrimSpace(req.Email) == "" {
		return NewAPIError(CodeValidationError, "The request is invalid.", FieldErrors{"email": "REQUIRED"})
	}
	if err := h.service.ForgotPassword(c.Context(), req.Email); err != nil {
		return authAPIError(err)
	}
	return c.SendStatus(fiber.StatusAccepted)
}
func (h *AuthHandler) ResetPassword(c fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.Bind().Body(&req); err != nil || req.Token == "" || len(req.Password) < 12 {
		return NewAPIError(CodeValidationError, "The request is invalid.", FieldErrors{"password": "MIN_LENGTH_12"})
	}
	if err := h.service.ResetPassword(c.Context(), req.Token, req.Password); err != nil {
		return authAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *AuthHandler) Me(c fiber.Ctx) error {
	principal, ok := c.Locals(principalLocal).(auth.Principal)
	if !ok {
		return NewAPIError(CodeAuthenticationRequired, "Authentication is required.", nil)
	}
	user, err := h.service.Me(c.Context(), principal)
	if err != nil {
		return authAPIError(err)
	}
	return c.JSON(userResponseFrom(user))
}
func (h *AuthHandler) RequirePrincipal(c fiber.Ctx) error {
	header := c.Get(fiber.HeaderAuthorization)
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok || token == "" {
		return NewAPIError(CodeAuthenticationRequired, "Authentication is required.", nil)
	}
	principal, err := h.service.Authenticate(c.Context(), token)
	if err != nil {
		return authAPIError(err)
	}
	c.Locals(principalLocal, principal)
	return c.Next()
}
func authAPIError(err error) error {
	switch {
	case errors.Is(err, auth.ErrValidation):
		return NewAPIError(CodeValidationError, "The request is invalid.", nil)
	case errors.Is(err, auth.ErrEmailExists):
		return NewAPIError(CodeConflict, "An account with this email already exists.", FieldErrors{"email": "EMAIL_ALREADY_EXISTS"})
	case errors.Is(err, auth.ErrInvalidCredentials), errors.Is(err, auth.ErrInvalidToken):
		return NewAPIError(CodeInvalidCredentials, "The credentials are invalid.", nil)
	case errors.Is(err, auth.ErrUserInactive):
		return NewAPIError(CodeForbidden, "This account is not active.", nil)
	default:
		return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
	}
}
