package customer

import "github.com/aisha-platform/aisha/apps/api/internal/features/user/domain"

type Profile = domain.Profile
type Address = domain.Address
type AddressInput = domain.AddressInput
type AccountState = domain.AccountState
type AccountSummary = domain.AccountSummary
type Repository = domain.Repository
type AccountRepository = domain.AccountRepository
type Authorizer = domain.Authorizer

var (
	ErrNotFound   = domain.ErrNotFound
	ErrValidation = domain.ErrValidation
)
