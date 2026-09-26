package moderation

import (
	application "github.com/aisha-platform/aisha/apps/api/internal/features/moderation/application"
	"github.com/aisha-platform/aisha/apps/api/internal/features/moderation/data/repositories"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service = application.Service
type QueueItem = application.QueueItem
type DecisionInput = application.DecisionInput
type Media = application.Media

var (
	ErrValidation         = application.ErrValidation
	ErrNotFound           = application.ErrNotFound
	ErrInvalidTransition  = application.ErrInvalidTransition
	ErrActivationNotReady = application.ErrActivationNotReady
)

func NewService(pool *pgxpool.Pool, authorizer application.Authorizer) *Service {
	return application.NewService(repositories.NewPostgresRepository(pool), authorizer)
}

func NewServiceWithMedia(pool *pgxpool.Pool, authorizer application.Authorizer, store storage.ObjectStore, bucket string) *Service {
	return application.NewServiceWithMedia(repositories.NewPostgresRepository(pool), authorizer, store, bucket)
}
