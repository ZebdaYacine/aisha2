package order

import (
	application "github.com/aisha-platform/aisha/apps/api/internal/features/order/application"
	"github.com/aisha-platform/aisha/apps/api/internal/features/order/data/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service = application.Service
type CartItem = application.CartItem
type Order = application.Order
type SellerItem = application.SellerItem
type Return = application.Return
type Authorizer = application.Authorizer

var (
	ErrValidation            = application.ErrValidation
	ErrNotFound              = application.ErrNotFound
	ErrOutOfStock            = application.ErrOutOfStock
	ErrPriceChanged          = application.ErrPriceChanged
	ErrPaymentAmountMismatch = application.ErrPaymentAmountMismatch
	ErrInvalidTransition     = application.ErrInvalidTransition
	ErrIdempotencyConflict   = application.ErrIdempotencyConflict
)

func NewService(pool *pgxpool.Pool, a Authorizer) *Service {
	return application.NewService(repositories.NewPostgresRepository(pool), a)
}
