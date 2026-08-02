package httpapi

import (
	"context"
	"encoding/json"
	"github.com/aisha-platform/aisha/backend/features/catalogue"
	"github.com/gofiber/fiber/v3"
	"net/http/httptest"
	"testing"
)

type catalogueRepoStub struct{ request catalogue.PageRequest }

func (r *catalogueRepoStub) Categories(_ context.Context, p catalogue.PageRequest) (catalogue.Page[catalogue.Category], error) {
	r.request = p
	return catalogue.Page[catalogue.Category]{Items: []catalogue.Category{{ID: "id", Slug: "pottery", Name: "Pottery"}}, Page: p.Page, PageSize: p.PageSize, Total: 1}, nil
}
func (r *catalogueRepoStub) Products(context.Context, catalogue.PageRequest, string) (catalogue.Page[catalogue.Product], error) {
	return catalogue.Page[catalogue.Product]{}, nil
}
func (r *catalogueRepoStub) Product(context.Context, string, string) (catalogue.Product, error) {
	return catalogue.Product{}, catalogue.ErrNotFound
}
func (r *catalogueRepoStub) Artisans(context.Context, catalogue.PageRequest) (catalogue.Page[catalogue.Artisan], error) {
	return catalogue.Page[catalogue.Artisan]{}, nil
}
func (r *catalogueRepoStub) Artisan(context.Context, string, string) (catalogue.Artisan, error) {
	return catalogue.Artisan{}, catalogue.ErrNotFound
}
func TestCatalogueHandlerReturnsBoundedLocalizedPage(t *testing.T) {
	repo := &catalogueRepoStub{}
	handler := NewCatalogueHandler(catalogue.NewService(repo))
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Get("/categories", handler.Categories)
	response, err := app.Test(httptest.NewRequest("GET", "/categories?locale=fr&pageSize=999", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("status=%d", response.StatusCode)
	}
	var body PageDTO[CategoryDTO]
	if err = json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if repo.request.Locale != "fr" || repo.request.PageSize != 100 || body.Items[0].Slug != "pottery" {
		t.Fatalf("request=%#v body=%#v", repo.request, body)
	}
}
func TestCatalogueHandlerMapsNotFound(t *testing.T) {
	handler := NewCatalogueHandler(catalogue.NewService(&catalogueRepoStub{}))
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Get("/products/:id", handler.Product)
	response, err := app.Test(httptest.NewRequest("GET", "/products/missing", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 404 {
		t.Fatalf("status=%d", response.StatusCode)
	}
}
