package container

import (
	"github.com/aisha-platform/aisha/apps/api/internal/features/admin"
	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/cart"
	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue"
	"github.com/aisha-platform/aisha/apps/api/internal/features/inventory"
	"github.com/aisha-platform/aisha/apps/api/internal/features/moderation"
	"github.com/aisha-platform/aisha/apps/api/internal/features/order"
	"github.com/aisha-platform/aisha/apps/api/internal/features/product"
	customer "github.com/aisha-platform/aisha/apps/api/internal/features/user"
	"github.com/aisha-platform/aisha/apps/api/internal/features/warehouse"
	"github.com/aisha-platform/aisha/apps/api/internal/features/wishlist"
)

// Features is the application-facing feature graph assembled by the
// container. Concrete repositories never cross this boundary.
type Features struct {
	Auth       *auth.Service
	Customer   *customer.Service
	Artisan    *artisan.Service
	Catalogue  *catalogue.Service
	Product    *product.Service
	Admin      *admin.Service
	Moderation *moderation.Service
	Order      *order.Service
	Warehouse  *warehouse.Service
	Inventory  *inventory.Service
	Cart       *cart.Service
	Wishlist   *wishlist.Service
}
