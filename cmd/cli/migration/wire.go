//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"doan/internal/infrastructure/database"
	"doan/pkg/config"
	"doan/pkg/logger"
	"github.com/google/wire"
)

func wireApp(app *App) error {
	wire.Build(
		// Provide context
		provideContext,

		// Provide config manager
		provideConfigManager,

		// Provide logger
		provideLogger,

		// Database providers (includes GetDBContext and NewMigration)
		database.DBProvider,

		// Inject into app
		inject,
	)
	return nil
}

// provideContext creates application context with trace ID
func provideContext() context.Context {
	return logger.NewBackgroundContextWithTraceID("migration")
}

// provideConfigManager returns the config manager (viper already loaded)
func provideConfigManager() config.Manager {
	return config.GetManager()
}

// provideLogger creates logger with context
func provideLogger(ctx context.Context) logger.Logger {
	return logger.NewZapLogger(logger.Config{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		ServiceName: "doan-migration",
		Environment: "development",
	})
}
