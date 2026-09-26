package warehouse

import (
	application "github.com/aisha-platform/aisha/apps/api/internal/features/warehouse/application"
	"github.com/aisha-platform/aisha/apps/api/internal/features/warehouse/data/repositories"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service = application.Service
type Reception = application.Reception
type Inspection = application.Inspection
type Evidence = application.Evidence
type ReceptionInput = application.ReceptionInput
type InspectionInput = application.InspectionInput

var (
	ErrValidation        = application.ErrValidation
	ErrNotFound          = application.ErrNotFound
	ErrDuplicate         = application.ErrDuplicate
	ErrInvalidTransition = application.ErrInvalidTransition
	ErrQuantityMismatch  = application.ErrQuantityMismatch
	ErrEvidenceRequired  = application.ErrEvidenceRequired
	ErrAlreadyInspected  = application.ErrAlreadyInspected
)

func NewService(pool *pgxpool.Pool, authorizer application.Authorizer, store storage.ObjectStore, bucket string, maxUpload int64) *Service {
	return application.NewService(repositories.NewPostgresRepository(pool), authorizer, store, bucket, maxUpload)
}
