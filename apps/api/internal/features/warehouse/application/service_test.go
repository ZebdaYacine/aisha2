package application

import (
	"context"
	"testing"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/warehouse/domain"
)

type warehouseAuthorizer struct{ err error }

func (a warehouseAuthorizer) Authorize(context.Context, auth.Principal, string, string) error {
	return a.err
}

type warehouseRepository struct{}

func (warehouseRepository) CreateReception(context.Context, string, domain.ReceptionInput) (domain.Reception, error) {
	return domain.Reception{}, nil
}
func (warehouseRepository) ListReceptions(context.Context, string, string, int, int) ([]domain.Reception, int, error) {
	return nil, 0, nil
}
func (warehouseRepository) ListValidatedProducts(context.Context, string, string, string, string, int, int) ([]domain.ValidatedProduct, int, error) {
	return nil, 0, nil
}
func (warehouseRepository) Inspect(context.Context, string, string, domain.InspectionInput) (domain.Inspection, error) {
	return domain.Inspection{}, nil
}
func (warehouseRepository) AddEvidence(context.Context, string, string, domain.EvidenceInput) (domain.Evidence, error) {
	return domain.Evidence{}, nil
}
func (warehouseRepository) ListEvidence(context.Context, string) ([]domain.Evidence, error) {
	return nil, nil
}

func TestCreateReceptionRejectsInvalidInput(t *testing.T) {
	s := NewService(warehouseRepository{}, warehouseAuthorizer{}, nil, "", 0)
	if _, err := s.CreateReception(context.Background(), auth.Principal{UserID: "user"}, domain.ReceptionInput{ProductID: "not-a-uuid", ReceivedQuantity: 1, ReferenceKey: "ref"}); err != domain.ErrValidation {
		t.Fatalf("err=%v", err)
	}
}

func TestListReceptionsRejectsUnknownStatus(t *testing.T) {
	s := NewService(warehouseRepository{}, warehouseAuthorizer{}, nil, "", 0)
	if _, _, err := s.ListReceptions(context.Background(), auth.Principal{UserID: "user"}, "UNKNOWN", 1, 20); err != domain.ErrValidation {
		t.Fatalf("err=%v", err)
	}
}

func TestListValidatedProductsRequiresArtisanPhone(t *testing.T) {
	s := NewService(warehouseRepository{}, warehouseAuthorizer{}, nil, "", 0)
	if _, _, err := s.ListValidatedProducts(context.Background(), auth.Principal{UserID: "user"}, "", "", "", 1, 20); err != domain.ErrValidation {
		t.Fatalf("err=%v", err)
	}
}

func TestInspectRejectsEmptyOutcome(t *testing.T) {
	s := NewService(warehouseRepository{}, warehouseAuthorizer{}, nil, "", 0)
	if _, err := s.Inspect(context.Background(), auth.Principal{UserID: "user"}, "00000000-0000-0000-0000-000000000001", domain.InspectionInput{Reason: "checked"}); err != domain.ErrQuantityMismatch {
		t.Fatalf("err=%v", err)
	}
}
