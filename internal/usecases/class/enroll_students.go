package class

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
)

type EnrollStudentsInput struct {
	ClassID    string
	StudentIDs []string
}

type EnrollStudentsOutput struct {
	Enrolled int
}

type EnrollStudentsUseCase interface {
	Execute(ctx context.Context, input EnrollStudentsInput) (*EnrollStudentsOutput, error)
}

type enrollStudentsUseCase struct {
	classRepo      repointerface.ClassRepository
	enrollmentRepo repointerface.EnrollmentRepository
	roomRepo       repointerface.RoomRepository
}

func NewEnrollStudentsUseCase(cRepo repointerface.ClassRepository, eRepo repointerface.EnrollmentRepository, rRepo repointerface.RoomRepository) EnrollStudentsUseCase {
	return &enrollStudentsUseCase{classRepo: cRepo, enrollmentRepo: eRepo, roomRepo: rRepo}
}

func (uc *enrollStudentsUseCase) Execute(ctx context.Context, input EnrollStudentsInput) (*EnrollStudentsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.ClassID == "" || len(input.StudentIDs) == 0 {
		return nil, errors.New("class_id and student_ids are required")
	}

	classEntity, err := uc.classRepo.GetByID(ctx, input.ClassID)
	if err != nil {
		ctxLogger.Errorf("Class not found: %v", err)
		return nil, err
	}

	// Calculate maximum capacity constraints
	maxAllowed := classEntity.MaxStudents
	if classEntity.RoomID != nil && *classEntity.RoomID != "" {
		roomEntity, rErr := uc.roomRepo.GetByID(ctx, *classEntity.RoomID)
		if rErr == nil && roomEntity.Capacity < maxAllowed {
			maxAllowed = roomEntity.Capacity // room limit is stricter
		}
	}

	// We would ideally query existing enrollments count here to see if remaining seats exist
	// For now, we will just assume we can enroll up to maxAllowed if provided.
	// In a real system you'd count `SELECT COUNT(*) FROM enrollments WHERE class_id = ?`
	// but currently the base repo may not support Count natively easily without extending.

	// Create enrollments
	enrolledCount := 0
	for _, sID := range input.StudentIDs {
		// Basic naive capacity check logic (assuming 0 existing enrollments for simplicity in this demo)
		// If needed, we'll implement a custom count in repo later over-budget.
		if maxAllowed > 0 && enrolledCount >= maxAllowed {
			ctxLogger.Warnf("Capacity reached (%d). Skipping remaining students.", maxAllowed)
			break
		}

		enroll := &entities.Enrollment{
			ClassID:   input.ClassID,
			StudentID: sID,
			Status:    "ENROLLED",
		}
		_, err := uc.enrollmentRepo.Create(ctx, enroll)
		if err != nil {
			ctxLogger.Errorf("Error enrolling student %s: %v", sID, err)
			continue
		}
		enrolledCount++
	}

	return &EnrollStudentsOutput{Enrolled: enrolledCount}, nil
}
