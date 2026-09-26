package httpapi

import (
	"errors"
	"github.com/aisha-platform/aisha/apps/api/internal/features/wishlist"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/gofiber/fiber/v3"
)

type WishlistHandler struct{ service *wishlist.Service }

func NewWishlistHandler(s *wishlist.Service) *WishlistHandler { return &WishlistHandler{service: s} }
func (h *WishlistHandler) List(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	items, e := h.service.List(c.Context(), p)
	if e != nil {
		return wishlistAPIError(e)
	}
	return c.JSON(items)
}
func (h *WishlistHandler) Add(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	items, e := h.service.Add(c.Context(), p, c.Params("productId"))
	if e != nil {
		return wishlistAPIError(e)
	}
	return c.Status(fiber.StatusCreated).JSON(items)
}
func (h *WishlistHandler) Remove(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	items, e := h.service.Remove(c.Context(), p, c.Params("productId"))
	if e != nil {
		return wishlistAPIError(e)
	}
	return c.JSON(items)
}
func wishlistAPIError(e error) error {
	switch {
	case errors.Is(e, wishlist.ErrValidation):
		return NewAPIError(CodeValidationError, "The wishlist request is invalid.", nil)
	case errors.Is(e, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	default:
		return WrapAPIError(e, CodeInternalError, "Unable to process the wishlist.")
	}
}
