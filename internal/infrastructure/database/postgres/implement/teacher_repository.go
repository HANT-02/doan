package implement

import (
	"context"
	"doan/internal/entities"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/base_struct"
	"doan/pkg/config"
	"doan/pkg/logger"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type teacherRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Teacher]
	db *gorm.DB
}

// NewTeacherRepository creates a new teacher repository instance
func NewTeacherRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.TeacherRepository {
	modelRepo := postgres.NewBaseRepository[entities.Teacher](log, manager, db, "teachers")
	return &teacherRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

// ExistsByEmail checks if a teacher with the given email exists
func (r *teacherRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Teacher{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByCode checks if a teacher with the given code exists
func (r *teacherRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Teacher{}).Where("code = ?", code).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetTeacherLessons retrieves all lessons for a teacher within a date range
func (r *teacherRepository) GetTeacherLessons(ctx context.Context, teacherID string, from, to time.Time) ([]entities.Lesson, error) {
	var lessons []entities.Lesson

	query := r.db.WithContext(ctx).
		Preload("Class").
		Preload("Room").
		Preload("Teacher").
		Where("teacher_id = ?", teacherID)

	// Apply date range filter if provided
	if !from.IsZero() {
		query = query.Where("date_start >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("date_end <= ?", to)
	}

	err := query.Order("date_start ASC").Find(&lessons).Error
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

// GetTeachingHoursStats calculates teaching hours statistics grouped by period
func (r *teacherRepository) GetTeachingHoursStats(ctx context.Context, teacherID string, from, to time.Time, groupBy string) ([]repointerface.TeachingHoursStat, error) {
	var stats []repointerface.TeachingHoursStat

	// Determine the date format based on groupBy
	var dateFormat string
	switch groupBy {
	case "day":
		dateFormat = "YYYY-MM-DD"
	case "week":
		dateFormat = "IYYY-IW" // ISO week format
	case "month":
		dateFormat = "YYYY-MM"
	default:
		dateFormat = "YYYY-MM-DD" // Default to day
	}

	// Build the query to calculate hours per period
	query := r.db.WithContext(ctx).
		Model(&entities.Lesson{}).
		Select(fmt.Sprintf(`
			TO_CHAR(date_start, '%s') as period,
			SUM(EXTRACT(EPOCH FROM (date_end - date_start)) / 3600.0) as hours
		`, dateFormat)).
		Where("teacher_id = ?", teacherID)

	// Apply date range filter
	if !from.IsZero() {
		query = query.Where("date_start >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("date_end <= ?", to)
	}

	// Group by period and order
	query = query.Group(fmt.Sprintf("TO_CHAR(date_start, '%s')", dateFormat)).
		Order("period ASC")

	err := query.Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats, nil
}
