package catalogue

import "github.com/aisha-platform/aisha/apps/api/internal/features/catalogue/domain"

type PageRequest = domain.PageRequest
type Category = domain.Category
type Product = domain.Product
type Artisan = domain.Artisan
type Page[T any] = domain.Page[T]
type Repository = domain.Repository

var ErrNotFound = domain.ErrNotFound
