package httpapi

import (
	"errors"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/gofiber/fiber/v3"
)

type Authorizer interface {
	Authorize(ctx fiber.Ctx, principal auth.Principal, resource, action string) error
}

type CasbinMiddleware struct{ service *authorization.Service }

func NewCasbinMiddleware(service *authorization.Service) *CasbinMiddleware {
	return &CasbinMiddleware{service: service}
}

func (m *CasbinMiddleware) Require(resource, action string) fiber.Handler {
	return func(c fiber.Ctx) error {
		principal, ok := c.Locals(principalLocal).(auth.Principal)
		if !ok {
			return NewAPIError(CodeAuthenticationRequired, "Authentication is required.", nil)
		}
		if err := m.service.Authorize(c.Context(), principal, resource, action); err != nil {
			if errors.Is(err, authorization.ErrForbidden) {
				return NewAPIError(CodeForbidden, "You are not allowed to perform this action.", nil)
			}
			return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
		}
		return c.Next()
	}
}
