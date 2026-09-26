package routes

import "github.com/gofiber/fiber/v3"

type PublicRoutes struct {
	Register, Login, Refresh, Logout, ForgotPassword, ResetPassword       []fiber.Handler
	Categories, Products, Product, Artisans, Artisan, Workshops, Workshop []fiber.Handler
}

func RegisterPublic(api fiber.Router, r PublicRoutes) {
	add(api, "POST", "/auth/register", r.Register...)
	add(api, "POST", "/auth/login", r.Login...)
	add(api, "POST", "/auth/refresh", r.Refresh...)
	add(api, "POST", "/auth/logout", r.Logout...)
	add(api, "POST", "/auth/forgot-password", r.ForgotPassword...)
	add(api, "POST", "/auth/reset-password", r.ResetPassword...)
	add(api, "GET", "/categories", r.Categories...)
	add(api, "GET", "/products", r.Products...)
	add(api, "GET", "/products/:id", r.Product...)
	add(api, "GET", "/artisans", r.Artisans...)
	add(api, "GET", "/artisans/:id", r.Artisan...)
	add(api, "GET", "/workshops", r.Workshops...)
	add(api, "GET", "/workshops/:id", r.Workshop...)
}

func add(api fiber.Router, method, path string, handlers ...fiber.Handler) {
	items := make([]any, 0, len(handlers))
	for _, handler := range handlers {
		if handler != nil {
			items = append(items, handler)
		}
	}
	if len(items) == 0 {
		return
	}
	api.Add([]string{method}, path, items[0], items[1:]...)
}
