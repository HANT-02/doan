package class

import (
	"context"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
)

type RemoveStudentsInput struct {
	ClassID    string
	StudentIDs []string
}

type RemoveStudentsOutput struct {
	Message string
}

type RemoveStudentsUseCase interface {
	Execute(ctx context.Context, input RemoveStudentsInput) (*RemoveStudentsOutput, error)
}

type removeStudentsUseCase struct {
	enrollmentRepo repointerface.EnrollmentRepository
}

func NewRemoveStudentsUseCase(eRepo repointerface.EnrollmentRepository) RemoveStudentsUseCase {
	return &removeStudentsUseCase{enrollmentRepo: eRepo}
}

func (uc *removeStudentsUseCase) Execute(ctx context.Context, input RemoveStudentsInput) (*RemoveStudentsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.ClassID == "" || len(input.StudentIDs) == 0 {
		return nil, errors.New("class_id and student_ids are required")
	}

	// This is slightly tricky with BaseRepo because we need to delete by multiple criteria (ClassID AND StudentID)
	// For simplicity, we can fetch all enrollments and delete them, or extend BaseRepo.
	// Let's assume we can fetch all enrollments for a class.
	// We'll iterate the students, find the associated enrollment, and delete it.

	// Find all enrollments for a class
	condition := &repositories.CommonCondition{
		Conditions: []repositories.Condition{
			{Field: "class_id", Op: repositories.Equal, Value: input.ClassID},
		},
		Paging: &repositories.Paging{
			Page:  1,
			Limit: 1000,
		},
	}
	enrollmentsResult, err := uc.enrollmentRepo.GetByCondition(ctx, condition)
	if err != nil {
		ctxLogger.Errorf("Failed to retrieve enrollments: %v", err)
		return nil, err
	}

	if enrollmentsResult != nil && len(enrollmentsResult.Data) > 0 {
		studentsToDel := make(map[string]bool)
		for _, sid := range input.StudentIDs {
			studentsToDel[sid] = true
		}

		for _, e := range enrollmentsResult.Data {
			if studentsToDel[e.StudentID] {
				err = uc.enrollmentRepo.SoftDelete(ctx, e.ID)
				if err != nil {
					ctxLogger.Errorf("Failed to delete enrollment %s: %v", e.ID, err)
				}
			}
		}
	}

	return &RemoveStudentsOutput{Message: "Removed students from class"}, nil
}
