package repositoryinterface

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
	"time"
)

// TeacherRepository defines the interface for teacher data access
type TeacherRepository interface {
	repositories.BaseRepository[entities.Teacher]

	// Additional teacher-specific methods
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// Timetable: Get lessons for a teacher in date range
	GetTeacherLessons(ctx context.Context, teacherID string, from, to time.Time) ([]entities.Lesson, error)

	// Stats: Calculate teaching hours grouped by period
	GetTeachingHoursStats(ctx context.Context, teacherID string, from, to time.Time, groupBy string) ([]TeachingHoursStat, error)
}

// TeachingHoursStat represents teaching hours statistics for a period
type TeachingHoursStat struct {
	Period string  `json:"period"` // Date/Week/Month depending on groupBy
	Hours  float64 `json:"hours"`  // Total hours taught
}
