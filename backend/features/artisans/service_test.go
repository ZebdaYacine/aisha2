package artisan

import (
	"context"
	"errors"
	"github.com/aisha-platform/aisha/backend/features/auth"
	"testing"
)

type repoStub struct {
	user, decision, reason string
	called                 bool
}

func (r *repoStub) Submit(_ context.Context, user string, _ ApplicationInput) (Application, error) {
	r.user = user
	r.called = true
	return Application{}, nil
}
func (r *repoStub) Mine(_ context.Context, user string) (Application, error) {
	r.user = user
	return Application{}, nil
}
func (r *repoStub) UpdateApproved(_ context.Context, user string, _ ApplicationInput) (Application, error) {
	r.user = user
	return Application{}, nil
}
func (r *repoStub) List(context.Context, string, int, int) ([]Application, int, error) {
	return nil, 0, nil
}
func (r *repoStub) Decide(_ context.Context, _ string, _ string, decision, reason string) (Application, error) {
	r.decision = decision
	r.reason = reason
	return Application{}, nil
}
func (r *repoStub) Documents(context.Context, string) ([]Document, error) { return nil, nil }

type authzStub struct{ err error }

func (a authzStub) Authorize(context.Context, auth.Principal, string, string) error { return a.err }
func validInput() ApplicationInput {
	return ApplicationInput{PublicDisplayName: "Atelier Tala", Wilaya: "Tizi Ouzou", ContactVisibility: "PRIVATE", CategoryIDs: []string{"category-id"}, Translations: []Translation{{Locale: "en", Biography: "Story"}}}
}
func TestSubmitUsesPrincipalOwnership(t *testing.T) {
	r := &repoStub{}
	s := NewService(r, authzStub{})
	_, err := s.Submit(context.Background(), auth.Principal{UserID: "user-1"}, validInput())
	if err != nil {
		t.Fatal(err)
	}
	if r.user != "user-1" {
		t.Fatalf("user=%q", r.user)
	}
}
func TestSubmitStopsWhenForbidden(t *testing.T) {
	r := &repoStub{}
	denied := errors.New("denied")
	s := NewService(r, authzStub{denied})
	_, err := s.Submit(context.Background(), auth.Principal{UserID: "user-1"}, validInput())
	if !errors.Is(err, denied) || r.called {
		t.Fatalf("err=%v called=%v", err, r.called)
	}
}
func TestDecisionReasonRequired(t *testing.T) {
	r := &repoStub{}
	s := NewService(r, authzStub{})
	_, err := s.Decide(context.Background(), auth.Principal{}, "id", "REJECTED", "")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("err=%v", err)
	}
}
func TestApproveDoesNotInventReason(t *testing.T) {
	r := &repoStub{}
	s := NewService(r, authzStub{})
	_, err := s.Decide(context.Background(), auth.Principal{}, "id", "APPROVED", "")
	if err != nil {
		t.Fatal(err)
	}
	if r.decision != "APPROVED" || r.reason != "" {
		t.Fatalf("decision=%q reason=%q", r.decision, r.reason)
	}
}
