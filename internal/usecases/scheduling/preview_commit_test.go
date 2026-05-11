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
	classScheduleRepo := &previewClassScheduleRepoStub{}
	enrollmentRepo := &previewEnrollmentRepoStub{}
	uow := &previewUnitOfWorkStub{}

	previewUseCase := NewPreviewUseCase(
		classRepo,
		roomRepo,
		shiftRepo,
		lessonRepo,
		enrollmentRepo,
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
	if !uow.beginCalled || !uow.commitCalled || uow.rollbackCalled {
		t.Fatalf("unexpected transaction lifecycle: begin=%t commit=%t rollback=%t", uow.beginCalled, uow.commitCalled, uow.rollbackCalled)
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
		store,
		schedulingservice.NewDefaultSchedulingSolver(schedulingservice.NewCPSATSolver()),
	)

	firstPreview, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
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
			ClassID:   firstAssignment.ClassID,
			TeacherID: &firstAssignment.TeacherID,
			RoomID:    &firstAssignment.RoomID,
			DateStart: firstAssignment.StartTime,
			DateEnd:   firstAssignment.EndTime,
		},
	}

	result, err := previewUseCase.Execute(context.Background(), PreviewInput{
		DateFrom: time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC),
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

func TestPreviewUseCase_ReturnsSpecificConflictWhenDateRangeHasTooFewScheduleSlots(t *testing.T) {
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
) ([]entities.Lesson, error) {
	return append([]entities.Lesson(nil), s.existingLessons...), nil
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
	return &repositories.Pagination[entities.Enrollment]{}, nil
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
