package application

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/warehouse/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"github.com/google/uuid"
)

const warehouseResource = "/api/v1/warehouse/receptions"

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
	store      storage.ObjectStore
	bucket     string
	maxUpload  int64
}

type Reception = domain.Reception
type Inspection = domain.Inspection
type Evidence = domain.Evidence
type ReceptionInput = domain.ReceptionInput
type InspectionInput = domain.InspectionInput
type Authorizer = domain.Authorizer

var (
	ErrValidation        = domain.ErrValidation
	ErrNotFound          = domain.ErrNotFound
	ErrDuplicate         = domain.ErrDuplicate
	ErrInvalidTransition = domain.ErrInvalidTransition
	ErrQuantityMismatch  = domain.ErrQuantityMismatch
	ErrEvidenceRequired  = domain.ErrEvidenceRequired
	ErrAlreadyInspected  = domain.ErrAlreadyInspected
)

func NewService(repository domain.Repository, authorizer domain.Authorizer, store storage.ObjectStore, bucket string, maxUpload int64) *Service {
	return &Service{repository: repository, authorizer: authorizer, store: store, bucket: bucket, maxUpload: maxUpload}
}

func (s *Service) CreateReception(ctx context.Context, principal auth.Principal, input domain.ReceptionInput) (domain.Reception, error) {
	if err := s.authorizer.Authorize(ctx, principal, warehouseResource, "write"); err != nil {
		return domain.Reception{}, err
	}
	input.ReferenceKey = strings.TrimSpace(input.ReferenceKey)
	input.SupplierName = strings.TrimSpace(input.SupplierName)
	input.ParcelReference = strings.TrimSpace(input.ParcelReference)
	input.Notes = strings.TrimSpace(input.Notes)
	if _, err := uuid.Parse(input.ProductID); err != nil || input.ReceivedQuantity <= 0 || input.ReferenceKey == "" || len(input.ReferenceKey) > 160 || len(input.SupplierName) > 200 || len(input.ParcelReference) > 200 || len(input.Notes) > 4000 {
		return domain.Reception{}, domain.ErrValidation
	}
	return s.repository.CreateReception(ctx, principal.UserID, input)
}

func (s *Service) ListReceptions(ctx context.Context, principal auth.Principal, status string, page, size int) ([]domain.Reception, int, error) {
	if err := s.authorizer.Authorize(ctx, principal, warehouseResource, "read"); err != nil {
		return nil, 0, err
	}
	status = strings.ToUpper(strings.TrimSpace(status))
	switch status {
	case "", "RECEIVED_PENDING_INSPECTION", "INSPECTED":
	default:
		return nil, 0, domain.ErrValidation
	}
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	items, total, err := s.repository.ListReceptions(ctx, principal.UserID, status, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].Evidence, err = s.evidenceURLs(ctx, items[i].Evidence)
		if err != nil {
			return nil, 0, err
		}
	}
	return items, total, nil
}

func (s *Service) Inspect(ctx context.Context, principal auth.Principal, id string, input domain.InspectionInput) (domain.Inspection, error) {
	if err := s.authorizer.Authorize(ctx, principal, warehouseResource, "write"); err != nil {
		return domain.Inspection{}, err
	}
	input.Reason = strings.TrimSpace(input.Reason)
	if _, err := uuid.Parse(id); err != nil || input.AcceptedQuantity < 0 || input.RejectedQuantity < 0 || input.QuarantinedQuantity < 0 || input.DamagedQuantity < 0 || input.Reason == "" || len(input.Reason) > 2000 {
		return domain.Inspection{}, domain.ErrValidation
	}
	if input.AcceptedQuantity+input.RejectedQuantity+input.QuarantinedQuantity+input.DamagedQuantity <= 0 {
		return domain.Inspection{}, domain.ErrQuantityMismatch
	}
	return s.repository.Inspect(ctx, principal.UserID, id, input)
}

func (s *Service) UploadEvidence(ctx context.Context, principal auth.Principal, receptionID, filename, declaredType string, data []byte) (domain.Evidence, error) {
	if err := s.authorizer.Authorize(ctx, principal, warehouseResource, "write"); err != nil {
		return domain.Evidence{}, err
	}
	if s.store == nil {
		return domain.Evidence{}, errors.New("warehouse evidence storage is unavailable")
	}
	if _, err := uuid.Parse(receptionID); err != nil {
		return domain.Evidence{}, domain.ErrValidation
	}
	if s.maxUpload <= 0 {
		s.maxUpload = 50 * 1024 * 1024
	}
	file, err := storage.ValidateFile(data, declaredType, s.maxUpload, map[string]bool{"IMAGE": true, "DOCUMENT": true, "VIDEO": true})
	if err != nil {
		return domain.Evidence{}, err
	}
	key := fmt.Sprintf("warehouse/receptions/%s/%s", receptionID, uuid.NewString())
	if err = s.store.Put(ctx, s.bucket, key, file.MediaType, bytes.NewReader(file.Bytes), int64(len(file.Bytes))); err != nil {
		return domain.Evidence{}, fmt.Errorf("store warehouse evidence: %w", err)
	}
	item, err := s.repository.AddEvidence(ctx, principal.UserID, receptionID, domain.EvidenceInput{ObjectKey: key, OriginalFilename: strings.TrimSpace(filename), MediaType: file.MediaType, SizeBytes: int64(len(file.Bytes)), Checksum: file.Checksum})
	if err != nil {
		_ = s.store.Delete(ctx, s.bucket, key)
		return domain.Evidence{}, err
	}
	return s.evidenceURL(ctx, item)
}

func (s *Service) Evidence(ctx context.Context, principal auth.Principal, receptionID string) ([]domain.Evidence, error) {
	if err := s.authorizer.Authorize(ctx, principal, warehouseResource, "read"); err != nil {
		return nil, err
	}
	if _, err := uuid.Parse(receptionID); err != nil {
		return nil, domain.ErrValidation
	}
	items, err := s.repository.ListEvidence(ctx, receptionID)
	if err != nil {
		return nil, err
	}
	return s.evidenceURLs(ctx, items)
}

func (s *Service) evidenceURLs(ctx context.Context, items []domain.Evidence) ([]domain.Evidence, error) {
	for i := range items {
		var err error
		items[i], err = s.evidenceURL(ctx, items[i])
		if err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (s *Service) evidenceURL(ctx context.Context, item domain.Evidence) (domain.Evidence, error) {
	if s.store == nil {
		return item, nil
	}
	value, err := s.store.PresignedGet(ctx, s.bucket, item.ObjectKey, 15*time.Minute)
	if err != nil {
		return item, err
	}
	item.URL = value.String()
	return item, nil
}
