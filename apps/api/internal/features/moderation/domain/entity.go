package domain

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

var (
	ErrValidation         = errors.New("moderation validation failed")
	ErrNotFound           = errors.New("moderation target not found")
	ErrInvalidTransition  = errors.New("invalid moderation transition")
	ErrActivationNotReady = errors.New("product activation gates are not satisfied")
)

type QueueItem struct {
	SubmissionID  string          `json:"submissionId"`
	ProductID     string          `json:"productId"`
	ProductCode   string          `json:"productCode,omitempty"`
	ProductName   string          `json:"productName"`
	Version       int             `json:"version"`
	Snapshot      json.RawMessage `json:"snapshot"`
	ProductStatus string          `json:"productStatus"`
	PriceMinor    int64           `json:"priceMinor"`
	Currency      string          `json:"currency"`
	ArtisanName   string          `json:"artisanName"`
	WorkshopName  string          `json:"workshopName"`
	Media         []Media         `json:"media"`
}

type Media struct {
	ID               string `json:"id"`
	MediaKind        string `json:"mediaKind"`
	OriginalFilename string `json:"originalFilename"`
	MediaType        string `json:"mediaType"`
	SizeBytes        int64  `json:"sizeBytes"`
	AltText          string `json:"altText"`
	Visibility       string `json:"visibility"`
	URL              string `json:"url,omitempty"`
	ObjectKey        string `json:"-"`
}
type DecisionInput struct {
	ID     string
	Action string
	Reason string
}
type Repository interface {
	ListQueue(context.Context, string, int, int) ([]QueueItem, int, error)
	Decide(context.Context, string, DecisionInput) (QueueItem, error)
}
type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
