package teacher

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
)

// CreateTeacherInput represents the input for creating a teacher
type CreateTeacherInput struct {
	Code            string `json:"code"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	IsSchoolTeacher bool   `json:"is_school_teacher"`
	SchoolName      string `json:"school_name"`
	EmploymentType  string `json:"employment_type"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
}

// CreateTeacherOutput represents the output after creating a teacher
type CreateTeacherOutput struct {
	Teacher *entities.Teacher `json:"teacher"`
}

// CreateTeacherUseCase defines the interface for creating a teacher
type CreateTeacherUseCase interface {
	Execute(ctx context.Context, input CreateTeacherInput) (*CreateTeacherOutput, error)
}

type createTeacherUseCase struct {
	teacherRepo repointerface.TeacherRepository
}

// NewCreateTeacherUseCase creates a new instance of CreateTeacherUseCase
func NewCreateTeacherUseCase(teacherRepo repointerface.TeacherRepository) CreateTeacherUseCase {
	return &createTeacherUseCase{
		teacherRepo: teacherRepo,
	}
}

func (uc *createTeacherUseCase) Execute(ctx context.Context, input CreateTeacherInput) (*CreateTeacherOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Validate required fields
	if input.FullName == "" {
		return nil, errors.New("full_name is required")
	}

	// Check email uniqueness if provided
	if input.Email != "" {
		exists, err := uc.teacherRepo.ExistsByEmail(ctx, input.Email)
		if err != nil {
			ctxLogger.Errorf("Failed to check email existence: %v", err)
			return nil, err
		}
		if exists {
			return nil, errors.New("email already exists")
		}
	}

	// Check code uniqueness if provided
	if input.Code != "" {
		exists, err := uc.teacherRepo.ExistsByCode(ctx, input.Code)
		if err != nil {
			ctxLogger.Errorf("Failed to check code existence: %v", err)
			return nil, err
		}
		if exists {
			return nil, errors.New("code already exists")
		}
	}

	// Set default values
	if input.EmploymentType == "" {
		input.EmploymentType = "PART_TIME"
	}
	if input.Status == "" {
		input.Status = "ACTIVE"
	}

	// Create teacher entity
	teacher := &entities.Teacher{
		Code:            input.Code,
		FullName:        input.FullName,
		Email:           input.Email,
		Phone:           input.Phone,
		IsSchoolTeacher: input.IsSchoolTeacher,
		SchoolName:      input.SchoolName,
		EmploymentType:  input.EmploymentType,
		Status:          input.Status,
		Notes:           input.Notes,
	}

	// Save to database
	createdTeacher, err := uc.teacherRepo.Create(ctx, teacher)
	if err != nil {
		ctxLogger.Errorf("Failed to create teacher: %v", err)
		return nil, err
	}

	return &CreateTeacherOutput{Teacher: createdTeacher}, nil
}
