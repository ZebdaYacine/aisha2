package customer

import application "github.com/aisha-platform/aisha/apps/api/internal/features/user/application"

type Service = application.Service

func NewService(repository Repository, authorizer Authorizer) *Service {
	return application.NewService(repository, authorizer)
}
