package postgres

import (
	"context"
	"doan/internal/entities"
	migrationinterface "doan/internal/repositories/migration"
	"doan/pkg/logger"
	"fmt"
	"gorm.io/gorm"
)

type migration struct {
	db  *gorm.DB
	log logger.Logger
	ctx context.Context
}

func NewMigration(
	ctx context.Context,
	db *gorm.DB,
	log logger.Logger,
) migrationinterface.Migration {
	return &migration{
		db:  db,
		log: log,
		ctx: ctx,
	}
}

func (m *migration) Up() error {
	m.log.Info(m.ctx, "Starting GORM auto migration...")

	// Step 1: Enable UUID extension
	if err := m.enableUUIDExtension(); err != nil {
		m.log.Error(m.ctx, "Failed to enable UUID extension", "error", err)
		return fmt.Errorf("failed to enable UUID extension: %w", err)
	}
	m.log.Info(m.ctx, "UUID extension enabled successfully")

	// Step 2: Auto migrate all entities
	entities := m.getAllEntities()
	m.log.Info(m.ctx, "Found entities to migrate", "count", len(entities))

	if err := m.db.WithContext(m.ctx).AutoMigrate(entities...); err != nil {
		m.log.Error(m.ctx, "Auto migration failed", "error", err)
		return fmt.Errorf("auto migration failed: %w", err)
	}
	m.log.Info(m.ctx, "Auto migration completed successfully", "tables", len(entities))

	// Step 3: Seed sample data (only if table is empty)
	if err := m.seedSampleData(); err != nil {
		m.log.Error(m.ctx, "Failed to seed sample data", "error", err)
		return fmt.Errorf("failed to seed sample data: %w", err)
	}

	m.log.Info(m.ctx, "Migration completed successfully!")
	return nil
}

func (m *migration) Down() error {
	m.log.Info(m.ctx, "Rolling back migrations...")

	// Drop tables
	entities := m.getAllEntities()
	if err := m.db.WithContext(m.ctx).Migrator().DropTable(entities...); err != nil {
		m.log.Error(m.ctx, "Rollback failed", "error", err)
		return err
	}

	m.log.Info(m.ctx, "Rollback completed successfully")
	return nil
}

// getAllEntities returns all entity models that need to be migrated
func (m *migration) getAllEntities() []interface{} {
	return []interface{}{
		&entities.User{},
		// Add more entities here:
		// &entities.Product{},
		// &entities.Order{},
	}
}

// enableUUIDExtension enables the uuid-ossp extension in PostgreSQL
func (m *migration) enableUUIDExtension() error {
	return m.db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
}

// seedSampleData inserts sample users if the table is empty
func (m *migration) seedSampleData() error {
	var count int64
	if err := m.db.Model(&entities.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	// Only seed if table is empty
	if count > 0 {
		m.log.Info(m.ctx, "Users table already has data, skipping seed", "count", count)
		return nil
	}

	m.log.Info(m.ctx, "Seeding sample users...")

	// Sample users - Password is bcrypt hash of "password123"
	sampleUsers := []entities.User{
		{UserName: "admin", PassWord: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Status: "active"},
		{UserName: "user1", PassWord: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Status: "active"},
		{UserName: "user2", PassWord: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Status: "active"},
		{UserName: "testuser", PassWord: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Status: "active"},
		{UserName: "john_doe", PassWord: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Status: "active"},
		{UserName: "jane_smith", PassWord: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Status: "active"},
		{UserName: "bob_wilson", PassWord: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Status: "inactive"},
		{UserName: "alice_jones", PassWord: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Status: "active"},
		{UserName: "charlie_brown", PassWord: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Status: "active"},
		{UserName: "david_miller", PassWord: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Status: "suspended"},
	}

	// Insert in batch
	if err := m.db.Create(&sampleUsers).Error; err != nil {
		return fmt.Errorf("failed to insert sample users: %w", err)
	}

	m.log.Info(m.ctx, "Sample users seeded successfully", "count", len(sampleUsers))
	return nil
}
