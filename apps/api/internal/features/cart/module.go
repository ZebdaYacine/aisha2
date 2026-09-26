package cart

import (
	application "github.com/aisha-platform/aisha/apps/api/internal/features/cart/application"
	"github.com/aisha-platform/aisha/apps/api/internal/features/cart/data/repositories"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service = application.Service
type Item = application.Item
type Input = application.Input
type Authorizer = application.Authorizer

var ErrValidation = application.ErrValidation
var ErrUnavailable = application.ErrUnavailable

func NewService(pool *pgxpool.Pool, a *authorization.Service) *Service {
	return application.NewService(repositories.NewPostgresRepository(pool), a)
}
