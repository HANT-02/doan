package teacher

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
)

// UpdateTeacherInput represents the input for updating a teacher
type UpdateTeacherInput struct {
	ID              string  `json:"id"`
	Code            *string `json:"code"`
	FullName        *string `json:"full_name"`
	Email           *string `json:"email"`
	Phone           *string `json:"phone"`
	IsSchoolTeacher *bool   `json:"is_school_teacher"`
	SchoolName      *string `json:"school_name"`
	EmploymentType  *string `json:"employment_type"`
	Status          *string `json:"status"`
	Notes           *string `json:"notes"`
}

// UpdateTeacherOutput represents the output after updating a teacher
type UpdateTeacherOutput struct {
	Teacher *entities.Teacher `json:"teacher"`
}

// UpdateTeacherUseCase defines the interface for updating a teacher
type UpdateTeacherUseCase interface {
	Execute(ctx context.Context, input UpdateTeacherInput) (*UpdateTeacherOutput, error)
}

type updateTeacherUseCase struct {
	teacherRepo repointerface.TeacherRepository
}

// NewUpdateTeacherUseCase creates a new instance of UpdateTeacherUseCase
func NewUpdateTeacherUseCase(teacherRepo repointerface.TeacherRepository) UpdateTeacherUseCase {
	return &updateTeacherUseCase{
		teacherRepo: teacherRepo,
	}
}

func (uc *updateTeacherUseCase) Execute(ctx context.Context, input UpdateTeacherInput) (*UpdateTeacherOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.ID == "" {
		return nil, errors.New("teacher ID is required")
	}

	// Check if teacher exists
	existingTeacher, err := uc.teacherRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to get teacher: %v", err)
		return nil, err
	}

	// Build update map
	updateData := make(map[string]interface{})

	if input.Code != nil {
		// Check code uniqueness if changing
		if *input.Code != existingTeacher.Code {
			exists, err := uc.teacherRepo.ExistsByCode(ctx, *input.Code)
			if err != nil {
				ctxLogger.Errorf("Failed to check code existence: %v", err)
				return nil, err
			}
			if exists {
				return nil, errors.New("code already exists")
			}
		}
		updateData["code"] = *input.Code
	}

	if input.Email != nil {
		// Check email uniqueness if changing
		if *input.Email != existingTeacher.Email {
			exists, err := uc.teacherRepo.ExistsByEmail(ctx, *input.Email)
			if err != nil {
				ctxLogger.Errorf("Failed to check email existence: %v", err)
				return nil, err
			}
			if exists {
				return nil, errors.New("email already exists")
			}
		}
		updateData["email"] = *input.Email
	}

	if input.FullName != nil {
		updateData["full_name"] = *input.FullName
	}
	if input.Phone != nil {
		updateData["phone"] = *input.Phone
	}
	if input.IsSchoolTeacher != nil {
		updateData["is_school_teacher"] = *input.IsSchoolTeacher
	}
	if input.SchoolName != nil {
		updateData["school_name"] = *input.SchoolName
	}
	if input.EmploymentType != nil {
		updateData["employment_type"] = *input.EmploymentType
	}
	if input.Status != nil {
		updateData["status"] = *input.Status
	}
	if input.Notes != nil {
		updateData["notes"] = *input.Notes
	}

	// Update teacher
	err = uc.teacherRepo.Update(ctx, input.ID, updateData)
	if err != nil {
		ctxLogger.Errorf("Failed to update teacher: %v", err)
		return nil, err
	}

	// Fetch updated teacher
	updatedTeacher, err := uc.teacherRepo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to get updated teacher: %v", err)
		return nil, err
	}

	return &UpdateTeacherOutput{Teacher: updatedTeacher}, nil
}
