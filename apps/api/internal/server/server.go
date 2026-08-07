package httpapi

import (
	"github.com/aisha-platform/aisha/apps/api/internal/config"
	"github.com/aisha-platform/aisha/apps/api/internal/features/admin"
	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue"
	"github.com/aisha-platform/aisha/apps/api/internal/features/product"
	"github.com/aisha-platform/aisha/apps/api/internal/features/user"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/health"
	servermiddleware "github.com/aisha-platform/aisha/apps/api/internal/server/middleware"
	"github.com/aisha-platform/aisha/apps/api/internal/server/routes"
	"github.com/gofiber/fiber/v3"
)

func New(cfg config.Config, healthService *health.Service, authService *auth.Service, authorizationService *authorization.Service, rateLimiter RateLimiter, artisanService *artisan.Service, customerService *customer.Service, catalogueServices ...*catalogue.Service) *fiber.App {
	return newServer(cfg, healthService, authService, authorizationService, rateLimiter, artisanService, customerService, nil, nil, nil, catalogueServices...)
}

func NewWithProduct(cfg config.Config, healthService *health.Service, authService *auth.Service, authorizationService *authorization.Service, rateLimiter RateLimiter, artisanService *artisan.Service, customerService *customer.Service, productService *product.Service, artisanMediaService *artisan.MediaService, adminService *admin.Service, catalogueServices ...*catalogue.Service) *fiber.App {
	return newServer(cfg, healthService, authService, authorizationService, rateLimiter, artisanService, customerService, productService, artisanMediaService, adminService, catalogueServices...)
}

func newServer(cfg config.Config, healthService *health.Service, authService *auth.Service, authorizationService *authorization.Service, rateLimiter RateLimiter, artisanService *artisan.Service, customerService *customer.Service, productService *product.Service, artisanMediaService *artisan.MediaService, adminService *admin.Service, catalogueServices ...*catalogue.Service) *fiber.App {
	bodyLimit := cfg.UploadMaxBytes + 2*1024*1024
	if bodyLimit <= 2*1024*1024 {
		bodyLimit = 52 * 1024 * 1024
	}
	app := fiber.New(fiber.Config{AppName: cfg.AppName + " API", BodyLimit: int(bodyLimit), ErrorHandler: errorHandler})
	app.Use(servermiddleware.Recovery())
	app.Use(servermiddleware.RequestID())
	app.Use(servermiddleware.SecurityHeaders())
	app.Use(servermiddleware.CORS(cfg.AllowedOrigins))
	app.Use(servermiddleware.Logging())
	routes.RegisterHealth(app, healthService)

	api := APIGroup(app)
	var authHandler *AuthHandler
	if authService != nil {
		authHandler = NewAuthHandler(authService, NewRequestValidator())
	}

	public := routes.PublicRoutes{}
	if authHandler != nil {
		public.Register = []fiber.Handler{authRateLimit(rateLimiter, cfg.AuthRateLimitMax, cfg.AuthRateLimitWindow, "register"), authHandler.Register}
		public.Login = []fiber.Handler{authRateLimit(rateLimiter, cfg.AuthRateLimitMax, cfg.AuthRateLimitWindow, "login"), authHandler.Login}
		public.Refresh = []fiber.Handler{authHandler.Refresh}
		public.Logout = []fiber.Handler{authHandler.Logout}
		public.ForgotPassword = []fiber.Handler{authRateLimit(rateLimiter, cfg.AuthRateLimitMax, cfg.AuthRateLimitWindow, "forgot-password"), authHandler.ForgotPassword}
		public.ResetPassword = []fiber.Handler{authHandler.ResetPassword}
	}
	if len(catalogueServices) > 0 && catalogueServices[0] != nil {
		handler := NewCatalogueHandler(catalogueServices[0])
		public.Categories = []fiber.Handler{handler.Categories}
		public.Products = []fiber.Handler{handler.Products}
		public.Product = []fiber.Handler{handler.Product}
		public.Artisans = []fiber.Handler{handler.Artisans}
		public.Artisan = []fiber.Handler{handler.Artisan}
	}
	routes.RegisterPublic(api, public)

	if authHandler != nil && authorizationService != nil {
		casbin := NewCasbinMiddleware(authorizationService)
		var customerHandler *CustomerHandler
		if customerService != nil {
			customerHandler = NewCustomerHandler(customerService, NewRequestValidator())
		}
		var artisanHandler *ArtisanHandler
		if artisanService != nil {
			artisanHandler = NewArtisanHandler(artisanService, NewRequestValidator(), artisanMediaService)
		}
		var artisanMediaHandler *ArtisanMediaHandler
		if artisanMediaService != nil {
			artisanMediaHandler = NewArtisanMediaHandler(artisanMediaService, cfg.UploadMaxBytes)
		}
		var adminHandler *AdminHandler
		if adminService != nil {
			adminHandler = NewAdminHandler(adminService, NewRequestValidator())
		}
		var productHandler *ProductHandler
		if productService != nil {
			productHandler = NewProductHandler(productService, NewRequestValidator(), cfg.ProductMediaMaxBytes)
		}
		routes.RegisterAuthenticated(api, routes.AuthenticatedRoutes{
			Authenticate: authHandler.RequirePrincipal, Authorize: casbin.Require, Me: authHandler.Me,
			Profile: customerProfile(customerHandler), UpdateProfile: customerUpdateProfile(customerHandler),
			Addresses: customerAddresses(customerHandler), CreateAddress: customerCreateAddress(customerHandler),
			UpdateAddress: customerUpdateAddress(customerHandler), DeleteAddress: customerDeleteAddress(customerHandler),
			SubmitArtisan: artisanSubmit(artisanHandler), MineArtisan: artisanMine(artisanHandler),
			UpdateArtisanProfile: artisanUpdateProfile(artisanHandler),
			ListProducts:         productList(productHandler), CreateProduct: productCreate(productHandler), GetProduct: productGet(productHandler),
			UpdateProduct: productUpdate(productHandler), SubmitProduct: productSubmit(productHandler), UploadProductMedia: productUploadMedia(productHandler), DeleteProductMedia: productDeleteMedia(productHandler),
			ArtisanDocuments: artisanDocumentsMine(artisanMediaHandler), UploadArtisanDocument: artisanDocumentUpload(artisanMediaHandler), ArtisanMedia: artisanMediaList(artisanMediaHandler), UploadArtisanMedia: artisanMediaUpload(artisanMediaHandler),
		})
		routes.RegisterAdmin(api, routes.AdminRoutes{
			Authenticate: authHandler.RequirePrincipal, Authorize: casbin.Require,
			ListApplications: artisanList(artisanHandler), DecideApplication: artisanDecide(artisanHandler),
			ApplicationDocuments: artisanDocuments(artisanHandler),
			ListUsers:            adminUsers(adminHandler), UpdateUserRoles: adminRoles(adminHandler), AuditEvents: adminAudit(adminHandler),
		})
	}
	return app
}

func productList(h *ProductHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.List
}
func productCreate(h *ProductHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Create
}
func productGet(h *ProductHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Get
}
func productUpdate(h *ProductHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Update
}
func productSubmit(h *ProductHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Submit
}
func productUploadMedia(h *ProductHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UploadMedia
}
func productDeleteMedia(h *ProductHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.DeleteMedia
}
func artisanDocumentsMine(h *ArtisanMediaHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Documents
}
func artisanDocumentUpload(h *ArtisanMediaHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UploadDocument
}
func artisanMediaList(h *ArtisanMediaHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Media
}
func artisanMediaUpload(h *ArtisanMediaHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UploadMedia
}
func adminUsers(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Users
}
func adminRoles(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Roles
}
func adminAudit(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Audit
}

func customerProfile(h *CustomerHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Profile
}
func customerUpdateProfile(h *CustomerHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UpdateProfile
}
func customerAddresses(h *CustomerHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Addresses
}
func customerCreateAddress(h *CustomerHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.CreateAddress
}
func customerUpdateAddress(h *CustomerHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UpdateAddress
}
func customerDeleteAddress(h *CustomerHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.DeleteAddress
}
func artisanSubmit(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Submit
}
func artisanMine(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Mine
}
func artisanUpdateProfile(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UpdateProfile
}
func artisanList(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.List
}
func artisanDecide(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Decide
}
func artisanDocuments(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Documents
}
