package admin

import "github.com/aisha-platform/aisha/apps/api/internal/features/admin/domain"

type User = domain.User
type UserInput = domain.UserInput
type Category = domain.Category
type CategoryInput = domain.CategoryInput
type Order = domain.Order
type AuditEvent = domain.AuditEvent
type AuditFilter = domain.AuditFilter
type Repository = domain.Repository
type Authorizer = domain.Authorizer

var (
	ErrNotFound   = domain.ErrNotFound
	ErrValidation = domain.ErrValidation
)
