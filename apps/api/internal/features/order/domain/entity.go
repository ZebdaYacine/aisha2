package domain

import (
	"context"
	"errors"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"time"
)

var (
	ErrValidation            = errors.New("order validation failed")
	ErrNotFound              = errors.New("order not found")
	ErrOutOfStock            = errors.New("product is out of stock")
	ErrPriceChanged          = errors.New("product price changed")
	ErrPaymentAmountMismatch = errors.New("payment amount does not match order")
	ErrInvalidTransition     = errors.New("invalid order transition")
	ErrIdempotencyConflict   = errors.New("idempotency conflict")
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
	ID                string `json:"id"`
	Provider          string `json:"provider"`
	ProviderReference string `json:"providerReference,omitempty"`
	Status            string `json:"status"`
	AmountMinor       int64  `json:"amountMinor"`
	Currency          string `json:"currency"`
	FailureReason     string `json:"failureReason,omitempty"`
}
type ShipmentEvent struct {
	ID                string    `json:"id"`
	Status            string    `json:"status"`
	TrackingReference string    `json:"trackingReference,omitempty"`
	OccurredAt        time.Time `json:"occurredAt"`
}
type Order struct {
	ID             string          `json:"id"`
	OrderNumber    string          `json:"orderNumber"`
	UserID         string          `json:"userId,omitempty"`
	Status         string          `json:"status"`
	Currency       string          `json:"currency"`
	SubtotalMinor  int64           `json:"subtotalMinor"`
	ShippingMinor  int64           `json:"shippingMinor"`
	TotalMinor     int64           `json:"totalMinor"`
	Address        AddressSnapshot `json:"address"`
	Items          []OrderItem     `json:"items"`
	Payment        *PaymentAttempt `json:"payment,omitempty"`
	ShipmentEvents []ShipmentEvent `json:"shipmentEvents"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
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
type ShipmentInput struct {
	Carrier           string `json:"carrier"`
	TrackingReference string `json:"trackingReference"`
}

type FulfilmentOrder struct {
	ID            string    `json:"id"`
	OrderNumber   string    `json:"orderNumber"`
	Status        string    `json:"status"`
	Currency      string    `json:"currency"`
	TotalMinor    int64     `json:"totalMinor"`
	CustomerName  string    `json:"customerName"`
	CustomerEmail string    `json:"customerEmail"`
	CreatedAt     time.Time `json:"createdAt"`
}
type Repository interface {
	Checkout(context.Context, string, string, []CartItem, string, string) (Order, error)
	ListMine(context.Context, string, int, int) ([]Order, int, error)
	GetMine(context.Context, string, string) (Order, error)
	ListSeller(context.Context, string, int, int) ([]SellerItem, int, error)
	Cancel(context.Context, string, string) (Order, error)
	RecordReturn(context.Context, string, string, string) (Return, error)
	GetPayment(context.Context, string, string) (PaymentAttempt, error)
	ConfirmPayment(context.Context, string, string, string) (Order, error)
	FailPayment(context.Context, string, string, string) (Order, error)
	ListFulfilment(context.Context, int, int) ([]FulfilmentOrder, int, error)
	Prepare(context.Context, string, string) (Order, error)
	Ship(context.Context, string, string, ShipmentInput) (Order, error)
	Deliver(context.Context, string, string) (Order, error)
	Refund(context.Context, string, string, int64, string, string) (Order, error)
}
type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
