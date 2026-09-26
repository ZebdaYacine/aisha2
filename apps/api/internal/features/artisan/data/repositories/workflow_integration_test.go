package repositories

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/artisan/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresMembershipAndWorkshopWorkflow(t *testing.T) {
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

	var userID, otherUserID, profileID, membershipID, additionalWorkshopID string
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES('artisan-workflow-'||gen_random_uuid()::text,'Workflow Artisan') RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO artisan_profiles(user_id,public_display_name,wilaya,contact_visibility,status) VALUES($1,'Workflow Artisan','Tizi Ouzou','PRIVATE','APPROVED') RETURNING id`, userID).Scan(&profileID); err != nil {
		t.Fatal(err)
	}
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES('artisan-workflow-other-'||gen_random_uuid()::text,'Other Workflow Artisan') RETURNING id`).Scan(&otherUserID); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_verifications WHERE artisan_membership_id=$1`, membershipID)
		_, _ = pool.Exec(ctx, `DELETE FROM idempotency_keys WHERE actor_user_id=$1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM audit_events WHERE actor_user_id=$1 OR target_id IN ($1::uuid,$2::uuid,$3::uuid)`, userID, profileID, additionalWorkshopID)
		_, _ = pool.Exec(ctx, `DELETE FROM outbox_events WHERE aggregate_id IN ($1::uuid,$2::uuid,$3::uuid)`, profileID, membershipID, additionalWorkshopID)
		_, _ = pool.Exec(ctx, `DELETE FROM workshops WHERE artisan_profile_id=$1`, profileID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_memberships WHERE id=$1`, membershipID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_profiles WHERE id=$1`, profileID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id IN ($1,$2)`, userID, otherUserID)
	}()

	repository := NewPostgresRepository(pool)
	input := domain.WorkshopInput{Name: "First Workshop", Wilaya: "Tizi Ouzou", IsPublic: true}
	results := make(chan struct {
		application domain.Application
		err         error
	}, 2)
	var group sync.WaitGroup
	for range 2 {
		group.Add(1)
		go func() {
			defer group.Done()
			application, err := repository.ActivateMembership(ctx, userID, input, "workflow-activation-1", "1111111111111111111111111111111111111111111111111111111111111111")
			results <- struct {
				application domain.Application
				err         error
			}{application: application, err: err}
		}()
	}
	group.Wait()
	close(results)
	var application domain.Application
	for result := range results {
		if result.err != nil {
			t.Fatal(result.err)
		}
		application = result.application
	}
	if application.MembershipStatus != "ACTIVE" {
		t.Fatalf("membership status=%q", application.MembershipStatus)
	}
	if err = pool.QueryRow(ctx, `SELECT id FROM artisan_memberships WHERE user_id=$1`, userID).Scan(&membershipID); err != nil {
		t.Fatal(err)
	}
	workshops, err := repository.ListWorkshops(ctx, userID)
	if err != nil || len(workshops) != 1 || !workshops[0].IsDefault {
		t.Fatalf("first workshops=%#v err=%v", workshops, err)
	}
	replayed, err := repository.ActivateMembership(ctx, userID, input, "workflow-activation-1", "1111111111111111111111111111111111111111111111111111111111111111")
	if err != nil || replayed.MembershipStatus != "ACTIVE" {
		t.Fatalf("replay=%#v err=%v", replayed, err)
	}
	workshops, err = repository.ListWorkshops(ctx, userID)
	if err != nil || len(workshops) != 1 {
		t.Fatalf("replayed workshops=%#v err=%v", workshops, err)
	}
	if _, err = repository.ActivateMembership(ctx, userID, input, "workflow-activation-1", "different-hash"); !errors.Is(err, domain.ErrInvalidTransition) {
		t.Fatalf("different payload error=%v", err)
	}

	additional, err := repository.CreateWorkshop(ctx, userID, domain.WorkshopInput{Name: "Second Workshop", Wilaya: "Blida", IsPublic: false}, "workflow-workshop-1", "2222222222222222222222222222222222222222222222222222222222222222")
	if err != nil {
		t.Fatal(err)
	}
	additionalWorkshopID = additional.ID
	if _, err = repository.CreateWorkshop(ctx, userID, domain.WorkshopInput{Name: "Second Workshop", Wilaya: "Blida", IsPublic: false}, "workflow-workshop-1", "2222222222222222222222222222222222222222222222222222222222222222"); err != nil {
		t.Fatal(err)
	}
	if _, err = repository.SetWorkshopStatus(ctx, userID, additionalWorkshopID, "INACTIVE"); err != nil {
		t.Fatal(err)
	}
	if _, err = repository.ListWorkshops(ctx, otherUserID); err != nil {
		t.Fatal(err)
	}
	if _, err = repository.UpdateWorkshop(ctx, otherUserID, additionalWorkshopID, input); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("cross-owner update error=%v", err)
	}
	if _, err = repository.SetMembershipStatus(ctx, userID, membershipID, "SUSPENDED", "policy review"); err != nil {
		t.Fatal(err)
	}
	if _, err = repository.CreateWorkshop(ctx, userID, domain.WorkshopInput{Name: "Blocked Workshop", Wilaya: "Blida"}, "workflow-workshop-blocked", "3333333333333333333333333333333333333333333333333333333333333333"); !errors.Is(err, domain.ErrInvalidTransition) {
		t.Fatalf("suspended create error=%v", err)
	}
	if _, err = repository.SetMembershipStatus(ctx, userID, membershipID, "ACTIVE", ""); err != nil {
		t.Fatal(err)
	}
	if err = repository.DeleteWorkshop(ctx, userID, additionalWorkshopID); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresArtisanApplicationRequiresFilesBeforeSubmit(t *testing.T) {
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

	var userID, categoryID, profileID string
	if err = pool.QueryRow(ctx, `INSERT INTO users(email,display_name) VALUES('artisan-draft-'||gen_random_uuid()::text,'Draft Artisan') RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_documents WHERE artisan_profile_id=$1`, profileID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_media WHERE artisan_profile_id=$1`, profileID)
		_, _ = pool.Exec(ctx, `DELETE FROM audit_events WHERE actor_user_id=$1 OR target_id=$1::uuid`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM outbox_events WHERE aggregate_id=$1::uuid`, profileID)
		_, _ = pool.Exec(ctx, `DELETE FROM artisan_profiles WHERE id=$1`, profileID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, userID)
	}()
	if err = pool.QueryRow(ctx, `SELECT id FROM categories WHERE is_active=true ORDER BY id LIMIT 1`).Scan(&categoryID); err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresRepository(pool)
	input := domain.ApplicationInput{
		PublicDisplayName: "Draft Artisan",
		WorkshopName:      "Atelier Draft",
		Wilaya:            "Tizi Ouzou",
		ContactVisibility: "PRIVATE",
		CategoryIDs:       []string{categoryID},
		Translations:      []domain.Translation{{Locale: "en", Biography: "A draft"}},
	}
	application, err := repository.SaveDraft(ctx, userID, input)
	if err != nil {
		t.Fatal(err)
	}
	profileID = application.ID
	if _, err = repository.FinalizeSubmission(ctx, userID); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("submit without files error=%v", err)
	}
	if _, err = repository.AddDocument(ctx, userID, domain.DocumentUploadInput{DocumentType: "IDENTITY", ObjectKey: "tests/" + userID + "/identity.pdf", OriginalFilename: "identity.pdf", MediaType: "application/pdf", SizeBytes: 10}); err != nil {
		t.Fatal(err)
	}
	if _, err = repository.AddMedia(ctx, userID, domain.MediaUploadInput{MediaKind: "IMAGE", ObjectKey: "tests/" + userID + "/profile.jpg", OriginalFilename: "profile.jpg", MediaType: "image/jpeg", SizeBytes: 10}); err != nil {
		t.Fatal(err)
	}
	application, err = repository.FinalizeSubmission(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	if application.Status != "SUBMITTED" {
		t.Fatalf("status=%q", application.Status)
	}
}
