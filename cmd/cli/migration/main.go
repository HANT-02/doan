package main

import (
	"doan/internal/repositories/migration"
	"doan/pkg/config"
	"doan/pkg/constants"
	"doan/pkg/logger"
	"flag"
	"fmt"
	"os"
)

type App struct {
	Name           string
	Version        string
	ConfigFilePath string
	ConfigFile     string
	migration      migration.Migration
	MigrateVersion uint
}

func (a *App) initFlag() {
	flag.StringVar(&a.Name, "name", "service-name", "")
	flag.StringVar(&a.Version, "version", "1.0.0", "")
	flag.StringVar(&a.ConfigFilePath, "config-file-path", "./configs", "Config file path: path to config dir")
	flag.StringVar(&a.ConfigFile, "config-file", "config", "Config file path: path to config dir")
	flag.Parse()
}

func (a *App) initConfig() {
	configSource := &config.Viper{
		ConfigType: constants.ConfigTypeFile,
		FilePath:   a.ConfigFilePath,
		ConfigFile: a.ConfigFile,
	}
	err := configSource.InitConfig()
	if err != nil {
		panic(err)
	}
}

func inject(
	app *App,
	migration migration.Migration,
) error {
	app.migration = migration
	return nil
}

func (a *App) Run() error {
	ctx := logger.NewBackgroundContextWithTraceID("migration")
	ctxLogger := logger.NewLogger(ctx)
	ctxLogger.Info(ctx, "Migration start")

	err := a.migration.Up()
	if err != nil {
		ctxLogger.Error(ctx, "Migration failed", "error", err)
		return err
	}

	ctxLogger.Info(ctx, "Migration success")
	return nil
}

func main() {
	app := &App{}
	app.initFlag()
	app.initConfig()

	// Wire dependencies
	err := wireApp(app)
	if err != nil {
		fmt.Printf("Failed to wire dependencies: %v\n", err)
		os.Exit(1)
	}

	// Run migration
	if err := app.Run(); err != nil {
		fmt.Printf("Migration failed: %v\n", err)
		os.Exit(1)
	}
}
