package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/aisha-platform/aisha/apps/api/internal/bootstrap"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	app, err := bootstrap.NewApp()
	if err != nil {
		slog.Error("initialize api", "error", err)
		os.Exit(1)
	}
	defer app.Close()
	if err := app.Run(); err != nil && !errors.Is(err, bootstrap.ErrServerStopped) {
		slog.Error("api stopped", "error", err)
		os.Exit(1)
	}
}
