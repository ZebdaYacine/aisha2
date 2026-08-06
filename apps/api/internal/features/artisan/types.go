package artisan

import "github.com/aisha-platform/aisha/apps/api/internal/features/artisan/domain"

type Translation = domain.Translation
type ApplicationInput = domain.ApplicationInput
type Application = domain.Application
type Document = domain.Document
type Repository = domain.Repository
type Authorizer = domain.Authorizer

var (
	ErrNotFound          = domain.ErrNotFound
	ErrValidation        = domain.ErrValidation
	ErrInvalidTransition = domain.ErrInvalidTransition
)
