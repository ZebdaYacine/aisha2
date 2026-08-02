package catalogue

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("catalogue resource not found")

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Categories(ctx context.Context, p PageRequest) (Page[Category], error) {
	rows, err := r.pool.Query(ctx, `SELECT id,slug,display_name,count(*) OVER() FROM categories WHERE is_active=true ORDER BY sort_order,slug LIMIT $1 OFFSET $2`, p.PageSize, (p.Page-1)*p.PageSize)
	if err != nil {
		return Page[Category]{}, fmt.Errorf("query public categories: %w", err)
	}
	defer rows.Close()
	result := Page[Category]{Items: []Category{}, Page: p.Page, PageSize: p.PageSize}
	for rows.Next() {
		var item Category
		if err := rows.Scan(&item.ID, &item.Slug, &item.Name, &result.Total); err != nil {
			return Page[Category]{}, fmt.Errorf("scan public category: %w", err)
		}
		result.Items = append(result.Items, item)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) Products(ctx context.Context, p PageRequest, category string) (Page[Product], error) {
	rows, err := r.pool.Query(ctx, `SELECT p.id,p.artisan_profile_id,a.public_display_name,c.id,c.slug,COALESCE(t.name,e.name,''),COALESCE(t.description,e.description,''),COALESCE(t.story,e.story,''),COALESCE(p.materials,''),COALESCE(p.production_method,''),COALESCE(p.region_of_origin,''),p.currency,p.price_minor,p.status,p.published_at,count(*) OVER() FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN categories c ON c.id=p.category_id LEFT JOIN product_translations t ON t.product_id=p.id AND t.locale=$1 LEFT JOIN product_translations e ON e.product_id=p.id AND e.locale='en' WHERE p.status='ACTIVE' AND p.published_at IS NOT NULL AND a.status='APPROVED' AND ($2='' OR c.slug=$2) ORDER BY p.published_at DESC,p.id LIMIT $3 OFFSET $4`, p.Locale, category, p.PageSize, (p.Page-1)*p.PageSize)
	if err != nil {
		return Page[Product]{}, fmt.Errorf("query public products: %w", err)
	}
	defer rows.Close()
	result := Page[Product]{Items: []Product{}, Page: p.Page, PageSize: p.PageSize}
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.ArtisanID, &item.ArtisanName, &item.CategoryID, &item.CategorySlug, &item.Name, &item.Description, &item.Story, &item.Materials, &item.ProductionMethod, &item.Region, &item.Currency, &item.PriceMinor, &item.Status, &item.PublishedAt, &result.Total); err != nil {
			return Page[Product]{}, fmt.Errorf("scan public product: %w", err)
		}
		item.Media, err = r.productMedia(ctx, item.ID)
		if err != nil {
			return Page[Product]{}, err
		}
		result.Items = append(result.Items, item)
	}
	return result, rows.Err()
}
func (r *PostgresRepository) Product(ctx context.Context, id, locale string) (Product, error) {
	var item Product
	err := r.pool.QueryRow(ctx, `SELECT p.id,p.artisan_profile_id,a.public_display_name,c.id,c.slug,COALESCE(t.name,e.name,''),COALESCE(t.description,e.description,''),COALESCE(t.story,e.story,''),COALESCE(p.materials,''),COALESCE(p.production_method,''),COALESCE(p.region_of_origin,''),p.currency,p.price_minor,p.status,p.published_at FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN categories c ON c.id=p.category_id LEFT JOIN product_translations t ON t.product_id=p.id AND t.locale=$2 LEFT JOIN product_translations e ON e.product_id=p.id AND e.locale='en' WHERE p.id::text=$1 AND p.status='ACTIVE' AND p.published_at IS NOT NULL AND a.status='APPROVED'`, id, locale).Scan(&item.ID, &item.ArtisanID, &item.ArtisanName, &item.CategoryID, &item.CategorySlug, &item.Name, &item.Description, &item.Story, &item.Materials, &item.ProductionMethod, &item.Region, &item.Currency, &item.PriceMinor, &item.Status, &item.PublishedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrNotFound
	}
	if err != nil {
		return Product{}, fmt.Errorf("query public product: %w", err)
	}
	item.Media, err = r.productMedia(ctx, item.ID)
	return item, err
}
func (r *PostgresRepository) productMedia(ctx context.Context, id string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT object_key FROM product_media WHERE product_id=$1 AND visibility='PUBLIC' AND media_kind='IMAGE' ORDER BY sort_order,created_at`, id)
	if err != nil {
		return nil, fmt.Errorf("query public product media: %w", err)
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		items = append(items, key)
	}
	return items, rows.Err()
}
func (r *PostgresRepository) Artisans(ctx context.Context, p PageRequest) (Page[Artisan], error) {
	rows, err := r.pool.Query(ctx, `SELECT a.id,a.public_display_name,COALESCE(a.workshop_name,''),COALESCE(a.wilaya,''),COALESCE(a.location_text,''),COALESCE(t.biography,e.biography,''),(SELECT count(*) FROM products p WHERE p.artisan_profile_id=a.id AND p.status='ACTIVE' AND p.published_at IS NOT NULL),count(*) OVER() FROM artisan_profiles a LEFT JOIN artisan_profile_translations t ON t.artisan_profile_id=a.id AND t.locale=$1 LEFT JOIN artisan_profile_translations e ON e.artisan_profile_id=a.id AND e.locale='en' WHERE a.status='APPROVED' ORDER BY a.public_display_name LIMIT $2 OFFSET $3`, p.Locale, p.PageSize, (p.Page-1)*p.PageSize)
	if err != nil {
		return Page[Artisan]{}, fmt.Errorf("query public artisans: %w", err)
	}
	defer rows.Close()
	result := Page[Artisan]{Items: []Artisan{}, Page: p.Page, PageSize: p.PageSize}
	for rows.Next() {
		var item Artisan
		if err := rows.Scan(&item.ID, &item.Name, &item.Workshop, &item.Wilaya, &item.Location, &item.Biography, &item.ProductCount, &result.Total); err != nil {
			return Page[Artisan]{}, err
		}
		result.Items = append(result.Items, item)
	}
	return result, rows.Err()
}
func (r *PostgresRepository) Artisan(ctx context.Context, id, locale string) (Artisan, error) {
	var item Artisan
	err := r.pool.QueryRow(ctx, `SELECT a.id,a.public_display_name,COALESCE(a.workshop_name,''),COALESCE(a.wilaya,''),COALESCE(a.location_text,''),COALESCE(t.biography,e.biography,''),(SELECT count(*) FROM products p WHERE p.artisan_profile_id=a.id AND p.status='ACTIVE' AND p.published_at IS NOT NULL) FROM artisan_profiles a LEFT JOIN artisan_profile_translations t ON t.artisan_profile_id=a.id AND t.locale=$2 LEFT JOIN artisan_profile_translations e ON e.artisan_profile_id=a.id AND e.locale='en' WHERE a.id::text=$1 AND a.status='APPROVED'`, id, locale).Scan(&item.ID, &item.Name, &item.Workshop, &item.Wilaya, &item.Location, &item.Biography, &item.ProductCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return Artisan{}, ErrNotFound
	}
	if err != nil {
		return Artisan{}, fmt.Errorf("query public artisan: %w", err)
	}
	return item, nil
}
