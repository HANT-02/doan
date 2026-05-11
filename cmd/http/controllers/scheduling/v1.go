package scheduling

import (
	"errors"
	"net/http"
	"time"

	"doan/cmd/http/rest"
	"doan/internal/usecases/scheduling"
	"doan/pkg/utils"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	previewUseCase    scheduling.PreviewUseCase
	benchmarkUseCase  scheduling.BenchmarkUseCase
	getPreviewUseCase scheduling.GetPreviewUseCase
	commitUseCase     scheduling.CommitPreviewUseCase
}

func NewSchedulingControllerV1(
	previewUseCase scheduling.PreviewUseCase,
	benchmarkUseCase scheduling.BenchmarkUseCase,
	getPreviewUseCase scheduling.GetPreviewUseCase,
	commitUseCase scheduling.CommitPreviewUseCase,
) *ControllerV1 {
	return &ControllerV1{
		previewUseCase:    previewUseCase,
		benchmarkUseCase:  benchmarkUseCase,
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

	dateFrom, err := parsePreviewDate(req.DateFrom)
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid date_from. Expected YYYY-MM-DD or RFC3339", err)
		return
	}

	dateTo, err := parsePreviewDate(req.DateTo)
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid date_to. Expected YYYY-MM-DD or RFC3339", err)
		return
	}

	if dateTo.Before(dateFrom) {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid date range. date_to must be greater than or equal to date_from", nil)
		return
	}

	output, err := ctrl.previewUseCase.Execute(c.Request.Context(), scheduling.PreviewInput{
		DateFrom:   dateFrom,
		DateTo:     dateTo,
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

func (ctrl *ControllerV1) Benchmark(c *gin.Context) {
	var req PreviewScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	dateFrom, err := parsePreviewDate(req.DateFrom)
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid date_from. Expected YYYY-MM-DD or RFC3339", err)
		return
	}

	dateTo, err := parsePreviewDate(req.DateTo)
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid date_to. Expected YYYY-MM-DD or RFC3339", err)
		return
	}

	if dateTo.Before(dateFrom) {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid date range. date_to must be greater than or equal to date_from", nil)
		return
	}

	output, err := ctrl.benchmarkUseCase.Execute(c.Request.Context(), scheduling.BenchmarkInput{
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		ClassIDs:   req.ClassIDs,
		TeacherIDs: req.TeacherIDs,
		RoomIDs:    req.RoomIDs,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to execute scheduling benchmark", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Scheduling benchmark executed successfully", output)
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
		RunID:             req.RunID,
		ManualAssignments: mapManualAssignments(req.ManualAssignments),
	})
	if err != nil {
		var conflictErr *scheduling.CommitPreviewConflictError
		if errors.As(err, &conflictErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, &rest.BaseResponse{
				Success:   false,
				Message:   utils.NewStringPtr("Failed to commit scheduling preview"),
				ErrorCode: utils.NewStringPtr(conflictErr.Error()),
				Data:      conflictErr.Preview,
			})
			return
		}
		rest.ResponseError(c, http.StatusBadRequest, "Failed to commit scheduling preview", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Scheduling preview committed successfully", output)
}

func mapManualAssignments(items []CommitAssignmentChoice) []scheduling.ManualAssignmentOverride {
	if len(items) == 0 {
		return nil
	}

	overrides := make([]scheduling.ManualAssignmentOverride, 0, len(items))
	for _, item := range items {
		if item.VariableID == "" || item.OptionKey == "" {
			continue
		}

		overrides = append(overrides, scheduling.ManualAssignmentOverride{
			VariableID: item.VariableID,
			OptionKey:  item.OptionKey,
		})
	}

	return overrides
}

func parsePreviewDate(raw string) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		location = time.FixedZone("ICT", 7*60*60)
	}

	if parsedDate, parseErr := time.ParseInLocation("2006-01-02", raw, location); parseErr == nil {
		return parsedDate, nil
	}

	return time.Parse(time.RFC3339, raw)
}
