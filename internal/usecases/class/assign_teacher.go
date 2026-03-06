package class

import (
	"context"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
)

type AssignTeacherInput struct {
	ClassID   string
	TeacherID string
}

type AssignTeacherOutput struct {
	Message string
}

type AssignTeacherUseCase interface {
	Execute(ctx context.Context, input AssignTeacherInput) (*AssignTeacherOutput, error)
}

type assignTeacherUseCase struct {
	classRepo repointerface.ClassRepository
}

func NewAssignTeacherUseCase(cRepo repointerface.ClassRepository) AssignTeacherUseCase {
	return &assignTeacherUseCase{classRepo: cRepo}
}

func (uc *assignTeacherUseCase) Execute(ctx context.Context, input AssignTeacherInput) (*AssignTeacherOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.ClassID == "" || input.TeacherID == "" {
		return nil, errors.New("class_id and teacher_id are required")
	}

	classEntity, err := uc.classRepo.GetByID(ctx, input.ClassID)
	if err != nil {
		ctxLogger.Errorf("Class not found: %v", err)
		return nil, err
	}

	tID := input.TeacherID
	updateData := map[string]interface{}{
		"teacher_id": &tID,
	}

	err = uc.classRepo.Update(ctx, classEntity.ID, updateData)
	if err != nil {
		ctxLogger.Errorf("Failed to assign teacher: %v", err)
		return nil, err
	}

	return &AssignTeacherOutput{Message: "Teacher assigned to class successfully"}, nil
}
