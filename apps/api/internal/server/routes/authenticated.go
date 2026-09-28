package routes

import "github.com/gofiber/fiber/v3"

type AuthenticatedRoutes struct {
	Authenticate                   fiber.Handler
	Authorize                      func(resource, action string) fiber.Handler
	Me                             fiber.Handler
	Profile                        fiber.Handler
	UpdateProfile                  fiber.Handler
	ChangePassword                 fiber.Handler
	Addresses                      fiber.Handler
	CreateAddress                  fiber.Handler
	UpdateAddress                  fiber.Handler
	DeleteAddress                  fiber.Handler
	SubmitArtisan                  fiber.Handler
	SaveArtisanDraft               fiber.Handler
	MineArtisan                    fiber.Handler
	FinalizeArtisanSubmission      fiber.Handler
	UpdateArtisanProfile           fiber.Handler
	ListWorkshops                  fiber.Handler
	ListProducts                   fiber.Handler
	CreateProduct                  fiber.Handler
	GetProduct                     fiber.Handler
	UpdateProduct                  fiber.Handler
	SubmitProduct                  fiber.Handler
	ArchiveProduct                 fiber.Handler
	UploadProductMedia             fiber.Handler
	DeleteProductMedia             fiber.Handler
	ArtisanDocuments               fiber.Handler
	UploadArtisanDocument          fiber.Handler
	ArtisanMedia                   fiber.Handler
	UploadArtisanMedia             fiber.Handler
	ReplaceArtisanMedia            fiber.Handler
	DeleteArtisanMedia             fiber.Handler
	Checkout                       fiber.Handler
	ListOrders                     fiber.Handler
	GetOrder                       fiber.Handler
	CancelOrder                    fiber.Handler
	SellerOrders                   fiber.Handler
	ActivateMembership             fiber.Handler
	CreateWorkshop                 fiber.Handler
	UpdateWorkshop                 fiber.Handler
	WorkshopStatus                 fiber.Handler
	DeleteWorkshop                 fiber.Handler
	Verification                   fiber.Handler
	WarehouseReceptions            fiber.Handler
	WarehouseProducts              fiber.Handler
	WarehouseReceptionCreate       fiber.Handler
	WarehouseReceptionEvidence     fiber.Handler
	WarehouseReceptionEvidenceList fiber.Handler
	WarehouseInspect               fiber.Handler
	Inventory                      fiber.Handler
	InventoryAdjust                fiber.Handler
	Cart                           fiber.Handler
	CartAdd                        fiber.Handler
	CartSet                        fiber.Handler
	CartRemove                     fiber.Handler
	CartMerge                      fiber.Handler
	Wishlist                       fiber.Handler
	WishlistAdd                    fiber.Handler
	WishlistRemove                 fiber.Handler
	Notifications                  fiber.Handler
	NotificationRead               fiber.Handler
	NotificationsReadAll           fiber.Handler
	NotificationSocketTicket       fiber.Handler
	NotificationSocket             fiber.Handler
}

func RegisterAuthenticated(api fiber.Router, r AuthenticatedRoutes) {
	if r.Authenticate == nil {
		return
	}
	add(api, "GET", "/me", r.Authenticate, authorize(r, "/api/v1/me", "read"), r.Me)
	add(api, "GET", "/me/profile", r.Authenticate, r.Profile)
	add(api, "PATCH", "/me/profile", r.Authenticate, r.UpdateProfile)
	add(api, "PATCH", "/me/password", r.Authenticate, r.ChangePassword)
	add(api, "GET", "/addresses", r.Authenticate, r.Addresses)
	add(api, "POST", "/addresses", r.Authenticate, r.CreateAddress)
	add(api, "PATCH", "/addresses/:id", r.Authenticate, r.UpdateAddress)
	add(api, "DELETE", "/addresses/:id", r.Authenticate, r.DeleteAddress)
	add(api, "POST", "/artisan-applications", r.Authenticate, r.SubmitArtisan)
	add(api, "POST", "/artisan-applications/draft", r.Authenticate, authorize(r, "/api/v1/artisan-applications/draft", "write"), r.SaveArtisanDraft)
	add(api, "GET", "/artisan-applications/me", r.Authenticate, r.MineArtisan)
	add(api, "POST", "/artisan-applications/me/submit", r.Authenticate, authorize(r, "/api/v1/artisan-applications/me/submit", "write"), r.FinalizeArtisanSubmission)
	add(api, "PATCH", "/artisan/profile", r.Authenticate, r.UpdateArtisanProfile)
	add(api, "GET", "/artisan/workshops", r.Authenticate, authorize(r, "/api/v1/artisan/workshops", "read"), r.ListWorkshops)
	add(api, "POST", "/artisan/workshops", r.Authenticate, authorize(r, "/api/v1/artisan/workshops", "write"), r.CreateWorkshop)
	add(api, "PATCH", "/artisan/workshops/:id", r.Authenticate, authorize(r, "/api/v1/artisan/workshops", "write"), r.UpdateWorkshop)
	add(api, "POST", "/artisan/workshops/:id/status", r.Authenticate, authorize(r, "/api/v1/artisan/workshops", "write"), r.WorkshopStatus)
	add(api, "DELETE", "/artisan/workshops/:id", r.Authenticate, authorize(r, "/api/v1/artisan/workshops", "write"), r.DeleteWorkshop)
	add(api, "POST", "/artisan/membership/activate", r.Authenticate, authorize(r, "/api/v1/artisan/membership", "write"), r.ActivateMembership)
	add(api, "GET", "/artisan/verification", r.Authenticate, authorize(r, "/api/v1/artisan/verification", "read"), r.Verification)
	add(api, "GET", "/artisan/products", r.Authenticate, authorize(r, "/api/v1/artisan/products", "read"), r.ListProducts)
	add(api, "POST", "/artisan/products", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.CreateProduct)
	add(api, "GET", "/artisan/products/:id", r.Authenticate, authorize(r, "/api/v1/artisan/products", "read"), r.GetProduct)
	add(api, "PATCH", "/artisan/products/:id", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.UpdateProduct)
	add(api, "POST", "/artisan/products/:id/submit", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.SubmitProduct)
	add(api, "POST", "/artisan/products/:id/archive", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.ArchiveProduct)
	add(api, "POST", "/artisan/products/:id/media", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.UploadProductMedia)
	add(api, "DELETE", "/artisan/products/:id/media/:mediaId", r.Authenticate, authorize(r, "/api/v1/artisan/products", "write"), r.DeleteProductMedia)
	add(api, "GET", "/artisan-applications/me/documents", r.Authenticate, authorize(r, "/api/v1/artisan-applications/me/documents", "read"), r.ArtisanDocuments)
	add(api, "POST", "/artisan-applications/me/documents", r.Authenticate, authorize(r, "/api/v1/artisan-applications/me/documents", "write"), r.UploadArtisanDocument)
	add(api, "GET", "/artisan/profile/media", r.Authenticate, authorize(r, "/api/v1/artisan/profile/media", "read"), r.ArtisanMedia)
	add(api, "POST", "/artisan/profile/media", r.Authenticate, authorize(r, "/api/v1/artisan/profile/media", "write"), r.UploadArtisanMedia)
	add(api, "PATCH", "/artisan/profile/media/:id", r.Authenticate, authorize(r, "/api/v1/artisan/profile/media", "write"), r.ReplaceArtisanMedia)
	add(api, "DELETE", "/artisan/profile/media/:id", r.Authenticate, authorize(r, "/api/v1/artisan/profile/media", "write"), r.DeleteArtisanMedia)
	add(api, "POST", "/checkout", r.Authenticate, authorize(r, "/api/v1/checkout", "write"), r.Checkout)
	add(api, "GET", "/orders/seller", r.Authenticate, authorize(r, "/api/v1/orders/seller", "read"), r.SellerOrders)
	add(api, "GET", "/orders", r.Authenticate, authorize(r, "/api/v1/orders", "read"), r.ListOrders)
	add(api, "GET", "/orders/:id", r.Authenticate, authorize(r, "/api/v1/orders", "read"), r.GetOrder)
	add(api, "POST", "/orders/:id/cancel", r.Authenticate, authorize(r, "/api/v1/orders", "write"), r.CancelOrder)
	add(api, "GET", "/warehouse/receptions", r.Authenticate, authorize(r, "/api/v1/warehouse/receptions", "read"), r.WarehouseReceptions)
	add(api, "GET", "/warehouse/products", r.Authenticate, authorize(r, "/api/v1/warehouse/receptions", "read"), r.WarehouseProducts)
	add(api, "POST", "/warehouse/receptions", r.Authenticate, authorize(r, "/api/v1/warehouse/receptions", "write"), r.WarehouseReceptionCreate)
	add(api, "POST", "/warehouse/receptions/:id/inspect", r.Authenticate, authorize(r, "/api/v1/warehouse/receptions", "write"), r.WarehouseInspect)
	add(api, "GET", "/warehouse/receptions/:id/evidence", r.Authenticate, authorize(r, "/api/v1/warehouse/receptions", "read"), r.WarehouseReceptionEvidenceList)
	add(api, "POST", "/warehouse/receptions/:id/evidence", r.Authenticate, authorize(r, "/api/v1/warehouse/receptions", "write"), r.WarehouseReceptionEvidence)
	add(api, "GET", "/warehouse/inventory", r.Authenticate, authorize(r, "/api/v1/warehouse/inventory", "read"), r.Inventory)
	add(api, "POST", "/warehouse/inventory/:productId/adjust", r.Authenticate, authorize(r, "/api/v1/warehouse/inventory", "write"), r.InventoryAdjust)
	add(api, "GET", "/cart", r.Authenticate, authorize(r, "/api/v1/cart", "read"), r.Cart)
	add(api, "POST", "/cart/items", r.Authenticate, authorize(r, "/api/v1/cart", "write"), r.CartAdd)
	add(api, "POST", "/cart/merge", r.Authenticate, authorize(r, "/api/v1/cart", "write"), r.CartMerge)
	add(api, "PATCH", "/cart/items/:productId", r.Authenticate, authorize(r, "/api/v1/cart", "write"), r.CartSet)
	add(api, "DELETE", "/cart/items/:productId", r.Authenticate, authorize(r, "/api/v1/cart", "write"), r.CartRemove)
	add(api, "GET", "/wishlist", r.Authenticate, authorize(r, "/api/v1/wishlist", "read"), r.Wishlist)
	add(api, "POST", "/wishlist/items/:productId", r.Authenticate, authorize(r, "/api/v1/wishlist", "write"), r.WishlistAdd)
	add(api, "DELETE", "/wishlist/items/:productId", r.Authenticate, authorize(r, "/api/v1/wishlist", "write"), r.WishlistRemove)
	add(api, "GET", "/notifications", r.Authenticate, authorize(r, "/api/v1/notifications", "read"), r.Notifications)
	add(api, "POST", "/notifications/:id/read", r.Authenticate, authorize(r, "/api/v1/notifications", "write"), r.NotificationRead)
	add(api, "POST", "/notifications/read-all", r.Authenticate, authorize(r, "/api/v1/notifications", "write"), r.NotificationsReadAll)
	add(api, "GET", "/notifications/ws-ticket", r.Authenticate, authorize(r, "/api/v1/notifications", "read"), r.NotificationSocketTicket)
	add(api, "GET", "/notifications/ws", r.NotificationSocket)
}

func authorize(r AuthenticatedRoutes, resource, action string) fiber.Handler {
	if r.Authorize == nil {
		return nil
	}
	return r.Authorize(resource, action)
}
