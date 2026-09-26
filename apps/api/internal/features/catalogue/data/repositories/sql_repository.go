package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/aisha-platform/aisha/apps/api/internal/features/catalogue/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PageRequest = domain.PageRequest
type Category = domain.Category
type Product = domain.Product
type Artisan = domain.Artisan
type Workshop = domain.Workshop
type Page[T any] = domain.Page[T]

var ErrNotFound = domain.ErrNotFound

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Categories(ctx context.Context, p PageRequest) (Page[Category], error) {
	rows, err := r.pool.Query(ctx, `SELECT c.id,c.slug,COALESCE(t.name,e.name,c.display_name),count(*) OVER() FROM categories c LEFT JOIN category_translations t ON t.category_id=c.id AND t.locale=$1 LEFT JOIN category_translations e ON e.category_id=c.id AND e.locale='en' WHERE c.is_active=true ORDER BY c.sort_order,c.slug LIMIT $2 OFFSET $3`, p.Locale, p.PageSize, (p.Page-1)*p.PageSize)
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

const publicWorkshopJoin = `LEFT JOIN workshops w ON w.id=COALESCE(p.workshop_id,(SELECT wd.id FROM workshops wd WHERE wd.artisan_profile_id=a.id AND wd.is_default=true LIMIT 1))`
const publicProductSelect = `p.id,p.artisan_profile_id,COALESCE(at.display_name,ae.display_name,a.public_display_name),COALESCE(w.id::text,''),COALESCE(wt.name,we.name,w.name,a.workshop_name,''),c.id,c.slug,COALESCE(t.name,e.name,''),COALESCE(t.description,e.description,''),COALESCE(t.story,e.story,''),COALESCE(p.materials,''),COALESCE(p.production_method,''),COALESCE(p.region_of_origin,''),p.currency,p.price_minor,p.status,p.published_at,p.made_to_order_eligible,GREATEST(COALESCE((SELECT SUM(quantity_delta) FROM inventory_movements im WHERE im.product_id=p.id AND im.stock_bucket='AVAILABLE'),0),0)`

func (r *PostgresRepository) Products(ctx context.Context, p PageRequest, category, query, workshop string) (Page[Product], error) {
	rows, err := r.pool.Query(ctx, `SELECT `+publicProductSelect+`,count(*) OVER() FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN categories c ON c.id=p.category_id `+publicWorkshopJoin+` LEFT JOIN workshop_translations wt ON wt.workshop_id=w.id AND wt.locale=$1 LEFT JOIN workshop_translations we ON we.workshop_id=w.id AND we.locale='en' LEFT JOIN product_translations t ON t.product_id=p.id AND t.locale=$1 LEFT JOIN product_translations e ON e.product_id=p.id AND e.locale='en' LEFT JOIN artisan_profile_translations at ON at.artisan_profile_id=a.id AND at.locale=$1 LEFT JOIN artisan_profile_translations ae ON ae.artisan_profile_id=a.id AND ae.locale='en' WHERE p.status='ACTIVE' AND p.published_at IS NOT NULL AND a.status='APPROVED' AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE') AND (w.id IS NULL OR (w.status='ACTIVE' AND w.is_public=true)) AND ($2='' OR c.slug=$2) AND ($3='' OR lower(COALESCE(t.name,e.name,'')) LIKE '%'||lower($3)||'%' OR lower(COALESCE(t.description,e.description,'')) LIKE '%'||lower($3)||'%' OR lower(COALESCE(p.region_of_origin,'')) LIKE '%'||lower($3)||'%' OR lower(COALESCE(p.materials,'')) LIKE '%'||lower($3)||'%') AND ($4='' OR w.id::text=$4) ORDER BY p.published_at DESC,p.id LIMIT $5 OFFSET $6`, p.Locale, category, query, workshop, p.PageSize, (p.Page-1)*p.PageSize)
	if err != nil {
		return Page[Product]{}, fmt.Errorf("query public products: %w", err)
	}
	defer rows.Close()
	result := Page[Product]{Items: []Product{}, Page: p.Page, PageSize: p.PageSize}
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.ArtisanID, &item.ArtisanName, &item.WorkshopID, &item.WorkshopName, &item.CategoryID, &item.CategorySlug, &item.Name, &item.Description, &item.Story, &item.Materials, &item.ProductionMethod, &item.Region, &item.Currency, &item.PriceMinor, &item.Status, &item.PublishedAt, &item.MadeToOrderEligible, &item.AvailableQuantity, &result.Total); err != nil {
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
	err := r.pool.QueryRow(ctx, `SELECT `+publicProductSelect+` FROM products p JOIN artisan_profiles a ON a.id=p.artisan_profile_id JOIN categories c ON c.id=p.category_id `+publicWorkshopJoin+` LEFT JOIN workshop_translations wt ON wt.workshop_id=w.id AND wt.locale=$2 LEFT JOIN workshop_translations we ON we.workshop_id=w.id AND we.locale='en' LEFT JOIN product_translations t ON t.product_id=p.id AND t.locale=$2 LEFT JOIN product_translations e ON e.product_id=p.id AND e.locale='en' LEFT JOIN artisan_profile_translations at ON at.artisan_profile_id=a.id AND at.locale=$2 LEFT JOIN artisan_profile_translations ae ON ae.artisan_profile_id=a.id AND ae.locale='en' WHERE p.id::text=$1 AND p.status='ACTIVE' AND p.published_at IS NOT NULL AND a.status='APPROVED' AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE') AND (w.id IS NULL OR (w.status='ACTIVE' AND w.is_public=true))`, id, locale).Scan(&item.ID, &item.ArtisanID, &item.ArtisanName, &item.WorkshopID, &item.WorkshopName, &item.CategoryID, &item.CategorySlug, &item.Name, &item.Description, &item.Story, &item.Materials, &item.ProductionMethod, &item.Region, &item.Currency, &item.PriceMinor, &item.Status, &item.PublishedAt, &item.MadeToOrderEligible, &item.AvailableQuantity)
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
	rows, err := r.pool.Query(ctx, `SELECT a.id,COALESCE(w.id::text,''),COALESCE(t.display_name,e.display_name,a.public_display_name),COALESCE(w.name,a.workshop_name,''),COALESCE(a.wilaya,''),COALESCE(a.location_text,''),COALESCE(t.biography,e.biography,''),COALESCE((SELECT string_agg(COALESCE(ct.name,ce.name,c.display_name), ', ' ORDER BY c.sort_order,c.slug) FROM artisan_profile_categories apc JOIN categories c ON c.id=apc.category_id LEFT JOIN category_translations ct ON ct.category_id=c.id AND ct.locale=$1 LEFT JOIN category_translations ce ON ce.category_id=c.id AND ce.locale='en' WHERE apc.artisan_profile_id=a.id AND c.is_active=true),''),COALESCE(a.profile_image_object_key,''),(SELECT count(*) FROM products p WHERE p.artisan_profile_id=a.id AND p.status='ACTIVE' AND p.published_at IS NOT NULL),count(*) OVER() FROM artisan_profiles a LEFT JOIN workshops w ON w.artisan_profile_id=a.id AND w.is_default=true LEFT JOIN artisan_profile_translations t ON t.artisan_profile_id=a.id AND t.locale=$1 LEFT JOIN artisan_profile_translations e ON e.artisan_profile_id=a.id AND e.locale='en' WHERE a.status='APPROVED' AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE') ORDER BY a.public_display_name LIMIT $2 OFFSET $3`, p.Locale, p.PageSize, (p.Page-1)*p.PageSize)
	if err != nil {
		return Page[Artisan]{}, fmt.Errorf("query public artisans: %w", err)
	}
	defer rows.Close()
	result := Page[Artisan]{Items: []Artisan{}, Page: p.Page, PageSize: p.PageSize}
	for rows.Next() {
		var item Artisan
		var imageKey string
		if err := rows.Scan(&item.ID, &item.WorkshopID, &item.Name, &item.Workshop, &item.Wilaya, &item.Location, &item.Biography, &item.Craft, &imageKey, &item.ProductCount, &result.Total); err != nil {
			return Page[Artisan]{}, err
		}
		if imageKey != "" {
			item.Media = []string{imageKey}
		}
		result.Items = append(result.Items, item)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) Artisan(ctx context.Context, id, locale string) (Artisan, error) {
	var item Artisan
	var imageKey string
	err := r.pool.QueryRow(ctx, `SELECT a.id,COALESCE(w.id::text,''),COALESCE(t.display_name,e.display_name,a.public_display_name),COALESCE(w.name,a.workshop_name,''),COALESCE(a.wilaya,''),COALESCE(a.location_text,''),COALESCE(t.biography,e.biography,''),COALESCE((SELECT string_agg(COALESCE(ct.name,ce.name,c.display_name), ', ' ORDER BY c.sort_order,c.slug) FROM artisan_profile_categories apc JOIN categories c ON c.id=apc.category_id LEFT JOIN category_translations ct ON ct.category_id=c.id AND ct.locale=$2 LEFT JOIN category_translations ce ON ce.category_id=c.id AND ce.locale='en' WHERE apc.artisan_profile_id=a.id AND c.is_active=true),''),COALESCE(a.profile_image_object_key,''),(SELECT count(*) FROM products p WHERE p.artisan_profile_id=a.id AND p.status='ACTIVE' AND p.published_at IS NOT NULL) FROM artisan_profiles a LEFT JOIN workshops w ON w.artisan_profile_id=a.id AND w.is_default=true LEFT JOIN artisan_profile_translations t ON t.artisan_profile_id=a.id AND t.locale=$2 LEFT JOIN artisan_profile_translations e ON e.artisan_profile_id=a.id AND e.locale='en' WHERE a.id::text=$1 AND a.status='APPROVED' AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE')`, id, locale).Scan(&item.ID, &item.WorkshopID, &item.Name, &item.Workshop, &item.Wilaya, &item.Location, &item.Biography, &item.Craft, &imageKey, &item.ProductCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return Artisan{}, ErrNotFound
	}
	if err != nil {
		return Artisan{}, fmt.Errorf("query public artisan: %w", err)
	}
	if imageKey != "" {
		item.Media = []string{imageKey}
	}
	return item, nil
}

const publicWorkshopSelect = `w.id,COALESCE(wt.name,we.name,w.name),COALESCE(wt.description,we.description,a.public_display_name),COALESCE(w.wilaya,a.wilaya,''),COALESCE(w.location_text,a.location_text,''),COALESCE((SELECT string_agg(COALESCE(ct.name,ce.name,c.display_name), ', ' ORDER BY c.sort_order,c.slug) FROM artisan_profile_categories apc JOIN categories c ON c.id=apc.category_id LEFT JOIN category_translations ct ON ct.category_id=c.id AND ct.locale=$1 LEFT JOIN category_translations ce ON ce.category_id=c.id AND ce.locale='en' WHERE apc.artisan_profile_id=a.id AND c.is_active=true),''),a.id,COALESCE(at.display_name,ae.display_name,a.public_display_name),COALESCE(a.profile_image_object_key,''),(SELECT count(*) FROM products p WHERE (p.workshop_id=w.id OR (p.workshop_id IS NULL AND w.is_default=true)) AND p.status='ACTIVE' AND p.published_at IS NOT NULL)`

func (r *PostgresRepository) Workshops(ctx context.Context, p PageRequest) (Page[Workshop], error) {
	rows, err := r.pool.Query(ctx, `SELECT `+publicWorkshopSelect+`,count(*) OVER() FROM workshops w JOIN artisan_profiles a ON a.id=w.artisan_profile_id LEFT JOIN workshop_translations wt ON wt.workshop_id=w.id AND wt.locale=$1 LEFT JOIN workshop_translations we ON we.workshop_id=w.id AND we.locale='en' LEFT JOIN artisan_profile_translations at ON at.artisan_profile_id=a.id AND at.locale=$1 LEFT JOIN artisan_profile_translations ae ON ae.artisan_profile_id=a.id AND ae.locale='en' WHERE w.status='ACTIVE' AND w.is_public=true AND a.status='APPROVED' AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE') ORDER BY w.name LIMIT $2 OFFSET $3`, p.Locale, p.PageSize, (p.Page-1)*p.PageSize)
	if err != nil {
		return Page[Workshop]{}, fmt.Errorf("query public workshops: %w", err)
	}
	defer rows.Close()
	result := Page[Workshop]{Items: []Workshop{}, Page: p.Page, PageSize: p.PageSize}
	for rows.Next() {
		var item Workshop
		var imageKey string
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Wilaya, &item.Location, &item.Craft, &item.ArtisanID, &item.ArtisanName, &imageKey, &item.ProductCount, &result.Total); err != nil {
			return Page[Workshop]{}, err
		}
		if imageKey != "" {
			item.Media = []string{imageKey}
		}
		result.Items = append(result.Items, item)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) Workshop(ctx context.Context, id, locale string) (Workshop, error) {
	var item Workshop
	var imageKey string
	err := r.pool.QueryRow(ctx, `SELECT `+publicWorkshopSelect+` FROM workshops w JOIN artisan_profiles a ON a.id=w.artisan_profile_id LEFT JOIN workshop_translations wt ON wt.workshop_id=w.id AND wt.locale=$2 LEFT JOIN workshop_translations we ON we.workshop_id=w.id AND we.locale='en' LEFT JOIN artisan_profile_translations at ON at.artisan_profile_id=a.id AND at.locale=$2 LEFT JOIN artisan_profile_translations ae ON ae.artisan_profile_id=a.id AND ae.locale='en' WHERE w.id::text=$1 AND w.status='ACTIVE' AND w.is_public=true AND a.status='APPROVED' AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE')`, id, locale).Scan(&item.ID, &item.Name, &item.Description, &item.Wilaya, &item.Location, &item.Craft, &item.ArtisanID, &item.ArtisanName, &imageKey, &item.ProductCount)
	if errors.Is(err, pgx.ErrNoRows) {
		return Workshop{}, ErrNotFound
	}
	if err != nil {
		return Workshop{}, fmt.Errorf("query public workshop: %w", err)
	}
	if imageKey != "" {
		item.Media = []string{imageKey}
	}
	return item, nil
}
