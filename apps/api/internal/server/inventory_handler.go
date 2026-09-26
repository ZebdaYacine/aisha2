package httpapi

import (
	"errors"

	"github.com/aisha-platform/aisha/apps/api/internal/features/inventory"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/gofiber/fiber/v3"
)

type InventoryHandler struct {
	service   *inventory.Service
	validator *RequestValidator
}

type inventoryAdjustmentRequest struct {
	QuantityDelta int64  `json:"quantityDelta" validate:"required,ne=0"`
	Reason        string `json:"reason" validate:"required,min=3,max=2000"`
	ReferenceKey  string `json:"referenceKey" validate:"required,max=160"`
}

func NewInventoryHandler(service *inventory.Service, validator *RequestValidator) *InventoryHandler {
	return &InventoryHandler{service: service, validator: validator}
}

func (h *InventoryHandler) List(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	page, size := queryPage(c)
	items, total, err := h.service.List(c.Context(), principal, c.Query("workshopId"), page, size)
	if err != nil {
		return inventoryAPIError(err)
	}
	return c.JSON(PageDTO[inventory.Balance]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}

func (h *InventoryHandler) Adjust(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request inventoryAdjustmentRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The inventory adjustment is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	item, err := h.service.Adjust(c.Context(), principal, c.Params("productId"), inventory.AdjustmentInput{QuantityDelta: request.QuantityDelta, Reason: request.Reason, ReferenceKey: request.ReferenceKey})
	if err != nil {
		return inventoryAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func inventoryAPIError(err error) error {
	switch {
	case errors.Is(err, inventory.ErrValidation):
		return NewAPIError(CodeValidationError, "The inventory request is invalid.", nil)
	case errors.Is(err, inventory.ErrNotFound):
		return NewAPIError(CodeResourceNotFound, "The inventory resource was not found.", nil)
	case errors.Is(err, inventory.ErrDuplicate):
		return NewAPIError(CodeConflict, "The inventory reference is already used.", nil)
	case errors.Is(err, inventory.ErrInsufficientStock):
		return NewAPIError(CodeOutOfStock, "The available inventory cannot become negative.", nil)
	case errors.Is(err, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	default:
		return WrapAPIError(err, CodeInternalError, "Unable to process inventory request.")
	}
}
