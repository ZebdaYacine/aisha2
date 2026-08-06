package container

import (
	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue"
	customer "github.com/aisha-platform/aisha/apps/api/internal/features/user"
)

// Features is the application-facing feature graph assembled by the
// container. Concrete repositories never cross this boundary.
type Features struct {
	Auth      *auth.Service
	Customer  *customer.Service
	Artisan   *artisan.Service
	Catalogue *catalogue.Service
}
