package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/warehouse/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateReception(ctx context.Context, actor string, input domain.ReceptionInput) (domain.Reception, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Reception{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var productID, artisanID, workshopID string
	if err = tx.QueryRow(ctx, `SELECT p.id,p.artisan_profile_id,p.workshop_id FROM products p WHERE (p.id=NULLIF($1,'')::uuid OR p.product_code=NULLIF($2,'')) AND p.status IN ('APPROVED','ACTIVE','SUSPENDED')`, input.ProductID, input.ProductCode).Scan(&productID, &artisanID, &workshopID); errors.Is(err, pgx.ErrNoRows) {
		return domain.Reception{}, domain.ErrNotFound
	} else if err != nil {
		return domain.Reception{}, err
	}
	var id string
	err = tx.QueryRow(ctx, `INSERT INTO warehouse_receptions(product_id,artisan_profile_id,workshop_id,supplier_name,received_quantity,reference_key,parcel_reference,notes,received_at,received_by_user_id) VALUES($1,$2,$3,NULLIF($4,''),$5,$6,NULLIF($7,''),NULLIF($8,''),COALESCE($9,CURRENT_TIMESTAMP),$10) ON CONFLICT(reference_key) DO NOTHING RETURNING id`, productID, artisanID, workshopID, input.SupplierName, input.ReceivedQuantity, input.ReferenceKey, input.ParcelReference, input.Notes, input.ReceivedAt, actor).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		var existingProduct string
		var existingQuantity int64
		if err = tx.QueryRow(ctx, `SELECT id,product_id,received_quantity FROM warehouse_receptions WHERE reference_key=$1`, input.ReferenceKey).Scan(&id, &existingProduct, &existingQuantity); err != nil {
			return domain.Reception{}, err
		}
		if existingProduct != productID || existingQuantity != input.ReceivedQuantity {
			return domain.Reception{}, domain.ErrDuplicate
		}
		if err = tx.Commit(ctx); err != nil {
			return domain.Reception{}, err
		}
		return r.getReception(ctx, id)
	}
	if err != nil {
		return domain.Reception{}, fmt.Errorf("insert warehouse reception: %w", err)
	}
	if err = warehouseEvent(ctx, tx, actor, "WAREHOUSE_RECEPTION_CREATED", id, "", "RECEIVED_PENDING_INSPECTION", map[string]any{"productId": productID, "productCode": input.ProductCode, "receivedQuantity": input.ReceivedQuantity, "referenceKey": input.ReferenceKey}); err != nil {
		return domain.Reception{}, err
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Reception{}, err
	}
	return r.getReception(ctx, id)
}

func (r *PostgresRepository) ListReceptions(ctx context.Context, actor, status string, limit, offset int) ([]domain.Reception, int, error) {
	rows, err := r.pool.Query(ctx, `SELECT wr.id,wr.product_id,COALESCE(p.product_code,''),COALESCE(pt.name,p.product_type),wr.artisan_profile_id,COALESCE(a.public_display_name,''),COALESCE(u.phone,''),wr.workshop_id,COALESCE(w.name,''),COALESCE(wr.supplier_name,''),wr.received_quantity,wr.reference_key,COALESCE(wr.parcel_reference,''),COALESCE(wr.notes,''),wr.status,wr.received_at,wr.received_by_user_id,wr.inspected_at,COALESCE(wi.id::text,''),COALESCE(wi.accepted_quantity,0),COALESCE(wi.rejected_quantity,0),COALESCE(wi.quarantined_quantity,0),COALESCE(wi.damaged_quantity,0),COALESCE(wi.reason,''),COALESCE(wi.inspected_by_user_id::text,''),COALESCE(wi.inspected_at,wr.inspected_at),COUNT(*) OVER() FROM warehouse_receptions wr JOIN products p ON p.id=wr.product_id JOIN artisan_profiles a ON a.id=wr.artisan_profile_id JOIN users u ON u.id=a.user_id JOIN workshops w ON w.id=wr.workshop_id LEFT JOIN LATERAL (SELECT t.name FROM product_translations t WHERE t.product_id=p.id ORDER BY CASE WHEN t.locale='en' THEN 0 ELSE 1 END,t.locale LIMIT 1) pt ON true LEFT JOIN warehouse_inspections wi ON wi.reception_id=wr.id WHERE ($1='' OR wr.status=$1) ORDER BY wr.received_at DESC,wr.id LIMIT $2 OFFSET $3`, status, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list warehouse receptions: %w", err)
	}
	defer rows.Close()
	items := []domain.Reception{}
	total := 0
	for rows.Next() {
		var item domain.Reception
		var inspectionID, inspectedBy string
		var inspectedAt *time.Time
		var accepted, rejected, quarantined, damaged int64
		var reason string
		if err = rows.Scan(&item.ID, &item.ProductID, &item.ProductCode, &item.ProductName, &item.ArtisanProfileID, &item.ArtisanName, &item.ArtisanPhone, &item.WorkshopID, &item.WorkshopName, &item.SupplierName, &item.ReceivedQuantity, &item.ReferenceKey, &item.ParcelReference, &item.Notes, &item.Status, &item.ReceivedAt, &item.ReceivedBy, &item.InspectedAt, &inspectionID, &accepted, &rejected, &quarantined, &damaged, &reason, &inspectedBy, &inspectedAt, &total); err != nil {
			return nil, 0, err
		}
		if inspectionID != "" {
			item.Inspection = &domain.Inspection{ID: inspectionID, ReceptionID: item.ID, AcceptedQuantity: accepted, RejectedQuantity: rejected, QuarantinedQuantity: quarantined, DamagedQuantity: damaged, Reason: reason, InspectedBy: inspectedBy}
			if inspectedAt != nil {
				item.Inspection.InspectedAt = *inspectedAt
			}
		}
		item.Evidence, err = r.ListEvidence(ctx, item.ID)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) ListValidatedProducts(ctx context.Context, actor, artisanPhone, workshopID, query string, limit, offset int) ([]domain.ValidatedProduct, int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT p.id,COALESCE(p.product_code,''),COALESCE(pt.name,p.product_type),p.status,
			a.id,COALESCE(a.public_display_name,''),COALESCE(u.phone,''),w.id,COALESCE(w.name,''),
			p.price_minor,p.currency,GREATEST(COALESCE((SELECT SUM(im.quantity_delta) FROM inventory_movements im WHERE im.product_id=p.id AND im.stock_bucket='AVAILABLE'),0),0),COUNT(*) OVER()
		FROM products p
		JOIN artisan_profiles a ON a.id=p.artisan_profile_id
		JOIN users u ON u.id=a.user_id
		JOIN workshops w ON w.id=p.workshop_id
		LEFT JOIN LATERAL (
			SELECT t.name FROM product_translations t
			WHERE t.product_id=p.id
			ORDER BY CASE WHEN t.locale='en' THEN 0 ELSE 1 END,t.locale
			LIMIT 1
		) pt ON true
		WHERE p.product_code IS NOT NULL
		  AND p.product_code<>''
		  AND p.status IN ('APPROVED','ACTIVE','SUSPENDED')
		  AND (
			CASE
				WHEN regexp_replace(COALESCE(u.phone,''),'[^0-9]','','g') LIKE '00213%'
					THEN '0'||substring(regexp_replace(COALESCE(u.phone,''),'[^0-9]','','g') FROM 6)
				WHEN regexp_replace(COALESCE(u.phone,''),'[^0-9]','','g') LIKE '213%'
					THEN '0'||substring(regexp_replace(COALESCE(u.phone,''),'[^0-9]','','g') FROM 4)
				ELSE regexp_replace(COALESCE(u.phone,''),'[^0-9]','','g')
			END
		  ) = (
			CASE
				WHEN regexp_replace($1,'[^0-9]','','g') LIKE '00213%'
					THEN '0'||substring(regexp_replace($1,'[^0-9]','','g') FROM 6)
				WHEN regexp_replace($1,'[^0-9]','','g') LIKE '213%'
					THEN '0'||substring(regexp_replace($1,'[^0-9]','','g') FROM 4)
				ELSE regexp_replace($1,'[^0-9]','','g')
			END
		  )
		  AND w.id=COALESCE(NULLIF($2,'')::uuid,w.id)
		  AND ($3='' OR p.product_code ILIKE '%'||$3||'%' OR COALESCE(pt.name,p.product_type) ILIKE '%'||$3||'%')
		ORDER BY w.name,COALESCE(pt.name,p.product_type),p.id
		LIMIT $4 OFFSET $5`, artisanPhone, workshopID, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list validated warehouse products: %w", err)
	}
	defer rows.Close()
	items := []domain.ValidatedProduct{}
	total := 0
	for rows.Next() {
		var item domain.ValidatedProduct
		if err = rows.Scan(&item.ProductID, &item.ProductCode, &item.ProductName, &item.ProductStatus, &item.ArtisanProfileID, &item.ArtisanName, &item.ArtisanPhone, &item.WorkshopID, &item.WorkshopName, &item.PriceMinor, &item.Currency, &item.AvailableQuantity, &total); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *PostgresRepository) Inspect(ctx context.Context, actor, id string, input domain.InspectionInput) (domain.Inspection, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Inspection{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var received int64
	var status string
	if err = tx.QueryRow(ctx, `SELECT received_quantity,status FROM warehouse_receptions WHERE id=$1 FOR UPDATE`, id).Scan(&received, &status); errors.Is(err, pgx.ErrNoRows) {
		return domain.Inspection{}, domain.ErrNotFound
	} else if err != nil {
		return domain.Inspection{}, err
	}
	if status == "INSPECTED" {
		return r.getInspection(ctx, id)
	}
	if status != "RECEIVED_PENDING_INSPECTION" {
		return domain.Inspection{}, domain.ErrInvalidTransition
	}
	if input.AcceptedQuantity+input.RejectedQuantity+input.QuarantinedQuantity+input.DamagedQuantity != received {
		return domain.Inspection{}, domain.ErrQuantityMismatch
	}
	var evidenceCount int
	if err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM warehouse_evidence WHERE reception_id=$1`, id).Scan(&evidenceCount); err != nil {
		return domain.Inspection{}, err
	}
	if evidenceCount == 0 {
		return domain.Inspection{}, domain.ErrEvidenceRequired
	}
	var inspectionID string
	if err = tx.QueryRow(ctx, `INSERT INTO warehouse_inspections(reception_id,accepted_quantity,rejected_quantity,quarantined_quantity,damaged_quantity,reason,inspected_by_user_id) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`, id, input.AcceptedQuantity, input.RejectedQuantity, input.QuarantinedQuantity, input.DamagedQuantity, input.Reason, actor).Scan(&inspectionID); err != nil {
		return domain.Inspection{}, err
	}
	for _, movement := range []struct {
		kind, bucket string
		quantity     int64
	}{
		{"INSPECTION_ACCEPT", "AVAILABLE", input.AcceptedQuantity},
		{"INSPECTION_REJECT", "REJECTED", input.RejectedQuantity},
		{"INSPECTION_QUARANTINE", "QUARANTINED", input.QuarantinedQuantity},
		{"INSPECTION_DAMAGE", "DAMAGED", input.DamagedQuantity},
	} {
		if movement.quantity == 0 {
			continue
		}
		if _, err = tx.Exec(ctx, `INSERT INTO inventory_movements(product_id,movement_type,quantity_delta,reference_key,reason,actor_user_id,stock_bucket) SELECT product_id,$2,$3,$4,$5,$6,$7 FROM warehouse_receptions WHERE id=$1`, id, movement.kind, movement.quantity, fmt.Sprintf("warehouse-inspection:%s:%s", inspectionID, strings.ToLower(movement.bucket)), input.Reason, actor, movement.bucket); err != nil {
			return domain.Inspection{}, err
		}
	}
	if _, err = tx.Exec(ctx, `UPDATE warehouse_receptions SET status='INSPECTED',inspected_at=CURRENT_TIMESTAMP,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, id); err != nil {
		return domain.Inspection{}, err
	}
	if input.AcceptedQuantity > 0 {
		var activatedProductID string
		err = tx.QueryRow(ctx, `UPDATE products p SET status='ACTIVE',published_at=COALESCE(p.published_at,CURRENT_TIMESTAMP),updated_at=CURRENT_TIMESTAMP
			WHERE p.id=(SELECT product_id FROM warehouse_receptions WHERE id=$1)
			AND p.status='APPROVED'
			AND p.price_minor>0
			AND p.currency ~ '^[A-Z]{3}$'
			AND EXISTS (SELECT 1 FROM product_media pm WHERE pm.product_id=p.id)
			AND EXISTS (SELECT 1 FROM artisan_profiles a WHERE a.id=p.artisan_profile_id AND a.status='APPROVED' AND EXISTS (SELECT 1 FROM artisan_memberships m WHERE m.artisan_profile_id=a.id AND m.status='ACTIVE'))
			AND EXISTS (SELECT 1 FROM workshops w WHERE w.id=p.workshop_id AND w.status='ACTIVE')
			AND COALESCE((SELECT SUM(quantity_delta) FROM inventory_movements im WHERE im.product_id=p.id AND im.stock_bucket='AVAILABLE'),0)>0
			RETURNING p.id`, id).Scan(&activatedProductID)
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		} else if err != nil {
			return domain.Inspection{}, err
		} else {
			if _, err = tx.Exec(ctx, `UPDATE product_media SET visibility='PUBLIC' WHERE product_id=$1`, activatedProductID); err != nil {
				return domain.Inspection{}, err
			}
			if _, err = tx.Exec(ctx, `INSERT INTO audit_events(event_type,actor_user_id,target_type,target_id,reason,previous_state,new_state) VALUES('PRODUCT_AUTO_ACTIVATED_AFTER_INSPECTION',$1,'product',$2,'accepted warehouse stock',jsonb_build_object('status','APPROVED'),jsonb_build_object('status','ACTIVE'))`, actor, activatedProductID); err != nil {
				return domain.Inspection{}, err
			}
			if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES('PRODUCT_AUTO_ACTIVATED_AFTER_INSPECTION','product',$1::uuid,jsonb_build_object('productId',$2::text,'status','ACTIVE','source','warehouse_inspection'))`, activatedProductID, activatedProductID); err != nil {
				return domain.Inspection{}, err
			}
		}
	}
	if err = warehouseEvent(ctx, tx, actor, "WAREHOUSE_RECEPTION_INSPECTED", id, "RECEIVED_PENDING_INSPECTION", "INSPECTED", map[string]any{"acceptedQuantity": input.AcceptedQuantity, "rejectedQuantity": input.RejectedQuantity, "quarantinedQuantity": input.QuarantinedQuantity, "damagedQuantity": input.DamagedQuantity}); err != nil {
		return domain.Inspection{}, err
	}
	if input.AcceptedQuantity > 0 {
		if err = warehouseEvent(ctx, tx, actor, "INVENTORY_ACCEPTED", id, "RECEIVED_PENDING_INSPECTION", "INSPECTED", map[string]any{"acceptedQuantity": input.AcceptedQuantity}); err != nil {
			return domain.Inspection{}, err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return domain.Inspection{}, err
	}
	return r.getInspection(ctx, id)
}

func (r *PostgresRepository) AddEvidence(ctx context.Context, actor, receptionID string, input domain.EvidenceInput) (domain.Evidence, error) {
	var item domain.Evidence
	err := r.pool.QueryRow(ctx, `INSERT INTO warehouse_evidence(reception_id,object_key,original_filename,media_type,size_bytes,checksum_sha256,uploaded_by_user_id) SELECT $1,$2,NULLIF($3,''),$4,$5,NULLIF($6,''),$7 FROM warehouse_receptions WHERE id=$1 AND status='RECEIVED_PENDING_INSPECTION' RETURNING id,reception_id,object_key,COALESCE(original_filename,''),media_type,size_bytes,COALESCE(checksum_sha256,''),uploaded_by_user_id,created_at`, receptionID, input.ObjectKey, input.OriginalFilename, input.MediaType, input.SizeBytes, input.Checksum, actor).Scan(&item.ID, &item.ReceptionID, &item.ObjectKey, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.Checksum, &item.UploadedBy, &item.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Evidence{}, domain.ErrInvalidTransition
	}
	if err != nil {
		return domain.Evidence{}, err
	}
	return item, nil
}

func (r *PostgresRepository) ListEvidence(ctx context.Context, receptionID string) ([]domain.Evidence, error) {
	rows, err := r.pool.Query(ctx, `SELECT e.id,e.reception_id,e.object_key,COALESCE(e.original_filename,''),e.media_type,e.size_bytes,COALESCE(e.checksum_sha256,''),e.uploaded_by_user_id,e.created_at FROM warehouse_evidence e WHERE e.reception_id=$1 ORDER BY e.created_at,e.id`, receptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.Evidence{}
	for rows.Next() {
		var item domain.Evidence
		if err = rows.Scan(&item.ID, &item.ReceptionID, &item.ObjectKey, &item.OriginalFilename, &item.MediaType, &item.SizeBytes, &item.Checksum, &item.UploadedBy, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PostgresRepository) getReception(ctx context.Context, id string) (domain.Reception, error) {
	items, _, err := r.ListReceptions(ctx, "", "", 100, 0)
	if err != nil {
		return domain.Reception{}, err
	}
	for _, item := range items {
		if item.ID == id {
			return item, nil
		}
	}
	return domain.Reception{}, domain.ErrNotFound
}

func (r *PostgresRepository) getInspection(ctx context.Context, receptionID string) (domain.Inspection, error) {
	var item domain.Inspection
	err := r.pool.QueryRow(ctx, `SELECT id,reception_id,accepted_quantity,rejected_quantity,quarantined_quantity,damaged_quantity,reason,inspected_by_user_id,inspected_at FROM warehouse_inspections WHERE reception_id=$1`, receptionID).Scan(&item.ID, &item.ReceptionID, &item.AcceptedQuantity, &item.RejectedQuantity, &item.QuarantinedQuantity, &item.DamagedQuantity, &item.Reason, &item.InspectedBy, &item.InspectedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Inspection{}, domain.ErrNotFound
	}
	return item, err
}
