package scheduling

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	schedulingservice "doan/internal/services/scheduling"
	"doan/pkg/logger"

	"github.com/lib/pq"
)

func TestPreviewAndCommitUseCase_CommitAssignmentsWithDefaultSolver(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{}
	classScheduleRepo := &previewClassScheduleRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{}
	uow := &previewUnitOfWorkStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)
	commitUseCase := NewCommitPreviewUseCase(
		lessonRepo,
		classScheduleRepo,
		uow,
		logger.NewZapLogger(logger.Config{Level: "error", Format: "json", Output: "stdout", ServiceName: "test", Environment: "test"}),
		store,
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	if result.Status != "COMPLETED" {
		t.Fatalf("expected preview COMPLETED, got %s", result.Status)
	}
	if len(result.Assignments) != 2 {
		t.Fatalf("expected 2 assignments, got %d", len(result.Assignments))
	}

	output, err := commitUseCase.Execute(context.Background(), CommitPreviewInput{RunID: result.RunID})
	if err != nil {
		t.Fatalf("unexpected commit error: %v", err)
	}

	if output.Status != "COMMITTED" {
		t.Fatalf("expected COMMITTED status, got %s", output.Status)
	}
	if output.ScheduledLessons != 2 {
		t.Fatalf("expected 2 committed lessons, got %d", output.ScheduledLessons)
	}
	if len(lessonRepo.createdLessons) != 2 {
		t.Fatalf("expected 2 lessons created in repository, got %d", len(lessonRepo.createdLessons))
	}
	for _, lesson := range lessonRepo.createdLessons {
		if lesson.Status != entities.LessonStatusPublished {
			t.Fatalf("expected committed lesson status %s, got %s", entities.LessonStatusPublished, lesson.Status)
		}
		if lesson.PublishedAt == nil {
			t.Fatalf("expected committed lesson to have published_at")
		}
		if lesson.SourcePreviewRun == nil || *lesson.SourcePreviewRun != result.RunID {
			t.Fatalf("expected committed lesson source_preview_run_id to be %s", result.RunID)
		}
		if lesson.ChangeReason != entities.LessonChangeReasonInitialSchedulingCommit {
			t.Fatalf("expected committed lesson change_reason %s, got %s", entities.LessonChangeReasonInitialSchedulingCommit, lesson.ChangeReason)
		}
	}
	if !uow.beginCalled || !uow.commitCalled || uow.rollbackCalled {
		t.Fatalf("unexpected transaction lifecycle: begin=%t commit=%t rollback=%t", uow.beginCalled, uow.commitCalled, uow.rollbackCalled)
	}
}

func TestPreviewUseCase_NormalizesEffectiveDateFromForPastRequest(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositoriesWithSessionCount(2)
	lessonRepo := &previewLessonRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	today := time.Now().UTC()
	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: today.AddDate(0, 0, -10),
		DateTo:   today.AddDate(0, 0, 7),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	expected := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	if !sameCalendarDay(result.EffectiveDateFrom, expected) {
		t.Fatalf("expected effective_date_from day %s, got %s", expected.Format(time.RFC3339), result.EffectiveDateFrom.Format(time.RFC3339))
	}
	if !sameCalendarDay(result.Filters.EffectiveDateFrom, expected) {
		t.Fatalf("expected filters.effective_date_from day %s, got %s", expected.Format(time.RFC3339), result.Filters.EffectiveDateFrom.Format(time.RFC3339))
	}
	for _, assignment := range result.Assignments {
		if assignment.StartTime.Before(expected) {
			t.Fatalf("did not expect assignment in the past: %s", assignment.StartTime.Format(time.RFC3339))
		}
	}
}

func TestCommitPreviewUseCase_ReturnsConflictWhenExistingLessonOverlaps(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{}
	classScheduleRepo := &previewClassScheduleRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{}
	uow := &previewUnitOfWorkStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)
	commitUseCase := NewCommitPreviewUseCase(
		lessonRepo,
		classScheduleRepo,
		uow,
		logger.NewZapLogger(logger.Config{Level: "error", Format: "json", Output: "stdout", ServiceName: "test", Environment: "test"}),
		store,
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	if len(result.Assignments) == 0 {
		t.Fatalf("expected assignments to create a conflict case")
	}

	firstAssignment := result.Assignments[0]
	lessonRepo.existingLessons = []entities.Lesson{
		{
			ID:        "lesson-existing-1",
			ClassID:   "class-existing",
			TeacherID: &firstAssignment.TeacherID,
			RoomID:    &firstAssignment.RoomID,
			DateStart: firstAssignment.StartTime.Add(-15 * time.Minute),
			DateEnd:   firstAssignment.EndTime.Add(15 * time.Minute),
			Status:    entities.LessonStatusPublished,
		},
	}

	_, err = commitUseCase.Execute(context.Background(), CommitPreviewInput{RunID: result.RunID})
	if err == nil {
		t.Fatalf("expected conflict error")
	}

	if !strings.Contains(err.Error(), "Khong the commit preview") {
		t.Fatalf("expected formatted conflict message, got %v", err)
	}
	if len(lessonRepo.createdLessons) != 0 {
		t.Fatalf("expected no new lessons to be created when conflict occurs")
	}
	storedPreview, ok := store.Get(result.RunID)
	if !ok {
		t.Fatalf("expected preview to remain available in store after conflict")
	}
	if storedPreview.Status != "PARTIAL" {
		t.Fatalf("expected stored preview status to become PARTIAL after conflict, got %s", storedPreview.Status)
	}
	if len(storedPreview.ExistingLessons) == 0 {
		t.Fatalf("expected stored preview to include existing lessons after conflict")
	}
	hasSystemConflict := false
	for _, conflict := range storedPreview.Conflicts {
		if conflict.Type == "SYSTEM_LESSON_CONFLICT" {
			hasSystemConflict = true
			break
		}
	}
	if !hasSystemConflict {
		t.Fatalf("expected stored preview to include SYSTEM_LESSON_CONFLICT after commit failure")
	}
	if !uow.beginCalled || uow.commitCalled || !uow.rollbackCalled {
		t.Fatalf("unexpected transaction lifecycle: begin=%t commit=%t rollback=%t", uow.beginCalled, uow.commitCalled, uow.rollbackCalled)
	}
}

func TestCommitPreviewUseCase_AcceptsCandidateOptionKeyFromPreview(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{}
	classScheduleRepo := &previewClassScheduleRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{}
	uow := &previewUnitOfWorkStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)
	commitUseCase := NewCommitPreviewUseCase(
		lessonRepo,
		classScheduleRepo,
		uow,
		logger.NewZapLogger(logger.Config{Level: "error", Format: "json", Output: "stdout", ServiceName: "test", Environment: "test"}),
		store,
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}
	if len(result.Assignments) == 0 {
		t.Fatalf("expected assignments")
	}

	assignment := result.Assignments[0]
	options := result.CandidateOptions[assignment.VariableID]
	var selectedKey string
	for _, option := range options {
		if option.RoomID == assignment.RoomID && option.StartTime.Equal(assignment.StartTime) && option.EndTime.Equal(assignment.EndTime) {
			selectedKey = option.Key
			break
		}
	}
	if selectedKey == "" {
		t.Fatalf("expected preview candidate option matching assignment")
	}

	_, err = commitUseCase.Execute(context.Background(), CommitPreviewInput{
		RunID: result.RunID,
		ManualAssignments: []ManualAssignmentOverride{
			{VariableID: assignment.VariableID, OptionKey: selectedKey},
		},
	})
	if err != nil {
		t.Fatalf("expected preview candidate option key to be accepted by commit, got %v", err)
	}
}

func TestPreviewUseCase_IncludesExistingLessonConflictBeforeCommit(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	firstPreview, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}
	if len(firstPreview.Assignments) == 0 {
		t.Fatalf("expected assignments to create a conflict fixture")
	}

	firstAssignment := firstPreview.Assignments[0]
	lessonRepo.existingLessons = []entities.Lesson{
		{
			ID:        "lesson-existing-preview",
			ClassID:   "class-other",
			TeacherID: &firstAssignment.TeacherID,
			RoomID:    &firstAssignment.RoomID,
			DateStart: firstAssignment.StartTime,
			DateEnd:   firstAssignment.EndTime,
			Status:    entities.LessonStatusPublished,
		},
	}

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	if result.Status != "PARTIAL" {
		t.Fatalf("expected preview status PARTIAL with existing lesson conflict, got %s", result.Status)
	}
	if len(result.ExistingLessons) == 0 {
		t.Fatalf("expected preview to include existing lesson cards")
	}
	hasSystemConflict := false
	for _, conflict := range result.Conflicts {
		if conflict.Type == "SYSTEM_LESSON_CONFLICT" {
			hasSystemConflict = true
			break
		}
	}
	if !hasSystemConflict {
		t.Fatalf("expected preview to include SYSTEM_LESSON_CONFLICT before commit")
	}
}

func TestPreviewUseCase_ReturnsSkillMismatchConflict(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositoriesWithSkills(2, []string{"TESOL"}, []string{"IELTS_8.0", "TESOL"})
	lessonRepo := &previewLessonRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	if len(result.Assignments) != 0 {
		t.Fatalf("expected no assignments when teacher lacks required skills, got %d", len(result.Assignments))
	}

	hasSkillMismatch := false
	for _, conflict := range result.Conflicts {
		if conflict.Type == "SKILL_MISMATCH" {
			hasSkillMismatch = true
			if !strings.Contains(conflict.Message, "IELTS_8.0") {
				t.Fatalf("expected missing required skill in conflict message, got %q", conflict.Message)
			}
		}
	}

	if !hasSkillMismatch {
		t.Fatalf("expected preview conflicts to include SKILL_MISMATCH")
	}
}

func TestPreviewUseCase_ColdStartIgnoresPublishedLessons(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{
		existingLessons: []entities.Lesson{
			{
				ID:        "lesson-published-1",
				ClassID:   "class-other",
				TeacherID: stringPtr("teacher-1"),
				RoomID:    stringPtr("room-1"),
				DateStart: time.Date(2026, 6, 15, 8, 0, 0, 0, time.UTC),
				DateEnd:   time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC),
				Status:    entities.LessonStatusPublished,
			},
		},
	}
	enrollmentRepo := &previewEnrollmentRepoStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC),
		Mode:     schedulingservice.PreviewModeColdStart,
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	if len(result.ExistingLessons) != 0 {
		t.Fatalf("expected cold_start to ignore existing published lessons, got %d", len(result.ExistingLessons))
	}
	for _, conflict := range result.Conflicts {
		if conflict.Type == "SYSTEM_LESSON_CONFLICT" {
			t.Fatalf("did not expect system lesson conflict in cold_start mode")
		}
	}
}

func TestPreviewUseCase_ReplanWithPublishedLockIgnoresHistoryButBlocksPublished(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{
		existingLessons: []entities.Lesson{
			{
				ID:        "lesson-history-1",
				ClassID:   "class-other-history",
				TeacherID: stringPtr("teacher-1"),
				RoomID:    stringPtr("room-1"),
				DateStart: time.Date(2026, 6, 15, 8, 0, 0, 0, time.UTC),
				DateEnd:   time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC),
				Status:    entities.LessonStatusHistory,
			},
			{
				ID:        "lesson-published-1",
				ClassID:   "class-other-published",
				TeacherID: stringPtr("teacher-1"),
				RoomID:    stringPtr("room-1"),
				DateStart: time.Date(2026, 6, 15, 8, 0, 0, 0, time.UTC),
				DateEnd:   time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC),
				Status:    entities.LessonStatusPublished,
			},
		},
	}
	enrollmentRepo := &previewEnrollmentRepoStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC),
		Mode:     schedulingservice.PreviewModeReplanWithPublishedLock,
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	if len(result.ExistingLessons) != 1 {
		t.Fatalf("expected only published lesson to be loaded as lock, got %d", len(result.ExistingLessons))
	}
	if result.ExistingLessons[0].Status != entities.LessonStatusPublished {
		t.Fatalf("expected loaded existing lesson status %s, got %s", entities.LessonStatusPublished, result.ExistingLessons[0].Status)
	}

	hasPublishedConflict := false
	for _, conflict := range result.Conflicts {
		if conflict.Type == "SYSTEM_LESSON_CONFLICT" {
			hasPublishedConflict = true
			if !strings.Contains(conflict.Message, "Published") {
				t.Fatalf("expected published lifecycle label in conflict message, got %q", conflict.Message)
			}
		}
		if strings.Contains(conflict.Message, "History") {
			t.Fatalf("did not expect history lesson to appear in conflict message")
		}
	}
	if !hasPublishedConflict {
		t.Fatalf("expected replan_with_published_lock to produce published lesson conflict")
	}
}

func TestCommitPreviewUseCase_FormatsPublishedConflictLifecycle(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{}
	classScheduleRepo := &previewClassScheduleRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{}
	uow := &previewUnitOfWorkStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)
	commitUseCase := NewCommitPreviewUseCase(
		lessonRepo,
		classScheduleRepo,
		uow,
		logger.NewZapLogger(logger.Config{Level: "error", Format: "json", Output: "stdout", ServiceName: "test", Environment: "test"}),
		store,
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	firstAssignment := result.Assignments[0]
	lessonRepo.existingLessons = []entities.Lesson{
		{
			ID:        "lesson-published-commit",
			ClassID:   "class-existing",
			TeacherID: &firstAssignment.TeacherID,
			RoomID:    &firstAssignment.RoomID,
			DateStart: firstAssignment.StartTime.Add(-15 * time.Minute),
			DateEnd:   firstAssignment.EndTime.Add(15 * time.Minute),
			Status:    entities.LessonStatusPublished,
		},
	}

	_, err = commitUseCase.Execute(context.Background(), CommitPreviewInput{RunID: result.RunID})
	if err == nil {
		t.Fatalf("expected commit conflict error")
	}
	if !strings.Contains(err.Error(), "lesson Published dang ton tai") {
		t.Fatalf("expected commit conflict message to mention published lifecycle, got %v", err)
	}
}

func TestPreviewUseCase_ComputesRequestedSessionsFromCourseSessionCount(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositoriesWithSessionCount(8)
	lessonRepo := &previewLessonRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 10, 16, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	if result.Status != "COMPLETED" {
		t.Fatalf("expected COMPLETED preview status, got %s", result.Status)
	}
	if result.Summary.RequestedSessions != 8 {
		t.Fatalf("expected requested sessions to match course session_count=8, got %d", result.Summary.RequestedSessions)
	}
	if result.Summary.ScheduledLessons != 8 {
		t.Fatalf("expected exactly 8 scheduled lessons from projected class window, got %d", result.Summary.ScheduledLessons)
	}
	if result.Summary.UnscheduledLessons != 0 {
		t.Fatalf("expected no unscheduled lessons, got %d", result.Summary.UnscheduledLessons)
	}
	if len(result.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(result.Conflicts))
	}
	if result.Filters.DateTo.Before(time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected projected preview window to extend beyond manual date_to, got %s", result.Filters.DateTo.Format(time.RFC3339))
	}
}

func TestPreviewUseCase_FiltersClassesWithInsufficientEnrollment(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{
		byClass: map[string][]entities.Enrollment{
			"class-1": {
				{ID: "enr-1", ClassID: "class-1", StudentID: "stu-1", Status: "APPROVED"},
				{ID: "enr-2", ClassID: "class-1", StudentID: "stu-2", Status: "APPROVED"},
				{ID: "enr-3", ClassID: "class-1", StudentID: "stu-3", Status: "APPROVED"},
			},
		},
	}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	if result.Summary.ScheduledLessons != 0 {
		t.Fatalf("expected no scheduled lessons for under-enrolled class, got %d", result.Summary.ScheduledLessons)
	}

	hasEnrollmentConflict := false
	for _, conflict := range result.Conflicts {
		if conflict.Type == "INSUFFICIENT_ENROLLMENT" {
			hasEnrollmentConflict = true
			break
		}
	}
	if !hasEnrollmentConflict {
		t.Fatalf("expected INSUFFICIENT_ENROLLMENT conflict")
	}
}

func TestPreviewUseCase_CountsEnrolledStudentsForEnrollmentThreshold(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{}

	enrollments := make([]entities.Enrollment, 0, 18)
	for i := 0; i < 18; i++ {
		enrollments = append(enrollments, entities.Enrollment{
			ID:        fmt.Sprintf("enr-%d", i+1),
			ClassID:   "class-1",
			StudentID: fmt.Sprintf("stu-%d", i+1),
			Status:    entities.EnrollmentStatusEnrolled,
		})
	}

	enrollmentRepo := &previewEnrollmentRepoStub{
		byClass: map[string][]entities.Enrollment{
			"class-1": enrollments,
		},
	}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
		ClassIDs: []string{"class-1"},
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	for _, conflict := range result.Conflicts {
		if conflict.Type == "INSUFFICIENT_ENROLLMENT" {
			t.Fatalf("did not expect INSUFFICIENT_ENROLLMENT conflict when class has 18 ENROLLED students: %s", conflict.Message)
		}
	}

	if result.Summary.RequestedClasses == 0 {
		t.Fatalf("expected preview to keep class eligible for scheduling")
	}
}

func TestCommitPreviewUseCase_BlocksExcessiveManualAdjustments(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositoriesWithSessionCount(8)
	lessonRepo := &previewLessonRepoStub{}
	classScheduleRepo := &previewClassScheduleRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{
		byClass: map[string][]entities.Enrollment{
			"class-1": buildApprovedEnrollments("class-1", 18),
		},
	}
	uow := &previewUnitOfWorkStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)
	commitUseCase := NewCommitPreviewUseCase(
		lessonRepo,
		classScheduleRepo,
		uow,
		logger.NewZapLogger(logger.Config{Level: "error", Format: "json", Output: "stdout", ServiceName: "test", Environment: "test"}),
		store,
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 10, 9, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}
	manualThreshold := allowedManualAdjustmentLimit(result.Summary.RequestedSessions)
	requiredOverrides := manualThreshold + 1
	if len(result.Assignments) < requiredOverrides {
		t.Fatalf("expected at least %d assignments to test manual adjustment threshold, got %d", requiredOverrides, len(result.Assignments))
	}

	manualAssignments := make([]ManualAssignmentOverride, 0, requiredOverrides)
	for index, assignment := range result.Assignments[:requiredOverrides] {
		options := result.CandidateOptions[assignment.VariableID]
		var currentKey string
		for _, option := range options {
			if option.RoomID == assignment.RoomID && option.StartTime.Equal(assignment.StartTime) && option.EndTime.Equal(assignment.EndTime) {
				currentKey = option.Key
				break
			}
		}
		if currentKey == "" {
			t.Fatalf("expected current candidate key for assignment %d", index)
		}
		manualAssignments = append(manualAssignments, ManualAssignmentOverride{
			VariableID: assignment.VariableID,
			OptionKey:  currentKey,
		})
	}

	_, err = commitUseCase.Execute(context.Background(), CommitPreviewInput{
		RunID:             result.RunID,
		ManualAssignments: manualAssignments,
	})
	if err == nil {
		t.Fatalf("expected excessive manual adjustment error")
	}

	var conflictErr *CommitPreviewConflictError
	if !strings.Contains(err.Error(), "chinh tay") && !strings.Contains(err.Error(), "commit") && !strings.Contains(err.Error(), "nguong") {
		t.Fatalf("expected commit error to mention manual adjustment threshold, got %v", err)
	}
	if !errors.As(err, &conflictErr) {
		t.Fatalf("expected CommitPreviewConflictError")
	}

	hasExcessiveConflict := false
	for _, conflict := range conflictErr.Preview.Conflicts {
		if conflict.Type == "EXCESSIVE_MANUAL_ADJUSTMENT" {
			hasExcessiveConflict = true
			break
		}
	}
	if !hasExcessiveConflict {
		t.Fatalf("expected EXCESSIVE_MANUAL_ADJUSTMENT conflict in preview")
	}
}

func buildApprovedEnrollments(classID string, count int) []entities.Enrollment {
	items := make([]entities.Enrollment, 0, count)
	now := time.Now()
	for index := 0; index < count; index++ {
		items = append(items, entities.Enrollment{
			ID:         fmt.Sprintf("enr-%s-%02d", classID, index+1),
			ClassID:    classID,
			StudentID:  fmt.Sprintf("student-%02d", index+1),
			Status:     "APPROVED",
			ApprovedAt: &now,
		})
	}
	return items
}

func previewFixtureRepositories() (repositoryinterface.ClassRepository, repositoryinterface.RoomRepository, repositoryinterface.ShiftRepository) {
	return previewFixtureRepositoriesWithSkills(2, []string{"MATH_CORE"}, []string{"MATH_CORE"})
}

func previewFixtureRepositoriesWithSessionCount(sessionCount int) (repositoryinterface.ClassRepository, repositoryinterface.RoomRepository, repositoryinterface.ShiftRepository) {
	return previewFixtureRepositoriesWithSkills(sessionCount, []string{"MATH_CORE"}, []string{"MATH_CORE"})
}

func previewFixtureRepositoriesWithSkills(
	sessionCount int,
	teacherSkills []string,
	requiredSkills []string,
) (repositoryinterface.ClassRepository, repositoryinterface.RoomRepository, repositoryinterface.ShiftRepository) {
	teacherID := "teacher-1"
	courseID := "course-1"
	roomID := "room-1"
	shiftMonday := entities.Shift{
		ID:          "shift-1",
		Code:        "S1",
		Name:        "Ca sang",
		StartTime:   "08:00",
		EndTime:     "10:00",
		SessionType: "MORNING",
		IsActive:    true,
	}
	shiftWednesday := entities.Shift{
		ID:          "shift-2",
		Code:        "S2",
		Name:        "Ca chieu",
		StartTime:   "13:30",
		EndTime:     "15:30",
		SessionType: "AFTERNOON",
		IsActive:    true,
	}

	classRepo := benchmarkClassRepoStub{
		data: []entities.Class{
			{
				ID:          "class-1",
				Code:        "L-001",
				Name:        "Lop Toan 1",
				MaxStudents: 20,
				Status:      "OPEN",
				TeacherID:   &teacherID,
				Teacher: entities.Teacher{
					ID:       teacherID,
					Code:     "GV-001",
					FullName: "Giao vien 1",
					Skills:   pq.StringArray(append([]string(nil), teacherSkills...)),
				},
				CourseID: &courseID,
				Course: entities.Course{
					ID:                     courseID,
					Code:                   "KH-001",
					Name:                   "Khoa hoc Toan",
					SessionCount:           sessionCount,
					SessionDurationMinutes: 120,
					RequiredSkills:         pq.StringArray(append([]string(nil), requiredSkills...)),
				},
				RoomID: &roomID,
				ClassSchedules: []entities.ClassSchedule{
					{
						ID:        "schedule-1",
						ClassID:   "class-1",
						DayOfWeek: "MONDAY",
						ShiftID:   shiftMonday.ID,
						Shift:     shiftMonday,
						RoomID:    &roomID,
					},
					{
						ID:        "schedule-2",
						ClassID:   "class-1",
						DayOfWeek: "WEDNESDAY",
						ShiftID:   shiftWednesday.ID,
						Shift:     shiftWednesday,
						RoomID:    &roomID,
					},
				},
			},
		},
	}

	roomRepo := benchmarkRoomRepoStub{
		data: []entities.Room{
			{ID: roomID, Code: "P-001", Name: "Phong 1", Capacity: 30},
		},
	}

	shiftRepo := benchmarkShiftRepoStub{
		data: []entities.Shift{shiftMonday, shiftWednesday},
	}

	return classRepo, roomRepo, shiftRepo
}

var _ repositoryinterface.LessonRepository = (*previewLessonRepoStub)(nil)

type previewLessonRepoStub struct {
	existingLessons []entities.Lesson
	createdLessons  []entities.Lesson
}

func (s *previewLessonRepoStub) GetTable() string { return "lessons" }

func (s *previewLessonRepoStub) GetByCondition(_ context.Context, _ *repositories.CommonCondition) (*repositories.Pagination[entities.Lesson], error) {
	return &repositories.Pagination[entities.Lesson]{}, nil
}

func (s *previewLessonRepoStub) GetTotal(_ context.Context, _ *repositories.CommonCondition) (uint64, error) {
	return uint64(len(s.createdLessons)), nil
}

func (s *previewLessonRepoStub) Create(_ context.Context, entity *entities.Lesson) (*entities.Lesson, error) {
	s.createdLessons = append(s.createdLessons, *entity)
	return entity, nil
}

func (s *previewLessonRepoStub) Update(_ context.Context, id interface{}, updatedData map[string]interface{}) error {
	for index := range s.existingLessons {
		if s.existingLessons[index].ID != fmt.Sprint(id) {
			continue
		}
		if status, ok := updatedData["status"].(string); ok {
			s.existingLessons[index].Status = status
		}
		if changeReason, ok := updatedData["change_reason"].(string); ok {
			s.existingLessons[index].ChangeReason = changeReason
		}
	}
	return nil
}

func (s *previewLessonRepoStub) UpdateWithIDs(_ context.Context, _ []string, _ map[string]interface{}) error {
	return nil
}

func (s *previewLessonRepoStub) SoftDelete(_ context.Context, _ interface{}) error {
	return nil
}

func (s *previewLessonRepoStub) HardDelete(_ context.Context, _ interface{}) error {
	return nil
}

func (s *previewLessonRepoStub) GetByID(_ context.Context, _ interface{}) (*entities.Lesson, error) {
	return nil, nil
}

func (s *previewLessonRepoStub) ListInRange(_ context.Context, from time.Time, to time.Time) ([]entities.Lesson, error) {
	filtered := make([]entities.Lesson, 0, len(s.existingLessons))
	for _, lesson := range s.existingLessons {
		if lesson.DateStart.Before(to) && lesson.DateEnd.After(from) {
			filtered = append(filtered, lesson)
		}
	}
	return filtered, nil
}

func (s *previewLessonRepoStub) FindOverlappingLessons(
	_ context.Context,
	_ time.Time,
	_ time.Time,
	_ []string,
	_ []string,
	_ []string,
	statuses []string,
) ([]entities.Lesson, error) {
	if len(statuses) == 0 {
		return append([]entities.Lesson(nil), s.existingLessons...), nil
	}

	statusSet := make(map[string]struct{}, len(statuses))
	for _, status := range statuses {
		statusSet[status] = struct{}{}
	}

	filtered := make([]entities.Lesson, 0, len(s.existingLessons))
	for _, lesson := range s.existingLessons {
		if _, ok := statusSet[lesson.Status]; !ok {
			continue
		}
		filtered = append(filtered, lesson)
	}

	return filtered, nil
}

func (s *previewLessonRepoStub) GetLessonWithRelations(_ context.Context, _ string) (*entities.Lesson, error) {
	return nil, nil
}

var _ repositoryinterface.ClassScheduleRepository = (*previewClassScheduleRepoStub)(nil)

type previewClassScheduleRepoStub struct {
	created []entities.ClassSchedule
}

func (s *previewClassScheduleRepoStub) GetTable() string { return "class_schedules" }
func (s *previewClassScheduleRepoStub) GetByCondition(_ context.Context, _ *repositories.CommonCondition) (*repositories.Pagination[entities.ClassSchedule], error) {
	return &repositories.Pagination[entities.ClassSchedule]{}, nil
}
func (s *previewClassScheduleRepoStub) GetTotal(_ context.Context, _ *repositories.CommonCondition) (uint64, error) {
	return uint64(len(s.created)), nil
}
func (s *previewClassScheduleRepoStub) Create(_ context.Context, entity *entities.ClassSchedule) (*entities.ClassSchedule, error) {
	s.created = append(s.created, *entity)
	return entity, nil
}
func (s *previewClassScheduleRepoStub) Update(_ context.Context, _ interface{}, _ map[string]interface{}) error {
	return nil
}
func (s *previewClassScheduleRepoStub) UpdateWithIDs(_ context.Context, _ []string, _ map[string]interface{}) error {
	return nil
}
func (s *previewClassScheduleRepoStub) SoftDelete(_ context.Context, _ interface{}) error { return nil }
func (s *previewClassScheduleRepoStub) HardDelete(_ context.Context, _ interface{}) error { return nil }
func (s *previewClassScheduleRepoStub) GetByID(_ context.Context, _ interface{}) (*entities.ClassSchedule, error) {
	return nil, nil
}
func (s *previewClassScheduleRepoStub) GetSchedulesByClassID(_ context.Context, _ string) ([]entities.ClassSchedule, error) {
	return append([]entities.ClassSchedule(nil), s.created...), nil
}

var _ repositoryinterface.EnrollmentRepository = (*previewEnrollmentRepoStub)(nil)

type previewEnrollmentRepoStub struct {
	byClass map[string][]entities.Enrollment
}

func (s *previewEnrollmentRepoStub) GetTable() string { return "enrollments" }
func (s *previewEnrollmentRepoStub) GetByCondition(_ context.Context, _ *repositories.CommonCondition) (*repositories.Pagination[entities.Enrollment], error) {
	data := make([]*entities.Enrollment, 0)
	if s.byClass != nil {
		for _, classEnrollments := range s.byClass {
			for i := range classEnrollments {
				data = append(data, &classEnrollments[i])
			}
		}
	}
	return &repositories.Pagination[entities.Enrollment]{
		Data: data,
		Meta: repositories.Meta{
			TotalItems: uint64(len(data)),
		},
	}, nil
}
func (s *previewEnrollmentRepoStub) GetTotal(_ context.Context, _ *repositories.CommonCondition) (uint64, error) {
	return 0, nil
}
func (s *previewEnrollmentRepoStub) Create(_ context.Context, entity *entities.Enrollment) (*entities.Enrollment, error) {
	return entity, nil
}
func (s *previewEnrollmentRepoStub) Update(_ context.Context, _ interface{}, _ map[string]interface{}) error {
	return nil
}
func (s *previewEnrollmentRepoStub) UpdateWithIDs(_ context.Context, _ []string, _ map[string]interface{}) error {
	return nil
}
func (s *previewEnrollmentRepoStub) SoftDelete(_ context.Context, _ interface{}) error { return nil }
func (s *previewEnrollmentRepoStub) HardDelete(_ context.Context, _ interface{}) error { return nil }
func (s *previewEnrollmentRepoStub) GetByID(_ context.Context, _ interface{}) (*entities.Enrollment, error) {
	return nil, nil
}
func (s *previewEnrollmentRepoStub) ListByClassID(_ context.Context, classID string) ([]entities.Enrollment, error) {
	if s.byClass == nil {
		return buildApprovedEnrollments(classID, 18), nil
	}
	return append([]entities.Enrollment(nil), s.byClass[classID]...), nil
}

type previewUnitOfWorkStub struct {
	beginCalled    bool
	commitCalled   bool
	rollbackCalled bool
}

func (s *previewUnitOfWorkStub) Begin(ctx context.Context) (context.Context, error) {
	s.beginCalled = true
	return ctx, nil
}

func (s *previewUnitOfWorkStub) Commit(_ context.Context) error {
	s.commitCalled = true
	return nil
}

func (s *previewUnitOfWorkStub) Rollback(_ context.Context) error {
	s.rollbackCalled = true
	return nil
}

type previewTravelRepoStub struct {
	repositories.BaseRepository[entities.CampusTravelTime]
}

func (s *previewTravelRepoStub) GetByCondition(ctx context.Context, condition *repositories.CommonCondition) (*repositories.Pagination[entities.CampusTravelTime], error) {
	return &repositories.Pagination[entities.CampusTravelTime]{Data: []*entities.CampusTravelTime{}}, nil
}

func TestPreviewUseCase_NonvolatileReplanning_MinimizesDisruption(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()

	// A Published lesson exactly matching the 1st session of the fixture
	lessonRepo := &previewLessonRepoStub{
		existingLessons: []entities.Lesson{
			{
				ID:        "lesson-published-target",
				ClassID:   "class-1",
				TeacherID: stringPtr("teacher-1"),
				RoomID:    stringPtr("room-1"),
				DateStart: time.Date(2026, 6, 15, 8, 0, 0, 0, time.UTC),
				DateEnd:   time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC),
				Status:    entities.LessonStatusPublished,
			},
		},
	}
	enrollmentRepo := &previewEnrollmentRepoStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	// In replan_draft mode, it should load the published lesson as target, and NO schedule change should occur
	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC),
		Mode:     schedulingservice.PreviewModeReplanDraft,
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	if result.Summary.ScheduleChangeCount != 0 {
		t.Fatalf("expected 0 schedule changes for stable replanning, got %d", result.Summary.ScheduleChangeCount)
	}
	if result.Summary.RoomChangeCount != 0 {
		t.Fatalf("expected 0 room changes, got %d", result.Summary.RoomChangeCount)
	}
	if result.Summary.TeacherChangeCount != 0 {
		t.Fatalf("expected 0 teacher changes, got %d", result.Summary.TeacherChangeCount)
	}
}

func TestPreviewUseCase_ReplanSelectedClassIgnoresOwnPublishedLessonsAsConflicts(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{
		existingLessons: []entities.Lesson{
			{
				ID:        "lesson-published-target",
				ClassID:   "class-1",
				TeacherID: stringPtr("teacher-1"),
				RoomID:    stringPtr("room-1"),
				DateStart: time.Date(2026, 4, 13, 8, 0, 0, 0, time.UTC),
				DateEnd:   time.Date(2026, 4, 13, 10, 0, 0, 0, time.UTC),
				Status:    entities.LessonStatusPublished,
			},
		},
	}
	enrollmentRepo := &previewEnrollmentRepoStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
		ClassIDs: []string{"class-1"},
		Mode:     schedulingservice.PreviewModeReplanWithPublishedLock,
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	for _, conflict := range result.Conflicts {
		if conflict.Type == "SYSTEM_LESSON_CONFLICT" {
			t.Fatalf("did not expect the selected class to conflict with its own published lessons: %s", conflict.Message)
		}
	}
	if len(result.ExistingLessons) != 0 {
		t.Fatalf("expected selected class published lessons to stay out of existing lock list, got %d", len(result.ExistingLessons))
	}
}

func TestCommitPreviewUseCase_ReplacesOwnPublishedLessonByArchivingOldOne(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{
		existingLessons: []entities.Lesson{
			{
				ID:        "lesson-published-target",
				ClassID:   "class-1",
				TeacherID: stringPtr("teacher-1"),
				RoomID:    stringPtr("room-1"),
				DateStart: time.Date(2026, 6, 15, 8, 0, 0, 0, time.UTC),
				DateEnd:   time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC),
				Status:    entities.LessonStatusPublished,
			},
		},
	}
	classScheduleRepo := &previewClassScheduleRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{}
	uow := &previewUnitOfWorkStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)
	commitUseCase := NewCommitPreviewUseCase(
		lessonRepo,
		classScheduleRepo,
		uow,
		logger.NewZapLogger(logger.Config{Level: "error", Format: "json", Output: "stdout", ServiceName: "test", Environment: "test"}),
		store,
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC),
		ClassIDs: []string{"class-1"},
		Mode:     schedulingservice.PreviewModeReplanWithPublishedLock,
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}
	if len(result.Assignments) == 0 {
		t.Fatalf("expected assignments for replan")
	}
	if result.Assignments[0].ReplaceLessonID != "lesson-published-target" {
		t.Fatalf("expected replace_lesson_id to map target lesson, got %q", result.Assignments[0].ReplaceLessonID)
	}

	output, err := commitUseCase.Execute(context.Background(), CommitPreviewInput{RunID: result.RunID})
	if err != nil {
		t.Fatalf("unexpected commit error: %v", err)
	}
	if output.Status != "COMMITTED" {
		t.Fatalf("expected COMMITTED status, got %s", output.Status)
	}
	if len(lessonRepo.createdLessons) == 0 {
		t.Fatalf("expected replacement lesson to be created")
	}
	if lessonRepo.createdLessons[0].ChangeReason != entities.LessonChangeReasonReplanReplacement {
		t.Fatalf("expected replacement lesson change_reason %s, got %s", entities.LessonChangeReasonReplanReplacement, lessonRepo.createdLessons[0].ChangeReason)
	}
	if lessonRepo.existingLessons[0].Status != entities.LessonStatusHistory {
		t.Fatalf("expected old lesson to be archived as HISTORY, got %s", lessonRepo.existingLessons[0].Status)
	}
	if lessonRepo.existingLessons[0].ChangeReason != entities.LessonChangeReasonReplanArchived {
		t.Fatalf("expected old lesson change_reason %s, got %s", entities.LessonChangeReasonReplanArchived, lessonRepo.existingLessons[0].ChangeReason)
	}
}

func TestPreviewUseCase_CapacityUtilization(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()

	lessonRepo := &previewLessonRepoStub{}
	enrollments := make([]entities.Enrollment, 0, 16)
	now := time.Now()
	for i := 0; i < 16; i++ {
		enrollments = append(enrollments, entities.Enrollment{
			StudentID:  fmt.Sprintf("student-%d", i),
			ClassID:    "class-1",
			Status:     "ENROLLED",
			ApprovedAt: &now,
		})
	}

	enrollmentRepo := &previewEnrollmentRepoStub{
		byClass: map[string][]entities.Enrollment{
			"class-1": enrollments,
		},
	}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
		&previewTravelRepoStub{},
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
		Mode:     schedulingservice.PreviewModeReplanDraft,
		ClassIDs: []string{"class-1"}, // only preview class 1
	})
	if err != nil {
		t.Fatalf("unexpected preview error: %v", err)
	}

	if result.Status == "FAILED" {
		t.Fatalf("preview failed, expected PARTIAL or COMPLETED")
	}

	for _, assignment := range result.Assignments {
		if assignment.ExpectedStudentCount != 16 {
			t.Errorf("expected student count 16 for class-1, got %d", assignment.ExpectedStudentCount)
		}
	}

	if result.Summary.AverageCapacityUtilization <= 0.0 {
		t.Errorf("expected average capacity utilization to be > 0, got %f", result.Summary.AverageCapacityUtilization)
	}
}
