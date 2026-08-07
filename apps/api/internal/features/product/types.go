package product

import "github.com/aisha-platform/aisha/apps/api/internal/features/product/domain"

type Translation = domain.Translation
type Input = domain.Input
type Media = domain.Media
type Product = domain.Product
type Repository = domain.Repository
type Authorizer = domain.Authorizer

var (
	ErrNotFound           = domain.ErrNotFound
	ErrValidation         = domain.ErrValidation
	ErrNotEditable        = domain.ErrNotEditable
	ErrArtisanNotApproved = domain.ErrArtisanNotApproved
	ErrInvalidTransition  = domain.ErrInvalidTransition
	ErrMediaNotFound      = domain.ErrMediaNotFound
)
