package leave

import (
	"net/http"
	"time"

	"doan/cmd/http/rest"
	leaveflow "doan/internal/usecases/leaveflow"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	listLeaveRequestsUseCase        leaveflow.ListLeaveRequestsUseCase
	createLeaveRequestUseCase       leaveflow.CreateLeaveRequestUseCase
	updateLeaveRequestStatusUseCase leaveflow.UpdateLeaveRequestStatusUseCase
	cancelLeaveRequestUseCase       leaveflow.CancelLeaveRequestUseCase
}

func NewLeaveControllerV1(
	listLeaveRequestsUseCase leaveflow.ListLeaveRequestsUseCase,
	createLeaveRequestUseCase leaveflow.CreateLeaveRequestUseCase,
	updateLeaveRequestStatusUseCase leaveflow.UpdateLeaveRequestStatusUseCase,
	cancelLeaveRequestUseCase leaveflow.CancelLeaveRequestUseCase,
) *ControllerV1 {
	return &ControllerV1{
		listLeaveRequestsUseCase:        listLeaveRequestsUseCase,
		createLeaveRequestUseCase:       createLeaveRequestUseCase,
		updateLeaveRequestStatusUseCase: updateLeaveRequestStatusUseCase,
		cancelLeaveRequestUseCase:       cancelLeaveRequestUseCase,
	}
}

func (ctrl *ControllerV1) ListLeaveRequests(c *gin.Context) {
	output, err := ctrl.listLeaveRequestsUseCase.Execute(c.Request.Context(), leaveflow.ListLeaveRequestsInput{
		Actor:     buildActor(c),
		Status:    c.Query("status"),
		ClassID:   c.Query("class_id"),
		StudentID: c.Query("student_id"),
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to list leave requests", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Leave requests retrieved successfully", output)
}

func (ctrl *ControllerV1) CreateLeaveRequest(c *gin.Context) {
	var req CreateLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	applyDate, err := time.Parse(time.RFC3339, req.ApplyDate)
	if err != nil {
		applyDate, err = time.Parse("2006-01-02T15:04", req.ApplyDate)
		if err != nil {
			rest.ResponseError(c, http.StatusBadRequest, "Invalid apply_date format", err)
			return
		}
	}

	output, err := ctrl.createLeaveRequestUseCase.Execute(c.Request.Context(), leaveflow.CreateLeaveRequestInput{
		Actor:        buildActor(c),
		LeaveType:    req.LeaveType,
		ApplyDate:    applyDate,
		LateMinutes:  req.LateMinutes,
		EarlyMinutes: req.EarlyMinutes,
		Reason:       req.Reason,
		Documents:    req.Documents,
		ClassID:      req.ClassID,
		LessonID:     req.LessonID,
		Subject:      req.Subject,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to create leave request", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusCreated, "Leave request created successfully", output)
}

func (ctrl *ControllerV1) ApproveLeaveRequest(c *gin.Context) {
	output, err := ctrl.updateLeaveRequestStatusUseCase.Execute(c.Request.Context(), leaveflow.UpdateLeaveRequestStatusInput{
		ID:     c.Param("id"),
		Actor:  buildActor(c),
		Status: "APPROVED",
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to approve leave request", err)
		return
	}
	rest.ResponseSuccess(c, http.StatusOK, "Leave request approved successfully", output)
}

func (ctrl *ControllerV1) RejectLeaveRequest(c *gin.Context) {
	var req RejectLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.updateLeaveRequestStatusUseCase.Execute(c.Request.Context(), leaveflow.UpdateLeaveRequestStatusInput{
		ID:              c.Param("id"),
		Actor:           buildActor(c),
		Status:          "REJECTED",
		RejectionReason: req.RejectionReason,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to reject leave request", err)
		return
	}
	rest.ResponseSuccess(c, http.StatusOK, "Leave request rejected successfully", output)
}

func (ctrl *ControllerV1) CancelLeaveRequest(c *gin.Context) {
	if err := ctrl.cancelLeaveRequestUseCase.Execute(c.Request.Context(), leaveflow.CancelLeaveRequestInput{
		ID:    c.Param("id"),
		Actor: buildActor(c),
	}); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to cancel leave request", err)
		return
	}
	rest.ResponseSuccess(c, http.StatusOK, "Leave request cancelled successfully", nil)
}

func buildActor(c *gin.Context) leaveflow.Actor {
	userRole, _ := c.Get("user_role")
	userEmail, _ := c.Get("user_email")
	userID, _ := c.Get("user_id")

	actor := leaveflow.Actor{}
	if userRole != nil {
		actor.Role, _ = userRole.(string)
	}
	if userEmail != nil {
		actor.Email, _ = userEmail.(string)
	}
	if userID != nil {
		actor.UserID, _ = userID.(string)
	}
	return actor
}
