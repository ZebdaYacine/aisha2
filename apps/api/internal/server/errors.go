package httpapi

import (
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/apperror"
	servererrors "github.com/aisha-platform/aisha/apps/api/internal/server/error_handler"
	"github.com/gofiber/fiber/v3"
)

type ErrorCode = servererrors.ErrorCode
type FieldErrors = servererrors.FieldErrors
type APIError = servererrors.APIError
type ErrorResponse = servererrors.ErrorResponse
type ErrorDTO = servererrors.ErrorDTO

const (
	CodeValidationError         = apperror.ValidationError
	CodeInvalidCredentials      = apperror.InvalidCredentials
	CodeAuthenticationRequired  = apperror.AuthenticationRequired
	CodeForbidden               = apperror.Forbidden
	CodeResourceNotFound        = apperror.ResourceNotFound
	CodeConflict                = apperror.Conflict
	CodeArtisanNotApproved      = apperror.ArtisanNotApproved
	CodeProductNotEditable      = apperror.ProductNotEditable
	CodeProductNotSellable      = apperror.ProductNotSellable
	CodeInvalidStateTransition  = apperror.InvalidStateTransition
	CodeOutOfStock              = apperror.OutOfStock
	CodeReservationExpired      = apperror.ReservationExpired
	CodeIdempotencyConflict     = apperror.IdempotencyConflict
	CodePaymentAmountMismatch   = apperror.PaymentAmountMismatch
	CodeInvalidWebhookSignature = apperror.InvalidWebhookSignature
	CodeDuplicateWebhook        = apperror.DuplicateWebhook
	CodeUnsupportedFileType     = apperror.UnsupportedFileType
	CodeFileTooLarge            = apperror.FileTooLarge
	CodeRateLimited             = apperror.RateLimited
	CodeInternalError           = apperror.InternalError
)

func NewAPIError(code ErrorCode, message string, fields FieldErrors) *APIError {
	return servererrors.NewAPIError(code, message, fields)
}

func WrapAPIError(cause error, code ErrorCode, message string) *APIError {
	return servererrors.WrapAPIError(cause, code, message)
}

func errorHandler(c fiber.Ctx, err error) error {
	return servererrors.Handler(c, err)
}

func statusForCode(code ErrorCode) int {
	return servererrors.StatusForCode(code)
}
