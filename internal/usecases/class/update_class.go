package class

import (
	"context"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type UpdateClassInput struct {
	ID          string
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

type UpdateClassOutput struct {
	Class *entities.Class
}

type UpdateClassUseCase interface {
	Execute(ctx context.Context, input UpdateClassInput) (*UpdateClassOutput, error)
}

type updateClassUseCase struct {
	classRepo repointerface.ClassRepository
}

func NewUpdateClassUseCase(classRepo repointerface.ClassRepository) UpdateClassUseCase {
	return &updateClassUseCase{
		classRepo: classRepo,
	}
}

func (uc *updateClassUseCase) Execute(ctx context.Context, input UpdateClassInput) (*UpdateClassOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Fetch existing class to ensure it exists
	classEntity, err := uc.classRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Class not found: %v", err)
		return nil, err
	}

	// Update Data mapping
	updateData := map[string]interface{}{
		"code":         input.Code,
		"name":         input.Name,
		"notes":        input.Notes,
		"start_date":   input.StartDate,
		"end_date":     input.EndDate,
		"max_students": input.MaxStudents,
		"status":       input.Status,
		"price":        input.Price,
		"program_id":   input.ProgramID,
		"course_id":    input.CourseID,
		"teacher_id":   input.TeacherID,
	}

	err = uc.classRepo.Update(ctx, input.ID, updateData)
	if err != nil {
		ctxLogger.Errorf("Failed to update class: %v", err)
		return nil, err
	}

	// Update instance values
	classEntity.Code = input.Code
	classEntity.Name = input.Name
	classEntity.Notes = input.Notes
	classEntity.StartDate = input.StartDate
	classEntity.EndDate = input.EndDate
	classEntity.MaxStudents = input.MaxStudents
	classEntity.Status = input.Status
	classEntity.Price = input.Price
	classEntity.ProgramID = input.ProgramID
	classEntity.CourseID = input.CourseID
	classEntity.TeacherID = input.TeacherID

	return &UpdateClassOutput{Class: classEntity}, nil
}
