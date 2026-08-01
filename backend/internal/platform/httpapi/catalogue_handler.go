package httpapi

import (
	"errors"
	"strconv"

	"github.com/aisha-platform/aisha/backend/internal/catalogue"
	"github.com/gofiber/fiber/v3"
)

type CatalogueHandler struct{ service *catalogue.Service }

func NewCatalogueHandler(service *catalogue.Service) *CatalogueHandler {
	return &CatalogueHandler{service: service}
}

type PageDTO[T any] struct {
	Items    []T `json:"items"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}
type CategoryDTO struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}
type ProductDTO struct {
	ID               string   `json:"id"`
	ArtisanID        string   `json:"artisanId"`
	ArtisanName      string   `json:"artisanName"`
	CategoryID       string   `json:"categoryId"`
	CategorySlug     string   `json:"categorySlug"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Story            string   `json:"story"`
	Materials        string   `json:"materials"`
	ProductionMethod string   `json:"productionMethod"`
	Region           string   `json:"region"`
	Currency         string   `json:"currency"`
	Status           string   `json:"status"`
	PriceMinor       int64    `json:"priceMinor"`
	Media            []string `json:"media"`
}
type ArtisanDTO struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Workshop     string   `json:"workshop"`
	Wilaya       string   `json:"wilaya"`
	Location     string   `json:"location"`
	Biography    string   `json:"biography"`
	Media        []string `json:"media"`
	ProductCount int      `json:"productCount"`
}

func pageRequest(c fiber.Ctx) catalogue.PageRequest {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("pageSize", "20"))
	return catalogue.PageRequest{Locale: c.Query("locale", "en"), Page: page, PageSize: size}
}
func (h *CatalogueHandler) Categories(c fiber.Ctx) error {
	result, err := h.service.Categories(c.Context(), pageRequest(c))
	if err != nil {
		return catalogueError(err)
	}
	items := make([]CategoryDTO, len(result.Items))
	for i, item := range result.Items {
		items[i] = CategoryDTO{ID: item.ID, Slug: item.Slug, Name: item.Name}
	}
	return c.JSON(PageDTO[CategoryDTO]{Items: items, Page: result.Page, PageSize: result.PageSize, Total: result.Total})
}
func (h *CatalogueHandler) Products(c fiber.Ctx) error {
	result, err := h.service.Products(c.Context(), pageRequest(c), c.Query("category"))
	if err != nil {
		return catalogueError(err)
	}
	items := make([]ProductDTO, len(result.Items))
	for i, item := range result.Items {
		items[i] = productDTO(item)
	}
	return c.JSON(PageDTO[ProductDTO]{Items: items, Page: result.Page, PageSize: result.PageSize, Total: result.Total})
}
func (h *CatalogueHandler) Product(c fiber.Ctx) error {
	item, err := h.service.Product(c.Context(), c.Params("id"), c.Query("locale", "en"))
	if err != nil {
		return catalogueError(err)
	}
	return c.JSON(productDTO(item))
}
func (h *CatalogueHandler) Artisans(c fiber.Ctx) error {
	result, err := h.service.Artisans(c.Context(), pageRequest(c))
	if err != nil {
		return catalogueError(err)
	}
	items := make([]ArtisanDTO, len(result.Items))
	for i, item := range result.Items {
		items[i] = artisanDTO(item)
	}
	return c.JSON(PageDTO[ArtisanDTO]{Items: items, Page: result.Page, PageSize: result.PageSize, Total: result.Total})
}
func (h *CatalogueHandler) Artisan(c fiber.Ctx) error {
	item, err := h.service.Artisan(c.Context(), c.Params("id"), c.Query("locale", "en"))
	if err != nil {
		return catalogueError(err)
	}
	return c.JSON(artisanDTO(item))
}
func productDTO(item catalogue.Product) ProductDTO {
	return ProductDTO{ID: item.ID, ArtisanID: item.ArtisanID, ArtisanName: item.ArtisanName, CategoryID: item.CategoryID, CategorySlug: item.CategorySlug, Name: item.Name, Description: item.Description, Story: item.Story, Materials: item.Materials, ProductionMethod: item.ProductionMethod, Region: item.Region, Currency: item.Currency, Status: item.Status, PriceMinor: item.PriceMinor, Media: item.Media}
}
func artisanDTO(item catalogue.Artisan) ArtisanDTO {
	return ArtisanDTO{ID: item.ID, Name: item.Name, Workshop: item.Workshop, Wilaya: item.Wilaya, Location: item.Location, Biography: item.Biography, Media: item.Media, ProductCount: item.ProductCount}
}
func catalogueError(err error) error {
	if errors.Is(err, catalogue.ErrNotFound) {
		return NewAPIError(CodeResourceNotFound, "The requested resource was not found.", nil)
	}
	return WrapAPIError(err, CodeInternalError, "An unexpected error occurred.")
}
