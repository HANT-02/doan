package leaveflow

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type ListLeaveRequestsInput struct {
	Actor     Actor
	Status    string
	ClassID   string
	StudentID string
}

type ListLeaveRequestsOutput struct {
	Requests []entities.LeaveRequest `json:"requests"`
}

type ListLeaveRequestsUseCase interface {
	Execute(ctx context.Context, input ListLeaveRequestsInput) (*ListLeaveRequestsOutput, error)
}

type listLeaveRequestsUseCase struct {
	leaveRepo   repointerface.LeaveRequestRepository
	studentRepo repointerface.StudentRepository
	teacherRepo repointerface.TeacherRepository
	classRepo   repointerface.ClassRepository
}

func NewListLeaveRequestsUseCase(
	leaveRepo repointerface.LeaveRequestRepository,
	studentRepo repointerface.StudentRepository,
	teacherRepo repointerface.TeacherRepository,
	classRepo repointerface.ClassRepository,
) ListLeaveRequestsUseCase {
	return &listLeaveRequestsUseCase{
		leaveRepo:   leaveRepo,
		studentRepo: studentRepo,
		teacherRepo: teacherRepo,
		classRepo:   classRepo,
	}
}

func (uc *listLeaveRequestsUseCase) Execute(ctx context.Context, input ListLeaveRequestsInput) (*ListLeaveRequestsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	filter := repointerface.LeaveRequestFilter{
		Status:    input.Status,
		ClassID:   input.ClassID,
		StudentID: input.StudentID,
	}

	switch input.Actor.Role {
	case "ADMIN", "SUPER_ADMIN":
	case "TEACHER":
		teacher, err := resolveTeacherByEmail(ctx, uc.teacherRepo, input.Actor.Email)
		if err != nil {
			return nil, err
		}
		condition := repositories.NewCommonCondition()
		condition.AddCondition("teacher_id", teacher.ID, repositories.Equal)
		condition.SetPaging(200, 1)
		classes, err := uc.classRepo.GetByCondition(ctx, condition)
		if err != nil {
			ctxLogger.Errorf("Failed to list teacher classes for leave requests: %v", err)
			return nil, err
		}
		classIDs := make([]string, 0)
		if classes != nil {
			for _, classItem := range classes.Data {
				if classItem != nil {
					classIDs = append(classIDs, classItem.ID)
				}
			}
		}
		if len(classIDs) == 0 {
			return &ListLeaveRequestsOutput{Requests: []entities.LeaveRequest{}}, nil
		}
		filter.ClassIDs = classIDs
	case "STUDENT", "PARENT":
		student, err := resolveStudentByEmail(ctx, uc.studentRepo, input.Actor.Email)
		if err != nil {
			return nil, err
		}
		filter.StudentID = student.ID
	default:
		return nil, ErrLeaveRequestForbidden
	}

	requests, err := uc.leaveRepo.ListWithRelations(ctx, filter)
	if err != nil {
		ctxLogger.Errorf("Failed to list leave requests: %v", err)
		return nil, err
	}
	return &ListLeaveRequestsOutput{Requests: requests}, nil
}
