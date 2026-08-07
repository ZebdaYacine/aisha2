package routes

import "github.com/gofiber/fiber/v3"

type AuthenticatedRoutes struct {
	Authenticate          fiber.Handler
	Authorize             func(resource, action string) fiber.Handler
	Me                    fiber.Handler
	Profile               fiber.Handler
	UpdateProfile         fiber.Handler
	Addresses             fiber.Handler
	CreateAddress         fiber.Handler
	UpdateAddress         fiber.Handler
	DeleteAddress         fiber.Handler
	SubmitArtisan         fiber.Handler
	MineArtisan           fiber.Handler
	UpdateArtisanProfile  fiber.Handler
	ListProducts          fiber.Handler
	CreateProduct         fiber.Handler
	GetProduct            fiber.Handler
	UpdateProduct         fiber.Handler
	SubmitProduct         fiber.Handler
	UploadProductMedia    fiber.Handler
	DeleteProductMedia    fiber.Handler
	ArtisanDocuments      fiber.Handler
	UploadArtisanDocument fiber.Handler
	ArtisanMedia          fiber.Handler
	UploadArtisanMedia    fiber.Handler
}

func RegisterAuthenticated(api fiber.Router, r AuthenticatedRoutes) {
	if r.Authenticate == nil {
		return
	}
	add(api, "GET", "/me", r.Authenticate, authorize(r, "/api/v1/me", "read"), r.Me)
	add(api, "GET", "/me/profile", r.Authenticate, r.Profile)
	add(api, "PATCH", "/me/profile", r.Authenticate, r.UpdateProfile)
	add(api, "GET", "/addresses", r.Authenticate, r.Addresses)
	add(api, "POST", "/addresses", r.Authenticate, r.CreateAddress)
	add(api, "PATCH", "/addresses/:id", r.Authenticate, r.UpdateAddress)
	add(api, "DELETE", "/addresses/:id", r.Authenticate, r.DeleteAddress)
	add(api, "POST", "/artisan-applications", r.Authenticate, r.SubmitArtisan)
	add(api, "GET", "/artisan-applications/me", r.Authenticate, r.MineArtisan)
	add(api, "PATCH", "/artisan/profile", r.Authenticate, r.UpdateArtisanProfile)
	add(api, "GET", "/artisan/products", r.Authenticate, authorize(r, "/api/v1/artisan/products", "read"), r.ListProducts)
	add(api, "POST", "/artisan/products", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.CreateProduct)
	add(api, "GET", "/artisan/products/:id", r.Authenticate, authorize(r, "/api/v1/artisan/products", "read"), r.GetProduct)
	add(api, "PATCH", "/artisan/products/:id", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.UpdateProduct)
	add(api, "POST", "/artisan/products/:id/submit", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.SubmitProduct)
	add(api, "POST", "/artisan/products/:id/media", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.UploadProductMedia)
	add(api, "DELETE", "/artisan/products/:id/media/:mediaId", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.DeleteProductMedia)
	add(api, "GET", "/artisan-applications/me/documents", r.Authenticate, authorize(r, "/api/v1/artisan-applications/me/documents", "read"), r.ArtisanDocuments)
	add(api, "POST", "/artisan-applications/me/documents", r.Authenticate, authorize(r, "/api/v1/artisan-applications/me/documents", "write"), r.UploadArtisanDocument)
	add(api, "GET", "/artisan/profile/media", r.Authenticate, authorize(r, "/api/v1/artisan/profile/media", "read"), r.ArtisanMedia)
	add(api, "POST", "/artisan/profile/media", r.Authenticate, authorize(r, "/api/v1/artisan/profile/media", "write"), r.UploadArtisanMedia)
}

func authorize(r AuthenticatedRoutes, resource, action string) fiber.Handler {
	if r.Authorize == nil {
		return nil
	}
	return r.Authorize(resource, action)
}
