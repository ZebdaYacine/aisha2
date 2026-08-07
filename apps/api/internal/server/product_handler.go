package httpapi

import (
	"errors"
	"io"
	"strconv"

	"github.com/aisha-platform/aisha/apps/api/internal/features/product"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"github.com/gofiber/fiber/v3"
)

type ProductHandler struct {
	service   *product.Service
	validator *RequestValidator
	maxUpload int64
}

type productTranslationRequest struct {
	Locale          string `json:"locale" validate:"required,oneof=ar fr en es"`
	Name            string `json:"name" validate:"max=200"`
	Description     string `json:"description" validate:"max=5000"`
	Story           string `json:"story" validate:"max=5000"`
	CulturalContext string `json:"culturalContext" validate:"max=5000"`
}

type productRequest struct {
	CategoryID          string                      `json:"categoryId" validate:"required,uuid"`
	ProductType         string                      `json:"productType" validate:"required,oneof=ARTISAN_SPECIFIC STANDARD_TRADITIONAL"`
	PriceMinor          int64                       `json:"priceMinor" validate:"gt=0"`
	Currency            string                      `json:"currency" validate:"required,len=3,uppercase"`
	Materials           string                      `json:"materials" validate:"max=2000"`
	ProductionMethod    string                      `json:"productionMethod" validate:"max=2000"`
	IntendedUse         string                      `json:"intendedUse" validate:"max=1000"`
	Dimensions          string                      `json:"dimensions" validate:"max=500"`
	WeightGrams         *int                        `json:"weightGrams" validate:"omitempty,gt=0"`
	CountryOfOrigin     string                      `json:"countryOfOrigin" validate:"max=120"`
	RegionOfOrigin      string                      `json:"regionOfOrigin" validate:"max=120"`
	EcoFriendlyVerified bool                        `json:"ecoFriendlyVerified"`
	FairTradeVerified   bool                        `json:"fairTradeVerified"`
	MadeToOrderEligible bool                        `json:"madeToOrderEligible"`
	Translations        []productTranslationRequest `json:"translations" validate:"max=4,dive"`
}

type productMediaResponse struct {
	ID               string `json:"id"`
	MediaKind        string `json:"mediaKind"`
	OriginalFilename string `json:"originalFilename"`
	MediaType        string `json:"mediaType"`
	SizeBytes        int64  `json:"sizeBytes"`
	AltText          string `json:"altText"`
	SortOrder        int    `json:"sortOrder"`
	Visibility       string `json:"visibility"`
	URL              string `json:"url,omitempty"`
}

func NewProductHandler(service *product.Service, validator *RequestValidator, maxUpload int64) *ProductHandler {
	return &ProductHandler{service: service, validator: validator, maxUpload: maxUpload}
}

func (h *ProductHandler) Create(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	input, err := h.input(c)
	if err != nil {
		return err
	}
	item, err := h.service.Create(c.Context(), p, input)
	if err != nil {
		return productAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *ProductHandler) List(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("pageSize", "20"))
	items, total, err := h.service.List(c.Context(), p, c.Query("status"), page, size)
	if err != nil {
		return productAPIError(err)
	}
	return c.JSON(PageDTO[product.Product]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}

func (h *ProductHandler) Get(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	item, err := h.service.Get(c.Context(), p, c.Params("id"))
	if err != nil {
		return productAPIError(err)
	}
	return c.JSON(item)
}

func (h *ProductHandler) Update(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	input, err := h.input(c)
	if err != nil {
		return err
	}
	item, err := h.service.Update(c.Context(), p, c.Params("id"), input)
	if err != nil {
		return productAPIError(err)
	}
	return c.JSON(item)
}

func (h *ProductHandler) Submit(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	item, err := h.service.Submit(c.Context(), p, c.Params("id"))
	if err != nil {
		return productAPIError(err)
	}
	return c.JSON(item)
}

func (h *ProductHandler) UploadMedia(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	file, err := c.FormFile("file")
	if err != nil {
		return NewAPIError(CodeValidationError, "A media file is required.", FieldErrors{"file": "REQUIRED"})
	}
	if h.maxUpload <= 0 {
		h.maxUpload = 50 * 1024 * 1024
	}
	if file.Size > h.maxUpload {
		return NewAPIError(CodeFileTooLarge, "The file is too large.", nil)
	}
	reader, err := file.Open()
	if err != nil {
		return WrapAPIError(err, CodeInternalError, "Unable to read the uploaded file.")
	}
	defer reader.Close()
	data, err := io.ReadAll(io.LimitReader(reader, h.maxUpload+1))
	if err != nil {
		return WrapAPIError(err, CodeInternalError, "Unable to read the uploaded file.")
	}
	if int64(len(data)) > h.maxUpload {
		return NewAPIError(CodeFileTooLarge, "The file is too large.", nil)
	}
	item, err := h.service.UploadMedia(c.Context(), p, c.Params("id"), file.Filename, file.Header.Get("Content-Type"), data, c.FormValue("altText"))
	if err != nil {
		return productAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(productMediaDTO(item))
}

func (h *ProductHandler) DeleteMedia(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	if err = h.service.DeleteMedia(c.Context(), p, c.Params("id"), c.Params("mediaId")); err != nil {
		return productAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ProductHandler) input(c fiber.Ctx) (product.Input, error) {
	var request productRequest
	if err := c.Bind().Body(&request); err != nil {
		return product.Input{}, NewAPIError(CodeValidationError, "The request body is invalid.", nil)
	}
	if err := h.validator.Validate(&request); err != nil {
		return product.Input{}, err
	}
	translations := make([]product.Translation, len(request.Translations))
	for i, item := range request.Translations {
		translations[i] = product.Translation{Locale: item.Locale, Name: item.Name, Description: item.Description, Story: item.Story, CulturalContext: item.CulturalContext}
	}
	return product.Input{CategoryID: request.CategoryID, ProductType: request.ProductType, PriceMinor: request.PriceMinor, Currency: request.Currency, Materials: request.Materials, ProductionMethod: request.ProductionMethod, IntendedUse: request.IntendedUse, Dimensions: request.Dimensions, WeightGrams: request.WeightGrams, CountryOfOrigin: request.CountryOfOrigin, RegionOfOrigin: request.RegionOfOrigin, EcoFriendlyVerified: request.EcoFriendlyVerified, FairTradeVerified: request.FairTradeVerified, MadeToOrderEligible: request.MadeToOrderEligible, Translations: translations}, nil
}

func productMediaDTO(media product.Media) productMediaResponse {
	return productMediaResponse{ID: media.ID, MediaKind: media.MediaKind, OriginalFilename: media.OriginalFilename, MediaType: media.MediaType, SizeBytes: media.SizeBytes, AltText: media.AltText, SortOrder: media.SortOrder, Visibility: media.Visibility, URL: media.URL}
}

func productAPIError(err error) error {
	switch {
	case errors.Is(err, product.ErrValidation):
		return NewAPIError(CodeValidationError, "The product is invalid.", nil)
	case errors.Is(err, product.ErrNotFound), errors.Is(err, product.ErrMediaNotFound):
		return NewAPIError(CodeResourceNotFound, "The requested product resource was not found.", nil)
	case errors.Is(err, product.ErrArtisanNotApproved):
		return NewAPIError(CodeArtisanNotApproved, "An approved artisan profile is required.", nil)
	case errors.Is(err, product.ErrNotEditable):
		return NewAPIError(CodeProductNotEditable, "The product cannot be edited in its current state.", nil)
	case errors.Is(err, product.ErrInvalidTransition):
		return NewAPIError(CodeInvalidStateTransition, "The product cannot be submitted from its current state.", nil)
	case errors.Is(err, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	case errors.Is(err, storage.ErrUnsupportedFileType), errors.Is(err, storage.ErrInvalidFileSignature):
		return NewAPIError(CodeUnsupportedFileType, "The uploaded file type is not supported.", nil)
	case errors.Is(err, storage.ErrFileTooLarge):
		return NewAPIError(CodeFileTooLarge, "The file is too large.", nil)
	default:
		return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
	}
}
