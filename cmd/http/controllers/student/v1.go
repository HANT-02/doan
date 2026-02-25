package student

import (
	"doan/cmd/http/rest"
	"doan/internal/usecases/student"
	"net/http"
	"strconv"

	"doan/pkg/logger"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	createStudentUseCase student.CreateStudentUseCase
	getStudentUseCase    student.GetStudentUseCase
	updateStudentUseCase student.UpdateStudentUseCase
	deleteStudentUseCase student.DeleteStudentUseCase
	listStudentsUseCase  student.ListStudentsUseCase
}

func NewStudentControllerV1(
	createStudentUseCase student.CreateStudentUseCase,
	getStudentUseCase student.GetStudentUseCase,
	updateStudentUseCase student.UpdateStudentUseCase,
	deleteStudentUseCase student.DeleteStudentUseCase,
	listStudentsUseCase student.ListStudentsUseCase,
) *ControllerV1 {
	return &ControllerV1{
		createStudentUseCase: createStudentUseCase,
		getStudentUseCase:    getStudentUseCase,
		updateStudentUseCase: updateStudentUseCase,
		deleteStudentUseCase: deleteStudentUseCase,
		listStudentsUseCase:  listStudentsUseCase,
	}
}

func (ctrl *ControllerV1) CreateStudent(c *gin.Context) {
	ctxLogger := logger.NewLogger(c.Request.Context())
	var req CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.createStudentUseCase.Execute(c.Request.Context(), student.CreateStudentInput{
		Code:          req.Code,
		FullName:      req.FullName,
		Email:         req.Email,
		Phone:         req.Phone,
		GuardianPhone: req.GuardianPhone,
		GradeLevel:    req.GradeLevel,
		Status:        req.Status,
		DateOfBirth:   req.DateOfBirth,
		Gender:        req.Gender,
		Address:       req.Address,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to create student", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusCreated, "Student created successfully", output.Student)
}

func (ctrl *ControllerV1) GetStudent(c *gin.Context) {
	id := c.Param("id")
	output, err := ctrl.getStudentUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Student not found", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Student retrieved successfully", output.Student)
}

func (ctrl *ControllerV1) UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var req UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.updateStudentUseCase.Execute(c.Request.Context(), student.UpdateStudentInput{
		ID:            id,
		Code:          req.Code,
		FullName:      req.FullName,
		Email:         req.Email,
		Phone:         req.Phone,
		GuardianPhone: req.GuardianPhone,
		GradeLevel:    req.GradeLevel,
		Status:        req.Status,
		DateOfBirth:   req.DateOfBirth,
		Gender:        req.Gender,
		Address:       req.Address,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to update student", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Student updated successfully", output.Student)
}

func (ctrl *ControllerV1) DeleteStudent(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.deleteStudentUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to delete student", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Student deleted successfully", nil)
}

func (ctrl *ControllerV1) ListStudents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")

	output, err := ctrl.listStudentsUseCase.Execute(c.Request.Context(), student.ListStudentsInput{
		Search:    search,
		Status:    status,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to list students", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Students retrieved successfully", output)
}
