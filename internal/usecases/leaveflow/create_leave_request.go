package leaveflow

import (
	"context"
	"strings"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"

	"github.com/lib/pq"
)

type CreateLeaveRequestInput struct {
	Actor        Actor
	LeaveType    string
	ApplyDate    time.Time
	LateMinutes  int
	EarlyMinutes int
	Reason       string
	Documents    []string
	ClassID      *string
	LessonID     *string
	Subject      string
}

type CreateLeaveRequestOutput struct {
	Request entities.LeaveRequest `json:"request"`
}

type CreateLeaveRequestUseCase interface {
	Execute(ctx context.Context, input CreateLeaveRequestInput) (*CreateLeaveRequestOutput, error)
}

type createLeaveRequestUseCase struct {
	leaveRepo      repointerface.LeaveRequestRepository
	studentRepo    repointerface.StudentRepository
	enrollmentRepo repointerface.EnrollmentRepository
	lessonRepo     repointerface.LessonRepository
}

func NewCreateLeaveRequestUseCase(
	leaveRepo repointerface.LeaveRequestRepository,
	studentRepo repointerface.StudentRepository,
	enrollmentRepo repointerface.EnrollmentRepository,
	lessonRepo repointerface.LessonRepository,
) CreateLeaveRequestUseCase {
	return &createLeaveRequestUseCase{
		leaveRepo:      leaveRepo,
		studentRepo:    studentRepo,
		enrollmentRepo: enrollmentRepo,
		lessonRepo:     lessonRepo,
	}
}

func (uc *createLeaveRequestUseCase) Execute(ctx context.Context, input CreateLeaveRequestInput) (*CreateLeaveRequestOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if !isValidLeaveType(input.LeaveType) {
		return nil, ErrInvalidLeaveType
	}

	student, err := resolveStudentByEmail(ctx, uc.studentRepo, input.Actor.Email)
	if err != nil {
		return nil, err
	}

	var classID *string = input.ClassID
	if input.LessonID != nil && *input.LessonID != "" {
		lesson, lessonErr := uc.lessonRepo.GetLessonWithRelations(ctx, *input.LessonID)
		if lessonErr != nil {
			return nil, lessonErr
		}
		if lesson == nil {
			return nil, ErrLeaveRequestNotFound
		}
		classID = &lesson.ClassID
	}

	if classID != nil && *classID != "" {
		enrollments, enrollErr := uc.enrollmentRepo.ListByClassID(ctx, *classID)
		if enrollErr != nil {
			return nil, enrollErr
		}
		isEnrolled := false
		for _, enrollment := range enrollments {
			if enrollment.StudentID == student.ID {
				isEnrolled = true
				break
			}
		}
		if !isEnrolled {
			return nil, ErrStudentNotInClass
		}
	}

	now := time.Now()
	entity := &entities.LeaveRequest{
		StudentID:    student.ID,
		LeaveType:    strings.ToUpper(strings.TrimSpace(input.LeaveType)),
		ApplyDate:    input.ApplyDate,
		LateMinutes:  input.LateMinutes,
		EarlyMinutes: input.EarlyMinutes,
		Reason:       strings.TrimSpace(input.Reason),
		Documents:    pq.StringArray(input.Documents),
		ClassID:      classID,
		LessonID:     input.LessonID,
		Subject:      strings.TrimSpace(input.Subject),
		Status:       "PENDING",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := uc.leaveRepo.Create(ctx, entity)
	if err != nil {
		ctxLogger.Errorf("Failed to create leave request for student %s: %v", student.ID, err)
		return nil, err
	}
	request, err := uc.leaveRepo.GetWithRelations(ctx, created.ID)
	if err != nil {
		return nil, err
	}
	return &CreateLeaveRequestOutput{Request: *request}, nil
}
