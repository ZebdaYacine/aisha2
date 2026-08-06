package httpapi

import (
	"github.com/aisha-platform/aisha/apps/api/internal/config"
	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue"
	"github.com/aisha-platform/aisha/apps/api/internal/features/user"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/health"
	servermiddleware "github.com/aisha-platform/aisha/apps/api/internal/server/middleware"
	"github.com/aisha-platform/aisha/apps/api/internal/server/routes"
	"github.com/gofiber/fiber/v3"
)

func New(cfg config.Config, healthService *health.Service, authService *auth.Service, authorizationService *authorization.Service, rateLimiter RateLimiter, artisanService *artisan.Service, customerService *customer.Service, catalogueServices ...*catalogue.Service) *fiber.App {
	app := fiber.New(fiber.Config{AppName: cfg.AppName + " API", BodyLimit: 2 * 1024 * 1024, ErrorHandler: errorHandler})
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
			artisanHandler = NewArtisanHandler(artisanService, NewRequestValidator())
		}
		routes.RegisterAuthenticated(api, routes.AuthenticatedRoutes{
			Authenticate: authHandler.RequirePrincipal, Authorize: casbin.Require, Me: authHandler.Me,
			Profile: customerProfile(customerHandler), UpdateProfile: customerUpdateProfile(customerHandler),
			Addresses: customerAddresses(customerHandler), CreateAddress: customerCreateAddress(customerHandler),
			UpdateAddress: customerUpdateAddress(customerHandler), DeleteAddress: customerDeleteAddress(customerHandler),
			SubmitArtisan: artisanSubmit(artisanHandler), MineArtisan: artisanMine(artisanHandler),
			UpdateArtisanProfile: artisanUpdateProfile(artisanHandler),
		})
		routes.RegisterAdmin(api, routes.AdminRoutes{
			Authenticate: authHandler.RequirePrincipal, Authorize: casbin.Require,
			ListApplications: artisanList(artisanHandler), DecideApplication: artisanDecide(artisanHandler),
			ApplicationDocuments: artisanDocuments(artisanHandler),
		})
	}
	return app
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
