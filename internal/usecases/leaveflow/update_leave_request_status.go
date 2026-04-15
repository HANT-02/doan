package leaveflow

import (
	"context"
	"time"

	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type UpdateLeaveRequestStatusInput struct {
	ID              string
	Actor           Actor
	Status          string
	RejectionReason string
}

type UpdateLeaveRequestStatusOutput struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}

type UpdateLeaveRequestStatusUseCase interface {
	Execute(ctx context.Context, input UpdateLeaveRequestStatusInput) (*UpdateLeaveRequestStatusOutput, error)
}

type updateLeaveRequestStatusUseCase struct {
	leaveRepo   repointerface.LeaveRequestRepository
	teacherRepo repointerface.TeacherRepository
	classRepo   repointerface.ClassRepository
}

func NewUpdateLeaveRequestStatusUseCase(
	leaveRepo repointerface.LeaveRequestRepository,
	teacherRepo repointerface.TeacherRepository,
	classRepo repointerface.ClassRepository,
) UpdateLeaveRequestStatusUseCase {
	return &updateLeaveRequestStatusUseCase{
		leaveRepo:   leaveRepo,
		teacherRepo: teacherRepo,
		classRepo:   classRepo,
	}
}

func (uc *updateLeaveRequestStatusUseCase) Execute(ctx context.Context, input UpdateLeaveRequestStatusInput) (*UpdateLeaveRequestStatusOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	request, err := uc.leaveRepo.GetWithRelations(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, ErrLeaveRequestNotFound
	}
	if request.Status != "PENDING" {
		return nil, ErrLeaveRequestNotPending
	}

	switch input.Actor.Role {
	case "ADMIN", "SUPER_ADMIN":
	case "TEACHER":
		teacher, resolveErr := resolveTeacherByEmail(ctx, uc.teacherRepo, input.Actor.Email)
		if resolveErr != nil {
			return nil, resolveErr
		}
		if request.ClassID == nil || *request.ClassID == "" {
			return nil, ErrLeaveRequestForbidden
		}
		classEntity, classErr := uc.classRepo.GetByID(ctx, *request.ClassID)
		if classErr != nil {
			return nil, classErr
		}
		if classEntity == nil || classEntity.TeacherID == nil || *classEntity.TeacherID != teacher.ID {
			return nil, ErrLeaveRequestForbidden
		}
	default:
		return nil, ErrLeaveRequestForbidden
	}

	updateData := map[string]interface{}{
		"status": input.Status,
	}
	if input.Status == "APPROVED" {
		now := time.Now()
		updateData["approved_by_id"] = input.Actor.UserID
		updateData["approved_at"] = now
		updateData["rejection_reason"] = ""
	} else {
		updateData["approved_by_id"] = nil
		updateData["approved_at"] = nil
		updateData["rejection_reason"] = input.RejectionReason
	}

	if err := uc.leaveRepo.Update(ctx, input.ID, updateData); err != nil {
		ctxLogger.Errorf("Failed to update leave request %s: %v", input.ID, err)
		return nil, err
	}

	return &UpdateLeaveRequestStatusOutput{RequestID: input.ID, Status: input.Status}, nil
}
