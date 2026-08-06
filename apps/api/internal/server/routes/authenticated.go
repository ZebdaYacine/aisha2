package routes

import "github.com/gofiber/fiber/v3"

type AuthenticatedRoutes struct {
	Authenticate         fiber.Handler
	Authorize            func(resource, action string) fiber.Handler
	Me                   fiber.Handler
	Profile              fiber.Handler
	UpdateProfile        fiber.Handler
	Addresses            fiber.Handler
	CreateAddress        fiber.Handler
	UpdateAddress        fiber.Handler
	DeleteAddress        fiber.Handler
	SubmitArtisan        fiber.Handler
	MineArtisan          fiber.Handler
	UpdateArtisanProfile fiber.Handler
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
}

func authorize(r AuthenticatedRoutes, resource, action string) fiber.Handler {
	if r.Authorize == nil {
		return nil
	}
	return r.Authorize(resource, action)
}
