package teacherportal

import (
	"context"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
)

type GetAttendanceSummaryByStudentInput struct {
	Actor    Actor
	ClassID  string
	DateFrom *time.Time
	DateTo   *time.Time
}

type AttendanceStudentSummary struct {
	Student        entities.Student `json:"student"`
	TotalLessons   int              `json:"total_lessons"`
	MarkedCount    int              `json:"marked_count"`
	PresentCount   int              `json:"present_count"`
	AbsentCount    int              `json:"absent_count"`
	LateCount      int              `json:"late_count"`
	ExcusedCount   int              `json:"excused_count"`
	UnmarkedCount  int              `json:"unmarked_count"`
	AttendanceRate float64          `json:"attendance_rate"`
}

type GetAttendanceSummaryByStudentOutput struct {
	TeacherID    string                     `json:"teacher_id"`
	ClassID      string                     `json:"class_id"`
	TotalLessons int                        `json:"total_lessons"`
	Students     []AttendanceStudentSummary `json:"students"`
}

type GetAttendanceSummaryByStudentUseCase interface {
	Execute(ctx context.Context, input GetAttendanceSummaryByStudentInput) (*GetAttendanceSummaryByStudentOutput, error)
}

type getAttendanceSummaryByStudentUseCase struct {
	teacherRepo    repointerface.TeacherRepository
	enrollmentRepo repointerface.EnrollmentRepository
	attendanceRepo repointerface.AttendanceRepository
}

func NewGetAttendanceSummaryByStudentUseCase(
	teacherRepo repointerface.TeacherRepository,
	enrollmentRepo repointerface.EnrollmentRepository,
	attendanceRepo repointerface.AttendanceRepository,
) GetAttendanceSummaryByStudentUseCase {
	return &getAttendanceSummaryByStudentUseCase{
		teacherRepo:    teacherRepo,
		enrollmentRepo: enrollmentRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (uc *getAttendanceSummaryByStudentUseCase) Execute(ctx context.Context, input GetAttendanceSummaryByStudentInput) (*GetAttendanceSummaryByStudentOutput, error) {
	if input.Actor.Role != "TEACHER" {
		return nil, ErrTeacherAccessDenied
	}

	teacher, err := resolveTeacherByEmail(ctx, uc.teacherRepo, input.Actor.Email)
	if err != nil {
		return nil, err
	}

	allLessons, err := uc.teacherRepo.GetTeacherLessons(ctx, teacher.ID, time.Time{}, time.Time{})
	if err != nil {
		return nil, err
	}

	hasClassAccess := false
	filteredLessons := make([]entities.Lesson, 0)
	for _, lesson := range allLessons {
		if lesson.ClassID != input.ClassID {
			continue
		}
		hasClassAccess = true
		if input.DateFrom != nil && lesson.DateStart.Before(*input.DateFrom) {
			continue
		}
		if input.DateTo != nil && lesson.DateEnd.After(*input.DateTo) {
			continue
		}
		filteredLessons = append(filteredLessons, lesson)
	}
	if !hasClassAccess {
		return nil, ErrTeacherAccessDenied
	}

	enrollments, err := uc.enrollmentRepo.ListByClassID(ctx, input.ClassID)
	if err != nil {
		return nil, err
	}

	lessonIDs := make([]string, 0, len(filteredLessons))
	for _, lesson := range filteredLessons {
		lessonIDs = append(lessonIDs, lesson.ID)
	}

	attendanceRecords, err := uc.attendanceRepo.ListByLessonIDs(ctx, lessonIDs)
	if err != nil {
		return nil, err
	}

	studentLessonStatus := make(map[string]map[string]int)
	for _, record := range attendanceRecords {
		if _, ok := studentLessonStatus[record.StudentID]; !ok {
			studentLessonStatus[record.StudentID] = make(map[string]int)
		}
		studentLessonStatus[record.StudentID][record.LessonID] = mapInternalAttendanceStatusToTeacher(record.Status)
	}

	summaries := make([]AttendanceStudentSummary, 0, len(enrollments))
	totalLessons := len(filteredLessons)
	for _, enrollment := range enrollments {
		summary := AttendanceStudentSummary{
			Student:      enrollment.Student,
			TotalLessons: totalLessons,
		}

		lessonStatusMap := studentLessonStatus[enrollment.StudentID]
		for _, status := range lessonStatusMap {
			summary.MarkedCount++
			switch status {
			case TeacherAttendanceStatusPresent:
				summary.PresentCount++
			case TeacherAttendanceStatusAbsent:
				summary.AbsentCount++
			case TeacherAttendanceStatusLate:
				summary.LateCount++
			case TeacherAttendanceStatusExcused:
				summary.ExcusedCount++
			}
		}

		if totalLessons > summary.MarkedCount {
			summary.UnmarkedCount = totalLessons - summary.MarkedCount
		}
		if totalLessons > 0 {
			summary.AttendanceRate = float64(summary.PresentCount+summary.LateCount) / float64(totalLessons)
		}

		summaries = append(summaries, summary)
	}

	return &GetAttendanceSummaryByStudentOutput{
		TeacherID:    teacher.ID,
		ClassID:      input.ClassID,
		TotalLessons: totalLessons,
		Students:     summaries,
	}, nil
}
