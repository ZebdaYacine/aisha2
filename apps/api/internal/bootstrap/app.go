package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/container"
)

// ErrServerStopped is returned when the HTTP server exits normally.
var ErrServerStopped = errors.New("server stopped")

// App is the process-level application assembled during bootstrap.
type App struct {
	runtime *container.Runtime
	port    string
}

// NewApp loads configuration and assembles infrastructure, feature modules,
// and the HTTP server. Database migrations remain an explicit deployment step.
func NewApp() (*App, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}
	runtime, err := container.Initialize(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("assemble application: %w", err)
	}
	return &App{runtime: runtime, port: cfg.APIPort}, nil
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- a.runtime.App.Listen(":" + a.port)
	}()

	select {
	case err := <-serverErrors:
		if err != nil {
			return err
		}
		return ErrServerStopped
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := a.runtime.App.ShutdownWithContext(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown server: %w", err)
		}
		return ErrServerStopped
	}
}

func (a *App) Close() {
	if a != nil && a.runtime != nil {
		a.runtime.Close()
	}
}
