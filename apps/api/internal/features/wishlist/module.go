package wishlist

import (
	application "github.com/aisha-platform/aisha/apps/api/internal/features/wishlist/application"
	"github.com/aisha-platform/aisha/apps/api/internal/features/wishlist/data/repositories"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service = application.Service
type Item = application.Item

var ErrValidation = application.ErrValidation

func NewService(pool *pgxpool.Pool, a *authorization.Service) *Service {
	return application.NewService(repositories.NewPostgresRepository(pool), a)
}
