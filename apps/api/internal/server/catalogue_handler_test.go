package httpapi

import (
	"context"
	"encoding/json"
	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue"
	"github.com/gofiber/fiber/v3"
	"net/http/httptest"
	"testing"
)

type catalogueRepoStub struct {
	request        catalogue.PageRequest
	productRequest catalogue.PageRequest
	category       string
	query          string
	workshop       string
	product        catalogue.Product
	workshopItem   catalogue.Workshop
}

func (r *catalogueRepoStub) Categories(_ context.Context, p catalogue.PageRequest) (catalogue.Page[catalogue.Category], error) {
	r.request = p
	return catalogue.Page[catalogue.Category]{Items: []catalogue.Category{{ID: "id", Slug: "pottery", Name: "Pottery"}}, Page: p.Page, PageSize: p.PageSize, Total: 1}, nil
}
func (r *catalogueRepoStub) Products(_ context.Context, p catalogue.PageRequest, category, query, workshop string) (catalogue.Page[catalogue.Product], error) {
	r.productRequest, r.category, r.query, r.workshop = p, category, query, workshop
	return catalogue.Page[catalogue.Product]{}, nil
}
func (r *catalogueRepoStub) Product(_ context.Context, id, _ string) (catalogue.Product, error) {
	if id == "mapped" {
		return r.product, nil
	}
	return catalogue.Product{}, catalogue.ErrNotFound
}
func (r *catalogueRepoStub) Artisans(context.Context, catalogue.PageRequest) (catalogue.Page[catalogue.Artisan], error) {
	return catalogue.Page[catalogue.Artisan]{}, nil
}
func (r *catalogueRepoStub) Artisan(context.Context, string, string) (catalogue.Artisan, error) {
	return catalogue.Artisan{}, catalogue.ErrNotFound
}
func (r *catalogueRepoStub) Workshops(context.Context, catalogue.PageRequest) (catalogue.Page[catalogue.Workshop], error) {
	return catalogue.Page[catalogue.Workshop]{Items: []catalogue.Workshop{r.workshopItem}, Page: 1, PageSize: 20, Total: 1}, nil
}
func (r *catalogueRepoStub) Workshop(_ context.Context, id, _ string) (catalogue.Workshop, error) {
	if id == "workshop-1" {
		return r.workshopItem, nil
	}
	return catalogue.Workshop{}, catalogue.ErrNotFound
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

func TestCatalogueHandlerPassesProductSearchFilters(t *testing.T) {
	repo := &catalogueRepoStub{}
	handler := NewCatalogueHandler(catalogue.NewService(repo))
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Get("/products", handler.Products)
	response, err := app.Test(httptest.NewRequest("GET", "/products?locale=ar&page=2&pageSize=12&category=pottery&q=clay&workshop=workshop-1", nil))
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("status=%d", response.StatusCode)
	}
	if repo.productRequest.Locale != "ar" || repo.productRequest.Page != 2 || repo.productRequest.PageSize != 12 || repo.category != "pottery" || repo.query != "clay" || repo.workshop != "workshop-1" {
		t.Fatalf("request=%#v category=%q query=%q", repo.productRequest, repo.category, repo.query)
	}
}

func TestCatalogueHandlerMapsWorkshopAndAvailability(t *testing.T) {
	repo := &catalogueRepoStub{
		product:      catalogue.Product{ID: "mapped", ArtisanID: "artisan-1", ArtisanName: "Maker", WorkshopID: "workshop-1", WorkshopName: "Atelier", AvailableQuantity: 1},
		workshopItem: catalogue.Workshop{ID: "workshop-1", Name: "Atelier", ArtisanID: "artisan-1", ArtisanName: "Maker", ProductCount: 2},
	}
	handler := NewCatalogueHandler(catalogue.NewService(repo))
	app := fiber.New(fiber.Config{ErrorHandler: errorHandler})
	app.Get("/products/:id", handler.Product)
	app.Get("/workshops/:id", handler.Workshop)
	response, err := app.Test(httptest.NewRequest("GET", "/products/mapped", nil))
	if err != nil {
		t.Fatal(err)
	}
	var product ProductDTO
	if err := json.NewDecoder(response.Body).Decode(&product); err != nil {
		t.Fatal(err)
	}
	if product.WorkshopID != "workshop-1" || product.Workshop != "Atelier" || product.Availability != "LOW_STOCK" || product.AvailableQuantity != 1 {
		t.Fatalf("product=%#v", product)
	}
	response, err = app.Test(httptest.NewRequest("GET", "/workshops/workshop-1", nil))
	if err != nil {
		t.Fatal(err)
	}
	var workshop WorkshopDTO
	if err := json.NewDecoder(response.Body).Decode(&workshop); err != nil {
		t.Fatal(err)
	}
	if workshop.ID != "workshop-1" || workshop.ArtisanID != "artisan-1" {
		t.Fatalf("workshop=%#v", workshop)
	}
}
