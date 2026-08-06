package error_handler

import (
	"errors"
	"log/slog"

	"github.com/aisha-platform/aisha/apps/api/internal/pkg/apperror"
	"github.com/gofiber/fiber/v3"
)

type ErrorCode = apperror.Code
type FieldErrors = apperror.FieldErrors
type APIError = apperror.Error

func NewAPIError(code ErrorCode, message string, fields FieldErrors) *APIError {
	return apperror.New(code, message, fields)
}

func WrapAPIError(cause error, code ErrorCode, message string) *APIError {
	return apperror.Wrap(cause, code, message)
}

type ErrorResponse struct {
	Error ErrorDTO `json:"error"`
}

type ErrorDTO struct {
	Code      ErrorCode   `json:"code"`
	Message   string      `json:"message"`
	Fields    FieldErrors `json:"fields,omitempty"`
	RequestID string      `json:"requestId"`
}

func Handler(c fiber.Ctx, err error) error {
	requestID := c.GetRespHeader(fiber.HeaderXRequestID)
	apiErr := normalize(err)
	status := StatusForCode(apiErr.Code)
	if status >= fiber.StatusInternalServerError {
		logErr := errors.Unwrap(apiErr)
		if logErr == nil {
			logErr = err
		}
		slog.Error("request failed", "request_id", requestID, "method", c.Method(), "path", c.Path(), "status", status, "error", logErr)
	} else {
		slog.Warn("request rejected", "request_id", requestID, "method", c.Method(), "path", c.Path(), "status", status, "code", apiErr.Code, "error", err)
	}
	return c.Status(status).JSON(ErrorResponse{Error: ErrorDTO{
		Code: apiErr.Code, Message: apiErr.Message, Fields: apiErr.Fields, RequestID: requestID,
	}})
}

func normalize(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) && fiberErr.Code == fiber.StatusNotFound {
		return NewAPIError(apperror.ResourceNotFound, "The requested resource was not found.", nil)
	}
	return WrapAPIError(err, apperror.InternalError, "An unexpected error occurred.")
}

func StatusForCode(code ErrorCode) int {
	switch code {
	case apperror.ValidationError, apperror.UnsupportedFileType, apperror.FileTooLarge:
		return fiber.StatusBadRequest
	case apperror.InvalidCredentials, apperror.AuthenticationRequired, apperror.InvalidWebhookSignature:
		return fiber.StatusUnauthorized
	case apperror.Forbidden, apperror.ArtisanNotApproved:
		return fiber.StatusForbidden
	case apperror.ResourceNotFound:
		return fiber.StatusNotFound
	case apperror.Conflict, apperror.ProductNotEditable, apperror.ProductNotSellable,
		apperror.InvalidStateTransition, apperror.OutOfStock, apperror.ReservationExpired,
		apperror.IdempotencyConflict, apperror.PaymentAmountMismatch, apperror.DuplicateWebhook:
		return fiber.StatusConflict
	case apperror.RateLimited:
		return fiber.StatusTooManyRequests
	default:
		return fiber.StatusInternalServerError
	}
}
