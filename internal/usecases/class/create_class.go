package class

import (
	"context"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type CreateClassInput struct {
	Code        string
	Name        string
	Notes       string
	StartDate   time.Time
	EndDate     *time.Time
	MaxStudents int
	Status      string
	Price       float64
	ProgramID   *string
	CourseID    *string
	TeacherID   *string
}

type CreateClassOutput struct {
	Class *entities.Class
}

type CreateClassUseCase interface {
	Execute(ctx context.Context, input CreateClassInput) (*CreateClassOutput, error)
}

type createClassUseCase struct {
	classRepo   repointerface.ClassRepository
	teacherRepo repointerface.TeacherRepository
	courseRepo  repointerface.CourseRepository
}

func NewCreateClassUseCase(
	classRepo repointerface.ClassRepository,
	teacherRepo repointerface.TeacherRepository,
	courseRepo repointerface.CourseRepository,
) CreateClassUseCase {
	return &createClassUseCase{
		classRepo:   classRepo,
		teacherRepo: teacherRepo,
		courseRepo:  courseRepo,
	}
}

func (uc *createClassUseCase) Execute(ctx context.Context, input CreateClassInput) (*CreateClassOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.Status == "" {
		input.Status = "OPEN"
	}

	if err := validateTeacherCourseSkills(ctx, uc.teacherRepo, uc.courseRepo, input.TeacherID, input.CourseID); err != nil {
		return nil, err
	}

	classEntity := &entities.Class{
		Code:        input.Code,
		Name:        input.Name,
		Notes:       input.Notes,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		MaxStudents: input.MaxStudents,
		Status:      input.Status,
		Price:       input.Price,
		ProgramID:   input.ProgramID,
		CourseID:    input.CourseID,
		TeacherID:   input.TeacherID,
	}

	createdClass, err := uc.classRepo.Create(ctx, classEntity)
	if err != nil {
		ctxLogger.Errorf("Failed to create class: %v", err)
		return nil, err
	}

	return &CreateClassOutput{Class: createdClass}, nil
}
