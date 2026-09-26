package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/user/domain"
)

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
}

func NewService(repository domain.Repository, authorizer domain.Authorizer) *Service {
	return &Service{repository: repository, authorizer: authorizer}
}

const (
	capCustomerAccountRead    = "customer.account.read"
	capCustomerProfileRead    = "customer.profile.read"
	capCustomerProfileWrite   = "customer.profile.write"
	capCustomerAddressesRead  = "customer.addresses.read"
	capCustomerAddressesWrite = "customer.addresses.write"
	capCustomerPurchasesRead  = "customer.purchases.read"
	capCustomerTrackingRead   = "customer.tracking.read"
	capArtisanAccountRead     = "artisan.account.read"
	capArtisanProfileRead     = "artisan.profile.read"
	capArtisanProfileWrite    = "artisan.profile.write"
	capArtisanProductsRead    = "artisan.products.read"
	capArtisanProductsWrite   = "artisan.products.write"
	capArtisanOrdersRead      = "artisan.orders.read"
	capAdminDashboardRead     = "admin.dashboard.read"
	capAdminUsersRead         = "admin.users.read"
	capAdminUsersWrite        = "admin.users.write"
	capAdminApplicationsRead  = "admin.artisan_applications.read"
	capAdminApplicationsWrite = "admin.artisan_applications.write"
	capAdminMediaRead         = "admin.media.read"
	capAdminMediaWrite        = "admin.media.write"
	capAdminModerationRead    = "admin.product_moderation.read"
	capAdminModerationWrite   = "admin.product_moderation.write"
	capAdminAuditRead         = "admin.audit.read"
	capWarehouseRead          = "warehouse.read"
	capWarehouseWrite         = "warehouse.write"
	capInventoryRead          = "inventory.read"
	capInventoryWrite         = "inventory.write"
)

// AccountSummary returns trusted capability state for the authenticated
// account. Capabilities are derived on the backend; client-provided roles or
// modes are never consulted here.
func (s *Service) AccountSummary(ctx context.Context, p auth.Principal) (domain.AccountSummary, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/me", "read"); err != nil {
		return domain.AccountSummary{}, err
	}

	state := domain.AccountState{UserStatus: "ACTIVE"}
	if repository, ok := s.repository.(domain.AccountRepository); ok {
		var err error
		state, err = repository.AccountState(ctx, p.UserID)
		if err != nil {
			return domain.AccountSummary{}, fmt.Errorf("read account state: %w", err)
		}
	}

	isAdministrator := hasRole(p.Roles, "administrator")
	artisanStatus := artisanSummaryStatus(state.ArtisanStatus)
	if isAdministrator {
		artisanStatus = "NOT_STARTED"
	}
	summary := domain.AccountSummary{
		UserID:          p.UserID,
		UserStatus:      state.UserStatus,
		CustomerEnabled: state.UserStatus == "ACTIVE",
		ArtisanStatus:   artisanStatus,
	}
	if summary.CustomerEnabled {
		summary.Capabilities = append(summary.Capabilities,
			capCustomerAccountRead,
			capCustomerProfileRead,
			capCustomerProfileWrite,
			capCustomerAddressesRead,
			capCustomerAddressesWrite,
			capCustomerPurchasesRead,
			capCustomerTrackingRead,
		)
		if hasRole(p.Roles, "administrator") || hasRole(p.Roles, "moderator") || hasRole(p.Roles, "warehouse_agent") {
			summary.Capabilities = append(summary.Capabilities, capAdminDashboardRead)
		}
		if hasRole(p.Roles, "administrator") || hasRole(p.Roles, "moderator") || hasRole(p.Roles, "warehouse_agent") {
			summary.Capabilities = append(summary.Capabilities, capAdminUsersRead)
		}
		if hasRole(p.Roles, "administrator") || hasRole(p.Roles, "warehouse_agent") {
			summary.Capabilities = append(summary.Capabilities, capAdminApplicationsRead)
		}
		if hasRole(p.Roles, "administrator") || hasRole(p.Roles, "moderator") {
			summary.Capabilities = append(summary.Capabilities, capAdminMediaRead, capAdminModerationRead, capAdminModerationWrite)
		}
		if hasRole(p.Roles, "administrator") {
			summary.Capabilities = append(summary.Capabilities, capAdminUsersWrite, capAdminApplicationsWrite, capAdminMediaWrite, capAdminAuditRead)
		}
		if hasRole(p.Roles, "warehouse_agent") || hasRole(p.Roles, "administrator") {
			summary.Capabilities = append(summary.Capabilities, capWarehouseRead, capWarehouseWrite, capInventoryRead, capInventoryWrite)
		}
	}
	if !isAdministrator && summary.ArtisanStatus == "ACTIVE" && summary.CustomerEnabled {
		summary.ArtisanEnabled = true
		summary.Capabilities = append(summary.Capabilities,
			capArtisanAccountRead,
			capArtisanProfileRead,
			capArtisanProfileWrite,
			capArtisanProductsRead,
			capArtisanProductsWrite,
			capArtisanOrdersRead,
		)
	}
	return summary, nil
}

func artisanSummaryStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "APPROVED", "ACTIVE":
		return "ACTIVE"
	case "SUSPENDED":
		return "SUSPENDED"
	case "DRAFT", "SUBMITTED", "UNDER_REVIEW", "CHANGES_REQUESTED", "PENDING":
		return "PENDING"
	default:
		return "NOT_STARTED"
	}
}

func hasRole(roles []string, expected string) bool {
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
}
func (s *Service) Profile(ctx context.Context, p auth.Principal) (domain.Profile, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/me", "read"); err != nil {
		return domain.Profile{}, err
	}
	return s.repository.Profile(ctx, p.UserID)
}
func (s *Service) UpdateProfile(ctx context.Context, p auth.Principal, name, phone string) (domain.Profile, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/me", "write"); err != nil {
		return domain.Profile{}, err
	}
	name, phone = strings.TrimSpace(name), strings.TrimSpace(phone)
	if len(name) < 2 {
		return domain.Profile{}, domain.ErrValidation
	}
	return s.repository.UpdateProfile(ctx, p.UserID, name, phone)
}
func (s *Service) Addresses(ctx context.Context, p auth.Principal) ([]domain.Address, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "read"); err != nil {
		return nil, err
	}
	return s.repository.Addresses(ctx, p.UserID)
}
func (s *Service) CreateAddress(ctx context.Context, p auth.Principal, input domain.AddressInput) (domain.Address, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "write"); err != nil {
		return domain.Address{}, err
	}
	if !validAddress(input) {
		return domain.Address{}, domain.ErrValidation
	}
	return s.repository.CreateAddress(ctx, p.UserID, input)
}
func (s *Service) UpdateAddress(ctx context.Context, p auth.Principal, id string, input domain.AddressInput) (domain.Address, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "write"); err != nil {
		return domain.Address{}, err
	}
	if !validAddress(input) {
		return domain.Address{}, domain.ErrValidation
	}
	return s.repository.UpdateAddress(ctx, p.UserID, id, input)
}
func (s *Service) DeleteAddress(ctx context.Context, p auth.Principal, id string) error {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/addresses", "write"); err != nil {
		return err
	}
	return s.repository.DeleteAddress(ctx, p.UserID, id)
}
func validAddress(i domain.AddressInput) bool {
	return strings.TrimSpace(i.FullName) != "" && strings.TrimSpace(i.Line1) != "" && strings.TrimSpace(i.City) != "" && strings.TrimSpace(i.PostalCode) != "" && strings.TrimSpace(i.Country) != ""
}
