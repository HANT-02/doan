package predictive

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"doan/internal/entities"

	"gorm.io/gorm"
)

type TrainingSourceData struct {
	Students        []entities.Student
	Enrollments     []entities.Enrollment
	Lessons         []entities.Lesson
	Attendances     []entities.Attendance
	LessonSummaries []entities.LessonSummary
	AcademicRecords []entities.AcademicRecord
	LeaveRequests   []entities.LeaveRequest
}

type attendanceObservation struct {
	At     time.Time
	Status int
}

type academicObservation struct {
	At                time.Time
	HomeworkCompleted bool
	TotalScore        float64
	IsCompleted       bool
}

type leaveObservation struct {
	At      time.Time
	ClassID string
}

func LoadTrainingRowsFromDB(ctx context.Context, db *gorm.DB, definition DatasetDefinition) ([]TrainingRow, error) {
	if db == nil {
		return nil, errors.New("db connection is required")
	}

	data, err := LoadTrainingSourceDataFromDB(ctx, db)
	if err != nil {
		return nil, err
	}

	rows, err := BuildTrainingRowsFromSourceData(definition, data)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("khong tao duoc training row nao tu database; hay kiem tra du lieu attendance/academic_record/lesson")
	}

	return rows, nil
}

func LoadTrainingSourceDataFromDB(ctx context.Context, db *gorm.DB) (TrainingSourceData, error) {
	var students []entities.Student
	if err := db.WithContext(ctx).
		Where("status = ?", "ACTIVE").
		Find(&students).Error; err != nil {
		return TrainingSourceData{}, fmt.Errorf("load students: %w", err)
	}
	if len(students) == 0 {
		return TrainingSourceData{}, errors.New("khong tim thay student ACTIVE trong database")
	}

	studentIDs := collectStudentIDs(students)

	var enrollments []entities.Enrollment
	if err := db.WithContext(ctx).
		Preload("Class").
		Where("student_id IN ? AND status = ?", studentIDs, "ENROLLED").
		Find(&enrollments).Error; err != nil {
		return TrainingSourceData{}, fmt.Errorf("load enrollments: %w", err)
	}
	if len(enrollments) == 0 {
		return TrainingSourceData{}, errors.New("khong tim thay enrollment ENROLLED trong database")
	}

	classIDs := collectClassIDsFromEnrollments(enrollments)

	var lessons []entities.Lesson
	if err := db.WithContext(ctx).
		Where("class_id IN ?", classIDs).
		Order("date_start ASC").
		Find(&lessons).Error; err != nil {
		return TrainingSourceData{}, fmt.Errorf("load lessons: %w", err)
	}
	if len(lessons) == 0 {
		return TrainingSourceData{}, errors.New("khong tim thay lesson cho cac lop dang hoc; hay commit scheduling preview hoac seed lesson truoc khi train predictive tu DB")
	}

	lessonIDs := collectLessonIDs(lessons)

	var attendances []entities.Attendance
	if err := db.WithContext(ctx).
		Where("student_id IN ? AND lesson_id IN ?", studentIDs, lessonIDs).
		Find(&attendances).Error; err != nil {
		return TrainingSourceData{}, fmt.Errorf("load attendances: %w", err)
	}

	var lessonSummaries []entities.LessonSummary
	if err := db.WithContext(ctx).
		Where("lesson_id IN ?", lessonIDs).
		Find(&lessonSummaries).Error; err != nil {
		return TrainingSourceData{}, fmt.Errorf("load lesson summaries: %w", err)
	}

	summaryIDs := collectLessonSummaryIDs(lessonSummaries)

	var academicRecords []entities.AcademicRecord
	if len(summaryIDs) > 0 {
		if err := db.WithContext(ctx).
			Where("student_id IN ? AND lesson_summary_id IN ?", studentIDs, summaryIDs).
			Find(&academicRecords).Error; err != nil {
			return TrainingSourceData{}, fmt.Errorf("load academic records: %w", err)
		}
	}

	var leaveRequests []entities.LeaveRequest
	if err := db.WithContext(ctx).
		Where("student_id IN ? AND status = ?", studentIDs, "APPROVED").
		Find(&leaveRequests).Error; err != nil {
		return TrainingSourceData{}, fmt.Errorf("load leave requests: %w", err)
	}

	return TrainingSourceData{
		Students:        students,
		Enrollments:     enrollments,
		Lessons:         lessons,
		Attendances:     attendances,
		LessonSummaries: lessonSummaries,
		AcademicRecords: academicRecords,
		LeaveRequests:   leaveRequests,
	}, nil
}

func BuildTrainingRowsFromSourceData(definition DatasetDefinition, data TrainingSourceData) ([]TrainingRow, error) {
	if definition.ObservationWindowDays <= 0 {
		return nil, errors.New("definition observation window must be > 0")
	}
	if definition.PredictionHorizonDays <= 0 {
		return nil, errors.New("definition prediction horizon must be > 0")
	}

	studentsByID := make(map[string]entities.Student, len(data.Students))
	for _, student := range data.Students {
		studentsByID[student.ID] = student
	}

	enrollmentsByStudent := make(map[string][]entities.Enrollment)
	for _, enrollment := range data.Enrollments {
		if !strings.EqualFold(enrollment.Status, "ENROLLED") {
			continue
		}
		enrollmentsByStudent[enrollment.StudentID] = append(enrollmentsByStudent[enrollment.StudentID], enrollment)
	}

	lessonsByClass := make(map[string][]entities.Lesson)
	lessonByID := make(map[string]entities.Lesson, len(data.Lessons))
	for _, lesson := range data.Lessons {
		lessonsByClass[lesson.ClassID] = append(lessonsByClass[lesson.ClassID], lesson)
		lessonByID[lesson.ID] = lesson
	}
	for classID := range lessonsByClass {
		sort.Slice(lessonsByClass[classID], func(i, j int) bool {
			return lessonsByClass[classID][i].DateStart.Before(lessonsByClass[classID][j].DateStart)
		})
	}

	summaryByID := make(map[string]entities.LessonSummary, len(data.LessonSummaries))
	for _, summary := range data.LessonSummaries {
		summaryByID[summary.ID] = summary
	}

	attendanceByStudentClass := make(map[string][]attendanceObservation)
	for _, attendance := range data.Attendances {
		lesson, ok := lessonByID[attendance.LessonID]
		if !ok {
			continue
		}
		key := studentClassKey(attendance.StudentID, lesson.ClassID)
		attendanceByStudentClass[key] = append(attendanceByStudentClass[key], attendanceObservation{
			At:     lesson.DateStart,
			Status: attendance.Status,
		})
	}

	academicByStudentClass := make(map[string][]academicObservation)
	for _, record := range data.AcademicRecords {
		summary, ok := summaryByID[record.LessonSummaryID]
		if !ok {
			continue
		}
		lesson, ok := lessonByID[summary.LessonID]
		if !ok {
			continue
		}
		key := studentClassKey(record.StudentID, lesson.ClassID)
		academicByStudentClass[key] = append(academicByStudentClass[key], academicObservation{
			At:                lesson.DateStart,
			HomeworkCompleted: record.HomeworkCompleted,
			TotalScore:        record.TotalScore,
			IsCompleted:       record.IsCompleted,
		})
	}

	leaveByStudent := make(map[string][]leaveObservation)
	for _, leaveRequest := range data.LeaveRequests {
		if !strings.EqualFold(leaveRequest.Status, "APPROVED") {
			continue
		}
		classID := ""
		if leaveRequest.ClassID != nil {
			classID = *leaveRequest.ClassID
		}
		leaveByStudent[leaveRequest.StudentID] = append(leaveByStudent[leaveRequest.StudentID], leaveObservation{
			At:      leaveRequest.ApplyDate,
			ClassID: classID,
		})
	}

	dedupedRows := make(map[string]TrainingRow)
	for studentID, enrollments := range enrollmentsByStudent {
		student, ok := studentsByID[studentID]
		if !ok || !strings.EqualFold(student.Status, "ACTIVE") {
			continue
		}

		for _, enrollment := range enrollments {
			classLessons := lessonsByClass[enrollment.ClassID]
			if len(classLessons) == 0 {
				continue
			}

			classKey := studentClassKey(studentID, enrollment.ClassID)
			classAttendances := attendanceByStudentClass[classKey]
			classAcademic := academicByStudentClass[classKey]

			for _, snapshotLesson := range classLessons {
				snapshotAt := snapshotLesson.DateStart
				if !isEnrollmentActiveAt(enrollment, snapshotAt) {
					continue
				}

				observationStart := snapshotAt.AddDate(0, 0, -definition.ObservationWindowDays)
				futureEnd := snapshotAt.AddDate(0, 0, definition.PredictionHorizonDays)

				observationAttendance := filterAttendances(classAttendances, observationStart, snapshotAt)
				futureAttendance := filterFutureAttendances(classAttendances, snapshotAt, futureEnd)
				observationAcademic := completedAcademicRecords(filterAcademicRecords(classAcademic, observationStart, snapshotAt))
				futureAcademic := completedAcademicRecords(filterFutureAcademicRecords(classAcademic, snapshotAt, futureEnd))

				if len(futureAttendance) < definition.Label.MinimumFutureAttendanceRows &&
					len(futureAcademic) < definition.Label.MinimumFutureAcademicRows {
					continue
				}

				activeClassIDs := activeEnrollmentClassIDs(enrollmentsByStudent[studentID], snapshotAt)
				rowID := makeTrainingRowID(studentID, enrollment.ClassID, snapshotAt)
				dedupedRows[rowID] = TrainingRow{
					ID: rowID,
					Features: map[string]float64{
						"attendance_rate_28d":          observationAttendanceRate(observationAttendance),
						"absence_count_28d":            float64(countAbsent(observationAttendance)),
						"average_total_score_28d":      averageCompletedScore(observationAcademic),
						"homework_completion_rate_28d": homeworkCompletionRate(observationAcademic),
						"active_enrollment_count_28d":  float64(len(activeClassIDs)),
						"weekly_lesson_load_28d":       weeklyLessonLoad(lessonsByClass, activeClassIDs, observationStart, snapshotAt, definition.ObservationWindowDays),
						"approved_leave_count_28d":     float64(countApprovedLeaves(leaveByStudent[studentID], observationStart, snapshotAt, enrollment.ClassID)),
						"days_since_last_lesson":       daysSinceLastLesson(lessonsByClass, activeClassIDs, snapshotAt, definition.ObservationWindowDays),
					},
					Label: deriveFutureRiskLabel(futureAttendance, futureAcademic),
				}
			}
		}
	}

	rows := make([]TrainingRow, 0, len(dedupedRows))
	for _, row := range dedupedRows {
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].ID < rows[j].ID
	})

	return rows, nil
}

func collectStudentIDs(students []entities.Student) []string {
	ids := make([]string, 0, len(students))
	for _, student := range students {
		ids = append(ids, student.ID)
	}
	return ids
}

func collectClassIDsFromEnrollments(enrollments []entities.Enrollment) []string {
	seen := make(map[string]struct{})
	ids := make([]string, 0, len(enrollments))
	for _, enrollment := range enrollments {
		if _, ok := seen[enrollment.ClassID]; ok {
			continue
		}
		seen[enrollment.ClassID] = struct{}{}
		ids = append(ids, enrollment.ClassID)
	}
	return ids
}

func collectLessonIDs(lessons []entities.Lesson) []string {
	ids := make([]string, 0, len(lessons))
	for _, lesson := range lessons {
		ids = append(ids, lesson.ID)
	}
	return ids
}

func collectLessonSummaryIDs(summaries []entities.LessonSummary) []string {
	ids := make([]string, 0, len(summaries))
	for _, summary := range summaries {
		ids = append(ids, summary.ID)
	}
	return ids
}

func studentClassKey(studentID, classID string) string {
	return studentID + "::" + classID
}

func makeTrainingRowID(studentID, classID string, snapshotAt time.Time) string {
	return fmt.Sprintf("%s:%s:%s", studentID, classID, snapshotAt.Format("20060102"))
}

func isEnrollmentActiveAt(enrollment entities.Enrollment, snapshotAt time.Time) bool {
	if !strings.EqualFold(enrollment.Status, "ENROLLED") {
		return false
	}

	effectiveStart := enrollment.CreatedAt
	if enrollment.ApprovedAt != nil && enrollment.ApprovedAt.After(effectiveStart) {
		effectiveStart = *enrollment.ApprovedAt
	}
	if effectiveStart.After(snapshotAt) {
		return false
	}

	if !enrollment.Class.StartDate.IsZero() && enrollment.Class.StartDate.After(snapshotAt) {
		return false
	}
	if enrollment.Class.EndDate != nil && enrollment.Class.EndDate.Before(snapshotAt) {
		return false
	}

	return true
}

func filterAttendances(items []attendanceObservation, start, end time.Time) []attendanceObservation {
	filtered := make([]attendanceObservation, 0, len(items))
	for _, item := range items {
		if !item.At.Before(start) && item.At.Before(end) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func filterFutureAttendances(items []attendanceObservation, start, end time.Time) []attendanceObservation {
	filtered := make([]attendanceObservation, 0, len(items))
	for _, item := range items {
		if item.At.After(start) && item.At.Before(end) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func filterAcademicRecords(items []academicObservation, start, end time.Time) []academicObservation {
	filtered := make([]academicObservation, 0, len(items))
	for _, item := range items {
		if !item.At.Before(start) && item.At.Before(end) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func filterFutureAcademicRecords(items []academicObservation, start, end time.Time) []academicObservation {
	filtered := make([]academicObservation, 0, len(items))
	for _, item := range items {
		if item.At.After(start) && item.At.Before(end) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func completedAcademicRecords(items []academicObservation) []academicObservation {
	filtered := make([]academicObservation, 0, len(items))
	for _, item := range items {
		if item.IsCompleted {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func observationAttendanceRate(items []attendanceObservation) float64 {
	if len(items) == 0 {
		return 1.0
	}

	var attendedCount int
	for _, item := range items {
		if isAttendancePresent(item.Status) {
			attendedCount++
		}
	}
	return float64(attendedCount) / float64(len(items))
}

func countAbsent(items []attendanceObservation) int {
	var total int
	for _, item := range items {
		if item.Status == AttendanceStatusAbsent {
			total++
		}
	}
	return total
}

func averageCompletedScore(items []academicObservation) float64 {
	if len(items) == 0 {
		return 7.0
	}

	var total float64
	for _, item := range items {
		total += item.TotalScore
	}
	return total / float64(len(items))
}

func homeworkCompletionRate(items []academicObservation) float64 {
	if len(items) == 0 {
		return 1.0
	}

	var completedCount int
	for _, item := range items {
		if item.HomeworkCompleted {
			completedCount++
		}
	}
	return float64(completedCount) / float64(len(items))
}

func deriveFutureRiskLabel(attendance []attendanceObservation, academic []academicObservation) string {
	if len(attendance) > 0 && observationAttendanceRate(attendance) < 0.80 {
		return LabelAtRisk
	}
	if len(academic) > 0 && averageCompletedScore(academic) < 5.0 {
		return LabelAtRisk
	}
	if len(academic) > 0 && homeworkCompletionRate(academic) < 0.60 {
		return LabelAtRisk
	}
	return LabelNotAtRisk
}

func activeEnrollmentClassIDs(enrollments []entities.Enrollment, snapshotAt time.Time) []string {
	seen := make(map[string]struct{})
	classIDs := make([]string, 0, len(enrollments))
	for _, enrollment := range enrollments {
		if !isEnrollmentActiveAt(enrollment, snapshotAt) {
			continue
		}
		if _, ok := seen[enrollment.ClassID]; ok {
			continue
		}
		seen[enrollment.ClassID] = struct{}{}
		classIDs = append(classIDs, enrollment.ClassID)
	}
	return classIDs
}

func weeklyLessonLoad(lessonsByClass map[string][]entities.Lesson, classIDs []string, start, end time.Time, observationWindowDays int) float64 {
	if len(classIDs) == 0 || observationWindowDays <= 0 {
		return 0
	}

	var lessonCount int
	for _, classID := range classIDs {
		for _, lesson := range lessonsByClass[classID] {
			if !lesson.DateStart.Before(start) && lesson.DateStart.Before(end) {
				lessonCount++
			}
		}
	}

	weeks := float64(observationWindowDays) / 7.0
	if weeks == 0 {
		return 0
	}
	return float64(lessonCount) / weeks
}

func countApprovedLeaves(items []leaveObservation, start, end time.Time, classID string) int {
	var total int
	for _, item := range items {
		if item.ClassID != "" && item.ClassID != classID {
			continue
		}
		if !item.At.Before(start) && item.At.Before(end) {
			total++
		}
	}
	return total
}

func daysSinceLastLesson(lessonsByClass map[string][]entities.Lesson, classIDs []string, snapshotAt time.Time, observationWindowDays int) float64 {
	var lastLesson *time.Time
	for _, classID := range classIDs {
		for _, lesson := range lessonsByClass[classID] {
			if !lesson.DateStart.Before(snapshotAt) {
				continue
			}
			lessonDate := lesson.DateStart
			if lastLesson == nil || lessonDate.After(*lastLesson) {
				lastLesson = &lessonDate
			}
		}
	}
	if lastLesson == nil {
		return float64(observationWindowDays + 1)
	}
	return snapshotAt.Sub(*lastLesson).Hours() / 24
}

func isAttendancePresent(status int) bool {
	return status == AttendanceStatusPresent || status == AttendanceStatusLate || status == AttendanceStatusEarly
}
