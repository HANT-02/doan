package studentportal

import (
	"errors"
	"net/http"
	"time"

	"doan/cmd/http/rest"
	studentportaluc "doan/internal/usecases/studentportal"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	getStudentTimetableUseCase  studentportaluc.GetStudentTimetableUseCase
	getMyAttendanceUseCase      studentportaluc.GetMyAttendanceUseCase
	getMyAcademicRecordsUseCase studentportaluc.GetMyAcademicRecordsUseCase
}

func NewStudentPortalControllerV1(
	getStudentTimetableUseCase studentportaluc.GetStudentTimetableUseCase,
	getMyAttendanceUseCase studentportaluc.GetMyAttendanceUseCase,
	getMyAcademicRecordsUseCase studentportaluc.GetMyAcademicRecordsUseCase,
) *ControllerV1 {
	return &ControllerV1{
		getStudentTimetableUseCase:  getStudentTimetableUseCase,
		getMyAttendanceUseCase:      getMyAttendanceUseCase,
		getMyAcademicRecordsUseCase: getMyAcademicRecordsUseCase,
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

func handleStudentPortalError(c *gin.Context, message string, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, studentportaluc.ErrStudentAccessDenied):
		status = http.StatusForbidden
	case errors.Is(err, studentportaluc.ErrStudentNotFound):
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
