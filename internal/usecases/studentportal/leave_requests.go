package studentportal

import (
	"context"
	"time"

	leaveflow "doan/internal/usecases/leaveflow"
)

type StudentLeaveRequestStudentItem struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"full_name"`
}

type StudentLeaveRequestClassItem struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type StudentLeaveRequestLessonItem struct {
	ID        string    `json:"id"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
}

type StudentLeaveRequestItem struct {
	ID              string                         `json:"id"`
	Student         StudentLeaveRequestStudentItem `json:"student"`
	LeaveType       string                         `json:"leave_type"`
	ApplyDate       time.Time                      `json:"apply_date"`
	LateMinutes     int                            `json:"late_minutes"`
	EarlyMinutes    int                            `json:"early_minutes"`
	Reason          string                         `json:"reason"`
	Documents       []string                       `json:"documents"`
	Class           *StudentLeaveRequestClassItem  `json:"class,omitempty"`
	Lesson          *StudentLeaveRequestLessonItem `json:"lesson,omitempty"`
	Subject         string                         `json:"subject"`
	Status          string                         `json:"status"`
	ApprovedByID    *string                        `json:"approved_by_id,omitempty"`
	ApprovedAt      *time.Time                     `json:"approved_at,omitempty"`
	RejectionReason string                         `json:"rejection_reason"`
	CreatedAt       time.Time                      `json:"created_at"`
	UpdatedAt       time.Time                      `json:"updated_at"`
}

type ListMyLeaveRequestsInput struct {
	Actor   Actor
	Status  string
	ClassID string
}

type ListMyLeaveRequestsOutput struct {
	Requests []StudentLeaveRequestItem `json:"requests"`
}

type CreateMyLeaveRequestInput struct {
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

type CreateMyLeaveRequestOutput struct {
	Request StudentLeaveRequestItem `json:"request"`
}

type CancelMyLeaveRequestInput struct {
	Actor Actor
	ID    string
}

type CancelMyLeaveRequestOutput struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}

type ListMyLeaveRequestsUseCase interface {
	Execute(ctx context.Context, input ListMyLeaveRequestsInput) (*ListMyLeaveRequestsOutput, error)
}

type CreateMyLeaveRequestUseCase interface {
	Execute(ctx context.Context, input CreateMyLeaveRequestInput) (*CreateMyLeaveRequestOutput, error)
}

type CancelMyLeaveRequestUseCase interface {
	Execute(ctx context.Context, input CancelMyLeaveRequestInput) (*CancelMyLeaveRequestOutput, error)
}

type listMyLeaveRequestsUseCase struct {
	listLeaveRequestsUseCase leaveflow.ListLeaveRequestsUseCase
}

type createMyLeaveRequestUseCase struct {
	createLeaveRequestUseCase leaveflow.CreateLeaveRequestUseCase
}

type cancelMyLeaveRequestUseCase struct {
	cancelLeaveRequestUseCase leaveflow.CancelLeaveRequestUseCase
}

func NewListMyLeaveRequestsUseCase(
	listLeaveRequestsUseCase leaveflow.ListLeaveRequestsUseCase,
) ListMyLeaveRequestsUseCase {
	return &listMyLeaveRequestsUseCase{
		listLeaveRequestsUseCase: listLeaveRequestsUseCase,
	}
}

func NewCreateMyLeaveRequestUseCase(
	createLeaveRequestUseCase leaveflow.CreateLeaveRequestUseCase,
) CreateMyLeaveRequestUseCase {
	return &createMyLeaveRequestUseCase{
		createLeaveRequestUseCase: createLeaveRequestUseCase,
	}
}

func NewCancelMyLeaveRequestUseCase(
	cancelLeaveRequestUseCase leaveflow.CancelLeaveRequestUseCase,
) CancelMyLeaveRequestUseCase {
	return &cancelMyLeaveRequestUseCase{
		cancelLeaveRequestUseCase: cancelLeaveRequestUseCase,
	}
}

func (uc *listMyLeaveRequestsUseCase) Execute(ctx context.Context, input ListMyLeaveRequestsInput) (*ListMyLeaveRequestsOutput, error) {
	output, err := uc.listLeaveRequestsUseCase.Execute(ctx, leaveflow.ListLeaveRequestsInput{
		Actor: leaveflow.Actor{
			UserID: input.Actor.UserID,
			Email:  input.Actor.Email,
			Role:   input.Actor.Role,
		},
		Status:  input.Status,
		ClassID: input.ClassID,
	})
	if err != nil {
		return nil, err
	}

	requests := make([]StudentLeaveRequestItem, 0, len(output.Requests))
	for _, request := range output.Requests {
		requests = append(requests, buildStudentLeaveRequestItem(request))
	}

	return &ListMyLeaveRequestsOutput{Requests: requests}, nil
}

func (uc *createMyLeaveRequestUseCase) Execute(ctx context.Context, input CreateMyLeaveRequestInput) (*CreateMyLeaveRequestOutput, error) {
	output, err := uc.createLeaveRequestUseCase.Execute(ctx, leaveflow.CreateLeaveRequestInput{
		Actor: leaveflow.Actor{
			UserID: input.Actor.UserID,
			Email:  input.Actor.Email,
			Role:   input.Actor.Role,
		},
		LeaveType:    input.LeaveType,
		ApplyDate:    input.ApplyDate,
		LateMinutes:  input.LateMinutes,
		EarlyMinutes: input.EarlyMinutes,
		Reason:       input.Reason,
		Documents:    input.Documents,
		ClassID:      input.ClassID,
		LessonID:     input.LessonID,
		Subject:      input.Subject,
	})
	if err != nil {
		return nil, err
	}

	return &CreateMyLeaveRequestOutput{
		Request: buildStudentLeaveRequestItem(output.Request),
	}, nil
}

func (uc *cancelMyLeaveRequestUseCase) Execute(ctx context.Context, input CancelMyLeaveRequestInput) (*CancelMyLeaveRequestOutput, error) {
	err := uc.cancelLeaveRequestUseCase.Execute(ctx, leaveflow.CancelLeaveRequestInput{
		ID: input.ID,
		Actor: leaveflow.Actor{
			UserID: input.Actor.UserID,
			Email:  input.Actor.Email,
			Role:   input.Actor.Role,
		},
	})
	if err != nil {
		return nil, err
	}

	return &CancelMyLeaveRequestOutput{
		RequestID: input.ID,
		Status:    "CANCELLED",
	}, nil
}
