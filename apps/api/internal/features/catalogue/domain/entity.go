package domain

import (
	"context"
	"errors"
	"time"
)

var supportedLocales = map[string]bool{"ar": true, "en": true, "es": true, "fr": true}

var ErrNotFound = errors.New("catalogue resource not found")

type PageRequest struct {
	Locale   string
	Page     int
	PageSize int
}

func (p PageRequest) Normalized() PageRequest {
	if !supportedLocales[p.Locale] {
		p.Locale = "en"
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p
}

type Category struct{ ID, Slug, Name string }
type Product struct {
	ID, ArtisanID, ArtisanName, WorkshopID, WorkshopName, CategoryID, CategorySlug, Name, Description, Story string
	Materials, ProductionMethod, Region, Currency, Status                                                    string
	PriceMinor, AvailableQuantity                                                                            int64
	MadeToOrderEligible                                                                                      bool
	Media                                                                                                    []string
	PublishedAt                                                                                              *time.Time
}
type Artisan struct {
	ID, Name, WorkshopID, Workshop, Wilaya, Location, Biography, Craft string
	Media                                                              []string
	ProductCount                                                       int
}

type Workshop struct {
	ID, Name, Description, Wilaya, Location, Craft, ArtisanID, ArtisanName string
	Media                                                                  []string
	ProductCount                                                           int
}
type Page[T any] struct {
	Items                 []T
	Page, PageSize, Total int
}

type Repository interface {
	Categories(context.Context, PageRequest) (Page[Category], error)
	Products(context.Context, PageRequest, string, string, string) (Page[Product], error)
	Product(context.Context, string, string) (Product, error)
	Artisans(context.Context, PageRequest) (Page[Artisan], error)
	Artisan(context.Context, string, string) (Artisan, error)
	Workshops(context.Context, PageRequest) (Page[Workshop], error)
	Workshop(context.Context, string, string) (Workshop, error)
}
