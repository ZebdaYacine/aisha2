package domain

import (
	"context"
	"errors"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

var (
	ErrValidation  = errors.New("cart validation failed")
	ErrNotFound    = errors.New("cart item not found")
	ErrUnavailable = errors.New("product is not available for cart")
)

type Item struct {
	ProductID    string    `json:"productId"`
	ProductName  string    `json:"productName"`
	ArtisanName  string    `json:"artisanName"`
	WorkshopName string    `json:"workshopName"`
	Image        string    `json:"image,omitempty"`
	Quantity     int       `json:"quantity"`
	PriceMinor   int64     `json:"priceMinor"`
	Currency     string    `json:"currency"`
	Available    int64     `json:"availableQuantity"`
	Active       bool      `json:"active"`
	Warning      string    `json:"warning,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Input struct {
	ProductID string
	Quantity  int
}

type Repository interface {
	List(context.Context, string) ([]Item, error)
	Add(context.Context, string, Input) ([]Item, error)
	Set(context.Context, string, string, int) ([]Item, error)
	Remove(context.Context, string, string) ([]Item, error)
	Merge(context.Context, string, []Input) ([]Item, error)
}

type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
