package httpapi

import (
	"github.com/aisha-platform/aisha/apps/api/internal/config"
	"github.com/aisha-platform/aisha/apps/api/internal/features/admin"
	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/cart"
	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue"
	"github.com/aisha-platform/aisha/apps/api/internal/features/inventory"
	"github.com/aisha-platform/aisha/apps/api/internal/features/moderation"
	"github.com/aisha-platform/aisha/apps/api/internal/features/notification"
	"github.com/aisha-platform/aisha/apps/api/internal/features/order"
	"github.com/aisha-platform/aisha/apps/api/internal/features/product"
	"github.com/aisha-platform/aisha/apps/api/internal/features/user"
	"github.com/aisha-platform/aisha/apps/api/internal/features/warehouse"
	"github.com/aisha-platform/aisha/apps/api/internal/features/wishlist"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/authorization"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/health"
	servermiddleware "github.com/aisha-platform/aisha/apps/api/internal/server/middleware"
	"github.com/aisha-platform/aisha/apps/api/internal/server/routes"
	"github.com/gofiber/fiber/v3"
)

func New(cfg config.Config, healthService *health.Service, authService *auth.Service, authorizationService *authorization.Service, rateLimiter RateLimiter, artisanService *artisan.Service, customerService *customer.Service, catalogueServices ...*catalogue.Service) *fiber.App {
	return newServer(cfg, healthService, authService, authorizationService, rateLimiter, artisanService, customerService, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, catalogueServices...)
}

func NewWithProduct(cfg config.Config, healthService *health.Service, authService *auth.Service, authorizationService *authorization.Service, rateLimiter RateLimiter, artisanService *artisan.Service, customerService *customer.Service, productService *product.Service, artisanMediaService *artisan.MediaService, adminService *admin.Service, catalogueServices ...*catalogue.Service) *fiber.App {
	return newServer(cfg, healthService, authService, authorizationService, rateLimiter, artisanService, customerService, productService, artisanMediaService, adminService, nil, nil, nil, nil, nil, nil, nil, nil, catalogueServices...)
}

func NewWithWorkflows(cfg config.Config, healthService *health.Service, authService *auth.Service, authorizationService *authorization.Service, rateLimiter RateLimiter, artisanService *artisan.Service, customerService *customer.Service, productService *product.Service, artisanMediaService *artisan.MediaService, adminService *admin.Service, moderationService *moderation.Service, orderService *order.Service, warehouseService *warehouse.Service, inventoryService *inventory.Service, cartService *cart.Service, wishlistService *wishlist.Service, notificationService *notification.Service, notificationHub *notification.Hub, catalogueServices ...*catalogue.Service) *fiber.App {
	return newServer(cfg, healthService, authService, authorizationService, rateLimiter, artisanService, customerService, productService, artisanMediaService, adminService, moderationService, orderService, warehouseService, inventoryService, cartService, wishlistService, notificationService, notificationHub, catalogueServices...)
}

func newServer(cfg config.Config, healthService *health.Service, authService *auth.Service, authorizationService *authorization.Service, rateLimiter RateLimiter, artisanService *artisan.Service, customerService *customer.Service, productService *product.Service, artisanMediaService *artisan.MediaService, adminService *admin.Service, moderationService *moderation.Service, orderService *order.Service, warehouseService *warehouse.Service, inventoryService *inventory.Service, cartService *cart.Service, wishlistService *wishlist.Service, notificationService *notification.Service, notificationHub *notification.Hub, catalogueServices ...*catalogue.Service) *fiber.App {
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
		if customerService != nil {
			authHandler = NewAuthHandler(authService, NewRequestValidator(), customerService)
		} else {
			authHandler = NewAuthHandler(authService, NewRequestValidator())
		}
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
		public.Workshops = []fiber.Handler{handler.Workshops}
		public.Workshop = []fiber.Handler{handler.Workshop}
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
		var moderationHandler *ModerationHandler
		if moderationService != nil {
			moderationHandler = NewModerationHandler(moderationService, NewRequestValidator())
		}
		var orderHandler *OrderHandler
		if orderService != nil {
			orderHandler = NewOrderHandler(orderService, NewRequestValidator())
		}
		var warehouseHandler *WarehouseHandler
		if warehouseService != nil {
			warehouseHandler = NewWarehouseHandler(warehouseService, NewRequestValidator(), cfg.UploadMaxBytes)
		}
		var inventoryHandler *InventoryHandler
		if inventoryService != nil {
			inventoryHandler = NewInventoryHandler(inventoryService, NewRequestValidator())
		}
		var cartHandler *CartHandler
		if cartService != nil {
			cartHandler = NewCartHandler(cartService, NewRequestValidator())
		}
		var wishlistHandler *WishlistHandler
		if wishlistService != nil {
			wishlistHandler = NewWishlistHandler(wishlistService)
		}
		var notificationHandler *NotificationHandler
		if notificationService != nil && notificationHub != nil {
			notificationHandler = NewNotificationHandler(notificationService, notificationHub, cfg.AllowedOrigins)
		}
		routes.RegisterAuthenticated(api, routes.AuthenticatedRoutes{
			Authenticate: authHandler.RequirePrincipal, Authorize: casbin.Require, Me: authHandler.Me,
			Profile: customerProfile(customerHandler), UpdateProfile: customerUpdateProfile(customerHandler), ChangePassword: authChangePassword(authHandler),
			Addresses: customerAddresses(customerHandler), CreateAddress: customerCreateAddress(customerHandler),
			UpdateAddress: customerUpdateAddress(customerHandler), DeleteAddress: customerDeleteAddress(customerHandler),
			SubmitArtisan: artisanSubmit(artisanHandler), SaveArtisanDraft: artisanSaveDraft(artisanHandler), MineArtisan: artisanMine(artisanHandler), FinalizeArtisanSubmission: artisanFinalizeSubmission(artisanHandler),
			UpdateArtisanProfile: artisanUpdateProfile(artisanHandler),
			ListWorkshops:        artisanWorkshops(artisanHandler), CreateWorkshop: artisanCreateWorkshop(artisanHandler), UpdateWorkshop: artisanUpdateWorkshop(artisanHandler), WorkshopStatus: artisanWorkshopStatus(artisanHandler), DeleteWorkshop: artisanDeleteWorkshop(artisanHandler), ActivateMembership: artisanActivateMembership(artisanHandler), Verification: artisanVerification(artisanHandler),
			ListProducts: productList(productHandler), CreateProduct: productCreate(productHandler), GetProduct: productGet(productHandler),
			UpdateProduct: productUpdate(productHandler), SubmitProduct: productSubmit(productHandler), ArchiveProduct: productArchive(productHandler), UploadProductMedia: productUploadMedia(productHandler), DeleteProductMedia: productDeleteMedia(productHandler),
			ArtisanDocuments: artisanDocumentsMine(artisanMediaHandler), UploadArtisanDocument: artisanDocumentUpload(artisanMediaHandler), ArtisanMedia: artisanMediaList(artisanMediaHandler), UploadArtisanMedia: artisanMediaUpload(artisanMediaHandler), ReplaceArtisanMedia: artisanMediaReplace(artisanMediaHandler), DeleteArtisanMedia: artisanMediaDelete(artisanMediaHandler),
			Checkout: orderCheckout(orderHandler), ListOrders: orderList(orderHandler), GetOrder: orderGet(orderHandler), CancelOrder: orderCancel(orderHandler), SellerOrders: orderSeller(orderHandler),
			WarehouseReceptions: warehouseList(warehouseHandler), WarehouseReceptionCreate: warehouseCreate(warehouseHandler), WarehouseInspect: warehouseInspect(warehouseHandler), WarehouseReceptionEvidence: warehouseUploadEvidence(warehouseHandler), WarehouseReceptionEvidenceList: warehouseListEvidence(warehouseHandler), Inventory: inventoryList(inventoryHandler), InventoryAdjust: inventoryAdjust(inventoryHandler),
			WarehouseProducts: warehouseProducts(warehouseHandler),
			Cart:              cartList(cartHandler), CartAdd: cartAdd(cartHandler), CartSet: cartSet(cartHandler), CartRemove: cartRemove(cartHandler), CartMerge: cartMerge(cartHandler), Wishlist: wishlistList(wishlistHandler), WishlistAdd: wishlistAdd(wishlistHandler), WishlistRemove: wishlistRemove(wishlistHandler),
			Notifications: notificationList(notificationHandler), NotificationRead: notificationRead(notificationHandler), NotificationsReadAll: notificationReadAll(notificationHandler), NotificationSocketTicket: notificationSocketTicket(notificationHandler), NotificationSocket: notificationSocket(notificationHandler),
		})
		routes.RegisterAdmin(api, routes.AdminRoutes{
			Authenticate: authHandler.RequirePrincipal, Authorize: casbin.Require,
			ListApplications: artisanList(artisanHandler), DecideApplication: artisanDecide(artisanHandler),
			ApplicationDocuments: artisanDocuments(artisanHandler), ApplicationMedia: artisanApplicationMedia(artisanHandler),
			ListUsers: adminUsers(adminHandler), CreateUser: adminCreateUser(adminHandler), UpdateUser: adminUpdateUser(adminHandler), UpdateUserRoles: adminRoles(adminHandler), UserStatus: adminUserStatus(adminHandler), AuditEvents: adminAudit(adminHandler),
			ProductSubmissions: moderationQueue(moderationHandler), ProductDecision: moderationDecision(moderationHandler), RecordReturn: orderReturn(orderHandler), MembershipStatus: artisanMembershipStatus(artisanHandler), WorkshopStatus: artisanAdminWorkshopStatus(artisanHandler), ListVerifications: artisanVerifications(artisanHandler), DecideVerification: artisanDecideVerification(artisanHandler),
			UserMedia: adminUserMedia(adminHandler), ProductMedia: adminProductMedia(adminHandler), DeleteUserMedia: adminDeleteUserMedia(adminHandler), DeleteProductMedia: adminDeleteProductMedia(adminHandler),
			Categories: adminCategories(adminHandler), CreateCategory: adminCreateCategory(adminHandler), UpdateCategory: adminUpdateCategory(adminHandler), DeleteCategory: adminDeleteCategory(adminHandler), Orders: adminOrders(adminHandler),
		})
	}
	return app
}

func moderationQueue(h *ModerationHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Queue
}
func moderationDecision(h *ModerationHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Decide
}
func orderCheckout(h *OrderHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Checkout
}
func orderList(h *OrderHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.List
}
func orderGet(h *OrderHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Get
}
func orderCancel(h *OrderHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Cancel
}
func orderSeller(h *OrderHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Seller
}
func orderReturn(h *OrderHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Return
}

func warehouseList(h *WarehouseHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.List
}
func warehouseProducts(h *WarehouseHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Products
}
func warehouseCreate(h *WarehouseHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Create
}
func warehouseInspect(h *WarehouseHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Inspect
}
func warehouseUploadEvidence(h *WarehouseHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UploadEvidence
}
func warehouseListEvidence(h *WarehouseHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Evidence
}

func adminUserMedia(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UserMedia
}
func adminProductMedia(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.ProductMedia
}
func adminDeleteUserMedia(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.DeleteUserMedia
}
func adminDeleteProductMedia(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.DeleteProductMedia
}

func inventoryList(h *InventoryHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.List
}
func inventoryAdjust(h *InventoryHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Adjust
}

func cartList(h *CartHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.List
}
func cartAdd(h *CartHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Add
}
func cartSet(h *CartHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Set
}
func cartRemove(h *CartHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Remove
}
func cartMerge(h *CartHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Merge
}
func wishlistList(h *WishlistHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.List
}
func wishlistAdd(h *WishlistHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Add
}
func wishlistRemove(h *WishlistHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Remove
}

func artisanWorkshops(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Workshops
}
func artisanCreateWorkshop(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.CreateWorkshop
}
func artisanUpdateWorkshop(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UpdateWorkshop
}
func artisanWorkshopStatus(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.WorkshopStatus
}
func artisanDeleteWorkshop(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.DeleteWorkshop
}
func artisanActivateMembership(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.ActivateMembership
}
func artisanVerification(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Verification
}
func artisanAdminWorkshopStatus(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.AdminWorkshopStatus
}
func artisanMembershipStatus(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.MembershipStatus
}
func artisanVerifications(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Verifications
}
func artisanDecideVerification(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.DecideVerification
}

func productWorkshops(h *ProductHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.ListWorkshops
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
func productArchive(h *ProductHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Archive
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
func artisanMediaReplace(h *ArtisanMediaHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.ReplaceMedia
}
func artisanMediaDelete(h *ArtisanMediaHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.DeleteMedia
}
func adminUsers(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Users
}

func adminCreateUser(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.CreateUser
}

func adminUpdateUser(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UpdateUser
}

func adminOrders(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Orders
}
func adminRoles(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Roles
}
func adminUserStatus(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UserStatus
}
func adminAudit(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Audit
}

func adminCategories(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Categories
}
func adminCreateCategory(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.CreateCategory
}
func adminUpdateCategory(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.UpdateCategory
}
func adminDeleteCategory(h *AdminHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.DeleteCategory
}

func notificationList(h *NotificationHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.List
}

func notificationRead(h *NotificationHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.MarkRead
}

func notificationReadAll(h *NotificationHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.MarkAllRead
}

func notificationSocketTicket(h *NotificationHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.SocketTicket
}

func notificationSocket(h *NotificationHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.WebSocket
}

func authChangePassword(h *AuthHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.ChangePassword
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
func artisanSaveDraft(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.SaveDraft
}
func artisanMine(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Mine
}
func artisanFinalizeSubmission(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.FinalizeSubmission
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

func artisanApplicationMedia(h *ArtisanHandler) fiber.Handler {
	if h == nil {
		return nil
	}
	return h.Media
}
