package repositories

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/warehouse/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresReceptionInspectionAndStockBuckets(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	var actorID, artisanUserID, artisanID, workshopID, categoryID, productID, receptionID, artisanPhone string
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM audit_events WHERE target_type='warehouse_reception' AND target_id=$1`, receptionID)
		_, _ = pool.Exec(ctx, `DELETE FROM audit_events WHERE target_type='product' AND target_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM outbox_events WHERE aggregate_type='warehouse_reception' AND aggregate_id=$1`, receptionID)
		_, _ = pool.Exec(ctx, `DELETE FROM outbox_events WHERE aggregate_type='product' AND aggregate_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM inventory_movements WHERE product_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM product_media WHERE product_id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM warehouse_evidence WHERE reception_id=$1`, receptionID)
		_, _ = pool.Exec(ctx, `DELETE FROM warehouse_inspections WHERE reception_id=$1`, receptionID)
		_, _ = pool.Exec(ctx, `DELETE FROM warehouse_receptions WHERE id=$1`, receptionID)
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, productID)
		_, _ = pool.Exec(ctx, `DELETE FROM workshops WHERE id=$1`, workshopID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_memberships WHERE artisan_profile_id=$1`, artisanID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_profiles WHERE id=$1`, artisanID)
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE id=$1`, categoryID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1,$2)`, actorID, artisanUserID)
	}()

	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES ('warehouse-repository-actor-'||gen_random_uuid()::text,'Warehouse Agent') RETURNING id`).Scan(&actorID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES ('warehouse-repository-artisan-'||gen_random_uuid()::text,'Warehouse Artisan') RETURNING id`).Scan(&artisanUserID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `UPDATE users SET phone='055'||lpad((floor(random()*10000000))::bigint::text,7,'0') WHERE id=$1 RETURNING phone`, artisanUserID).Scan(&artisanPhone); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO artisan_profiles(user_id,public_display_name,status) VALUES ($1,'Warehouse Artisan','APPROVED') RETURNING id`, artisanUserID).Scan(&artisanID); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `INSERT INTO artisan_memberships(user_id,artisan_profile_id,status,activated_at) VALUES ($1,$2,'ACTIVE',CURRENT_TIMESTAMP)`, artisanUserID, artisanID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO workshops(artisan_profile_id,name,status,is_default,is_public) VALUES ($1,'Warehouse Workshop','ACTIVE',true,true) RETURNING id`, artisanID).Scan(&workshopID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO categories(slug,display_name) VALUES ('warehouse-repository-'||gen_random_uuid()::text,'Warehouse Category') RETURNING id`).Scan(&categoryID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO products(artisan_profile_id,category_id,workshop_id,product_type,product_code,status,price_minor,currency) VALUES ($1,$2,$3,'ARTISAN_SPECIFIC','AISHA-'||upper(substr(replace(gen_random_uuid()::text,'-',''),1,10)),'APPROVED',2500,'EUR') RETURNING id`, artisanID, categoryID, workshopID).Scan(&productID); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `INSERT INTO product_media(product_id,media_kind,object_key,original_filename,media_type,size_bytes,checksum_sha256,visibility) VALUES ($1,'IMAGE','warehouse-repository/'||$2||'/product.jpg','product.jpg','image/jpeg',100,'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa','PRIVATE')`, productID, productID); err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresRepository(pool)
	validated, total, err := repository.ListValidatedProducts(ctx, actorID, "+213 "+artisanPhone[1:4]+" "+artisanPhone[4:7]+" "+artisanPhone[7:], "", "", 10, 0)
	if err != nil || total != 1 || len(validated) != 1 || validated[0].WorkshopID != workshopID {
		t.Fatalf("validated products=%#v total=%d err=%v", validated, total, err)
	}
	filtered, total, err := repository.ListValidatedProducts(ctx, actorID, artisanPhone, workshopID, "", 10, 0)
	if err != nil || total != 1 || len(filtered) != 1 || filtered[0].ProductID != productID {
		t.Fatalf("filtered products=%#v total=%d err=%v", filtered, total, err)
	}

	reception, err := repository.CreateReception(ctx, actorID, domain.ReceptionInput{ProductID: productID, ReceivedQuantity: 5, ReferenceKey: "warehouse-repository:" + productID, SupplierName: "Warehouse Artisan"})
	if err != nil {
		t.Fatal(err)
	}
	receptionID = reception.ID
	if reception.Status != "RECEIVED_PENDING_INSPECTION" || reception.ReceivedQuantity != 5 {
		t.Fatalf("reception=%#v", reception)
	}
	duplicate, err := repository.CreateReception(ctx, actorID, domain.ReceptionInput{ProductID: productID, ReceivedQuantity: 5, ReferenceKey: reception.ReferenceKey})
	if err != nil || duplicate.ID != receptionID {
		t.Fatalf("duplicate reception=%#v err=%v", duplicate, err)
	}
	if _, _, err = repository.ListReceptions(ctx, "", "RECEIVED_PENDING_INSPECTION", 10, 0); err != nil {
		t.Fatal(err)
	}
	if _, err = repository.Inspect(ctx, actorID, receptionID, domain.InspectionInput{AcceptedQuantity: 5, Reason: "checked"}); !errors.Is(err, domain.ErrEvidenceRequired) {
		t.Fatalf("inspection without evidence error=%v", err)
	}
	if _, err = repository.AddEvidence(ctx, actorID, receptionID, domain.EvidenceInput{ObjectKey: "warehouse-repository/" + productID + "/evidence", OriginalFilename: "inspection.pdf", MediaType: "application/pdf", SizeBytes: 100, Checksum: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}); err != nil {
		t.Fatal(err)
	}
	if _, err = repository.Inspect(ctx, actorID, receptionID, domain.InspectionInput{AcceptedQuantity: 4, Reason: "checked"}); !errors.Is(err, domain.ErrQuantityMismatch) {
		t.Fatalf("mismatched inspection error=%v", err)
	}
	inspection, err := repository.Inspect(ctx, actorID, receptionID, domain.InspectionInput{AcceptedQuantity: 3, RejectedQuantity: 1, QuarantinedQuantity: 1, Reason: "checked"})
	if err != nil {
		t.Fatal(err)
	}
	if inspection.AcceptedQuantity != 3 || inspection.RejectedQuantity != 1 || inspection.QuarantinedQuantity != 1 {
		t.Fatalf("inspection=%#v", inspection)
	}
	again, err := repository.Inspect(ctx, actorID, receptionID, domain.InspectionInput{AcceptedQuantity: 3, RejectedQuantity: 1, QuarantinedQuantity: 1, Reason: "different retry"})
	if err != nil || again.ID != inspection.ID {
		t.Fatalf("idempotent inspection=%#v err=%v", again, err)
	}
	items, total, err := repository.ListReceptions(ctx, "", "INSPECTED", 10, 0)
	if err != nil || total < 1 {
		t.Fatalf("inspected list total=%d err=%v", total, err)
	}
	var found bool
	for _, item := range items {
		if item.ID == receptionID {
			found = item.Inspection != nil && len(item.Evidence) == 1
		}
	}
	if !found {
		t.Fatalf("reception not found in inspected list: %#v", items)
	}
	var available, rejected, quarantined int64
	if err = pool.QueryRow(ctx, `SELECT COALESCE(SUM(quantity_delta) FILTER (WHERE stock_bucket='AVAILABLE'),0),COALESCE(SUM(quantity_delta) FILTER (WHERE stock_bucket='REJECTED'),0),COALESCE(SUM(quantity_delta) FILTER (WHERE stock_bucket='QUARANTINED'),0) FROM inventory_movements WHERE product_id=$1`, productID).Scan(&available, &rejected, &quarantined); err != nil {
		t.Fatal(err)
	}
	if available != 3 || rejected != 1 || quarantined != 1 {
		t.Fatalf("stock buckets available=%d rejected=%d quarantined=%d", available, rejected, quarantined)
	}
	var productStatus, mediaVisibility string
	if err = pool.QueryRow(ctx, `SELECT p.status,pm.visibility FROM products p JOIN product_media pm ON pm.product_id=p.id WHERE p.id=$1`, productID).Scan(&productStatus, &mediaVisibility); err != nil {
		t.Fatal(err)
	}
	if productStatus != "ACTIVE" || mediaVisibility != "PUBLIC" {
		t.Fatalf("expected accepted inspection to publish product, status=%s visibility=%s", productStatus, mediaVisibility)
	}
}
