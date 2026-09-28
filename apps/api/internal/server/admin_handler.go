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
type userStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=ACTIVE SUSPENDED DISABLED"`
	Reason string `json:"reason" validate:"max=1000"`
}
type userRequest struct {
	Email       string   `json:"email" validate:"required,email,max=254"`
	Phone       string   `json:"phone" validate:"max=32"`
	DisplayName string   `json:"displayName" validate:"required,min=2,max=120"`
	Password    string   `json:"password" validate:"max=128"`
	Roles       []string `json:"roles" validate:"required,min=1,max=10,dive,min=1,max=64"`
}
type categoryRequest struct {
	Slug                   string            `json:"slug" validate:"required,min=2,max=80"`
	DisplayName            string            `json:"displayName" validate:"required,min=2,max=160"`
	Translations           map[string]string `json:"translations" validate:"required"`
	BenefitRateBasisPoints int64             `json:"benefitRateBasisPoints" validate:"gte=0,lte=10000"`
	IsActive               bool              `json:"isActive"`
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

func (h *AdminHandler) CreateUser(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request userRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The user request is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	item, err := h.service.CreateUser(c.Context(), p, admin.UserInput{Email: request.Email, Phone: request.Phone, DisplayName: request.DisplayName, Password: request.Password, Roles: request.Roles})
	if err != nil {
		return adminAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *AdminHandler) UpdateUser(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request userRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The user request is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	item, err := h.service.UpdateUser(c.Context(), p, c.Params("id"), admin.UserInput{Email: request.Email, Phone: request.Phone, DisplayName: request.DisplayName, Password: request.Password, Roles: request.Roles})
	if err != nil {
		return adminAPIError(err)
	}
	return c.JSON(item)
}

func (h *AdminHandler) UserStatus(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request userStatusRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	item, err := h.service.SetUserStatus(c.Context(), p, c.Params("id"), request.Status, request.Reason)
	if err != nil {
		return adminAPIError(err)
	}
	return c.JSON(item)
}

func (h *AdminHandler) Categories(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	page, size := queryPage(c)
	items, total, err := h.service.ListCategories(c.Context(), p, page, size)
	if err != nil {
		return adminAPIError(err)
	}
	return c.JSON(PageDTO[admin.Category]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}

func (h *AdminHandler) Orders(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	page, size := queryPage(c)
	items, total, err := h.service.ListOrders(c.Context(), p, page, size)
	if err != nil {
		return adminAPIError(err)
	}
	return c.JSON(PageDTO[admin.Order]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}

func (h *AdminHandler) CreateCategory(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request categoryRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The category request is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	item, err := h.service.CreateCategory(c.Context(), p, admin.CategoryInput{Slug: request.Slug, DisplayName: request.DisplayName, Translations: request.Translations, BenefitRateBasisPoints: request.BenefitRateBasisPoints, IsActive: request.IsActive})
	if err != nil {
		return adminAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *AdminHandler) UpdateCategory(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request categoryRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The category request is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	item, err := h.service.UpdateCategory(c.Context(), p, c.Params("id"), admin.CategoryInput{Slug: request.Slug, DisplayName: request.DisplayName, Translations: request.Translations, BenefitRateBasisPoints: request.BenefitRateBasisPoints, IsActive: request.IsActive})
	if err != nil {
		return adminAPIError(err)
	}
	return c.JSON(item)
}

func (h *AdminHandler) DeleteCategory(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	if err = h.service.DeleteCategory(c.Context(), p, c.Params("id")); err != nil {
		return adminAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
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

func (h *AdminHandler) UserMedia(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	page, size := queryPage(c)
	items, total, err := h.service.ListUserMedia(c.Context(), p, page, size)
	if err != nil {
		return adminAPIError(err)
	}
	return c.JSON(PageDTO[admin.UserMedia]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}

func (h *AdminHandler) ProductMedia(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	page, size := queryPage(c)
	items, total, err := h.service.ListProductMedia(c.Context(), p, page, size)
	if err != nil {
		return adminAPIError(err)
	}
	return c.JSON(PageDTO[admin.ProductMedia]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}

func (h *AdminHandler) DeleteUserMedia(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	if err = h.service.DeleteUserMedia(c.Context(), p, c.Params("id")); err != nil {
		return adminAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AdminHandler) DeleteProductMedia(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	if err = h.service.DeleteProductMedia(c.Context(), p, c.Params("id")); err != nil {
		return adminAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
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
		return NewAPIError(CodeResourceNotFound, "The requested administration resource was not found.", nil)
	case errors.Is(err, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	default:
		return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
	}
}
