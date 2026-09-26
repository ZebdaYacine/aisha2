package application

import (
	"context"
	"fmt"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/moderation/domain"
	"github.com/aisha-platform/aisha/apps/api/internal/pkg/storage"
	"strings"
)

type Service struct {
	repository domain.Repository
	authorizer domain.Authorizer
	store      storage.ObjectStore
	bucket     string
}
type QueueItem = domain.QueueItem
type DecisionInput = domain.DecisionInput
type Media = domain.Media
type Authorizer = domain.Authorizer

var (
	ErrValidation         = domain.ErrValidation
	ErrNotFound           = domain.ErrNotFound
	ErrInvalidTransition  = domain.ErrInvalidTransition
	ErrActivationNotReady = domain.ErrActivationNotReady
)

func NewService(r domain.Repository, a domain.Authorizer) *Service {
	return &Service{repository: r, authorizer: a}
}

func NewServiceWithMedia(r domain.Repository, a domain.Authorizer, store storage.ObjectStore, bucket string) *Service {
	return &Service{repository: r, authorizer: a, store: store, bucket: bucket}
}
func (s *Service) ListQueue(ctx context.Context, p auth.Principal, status string, page, size int) ([]domain.QueueItem, int, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/product-submissions", "read"); err != nil {
		return nil, 0, err
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
	status = strings.ToUpper(strings.TrimSpace(status))
	switch status {
	case "", "PENDING_REVIEW", "CHANGES_REQUESTED", "APPROVED", "ACTIVE", "SUSPENDED", "ARCHIVED":
	default:
		return nil, 0, domain.ErrValidation
	}
	items, total, err := s.repository.ListQueue(ctx, status, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		for j := range items[i].Media {
			if items[i].Media[j].ObjectKey == "" || s.store == nil {
				continue
			}
			url, signErr := s.store.PresignedGet(ctx, s.bucket, items[i].Media[j].ObjectKey, 15*time.Minute)
			if signErr != nil {
				return nil, 0, fmt.Errorf("create moderation media URL: %w", signErr)
			}
			items[i].Media[j].URL = url.String()
		}
	}
	return items, total, nil
}
func (s *Service) Decide(ctx context.Context, p auth.Principal, in domain.DecisionInput) (domain.QueueItem, error) {
	if err := s.authorizer.Authorize(ctx, p, "/api/v1/admin/product-submissions", "write"); err != nil {
		return domain.QueueItem{}, err
	}
	in.Action = strings.ToUpper(strings.TrimSpace(in.Action))
	in.Reason = strings.TrimSpace(in.Reason)
	switch in.Action {
	case "APPROVE", "REQUEST_CHANGES", "REJECT", "SUSPEND", "ACTIVATE", "ARCHIVE":
	default:
		return domain.QueueItem{}, domain.ErrValidation
	}
	if in.ID == "" || ((in.Action == "REQUEST_CHANGES" || in.Action == "REJECT" || in.Action == "SUSPEND") && in.Reason == "") {
		return domain.QueueItem{}, domain.ErrValidation
	}
	return s.repository.Decide(ctx, p.UserID, in)
}
