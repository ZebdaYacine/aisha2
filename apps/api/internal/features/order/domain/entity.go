package domain

import (
	"context"
	"errors"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"time"
)

var (
	ErrValidation          = errors.New("order validation failed")
	ErrNotFound            = errors.New("order not found")
	ErrOutOfStock          = errors.New("product is out of stock")
	ErrPriceChanged        = errors.New("product price changed")
	ErrInvalidTransition   = errors.New("invalid order transition")
	ErrIdempotencyConflict = errors.New("idempotency conflict")
)

type CartItem struct {
	ProductID          string `json:"productId"`
	Quantity           int    `json:"quantity"`
	ExpectedPriceMinor *int64 `json:"expectedPriceMinor,omitempty"`
	ExpectedCurrency   string `json:"currency,omitempty"`
}
type AddressSnapshot struct {
	FullName   string `json:"fullName"`
	Phone      string `json:"phone,omitempty"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}
type OrderItem struct {
	ID             string `json:"id"`
	ProductID      string `json:"productId"`
	ProductName    string `json:"productName"`
	ArtisanID      string `json:"artisanId"`
	ArtisanName    string `json:"artisanName"`
	WorkshopID     string `json:"workshopId"`
	WorkshopName   string `json:"workshopName"`
	UnitPriceMinor int64  `json:"unitPriceMinor"`
	Currency       string `json:"currency"`
	Quantity       int    `json:"quantity"`
	SubtotalMinor  int64  `json:"subtotalMinor"`
}
type PaymentAttempt struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	AmountMinor int64  `json:"amountMinor"`
	Currency    string `json:"currency"`
}
type ShipmentEvent struct {
	ID               string    `json:"id"`
	Status           string    `json:"status"`
	TrackingReference string   `json:"trackingReference,omitempty"`
	OccurredAt       time.Time `json:"occurredAt"`
}
type Order struct {
	ID            string          `json:"id"`
	OrderNumber   string          `json:"orderNumber"`
	UserID        string          `json:"userId,omitempty"`
	Status        string          `json:"status"`
	Currency      string          `json:"currency"`
	SubtotalMinor int64           `json:"subtotalMinor"`
	ShippingMinor int64           `json:"shippingMinor"`
	TotalMinor    int64           `json:"totalMinor"`
	Address       AddressSnapshot `json:"address"`
	Items         []OrderItem     `json:"items"`
	Payment       *PaymentAttempt `json:"payment,omitempty"`
	ShipmentEvents []ShipmentEvent `json:"shipmentEvents"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}
type SellerItem struct {
	OrderID       string    `json:"orderId"`
	OrderNumber   string    `json:"orderNumber"`
	OrderStatus   string    `json:"orderStatus"`
	ProductID     string    `json:"productId"`
	ProductName   string    `json:"productName"`
	Quantity      int       `json:"quantity"`
	SubtotalMinor int64     `json:"subtotalMinor"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"createdAt"`
}
type Return struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"orderId"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"createdAt"`
}
type Repository interface {
	Checkout(context.Context, string, string, []CartItem, string, string) (Order, error)
	ListMine(context.Context, string, int, int) ([]Order, int, error)
	GetMine(context.Context, string, string) (Order, error)
	ListSeller(context.Context, string, int, int) ([]SellerItem, int, error)
	Cancel(context.Context, string, string) (Order, error)
	RecordReturn(context.Context, string, string, string) (Return, error)
}
type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
