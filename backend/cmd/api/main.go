package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/aisha-platform/aisha/backend/internal/platform/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	runtime, err := initializeAPI(ctx, cfg)
	if err != nil {
		slog.Error("initialize api", "error", err)
		os.Exit(1)
	}
	defer runtime.Close()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := runtime.App.ShutdownWithContext(shutdownCtx); err != nil {
			slog.Error("shutdown api", "error", err)
		}
	}()
	if err := runtime.App.Listen(":" + cfg.APIPort); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("api stopped", "error", err)
		os.Exit(1)
	}
}
