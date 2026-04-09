package shift

import (
	"net/http"
	"strconv"

	"doan/cmd/http/rest"
	"doan/internal/entities"
	shiftusecase "doan/internal/usecases/shift"
	"doan/pkg/logger"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	createShiftUseCase shiftusecase.CreateShiftUseCase
	getShiftUseCase    shiftusecase.GetShiftUseCase
	updateShiftUseCase shiftusecase.UpdateShiftUseCase
	deleteShiftUseCase shiftusecase.DeleteShiftUseCase
	listShiftsUseCase  shiftusecase.ListShiftsUseCase
}

func NewShiftControllerV1(
	createShiftUseCase shiftusecase.CreateShiftUseCase,
	getShiftUseCase shiftusecase.GetShiftUseCase,
	updateShiftUseCase shiftusecase.UpdateShiftUseCase,
	deleteShiftUseCase shiftusecase.DeleteShiftUseCase,
	listShiftsUseCase shiftusecase.ListShiftsUseCase,
) *ControllerV1 {
	return &ControllerV1{
		createShiftUseCase: createShiftUseCase,
		getShiftUseCase:    getShiftUseCase,
		updateShiftUseCase: updateShiftUseCase,
		deleteShiftUseCase: deleteShiftUseCase,
		listShiftsUseCase:  listShiftsUseCase,
	}
}

func (ctrl *ControllerV1) CreateShift(c *gin.Context) {
	ctxLogger := logger.NewLogger(c.Request.Context())
	var req CreateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind shift create request: %v", err)
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.createShiftUseCase.Execute(c.Request.Context(), shiftusecase.CreateShiftInput{
		Code:            req.Code,
		Name:            req.Name,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DurationMinutes: req.DurationMinutes,
		SessionType:     req.SessionType,
		IsActive:        req.IsActive,
		Notes:           req.Notes,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to create shift", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusCreated, "Shift created successfully", mapShiftToResponse(output.Shift))
}

func (ctrl *ControllerV1) GetShift(c *gin.Context) {
	output, err := ctrl.getShiftUseCase.Execute(c.Request.Context(), shiftusecase.GetShiftInput{
		ID: c.Param("id"),
	})
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Shift not found", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Shift retrieved successfully", mapShiftToResponse(output.Shift))
}

func (ctrl *ControllerV1) UpdateShift(c *gin.Context) {
	var req UpdateShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.updateShiftUseCase.Execute(c.Request.Context(), shiftusecase.UpdateShiftInput{
		ID:              c.Param("id"),
		Code:            req.Code,
		Name:            req.Name,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DurationMinutes: req.DurationMinutes,
		SessionType:     req.SessionType,
		IsActive:        req.IsActive,
		Notes:           req.Notes,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to update shift", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Shift updated successfully", mapShiftToResponse(output.Shift))
}

func (ctrl *ControllerV1) DeleteShift(c *gin.Context) {
	output, err := ctrl.deleteShiftUseCase.Execute(c.Request.Context(), shiftusecase.DeleteShiftInput{
		ID: c.Param("id"),
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to delete shift", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, output.Message, nil)
}

func (ctrl *ControllerV1) ListShifts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	sessionType := c.Query("session_type")
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")

	var isActive *bool
	if raw := c.Query("is_active"); raw != "" {
		parsed := raw == "true" || raw == "1" || raw == "TRUE"
		isActive = &parsed
	}

	output, err := ctrl.listShiftsUseCase.Execute(c.Request.Context(), shiftusecase.ListShiftsInput{
		Search:      search,
		SessionType: sessionType,
		IsActive:    isActive,
		Page:        page,
		Limit:       limit,
		SortBy:      sortBy,
		SortOrder:   sortOrder,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to list shifts", err)
		return
	}

	items := make([]ShiftResponse, 0, len(output.Shifts))
	for index := range output.Shifts {
		items = append(items, mapShiftToResponse(&output.Shifts[index]))
	}

	rest.ResponseSuccess(c, http.StatusOK, "Shifts retrieved successfully", ListShiftsResponse{
		Shifts: items,
		Pagination: PaginationMeta{
			ItemsPerPage: output.Pagination.ItemsPerPage,
			TotalItems:   output.Pagination.TotalItems,
			CurrentPage:  output.Pagination.CurrentPage,
			TotalPages:   output.Pagination.TotalPages,
		},
	})
}

func mapShiftToResponse(shiftEntity *entities.Shift) ShiftResponse {
	return ShiftResponse{
		ID:              shiftEntity.ID,
		Code:            shiftEntity.Code,
		Name:            shiftEntity.Name,
		StartTime:       shiftEntity.StartTime,
		EndTime:         shiftEntity.EndTime,
		DurationMinutes: shiftEntity.DurationMinutes,
		SessionType:     shiftEntity.SessionType,
		IsActive:        shiftEntity.IsActive,
		Notes:           shiftEntity.Notes,
		CreatedAt:       shiftEntity.CreatedAt,
		UpdatedAt:       shiftEntity.UpdatedAt,
	}
}
