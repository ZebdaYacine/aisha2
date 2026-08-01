//go:build wireinject

package main

import (
	"context"

	"github.com/aisha-platform/aisha/backend/internal/platform/config"
	"github.com/aisha-platform/aisha/backend/internal/platform/httpapi"
	"github.com/google/wire"
)

func initializeAPI(context.Context, config.Config) (*httpapi.Runtime, error) {
	wire.Build(httpapi.ProviderSet)
	return nil, nil
}
