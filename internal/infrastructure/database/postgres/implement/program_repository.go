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

type programRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Program]
	db *gorm.DB
}

func NewProgramRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.ProgramRepository {
	modelRepo := postgres.NewBaseRepository[entities.Program](log, manager, db, "programs")
	return &programRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *programRepository) AddCourses(ctx context.Context, programID string, courseIDs []string) error {
	var programCourses []entities.ProgramCourse
	for _, courseID := range courseIDs {
		programCourses = append(programCourses, entities.ProgramCourse{
			ProgramID: programID,
			CourseID:  courseID,
		})
	}
	return r.db.WithContext(ctx).Create(&programCourses).Error
}

func (r *programRepository) RemoveCourses(ctx context.Context, programID string, courseIDs []string) error {
	return r.db.WithContext(ctx).
		Where("program_id = ? AND course_id IN ?", programID, courseIDs).
		Delete(&entities.ProgramCourse{}).Error
}

func (r *programRepository) GetProgramWithCourses(ctx context.Context, id string) (*entities.Program, error) {
	var program entities.Program
	err := r.db.WithContext(ctx).Preload("Courses").First(&program, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &program, nil
}
