package predictive

import (
	"testing"
	"time"

	"doan/internal/entities"
)

func TestBuildTrainingRowsFromSourceData_BuildsSnapshotFeaturesFromEntitySlices(t *testing.T) {
	t.Parallel()

	baseDate := time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC)
	classID := "class-1"
	studentID := "student-1"

	lessons := make([]entities.Lesson, 0, 9)
	summaries := make([]entities.LessonSummary, 0, 9)
	attendances := make([]entities.Attendance, 0, 9)
	records := make([]entities.AcademicRecord, 0, 9)

	attendanceStatuses := []int{
		AttendanceStatusPresent,
		AttendanceStatusAbsent,
		AttendanceStatusPresent,
		AttendanceStatusPresent,
		AttendanceStatusPresent,
		AttendanceStatusAbsent,
		AttendanceStatusPresent,
		AttendanceStatusAbsent,
		AttendanceStatusPresent,
	}
	scores := []float64{6.0, 7.0, 8.0, 5.0, 7.0, 4.0, 6.0, 4.0, 6.0}
	homeworkCompleted := []bool{true, true, false, true, true, false, true, false, true}

	for index := 0; index < 9; index++ {
		lessonID := "lesson-" + string(rune('1'+index))
		summaryID := "summary-" + string(rune('1'+index))
		lessonDate := baseDate.AddDate(0, 0, index*7)

		lessons = append(lessons, entities.Lesson{
			ID:        lessonID,
			ClassID:   classID,
			DateStart: lessonDate,
			DateEnd:   lessonDate.Add(90 * time.Minute),
		})
		summaries = append(summaries, entities.LessonSummary{
			ID:       summaryID,
			LessonID: lessonID,
		})
		attendances = append(attendances, entities.Attendance{
			ID:        "attendance-" + lessonID,
			LessonID:  lessonID,
			StudentID: studentID,
			Status:    attendanceStatuses[index],
		})
		records = append(records, entities.AcademicRecord{
			ID:                "record-" + lessonID,
			LessonSummaryID:   summaryID,
			StudentID:         studentID,
			HomeworkCompleted: homeworkCompleted[index],
			TotalScore:        scores[index],
			IsCompleted:       true,
		})
	}

	classStart := baseDate.AddDate(0, 0, -7)
	rows, err := BuildTrainingRowsFromSourceData(DefaultAtRiskDatasetDefinition(), TrainingSourceData{
		Students: []entities.Student{
			{
				ID:     studentID,
				Status: "ACTIVE",
			},
		},
		Enrollments: []entities.Enrollment{
			{
				ID:        "enrollment-1",
				ClassID:   classID,
				StudentID: studentID,
				Status:    "ENROLLED",
				CreatedAt: classStart,
				Class: entities.Class{
					ID:        classID,
					StartDate: classStart,
				},
			},
		},
		Lessons:         lessons,
		Attendances:     attendances,
		LessonSummaries: summaries,
		AcademicRecords: records,
		LeaveRequests: []entities.LeaveRequest{
			{
				ID:        "leave-1",
				StudentID: studentID,
				ClassID:   stringPtr(classID),
				Status:    "APPROVED",
				ApplyDate: baseDate.AddDate(0, 0, 14),
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snapshotID := makeTrainingRowID(studentID, classID, baseDate.AddDate(0, 0, 28))
	row, ok := findTrainingRow(rows, snapshotID)
	if !ok {
		t.Fatalf("expected snapshot row %s, got %d rows", snapshotID, len(rows))
	}

	if row.Label != LabelAtRisk {
		t.Fatalf("expected AT_RISK label, got %s", row.Label)
	}

	assertFloatEquals(t, row.Features["attendance_rate_28d"], 0.75)
	assertFloatEquals(t, row.Features["absence_count_28d"], 1.0)
	assertFloatEquals(t, row.Features["average_total_score_28d"], 6.5)
	assertFloatEquals(t, row.Features["homework_completion_rate_28d"], 0.75)
	assertFloatEquals(t, row.Features["active_enrollment_count_28d"], 1.0)
	assertFloatEquals(t, row.Features["weekly_lesson_load_28d"], 1.0)
	assertFloatEquals(t, row.Features["approved_leave_count_28d"], 1.0)
	assertFloatEquals(t, row.Features["days_since_last_lesson"], 7.0)
}

func TestBuildTrainingRowsFromSourceData_SkipsSnapshotsWithoutFutureEvidence(t *testing.T) {
	t.Parallel()

	baseDate := time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC)
	classID := "class-1"
	studentID := "student-1"

	lessons := make([]entities.Lesson, 0, 8)
	summaries := make([]entities.LessonSummary, 0, 8)
	attendances := make([]entities.Attendance, 0, 8)
	for index := 0; index < 8; index++ {
		lessonID := "lesson-short-" + string(rune('1'+index))
		summaryID := "summary-short-" + string(rune('1'+index))
		lessonDate := baseDate.AddDate(0, 0, index*7)
		lessons = append(lessons, entities.Lesson{
			ID:        lessonID,
			ClassID:   classID,
			DateStart: lessonDate,
			DateEnd:   lessonDate.Add(90 * time.Minute),
		})
		summaries = append(summaries, entities.LessonSummary{ID: summaryID, LessonID: lessonID})
		attendances = append(attendances, entities.Attendance{
			ID:        "attendance-" + lessonID,
			LessonID:  lessonID,
			StudentID: studentID,
			Status:    AttendanceStatusPresent,
		})
	}

	rows, err := BuildTrainingRowsFromSourceData(DefaultAtRiskDatasetDefinition(), TrainingSourceData{
		Students: []entities.Student{{ID: studentID, Status: "ACTIVE"}},
		Enrollments: []entities.Enrollment{{
			ID:        "enrollment-1",
			ClassID:   classID,
			StudentID: studentID,
			Status:    "ENROLLED",
			CreatedAt: baseDate.AddDate(0, 0, -7),
			Class: entities.Class{
				ID:        classID,
				StartDate: baseDate.AddDate(0, 0, -7),
			},
		}},
		Lessons:         lessons,
		Attendances:     attendances,
		LessonSummaries: summaries,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snapshotID := makeTrainingRowID(studentID, classID, baseDate.AddDate(0, 0, 28))
	if _, ok := findTrainingRow(rows, snapshotID); ok {
		t.Fatalf("did not expect snapshot %s because future attendance rows < 4", snapshotID)
	}
}

func findTrainingRow(rows []TrainingRow, rowID string) (TrainingRow, bool) {
	for _, row := range rows {
		if row.ID == rowID {
			return row, true
		}
	}
	return TrainingRow{}, false
}

func stringPtr(value string) *string {
	return &value
}

func assertFloatEquals(t *testing.T, actual, expected float64) {
	t.Helper()

	const epsilon = 0.0001
	diff := actual - expected
	if diff < 0 {
		diff = -diff
	}
	if diff > epsilon {
		t.Fatalf("expected %.4f, got %.4f", expected, actual)
	}
}
