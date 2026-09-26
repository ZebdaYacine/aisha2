package artisan

import "github.com/aisha-platform/aisha/apps/api/internal/features/artisan/domain"

type Translation = domain.Translation
type ApplicationInput = domain.ApplicationInput
type Application = domain.Application
type Document = domain.Document
type Media = domain.Media
type DocumentUploadInput = domain.DocumentUploadInput
type MediaUploadInput = domain.MediaUploadInput
type Workshop = domain.Workshop
type WorkshopInput = domain.WorkshopInput
type Verification = domain.Verification
type VerificationDecision = domain.VerificationDecision
type Repository = domain.Repository
type Authorizer = domain.Authorizer
type MediaRepository = domain.MediaRepository

var (
	ErrNotFound          = domain.ErrNotFound
	ErrValidation        = domain.ErrValidation
	ErrInvalidTransition = domain.ErrInvalidTransition
)
