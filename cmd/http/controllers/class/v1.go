package class

import (
	"doan/cmd/http/rest"
	"doan/internal/usecases/class"
	"net/http"
	"strconv"

	"doan/pkg/logger"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	createClassUseCase class.CreateClassUseCase
	getClassUseCase    class.GetClassUseCase
	updateClassUseCase class.UpdateClassUseCase
	deleteClassUseCase class.DeleteClassUseCase
	listClassesUseCase class.ListClassesUseCase
}

func NewClassControllerV1(
	createClassUseCase class.CreateClassUseCase,
	getClassUseCase class.GetClassUseCase,
	updateClassUseCase class.UpdateClassUseCase,
	deleteClassUseCase class.DeleteClassUseCase,
	listClassesUseCase class.ListClassesUseCase,
) *ControllerV1 {
	return &ControllerV1{
		createClassUseCase: createClassUseCase,
		getClassUseCase:    getClassUseCase,
		updateClassUseCase: updateClassUseCase,
		deleteClassUseCase: deleteClassUseCase,
		listClassesUseCase: listClassesUseCase,
	}
}

func (ctrl *ControllerV1) CreateClass(c *gin.Context) {
	ctxLogger := logger.NewLogger(c.Request.Context())
	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.createClassUseCase.Execute(c.Request.Context(), class.CreateClassInput{
		Code:        req.Code,
		Name:        req.Name,
		Notes:       req.Notes,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		MaxStudents: req.MaxStudents,
		Status:      req.Status,
		Price:       req.Price,
		ProgramID:   req.ProgramID,
		CourseID:    req.CourseID,
		TeacherID:   req.TeacherID,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to create class", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusCreated, "Class created successfully", output.Class)
}

func (ctrl *ControllerV1) GetClass(c *gin.Context) {
	id := c.Param("id")
	output, err := ctrl.getClassUseCase.Execute(c.Request.Context(), class.GetClassInput{ID: id})
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Class not found", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Class retrieved successfully", output.Class)
}

func (ctrl *ControllerV1) UpdateClass(c *gin.Context) {
	id := c.Param("id")
	var req UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.updateClassUseCase.Execute(c.Request.Context(), class.UpdateClassInput{
		ID:          id,
		Code:        req.Code,
		Name:        req.Name,
		Notes:       req.Notes,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		MaxStudents: req.MaxStudents,
		Status:      req.Status,
		Price:       req.Price,
		ProgramID:   req.ProgramID,
		CourseID:    req.CourseID,
		TeacherID:   req.TeacherID,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to update class", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Class updated successfully", output.Class)
}

func (ctrl *ControllerV1) DeleteClass(c *gin.Context) {
	id := c.Param("id")

	output, err := ctrl.deleteClassUseCase.Execute(c.Request.Context(), class.DeleteClassInput{ID: id})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to delete class", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, output.Message, nil)
}

func (ctrl *ControllerV1) ListClasses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")

	output, err := ctrl.listClassesUseCase.Execute(c.Request.Context(), class.ListClassesInput{
		Search:    search,
		Status:    status,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to list classes", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Classes retrieved successfully", output)
}
