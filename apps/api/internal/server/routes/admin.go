package routes

import "github.com/gofiber/fiber/v3"

type AdminRoutes struct {
	Authenticate         fiber.Handler
	Authorize            func(resource, action string) fiber.Handler
	ListApplications     fiber.Handler
	DecideApplication    fiber.Handler
	ApplicationDocuments fiber.Handler
	ListUsers            fiber.Handler
	UpdateUserRoles      fiber.Handler
	AuditEvents          fiber.Handler
}

func RegisterAdmin(api fiber.Router, r AdminRoutes) {
	if r.Authenticate == nil {
		return
	}
	add(api, "GET", "/admin/artisan-applications", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/artisan-applications", "read"), r.ListApplications)
	add(api, "POST", "/admin/artisan-applications/:id/:decision", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/artisan-applications/*", "write"), r.DecideApplication)
	add(api, "GET", "/admin/artisan-applications/:id/documents", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/artisan-applications/*/documents", "read"), r.ApplicationDocuments)
	add(api, "GET", "/admin/users", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/users", "read"), r.ListUsers)
	add(api, "PATCH", "/admin/users/:id/roles", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/users", "write"), r.UpdateUserRoles)
	add(api, "GET", "/admin/audit-events", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/audit-events", "read"), r.AuditEvents)
}

func authorizeAdmin(r AdminRoutes, resource, action string) fiber.Handler {
	if r.Authorize == nil {
		return nil
	}
	return r.Authorize(resource, action)
}
