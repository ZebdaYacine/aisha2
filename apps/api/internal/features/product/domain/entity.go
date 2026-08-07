package domain

import (
	"context"
	"errors"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
)

var (
	ErrNotFound           = errors.New("product not found")
	ErrValidation         = errors.New("product validation failed")
	ErrNotEditable        = errors.New("product is not editable")
	ErrArtisanNotApproved = errors.New("artisan is not approved")
	ErrInvalidTransition  = errors.New("invalid product transition")
	ErrMediaNotFound      = errors.New("product media not found")
)

var SupportedLocales = []string{"ar", "fr", "en", "es"}

type Translation struct {
	Locale          string `json:"locale"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Story           string `json:"story"`
	CulturalContext string `json:"culturalContext"`
}

type Input struct {
	CategoryID          string
	ProductType         string
	PriceMinor          int64
	Currency            string
	Materials           string
	ProductionMethod    string
	IntendedUse         string
	Dimensions          string
	WeightGrams         *int
	CountryOfOrigin     string
	RegionOfOrigin      string
	EcoFriendlyVerified bool
	FairTradeVerified   bool
	MadeToOrderEligible bool
	Translations        []Translation
}

type Media struct {
	ID               string `json:"id"`
	ProductID        string `json:"-"`
	MediaKind        string `json:"mediaKind"`
	ObjectKey        string `json:"-"`
	Checksum         string `json:"-"`
	OriginalFilename string `json:"originalFilename"`
	MediaType        string `json:"mediaType"`
	SizeBytes        int64  `json:"sizeBytes"`
	AltText          string `json:"altText"`
	SortOrder        int    `json:"sortOrder"`
	Visibility       string `json:"visibility"`
	URL              string `json:"url,omitempty"`
}

type Product struct {
	ID                  string        `json:"id"`
	ArtisanID           string        `json:"artisanId"`
	CategoryID          string        `json:"categoryId"`
	ProductType         string        `json:"productType"`
	Status              string        `json:"status"`
	PriceMinor          int64         `json:"priceMinor"`
	Currency            string        `json:"currency"`
	Materials           string        `json:"materials"`
	ProductionMethod    string        `json:"productionMethod"`
	IntendedUse         string        `json:"intendedUse"`
	Dimensions          string        `json:"dimensions"`
	WeightGrams         *int          `json:"weightGrams,omitempty"`
	CountryOfOrigin     string        `json:"countryOfOrigin"`
	RegionOfOrigin      string        `json:"regionOfOrigin"`
	EcoFriendlyVerified bool          `json:"ecoFriendlyVerified"`
	FairTradeVerified   bool          `json:"fairTradeVerified"`
	MadeToOrderEligible bool          `json:"madeToOrderEligible"`
	Translations        []Translation `json:"translations"`
	Media               []Media       `json:"media"`
}

type Repository interface {
	Create(context.Context, string, Input) (Product, error)
	ListOwned(context.Context, string, string, int, int) ([]Product, int, error)
	GetOwned(context.Context, string, string) (Product, error)
	Update(context.Context, string, string, Input) (Product, error)
	Submit(context.Context, string, string) (Product, error)
	AddMedia(context.Context, string, string, Media) (Media, error)
	DeleteMedia(context.Context, string, string, string) (Media, error)
}

type Authorizer interface {
	Authorize(context.Context, auth.Principal, string, string) error
}
