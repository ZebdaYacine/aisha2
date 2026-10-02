package httpapi

import (
	"errors"
	"github.com/aisha-platform/aisha/apps/api/internal/features/order"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/gofiber/fiber/v3"
)

type OrderHandler struct {
	service   *order.Service
	validator *RequestValidator
}
type checkoutItemRequest struct {
	ProductID          string `json:"productId" validate:"required,uuid"`
	Quantity           int    `json:"quantity" validate:"gt=0,lte=100"`
	ExpectedPriceMinor *int64 `json:"expectedPriceMinor,omitempty" validate:"omitempty,gt=0"`
	Currency           string `json:"currency,omitempty" validate:"omitempty,len=3,uppercase"`
}
type checkoutRequest struct {
	AddressID string                `json:"addressId" validate:"required,uuid"`
	Items     []checkoutItemRequest `json:"items" validate:"required,min=1,max=50,dive"`
}
type checkoutResponse struct {
	Order order.Order `json:"order"`
}
type returnRequest struct {
	Reason string `json:"reason" validate:"required,min=3,max=1000"`
}

func NewOrderHandler(s *order.Service, v *RequestValidator) *OrderHandler {
	return &OrderHandler{service: s, validator: v}
}
func (h *OrderHandler) Checkout(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	var req checkoutRequest
	if e = c.Bind().Body(&req); e != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if e = h.validator.Validate(&req); e != nil {
		return e
	}
	items := make([]order.CartItem, len(req.Items))
	for i, v := range req.Items {
		items[i] = order.CartItem{ProductID: v.ProductID, Quantity: v.Quantity, ExpectedPriceMinor: v.ExpectedPriceMinor, ExpectedCurrency: v.Currency}
	}
	o, e := h.service.Checkout(c.Context(), p, req.AddressID, items, c.Get("Idempotency-Key"))
	if e != nil {
		return orderAPIError(e)
	}
	return c.Status(fiber.StatusCreated).JSON(checkoutResponse{Order: o})
}
func (h *OrderHandler) Payment(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	payment, e := h.service.Payment(c.Context(), p, c.Params("id"))
	if e != nil {
		return orderAPIError(e)
	}
	return c.JSON(payment)
}
func (h *OrderHandler) ConfirmPayment(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	o, e := h.service.ConfirmPayment(c.Context(), p, c.Params("id"), c.Get("Idempotency-Key"))
	if e != nil {
		return orderAPIError(e)
	}
	return c.JSON(o)
}
func (h *OrderHandler) List(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	page, size := queryPage(c)
	items, total, e := h.service.ListMine(c.Context(), p, page, size)
	if e != nil {
		return orderAPIError(e)
	}
	return c.JSON(PageDTO[order.Order]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}
func (h *OrderHandler) Get(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	o, e := h.service.GetMine(c.Context(), p, c.Params("id"))
	if e != nil {
		return orderAPIError(e)
	}
	return c.JSON(o)
}
func (h *OrderHandler) Cancel(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	o, e := h.service.Cancel(c.Context(), p, c.Params("id"))
	if e != nil {
		return orderAPIError(e)
	}
	return c.JSON(o)
}
func (h *OrderHandler) Seller(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	page, size := queryPage(c)
	items, total, e := h.service.ListSeller(c.Context(), p, page, size)
	if e != nil {
		return orderAPIError(e)
	}
	return c.JSON(PageDTO[order.SellerItem]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}
func (h *OrderHandler) Return(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	var req returnRequest
	if e = c.Bind().Body(&req); e != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if e = h.validator.Validate(&req); e != nil {
		return e
	}
	v, e := h.service.RecordReturn(c.Context(), p, c.Params("id"), req.Reason)
	if e != nil {
		return orderAPIError(e)
	}
	return c.Status(fiber.StatusCreated).JSON(v)
}
func orderAPIError(e error) error {
	switch {
	case errors.Is(e, order.ErrValidation):
		return NewAPIError(CodeValidationError, "The order request is invalid.", nil)
	case errors.Is(e, order.ErrNotFound):
		return NewAPIError(CodeResourceNotFound, "The order was not found.", nil)
	case errors.Is(e, order.ErrOutOfStock):
		return NewAPIError(CodeOutOfStock, "The requested quantity is no longer available.", nil)
	case errors.Is(e, order.ErrPriceChanged):
		return NewAPIError(CodeConflict, "The product price changed.", nil)
	case errors.Is(e, order.ErrPaymentAmountMismatch):
		return NewAPIError(CodePaymentAmountMismatch, "The payment amount does not match the order.", nil)
	case errors.Is(e, order.ErrIdempotencyConflict):
		return NewAPIError(CodeIdempotencyConflict, "The idempotency key cannot be reused for another request.", nil)
	case errors.Is(e, order.ErrInvalidTransition):
		return NewAPIError(CodeInvalidStateTransition, "The order cannot enter the requested state.", nil)
	case errors.Is(e, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	default:
		return WrapAPIError(e, CodeInternalError, "Unable to process the order.")
	}
}
