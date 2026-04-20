package studentportal

import (
	"context"
	"sort"
	"strings"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

const attendanceWarningThreshold = 0.20

type GetMyAttendanceInput struct {
	Actor    Actor
	ClassID  string
	DateFrom *time.Time
	DateTo   *time.Time
}

type StudentAttendanceShift struct {
	ID              string `json:"id"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	DurationMinutes int    `json:"duration_minutes"`
	SessionType     string `json:"session_type"`
}

type StudentAttendanceTeacher struct {
	ID       *string `json:"id,omitempty"`
	Code     *string `json:"code,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

type StudentAttendanceLesson struct {
	ID        string                   `json:"id"`
	ClassID   string                   `json:"class_id"`
	ClassName string                   `json:"class_name"`
	ClassCode string                   `json:"class_code"`
	Teacher   StudentAttendanceTeacher `json:"teacher"`
	RoomID    *string                  `json:"room_id,omitempty"`
	RoomName  *string                  `json:"room_name,omitempty"`
	DateStart time.Time                `json:"date_start"`
	DateEnd   time.Time                `json:"date_end"`
	Notes     string                   `json:"notes"`
	Shift     *StudentAttendanceShift  `json:"shift,omitempty"`
}

type StudentAttendanceRecord struct {
	Lesson   StudentAttendanceLesson `json:"lesson"`
	Status   *int                    `json:"status,omitempty"`
	Note     string                  `json:"note"`
	MarkedAt *time.Time              `json:"marked_at,omitempty"`
}

type StudentAttendanceSummary struct {
	TotalLessons   int     `json:"total_lessons"`
	MarkedCount    int     `json:"marked_count"`
	PresentCount   int     `json:"present_count"`
	AbsentCount    int     `json:"absent_count"`
	LateCount      int     `json:"late_count"`
	ExcusedCount   int     `json:"excused_count"`
	UnmarkedCount  int     `json:"unmarked_count"`
	AttendanceRate float64 `json:"attendance_rate"`
	AbsentRate     float64 `json:"absent_rate"`
	Warning        bool    `json:"warning"`
	WarningMessage string  `json:"warning_message,omitempty"`
}

type GetMyAttendanceOutput struct {
	StudentID string                    `json:"student_id"`
	ClassID   string                    `json:"class_id,omitempty"`
	Summary   StudentAttendanceSummary  `json:"summary"`
	Records   []StudentAttendanceRecord `json:"records"`
}

type GetMyAttendanceUseCase interface {
	Execute(ctx context.Context, input GetMyAttendanceInput) (*GetMyAttendanceOutput, error)
}

type getMyAttendanceUseCase struct {
	studentRepo       repointerface.StudentRepository
	enrollmentRepo    repointerface.EnrollmentRepository
	lessonRepo        repointerface.LessonRepository
	classScheduleRepo repointerface.ClassScheduleRepository
	attendanceRepo    repointerface.AttendanceRepository
}

func NewGetMyAttendanceUseCase(
	studentRepo repointerface.StudentRepository,
	enrollmentRepo repointerface.EnrollmentRepository,
	lessonRepo repointerface.LessonRepository,
	classScheduleRepo repointerface.ClassScheduleRepository,
	attendanceRepo repointerface.AttendanceRepository,
) GetMyAttendanceUseCase {
	return &getMyAttendanceUseCase{
		studentRepo:       studentRepo,
		enrollmentRepo:    enrollmentRepo,
		lessonRepo:        lessonRepo,
		classScheduleRepo: classScheduleRepo,
		attendanceRepo:    attendanceRepo,
	}
}

func (uc *getMyAttendanceUseCase) Execute(ctx context.Context, input GetMyAttendanceInput) (*GetMyAttendanceOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if strings.TrimSpace(input.Actor.Role) != "STUDENT" {
		return nil, ErrStudentAccessDenied
	}

	student, err := resolveStudentByEmail(ctx, uc.studentRepo, input.Actor.Email)
	if err != nil {
		ctxLogger.Errorf("Failed to resolve student from actor email %s: %v", input.Actor.Email, err)
		return nil, err
	}

	enrollmentCondition := repositories.NewCommonCondition()
	enrollmentCondition.AddCondition("student_id", student.ID, repositories.Equal)
	enrollmentCondition.AddCondition("status", "ENROLLED", repositories.Equal)
	enrollmentCondition.SetPaging(500, 1)

	enrollments, err := uc.enrollmentRepo.GetByCondition(ctx, enrollmentCondition)
	if err != nil {
		ctxLogger.Errorf("Failed to get enrollments for student %s: %v", student.ID, err)
		return nil, err
	}

	classIDs := make([]string, 0)
	classIDMap := make(map[string]struct{})
	if enrollments != nil {
		for _, enrollment := range enrollments.Data {
			if enrollment == nil {
				continue
			}
			if input.ClassID != "" && enrollment.ClassID != input.ClassID {
				continue
			}
			if _, exists := classIDMap[enrollment.ClassID]; exists {
				continue
			}
			classIDMap[enrollment.ClassID] = struct{}{}
			classIDs = append(classIDs, enrollment.ClassID)
		}
	}

	if len(classIDs) == 0 {
		return &GetMyAttendanceOutput{
			StudentID: student.ID,
			ClassID:   input.ClassID,
			Summary:   StudentAttendanceSummary{},
			Records:   []StudentAttendanceRecord{},
		}, nil
	}

	lessonCondition := repositories.NewCommonCondition()
	lessonCondition.AddCondition("class_id", classIDs, repositories.In)
	if input.DateFrom != nil {
		lessonCondition.AddCondition("date_end", input.DateFrom.UTC(), repositories.GreaterThanOrEqual)
	}
	if input.DateTo != nil {
		lessonCondition.AddCondition("date_start", input.DateTo.UTC(), repositories.LessThanOrEqual)
	}
	lessonCondition.SetPreload([]string{"Class", "Teacher", "Room"})
	lessonCondition.SetPaging(1000, 1)
	lessonCondition.AddSorting("date_start", repositories.Asc)

	lessonsResult, err := uc.lessonRepo.GetByCondition(ctx, lessonCondition)
	if err != nil {
		ctxLogger.Errorf("Failed to get lessons for student attendance %s: %v", student.ID, err)
		return nil, err
	}

	lessonIDs := make([]string, 0)
	classScheduleMap := make(map[string][]entities.ClassSchedule)
	lessonMap := make(map[string]*entities.Lesson)

	if lessonsResult != nil {
		for _, lesson := range lessonsResult.Data {
			if lesson == nil {
				continue
			}
			lessonIDs = append(lessonIDs, lesson.ID)
			lessonMap[lesson.ID] = lesson
		}
	}

	attendanceRecords, err := uc.attendanceRepo.ListByLessonIDs(ctx, lessonIDs)
	if err != nil {
		ctxLogger.Errorf("Failed to get attendance records for student %s: %v", student.ID, err)
		return nil, err
	}

	myAttendanceMap := make(map[string]entities.Attendance)
	for _, record := range attendanceRecords {
		if record.StudentID == student.ID {
			myAttendanceMap[record.LessonID] = record
		}
	}

	items := make([]StudentAttendanceRecord, 0, len(lessonMap))
	summary := StudentAttendanceSummary{
		TotalLessons: len(lessonMap),
	}

	for _, lesson := range lessonMap {
		schedules, exists := classScheduleMap[lesson.ClassID]
		if !exists {
			schedules, err = uc.classScheduleRepo.GetSchedulesByClassID(ctx, lesson.ClassID)
			if err != nil {
				ctxLogger.Errorf("Failed to load class schedules for class %s: %v", lesson.ClassID, err)
				return nil, err
			}
			classScheduleMap[lesson.ClassID] = schedules
		}

		item := StudentAttendanceRecord{
			Lesson: StudentAttendanceLesson{
				ID:        lesson.ID,
				ClassID:   lesson.ClassID,
				ClassName: lesson.Class.Name,
				ClassCode: lesson.Class.Code,
				DateStart: lesson.DateStart,
				DateEnd:   lesson.DateEnd,
				Notes:     lesson.Notes,
				RoomID:    lesson.RoomID,
				Teacher:   StudentAttendanceTeacher{},
			},
		}

		if lesson.TeacherID != nil {
			item.Lesson.Teacher.ID = lesson.TeacherID
		}
		if lesson.Teacher.Code != "" {
			code := lesson.Teacher.Code
			item.Lesson.Teacher.Code = &code
		}
		if lesson.Teacher.FullName != "" {
			name := lesson.Teacher.FullName
			item.Lesson.Teacher.FullName = &name
		}
		if lesson.Room.Name != "" {
			roomName := lesson.Room.Name
			item.Lesson.RoomName = &roomName
		}

		if matched := matchShiftForLesson(*lesson, schedules); matched != nil {
			item.Lesson.Shift = &StudentAttendanceShift{
				ID:              matched.Shift.ID,
				Code:            matched.Shift.Code,
				Name:            matched.Shift.Name,
				StartTime:       matched.Shift.StartTime,
				EndTime:         matched.Shift.EndTime,
				DurationMinutes: matched.Shift.DurationMinutes,
				SessionType:     matched.Shift.SessionType,
			}
			if item.Lesson.RoomID == nil && matched.RoomID != nil {
				item.Lesson.RoomID = matched.RoomID
			}
			if item.Lesson.RoomName == nil && matched.Room.Name != "" {
				roomName := matched.Room.Name
				item.Lesson.RoomName = &roomName
			}
		}

		if attendance, ok := myAttendanceMap[lesson.ID]; ok {
			status := attendance.Status
			item.Status = &status
			item.Note = attendance.Note
			item.MarkedAt = &attendance.MarkedAt
			summary.MarkedCount++
			switch attendance.Status {
			case 1:
				summary.PresentCount++
			case 0:
				summary.AbsentCount++
			case 2:
				summary.LateCount++
			case 3:
				summary.ExcusedCount++
			}
		}

		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Lesson.DateStart.Before(items[j].Lesson.DateStart)
	})

	if summary.TotalLessons > summary.MarkedCount {
		summary.UnmarkedCount = summary.TotalLessons - summary.MarkedCount
	}
	if summary.TotalLessons > 0 {
		summary.AttendanceRate = float64(summary.PresentCount+summary.LateCount) / float64(summary.TotalLessons)
		summary.AbsentRate = float64(summary.AbsentCount) / float64(summary.TotalLessons)
		summary.Warning = summary.AbsentRate > attendanceWarningThreshold
		if summary.Warning {
			summary.WarningMessage = "Tỷ lệ vắng của bạn đang vượt ngưỡng 20%. Hãy liên hệ giáo viên nếu cần hỗ trợ."
		}
	}

	return &GetMyAttendanceOutput{
		StudentID: student.ID,
		ClassID:   input.ClassID,
		Summary:   summary,
		Records:   items,
	}, nil
}
