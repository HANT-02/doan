package studentportal

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"doan/cmd/http/rest"
	leaveflow "doan/internal/usecases/leaveflow"
	studentportaluc "doan/internal/usecases/studentportal"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	getStudentTimetableUseCase   studentportaluc.GetStudentTimetableUseCase
	getMyAttendanceUseCase       studentportaluc.GetMyAttendanceUseCase
	getMyAcademicRecordsUseCase  studentportaluc.GetMyAcademicRecordsUseCase
	getMyAtRiskPredictionUseCase studentportaluc.GetMyAtRiskPredictionUseCase
	listMyLeaveRequestsUseCase   studentportaluc.ListMyLeaveRequestsUseCase
	createMyLeaveRequestUseCase  studentportaluc.CreateMyLeaveRequestUseCase
	cancelMyLeaveRequestUseCase  studentportaluc.CancelMyLeaveRequestUseCase
}

func NewStudentPortalControllerV1(
	getStudentTimetableUseCase studentportaluc.GetStudentTimetableUseCase,
	getMyAttendanceUseCase studentportaluc.GetMyAttendanceUseCase,
	getMyAcademicRecordsUseCase studentportaluc.GetMyAcademicRecordsUseCase,
	getMyAtRiskPredictionUseCase studentportaluc.GetMyAtRiskPredictionUseCase,
	listMyLeaveRequestsUseCase studentportaluc.ListMyLeaveRequestsUseCase,
	createMyLeaveRequestUseCase studentportaluc.CreateMyLeaveRequestUseCase,
	cancelMyLeaveRequestUseCase studentportaluc.CancelMyLeaveRequestUseCase,
) *ControllerV1 {
	return &ControllerV1{
		getStudentTimetableUseCase:   getStudentTimetableUseCase,
		getMyAttendanceUseCase:       getMyAttendanceUseCase,
		getMyAcademicRecordsUseCase:  getMyAcademicRecordsUseCase,
		getMyAtRiskPredictionUseCase: getMyAtRiskPredictionUseCase,
		listMyLeaveRequestsUseCase:   listMyLeaveRequestsUseCase,
		createMyLeaveRequestUseCase:  createMyLeaveRequestUseCase,
		cancelMyLeaveRequestUseCase:  cancelMyLeaveRequestUseCase,
	}
}

func (ctrl *ControllerV1) GetStudentTimetable(c *gin.Context) {
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

	output, err := ctrl.getStudentTimetableUseCase.Execute(c.Request.Context(), studentportaluc.GetStudentTimetableInput{
		Actor:    buildActor(c),
		ClassID:  c.Query("class_id"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
	})
	if err != nil {
		handleStudentPortalError(c, "Failed to get student timetable", err)
		return
	}

	lessons := make([]StudentTimetableLessonResponse, 0, len(output.Lessons))
	for _, lesson := range output.Lessons {
		item := StudentTimetableLessonResponse{
			ID:        lesson.ID,
			ClassID:   lesson.ClassID,
			ClassName: lesson.ClassName,
			ClassCode: lesson.ClassCode,
			Teacher: StudentTimetableTeacherResponse{
				ID:       lesson.Teacher.ID,
				Code:     lesson.Teacher.Code,
				FullName: lesson.Teacher.FullName,
			},
			RoomID:    lesson.RoomID,
			RoomName:  lesson.RoomName,
			DateStart: lesson.DateStart,
			DateEnd:   lesson.DateEnd,
			Notes:     lesson.Notes,
		}
		if lesson.Shift != nil {
			item.Shift = &StudentTimetableShiftResponse{
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

	rest.ResponseSuccess(c, http.StatusOK, "Student timetable retrieved successfully", GetStudentTimetableResponse{
		StudentID: output.StudentID,
		Lessons:   lessons,
	})
}

func (ctrl *ControllerV1) GetMyAttendance(c *gin.Context) {
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

	output, err := ctrl.getMyAttendanceUseCase.Execute(c.Request.Context(), studentportaluc.GetMyAttendanceInput{
		Actor:    buildActor(c),
		ClassID:  c.Query("class_id"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
	})
	if err != nil {
		handleStudentPortalError(c, "Failed to get student attendance", err)
		return
	}

	records := make([]StudentAttendanceRecordResponse, 0, len(output.Records))
	for _, record := range output.Records {
		item := StudentAttendanceRecordResponse{
			Lesson: StudentTimetableLessonResponse{
				ID:        record.Lesson.ID,
				ClassID:   record.Lesson.ClassID,
				ClassName: record.Lesson.ClassName,
				ClassCode: record.Lesson.ClassCode,
				Teacher: StudentTimetableTeacherResponse{
					ID:       record.Lesson.Teacher.ID,
					Code:     record.Lesson.Teacher.Code,
					FullName: record.Lesson.Teacher.FullName,
				},
				RoomID:    record.Lesson.RoomID,
				RoomName:  record.Lesson.RoomName,
				DateStart: record.Lesson.DateStart,
				DateEnd:   record.Lesson.DateEnd,
				Notes:     record.Lesson.Notes,
			},
			Status:   record.Status,
			Note:     record.Note,
			MarkedAt: record.MarkedAt,
		}
		if record.Lesson.Shift != nil {
			item.Lesson.Shift = &StudentTimetableShiftResponse{
				ID:              record.Lesson.Shift.ID,
				Code:            record.Lesson.Shift.Code,
				Name:            record.Lesson.Shift.Name,
				StartTime:       record.Lesson.Shift.StartTime,
				EndTime:         record.Lesson.Shift.EndTime,
				DurationMinutes: record.Lesson.Shift.DurationMinutes,
				SessionType:     record.Lesson.Shift.SessionType,
			}
		}
		records = append(records, item)
	}

	rest.ResponseSuccess(c, http.StatusOK, "Student attendance retrieved successfully", GetMyAttendanceResponse{
		StudentID: output.StudentID,
		ClassID:   output.ClassID,
		Summary: StudentAttendanceSummaryResponse{
			TotalLessons:   output.Summary.TotalLessons,
			MarkedCount:    output.Summary.MarkedCount,
			PresentCount:   output.Summary.PresentCount,
			AbsentCount:    output.Summary.AbsentCount,
			LateCount:      output.Summary.LateCount,
			ExcusedCount:   output.Summary.ExcusedCount,
			UnmarkedCount:  output.Summary.UnmarkedCount,
			AttendanceRate: output.Summary.AttendanceRate,
			AbsentRate:     output.Summary.AbsentRate,
			Warning:        output.Summary.Warning,
			WarningMessage: output.Summary.WarningMessage,
		},
		Records: records,
	})
}

func (ctrl *ControllerV1) GetMyAcademicRecords(c *gin.Context) {
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

	output, err := ctrl.getMyAcademicRecordsUseCase.Execute(c.Request.Context(), studentportaluc.GetMyAcademicRecordsInput{
		Actor:    buildActor(c),
		ClassID:  c.Query("class_id"),
		DateFrom: dateFrom,
		DateTo:   dateTo,
	})
	if err != nil {
		handleStudentPortalError(c, "Failed to get student academic records", err)
		return
	}

	classSummaries := make([]StudentAcademicClassSummaryResponse, 0, len(output.ClassSummaries))
	for _, summary := range output.ClassSummaries {
		classSummaries = append(classSummaries, StudentAcademicClassSummaryResponse{
			ClassID:           summary.ClassID,
			ClassName:         summary.ClassName,
			ClassCode:         summary.ClassCode,
			RecordsCount:      summary.RecordsCount,
			CompletedCount:    summary.CompletedCount,
			AverageTotalScore: summary.AverageTotalScore,
		})
	}

	records := make([]StudentAcademicRecordItemResponse, 0, len(output.Records))
	for _, record := range output.Records {
		item := StudentAcademicRecordItemResponse{
			RecordID:        record.RecordID,
			LessonSummaryID: record.LessonSummaryID,
			Lesson: StudentAcademicRecordLessonResponse{
				ID:        record.Lesson.ID,
				ClassID:   record.Lesson.ClassID,
				ClassName: record.Lesson.ClassName,
				ClassCode: record.Lesson.ClassCode,
				Teacher: StudentTimetableTeacherResponse{
					ID:       record.Lesson.Teacher.ID,
					Code:     record.Lesson.Teacher.Code,
					FullName: record.Lesson.Teacher.FullName,
				},
				RoomID:    record.Lesson.RoomID,
				RoomName:  record.Lesson.RoomName,
				DateStart: record.Lesson.DateStart,
				DateEnd:   record.Lesson.DateEnd,
				Notes:     record.Lesson.Notes,
				Summary: StudentAcademicRecordLessonSummaryResponse{
					ID:       record.Lesson.Summary.ID,
					Topic:    record.Lesson.Summary.Topic,
					Homework: record.Lesson.Summary.Homework,
				},
			},
			HomeworkCompleted:  record.HomeworkCompleted,
			HomeworkScore:      record.HomeworkScore,
			AttitudeRating:     record.AttitudeRating,
			ParticipationScore: record.ParticipationScore,
			PersonalComment:    record.PersonalComment,
			TotalScore:         record.TotalScore,
			IsCompleted:        record.IsCompleted,
			CreatedAt:          record.CreatedAt,
			UpdatedAt:          record.UpdatedAt,
		}
		if record.Lesson.Shift != nil {
			item.Lesson.Shift = &StudentTimetableShiftResponse{
				ID:              record.Lesson.Shift.ID,
				Code:            record.Lesson.Shift.Code,
				Name:            record.Lesson.Shift.Name,
				StartTime:       record.Lesson.Shift.StartTime,
				EndTime:         record.Lesson.Shift.EndTime,
				DurationMinutes: record.Lesson.Shift.DurationMinutes,
				SessionType:     record.Lesson.Shift.SessionType,
			}
		}
		records = append(records, item)
	}

	rest.ResponseSuccess(c, http.StatusOK, "Student academic records retrieved successfully", GetMyAcademicRecordsResponse{
		StudentID:      output.StudentID,
		ClassID:        output.ClassID,
		ClassSummaries: classSummaries,
		Records:        records,
	})
}

func (ctrl *ControllerV1) GetMyAtRiskPrediction(c *gin.Context) {
	refresh, _ := strconv.ParseBool(c.DefaultQuery("refresh", "false"))

	output, err := ctrl.getMyAtRiskPredictionUseCase.Execute(c.Request.Context(), studentportaluc.GetMyAtRiskPredictionInput{
		Actor:   buildActor(c),
		Refresh: refresh,
	})
	if err != nil {
		handleStudentPortalError(c, "Failed to get student at-risk prediction", err)
		return
	}

	response := GetMyAtRiskPredictionResponse{
		StudentID: output.StudentID,
	}
	if output.Prediction != nil {
		topFeatures := make([]StudentAtRiskTopFeatureResponse, 0, len(output.Prediction.TopFeatures))
		for _, item := range output.Prediction.TopFeatures {
			topFeatures = append(topFeatures, StudentAtRiskTopFeatureResponse{
				Key:          item.Key,
				Label:        item.Label,
				Value:        item.Value,
				DisplayValue: item.DisplayValue,
			})
		}

		response.Prediction = &StudentAtRiskPredictionResponse{
			StudentID:     output.Prediction.StudentID,
			StudentCode:   output.Prediction.StudentCode,
			StudentName:   output.Prediction.StudentName,
			GradeLevel:    output.Prediction.GradeLevel,
			ClassID:       output.Prediction.ClassID,
			ClassCode:     output.Prediction.ClassCode,
			ClassName:     output.Prediction.ClassName,
			SnapshotAt:    output.Prediction.SnapshotAt,
			RiskLabel:     output.Prediction.RiskLabel,
			RiskScore:     output.Prediction.RiskScore,
			RiskBand:      output.Prediction.RiskBand,
			ModelName:     output.Prediction.ModelName,
			ModelVersion:  output.Prediction.ModelVersion,
			PrimaryReason: output.Prediction.PrimaryReason,
			Reasons:       append([]string(nil), output.Prediction.Reasons...),
			TopFeatures:   topFeatures,
			FeatureSummary: StudentAtRiskFeatureSummaryResponse{
				AttendanceRate28d:         output.Prediction.FeatureSummary.AttendanceRate28d,
				AbsenceCount28d:           output.Prediction.FeatureSummary.AbsenceCount28d,
				AverageTotalScore28d:      output.Prediction.FeatureSummary.AverageTotalScore28d,
				HomeworkCompletionRate28d: output.Prediction.FeatureSummary.HomeworkCompletionRate28d,
				ActiveEnrollmentCount28d:  output.Prediction.FeatureSummary.ActiveEnrollmentCount28d,
				WeeklyLessonLoad28d:       output.Prediction.FeatureSummary.WeeklyLessonLoad28d,
				ApprovedLeaveCount28d:     output.Prediction.FeatureSummary.ApprovedLeaveCount28d,
				DaysSinceLastLesson:       output.Prediction.FeatureSummary.DaysSinceLastLesson,
			},
		}
	}

	rest.ResponseSuccess(c, http.StatusOK, "Student at-risk prediction retrieved successfully", response)
}

func (ctrl *ControllerV1) ListMyLeaveRequests(c *gin.Context) {
	output, err := ctrl.listMyLeaveRequestsUseCase.Execute(c.Request.Context(), studentportaluc.ListMyLeaveRequestsInput{
		Actor:   buildActor(c),
		Status:  c.Query("status"),
		ClassID: c.Query("class_id"),
	})
	if err != nil {
		handleStudentPortalError(c, "Failed to list student leave requests", err)
		return
	}

	requests := make([]StudentLeaveRequestResponse, 0, len(output.Requests))
	for _, item := range output.Requests {
		request := StudentLeaveRequestResponse{
			ID: item.ID,
			Student: StudentLeaveRequestStudentResponse{
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
			request.Class = &StudentLeaveRequestClassResponse{
				ID:   item.Class.ID,
				Code: item.Class.Code,
				Name: item.Class.Name,
			}
		}
		if item.Lesson != nil {
			request.Lesson = &StudentLeaveRequestLessonResponse{
				ID:        item.Lesson.ID,
				DateStart: item.Lesson.DateStart,
				DateEnd:   item.Lesson.DateEnd,
			}
		}
		requests = append(requests, request)
	}

	rest.ResponseSuccess(c, http.StatusOK, "Student leave requests retrieved successfully", ListMyLeaveRequestsResponse{
		Requests: requests,
	})
}

func (ctrl *ControllerV1) CreateMyLeaveRequest(c *gin.Context) {
	var req CreateMyLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	applyDate, err := time.Parse(time.RFC3339, req.ApplyDate)
	if err != nil {
		applyDate, err = time.Parse("2006-01-02T15:04", req.ApplyDate)
		if err != nil {
			rest.ResponseError(c, http.StatusBadRequest, "Invalid apply_date format", err)
			return
		}
	}

	output, err := ctrl.createMyLeaveRequestUseCase.Execute(c.Request.Context(), studentportaluc.CreateMyLeaveRequestInput{
		Actor:        buildActor(c),
		LeaveType:    req.LeaveType,
		ApplyDate:    applyDate,
		LateMinutes:  req.LateMinutes,
		EarlyMinutes: req.EarlyMinutes,
		Reason:       req.Reason,
		Documents:    req.Documents,
		ClassID:      req.ClassID,
		LessonID:     req.LessonID,
		Subject:      req.Subject,
	})
	if err != nil {
		handleStudentPortalError(c, "Failed to create student leave request", err)
		return
	}

	request := StudentLeaveRequestResponse{
		ID: output.Request.ID,
		Student: StudentLeaveRequestStudentResponse{
			ID:       output.Request.Student.ID,
			Code:     output.Request.Student.Code,
			FullName: output.Request.Student.FullName,
		},
		LeaveType:       output.Request.LeaveType,
		ApplyDate:       output.Request.ApplyDate,
		LateMinutes:     output.Request.LateMinutes,
		EarlyMinutes:    output.Request.EarlyMinutes,
		Reason:          output.Request.Reason,
		Documents:       output.Request.Documents,
		Subject:         output.Request.Subject,
		Status:          output.Request.Status,
		ApprovedByID:    output.Request.ApprovedByID,
		ApprovedAt:      output.Request.ApprovedAt,
		RejectionReason: output.Request.RejectionReason,
		CreatedAt:       output.Request.CreatedAt,
		UpdatedAt:       output.Request.UpdatedAt,
	}
	if output.Request.Class != nil {
		request.Class = &StudentLeaveRequestClassResponse{
			ID:   output.Request.Class.ID,
			Code: output.Request.Class.Code,
			Name: output.Request.Class.Name,
		}
	}
	if output.Request.Lesson != nil {
		request.Lesson = &StudentLeaveRequestLessonResponse{
			ID:        output.Request.Lesson.ID,
			DateStart: output.Request.Lesson.DateStart,
			DateEnd:   output.Request.Lesson.DateEnd,
		}
	}

	rest.ResponseSuccess(c, http.StatusCreated, "Student leave request created successfully", request)
}

func (ctrl *ControllerV1) CancelMyLeaveRequest(c *gin.Context) {
	output, err := ctrl.cancelMyLeaveRequestUseCase.Execute(c.Request.Context(), studentportaluc.CancelMyLeaveRequestInput{
		Actor: buildActor(c),
		ID:    c.Param("id"),
	})
	if err != nil {
		handleStudentPortalError(c, "Failed to cancel student leave request", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Student leave request cancelled successfully", CancelMyLeaveRequestResponse{
		RequestID: output.RequestID,
		Status:    output.Status,
	})
}

func handleStudentPortalError(c *gin.Context, message string, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, studentportaluc.ErrStudentAccessDenied):
		status = http.StatusForbidden
	case errors.Is(err, studentportaluc.ErrStudentNotFound):
		status = http.StatusNotFound
	case errors.Is(err, leaveflow.ErrInvalidLeaveType), errors.Is(err, leaveflow.ErrLeaveRequestNotPending), errors.Is(err, leaveflow.ErrStudentNotInClass):
		status = http.StatusBadRequest
	case errors.Is(err, leaveflow.ErrLeaveRequestForbidden):
		status = http.StatusForbidden
	case errors.Is(err, leaveflow.ErrLeaveRequestNotFound):
		status = http.StatusNotFound
	}
	rest.ResponseError(c, status, message, err)
}

func buildActor(c *gin.Context) studentportaluc.Actor {
	userRole, _ := c.Get("user_role")
	userEmail, _ := c.Get("user_email")
	userID, _ := c.Get("user_id")

	actor := studentportaluc.Actor{}
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
