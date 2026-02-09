package teacher

import (
	"doan/cmd/http/rest"
	"doan/internal/entities"
	"doan/internal/usecases/teacher"
	"doan/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	createTeacherUseCase         teacher.CreateTeacherUseCase
	getTeacherUseCase            teacher.GetTeacherUseCase
	updateTeacherUseCase         teacher.UpdateTeacherUseCase
	deleteTeacherUseCase         teacher.DeleteTeacherUseCase
	listTeachersUseCase          teacher.ListTeachersUseCase
	getTeacherTimetableUseCase   teacher.GetTeacherTimetableUseCase
	getTeachingHoursStatsUseCase teacher.GetTeachingHoursStatsUseCase
}

func NewTeacherControllerV1(
	createTeacherUseCase teacher.CreateTeacherUseCase,
	getTeacherUseCase teacher.GetTeacherUseCase,
	updateTeacherUseCase teacher.UpdateTeacherUseCase,
	deleteTeacherUseCase teacher.DeleteTeacherUseCase,
	listTeachersUseCase teacher.ListTeachersUseCase,
	getTeacherTimetableUseCase teacher.GetTeacherTimetableUseCase,
	getTeachingHoursStatsUseCase teacher.GetTeachingHoursStatsUseCase,
) *ControllerV1 {
	return &ControllerV1{
		createTeacherUseCase:         createTeacherUseCase,
		getTeacherUseCase:            getTeacherUseCase,
		updateTeacherUseCase:         updateTeacherUseCase,
		deleteTeacherUseCase:         deleteTeacherUseCase,
		listTeachersUseCase:          listTeachersUseCase,
		getTeacherTimetableUseCase:   getTeacherTimetableUseCase,
		getTeachingHoursStatsUseCase: getTeachingHoursStatsUseCase,
	}
}

// CreateTeacher godoc
// @Summary Create a new teacher
// @Description Create a new teacher (Admin only)
// @Tags Teachers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body CreateTeacherRequest true "Teacher data"
// @Success 201 {object} rest.BaseResponse{data=TeacherResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 401 {object} rest.BaseResponse
// @Failure 403 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/teachers [post]
func (c *ControllerV1) CreateTeacher(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	var req CreateTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := c.createTeacherUseCase.Execute(ctx, teacher.CreateTeacherInput{
		Code:            req.Code,
		FullName:        req.FullName,
		Email:           req.Email,
		Phone:           req.Phone,
		IsSchoolTeacher: req.IsSchoolTeacher,
		SchoolName:      req.SchoolName,
		EmploymentType:  req.EmploymentType,
		Status:          req.Status,
		Notes:           req.Notes,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to create teacher: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to create teacher", err)
		return
	}

	response := mapTeacherToResponse(output.Teacher)
	rest.ResponseSuccess(ctx, http.StatusCreated, "Teacher created successfully", response)
}

// GetTeacher godoc
// @Summary Get teacher by ID
// @Description Get teacher details by ID
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path string true "Teacher ID"
// @Success 200 {object} rest.BaseResponse{data=TeacherResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 404 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/teachers/{id} [get]
func (c *ControllerV1) GetTeacher(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	teacherID := ctx.Param("id")
	if teacherID == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Teacher ID is required", nil)
		return
	}

	output, err := c.getTeacherUseCase.Execute(ctx, teacher.GetTeacherInput{
		ID: teacherID,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to get teacher: %v", err)
		rest.ResponseError(ctx, http.StatusNotFound, "Teacher not found", err)
		return
	}

	response := mapTeacherToResponse(output.Teacher)
	rest.ResponseSuccess(ctx, http.StatusOK, "Teacher retrieved successfully", response)
}

// UpdateTeacher godoc
// @Summary Update teacher
// @Description Update teacher information (Admin only)
// @Tags Teachers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Teacher ID"
// @Param payload body UpdateTeacherRequest true "Updated teacher data"
// @Success 200 {object} rest.BaseResponse{data=TeacherResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 401 {object} rest.BaseResponse
// @Failure 403 {object} rest.BaseResponse
// @Failure 404 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/teachers/{id} [put]
func (c *ControllerV1) UpdateTeacher(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	teacherID := ctx.Param("id")
	if teacherID == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Teacher ID is required", nil)
		return
	}

	var req UpdateTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := c.updateTeacherUseCase.Execute(ctx, teacher.UpdateTeacherInput{
		ID:              teacherID,
		Code:            req.Code,
		FullName:        req.FullName,
		Email:           req.Email,
		Phone:           req.Phone,
		IsSchoolTeacher: req.IsSchoolTeacher,
		SchoolName:      req.SchoolName,
		EmploymentType:  req.EmploymentType,
		Status:          req.Status,
		Notes:           req.Notes,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to update teacher: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to update teacher", err)
		return
	}

	response := mapTeacherToResponse(output.Teacher)
	rest.ResponseSuccess(ctx, http.StatusOK, "Teacher updated successfully", response)
}

// DeleteTeacher godoc
// @Summary Delete teacher
// @Description Soft delete a teacher (Admin only)
// @Tags Teachers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Teacher ID"
// @Success 200 {object} rest.BaseResponse{data=MessageResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 401 {object} rest.BaseResponse
// @Failure 403 {object} rest.BaseResponse
// @Failure 404 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/teachers/{id} [delete]
func (c *ControllerV1) DeleteTeacher(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	teacherID := ctx.Param("id")
	if teacherID == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Teacher ID is required", nil)
		return
	}

	output, err := c.deleteTeacherUseCase.Execute(ctx, teacher.DeleteTeacherInput{
		ID: teacherID,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to delete teacher: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to delete teacher", err)
		return
	}

	rest.ResponseSuccess(ctx, http.StatusOK, output.Message, MessageResponse{Message: output.Message})
}

// ListTeachers godoc
// @Summary List teachers
// @Description Get a list of teachers with filtering and pagination
// @Tags Teachers
// @Accept json
// @Produce json
// @Param search query string false "Search in name, email, phone, code"
// @Param status query string false "Filter by status (ACTIVE, INACTIVE)"
// @Param employment_type query string false "Filter by employment type (PART_TIME, FULL_TIME)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param sort_by query string false "Sort field"
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Success 200 {object} rest.BaseResponse{data=ListTeachersResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/teachers [get]
func (c *ControllerV1) ListTeachers(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	// Parse query parameters
	search := ctx.Query("search")
	status := ctx.Query("status")
	employmentType := ctx.Query("employment_type")
	sortBy := ctx.Query("sort_by")
	sortOrder := ctx.Query("sort_order")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	output, err := c.listTeachersUseCase.Execute(ctx, teacher.ListTeachersInput{
		Search:         search,
		Status:         status,
		EmploymentType: employmentType,
		Page:           page,
		Limit:          limit,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to list teachers: %v", err)
		rest.ResponseError(ctx, http.StatusInternalServerError, "Failed to list teachers", err)
		return
	}

	// Map to response
	teachers := make([]TeacherResponse, 0, len(output.Teachers))
	for _, t := range output.Teachers {
		teachers = append(teachers, mapTeacherToResponse(t))
	}

	response := ListTeachersResponse{
		Teachers: teachers,
		Pagination: PaginationMeta{
			ItemsPerPage: output.Pagination.ItemsPerPage,
			TotalItems:   output.Pagination.TotalItems,
			CurrentPage:  output.Pagination.CurrentPage,
			TotalPages:   output.Pagination.TotalPages,
		},
	}

	rest.ResponseSuccess(ctx, http.StatusOK, "Teachers retrieved successfully", response)
}

// GetTeacherTimetable godoc
// @Summary Get teacher timetable
// @Description Get teacher's lesson schedule within a date range
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path string true "Teacher ID"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} rest.BaseResponse{data=TimetableResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 404 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/teachers/{id}/timetable [get]
func (c *ControllerV1) GetTeacherTimetable(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	teacherID := ctx.Param("id")
	if teacherID == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Teacher ID is required", nil)
		return
	}

	// Parse date parameters
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			rest.ResponseError(ctx, http.StatusBadRequest, "Invalid 'from' date format. Use YYYY-MM-DD", err)
			return
		}
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			rest.ResponseError(ctx, http.StatusBadRequest, "Invalid 'to' date format. Use YYYY-MM-DD", err)
			return
		}
	}

	output, err := c.getTeacherTimetableUseCase.Execute(ctx, teacher.GetTeacherTimetableInput{
		TeacherID: teacherID,
		From:      from,
		To:        to,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to get teacher timetable: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to get teacher timetable", err)
		return
	}

	// Map to response
	lessons := make([]TimetableLesson, 0, len(output.Lessons))
	for _, l := range output.Lessons {
		lessons = append(lessons, TimetableLesson{
			ID:        l.ID,
			ClassID:   l.ClassID,
			ClassName: l.ClassName,
			RoomID:    l.RoomID,
			RoomName:  l.RoomName,
			StartTime: l.StartTime,
			EndTime:   l.EndTime,
			Notes:     l.Notes,
		})
	}

	response := TimetableResponse{Lessons: lessons}
	rest.ResponseSuccess(ctx, http.StatusOK, "Timetable retrieved successfully", response)
}

// GetTeachingHoursStats godoc
// @Summary Get teaching hours statistics
// @Description Get teacher's teaching hours statistics grouped by period
// @Tags Teachers
// @Accept json
// @Produce json
// @Param id path string true "Teacher ID"
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Param group_by query string false "Group by period (day, week, month)" default(day)
// @Success 200 {object} rest.BaseResponse{data=TeachingHoursStatsResponse}
// @Failure 400 {object} rest.BaseResponse
// @Failure 404 {object} rest.BaseResponse
// @Failure 500 {object} rest.BaseResponse
// @Router /v1/teachers/{id}/stats/teaching-hours [get]
func (c *ControllerV1) GetTeachingHoursStats(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	teacherID := ctx.Param("id")
	if teacherID == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Teacher ID is required", nil)
		return
	}

	// Parse parameters
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	groupBy := ctx.DefaultQuery("group_by", "day")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			rest.ResponseError(ctx, http.StatusBadRequest, "Invalid 'from' date format. Use YYYY-MM-DD", err)
			return
		}
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			rest.ResponseError(ctx, http.StatusBadRequest, "Invalid 'to' date format. Use YYYY-MM-DD", err)
			return
		}
	}

	output, err := c.getTeachingHoursStatsUseCase.Execute(ctx, teacher.GetTeachingHoursStatsInput{
		TeacherID: teacherID,
		From:      from,
		To:        to,
		GroupBy:   groupBy,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to get teaching hours stats: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to get teaching hours stats", err)
		return
	}

	// Map to response
	breakdown := make([]TeachingHoursStat, 0, len(output.Breakdown))
	for _, stat := range output.Breakdown {
		breakdown = append(breakdown, TeachingHoursStat{
			Period: stat.Period,
			Hours:  stat.Hours,
		})
	}

	response := TeachingHoursStatsResponse{
		TotalHours: output.TotalHours,
		Breakdown:  breakdown,
	}

	rest.ResponseSuccess(ctx, http.StatusOK, "Teaching hours stats retrieved successfully", response)
}

// Helper function to map entity to response
func mapTeacherToResponse(t *entities.Teacher) TeacherResponse {
	return TeacherResponse{
		ID:              t.ID,
		Code:            t.Code,
		FullName:        t.FullName,
		Email:           t.Email,
		Phone:           t.Phone,
		IsSchoolTeacher: t.IsSchoolTeacher,
		SchoolName:      t.SchoolName,
		EmploymentType:  t.EmploymentType,
		Status:          t.Status,
		Notes:           t.Notes,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}
