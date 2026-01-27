package main

import (
	"context"
	"doan/internal/entities"
	"doan/internal/infrastructure/database/postgres"
	"doan/pkg/config"
	"doan/pkg/constants"
	"doan/pkg/logger"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type App struct {
	ConfigFilePath string
	ConfigFile     string
	db             *gorm.DB
	log            logger.Logger
}

func (a *App) initFlag() {
	flag.StringVar(&a.ConfigFilePath, "config-file-path", "./configs", "Config file path")
	flag.StringVar(&a.ConfigFile, "config-file", "config", "Config file name")
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

func (a *App) initLogger() {
	a.log = logger.NewZapLogger(logger.Config{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		ServiceName: "doan-seeder",
		Environment: "development",
	})
}

func (a *App) initDB(ctx context.Context) error {
	db, err := postgres.GetDBContext(ctx, a.log, config.GetManager())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	a.db = db
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (a *App) seedUsers(ctx context.Context) error {
	a.log.Info(ctx, "Starting to seed users...")

	// Sample users
	users := []entities.User{
		{
			ID:        uuid.New().String(),
			Code:      "ADMIN001",
			FullName:  "Admin User",
			Email:     "admin@example.com",
			Role:      "ADMIN",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Code:      "TEACHER001",
			FullName:  "John Teacher",
			Email:     "teacher@example.com",
			Role:      "TEACHER",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Code:      "STUDENT001",
			FullName:  "Alice Student",
			Email:     "student1@example.com",
			Role:      "STUDENT",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Code:      "STUDENT002",
			FullName:  "Bob Student",
			Email:     "student2@example.com",
			Role:      "STUDENT",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Code:      "STUDENT003",
			FullName:  "Charlie Student",
			Email:     "student3@example.com",
			Role:      "STUDENT",
			IsActive:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Default password for all users
	defaultPassword := "password123"
	hashedPassword, err := hashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert users
	for _, user := range users {
		user.Password = hashedPassword

		// Check if user already exists
		var existingUser entities.User
		result := a.db.Where("email = ?", user.Email).First(&existingUser)

		if result.Error == gorm.ErrRecordNotFound {
			// User doesn't exist, create it
			if err := a.db.Create(&user).Error; err != nil {
				a.log.Error(ctx, "Failed to create user %s: %v", user.Email, err)
				continue
			}
			a.log.Error(ctx, "✓ Created user: %s (%s) - Role: %s", user.FullName, user.Email, user.Role)
		} else if result.Error != nil {
			a.log.Error(ctx, "Error checking user %s: %v", user.Email, result.Error)
			continue
		} else {
			a.log.Info(ctx, "- User already exists: %s (%s)", user.FullName, user.Email)
		}
	}

	a.log.Info(ctx, "User seeding completed!")
	a.log.Info(ctx, "\n=== Sample Login Credentials ===")
	a.log.Info(ctx, "Default password for all users: password123")
	a.log.Info(ctx, "\nAdmin:")
	a.log.Info(ctx, "  Email: admin@example.com")
	a.log.Info(ctx, "  Role: ADMIN")
	a.log.Info(ctx, "\nTeacher:")
	a.log.Info(ctx, "  Email: teacher@example.com")
	a.log.Info(ctx, "  Role: TEACHER")
	a.log.Info(ctx, "\nStudents:")
	a.log.Info(ctx, "  Email: student1@example.com (Active)")
	a.log.Info(ctx, "  Email: student2@example.com (Active)")
	a.log.Info(ctx, "  Email: student3@example.com (Inactive)")
	a.log.Info(ctx, "================================")

	return nil
}

func main() {
	app := &App{}
	app.initFlag()
	app.initConfig()
	app.initLogger()

	ctx := context.Background()

	// Initialize database connection
	if err := app.initDB(ctx); err != nil {
		app.log.Error(ctx, "Failed to initialize database: %v", err)
		panic(err)
	}

	// Seed users
	if err := app.seedUsers(ctx); err != nil {
		app.log.Error(ctx, "Failed to seed users: %v", err)
		panic(err)
	}

	app.log.Info(ctx, "Seeding completed successfully!")
}
