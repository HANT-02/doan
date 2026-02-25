package student

import (
	"context"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type CreateStudentInput struct {
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

type CreateStudentOutput struct {
	Student *entities.Student
}

type CreateStudentUseCase interface {
	Execute(ctx context.Context, input CreateStudentInput) (*CreateStudentOutput, error)
}

type createStudentUseCase struct {
	studentRepo repointerface.StudentRepository
}

func NewCreateStudentUseCase(studentRepo repointerface.StudentRepository) CreateStudentUseCase {
	return &createStudentUseCase{
		studentRepo: studentRepo,
	}
}

func (uc *createStudentUseCase) Execute(ctx context.Context, input CreateStudentInput) (*CreateStudentOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	student := &entities.Student{
		Code:          input.Code,
		FullName:      input.FullName,
		Email:         input.Email,
		Phone:         input.Phone,
		GuardianPhone: input.GuardianPhone,
		GradeLevel:    input.GradeLevel,
		Status:        input.Status,
		DateOfBirth:   input.DateOfBirth,
		Gender:        input.Gender,
		Address:       input.Address,
	}

	createdStudent, err := uc.studentRepo.Create(ctx, student)
	if err != nil {
		ctxLogger.Errorf("Failed to create student: %v", err)
		return nil, err
	}

	return &CreateStudentOutput{Student: createdStudent}, nil
}
