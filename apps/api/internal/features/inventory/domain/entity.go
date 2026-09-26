package domain

import (
	"context"
	"errors"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

var (
	ErrValidation        = errors.New("inventory validation failed")
	ErrNotFound          = errors.New("inventory resource not found")
	ErrDuplicate         = errors.New("inventory reference already exists")
	ErrInsufficientStock = errors.New("inventory available quantity is insufficient")
)

type Balance struct {
	ProductID    string    `json:"productId"`
	ProductName  string    `json:"productName"`
	WorkshopID   string    `json:"workshopId"`
	WorkshopName string    `json:"workshopName"`
	ArtisanID    string    `json:"artisanId"`
	ArtisanName  string    `json:"artisanName"`
	OnHand       int64     `json:"onHand"`
	Available    int64     `json:"available"`
	Reserved     int64     `json:"reserved"`
	Quarantined  int64     `json:"quarantined"`
	Damaged      int64     `json:"damaged"`
	Rejected     int64     `json:"rejected"`
	Shipped      int64     `json:"shipped"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Movement struct {
	ID            string    `json:"id"`
	ProductID     string    `json:"productId"`
	MovementType  string    `json:"movementType"`
	QuantityDelta int64     `json:"quantityDelta"`
	StockBucket   string    `json:"stockBucket"`
	ReferenceKey  string    `json:"referenceKey"`
	Reason        string    `json:"reason"`
	ActorUserID   string    `json:"actorUserId,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

type AdjustmentInput struct {
	QuantityDelta int64
	Reason        string
	ReferenceKey  string
}

type Repository interface {
	ListBalances(context.Context, string, bool, string, int, int) ([]Balance, int, error)
	Adjust(context.Context, string, bool, string, AdjustmentInput) (Movement, error)
}

type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
