package program

import (
	"doan/cmd/http/controllers/course"
	"doan/cmd/http/rest"
	"doan/internal/entities"
	"doan/internal/usecases/program"
	"doan/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	createProgramUseCase program.CreateProgramUseCase
	getProgramUseCase    program.GetProgramUseCase
	updateProgramUseCase program.UpdateProgramUseCase
	deleteProgramUseCase program.DeleteProgramUseCase
	listProgramsUseCase  program.ListProgramsUseCase
	addCoursesUseCase    program.AddCoursesUseCase
	removeCoursesUseCase program.RemoveCoursesUseCase
}

func NewProgramControllerV1(
	createProgramUseCase program.CreateProgramUseCase,
	getProgramUseCase program.GetProgramUseCase,
	updateProgramUseCase program.UpdateProgramUseCase,
	deleteProgramUseCase program.DeleteProgramUseCase,
	listProgramsUseCase program.ListProgramsUseCase,
	addCoursesUseCase program.AddCoursesUseCase,
	removeCoursesUseCase program.RemoveCoursesUseCase,
) *ControllerV1 {
	return &ControllerV1{
		createProgramUseCase: createProgramUseCase,
		getProgramUseCase:    getProgramUseCase,
		updateProgramUseCase: updateProgramUseCase,
		deleteProgramUseCase: deleteProgramUseCase,
		listProgramsUseCase:  listProgramsUseCase,
		addCoursesUseCase:    addCoursesUseCase,
		removeCoursesUseCase: removeCoursesUseCase,
	}
}

// CreateProgram creates a new program
func (c *ControllerV1) CreateProgram(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	var req CreateProgramRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := c.createProgramUseCase.Execute(ctx, program.CreateProgramInput{
		Code:          req.Code,
		Name:          req.Name,
		Track:         req.Track,
		EffectiveFrom: req.EffectiveFrom,
		EffectiveTo:   req.EffectiveTo,
		ApprovalNote:  req.ApprovalNote,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to create program: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to create program", err)
		return
	}

	response := mapProgramToResponse(output.Program)
	rest.ResponseSuccess(ctx, http.StatusCreated, "Program created successfully", response)
}

// GetProgram gets a program by ID
func (c *ControllerV1) GetProgram(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	id := ctx.Param("id")
	if id == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Program ID is required", nil)
		return
	}

	output, err := c.getProgramUseCase.Execute(ctx, program.GetProgramInput{ID: id})
	if err != nil {
		ctxLogger.Errorf("Failed to get program: %v", err)
		rest.ResponseError(ctx, http.StatusNotFound, "Program not found", err)
		return
	}

	response := mapProgramToResponse(output.Program)
	rest.ResponseSuccess(ctx, http.StatusOK, "Program retrieved successfully", response)
}

// UpdateProgram updates a program
func (c *ControllerV1) UpdateProgram(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	id := ctx.Param("id")
	if id == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Program ID is required", nil)
		return
	}

	var req UpdateProgramRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := c.updateProgramUseCase.Execute(ctx, program.UpdateProgramInput{
		ID:            id,
		Code:          req.Code,
		Name:          req.Name,
		Track:         req.Track,
		EffectiveFrom: req.EffectiveFrom,
		EffectiveTo:   req.EffectiveTo,
		ApprovalNote:  req.ApprovalNote,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to update program: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to update program", err)
		return
	}

	response := mapProgramToResponse(output.Program)
	rest.ResponseSuccess(ctx, http.StatusOK, "Program updated successfully", response)
}

// DeleteProgram deletes a program
func (c *ControllerV1) DeleteProgram(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	id := ctx.Param("id")
	if id == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Program ID is required", nil)
		return
	}

	output, err := c.deleteProgramUseCase.Execute(ctx, program.DeleteProgramInput{ID: id})
	if err != nil {
		ctxLogger.Errorf("Failed to delete program: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to delete program", err)
		return
	}

	rest.ResponseSuccess(ctx, http.StatusOK, output.Message, MessageResponse{Message: output.Message})
}

// ListPrograms lists programs with filters
func (c *ControllerV1) ListPrograms(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)

	search := ctx.Query("search")
	track := ctx.Query("track")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	output, err := c.listProgramsUseCase.Execute(ctx, program.ListProgramsInput{
		Search: search,
		Track:  track,
		Page:   page,
		Limit:  limit,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to list programs: %v", err)
		rest.ResponseError(ctx, http.StatusInternalServerError, "Failed to list programs", err)
		return
	}

	programs := make([]ProgramResponse, 0, len(output.Programs))
	for _, pr := range output.Programs {
		programs = append(programs, mapProgramToResponse(pr))
	}

	response := ListProgramsResponse{
		Programs: programs,
		Pagination: PaginationMeta{
			ItemsPerPage: output.Pagination.ItemsPerPage,
			TotalItems:   output.Pagination.TotalItems,
			CurrentPage:  output.Pagination.CurrentPage,
			TotalPages:   output.Pagination.TotalPages,
		},
	}

	rest.ResponseSuccess(ctx, http.StatusOK, "Programs retrieved successfully", response)
}

func (c *ControllerV1) AddCourses(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)
	id := ctx.Param("id")
	if id == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Program ID is required", nil)
		return
	}

	var req AddRemoveCoursesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := c.addCoursesUseCase.Execute(ctx, program.AddCoursesInput{
		ProgramID: id,
		CourseIDs: req.CourseIDs,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to add courses: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to add courses", err)
		return
	}

	rest.ResponseSuccess(ctx, http.StatusOK, output.Message, MessageResponse{Message: output.Message})
}

func (c *ControllerV1) RemoveCourses(ctx *gin.Context) {
	ctxLogger := logger.NewLogger(ctx)
	id := ctx.Param("id")
	if id == "" {
		rest.ResponseError(ctx, http.StatusBadRequest, "Program ID is required", nil)
		return
	}

	var req AddRemoveCoursesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := c.removeCoursesUseCase.Execute(ctx, program.RemoveCoursesInput{
		ProgramID: id,
		CourseIDs: req.CourseIDs,
	})

	if err != nil {
		ctxLogger.Errorf("Failed to remove courses: %v", err)
		rest.ResponseError(ctx, http.StatusBadRequest, "Failed to remove courses", err)
		return
	}

	rest.ResponseSuccess(ctx, http.StatusOK, output.Message, MessageResponse{Message: output.Message})
}

// Helper function to map entity to response
func mapProgramToResponse(p *entities.Program) ProgramResponse {
	resp := ProgramResponse{
		ID:            p.ID,
		Code:          p.Code,
		Name:          p.Name,
		Track:         p.Track,
		EffectiveFrom: p.EffectiveFrom,
		EffectiveTo:   p.EffectiveTo,
		ApprovalNote:  p.ApprovalNote,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}

	if len(p.Courses) > 0 {
		var courseResps []course.CourseResponse
		for _, c := range p.Courses {
			courseResps = append(courseResps, course.CourseResponse{
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
			})
		}
		resp.Courses = courseResps
	}

	return resp
}
