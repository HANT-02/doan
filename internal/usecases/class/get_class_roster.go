package class

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetClassRosterInput struct {
	ClassID string
}

type GetClassRosterOutput struct {
	ClassID          string             `json:"class_id"`
	MaxStudents      int                `json:"max_students"`
	CapacityLimit    int                `json:"capacity_limit"`
	CurrentCount     int                `json:"current_count"`
	Students         []entities.Student `json:"students"`
	ReservedStudents []entities.Student `json:"reserved_students"`
}

type GetClassRosterUseCase interface {
	Execute(ctx context.Context, input GetClassRosterInput) (*GetClassRosterOutput, error)
}

type getClassRosterUseCase struct {
	classRepo      repointerface.ClassRepository
	enrollmentRepo repointerface.EnrollmentRepository
	roomRepo       repointerface.RoomRepository
}

func NewGetClassRosterUseCase(
	classRepo repointerface.ClassRepository,
	enrollmentRepo repointerface.EnrollmentRepository,
	roomRepo repointerface.RoomRepository,
) GetClassRosterUseCase {
	return &getClassRosterUseCase{
		classRepo:      classRepo,
		enrollmentRepo: enrollmentRepo,
		roomRepo:       roomRepo,
	}
}

func (uc *getClassRosterUseCase) Execute(ctx context.Context, input GetClassRosterInput) (*GetClassRosterOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	classEntity, err := uc.classRepo.GetByID(ctx, input.ClassID)
	if err != nil {
		ctxLogger.Errorf("Failed to load class for roster: %v", err)
		return nil, err
	}

	enrollments, err := uc.enrollmentRepo.ListByClassID(ctx, input.ClassID)
	if err != nil {
		ctxLogger.Errorf("Failed to load class roster: %v", err)
		return nil, err
	}

	students := make([]entities.Student, 0, len(enrollments))
	reservedStudents := make([]entities.Student, 0)
	for _, enrollment := range enrollments {
		if isActiveEnrollmentStatus(enrollment.Status) {
			students = append(students, enrollment.Student)
			continue
		}
		if enrollment.Status == entities.EnrollmentStatusSuspended {
			reservedStudents = append(reservedStudents, enrollment.Student)
		}
	}

	capacityLimit := classEntity.MaxStudents
	if classEntity.RoomID != nil && *classEntity.RoomID != "" {
		roomEntity, roomErr := uc.roomRepo.GetByID(ctx, *classEntity.RoomID)
		if roomErr == nil && roomEntity != nil && roomEntity.Capacity > 0 && roomEntity.Capacity < capacityLimit {
			capacityLimit = roomEntity.Capacity
		}
	}

	return &GetClassRosterOutput{
		ClassID:          classEntity.ID,
		MaxStudents:      classEntity.MaxStudents,
		CapacityLimit:    capacityLimit,
		CurrentCount:     len(students),
		Students:         students,
		ReservedStudents: reservedStudents,
	}, nil
}
