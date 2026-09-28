package httpapi

import (
	"context"
	"errors"
	"strings"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	customer "github.com/aisha-platform/aisha/apps/api/internal/features/user"
	"github.com/gofiber/fiber/v3"
)

const principalLocal = "authenticated_principal"

type AuthHandler struct {
	service        *auth.Service
	validator      *RequestValidator
	accountSummary accountSummaryReader
}

type accountSummaryReader interface {
	AccountSummary(context.Context, auth.Principal) (customer.AccountSummary, error)
}

func NewAuthHandler(service *auth.Service, requestValidator *RequestValidator, accountSummaries ...accountSummaryReader) *AuthHandler {
	var summary accountSummaryReader
	if len(accountSummaries) > 0 {
		summary = accountSummaries[0]
	}
	return &AuthHandler{service: service, validator: requestValidator, accountSummary: summary}
}
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := h.bindAndValidate(c, &req); err != nil {
		return err
	}
	user, tokens, err := h.service.Register(c.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		return authAPIError(err)
	}
	response, err := h.userResponse(c.Context(), user)
	if err != nil {
		return authAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(AuthenticationResponse{User: response, Tokens: tokenResponseFrom(tokens)})
}
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := h.bindAndValidate(c, &req); err != nil {
		return err
	}
	identifier := req.Identifier
	if strings.TrimSpace(identifier) == "" {
		identifier = req.Email
	}
	user, tokens, err := h.service.Login(c.Context(), identifier, req.Password)
	if err != nil {
		return authAPIError(err)
	}
	response, err := h.userResponse(c.Context(), user)
	if err != nil {
		return authAPIError(err)
	}
	return c.JSON(AuthenticationResponse{User: response, Tokens: tokenResponseFrom(tokens)})
}
func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req RefreshRequest
	if err := h.bindAndValidate(c, &req); err != nil {
		return err
	}
	tokens, err := h.service.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return authAPIError(err)
	}
	return c.JSON(tokenResponseFrom(tokens))
}
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	var req LogoutRequest
	if err := h.bindAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.service.Logout(c.Context(), req.RefreshToken); err != nil {
		return authAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *AuthHandler) ChangePassword(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var req ChangePasswordRequest
	if err = h.bindAndValidate(c, &req); err != nil {
		return err
	}
	if err = h.service.ChangePassword(c.Context(), p, req.CurrentPassword, req.NewPassword); err != nil {
		return authAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) ForgotPassword(c fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := h.bindAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.service.ForgotPassword(c.Context(), req.Email); err != nil {
		return authAPIError(err)
	}
	return c.SendStatus(fiber.StatusAccepted)
}
func (h *AuthHandler) ResetPassword(c fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := h.bindAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.service.ResetPassword(c.Context(), req.Token, req.Password); err != nil {
		return authAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) bindAndValidate(c fiber.Ctx, request any) error {
	if err := c.Bind().Body(request); err != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	return h.validator.Validate(request)
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
	response, err := h.userResponse(c.Context(), user)
	if err != nil {
		return authAPIError(err)
	}
	return c.JSON(response)
}

func (h *AuthHandler) userResponse(ctx context.Context, user auth.User) (UserResponse, error) {
	response := userResponseFrom(user)
	if h.accountSummary == nil {
		return response, nil
	}
	summary, err := h.accountSummary.AccountSummary(ctx, auth.Principal{UserID: user.ID, Roles: user.Roles})
	if err != nil {
		return UserResponse{}, err
	}
	return applyAccountSummary(response, summary), nil
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
