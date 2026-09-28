package application

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth"
	"github.com/aisha-platform/aisha/apps/api/internal/features/notification/domain"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
	ticketTTL       = time.Minute
)

var ErrInvalidTicket = errors.New("invalid notification socket ticket")

type Service struct {
	repository domain.Repository
	tickets    *TicketStore
}

func NewService(repository domain.Repository) *Service {
	return &Service{repository: repository, tickets: NewTicketStore()}
}

func (s *Service) List(ctx context.Context, principal auth.Principal, page, pageSize int) ([]domain.Notification, int, int, error) {
	page, pageSize = normalizePage(page, pageSize)
	items, total, err := s.repository.List(ctx, principal.UserID, pageSize, (page-1)*pageSize)
	return items, total, pageSize, err
}

func (s *Service) UnreadCount(ctx context.Context, principal auth.Principal) (int, error) {
	return s.repository.UnreadCount(ctx, principal.UserID)
}

func (s *Service) MarkRead(ctx context.Context, principal auth.Principal, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.ErrValidation
	}
	return s.repository.MarkRead(ctx, principal.UserID, id)
}

func (s *Service) MarkAllRead(ctx context.Context, principal auth.Principal) error {
	return s.repository.MarkAllRead(ctx, principal.UserID)
}

func (s *Service) IssueSocketTicket(principal auth.Principal) string {
	return s.tickets.Issue(principal.UserID)
}

func (s *Service) ConsumeSocketTicket(ticket string) (string, error) {
	return s.tickets.Consume(ticket)
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	return page, pageSize
}

type TicketStore struct {
	mu      sync.Mutex
	tickets map[string]ticket
}

type ticket struct {
	userID  string
	expires time.Time
}

func NewTicketStore() *TicketStore {
	return &TicketStore{tickets: make(map[string]ticket)}
}

func (s *TicketStore) Issue(userID string) string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	token := base64.RawURLEncoding.EncodeToString(bytes)
	s.mu.Lock()
	s.pruneLocked(time.Now().UTC())
	s.tickets[token] = ticket{userID: userID, expires: time.Now().UTC().Add(ticketTTL)}
	s.mu.Unlock()
	return token
}

func (s *TicketStore) Consume(token string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.tickets[token]
	delete(s.tickets, token)
	if !ok || token == "" || !value.expires.After(time.Now().UTC()) {
		return "", ErrInvalidTicket
	}
	return value.userID, nil
}

func (s *TicketStore) pruneLocked(now time.Time) {
	for token, value := range s.tickets {
		if !value.expires.After(now) {
			delete(s.tickets, token)
		}
	}
}
