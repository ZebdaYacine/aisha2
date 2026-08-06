package catalogue

import application "github.com/aisha-platform/aisha/apps/api/internal/features/catalogue/application"

type Service = application.Service

func NewService(repository Repository) *Service { return application.NewService(repository) }
