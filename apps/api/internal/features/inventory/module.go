package inventory

import (
	application "github.com/aisha-platform/aisha/apps/api/internal/features/inventory/application"
	"github.com/aisha-platform/aisha/apps/api/internal/features/inventory/data/repositories"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service = application.Service
type Balance = application.Balance
type Movement = application.Movement
type AdjustmentInput = application.AdjustmentInput
type Authorizer = application.Authorizer

var (
	ErrValidation        = application.ErrValidation
	ErrNotFound          = application.ErrNotFound
	ErrDuplicate         = application.ErrDuplicate
	ErrInsufficientStock = application.ErrInsufficientStock
)

func NewService(pool *pgxpool.Pool, authorizer *authorization.Service) *Service {
	return application.NewService(repositories.NewPostgresRepository(pool), authorizer)
}

// Module is the composition boundary for this feature. Business behavior is
// intentionally added through the documented domain, application, data, and
// server layers as the feature is implemented.
type Module struct{}

func NewModule() Module { return Module{} }
