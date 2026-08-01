package httpapi

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

type ErrorCode string

const (
	CodeValidationError         ErrorCode = "VALIDATION_ERROR"
	CodeInvalidCredentials      ErrorCode = "INVALID_CREDENTIALS"
	CodeAuthenticationRequired  ErrorCode = "AUTHENTICATION_REQUIRED"
	CodeForbidden               ErrorCode = "FORBIDDEN"
	CodeResourceNotFound        ErrorCode = "RESOURCE_NOT_FOUND"
	CodeConflict                ErrorCode = "CONFLICT"
	CodeArtisanNotApproved      ErrorCode = "ARTISAN_NOT_APPROVED"
	CodeProductNotEditable      ErrorCode = "PRODUCT_NOT_EDITABLE"
	CodeProductNotSellable      ErrorCode = "PRODUCT_NOT_SELLABLE"
	CodeInvalidStateTransition  ErrorCode = "INVALID_STATE_TRANSITION"
	CodeOutOfStock              ErrorCode = "OUT_OF_STOCK"
	CodeReservationExpired      ErrorCode = "RESERVATION_EXPIRED"
	CodeIdempotencyConflict     ErrorCode = "IDEMPOTENCY_CONFLICT"
	CodePaymentAmountMismatch   ErrorCode = "PAYMENT_AMOUNT_MISMATCH"
	CodeInvalidWebhookSignature ErrorCode = "INVALID_WEBHOOK_SIGNATURE"
	CodeDuplicateWebhook        ErrorCode = "DUPLICATE_WEBHOOK"
	CodeUnsupportedFileType     ErrorCode = "UNSUPPORTED_FILE_TYPE"
	CodeFileTooLarge            ErrorCode = "FILE_TOO_LARGE"
	CodeRateLimited             ErrorCode = "RATE_LIMITED"
	CodeInternalError           ErrorCode = "INTERNAL_ERROR"
)

type FieldErrors map[string]string

type APIError struct {
	Code    ErrorCode
	Message string
	Fields  FieldErrors
	cause   error
}

func NewAPIError(code ErrorCode, message string, fields FieldErrors) *APIError {
	return &APIError{Code: code, Message: message, Fields: fields}
}

func WrapAPIError(cause error, code ErrorCode, message string) *APIError {
	return &APIError{Code: code, Message: message, cause: cause}
}

func (e *APIError) Error() string { return e.Message }
func (e *APIError) Unwrap() error { return e.cause }

type ErrorResponse struct {
	Error ErrorDTO `json:"error"`
}

type ErrorDTO struct {
	Code      ErrorCode   `json:"code"`
	Message   string      `json:"message"`
	Fields    FieldErrors `json:"fields,omitempty"`
	RequestID string      `json:"requestId"`
}

func errorHandler(c fiber.Ctx, err error) error {
	requestID := c.GetRespHeader(fiber.HeaderXRequestID)
	apiErr := normalizeError(err)
	status := statusForCode(apiErr.Code)
	if status >= fiber.StatusInternalServerError {
		slog.Error("request failed", "request_id", requestID, "error", err)
	} else {
		slog.Warn("request rejected", "request_id", requestID, "code", apiErr.Code, "error", err)
	}
	return c.Status(status).JSON(ErrorResponse{Error: ErrorDTO{Code: apiErr.Code, Message: apiErr.Message, Fields: apiErr.Fields, RequestID: requestID}})
}

func normalizeError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) && fiberErr.Code == fiber.StatusNotFound {
		return NewAPIError(CodeResourceNotFound, "The requested resource was not found.", nil)
	}
	return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
}

func statusForCode(code ErrorCode) int {
	switch code {
	case CodeValidationError, CodeUnsupportedFileType, CodeFileTooLarge:
		return fiber.StatusBadRequest
	case CodeInvalidCredentials, CodeAuthenticationRequired, CodeInvalidWebhookSignature:
		return fiber.StatusUnauthorized
	case CodeForbidden, CodeArtisanNotApproved:
		return fiber.StatusForbidden
	case CodeResourceNotFound:
		return fiber.StatusNotFound
	case CodeConflict, CodeProductNotEditable, CodeProductNotSellable, CodeInvalidStateTransition, CodeOutOfStock, CodeReservationExpired, CodeIdempotencyConflict, CodePaymentAmountMismatch, CodeDuplicateWebhook:
		return fiber.StatusConflict
	case CodeRateLimited:
		return fiber.StatusTooManyRequests
	default:
		return fiber.StatusInternalServerError
	}
}
