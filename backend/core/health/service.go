package health

import (
	"context"
	"fmt"
	"sync"
)

type Checker interface {
	Name() string
	Check(context.Context) error
}

type Result struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services,omitempty"`
}

type Service struct {
	checkers []Checker
}

func New(checkers ...Checker) *Service {
	return &Service{checkers: checkers}
}

func (s *Service) Live() Result {
	return Result{Status: "ok"}
}

func (s *Service) Ready(ctx context.Context) (Result, error) {
	services := make(map[string]string, len(s.checkers))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for _, checker := range s.checkers {
		checker := checker
		wg.Add(1)
		go func() {
			defer wg.Done()
			status := "ok"
			if err := checker.Check(ctx); err != nil {
				status = "unavailable"
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("%s readiness check: %w", checker.Name(), err)
				}
				mu.Unlock()
			}
			mu.Lock()
			services[checker.Name()] = status
			mu.Unlock()
		}()
	}
	wg.Wait()
	if firstErr != nil {
		return Result{Status: "unavailable", Services: services}, firstErr
	}
	return Result{Status: "ok", Services: services}, nil
}
