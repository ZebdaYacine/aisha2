package customer

import (
	"context"
	"errors"
	"testing"

	"github.com/aisha-platform/aisha/backend/features/auth"
)

type repositoryStub struct {
	profileUser string
	addressUser string
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
