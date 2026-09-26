package httpapi

import (
	"errors"

	"github.com/aisha-platform/aisha/apps/api/internal/features/cart"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/gofiber/fiber/v3"
)

type CartHandler struct {
	service   *cart.Service
	validator *RequestValidator
}
type cartItemRequest struct {
	ProductID string `json:"productId" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"gt=0,lte=100"`
}
type cartItemsRequest struct {
	Items []cartItemRequest `json:"items" validate:"max=100,dive"`
}

func NewCartHandler(s *cart.Service, v *RequestValidator) *CartHandler {
	return &CartHandler{service: s, validator: v}
}
func (h *CartHandler) List(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	items, e := h.service.List(c.Context(), p)
	if e != nil {
		return cartAPIError(e)
	}
	return c.JSON(PageDTO[cart.Item]{Items: items, Page: 1, PageSize: len(items), Total: len(items)})
}
func (h *CartHandler) Add(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	var req cartItemRequest
	if e = c.Bind().Body(&req); e != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if e = h.validator.Validate(&req); e != nil {
		return e
	}
	items, e := h.service.Add(c.Context(), p, cart.Input{ProductID: req.ProductID, Quantity: req.Quantity})
	if e != nil {
		return cartAPIError(e)
	}
	return c.Status(fiber.StatusCreated).JSON(items)
}
func (h *CartHandler) Set(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	var req cartItemRequest
	if e = c.Bind().Body(&req); e != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if e = h.validator.Validate(&req); e != nil {
		return e
	}
	items, e := h.service.Set(c.Context(), p, c.Params("productId"), req.Quantity)
	if e != nil {
		return cartAPIError(e)
	}
	return c.JSON(items)
}
func (h *CartHandler) Remove(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	items, e := h.service.Remove(c.Context(), p, c.Params("productId"))
	if e != nil {
		return cartAPIError(e)
	}
	return c.JSON(items)
}
func (h *CartHandler) Merge(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	var req cartItemsRequest
	if e = c.Bind().Body(&req); e != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if e = h.validator.Validate(&req); e != nil {
		return e
	}
	inputs := make([]cart.Input, len(req.Items))
	for i, v := range req.Items {
		inputs[i] = cart.Input{ProductID: v.ProductID, Quantity: v.Quantity}
	}
	items, e := h.service.Merge(c.Context(), p, inputs)
	if e != nil {
		return cartAPIError(e)
	}
	return c.JSON(items)
}
func cartAPIError(e error) error {
	switch {
	case errors.Is(e, cart.ErrValidation):
		return NewAPIError(CodeValidationError, "The cart request is invalid.", nil)
	case errors.Is(e, cart.ErrUnavailable):
		return NewAPIError(CodeConflict, "The product is no longer available.", nil)
	case errors.Is(e, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	default:
		return WrapAPIError(e, CodeInternalError, "Unable to process the cart.")
	}
}
