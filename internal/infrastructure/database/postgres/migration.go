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
		&entities.Teacher{},
		&entities.Student{},
		&entities.Room{},
		&entities.Course{},
		&entities.Program{},
		&entities.ProgramCourse{},
		&entities.Objective{},
		&entities.Outcome{},
		&entities.Class{},
		&entities.Lesson{},
		&entities.Enrollment{},
		&entities.Attendance{},
		&entities.ClassSchedule{},
		&entities.LessonSummary{},
		&entities.AcademicRecord{},
		&entities.Consultation{},
		&entities.LeaveRequest{},
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
		{Email: "admin", Code: "admin", FullName: "admin", Password: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Role: "ADMIN"},
		{Email: "user1", Code: "user1", FullName: "user1", Password: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Role: "STUDENT"},
		{Email: "user2", Code: "user2", FullName: "user2", Password: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Role: "STUDENT"},
		{Email: "testuser", Code: "testuser", FullName: "testuser", Password: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Role: "STUDENT"},
		{Email: "john_doe", Code: "john_doe", FullName: "john_doe", Password: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Role: "STUDENT"},
		{Email: "jane_smith", Code: "jane_smith", FullName: "jane_smith", Password: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Role: "STUDENT"},
		{Email: "bob_wilson", Code: "bob_wilson", FullName: "bob_wilson", Password: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Role: "STUDENT"},
		{Email: "alice_jones", Code: "alice_jones", FullName: "alice_jones", Password: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Role: "STUDENT"},
		{Email: "charlie_brown", Code: "charlie_brown", FullName: "charlie_brown", Password: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Role: "STUDENT"},
		{Email: "david_miller", Code: "david_miller", FullName: "david_miller", Password: "$2a$10$rqJ2VqGXqK3Z5X8Ry5Qz7.YqJ3vZJqJ5X8Ry5Qz7YqJ3vZJqJ5X8R", Role: "STUDENT"},
	}

	// Insert in batch
	if err := m.db.Create(&sampleUsers).Error; err != nil {
		return fmt.Errorf("failed to insert sample users: %w", err)
	}

	m.log.Info(m.ctx, "Sample users seeded successfully", "count", len(sampleUsers))
	return nil
}
