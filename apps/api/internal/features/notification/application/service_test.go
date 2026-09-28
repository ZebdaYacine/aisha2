package application

import (
	"testing"
	"time"
)

func TestTicketStoreConsumesTicketsOnce(t *testing.T) {
	store := NewTicketStore()
	ticket := store.Issue("user-1")
	if ticket == "" {
		t.Fatal("expected a ticket")
	}
	userID, err := store.Consume(ticket)
	if err != nil {
		t.Fatalf("consume ticket: %v", err)
	}
	if userID != "user-1" {
		t.Fatalf("user id = %q, want user-1", userID)
	}
	if _, err := store.Consume(ticket); err == nil {
		t.Fatal("expected a consumed ticket to be rejected")
	}
}

func TestTicketStoreRejectsUnknownAndExpiredTickets(t *testing.T) {
	store := NewTicketStore()
	if _, err := store.Consume("unknown"); err == nil {
		t.Fatal("expected unknown ticket to be rejected")
	}
	store.tickets["expired"] = ticket{userID: "user-1", expires: time.Now().UTC().Add(-time.Second)}
	if _, err := store.Consume("expired"); err == nil {
		t.Fatal("expected expired ticket to be rejected")
	}
}
