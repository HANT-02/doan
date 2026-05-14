package implement

import (
	"context"
	"strings"
	"time"

	"doan/internal/entities"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/pkg/base_struct"
	"doan/pkg/config"
	"doan/pkg/logger"

	"gorm.io/gorm"
)

type lessonRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Lesson]
	db *gorm.DB
}

func NewLessonRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repositoryinterface.LessonRepository {
	modelRepo := postgres.NewBaseRepository[entities.Lesson](log, manager, db, "lessons")
	return &lessonRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *lessonRepository) FindOverlappingLessons(
	ctx context.Context,
	from time.Time,
	to time.Time,
	classIDs []string,
	teacherIDs []string,
	roomIDs []string,
	statuses []string,
) ([]entities.Lesson, error) {
	if to.Before(from) {
		to = from
	}

	filters := make([]string, 0, 3)
	args := make([]interface{}, 0, 5)

	if len(classIDs) > 0 {
		filters = append(filters, "class_id IN ?")
		args = append(args, classIDs)
	}
	if len(teacherIDs) > 0 {
		filters = append(filters, "teacher_id IN ?")
		args = append(args, teacherIDs)
	}
	if len(roomIDs) > 0 {
		filters = append(filters, "room_id IN ?")
		args = append(args, roomIDs)
	}

	if len(filters) == 0 {
		return []entities.Lesson{}, nil
	}

	lessons := make([]entities.Lesson, 0)
	query := postgres.GetDb(ctx, r.db).
		Preload("Class").
		Preload("Teacher").
		Preload("Room").
		Where("date_start < ? AND date_end > ?", to, from).
		Where("("+strings.Join(filters, " OR ")+")", args...).
		Order("date_start ASC")
	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}
	err := query.Find(&lessons).Error
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (r *lessonRepository) ListInRange(
	ctx context.Context,
	from time.Time,
	to time.Time,
) ([]entities.Lesson, error) {
	if to.Before(from) {
		to = from
	}

	lessons := make([]entities.Lesson, 0)
	err := postgres.GetDb(ctx, r.db).
		Preload("Class").
		Preload("Teacher").
		Preload("Room").
		Where("date_start < ? AND date_end > ?", to, from).
		Order("date_start ASC").
		Find(&lessons).Error
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (r *lessonRepository) GetLessonWithRelations(ctx context.Context, id string) (*entities.Lesson, error) {
	var lesson entities.Lesson
	err := r.db.WithContext(ctx).
		Preload("Class").
		Preload("Teacher").
		Preload("Room").
		First(&lesson, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}
