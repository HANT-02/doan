//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"doan/cmd/http/controllers"
	"doan/internal/infrastructure"
	"doan/internal/services"
	"doan/internal/usecases"
	"doan/pkg/config"
	"doan/pkg/logger"

	"github.com/google/wire"
)

// wireApp wires all dependencies using Clean Architecture layers:
// App Layer -> Controllers -> UseCases -> Services -> Infrastructure
func wireApp(app *App) error {
	wire.Build(
		// Layer 1: Application fundamentals (Context, Config, Logger)
		AppProviders,

		// Layer 2: Infrastructure (Database, Queue, External services)
		infrastructure.InfrastructureProviders,

		// Layer 3: Services (Domain services, Security, Mailer)
		services.ServiceProviders,

		// Layer 4: Use Cases (Business logic)
		usecases.UseCaseProviders,

		// Layer 5: Controllers (HTTP handlers)
		controllers.ControllerProviders,

		// Layer 6: Application injection
		inject,
	)
	return nil
}

// AppProviders provides application-level dependencies
var AppProviders = wire.NewSet(
	provideContext,
	provideConfigManager,
	provideLogger,
)

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
