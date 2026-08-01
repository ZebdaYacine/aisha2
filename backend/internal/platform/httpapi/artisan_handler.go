package httpapi

import (
	"errors"
	"strconv"
	"strings"

	"github.com/aisha-platform/aisha/backend/internal/artisan"
	"github.com/aisha-platform/aisha/backend/internal/authorization"
	"github.com/gofiber/fiber/v3"
)

type ArtisanHandler struct {
	service   *artisan.Service
	validator *RequestValidator
}
type artisanTranslationRequest struct {
	Locale    string `json:"locale" validate:"required,oneof=ar fr en es"`
	Biography string `json:"biography" validate:"max=5000"`
}
type artisanApplicationRequest struct {
	PublicDisplayName string                      `json:"publicDisplayName" validate:"required,min=2,max=120"`
	InternalName      string                      `json:"internalName" validate:"max=160"`
	WorkshopName      string                      `json:"workshopName" validate:"max=160"`
	Wilaya            string                      `json:"wilaya" validate:"required,max=100"`
	Location          string                      `json:"location" validate:"max=200"`
	ContactEmail      string                      `json:"contactEmail" validate:"omitempty,email,max=254"`
	ContactPhone      string                      `json:"contactPhone" validate:"omitempty,min=6,max=32"`
	ContactVisibility string                      `json:"contactVisibility" validate:"required,oneof=PRIVATE PUBLIC"`
	CategoryIDs       []string                    `json:"categoryIds" validate:"required,min=1,dive,uuid"`
	Translations      []artisanTranslationRequest `json:"translations" validate:"max=4,dive"`
}
type decisionRequest struct {
	Reason string `json:"reason" validate:"max=1000"`
}
type artisanApplicationResponse struct {
	ID                string                `json:"id"`
	PublicDisplayName string                `json:"publicDisplayName"`
	InternalName      string                `json:"internalName"`
	WorkshopName      string                `json:"workshopName"`
	Wilaya            string                `json:"wilaya"`
	Location          string                `json:"location"`
	ContactEmail      string                `json:"contactEmail"`
	ContactPhone      string                `json:"contactPhone"`
	ContactVisibility string                `json:"contactVisibility"`
	Status            string                `json:"status"`
	ReviewReason      string                `json:"reviewReason,omitempty"`
	CategoryIDs       []string              `json:"categoryIds"`
	Translations      []artisan.Translation `json:"translations"`
}

func NewArtisanHandler(service *artisan.Service, validator *RequestValidator) *ArtisanHandler {
	return &ArtisanHandler{service: service, validator: validator}
}
func (h *ArtisanHandler) input(c fiber.Ctx) (artisan.ApplicationInput, error) {
	var r artisanApplicationRequest
	if err := c.Bind().Body(&r); err != nil {
		return artisan.ApplicationInput{}, NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if err := h.validator.Validate(&r); err != nil {
		return artisan.ApplicationInput{}, err
	}
	translations := make([]artisan.Translation, len(r.Translations))
	for i, t := range r.Translations {
		translations[i] = artisan.Translation{Locale: t.Locale, Biography: t.Biography}
	}
	return artisan.ApplicationInput{PublicDisplayName: r.PublicDisplayName, InternalName: r.InternalName, WorkshopName: r.WorkshopName, Wilaya: r.Wilaya, Location: r.Location, ContactEmail: r.ContactEmail, ContactPhone: r.ContactPhone, ContactVisibility: r.ContactVisibility, CategoryIDs: r.CategoryIDs, Translations: translations}, nil
}
func (h *ArtisanHandler) Submit(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	input, err := h.input(c)
	if err != nil {
		return err
	}
	result, err := h.service.Submit(c.Context(), p, input)
	if err != nil {
		return artisanAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(artisanApplicationDTO(result))
}
func (h *ArtisanHandler) Mine(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	result, err := h.service.Mine(c.Context(), p)
	if err != nil {
		return artisanAPIError(err)
	}
	return c.JSON(artisanApplicationDTO(result))
}
func (h *ArtisanHandler) UpdateProfile(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	input, err := h.input(c)
	if err != nil {
		return err
	}
	result, err := h.service.UpdateProfile(c.Context(), p, input)
	if err != nil {
		return artisanAPIError(err)
	}
	return c.JSON(artisanApplicationDTO(result))
}
func (h *ArtisanHandler) List(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("pageSize", "20"))
	items, total, err := h.service.List(c.Context(), p, c.Query("status"), page, size)
	if err != nil {
		return artisanAPIError(err)
	}
	response := make([]artisanApplicationResponse, len(items))
	for i, item := range items {
		response[i] = artisanApplicationDTO(item)
	}
	return c.JSON(PageDTO[artisanApplicationResponse]{Items: response, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}
func (h *ArtisanHandler) Decide(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request decisionRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	decisions := map[string]string{"approve": "APPROVED", "request-changes": "CHANGES_REQUESTED", "reject": "REJECTED"}
	decision := decisions[strings.ToLower(c.Params("decision"))]
	result, err := h.service.Decide(c.Context(), p, c.Params("id"), decision, request.Reason)
	if err != nil {
		return artisanAPIError(err)
	}
	return c.JSON(artisanApplicationDTO(result))
}
func (h *ArtisanHandler) Documents(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	items, err := h.service.Documents(c.Context(), p, c.Params("id"))
	if err != nil {
		return artisanAPIError(err)
	}
	return c.JSON(items)
}
func artisanApplicationDTO(a artisan.Application) artisanApplicationResponse {
	return artisanApplicationResponse{ID: a.ID, PublicDisplayName: a.PublicDisplayName, InternalName: a.InternalName, WorkshopName: a.WorkshopName, Wilaya: a.Wilaya, Location: a.Location, ContactEmail: a.ContactEmail, ContactPhone: a.ContactPhone, ContactVisibility: a.ContactVisibility, Status: a.Status, ReviewReason: a.ReviewReason, CategoryIDs: a.CategoryIDs, Translations: a.Translations}
}
func artisanAPIError(err error) error {
	switch {
	case errors.Is(err, artisan.ErrValidation):
		return NewAPIError(CodeValidationError, "The request is invalid.", nil)
	case errors.Is(err, artisan.ErrNotFound):
		return NewAPIError(CodeResourceNotFound, "The requested resource was not found.", nil)
	case errors.Is(err, artisan.ErrInvalidTransition):
		return NewAPIError(CodeInvalidStateTransition, "The artisan application cannot be changed from its current state.", nil)
	case errors.Is(err, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	default:
		return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
	}
}
