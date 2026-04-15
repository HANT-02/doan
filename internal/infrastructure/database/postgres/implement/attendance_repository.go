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

	"gorm.io/gorm"
)

type attendanceRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Attendance]
	db *gorm.DB
}

func NewAttendanceRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.AttendanceRepository {
	modelRepo := postgres.NewBaseRepository[entities.Attendance](log, manager, db, "attendances")
	return &attendanceRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *attendanceRepository) ListByLessonID(ctx context.Context, lessonID string) ([]entities.Attendance, error) {
	records := make([]entities.Attendance, 0)
	err := r.db.WithContext(ctx).
		Preload("Student").
		Where("lesson_id = ?", lessonID).
		Order("created_at ASC").
		Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *attendanceRepository) GetByLessonAndStudent(ctx context.Context, lessonID, studentID string) (*entities.Attendance, error) {
	var record entities.Attendance
	err := r.db.WithContext(ctx).
		Preload("Student").
		Where("lesson_id = ? AND student_id = ?", lessonID, studentID).
		First(&record).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}
