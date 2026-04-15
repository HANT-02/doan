package leaveflow

import (
	"context"

	repointerface "doan/internal/repositories/interface"
)

type CancelLeaveRequestInput struct {
	ID    string
	Actor Actor
}

type CancelLeaveRequestUseCase interface {
	Execute(ctx context.Context, input CancelLeaveRequestInput) error
}

type cancelLeaveRequestUseCase struct {
	leaveRepo   repointerface.LeaveRequestRepository
	studentRepo repointerface.StudentRepository
}

func NewCancelLeaveRequestUseCase(
	leaveRepo repointerface.LeaveRequestRepository,
	studentRepo repointerface.StudentRepository,
) CancelLeaveRequestUseCase {
	return &cancelLeaveRequestUseCase{
		leaveRepo:   leaveRepo,
		studentRepo: studentRepo,
	}
}

func (uc *cancelLeaveRequestUseCase) Execute(ctx context.Context, input CancelLeaveRequestInput) error {
	request, err := uc.leaveRepo.GetWithRelations(ctx, input.ID)
	if err != nil {
		return err
	}
	if request == nil {
		return ErrLeaveRequestNotFound
	}
	if request.Status != "PENDING" {
		return ErrLeaveRequestNotPending
	}

	student, err := resolveStudentByEmail(ctx, uc.studentRepo, input.Actor.Email)
	if err != nil {
		return err
	}
	if request.StudentID != student.ID {
		return ErrLeaveRequestForbidden
	}

	return uc.leaveRepo.Update(ctx, input.ID, map[string]interface{}{
		"status": "CANCELLED",
	})
}
