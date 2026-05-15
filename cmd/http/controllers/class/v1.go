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
	createClassUseCase         class.CreateClassUseCase
	getClassUseCase            class.GetClassUseCase
	getClassRosterUseCase      class.GetClassRosterUseCase
	updateClassUseCase         class.UpdateClassUseCase
	deleteClassUseCase         class.DeleteClassUseCase
	listClassesUseCase         class.ListClassesUseCase
	enrollStudentsUseCase      class.EnrollStudentsUseCase
	removeStudentsUseCase      class.RemoveStudentsUseCase
	reserveStudentUseCase      class.ReserveStudentUseCase
	resumeStudentUseCase       class.ResumeStudentUseCase
	transferStudentUseCase     class.TransferStudentUseCase
	assignTeacherUseCase       class.AssignTeacherUseCase
	getClassSchedulesUseCase   class.GetClassSchedulesUseCase
	createClassScheduleUseCase class.CreateClassScheduleUseCase
	deleteClassScheduleUseCase class.DeleteClassScheduleUseCase
}

func NewClassControllerV1(
	createClassUseCase class.CreateClassUseCase,
	getClassUseCase class.GetClassUseCase,
	getClassRosterUseCase class.GetClassRosterUseCase,
	updateClassUseCase class.UpdateClassUseCase,
	deleteClassUseCase class.DeleteClassUseCase,
	listClassesUseCase class.ListClassesUseCase,
	enrollStudentsUseCase class.EnrollStudentsUseCase,
	removeStudentsUseCase class.RemoveStudentsUseCase,
	reserveStudentUseCase class.ReserveStudentUseCase,
	resumeStudentUseCase class.ResumeStudentUseCase,
	transferStudentUseCase class.TransferStudentUseCase,
	assignTeacherUseCase class.AssignTeacherUseCase,
	getClassSchedulesUseCase class.GetClassSchedulesUseCase,
	createClassScheduleUseCase class.CreateClassScheduleUseCase,
	deleteClassScheduleUseCase class.DeleteClassScheduleUseCase,
) *ControllerV1 {
	return &ControllerV1{
		createClassUseCase:         createClassUseCase,
		getClassUseCase:            getClassUseCase,
		getClassRosterUseCase:      getClassRosterUseCase,
		updateClassUseCase:         updateClassUseCase,
		deleteClassUseCase:         deleteClassUseCase,
		listClassesUseCase:         listClassesUseCase,
		enrollStudentsUseCase:      enrollStudentsUseCase,
		removeStudentsUseCase:      removeStudentsUseCase,
		reserveStudentUseCase:      reserveStudentUseCase,
		resumeStudentUseCase:       resumeStudentUseCase,
		transferStudentUseCase:     transferStudentUseCase,
		assignTeacherUseCase:       assignTeacherUseCase,
		getClassSchedulesUseCase:   getClassSchedulesUseCase,
		createClassScheduleUseCase: createClassScheduleUseCase,
		deleteClassScheduleUseCase: deleteClassScheduleUseCase,
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
		rest.ResponseError(c, http.StatusBadRequest, "Failed to create class", err)
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

func (ctrl *ControllerV1) GetClassRoster(c *gin.Context) {
	id := c.Param("id")
	output, err := ctrl.getClassRosterUseCase.Execute(c.Request.Context(), class.GetClassRosterInput{ClassID: id})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to retrieve class roster", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Class roster retrieved successfully", output)
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
		rest.ResponseError(c, http.StatusBadRequest, "Failed to update class", err)
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

func (ctrl *ControllerV1) EnrollStudents(c *gin.Context) {
	id := c.Param("id")
	var req EnrollStudentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.enrollStudentsUseCase.Execute(c.Request.Context(), class.EnrollStudentsInput{
		ClassID:    id,
		StudentIDs: req.StudentIDs,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to enroll students", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Students enrolled successfully", output)
}

func (ctrl *ControllerV1) RemoveStudents(c *gin.Context) {
	id := c.Param("id")
	var req RemoveStudentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.removeStudentsUseCase.Execute(c.Request.Context(), class.RemoveStudentsInput{
		ClassID:    id,
		StudentIDs: req.StudentIDs,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to remove students", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, output.Message, nil)
}

func (ctrl *ControllerV1) ReserveStudent(c *gin.Context) {
	id := c.Param("id")
	studentID := c.Param("studentId")
	var req ReserveStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.reserveStudentUseCase.Execute(c.Request.Context(), class.ReserveStudentInput{
		ClassID:     id,
		StudentID:   studentID,
		Reason:      req.Reason,
		EffectiveAt: req.EffectiveAt,
	})
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == class.ErrEnrollmentNotFound {
			statusCode = http.StatusNotFound
		}
		rest.ResponseError(c, statusCode, "Failed to reserve student", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, output.Message, output)
}

func (ctrl *ControllerV1) ResumeStudent(c *gin.Context) {
	id := c.Param("id")
	studentID := c.Param("studentId")
	var req ResumeStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.resumeStudentUseCase.Execute(c.Request.Context(), class.ResumeStudentInput{
		ClassID:     id,
		StudentID:   studentID,
		Reason:      req.Reason,
		EffectiveAt: req.EffectiveAt,
	})
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == class.ErrReservedEnrollmentNotFound {
			statusCode = http.StatusNotFound
		}
		rest.ResponseError(c, statusCode, "Failed to resume student", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, output.Message, output)
}

func (ctrl *ControllerV1) TransferStudent(c *gin.Context) {
	id := c.Param("id")
	studentID := c.Param("studentId")
	var req TransferStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.transferStudentUseCase.Execute(c.Request.Context(), class.TransferStudentInput{
		SourceClassID: id,
		TargetClassID: req.TargetClassID,
		StudentID:     studentID,
		Reason:        req.Reason,
		EffectiveAt:   req.EffectiveAt,
	})
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == class.ErrEnrollmentNotFound {
			statusCode = http.StatusNotFound
		}
		rest.ResponseError(c, statusCode, "Failed to transfer student", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, output.Message, output)
}

func (ctrl *ControllerV1) AssignTeacher(c *gin.Context) {
	id := c.Param("id")
	var req AssignTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.assignTeacherUseCase.Execute(c.Request.Context(), class.AssignTeacherInput{
		ClassID:   id,
		TeacherID: req.TeacherID,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to assign teacher", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, output.Message, nil)
}

func (ctrl *ControllerV1) GetClassSchedules(c *gin.Context) {
	id := c.Param("id")
	output, err := ctrl.getClassSchedulesUseCase.Execute(c.Request.Context(), class.GetClassSchedulesInput{ClassID: id})
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Failed to retrieve schedules", err)
		return
	}
	rest.ResponseSuccess(c, http.StatusOK, "Schedules retrieved successfully", output)
}

func (ctrl *ControllerV1) CreateClassSchedule(c *gin.Context) {
	id := c.Param("id")
	var req CreateClassScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.createClassScheduleUseCase.Execute(c.Request.Context(), class.CreateClassScheduleInput{
		ClassID:   id,
		ShiftID:   req.ShiftID,
		DayOfWeek: req.DayOfWeek,
		RoomID:    req.RoomID,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to create schedule", err)
		return
	}
	rest.ResponseSuccess(c, http.StatusCreated, "Schedule created successfully", output)
}

func (ctrl *ControllerV1) DeleteClassSchedule(c *gin.Context) {
	id := c.Param("id")
	scheduleId := c.Param("scheduleId")

	output, err := ctrl.deleteClassScheduleUseCase.Execute(c.Request.Context(), class.DeleteClassScheduleInput{
		ClassID:    id,
		ScheduleID: scheduleId,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to delete schedule", err)
		return
	}
	rest.ResponseSuccess(c, http.StatusOK, output.Message, nil)
}
