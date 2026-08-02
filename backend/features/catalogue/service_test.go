package catalogue

import (
	"context"
	"testing"
)

type repositoryStub struct{ request PageRequest }

func (r *repositoryStub) Categories(context.Context, PageRequest) (Page[Category], error) {
	return Page[Category]{}, nil
}
func (r *repositoryStub) Products(_ context.Context, request PageRequest, _ string) (Page[Product], error) {
	r.request = request
	return Page[Product]{}, nil
}
func (r *repositoryStub) Product(context.Context, string, string) (Product, error) {
	return Product{}, nil
}
func (r *repositoryStub) Artisans(context.Context, PageRequest) (Page[Artisan], error) {
	return Page[Artisan]{}, nil
}
func (r *repositoryStub) Artisan(context.Context, string, string) (Artisan, error) {
	return Artisan{}, nil
}

func TestProductsNormalizesPaginationAndLocale(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository)
	_, err := service.Products(context.Background(), PageRequest{Locale: "xx", Page: -1, PageSize: 1000}, "")
	if err != nil {
		t.Fatalf("Products() error = %v", err)
	}
	if repository.request.Locale != "en" || repository.request.Page != 1 || repository.request.PageSize != 100 {
		t.Fatalf("normalized request = %#v", repository.request)
	}
}
