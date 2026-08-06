package apperror

type Code string

const (
	ValidationError         Code = "VALIDATION_ERROR"
	InvalidCredentials      Code = "INVALID_CREDENTIALS"
	AuthenticationRequired  Code = "AUTHENTICATION_REQUIRED"
	Forbidden               Code = "FORBIDDEN"
	ResourceNotFound        Code = "RESOURCE_NOT_FOUND"
	Conflict                Code = "CONFLICT"
	ArtisanNotApproved      Code = "ARTISAN_NOT_APPROVED"
	ProductNotEditable      Code = "PRODUCT_NOT_EDITABLE"
	ProductNotSellable      Code = "PRODUCT_NOT_SELLABLE"
	InvalidStateTransition  Code = "INVALID_STATE_TRANSITION"
	OutOfStock              Code = "OUT_OF_STOCK"
	ReservationExpired      Code = "RESERVATION_EXPIRED"
	IdempotencyConflict     Code = "IDEMPOTENCY_CONFLICT"
	PaymentAmountMismatch   Code = "PAYMENT_AMOUNT_MISMATCH"
	InvalidWebhookSignature Code = "INVALID_WEBHOOK_SIGNATURE"
	DuplicateWebhook        Code = "DUPLICATE_WEBHOOK"
	UnsupportedFileType     Code = "UNSUPPORTED_FILE_TYPE"
	FileTooLarge            Code = "FILE_TOO_LARGE"
	RateLimited             Code = "RATE_LIMITED"
	InternalError           Code = "INTERNAL_ERROR"
)
