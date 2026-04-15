package repositoryinterface

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
)

type AttendanceRepository interface {
	repositories.BaseRepository[entities.Attendance]
	ListByLessonID(ctx context.Context, lessonID string) ([]entities.Attendance, error)
	GetByLessonAndStudent(ctx context.Context, lessonID, studentID string) (*entities.Attendance, error)
}
