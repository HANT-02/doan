package scheduling

import (
	"net/http"

	"doan/cmd/http/rest"
	"doan/internal/usecases/scheduling"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	previewUseCase    scheduling.PreviewUseCase
	getPreviewUseCase scheduling.GetPreviewUseCase
	commitUseCase     scheduling.CommitPreviewUseCase
}

func NewSchedulingControllerV1(
	previewUseCase scheduling.PreviewUseCase,
	getPreviewUseCase scheduling.GetPreviewUseCase,
	commitUseCase scheduling.CommitPreviewUseCase,
) *ControllerV1 {
	return &ControllerV1{
		previewUseCase:    previewUseCase,
		getPreviewUseCase: getPreviewUseCase,
		commitUseCase:     commitUseCase,
	}
}

func (ctrl *ControllerV1) Preview(c *gin.Context) {
	var req PreviewScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.previewUseCase.Execute(c.Request.Context(), scheduling.PreviewInput{
		DateFrom:   req.DateFrom,
		DateTo:     req.DateTo,
		ClassIDs:   req.ClassIDs,
		TeacherIDs: req.TeacherIDs,
		RoomIDs:    req.RoomIDs,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to generate scheduling preview", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Scheduling preview generated successfully", output)
}

func (ctrl *ControllerV1) GetPreview(c *gin.Context) {
	output, err := ctrl.getPreviewUseCase.Execute(c.Request.Context(), scheduling.GetPreviewInput{
		RunID: c.Param("id"),
	})
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Scheduling preview not found", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Scheduling preview retrieved successfully", output)
}

func (ctrl *ControllerV1) GetLatestPreview(c *gin.Context) {
	output, err := ctrl.getPreviewUseCase.GetLatest(c.Request.Context())
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Scheduling preview not found", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Latest scheduling preview retrieved successfully", output)
}

func (ctrl *ControllerV1) Commit(c *gin.Context) {
	var req CommitScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.commitUseCase.Execute(c.Request.Context(), scheduling.CommitPreviewInput{
		RunID: req.RunID,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to commit scheduling preview", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Scheduling commit scaffold executed successfully", output)
}
