package routes

import "github.com/gofiber/fiber/v3"

type AdminRoutes struct {
	Authenticate         fiber.Handler
	Authorize            func(resource, action string) fiber.Handler
	ListApplications     fiber.Handler
	DecideApplication    fiber.Handler
	ApplicationDocuments fiber.Handler
	ApplicationMedia     fiber.Handler
	ListUsers            fiber.Handler
	UpdateUserRoles      fiber.Handler
	UserStatus           fiber.Handler
	AuditEvents          fiber.Handler
	ProductSubmissions   fiber.Handler
	ProductDecision      fiber.Handler
	RecordReturn         fiber.Handler
	MembershipStatus     fiber.Handler
	WorkshopStatus       fiber.Handler
	ListVerifications    fiber.Handler
	DecideVerification   fiber.Handler
	UserMedia            fiber.Handler
	ProductMedia         fiber.Handler
	DeleteUserMedia      fiber.Handler
	DeleteProductMedia   fiber.Handler
}

func RegisterAdmin(api fiber.Router, r AdminRoutes) {
	if r.Authenticate == nil {
		return
	}
	add(api, "GET", "/admin/artisan-applications", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/artisan-applications", "read"), r.ListApplications)
	add(api, "POST", "/admin/artisan-applications/:id/:decision", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/artisan-applications/*", "write"), r.DecideApplication)
	add(api, "GET", "/admin/artisan-applications/:id/documents", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/artisan-applications/*/documents", "read"), r.ApplicationDocuments)
	add(api, "GET", "/admin/artisan-applications/:id/media", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/artisan-applications/*/media", "read"), r.ApplicationMedia)
	add(api, "GET", "/admin/users", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/users", "read"), r.ListUsers)
	add(api, "PATCH", "/admin/users/:id/roles", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/users", "write"), r.UpdateUserRoles)
	add(api, "POST", "/admin/users/:id/status", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/users", "write"), r.UserStatus)
	add(api, "GET", "/admin/audit-events", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/audit-events", "read"), r.AuditEvents)
	add(api, "GET", "/admin/product-submissions", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/product-submissions", "read"), r.ProductSubmissions)
	add(api, "POST", "/admin/product-submissions/:id/decision", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/product-submissions", "write"), r.ProductDecision)
	add(api, "POST", "/admin/orders/:id/returns", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/orders/*/returns", "write"), r.RecordReturn)
	add(api, "POST", "/admin/artisan-memberships/:id/status", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/artisan-memberships", "write"), r.MembershipStatus)
	add(api, "POST", "/admin/workshops/:id/status", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/workshops", "write"), r.WorkshopStatus)
	add(api, "GET", "/admin/artisan-verifications", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/artisan-verifications", "read"), r.ListVerifications)
	add(api, "POST", "/admin/artisan-verifications/:id/decision", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/artisan-verifications", "write"), r.DecideVerification)
	add(api, "GET", "/admin/media/users", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/media", "read"), r.UserMedia)
	add(api, "GET", "/admin/media/products", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/media", "read"), r.ProductMedia)
	add(api, "DELETE", "/admin/media/users/:id", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/media", "write"), r.DeleteUserMedia)
	add(api, "DELETE", "/admin/media/products/:id", r.Authenticate, authorizeAdmin(r, "/api/v1/admin/media", "write"), r.DeleteProductMedia)
}

func authorizeAdmin(r AdminRoutes, resource, action string) fiber.Handler {
	if r.Authorize == nil {
		return nil
	}
	return r.Authorize(resource, action)
}
