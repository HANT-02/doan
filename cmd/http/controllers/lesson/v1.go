package lesson

import (
	"doan/cmd/http/rest"
	"doan/internal/usecases/lesson"
	lessonactivity "doan/internal/usecases/lessonactivity"
	lessonrecord "doan/internal/usecases/lessonrecord"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	listLessonsUseCase                   lesson.ListLessonsUseCase
	getLessonUseCase                     lesson.GetLessonUseCase
	getLessonAttendanceUseCase           lessonactivity.GetLessonAttendanceUseCase
	upsertLessonAttendanceUseCase        lessonactivity.UpsertLessonAttendanceUseCase
	getLessonSummaryUseCase              lessonactivity.GetLessonSummaryUseCase
	upsertLessonSummaryUseCase           lessonactivity.UpsertLessonSummaryUseCase
	getLessonAcademicRecordsUseCase      lessonrecord.GetLessonAcademicRecordsUseCase
	upsertLessonAcademicRecordsUseCase   lessonrecord.UpsertLessonAcademicRecordsUseCase
	finalizeLessonAcademicRecordsUseCase lessonrecord.FinalizeLessonAcademicRecordsUseCase
}

func NewLessonControllerV1(
	listLessonsUseCase lesson.ListLessonsUseCase,
	getLessonUseCase lesson.GetLessonUseCase,
	getLessonAttendanceUseCase lessonactivity.GetLessonAttendanceUseCase,
	upsertLessonAttendanceUseCase lessonactivity.UpsertLessonAttendanceUseCase,
	getLessonSummaryUseCase lessonactivity.GetLessonSummaryUseCase,
	upsertLessonSummaryUseCase lessonactivity.UpsertLessonSummaryUseCase,
	getLessonAcademicRecordsUseCase lessonrecord.GetLessonAcademicRecordsUseCase,
	upsertLessonAcademicRecordsUseCase lessonrecord.UpsertLessonAcademicRecordsUseCase,
	finalizeLessonAcademicRecordsUseCase lessonrecord.FinalizeLessonAcademicRecordsUseCase,
) *ControllerV1 {
	return &ControllerV1{
		listLessonsUseCase:                   listLessonsUseCase,
		getLessonUseCase:                     getLessonUseCase,
		getLessonAttendanceUseCase:           getLessonAttendanceUseCase,
		upsertLessonAttendanceUseCase:        upsertLessonAttendanceUseCase,
		getLessonSummaryUseCase:              getLessonSummaryUseCase,
		upsertLessonSummaryUseCase:           upsertLessonSummaryUseCase,
		getLessonAcademicRecordsUseCase:      getLessonAcademicRecordsUseCase,
		upsertLessonAcademicRecordsUseCase:   upsertLessonAcademicRecordsUseCase,
		finalizeLessonAcademicRecordsUseCase: finalizeLessonAcademicRecordsUseCase,
	}
}

func (ctrl *ControllerV1) ListLessons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	classID := c.Query("class_id")
	teacherID := c.Query("teacher_id")
	status := c.Query("status")
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")

	var dateFrom, dateTo *time.Time
	if dateFromStr != "" {
		t, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			dateFrom = &t
		}
	}
	if dateToStr != "" {
		t, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			// End of day
			end := t.Add(24*time.Hour - time.Second)
			dateTo = &end
		}
	}

	output, err := ctrl.listLessonsUseCase.Execute(c.Request.Context(), lesson.ListLessonsInput{
		ClassID:   classID,
		TeacherID: teacherID,
		Status:    status,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to list lessons", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lessons retrieved successfully", output)
}

func (ctrl *ControllerV1) GetLesson(c *gin.Context) {
	id := c.Param("id")

	output, err := ctrl.getLessonUseCase.Execute(c.Request.Context(), lesson.GetLessonInput{ID: id})
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Lesson not found", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lesson retrieved successfully", output.Lesson)
}

func (ctrl *ControllerV1) GetLessonAttendance(c *gin.Context) {
	output, err := ctrl.getLessonAttendanceUseCase.Execute(c.Request.Context(), lessonactivity.GetLessonAttendanceInput{
		LessonID: c.Param("id"),
		Actor:    buildLessonActor(c),
	})
	if err != nil {
		handleLessonActivityError(c, "Failed to retrieve lesson attendance", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lesson attendance retrieved successfully", output)
}

func (ctrl *ControllerV1) UpsertLessonAttendance(c *gin.Context) {
	var req UpsertLessonAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	records := make([]lessonactivity.UpsertLessonAttendanceRecord, 0, len(req.Records))
	for _, record := range req.Records {
		records = append(records, lessonactivity.UpsertLessonAttendanceRecord{
			StudentID: record.StudentID,
			Status:    record.Status,
			Note:      record.Note,
		})
	}

	output, err := ctrl.upsertLessonAttendanceUseCase.Execute(c.Request.Context(), lessonactivity.UpsertLessonAttendanceInput{
		LessonID: c.Param("id"),
		Actor:    buildLessonActor(c),
		Records:  records,
	})
	if err != nil {
		handleLessonActivityError(c, "Failed to save lesson attendance", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lesson attendance saved successfully", output)
}

func (ctrl *ControllerV1) GetLessonSummary(c *gin.Context) {
	output, err := ctrl.getLessonSummaryUseCase.Execute(c.Request.Context(), lessonactivity.GetLessonSummaryInput{
		LessonID: c.Param("id"),
		Actor:    buildLessonActor(c),
	})
	if err != nil {
		handleLessonActivityError(c, "Failed to retrieve lesson summary", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lesson summary retrieved successfully", output)
}

func (ctrl *ControllerV1) UpsertLessonSummary(c *gin.Context) {
	var req UpsertLessonSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.upsertLessonSummaryUseCase.Execute(c.Request.Context(), lessonactivity.UpsertLessonSummaryInput{
		LessonID:         c.Param("id"),
		Actor:            buildLessonActor(c),
		Topic:            req.Topic,
		LessonContent:    req.LessonContent,
		ClassFeedback:    req.ClassFeedback,
		Homework:         req.Homework,
		HomeworkDeadline: req.HomeworkDeadline,
		TeacherNotes:     req.TeacherNotes,
	})
	if err != nil {
		handleLessonActivityError(c, "Failed to save lesson summary", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lesson summary saved successfully", output)
}

func (ctrl *ControllerV1) GetLessonAcademicRecords(c *gin.Context) {
	output, err := ctrl.getLessonAcademicRecordsUseCase.Execute(c.Request.Context(), lessonrecord.GetLessonAcademicRecordsInput{
		LessonID: c.Param("id"),
		Actor:    buildLessonRecordActor(c),
	})
	if err != nil {
		handleLessonRecordError(c, "Failed to retrieve lesson academic records", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lesson academic records retrieved successfully", output)
}

func (ctrl *ControllerV1) UpsertLessonAcademicRecords(c *gin.Context) {
	var req UpsertLessonAcademicRecordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	records := make([]lessonrecord.UpsertLessonAcademicRecordRow, 0, len(req.Records))
	for _, row := range req.Records {
		records = append(records, lessonrecord.UpsertLessonAcademicRecordRow{
			StudentID:          row.StudentID,
			HomeworkCompleted:  row.HomeworkCompleted,
			HomeworkScore:      row.HomeworkScore,
			AttitudeRating:     row.AttitudeRating,
			ParticipationScore: row.ParticipationScore,
			PersonalComment:    row.PersonalComment,
		})
	}

	output, err := ctrl.upsertLessonAcademicRecordsUseCase.Execute(c.Request.Context(), lessonrecord.UpsertLessonAcademicRecordsInput{
		LessonID: c.Param("id"),
		Actor:    buildLessonRecordActor(c),
		Records:  records,
	})
	if err != nil {
		handleLessonRecordError(c, "Failed to save lesson academic records", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lesson academic records saved successfully", output)
}

func (ctrl *ControllerV1) FinalizeLessonAcademicRecords(c *gin.Context) {
	output, err := ctrl.finalizeLessonAcademicRecordsUseCase.Execute(c.Request.Context(), lessonrecord.FinalizeLessonAcademicRecordsInput{
		LessonID: c.Param("id"),
		Actor:    buildLessonRecordActor(c),
	})
	if err != nil {
		handleLessonRecordError(c, "Failed to finalize lesson academic records", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Lesson academic records finalized successfully", output)
}

func buildLessonActor(c *gin.Context) lessonactivity.LessonActor {
	userRole, _ := c.Get("user_role")
	userEmail, _ := c.Get("user_email")
	userID, _ := c.Get("user_id")

	actor := lessonactivity.LessonActor{}
	if userRole != nil {
		actor.Role, _ = userRole.(string)
	}
	if userEmail != nil {
		actor.Email, _ = userEmail.(string)
	}
	if userID != nil {
		actor.UserID, _ = userID.(string)
	}
	return actor
}

func buildLessonRecordActor(c *gin.Context) lessonrecord.LessonActor {
	userRole, _ := c.Get("user_role")
	userEmail, _ := c.Get("user_email")
	userID, _ := c.Get("user_id")

	actor := lessonrecord.LessonActor{}
	if userRole != nil {
		actor.Role, _ = userRole.(string)
	}
	if userEmail != nil {
		actor.Email, _ = userEmail.(string)
	}
	if userID != nil {
		actor.UserID, _ = userID.(string)
	}
	return actor
}

func handleLessonActivityError(c *gin.Context, message string, err error) {
	switch {
	case errors.Is(err, lessonactivity.ErrLessonNotFound):
		rest.ResponseError(c, http.StatusNotFound, message, err)
	case errors.Is(err, lessonactivity.ErrLessonAccessDenied):
		rest.ResponseError(c, http.StatusForbidden, message, err)
	case errors.Is(err, lessonactivity.ErrTeacherProfileMissing),
		errors.Is(err, lessonactivity.ErrInvalidAttendanceRow):
		rest.ResponseError(c, http.StatusBadRequest, message, err)
	default:
		rest.ResponseError(c, http.StatusInternalServerError, message, err)
	}
}

func handleLessonRecordError(c *gin.Context, message string, err error) {
	switch {
	case errors.Is(err, lessonrecord.ErrLessonNotFound):
		rest.ResponseError(c, http.StatusNotFound, message, err)
	case errors.Is(err, lessonrecord.ErrLessonAccessDenied):
		rest.ResponseError(c, http.StatusForbidden, message, err)
	case errors.Is(err, lessonrecord.ErrInvalidRecordRow),
		errors.Is(err, lessonrecord.ErrStudentNotFound):
		rest.ResponseError(c, http.StatusBadRequest, message, err)
	default:
		rest.ResponseError(c, http.StatusInternalServerError, message, err)
	}
}
