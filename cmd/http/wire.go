//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"doan/cmd/http/controllers"
	"doan/internal/infrastructure/database"
	"doan/internal/services"
	"doan/internal/usecases"
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

		// Database providers (includes GetDBContext + NewUserRepository)
		database.DBProvider,

		// Service providers
		services.UserServiceProvider,

		// UseCase providers
		usecases.UserUseCaseProviders,

		// Controller providers
		controllers.ControllerProviders,

		// Inject
		inject,
	)
	return nil
}

// provideContext creates application context
func provideContext() context.Context {
	return context.Background()
}

// provideConfigManager returns the config manager
func provideConfigManager() config.Manager {
	return config.GetManager()
}

// provideLogger creates logger
func provideLogger(ctx context.Context) logger.Logger {
	return logger.NewZapLogger(logger.Config{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		ServiceName: "doan-http",
		Environment: "development",
	})
}
