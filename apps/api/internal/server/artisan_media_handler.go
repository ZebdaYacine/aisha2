package httpapi

import (
	"errors"
	"io"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"github.com/gofiber/fiber/v3"
)

type ArtisanMediaHandler struct {
	service   *artisan.MediaService
	maxUpload int64
}

func NewArtisanMediaHandler(service *artisan.MediaService, maxUpload int64) *ArtisanMediaHandler {
	return &ArtisanMediaHandler{service: service, maxUpload: maxUpload}
}

func (h *ArtisanMediaHandler) Documents(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	items, err := h.service.OwnDocuments(c.Context(), p)
	if err != nil {
		return artisanMediaAPIError(err)
	}
	return c.JSON(items)
}

func (h *ArtisanMediaHandler) UploadDocument(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	file, data, err := uploadedBytes(c, "file", h.maxUpload)
	if err != nil {
		return err
	}
	item, err := h.service.UploadDocument(c.Context(), p, c.FormValue("documentType"), file.Filename, file.ContentType, data)
	if err != nil {
		return artisanMediaAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *ArtisanMediaHandler) Media(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	items, err := h.service.OwnMedia(c.Context(), p)
	if err != nil {
		return artisanMediaAPIError(err)
	}
	response := make([]artisanMediaResponse, len(items))
	for i, item := range items {
		response[i] = artisanMediaDTO(item)
	}
	return c.JSON(response)
}

func (h *ArtisanMediaHandler) UploadMedia(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	file, data, err := uploadedBytes(c, "file", h.maxUpload)
	if err != nil {
		return err
	}
	filename := c.FormValue("mediaName")
	if filename == "" {
		filename = file.Filename
	}
	item, err := h.service.UploadProfileMediaWithOptions(c.Context(), p, c.FormValue("mediaKind"), filename, file.ContentType, data)
	if err != nil {
		return artisanMediaAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(artisanMediaDTO(item))
}

func (h *ArtisanMediaHandler) ReplaceMedia(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	file, data, err := uploadedBytes(c, "file", h.maxUpload)
	if err != nil {
		return err
	}
	filename := c.FormValue("mediaName")
	if filename == "" {
		filename = file.Filename
	}
	item, err := h.service.ReplaceProfileMedia(c.Context(), p, c.Params("id"), c.FormValue("mediaKind"), filename, file.ContentType, data)
	if err != nil {
		return artisanMediaAPIError(err)
	}
	return c.JSON(artisanMediaDTO(item))
}

func (h *ArtisanMediaHandler) DeleteMedia(c fiber.Ctx) error {
	p, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	if err = h.service.DeleteProfileMedia(c.Context(), p, c.Params("id")); err != nil {
		return artisanMediaAPIError(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func uploadedBytes(c fiber.Ctx, field string, maxBytes int64) (*multipartFile, []byte, error) {
	file, err := c.FormFile(field)
	if err != nil {
		return nil, nil, NewAPIError(CodeValidationError, "A file is required.", FieldErrors{field: "REQUIRED"})
	}
	if maxBytes <= 0 {
		maxBytes = 10 * 1024 * 1024
	}
	if file.Size > maxBytes {
		return nil, nil, NewAPIError(CodeFileTooLarge, "The file is too large.", nil)
	}
	reader, err := file.Open()
	if err != nil {
		return nil, nil, WrapAPIError(err, CodeInternalError, "Unable to read the uploaded file.")
	}
	defer reader.Close()
	data, err := io.ReadAll(io.LimitReader(reader, maxBytes+1))
	if err != nil {
		return nil, nil, WrapAPIError(err, CodeInternalError, "Unable to read the uploaded file.")
	}
	if int64(len(data)) > maxBytes {
		return nil, nil, NewAPIError(CodeFileTooLarge, "The file is too large.", nil)
	}
	return &multipartFile{Filename: file.Filename, ContentType: file.Header.Get("Content-Type")}, data, nil
}

type multipartFile struct {
	Filename    string
	ContentType string
}

type artisanMediaResponse struct {
	ID               string    `json:"id"`
	MediaKind        string    `json:"mediaKind"`
	OriginalFilename string    `json:"originalFilename"`
	MediaType        string    `json:"mediaType"`
	SizeBytes        int64     `json:"sizeBytes"`
	SortOrder        int       `json:"sortOrder"`
	Visibility       string    `json:"visibility"`
	CreatedAt        time.Time `json:"createdAt"`
	URL              string    `json:"url,omitempty"`
}

func artisanMediaDTO(item artisan.Media) artisanMediaResponse {
	return artisanMediaResponse{ID: item.ID, MediaKind: item.MediaKind, OriginalFilename: item.OriginalFilename, MediaType: item.MediaType, SizeBytes: item.SizeBytes, SortOrder: item.SortOrder, Visibility: item.Visibility, CreatedAt: item.CreatedAt, URL: item.URL}
}

func artisanMediaAPIError(err error) error {
	switch {
	case errors.Is(err, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	case errors.Is(err, artisan.ErrNotFound):
		return NewAPIError(CodeResourceNotFound, "The artisan profile was not found.", nil)
	case errors.Is(err, artisan.ErrInvalidTransition):
		return NewAPIError(CodeInvalidStateTransition, "The artisan profile cannot be changed in its current state.", nil)
	case errors.Is(err, artisan.ErrValidation):
		return NewAPIError(CodeValidationError, "The uploaded artisan media is invalid.", nil)
	case errors.Is(err, storage.ErrUnsupportedFileType), errors.Is(err, storage.ErrInvalidFileSignature):
		return NewAPIError(CodeUnsupportedFileType, "The uploaded file type is not supported.", nil)
	case errors.Is(err, storage.ErrFileTooLarge):
		return NewAPIError(CodeFileTooLarge, "The file is too large.", nil)
	default:
		return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
	}
}
