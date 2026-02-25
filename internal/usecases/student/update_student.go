package student

import (
	"context"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type UpdateStudentInput struct {
	ID            string
	Code          string
	FullName      string
	Email         string
	Phone         string
	GuardianPhone string
	GradeLevel    string
	Status        string
	DateOfBirth   *time.Time
	Gender        string
	Address       string
}

type UpdateStudentOutput struct {
	Student *entities.Student
}

type UpdateStudentUseCase interface {
	Execute(ctx context.Context, input UpdateStudentInput) (*UpdateStudentOutput, error)
}

type updateStudentUseCase struct {
	studentRepo repointerface.StudentRepository
}

func NewUpdateStudentUseCase(studentRepo repointerface.StudentRepository) UpdateStudentUseCase {
	return &updateStudentUseCase{
		studentRepo: studentRepo,
	}
}

func (uc *updateStudentUseCase) Execute(ctx context.Context, input UpdateStudentInput) (*UpdateStudentOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Build update map
	updateData := make(map[string]interface{})
	updateData["code"] = input.Code
	updateData["full_name"] = input.FullName
	updateData["email"] = input.Email
	updateData["phone"] = input.Phone
	updateData["guardian_phone"] = input.GuardianPhone
	updateData["grade_level"] = input.GradeLevel
	updateData["status"] = input.Status
	updateData["date_of_birth"] = input.DateOfBirth
	updateData["gender"] = input.Gender
	updateData["address"] = input.Address

	err := uc.studentRepo.Update(ctx, input.ID, updateData)
	if err != nil {
		ctxLogger.Errorf("Failed to update student: %v", err)
		return nil, err
	}

	// Fetch updated student
	updatedStudent, err := uc.studentRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to get updated student: %v", err)
		return nil, err
	}

	return &UpdateStudentOutput{Student: updatedStudent}, nil
}
