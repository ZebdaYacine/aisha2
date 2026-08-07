package httpapi

import (
	"errors"
	"strconv"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/admin"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/gofiber/fiber/v3"
)

type AdminHandler struct {
	service   *admin.Service
	validator *RequestValidator
}

type rolesRequest struct {
	Roles []string `json:"roles" validate:"required,min=1,max=10,dive,min=1,max=64"`
}

func NewAdminHandler(service *admin.Service, validator *RequestValidator) *AdminHandler {
	return &AdminHandler{service: service, validator: validator}
}

func (h *AdminHandler) Users(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	page, size := queryPage(c)
	items, total, err := h.service.ListUsers(c.Context(), p, c.Query("role"), c.Query("status"), page, size)
	if err != nil {
		return adminAPIError(err)
	}
	return c.JSON(PageDTO[admin.User]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}

func (h *AdminHandler) Roles(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request rolesRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	item, err := h.service.SetRoles(c.Context(), p, c.Params("id"), request.Roles)
	if err != nil {
		return adminAPIError(err)
	}
	return c.JSON(item)
}

func (h *AdminHandler) Audit(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	filter := admin.AuditFilter{ActorUserID: c.Query("actorUserId"), TargetType: c.Query("targetType"), TargetID: c.Query("targetId"), EventType: c.Query("eventType"), CorrelationID: c.Query("correlationId")}
	if value := c.Query("from"); value != "" {
		parsed, parseErr := time.Parse(time.RFC3339, value)
		if parseErr != nil {
			return NewAPIError(CodeValidationError, "The audit date filter is invalid.", FieldErrors{"from": "INVALID"})
		}
		filter.From = &parsed
	}
	if value := c.Query("to"); value != "" {
		parsed, parseErr := time.Parse(time.RFC3339, value)
		if parseErr != nil {
			return NewAPIError(CodeValidationError, "The audit date filter is invalid.", FieldErrors{"to": "INVALID"})
		}
		filter.To = &parsed
	}
	page, size := queryPage(c)
	items, total, err := h.service.AuditEvents(c.Context(), p, filter, page, size)
	if err != nil {
		return adminAPIError(err)
	}
	return c.JSON(PageDTO[admin.AuditEvent]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}

func queryPage(c fiber.Ctx) (int, int) {
	page := 1
	size := 20
	if value := c.Query("page"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			page = parsed
		}
	}
	if value := c.Query("pageSize"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			size = parsed
		}
	}
	return page, size
}

func adminAPIError(err error) error {
	switch {
	case errors.Is(err, admin.ErrValidation):
		return NewAPIError(CodeValidationError, "The administration request is invalid.", nil)
	case errors.Is(err, admin.ErrNotFound):
		return NewAPIError(CodeResourceNotFound, "The requested user was not found.", nil)
	case errors.Is(err, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	default:
		return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
	}
}
