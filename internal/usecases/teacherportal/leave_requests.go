package teacherportal

import (
	"context"

	leaveflow "doan/internal/usecases/leaveflow"
)

type ListLeaveRequestsForTeacherInput struct {
	Actor     Actor
	Status    string
	ClassID   string
	StudentID string
}

type ListLeaveRequestsForTeacherOutput struct {
	Requests []LeaveRequestItem `json:"requests"`
}

type ApproveLeaveRequestInput struct {
	Actor Actor
	ID    string
}

type ApproveLeaveRequestOutput struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}

type RejectLeaveRequestInput struct {
	Actor           Actor
	ID              string
	RejectionReason string
}

type RejectLeaveRequestOutput struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}

type ListLeaveRequestsForTeacherUseCase interface {
	Execute(ctx context.Context, input ListLeaveRequestsForTeacherInput) (*ListLeaveRequestsForTeacherOutput, error)
}

type ApproveLeaveRequestUseCase interface {
	Execute(ctx context.Context, input ApproveLeaveRequestInput) (*ApproveLeaveRequestOutput, error)
}

type RejectLeaveRequestUseCase interface {
	Execute(ctx context.Context, input RejectLeaveRequestInput) (*RejectLeaveRequestOutput, error)
}

type listLeaveRequestsForTeacherUseCase struct {
	listLeaveRequestsUseCase leaveflow.ListLeaveRequestsUseCase
}

type approveLeaveRequestUseCase struct {
	updateLeaveRequestStatusUseCase leaveflow.UpdateLeaveRequestStatusUseCase
}

type rejectLeaveRequestUseCase struct {
	updateLeaveRequestStatusUseCase leaveflow.UpdateLeaveRequestStatusUseCase
}

func NewListLeaveRequestsForTeacherUseCase(
	listLeaveRequestsUseCase leaveflow.ListLeaveRequestsUseCase,
) ListLeaveRequestsForTeacherUseCase {
	return &listLeaveRequestsForTeacherUseCase{
		listLeaveRequestsUseCase: listLeaveRequestsUseCase,
	}
}

func NewApproveLeaveRequestUseCase(
	updateLeaveRequestStatusUseCase leaveflow.UpdateLeaveRequestStatusUseCase,
) ApproveLeaveRequestUseCase {
	return &approveLeaveRequestUseCase{
		updateLeaveRequestStatusUseCase: updateLeaveRequestStatusUseCase,
	}
}

func NewRejectLeaveRequestUseCase(
	updateLeaveRequestStatusUseCase leaveflow.UpdateLeaveRequestStatusUseCase,
) RejectLeaveRequestUseCase {
	return &rejectLeaveRequestUseCase{
		updateLeaveRequestStatusUseCase: updateLeaveRequestStatusUseCase,
	}
}

func (uc *listLeaveRequestsForTeacherUseCase) Execute(ctx context.Context, input ListLeaveRequestsForTeacherInput) (*ListLeaveRequestsForTeacherOutput, error) {
	output, err := uc.listLeaveRequestsUseCase.Execute(ctx, leaveflow.ListLeaveRequestsInput{
		Actor: leaveflow.Actor{
			Role:   input.Actor.Role,
			Email:  input.Actor.Email,
			UserID: input.Actor.UserID,
		},
		Status:    input.Status,
		ClassID:   input.ClassID,
		StudentID: input.StudentID,
	})
	if err != nil {
		return nil, err
	}

	requests := make([]LeaveRequestItem, 0, len(output.Requests))
	for _, request := range output.Requests {
		requests = append(requests, buildLeaveRequestItem(request))
	}

	return &ListLeaveRequestsForTeacherOutput{Requests: requests}, nil
}

func (uc *approveLeaveRequestUseCase) Execute(ctx context.Context, input ApproveLeaveRequestInput) (*ApproveLeaveRequestOutput, error) {
	output, err := uc.updateLeaveRequestStatusUseCase.Execute(ctx, leaveflow.UpdateLeaveRequestStatusInput{
		ID:     input.ID,
		Actor:  leaveflow.Actor{Role: input.Actor.Role, Email: input.Actor.Email, UserID: input.Actor.UserID},
		Status: "APPROVED",
	})
	if err != nil {
		return nil, err
	}
	return &ApproveLeaveRequestOutput{RequestID: output.RequestID, Status: output.Status}, nil
}

func (uc *rejectLeaveRequestUseCase) Execute(ctx context.Context, input RejectLeaveRequestInput) (*RejectLeaveRequestOutput, error) {
	output, err := uc.updateLeaveRequestStatusUseCase.Execute(ctx, leaveflow.UpdateLeaveRequestStatusInput{
		ID:              input.ID,
		Actor:           leaveflow.Actor{Role: input.Actor.Role, Email: input.Actor.Email, UserID: input.Actor.UserID},
		Status:          "REJECTED",
		RejectionReason: input.RejectionReason,
	})
	if err != nil {
		return nil, err
	}
	return &RejectLeaveRequestOutput{RequestID: output.RequestID, Status: output.Status}, nil
}
