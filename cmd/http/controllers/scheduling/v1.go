package scheduling

import (
	"errors"
	"net/http"
	"strings"
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
	substituteUseCase scheduling.SubstituteUseCase
	makeupUseCase     scheduling.MakeupUseCase
}

func NewSchedulingControllerV1(
	previewUseCase scheduling.PreviewUseCase,
	benchmarkUseCase scheduling.BenchmarkUseCase,
	getPreviewUseCase scheduling.GetPreviewUseCase,
	commitUseCase scheduling.CommitPreviewUseCase,
	substituteUseCase scheduling.SubstituteUseCase,
	makeupUseCase scheduling.MakeupUseCase,
) *ControllerV1 {
	return &ControllerV1{
		previewUseCase:    previewUseCase,
		benchmarkUseCase:  benchmarkUseCase,
		getPreviewUseCase: getPreviewUseCase,
		commitUseCase:     commitUseCase,
		substituteUseCase: substituteUseCase,
		makeupUseCase:     makeupUseCase,
	}
}

func (ctrl *ControllerV1) Preview(c *gin.Context) {
	var req PreviewScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	var dateFrom time.Time
	if strings.TrimSpace(req.DateFrom) != "" {
		var err error
		dateFrom, err = parsePreviewDate(req.DateFrom)
		if err != nil {
			rest.ResponseError(c, http.StatusBadRequest, "Invalid date_from. Expected YYYY-MM-DD or RFC3339", err)
			return
		}
	}

	var dateTo time.Time
	if strings.TrimSpace(req.DateTo) != "" {
		var err error
		dateTo, err = parsePreviewDate(req.DateTo)
		if err != nil {
			rest.ResponseError(c, http.StatusBadRequest, "Invalid date_to. Expected YYYY-MM-DD or RFC3339", err)
			return
		}
	}

	if !dateFrom.IsZero() && !dateTo.IsZero() && dateTo.Before(dateFrom) {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid date range. date_to must be greater than or equal to date_from", nil)
		return
	}

	output, err := ctrl.previewUseCase.Execute(c.Request.Context(), scheduling.PreviewInput{
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		Mode:       req.Mode,
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

func (ctrl *ControllerV1) SuggestSubstitute(c *gin.Context) {
	lessonID := c.Param("id")
	if lessonID == "" {
		rest.ResponseError(c, http.StatusBadRequest, "Lesson ID is required", nil)
		return
	}

	actor := buildSchedulingActor(c)
	suggestions, err := ctrl.substituteUseCase.SuggestSubstituteTeachers(c.Request.Context(), actor, lessonID)
	if err != nil {
		switch {
		case errors.Is(err, scheduling.ErrSubstituteAccessDenied):
			rest.ResponseError(c, http.StatusForbidden, "Access denied for substitute suggestion", err)
		default:
			rest.ResponseError(c, http.StatusInternalServerError, "Failed to suggest substitute teachers", err)
		}
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Substitute teachers suggested successfully", suggestions)
}

func (ctrl *ControllerV1) AssignSubstitute(c *gin.Context) {
	lessonID := c.Param("id")
	if lessonID == "" {
		rest.ResponseError(c, http.StatusBadRequest, "Lesson ID is required", nil)
		return
	}

	var req AssignSubstituteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	actor := buildSchedulingActor(c)
	err := ctrl.substituteUseCase.AssignSubstitute(c.Request.Context(), actor, lessonID, req.NewTeacherID, req.Reason)
	if err != nil {
		switch {
		case errors.Is(err, scheduling.ErrSubstituteAccessDenied):
			rest.ResponseError(c, http.StatusForbidden, "Access denied for substitute assignment", err)
		case errors.Is(err, scheduling.ErrSubstituteNotEligible):
			rest.ResponseError(c, http.StatusBadRequest, "Selected substitute teacher is not eligible", err)
		default:
			rest.ResponseError(c, http.StatusInternalServerError, "Failed to assign substitute teacher", err)
		}
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Substitute teacher assigned successfully", nil)
}

func (ctrl *ControllerV1) FindMakeupSpots(c *gin.Context) {
	lessonID := c.Param("id")
	if lessonID == "" {
		rest.ResponseError(c, http.StatusBadRequest, "Lesson ID is required", nil)
		return
	}

	var req FindMakeupSpotsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	output, err := ctrl.makeupUseCase.FindMakeupSpots(c.Request.Context(), scheduling.FindMakeupSpotsInput{
		LessonID:  lessonID,
		StudentID: req.StudentID,
		Limit:     req.Limit,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to find makeup spots", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Makeup spots found successfully", output)
}

func buildSchedulingActor(c *gin.Context) scheduling.Actor {
	actor := scheduling.Actor{}
	if userID, ok := c.Get("user_id"); ok {
		if value, ok := userID.(string); ok {
			actor.UserID = value
		}
	}
	if userEmail, ok := c.Get("user_email"); ok {
		if value, ok := userEmail.(string); ok {
			actor.Email = value
		}
	}
	if userRole, ok := c.Get("user_role"); ok {
		if value, ok := userRole.(string); ok {
			actor.Role = value
		}
	}
	return actor
}
