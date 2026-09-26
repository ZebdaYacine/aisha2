package artisan

import (
	"context"
	"errors"
	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"testing"
)

type repoStub struct {
	user, decision, reason        string
	called                        bool
	draftCalled, finalizeCalled   bool
	activationKey, activationHash string
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
func (r *repoStub) SaveDraft(_ context.Context, user string, _ ApplicationInput) (Application, error) {
	r.user = user
	r.draftCalled = true
	return Application{Status: "DRAFT"}, nil
}
func (r *repoStub) FinalizeSubmission(_ context.Context, user string) (Application, error) {
	r.user = user
	r.finalizeCalled = true
	return Application{Status: "SUBMITTED"}, nil
}

func (r *repoStub) ActivateMembership(_ context.Context, userID string, _ WorkshopInput, key, requestHash string) (Application, error) {
	r.user, r.activationKey, r.activationHash = userID, key, requestHash
	return Application{MembershipStatus: "ACTIVE"}, nil
}
func (r *repoStub) ListWorkshops(context.Context, string) ([]Workshop, error) { return nil, nil }
func (r *repoStub) CreateWorkshop(context.Context, string, WorkshopInput, string, string) (Workshop, error) {
	return Workshop{}, nil
}
func (r *repoStub) UpdateWorkshop(context.Context, string, string, WorkshopInput) (Workshop, error) {
	return Workshop{}, nil
}
func (r *repoStub) SetWorkshopStatus(context.Context, string, string, string) (Workshop, error) {
	return Workshop{}, nil
}
func (r *repoStub) SetWorkshopStatusAdmin(context.Context, string, string, string, string) (Workshop, error) {
	return Workshop{}, nil
}
func (r *repoStub) DeleteWorkshop(context.Context, string, string) error { return nil }
func (r *repoStub) SetMembershipStatus(context.Context, string, string, string, string) (Application, error) {
	return Application{}, nil
}
func (r *repoStub) MineVerification(context.Context, string) (Verification, error) {
	return Verification{}, nil
}
func (r *repoStub) ListVerifications(context.Context, string, int, int) ([]Verification, int, error) {
	return nil, 0, nil
}
func (r *repoStub) DecideVerification(context.Context, string, VerificationDecision) (Verification, error) {
	return Verification{}, nil
}

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
func TestSaveDraftUsesPrincipalOwnership(t *testing.T) {
	r := &repoStub{}
	s := NewService(r, authzStub{})
	_, err := s.SaveDraft(context.Background(), auth.Principal{UserID: "user-1"}, validInput())
	if err != nil {
		t.Fatal(err)
	}
	if r.user != "user-1" || !r.draftCalled {
		t.Fatalf("user=%q draftCalled=%v", r.user, r.draftCalled)
	}
}
func TestFinalizeSubmissionUsesPrincipalOwnership(t *testing.T) {
	r := &repoStub{}
	s := NewService(r, authzStub{})
	result, err := s.FinalizeSubmission(context.Background(), auth.Principal{UserID: "user-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "SUBMITTED" || r.user != "user-1" || !r.finalizeCalled {
		t.Fatalf("result=%#v user=%q finalizeCalled=%v", result, r.user, r.finalizeCalled)
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

func TestAdminWorkshopStatusRequiresReasonForDeactivation(t *testing.T) {
	r := &repoStub{}
	s := NewService(r, authzStub{})
	_, err := s.SetWorkshopStatusAdmin(context.Background(), auth.Principal{UserID: "admin-1"}, "workshop-1", "INACTIVE", "")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("err=%v", err)
	}
}

func TestMembershipActivationRequiresIdempotencyKey(t *testing.T) {
	r := &repoStub{}
	s := NewService(r, authzStub{})
	_, err := s.ActivateMembership(context.Background(), auth.Principal{UserID: "user-1"}, WorkshopInput{Name: "Atelier", Wilaya: "Tizi Ouzou"}, "")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("err=%v", err)
	}
}

func TestMembershipActivationUsesPrincipalAndRequestHash(t *testing.T) {
	r := &repoStub{}
	s := NewService(r, authzStub{})
	_, err := s.ActivateMembership(context.Background(), auth.Principal{UserID: "user-1"}, WorkshopInput{Name: "Atelier", Wilaya: "Tizi Ouzou"}, "activation-1")
	if err != nil {
		t.Fatal(err)
	}
	if r.user != "user-1" || r.activationKey != "activation-1" || r.activationHash == "" {
		t.Fatalf("user=%q key=%q hash=%q", r.user, r.activationKey, r.activationHash)
	}
}

func TestSuspensionRequiresReason(t *testing.T) {
	r := &repoStub{}
	s := NewService(r, authzStub{})
	_, err := s.SetMembershipStatus(context.Background(), auth.Principal{UserID: "admin-1"}, "membership-1", "SUSPENDED", "")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("err=%v", err)
	}
}

func TestAdministratorCannotEnterArtisanSelfService(t *testing.T) {
	r := &repoStub{}
	s := NewService(r, authzStub{})
	p := auth.Principal{UserID: "admin-1", Roles: []string{"administrator"}}
	if _, err := s.Submit(context.Background(), p, validInput()); !errors.Is(err, ErrValidation) {
		t.Fatalf("Submit() error = %v", err)
	}
	if _, err := s.ActivateMembership(context.Background(), p, WorkshopInput{Name: "Atelier", Wilaya: "Algiers"}, "key"); !errors.Is(err, ErrValidation) {
		t.Fatalf("ActivateMembership() error = %v", err)
	}
	if r.called || r.user != "" {
		t.Fatal("administrator reached artisan repository")
	}
}
