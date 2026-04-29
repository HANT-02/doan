package scheduling

import (
	"context"
	"strings"
	"testing"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	schedulingservice "doan/internal/services/scheduling"
	"doan/pkg/logger"
)

func TestPreviewAndCommitUseCase_CommitAssignmentsWithDefaultSolver(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{}
	uow := &previewUnitOfWorkStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)
	commitUseCase := NewCommitPreviewUseCase(
		lessonRepo,
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
	if !uow.beginCalled || !uow.commitCalled || uow.rollbackCalled {
		t.Fatalf("unexpected transaction lifecycle: begin=%t commit=%t rollback=%t", uow.beginCalled, uow.commitCalled, uow.rollbackCalled)
	}
}

func TestCommitPreviewUseCase_ReturnsConflictWhenExistingLessonOverlaps(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositories()
	lessonRepo := &previewLessonRepoStub{}
	uow := &previewUnitOfWorkStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)
	commitUseCase := NewCommitPreviewUseCase(
		lessonRepo,
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
	if !uow.beginCalled || uow.commitCalled || !uow.rollbackCalled {
		t.Fatalf("unexpected transaction lifecycle: begin=%t commit=%t rollback=%t", uow.beginCalled, uow.commitCalled, uow.rollbackCalled)
	}
}

func TestPreviewUseCase_ReturnsSpecificConflictWhenDateRangeHasTooFewScheduleSlots(t *testing.T) {
	t.Parallel()

	store := schedulingservice.NewPreviewStore[PreviewResult]()
	classRepo, roomRepo, shiftRepo := previewFixtureRepositoriesWithSessionCount(8)

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
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

	if result.Status != "PARTIAL" {
		t.Fatalf("expected PARTIAL preview status, got %s", result.Status)
	}
	if result.Summary.ScheduledLessons != 2 {
		t.Fatalf("expected exactly 2 scheduled lessons in the selected range, got %d", result.Summary.ScheduledLessons)
	}
	if result.Summary.UnscheduledLessons != 6 {
		t.Fatalf("expected 6 unscheduled lessons, got %d", result.Summary.UnscheduledLessons)
	}

	for _, conflict := range result.Conflicts {
		if conflict.Type != "INSUFFICIENT_SCHEDULE_SLOTS" {
			t.Fatalf("expected INSUFFICIENT_SCHEDULE_SLOTS conflict, got %s", conflict.Type)
		}
	}
}

func previewFixtureRepositories() (repositoryinterface.ClassRepository, repositoryinterface.RoomRepository, repositoryinterface.ShiftRepository) {
	return previewFixtureRepositoriesWithSessionCount(2)
}

func previewFixtureRepositoriesWithSessionCount(sessionCount int) (repositoryinterface.ClassRepository, repositoryinterface.RoomRepository, repositoryinterface.ShiftRepository) {
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
				},
				CourseID: &courseID,
				Course: entities.Course{
					ID:                     courseID,
					Code:                   "KH-001",
					Name:                   "Khoa hoc Toan",
					SessionCount:           sessionCount,
					SessionDurationMinutes: 120,
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

func (s *previewLessonRepoStub) Update(_ context.Context, _ interface{}, _ map[string]interface{}) error {
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

func (s *previewLessonRepoStub) FindOverlappingLessons(
	_ context.Context,
	_ time.Time,
	_ time.Time,
	_ []string,
	_ []string,
	_ []string,
) ([]entities.Lesson, error) {
	return append([]entities.Lesson(nil), s.existingLessons...), nil
}

func (s *previewLessonRepoStub) GetLessonWithRelations(_ context.Context, _ string) (*entities.Lesson, error) {
	return nil, nil
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
