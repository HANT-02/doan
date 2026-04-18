package teacherportal

import (
	"errors"
	"net/http"
	"time"

	"doan/cmd/http/rest"
	leaveflow "doan/internal/usecases/leaveflow"
	lessonactivity "doan/internal/usecases/lessonactivity"
	lessonrecord "doan/internal/usecases/lessonrecord"
	teacherportaluc "doan/internal/usecases/teacherportal"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	getTeacherLessonsUseCase             teacherportaluc.GetTeacherLessonsUseCase
	getAttendanceByLessonUseCase         teacherportaluc.GetAttendanceByLessonUseCase
	submitLessonAttendanceUseCase        teacherportaluc.SubmitLessonAttendanceUseCase
	markAttendanceUseCase                teacherportaluc.MarkAttendanceUseCase
	getAttendanceSummaryByStudentUseCase teacherportaluc.GetAttendanceSummaryByStudentUseCase
	getLessonSummaryUseCase              teacherportaluc.GetLessonSummaryUseCase
	upsertLessonSummaryUseCase           teacherportaluc.UpsertLessonSummaryUseCase
	getAcademicRecordsByLessonUseCase    teacherportaluc.GetAcademicRecordsByLessonUseCase
	upsertAcademicRecordUseCase          teacherportaluc.UpsertAcademicRecordUseCase
	finalizeAcademicRecordUseCase        teacherportaluc.FinalizeAcademicRecordUseCase
	getAcademicRecordsByStudentUseCase   teacherportaluc.GetAcademicRecordsByStudentUseCase
	listLeaveRequestsForTeacherUseCase   teacherportaluc.ListLeaveRequestsForTeacherUseCase
	approveLeaveRequestUseCase           teacherportaluc.ApproveLeaveRequestUseCase
	rejectLeaveRequestUseCase            teacherportaluc.RejectLeaveRequestUseCase
}

func NewTeacherPortalControllerV1(
	getTeacherLessonsUseCase teacherportaluc.GetTeacherLessonsUseCase,
	getAttendanceByLessonUseCase teacherportaluc.GetAttendanceByLessonUseCase,
	submitLessonAttendanceUseCase teacherportaluc.SubmitLessonAttendanceUseCase,
	markAttendanceUseCase teacherportaluc.MarkAttendanceUseCase,
	getAttendanceSummaryByStudentUseCase teacherportaluc.GetAttendanceSummaryByStudentUseCase,
	getLessonSummaryUseCase teacherportaluc.GetLessonSummaryUseCase,
	upsertLessonSummaryUseCase teacherportaluc.UpsertLessonSummaryUseCase,
	getAcademicRecordsByLessonUseCase teacherportaluc.GetAcademicRecordsByLessonUseCase,
	upsertAcademicRecordUseCase teacherportaluc.UpsertAcademicRecordUseCase,
	finalizeAcademicRecordUseCase teacherportaluc.FinalizeAcademicRecordUseCase,
	getAcademicRecordsByStudentUseCase teacherportaluc.GetAcademicRecordsByStudentUseCase,
	listLeaveRequestsForTeacherUseCase teacherportaluc.ListLeaveRequestsForTeacherUseCase,
	approveLeaveRequestUseCase teacherportaluc.ApproveLeaveRequestUseCase,
	rejectLeaveRequestUseCase teacherportaluc.RejectLeaveRequestUseCase,
) *ControllerV1 {
	return &ControllerV1{
		getTeacherLessonsUseCase:             getTeacherLessonsUseCase,
		getAttendanceByLessonUseCase:         getAttendanceByLessonUseCase,
		submitLessonAttendanceUseCase:        submitLessonAttendanceUseCase,
		markAttendanceUseCase:                markAttendanceUseCase,
		getAttendanceSummaryByStudentUseCase: getAttendanceSummaryByStudentUseCase,
		getLessonSummaryUseCase:              getLessonSummaryUseCase,
		upsertLessonSummaryUseCase:           upsertLessonSummaryUseCase,
		getAcademicRecordsByLessonUseCase:    getAcademicRecordsByLessonUseCase,
		upsertAcademicRecordUseCase:          upsertAcademicRecordUseCase,
		finalizeAcademicRecordUseCase:        finalizeAcademicRecordUseCase,
		getAcademicRecordsByStudentUseCase:   getAcademicRecordsByStudentUseCase,
		listLeaveRequestsForTeacherUseCase:   listLeaveRequestsForTeacherUseCase,
		approveLeaveRequestUseCase:           approveLeaveRequestUseCase,
		rejectLeaveRequestUseCase:            rejectLeaveRequestUseCase,
	}
}

func (ctrl *ControllerV1) GetTeacherLessons(c *gin.Context) {
	var dateFrom *time.Time
	if value := c.Query("from"); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			rest.ResponseError(c, http.StatusBadRequest, "Invalid from date format", err)
			return
		}
		dateFrom = &parsed
	}

	var dateTo *time.Time
	if value := c.Query("to"); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			rest.ResponseError(c, http.StatusBadRequest, "Invalid to date format", err)
			return
		}
		endOfDay := parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		dateTo = &endOfDay
	}

	output, err := ctrl.getTeacherLessonsUseCase.Execute(c.Request.Context(), teacherportaluc.GetTeacherLessonsInput{
		Actor:    buildActor(c),
		ClassID:  c.Query("class_id"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
	})
	if err != nil {
		status := http.StatusBadRequest
		switch err {
		case teacherportaluc.ErrTeacherAccessDenied:
			status = http.StatusForbidden
		case teacherportaluc.ErrTeacherProfileMissing:
			status = http.StatusNotFound
		}
		rest.ResponseError(c, status, "Failed to get teacher lessons", err)
		return
	}

	lessons := make([]TeacherLessonResponse, 0, len(output.Lessons))
	for _, lesson := range output.Lessons {
		item := TeacherLessonResponse{
			ID:        lesson.ID,
			ClassID:   lesson.ClassID,
			ClassName: lesson.ClassName,
			ClassCode: lesson.ClassCode,
			RoomID:    lesson.RoomID,
			RoomName:  lesson.RoomName,
			DateStart: lesson.DateStart,
			DateEnd:   lesson.DateEnd,
			Notes:     lesson.Notes,
		}
		if lesson.Shift != nil {
			item.Shift = &TeacherLessonShiftResponse{
				ID:              lesson.Shift.ID,
				Code:            lesson.Shift.Code,
				Name:            lesson.Shift.Name,
				StartTime:       lesson.Shift.StartTime,
				EndTime:         lesson.Shift.EndTime,
				DurationMinutes: lesson.Shift.DurationMinutes,
				SessionType:     lesson.Shift.SessionType,
			}
		}
		lessons = append(lessons, item)
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher lessons retrieved successfully", GetTeacherLessonsResponse{
		TeacherID: output.TeacherID,
		Lessons:   lessons,
	})
}

func (ctrl *ControllerV1) GetTeacherLessonAttendance(c *gin.Context) {
	output, err := ctrl.getAttendanceByLessonUseCase.Execute(c.Request.Context(), teacherportaluc.GetAttendanceByLessonInput{
		Actor:    buildActor(c),
		LessonID: c.Param("lesson_id"),
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to get teacher lesson attendance", err)
		return
	}

	records := make([]TeacherAttendanceRecordResponse, 0, len(output.Records))
	for _, record := range output.Records {
		item := TeacherAttendanceRecordResponse{
			AttendanceID: record.AttendanceID,
			Status:       record.Status,
			Note:         record.Note,
			MarkedAt:     record.MarkedAt,
			Student: TeacherAttendanceStudentResponse{
				ID:       record.Student.ID,
				Code:     record.Student.Code,
				FullName: record.Student.FullName,
			},
		}
		records = append(records, item)
	}

	lesson := TeacherLessonResponse{
		ID:        output.Lesson.ID,
		ClassID:   output.Lesson.ClassID,
		ClassName: output.Lesson.Class.Name,
		ClassCode: output.Lesson.Class.Code,
		DateStart: output.Lesson.DateStart,
		DateEnd:   output.Lesson.DateEnd,
		Notes:     output.Lesson.Notes,
	}
	if output.Lesson.RoomID != nil {
		lesson.RoomID = output.Lesson.RoomID
		roomName := output.Lesson.Room.Name
		lesson.RoomName = &roomName
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher lesson attendance retrieved successfully", GetTeacherLessonAttendanceResponse{
		Lesson:  lesson,
		Records: records,
	})
}

func (ctrl *ControllerV1) SubmitTeacherLessonAttendance(c *gin.Context) {
	var req SubmitTeacherLessonAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	records := make([]teacherportaluc.SubmitLessonAttendanceRecord, 0, len(req.Records))
	for _, record := range req.Records {
		records = append(records, teacherportaluc.SubmitLessonAttendanceRecord{
			StudentID: record.StudentID,
			Status:    record.Status,
			Note:      record.Note,
		})
	}

	output, err := ctrl.submitLessonAttendanceUseCase.Execute(c.Request.Context(), teacherportaluc.SubmitLessonAttendanceInput{
		Actor:    buildActor(c),
		LessonID: c.Param("lesson_id"),
		Records:  records,
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to save teacher lesson attendance", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher lesson attendance saved successfully", TeacherAttendanceSaveResponse{
		SavedCount: output.SavedCount,
	})
}

func (ctrl *ControllerV1) UpdateTeacherLessonAttendance(c *gin.Context) {
	var req UpdateTeacherLessonAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.markAttendanceUseCase.Execute(c.Request.Context(), teacherportaluc.MarkAttendanceInput{
		Actor:     buildActor(c),
		LessonID:  c.Param("lesson_id"),
		StudentID: c.Param("student_id"),
		Status:    req.Status,
		Note:      req.Note,
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to update teacher lesson attendance", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher lesson attendance updated successfully", TeacherAttendanceSaveResponse{
		SavedCount: output.SavedCount,
	})
}

func (ctrl *ControllerV1) GetTeacherAttendanceSummary(c *gin.Context) {
	var dateFrom *time.Time
	if value := c.Query("from"); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			rest.ResponseError(c, http.StatusBadRequest, "Invalid from date format", err)
			return
		}
		dateFrom = &parsed
	}

	var dateTo *time.Time
	if value := c.Query("to"); value != "" {
		parsed, err := time.Parse("2006-01-02", value)
		if err != nil {
			rest.ResponseError(c, http.StatusBadRequest, "Invalid to date format", err)
			return
		}
		endOfDay := parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		dateTo = &endOfDay
	}

	output, err := ctrl.getAttendanceSummaryByStudentUseCase.Execute(c.Request.Context(), teacherportaluc.GetAttendanceSummaryByStudentInput{
		Actor:    buildActor(c),
		ClassID:  c.Param("class_id"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to get teacher attendance summary", err)
		return
	}

	students := make([]TeacherAttendanceStudentSummaryResponse, 0, len(output.Students))
	for _, item := range output.Students {
		students = append(students, TeacherAttendanceStudentSummaryResponse{
			Student: TeacherAttendanceStudentResponse{
				ID:       item.Student.ID,
				Code:     item.Student.Code,
				FullName: item.Student.FullName,
			},
			TotalLessons:   item.TotalLessons,
			MarkedCount:    item.MarkedCount,
			PresentCount:   item.PresentCount,
			AbsentCount:    item.AbsentCount,
			LateCount:      item.LateCount,
			ExcusedCount:   item.ExcusedCount,
			UnmarkedCount:  item.UnmarkedCount,
			AttendanceRate: item.AttendanceRate,
		})
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher attendance summary retrieved successfully", GetTeacherAttendanceSummaryResponse{
		TeacherID:    output.TeacherID,
		ClassID:      output.ClassID,
		TotalLessons: output.TotalLessons,
		Students:     students,
	})
}

func (ctrl *ControllerV1) GetTeacherLessonSummary(c *gin.Context) {
	output, err := ctrl.getLessonSummaryUseCase.Execute(c.Request.Context(), teacherportaluc.GetLessonSummaryInput{
		Actor:    buildActor(c),
		LessonID: c.Param("lesson_id"),
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to get teacher lesson summary", err)
		return
	}

	response := GetTeacherLessonSummaryResponse{
		Lesson: TeacherLessonResponse{
			ID:        output.Lesson.ID,
			ClassID:   output.Lesson.ClassID,
			ClassName: output.Lesson.ClassName,
			ClassCode: output.Lesson.ClassCode,
			RoomID:    output.Lesson.RoomID,
			RoomName:  output.Lesson.RoomName,
			DateStart: output.Lesson.DateStart,
			DateEnd:   output.Lesson.DateEnd,
			Notes:     output.Lesson.Notes,
		},
	}
	if output.Summary != nil {
		response.Summary = &TeacherLessonSummaryResponse{
			ID:               output.Summary.ID,
			LessonID:         output.Summary.LessonID,
			Topic:            output.Summary.Topic,
			LessonContent:    output.Summary.LessonContent,
			ClassFeedback:    output.Summary.ClassFeedback,
			Homework:         output.Summary.Homework,
			HomeworkDeadline: output.Summary.HomeworkDeadline,
			TeacherNotes:     output.Summary.TeacherNotes,
			CreatedByID:      output.Summary.CreatedByID,
			CreatedAt:        output.Summary.CreatedAt,
			UpdatedAt:        output.Summary.UpdatedAt,
		}
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher lesson summary retrieved successfully", response)
}

func (ctrl *ControllerV1) UpsertTeacherLessonSummary(c *gin.Context) {
	var req UpsertTeacherLessonSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.upsertLessonSummaryUseCase.Execute(c.Request.Context(), teacherportaluc.UpsertLessonSummaryInput{
		Actor:            buildActor(c),
		LessonID:         c.Param("lesson_id"),
		Topic:            req.Topic,
		LessonContent:    req.LessonContent,
		ClassFeedback:    req.ClassFeedback,
		Homework:         req.Homework,
		HomeworkDeadline: req.HomeworkDeadline,
		TeacherNotes:     req.TeacherNotes,
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to save teacher lesson summary", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher lesson summary saved successfully", TeacherLessonSummaryResponse{
		ID:               output.Summary.ID,
		LessonID:         output.Summary.LessonID,
		Topic:            output.Summary.Topic,
		LessonContent:    output.Summary.LessonContent,
		ClassFeedback:    output.Summary.ClassFeedback,
		Homework:         output.Summary.Homework,
		HomeworkDeadline: output.Summary.HomeworkDeadline,
		TeacherNotes:     output.Summary.TeacherNotes,
		CreatedByID:      output.Summary.CreatedByID,
		CreatedAt:        output.Summary.CreatedAt,
		UpdatedAt:        output.Summary.UpdatedAt,
	})
}

func (ctrl *ControllerV1) GetTeacherLessonAcademicRecords(c *gin.Context) {
	output, err := ctrl.getAcademicRecordsByLessonUseCase.Execute(c.Request.Context(), teacherportaluc.GetAcademicRecordsByLessonInput{
		Actor:    buildActor(c),
		LessonID: c.Param("lesson_id"),
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to get teacher lesson academic records", err)
		return
	}

	records := make([]TeacherAcademicRecordResponse, 0, len(output.Records))
	for _, item := range output.Records {
		records = append(records, TeacherAcademicRecordResponse{
			RecordID:           item.RecordID,
			LessonSummaryID:    item.LessonSummaryID,
			Student:            TeacherAttendanceStudentResponse{ID: item.Student.ID, Code: item.Student.Code, FullName: item.Student.FullName},
			HomeworkCompleted:  item.HomeworkCompleted,
			HomeworkScore:      item.HomeworkScore,
			AttitudeRating:     item.AttitudeRating,
			ParticipationScore: item.ParticipationScore,
			PersonalComment:    item.PersonalComment,
			TotalScore:         item.TotalScore,
			IsCompleted:        item.IsCompleted,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
		})
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher lesson academic records retrieved successfully", GetTeacherLessonAcademicRecordsResponse{
		Lesson: TeacherLessonResponse{
			ID:        output.Lesson.ID,
			ClassID:   output.Lesson.ClassID,
			ClassName: output.Lesson.ClassName,
			ClassCode: output.Lesson.ClassCode,
			RoomID:    output.Lesson.RoomID,
			RoomName:  output.Lesson.RoomName,
			DateStart: output.Lesson.DateStart,
			DateEnd:   output.Lesson.DateEnd,
			Notes:     output.Lesson.Notes,
		},
		Records: records,
	})
}

func (ctrl *ControllerV1) UpsertTeacherLessonAcademicRecord(c *gin.Context) {
	var req UpsertTeacherLessonAcademicRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.upsertAcademicRecordUseCase.Execute(c.Request.Context(), teacherportaluc.UpsertAcademicRecordInput{
		Actor:              buildActor(c),
		LessonID:           c.Param("lesson_id"),
		StudentID:          c.Param("student_id"),
		HomeworkCompleted:  req.HomeworkCompleted,
		HomeworkScore:      req.HomeworkScore,
		AttitudeRating:     req.AttitudeRating,
		ParticipationScore: req.ParticipationScore,
		PersonalComment:    req.PersonalComment,
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to save teacher lesson academic record", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher lesson academic record saved successfully", TeacherAcademicRecordSaveResponse{
		SavedCount: output.SavedCount,
	})
}

func (ctrl *ControllerV1) FinalizeTeacherLessonAcademicRecords(c *gin.Context) {
	output, err := ctrl.finalizeAcademicRecordUseCase.Execute(c.Request.Context(), teacherportaluc.FinalizeAcademicRecordInput{
		Actor:    buildActor(c),
		LessonID: c.Param("lesson_id"),
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to finalize teacher lesson academic records", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher lesson academic records finalized successfully", TeacherAcademicRecordFinalizeResponse{
		FinalizedCount: output.FinalizedCount,
	})
}

func (ctrl *ControllerV1) GetTeacherStudentAcademicRecords(c *gin.Context) {
	output, err := ctrl.getAcademicRecordsByStudentUseCase.Execute(c.Request.Context(), teacherportaluc.GetAcademicRecordsByStudentInput{
		Actor:     buildActor(c),
		ClassID:   c.Param("class_id"),
		StudentID: c.Param("student_id"),
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to get teacher student academic records", err)
		return
	}

	records := make([]TeacherAcademicRecordResponse, 0, len(output.History.Records))
	for _, item := range output.History.Records {
		records = append(records, TeacherAcademicRecordResponse{
			RecordID:           item.RecordID,
			LessonSummaryID:    item.LessonSummaryID,
			Student:            TeacherAttendanceStudentResponse{ID: item.Student.ID, Code: item.Student.Code, FullName: item.Student.FullName},
			HomeworkCompleted:  item.HomeworkCompleted,
			HomeworkScore:      item.HomeworkScore,
			AttitudeRating:     item.AttitudeRating,
			ParticipationScore: item.ParticipationScore,
			PersonalComment:    item.PersonalComment,
			TotalScore:         item.TotalScore,
			IsCompleted:        item.IsCompleted,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
		})
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher student academic records retrieved successfully", GetTeacherStudentAcademicRecordsResponse{
		TeacherID: output.TeacherID,
		ClassID:   output.ClassID,
		StudentID: output.StudentID,
		Student: TeacherAttendanceStudentResponse{
			ID:       output.History.Student.ID,
			Code:     output.History.Student.Code,
			FullName: output.History.Student.FullName,
		},
		Records: records,
	})
}

func (ctrl *ControllerV1) ListTeacherLeaveRequests(c *gin.Context) {
	output, err := ctrl.listLeaveRequestsForTeacherUseCase.Execute(c.Request.Context(), teacherportaluc.ListLeaveRequestsForTeacherInput{
		Actor:     buildActor(c),
		Status:    c.Query("status"),
		ClassID:   c.Query("class_id"),
		StudentID: c.Query("student_id"),
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to list teacher leave requests", err)
		return
	}

	requests := make([]TeacherLeaveRequestResponse, 0, len(output.Requests))
	for _, item := range output.Requests {
		request := TeacherLeaveRequestResponse{
			ID: item.ID,
			Student: TeacherLeaveRequestStudentResponse{
				ID:       item.Student.ID,
				Code:     item.Student.Code,
				FullName: item.Student.FullName,
			},
			LeaveType:       item.LeaveType,
			ApplyDate:       item.ApplyDate,
			LateMinutes:     item.LateMinutes,
			EarlyMinutes:    item.EarlyMinutes,
			Reason:          item.Reason,
			Documents:       item.Documents,
			Subject:         item.Subject,
			Status:          item.Status,
			ApprovedByID:    item.ApprovedByID,
			ApprovedAt:      item.ApprovedAt,
			RejectionReason: item.RejectionReason,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
		}
		if item.Class != nil {
			request.Class = &TeacherLeaveRequestClassResponse{
				ID:   item.Class.ID,
				Code: item.Class.Code,
				Name: item.Class.Name,
			}
		}
		if item.Lesson != nil {
			request.Lesson = &TeacherLeaveRequestLessonResponse{
				ID:        item.Lesson.ID,
				DateStart: item.Lesson.DateStart,
				DateEnd:   item.Lesson.DateEnd,
			}
		}
		requests = append(requests, request)
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher leave requests retrieved successfully", ListTeacherLeaveRequestsResponse{
		Requests: requests,
	})
}

func (ctrl *ControllerV1) ApproveTeacherLeaveRequest(c *gin.Context) {
	output, err := ctrl.approveLeaveRequestUseCase.Execute(c.Request.Context(), teacherportaluc.ApproveLeaveRequestInput{
		Actor: buildActor(c),
		ID:    c.Param("id"),
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to approve teacher leave request", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher leave request approved successfully", TeacherLeaveRequestStatusResponse{
		RequestID: output.RequestID,
		Status:    output.Status,
	})
}

func (ctrl *ControllerV1) RejectTeacherLeaveRequest(c *gin.Context) {
	var req RejectTeacherLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.rejectLeaveRequestUseCase.Execute(c.Request.Context(), teacherportaluc.RejectLeaveRequestInput{
		Actor:           buildActor(c),
		ID:              c.Param("id"),
		RejectionReason: req.RejectionReason,
	})
	if err != nil {
		handleTeacherPortalError(c, "Failed to reject teacher leave request", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Teacher leave request rejected successfully", TeacherLeaveRequestStatusResponse{
		RequestID: output.RequestID,
		Status:    output.Status,
	})
}

func handleTeacherPortalError(c *gin.Context, message string, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, teacherportaluc.ErrTeacherAccessDenied), errors.Is(err, lessonactivity.ErrLessonAccessDenied):
		status = http.StatusForbidden
	case errors.Is(err, teacherportaluc.ErrTeacherProfileMissing), errors.Is(err, lessonactivity.ErrTeacherProfileMissing), errors.Is(err, lessonactivity.ErrLessonNotFound):
		status = http.StatusNotFound
	case errors.Is(err, teacherportaluc.ErrInvalidTeacherAttendanceStatus), errors.Is(err, lessonactivity.ErrInvalidAttendanceRow):
		status = http.StatusBadRequest
	case errors.Is(err, lessonrecord.ErrInvalidRecordRow), errors.Is(err, lessonrecord.ErrStudentNotFound):
		status = http.StatusBadRequest
	case errors.Is(err, lessonrecord.ErrLessonAccessDenied):
		status = http.StatusForbidden
	case errors.Is(err, lessonrecord.ErrLessonNotFound):
		status = http.StatusNotFound
	case errors.Is(err, leaveflow.ErrLeaveRequestForbidden):
		status = http.StatusForbidden
	case errors.Is(err, leaveflow.ErrLeaveRequestNotFound):
		status = http.StatusNotFound
	case errors.Is(err, leaveflow.ErrLeaveRequestNotPending), errors.Is(err, leaveflow.ErrTeacherNotFound):
		status = http.StatusBadRequest
	}
	rest.ResponseError(c, status, message, err)
}

func buildActor(c *gin.Context) teacherportaluc.Actor {
	userRole, _ := c.Get("user_role")
	userEmail, _ := c.Get("user_email")
	userID, _ := c.Get("user_id")

	actor := teacherportaluc.Actor{}
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
