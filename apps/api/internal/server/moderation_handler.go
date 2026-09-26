package httpapi

import (
	"errors"
	"github.com/aisha-platform/aisha/apps/api/internal/features/moderation"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/gofiber/fiber/v3"
)

type ModerationHandler struct {
	service   *moderation.Service
	validator *RequestValidator
}
type moderationDecisionRequest struct {
	Action string `json:"action" validate:"required"`
	Reason string `json:"reason" validate:"max=1000"`
}

func NewModerationHandler(s *moderation.Service, v *RequestValidator) *ModerationHandler {
	return &ModerationHandler{service: s, validator: v}
}
func (h *ModerationHandler) Queue(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	page, size := queryPage(c)
	items, total, e := h.service.ListQueue(c.Context(), p, c.Query("status"), page, size)
	if e != nil {
		return moderationAPIError(e)
	}
	return c.JSON(PageDTO[moderation.QueueItem]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}
func (h *ModerationHandler) Decide(c fiber.Ctx) error {
	p, e := customerPrincipal(c)
	if e != nil {
		return e
	}
	var req moderationDecisionRequest
	if e = c.Bind().Body(&req); e != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if e = h.validator.Validate(&req); e != nil {
		return e
	}
	item, e := h.service.Decide(c.Context(), p, moderation.DecisionInput{ID: c.Params("id"), Action: req.Action, Reason: req.Reason})
	if e != nil {
		return moderationAPIError(e)
	}
	return c.JSON(item)
}
func moderationAPIError(e error) error {
	switch {
	case errors.Is(e, moderation.ErrValidation):
		return NewAPIError(CodeValidationError, "The moderation request is invalid.", nil)
	case errors.Is(e, moderation.ErrNotFound):
		return NewAPIError(CodeResourceNotFound, "The moderation target was not found.", nil)
	case errors.Is(e, moderation.ErrInvalidTransition), errors.Is(e, moderation.ErrActivationNotReady):
		return NewAPIError(CodeInvalidStateTransition, "The product cannot enter the requested state.", nil)
	case errors.Is(e, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	default:
		return WrapAPIError(e, CodeInternalError, "Unable to process moderation decision.")
	}
}
