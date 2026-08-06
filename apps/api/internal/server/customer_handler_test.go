package httpapi

import (
	"context"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/user"
	"github.com/gofiber/fiber/v3"
	"net/http/httptest"
	"strings"
	"testing"
)

type customerRepoStub struct{ user string }

func (r *customerRepoStub) Profile(context.Context, string) (customer.Profile, error) {
	return customer.Profile{}, nil
}
func (r *customerRepoStub) UpdateProfile(context.Context, string, string, string) (customer.Profile, error) {
	return customer.Profile{}, nil
}
func (r *customerRepoStub) Addresses(context.Context, string) ([]customer.Address, error) {
	return nil, nil
}
func (r *customerRepoStub) CreateAddress(_ context.Context, user string, i customer.AddressInput) (customer.Address, error) {
	r.user = user
	return customer.Address{ID: "address", FullName: i.FullName}, nil
}
func (r *customerRepoStub) UpdateAddress(context.Context, string, string, customer.AddressInput) (customer.Address, error) {
	return customer.Address{}, nil
}
func (r *customerRepoStub) DeleteAddress(context.Context, string, string) error { return nil }

type allowCustomer struct{}

func (allowCustomer) Authorize(context.Context, auth.Principal, string, string) error { return nil }
func TestCustomerHandlerUsesAuthenticatedPrincipal(t *testing.T) {
	repo := &customerRepoStub{}
	handler := NewCustomerHandler(customer.NewService(repo, allowCustomer{}), NewRequestValidator())
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Post("/addresses", func(c fiber.Ctx) error { c.Locals(principalLocal, auth.Principal{UserID: "owner"}); return c.Next() }, handler.CreateAddress)
	body := `{"fullName":"Amina","line1":"12 Rue","city":"Alger","postalCode":"16000","country":"Algeria","isDefault":true}`
	request := httptest.NewRequest("POST", "/addresses", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 201 || repo.user != "owner" {
		t.Fatalf("status=%d user=%q", response.StatusCode, repo.user)
	}
}
func TestCustomerHandlerReturnsFieldValidationErrors(t *testing.T) {
	handler := NewCustomerHandler(customer.NewService(&customerRepoStub{}, allowCustomer{}), NewRequestValidator())
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Post("/addresses", func(c fiber.Ctx) error { c.Locals(principalLocal, auth.Principal{UserID: "owner"}); return c.Next() }, handler.CreateAddress)
	request := httptest.NewRequest("POST", "/addresses", strings.NewReader(`{}`))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 400 {
		t.Fatalf("status=%d", response.StatusCode)
	}
}
