package application

import (
	"context"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/moderation/domain"
	"testing"
)

type moderationAuth struct{ err error }

func (a moderationAuth) Authorize(context.Context, auth.Principal, string, string) error {
	return a.err
}

type moderationRepo struct{ input domain.DecisionInput }

func (r *moderationRepo) ListQueue(context.Context, string, int, int) ([]domain.QueueItem, int, error) {
	return nil, 0, nil
}
func (r *moderationRepo) Decide(_ context.Context, _ string, i domain.DecisionInput) (domain.QueueItem, error) {
	r.input = i
	return domain.QueueItem{ProductID: "p"}, nil
}
func TestListQueueRejectsUnknownStatus(t *testing.T) {
	s := NewService(&moderationRepo{}, moderationAuth{})
	if _, _, err := s.ListQueue(context.Background(), auth.Principal{UserID: "u"}, "UNKNOWN", 1, 10); err != domain.ErrValidation {
		t.Fatalf("err=%v", err)
	}
}

func TestDecisionRequiresReasonForProtectedActions(t *testing.T) {
	r := &moderationRepo{}
	s := NewService(r, moderationAuth{})
	if _, err := s.Decide(context.Background(), auth.Principal{UserID: "u"}, domain.DecisionInput{ID: "x", Action: "REJECT"}); err != domain.ErrValidation {
		t.Fatalf("err=%v", err)
	}
	if _, err := s.Decide(context.Background(), auth.Principal{UserID: "u"}, domain.DecisionInput{ID: "x", Action: "reject", Reason: "incomplete"}); err != nil {
		t.Fatal(err)
	}
	if r.input.Action != "REJECT" {
		t.Fatalf("action=%q", r.input.Action)
	}
}
