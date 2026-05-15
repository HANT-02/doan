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

	existingEnrollments, err := uc.enrollmentRepo.ListByClassID(ctx, input.ClassID)
	if err != nil {
		ctxLogger.Errorf("Failed to load existing enrollments: %v", err)
		return nil, err
	}

	activeStudentIDs := make(map[string]struct{}, len(existingEnrollments))
	existingCount := 0
	for _, enrollment := range existingEnrollments {
		if !isActiveEnrollmentStatus(enrollment.Status) {
			continue
		}
		existingCount++
		activeStudentIDs[enrollment.StudentID] = struct{}{}
	}

	// Create enrollments
	enrolledCount := 0
	for _, sID := range input.StudentIDs {
		if _, exists := activeStudentIDs[sID]; exists {
			continue
		}
		if maxAllowed > 0 && existingCount+enrolledCount >= maxAllowed {
			ctxLogger.Warnf("Capacity reached (%d). Skipping remaining students.", maxAllowed)
			break
		}

		enroll := &entities.Enrollment{
			ClassID:   input.ClassID,
			StudentID: sID,
			Status:    entities.EnrollmentStatusEnrolled,
		}
		_, err := uc.enrollmentRepo.Create(ctx, enroll)
		if err != nil {
			ctxLogger.Errorf("Error enrolling student %s: %v", sID, err)
			continue
		}
		enrolledCount++
		activeStudentIDs[sID] = struct{}{}
	}

	return &EnrollStudentsOutput{Enrolled: enrolledCount}, nil
}
