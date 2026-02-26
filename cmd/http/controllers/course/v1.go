package course

import (
	"doan/cmd/http/rest"
	"doan/internal/entities"
	"doan/internal/usecases/course"
	"doan/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	createCourseUseCase course.CreateCourseUseCase
	getCourseUseCase    course.GetCourseUseCase
	updateCourseUseCase course.UpdateCourseUseCase
	deleteCourseUseCase course.DeleteCourseUseCase
	listCoursesUseCase  course.ListCoursesUseCase
}

func NewCourseControllerV1(
	createCourseUseCase course.CreateCourseUseCase,
	getCourseUseCase course.GetCourseUseCase,
	updateCourseUseCase course.UpdateCourseUseCase,
	deleteCourseUseCase course.DeleteCourseUseCase,
	listCoursesUseCase course.ListCoursesUseCase,
) *ControllerV1 {
	return &ControllerV1{
		createCourseUseCase: createCourseUseCase,
		getCourseUseCase:    getCourseUseCase,
		updateCourseUseCase: updateCourseUseCase,
		deleteCourseUseCase: deleteCourseUseCase,
		listCoursesUseCase:  listCoursesUseCase,
	}
}

// CreateCourse creates a new course
func (c *ControllerV1) CreateCourse(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	var req CreateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := c.createCourseUseCase.Execute(ctx, course.CreateCourseInput{
		Code:                   req.Code,
		Name:                   req.Name,
		Description:            req.Description,
		GradeLevel:             req.GradeLevel,
		Subject:                req.Subject,
		SessionCount:           req.SessionCount,
		SessionDurationMinutes: req.SessionDurationMinutes,
		TotalHours:             req.TotalHours,
		Price:                  req.Price,
		Status:                 req.Status,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to create course: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to create course", err)
		return
	}

	response := mapCourseToResponse(output.Course)
	rest.ResponseSuccess(ctx, http.StatusCreated, "Course created successfully", response)
}

// GetCourse gets a course by ID
func (c *ControllerV1) GetCourse(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	id := ctx.Param("id")
	if id == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Course ID is required", nil)
		return
	}

	output, err := c.getCourseUseCase.Execute(ctx, course.GetCourseInput{ID: id})
	if err != nil {
		ctxLogger.Errorf("Failed to get course: %v", err)
		rest.ResponseError(ctx, http.StatusNotFound, "Course not found", err)
		return
	}

	response := mapCourseToResponse(output.Course)
	rest.ResponseSuccess(ctx, http.StatusOK, "Course retrieved successfully", response)
}

// UpdateCourse updates a course
func (c *ControllerV1) UpdateCourse(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	id := ctx.Param("id")
	if id == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Course ID is required", nil)
		return
	}

	var req UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := c.updateCourseUseCase.Execute(ctx, course.UpdateCourseInput{
		ID:                     id,
		Code:                   req.Code,
		Name:                   req.Name,
		Description:            req.Description,
		GradeLevel:             req.GradeLevel,
		Subject:                req.Subject,
		SessionCount:           req.SessionCount,
		SessionDurationMinutes: req.SessionDurationMinutes,
		TotalHours:             req.TotalHours,
		Price:                  req.Price,
		Status:                 req.Status,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to update course: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to update course", err)
		return
	}

	response := mapCourseToResponse(output.Course)
	rest.ResponseSuccess(ctx, http.StatusOK, "Course updated successfully", response)
}

// DeleteCourse deletes a course
func (c *ControllerV1) DeleteCourse(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	id := ctx.Param("id")
	if id == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Course ID is required", nil)
		return
	}

	output, err := c.deleteCourseUseCase.Execute(ctx, course.DeleteCourseInput{ID: id})
	if err != nil {
		ctxLogger.Errorf("Failed to delete course: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to delete course", err)
		return
	}

	rest.ResponseSuccess(ctx, http.StatusOK, output.Message, MessageResponse{Message: output.Message})
}

// ListCourses lists courses with filters
func (c *ControllerV1) ListCourses(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	search := ctx.Query("search")
	status := ctx.Query("status")
	subject := ctx.Query("subject")
	sortBy := ctx.Query("sort_by")
	sortOrder := ctx.Query("sort_order")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	output, err := c.listCoursesUseCase.Execute(ctx, course.ListCoursesInput{
		Search:    search,
		Status:    status,
		Subject:   subject,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to list courses: %v", err)
		rest.ResponseError(ctx, http.StatusInternalServerError, "Failed to list courses", err)
		return
	}

	courses := make([]CourseResponse, 0, len(output.Courses))
	for _, cr := range output.Courses {
		courses = append(courses, mapCourseToResponse(cr))
	}

	response := ListCoursesResponse{
		Courses: courses,
		Pagination: PaginationMeta{
			ItemsPerPage: output.Pagination.ItemsPerPage,
			TotalItems:   output.Pagination.TotalItems,
			CurrentPage:  output.Pagination.CurrentPage,
			TotalPages:   output.Pagination.TotalPages,
		},
	}

	rest.ResponseSuccess(ctx, http.StatusOK, "Courses retrieved successfully", response)
}

// Helper function to map entity to response
func mapCourseToResponse(c *entities.Course) CourseResponse {
	return CourseResponse{
		ID:                     c.ID,
		Code:                   c.Code,
		Name:                   c.Name,
		Description:            c.Description,
		GradeLevel:             c.GradeLevel,
		Subject:                c.Subject,
		SessionCount:           c.SessionCount,
		SessionDurationMinutes: c.SessionDurationMinutes,
		TotalHours:             c.TotalHours,
		Price:                  c.Price,
		Status:                 c.Status,
		CreatedAt:              c.CreatedAt,
		UpdatedAt:              c.UpdatedAt,
	}
}
