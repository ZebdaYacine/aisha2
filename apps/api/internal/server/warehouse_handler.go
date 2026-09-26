package httpapi

import (
	"errors"
	"io"

	"github.com/aisha-platform/aisha/apps/api/internal/features/warehouse"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"github.com/gofiber/fiber/v3"
)

type WarehouseHandler struct {
	service   *warehouse.Service
	validator *RequestValidator
	maxUpload int64
}

type receptionRequest struct {
	ProductID        string `json:"productId" validate:"required,uuid"`
	ReceivedQuantity int64  `json:"receivedQuantity" validate:"gt=0"`
	ReferenceKey     string `json:"referenceKey" validate:"required,max=160"`
	SupplierName     string `json:"supplierName" validate:"max=200"`
	ParcelReference  string `json:"parcelReference" validate:"max=200"`
	Notes            string `json:"notes" validate:"max=4000"`
}

type inspectionRequest struct {
	AcceptedQuantity    int64  `json:"acceptedQuantity" validate:"gte=0"`
	RejectedQuantity    int64  `json:"rejectedQuantity" validate:"gte=0"`
	QuarantinedQuantity int64  `json:"quarantinedQuantity" validate:"gte=0"`
	DamagedQuantity     int64  `json:"damagedQuantity" validate:"gte=0"`
	Reason              string `json:"reason" validate:"required,min=3,max=2000"`
}

func NewWarehouseHandler(service *warehouse.Service, validator *RequestValidator, maxUpload int64) *WarehouseHandler {
	return &WarehouseHandler{service: service, validator: validator, maxUpload: maxUpload}
}

func (h *WarehouseHandler) List(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	page, size := queryPage(c)
	items, total, err := h.service.ListReceptions(c.Context(), principal, c.Query("status"), page, size)
	if err != nil {
		return warehouseAPIError(err)
	}
	return c.JSON(PageDTO[warehouse.Reception]{Items: items, Page: max(page, 1), PageSize: min(max(size, 1), 100), Total: total})
}

func (h *WarehouseHandler) Create(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request receptionRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The reception request is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	item, err := h.service.CreateReception(c.Context(), principal, warehouse.ReceptionInput{ProductID: request.ProductID, ReceivedQuantity: request.ReceivedQuantity, ReferenceKey: request.ReferenceKey, SupplierName: request.SupplierName, ParcelReference: request.ParcelReference, Notes: request.Notes})
	if err != nil {
		return warehouseAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *WarehouseHandler) Inspect(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	var request inspectionRequest
	if err = c.Bind().Body(&request); err != nil {
		return NewAPIError(CodeValidationError, "The inspection request is invalid.", nil)
	}
	if err = h.validator.Validate(&request); err != nil {
		return err
	}
	item, err := h.service.Inspect(c.Context(), principal, c.Params("id"), warehouse.InspectionInput{AcceptedQuantity: request.AcceptedQuantity, RejectedQuantity: request.RejectedQuantity, QuarantinedQuantity: request.QuarantinedQuantity, DamagedQuantity: request.DamagedQuantity, Reason: request.Reason})
	if err != nil {
		return warehouseAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *WarehouseHandler) UploadEvidence(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	file, err := c.FormFile("file")
	if err != nil {
		return NewAPIError(CodeValidationError, "An evidence file is required.", FieldErrors{"file": "REQUIRED"})
	}
	if h.maxUpload > 0 && file.Size > h.maxUpload {
		return NewAPIError(CodeFileTooLarge, "The evidence file is too large.", nil)
	}
	reader, err := file.Open()
	if err != nil {
		return WrapAPIError(err, CodeInternalError, "Unable to read the evidence file.")
	}
	defer reader.Close()
	maxUpload := h.maxUpload
	if maxUpload <= 0 {
		maxUpload = 50 * 1024 * 1024
	}
	data, err := io.ReadAll(io.LimitReader(reader, maxUpload+1))
	if err != nil {
		return WrapAPIError(err, CodeInternalError, "Unable to read the evidence file.")
	}
	if int64(len(data)) > maxUpload {
		return NewAPIError(CodeFileTooLarge, "The evidence file is too large.", nil)
	}
	item, err := h.service.UploadEvidence(c.Context(), principal, c.Params("id"), file.Filename, file.Header.Get("Content-Type"), data)
	if err != nil {
		return warehouseAPIError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *WarehouseHandler) Evidence(c fiber.Ctx) error {
	principal, err := customerPrincipal(c)
	if err != nil {
		return err
	}
	items, err := h.service.Evidence(c.Context(), principal, c.Params("id"))
	if err != nil {
		return warehouseAPIError(err)
	}
	return c.JSON(struct {
		Items []warehouse.Evidence `json:"items"`
	}{Items: items})
}

func warehouseAPIError(err error) error {
	switch {
	case errors.Is(err, warehouse.ErrValidation), errors.Is(err, warehouse.ErrQuantityMismatch), errors.Is(err, warehouse.ErrEvidenceRequired):
		return NewAPIError(CodeValidationError, "The warehouse request is invalid.", nil)
	case errors.Is(err, warehouse.ErrNotFound):
		return NewAPIError(CodeResourceNotFound, "The warehouse resource was not found.", nil)
	case errors.Is(err, warehouse.ErrDuplicate):
		return NewAPIError(CodeConflict, "The reception reference is already used for different stock.", nil)
	case errors.Is(err, warehouse.ErrInvalidTransition), errors.Is(err, warehouse.ErrAlreadyInspected):
		return NewAPIError(CodeInvalidStateTransition, "The warehouse resource cannot enter the requested state.", nil)
	case errors.Is(err, storage.ErrUnsupportedFileType), errors.Is(err, storage.ErrInvalidFileSignature):
		return NewAPIError(CodeUnsupportedFileType, "The evidence file type is not allowed.", nil)
	case errors.Is(err, storage.ErrFileTooLarge):
		return NewAPIError(CodeFileTooLarge, "The evidence file is too large.", nil)
	case errors.Is(err, authorization.ErrForbidden):
		return NewAPIError(CodeForbidden, "Access is forbidden.", nil)
	default:
		return WrapAPIError(err, CodeInternalError, "Unable to process warehouse request.")
	}
}
