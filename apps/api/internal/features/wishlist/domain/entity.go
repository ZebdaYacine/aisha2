package domain

import (
	"context"
	"errors"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"time"
)

var ErrValidation = errors.New("wishlist validation failed")

type Item struct {
	ProductID   string    `json:"productId"`
	ProductName string    `json:"productName"`
	PriceMinor  int64     `json:"priceMinor"`
	Currency    string    `json:"currency"`
	Image       string    `json:"image,omitempty"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"createdAt"`
}
type Repository interface {
	List(context.Context, string) ([]Item, error)
	Add(context.Context, string, string) ([]Item, error)
	Remove(context.Context, string, string) ([]Item, error)
}
type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
