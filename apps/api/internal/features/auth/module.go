package auth

import (
	"errors"

	application "github.com/aisha-platform/aisha/apps/api/internal/features/auth/application"
)

type Service = application.Service

func NewService(repository Repository, notifier ResetNotifier, secret string) *Service {
	return application.NewService(repository, notifier, secret)
}

func Is(err, target error) bool { return errors.Is(err, target) }
