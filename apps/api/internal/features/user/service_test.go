package customer

import (
	"context"
	"errors"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

type repositoryStub struct {
	profileUser string
	addressUser string
}

type accountRepositoryStub struct {
	repositoryStub
	state AccountState
}

func (r *accountRepositoryStub) AccountState(context.Context, string) (AccountState, error) {
	return r.state, nil
}

func (r *repositoryStub) Profile(_ context.Context, user string) (Profile, error) {
	r.profileUser = user
	return Profile{ID: user}, nil
}
func (r *repositoryStub) UpdateProfile(_ context.Context, user, name, phone string) (Profile, error) {
	r.profileUser = user
	return Profile{ID: user, DisplayName: name, Phone: phone}, nil
}
func (r *repositoryStub) Addresses(_ context.Context, user string) ([]Address, error) {
	r.addressUser = user
	return []Address{}, nil
}
func (r *repositoryStub) CreateAddress(_ context.Context, user string, input AddressInput) (Address, error) {
	r.addressUser = user
	return Address{FullName: input.FullName}, nil
}
func (r *repositoryStub) UpdateAddress(_ context.Context, user, id string, input AddressInput) (Address, error) {
	r.addressUser = user
	return Address{ID: id}, nil
}
func (r *repositoryStub) DeleteAddress(_ context.Context, user, id string) error {
	r.addressUser = user
	return nil
}

type authorizerStub struct{ err error }

func (a authorizerStub) Authorize(context.Context, auth.Principal, string, string) error {
	return a.err
}

func TestProfileUsesAuthenticatedUserID(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository, authorizerStub{})
	_, err := service.Profile(context.Background(), auth.Principal{UserID: "user-1", Roles: []string{"customer"}})
	if err != nil {
		t.Fatalf("Profile() error = %v", err)
	}
	if repository.profileUser != "user-1" {
		t.Fatalf("Profile() user = %q", repository.profileUser)
	}
}

func TestArtisanUsesTheSameCustomerProfileAndAddressServices(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository, authorizerStub{})
	principal := auth.Principal{UserID: "artisan-user", Roles: []string{"customer", "artisan"}}
	if _, err := service.UpdateProfile(context.Background(), principal, "Artisan Buyer", "0555000000"); err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
	if _, err := service.Addresses(context.Background(), principal); err != nil {
		t.Fatalf("Addresses() error = %v", err)
	}
	if repository.profileUser != principal.UserID || repository.addressUser != principal.UserID {
		t.Fatalf("customer resources did not use the same user ID: profile=%q addresses=%q", repository.profileUser, repository.addressUser)
	}
}

func TestAddressRejectsInvalidInputBeforeRepository(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository, authorizerStub{})
	_, err := service.CreateAddress(context.Background(), auth.Principal{UserID: "user-1"}, AddressInput{})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("CreateAddress() error = %v", err)
	}
	if repository.addressUser != "" {
		t.Fatal("repository called for invalid address")
	}
}

func TestAuthorizationFailureStopsCustomerRead(t *testing.T) {
	repository := &repositoryStub{}
	denied := errors.New("denied")
	service := NewService(repository, authorizerStub{err: denied})
	_, err := service.Addresses(context.Background(), auth.Principal{UserID: "user-1"})
	if !errors.Is(err, denied) {
		t.Fatalf("Addresses() error = %v", err)
	}
	if repository.addressUser != "" {
		t.Fatal("repository called after authorization failure")
	}
}

func TestAccountSummaryKeepsCustomerCapabilitiesForNormalUser(t *testing.T) {
	repository := &accountRepositoryStub{state: AccountState{UserStatus: "ACTIVE"}}
	service := NewService(repository, authorizerStub{})
	summary, err := service.AccountSummary(context.Background(), auth.Principal{UserID: "user-1", Roles: []string{"customer"}})
	if err != nil {
		t.Fatal(err)
	}
	if summary.UserID != "user-1" || !summary.CustomerEnabled || summary.ArtisanEnabled || summary.ArtisanStatus != "NOT_STARTED" {
		t.Fatalf("unexpected customer summary: %#v", summary)
	}
	if !containsCapability(summary.Capabilities, "customer.profile.write") {
		t.Fatalf("customer capability missing: %#v", summary.Capabilities)
	}
}

func TestAccountSummaryAddsArtisanCapabilitiesToSameUser(t *testing.T) {
	repository := &accountRepositoryStub{state: AccountState{UserStatus: "ACTIVE", ArtisanStatus: "APPROVED"}}
	service := NewService(repository, authorizerStub{})
	summary, err := service.AccountSummary(context.Background(), auth.Principal{UserID: "same-user", Roles: []string{"customer", "artisan"}})
	if err != nil {
		t.Fatal(err)
	}
	if summary.UserID != "same-user" || !summary.CustomerEnabled || !summary.ArtisanEnabled || summary.ArtisanStatus != "ACTIVE" {
		t.Fatalf("unexpected active artisan summary: %#v", summary)
	}
	if !containsCapability(summary.Capabilities, "customer.addresses.write") || !containsCapability(summary.Capabilities, "artisan.products.write") {
		t.Fatalf("combined capabilities missing: %#v", summary.Capabilities)
	}
}

func TestSuspendedArtisanKeepsCustomerCapabilities(t *testing.T) {
	repository := &accountRepositoryStub{state: AccountState{UserStatus: "ACTIVE", ArtisanStatus: "SUSPENDED"}}
	service := NewService(repository, authorizerStub{})
	summary, err := service.AccountSummary(context.Background(), auth.Principal{UserID: "artisan-user", Roles: []string{"customer", "artisan"}})
	if err != nil {
		t.Fatal(err)
	}
	if !summary.CustomerEnabled || summary.ArtisanEnabled || summary.ArtisanStatus != "SUSPENDED" {
		t.Fatalf("unexpected suspended summary: %#v", summary)
	}
	if !containsCapability(summary.Capabilities, "customer.profile.write") || containsCapability(summary.Capabilities, "artisan.products.write") {
		t.Fatalf("suspended capabilities are incorrect: %#v", summary.Capabilities)
	}
}

func TestSuspendedUserHasNoProtectedCapabilities(t *testing.T) {
	repository := &accountRepositoryStub{state: AccountState{UserStatus: "SUSPENDED", ArtisanStatus: "APPROVED"}}
	service := NewService(repository, authorizerStub{})
	summary, err := service.AccountSummary(context.Background(), auth.Principal{UserID: "suspended-user", Roles: []string{"customer", "artisan"}})
	if err != nil {
		t.Fatal(err)
	}
	if summary.CustomerEnabled || summary.ArtisanEnabled || len(summary.Capabilities) != 0 {
		t.Fatalf("suspended user still has capabilities: %#v", summary)
	}
}

func containsCapability(capabilities []string, expected string) bool {
	for _, capability := range capabilities {
		if capability == expected {
			return true
		}
	}
	return false
}

func TestAdministratorCannotBecomeArtisan(t *testing.T) {
	repository := &accountRepositoryStub{state: AccountState{UserStatus: "ACTIVE", ArtisanStatus: "APPROVED"}}
	service := NewService(repository, authorizerStub{})
	summary, err := service.AccountSummary(context.Background(), auth.Principal{UserID: "admin-1", Roles: []string{"administrator"}})
	if err != nil {
		t.Fatal(err)
	}
	if !summary.CustomerEnabled || summary.ArtisanEnabled || summary.ArtisanStatus != "NOT_STARTED" {
		t.Fatalf("administrator received artisan state: %#v", summary)
	}
	if containsCapability(summary.Capabilities, "artisan.account.read") {
		t.Fatalf("administrator received artisan capability: %#v", summary.Capabilities)
	}
}

func TestBackOfficeCapabilitiesFollowAssignedRole(t *testing.T) {
	for _, test := range []struct {
		role       string
		capability string
	}{
		{role: "moderator", capability: "admin.product_moderation.write"},
		{role: "warehouse_agent", capability: "warehouse.write"},
		{role: "administrator", capability: "admin.audit.read"},
	} {
		repository := &accountRepositoryStub{state: AccountState{UserStatus: "ACTIVE"}}
		service := NewService(repository, authorizerStub{})
		summary, err := service.AccountSummary(context.Background(), auth.Principal{UserID: test.role + "-1", Roles: []string{test.role}})
		if err != nil {
			t.Fatal(err)
		}
		if !containsCapability(summary.Capabilities, test.capability) {
			t.Fatalf("role %s missing %s: %#v", test.role, test.capability, summary.Capabilities)
		}
	}
}
