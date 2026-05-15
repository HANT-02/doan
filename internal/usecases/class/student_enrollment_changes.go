package class

import (
	"context"
	"fmt"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	schedulingservice "doan/internal/services/scheduling"
	scheduling "doan/internal/usecases/scheduling"
	"doan/pkg/logger"
)

const defaultTransferValidationWindowDays = 90

type ReserveStudentInput struct {
	ClassID     string
	StudentID   string
	Reason      string
	EffectiveAt *time.Time
}

type ReserveStudentOutput struct {
	Message             string    `json:"message"`
	ClassID             string    `json:"class_id"`
	StudentID           string    `json:"student_id"`
	Status              string    `json:"status"`
	ImpactedLessonCount int       `json:"impacted_lesson_count"`
	EffectiveAt         time.Time `json:"effective_at"`
}

type ReserveStudentUseCase interface {
	Execute(ctx context.Context, input ReserveStudentInput) (*ReserveStudentOutput, error)
}

type ResumeStudentInput struct {
	ClassID     string
	StudentID   string
	Reason      string
	EffectiveAt *time.Time
}

type ResumeStudentOutput struct {
	Message             string    `json:"message"`
	ClassID             string    `json:"class_id"`
	StudentID           string    `json:"student_id"`
	Status              string    `json:"status"`
	ImpactedLessonCount int       `json:"impacted_lesson_count"`
	EffectiveAt         time.Time `json:"effective_at"`
}

type ResumeStudentUseCase interface {
	Execute(ctx context.Context, input ResumeStudentInput) (*ResumeStudentOutput, error)
}

type TransferStudentInput struct {
	SourceClassID string
	TargetClassID string
	StudentID     string
	Reason        string
	EffectiveAt   *time.Time
}

type TransferStudentOutput struct {
	Message               string    `json:"message"`
	StudentID             string    `json:"student_id"`
	SourceClassID         string    `json:"source_class_id"`
	TargetClassID         string    `json:"target_class_id"`
	SourceEnrollmentState string    `json:"source_enrollment_state"`
	TargetEnrollmentState string    `json:"target_enrollment_state"`
	ImpactedLessonCount   int       `json:"impacted_lesson_count"`
	RemainingCapacity     int       `json:"remaining_capacity"`
	CapacityUtilization   float64   `json:"capacity_utilization"`
	EffectiveAt           time.Time `json:"effective_at"`
}

type TransferStudentUseCase interface {
	Execute(ctx context.Context, input TransferStudentInput) (*TransferStudentOutput, error)
}

type reserveStudentUseCase struct {
	enrollmentRepo repositoryinterface.EnrollmentRepository
	lessonRepo     repositoryinterface.LessonRepository
	uow            repositories.UnitOfWork
	log            logger.Logger
}

type resumeStudentUseCase struct {
	enrollmentRepo repositoryinterface.EnrollmentRepository
	lessonRepo     repositoryinterface.LessonRepository
	uow            repositories.UnitOfWork
	log            logger.Logger
}

type transferStudentUseCase struct {
	classRepo        repositoryinterface.ClassRepository
	enrollmentRepo   repositoryinterface.EnrollmentRepository
	lessonRepo       repositoryinterface.LessonRepository
	campusTravelRepo repositoryinterface.CampusTravelTimeRepository
	uow              repositories.UnitOfWork
	log              logger.Logger
}

func NewReserveStudentUseCase(
	enrollmentRepo repositoryinterface.EnrollmentRepository,
	lessonRepo repositoryinterface.LessonRepository,
	uow repositories.UnitOfWork,
	log logger.Logger,
) ReserveStudentUseCase {
	return &reserveStudentUseCase{
		enrollmentRepo: enrollmentRepo,
		lessonRepo:     lessonRepo,
		uow:            uow,
		log:            log,
	}
}

func NewResumeStudentUseCase(
	enrollmentRepo repositoryinterface.EnrollmentRepository,
	lessonRepo repositoryinterface.LessonRepository,
	uow repositories.UnitOfWork,
	log logger.Logger,
) ResumeStudentUseCase {
	return &resumeStudentUseCase{
		enrollmentRepo: enrollmentRepo,
		lessonRepo:     lessonRepo,
		uow:            uow,
		log:            log,
	}
}

func NewTransferStudentUseCase(
	classRepo repositoryinterface.ClassRepository,
	enrollmentRepo repositoryinterface.EnrollmentRepository,
	lessonRepo repositoryinterface.LessonRepository,
	campusTravelRepo repositoryinterface.CampusTravelTimeRepository,
	uow repositories.UnitOfWork,
	log logger.Logger,
) TransferStudentUseCase {
	return &transferStudentUseCase{
		classRepo:        classRepo,
		enrollmentRepo:   enrollmentRepo,
		lessonRepo:       lessonRepo,
		campusTravelRepo: campusTravelRepo,
		uow:              uow,
		log:              log,
	}
}

func (uc *reserveStudentUseCase) Execute(ctx context.Context, input ReserveStudentInput) (*ReserveStudentOutput, error) {
	if input.ClassID == "" || input.StudentID == "" {
		return nil, ErrEnrollmentNotFound
	}

	effectiveAt := resolveEffectiveAt(input.EffectiveAt)
	enrollment, err := findActiveEnrollment(ctx, uc.enrollmentRepo, input.ClassID, input.StudentID)
	if err != nil {
		return nil, err
	}

	impactedLessons, err := countFutureLessons(ctx, uc.lessonRepo, input.ClassID, effectiveAt)
	if err != nil {
		return nil, err
	}

	_, err = repositories.ExecuteInTransaction(ctx, uc.uow, uc.log, func(txCtx context.Context) (interface{}, error) {
		return nil, uc.enrollmentRepo.Update(txCtx, enrollment.ID, map[string]interface{}{
			"status":     entities.EnrollmentStatusSuspended,
			"updated_at": time.Now(),
		})
	})
	if err != nil {
		return nil, err
	}

	return &ReserveStudentOutput{
		Message:             "Đã bảo lưu học viên khỏi lớp hiện tại",
		ClassID:             input.ClassID,
		StudentID:           input.StudentID,
		Status:              entities.EnrollmentStatusSuspended,
		ImpactedLessonCount: impactedLessons,
		EffectiveAt:         effectiveAt,
	}, nil
}

func (uc *resumeStudentUseCase) Execute(ctx context.Context, input ResumeStudentInput) (*ResumeStudentOutput, error) {
	if input.ClassID == "" || input.StudentID == "" {
		return nil, ErrReservedEnrollmentNotFound
	}

	effectiveAt := resolveEffectiveAt(input.EffectiveAt)
	enrollment, err := findEnrollmentByStatuses(ctx, uc.enrollmentRepo, input.ClassID, input.StudentID, []string{
		entities.EnrollmentStatusSuspended,
	})
	if err != nil {
		return nil, err
	}

	impactedLessons, err := countFutureLessons(ctx, uc.lessonRepo, input.ClassID, effectiveAt)
	if err != nil {
		return nil, err
	}

	_, err = repositories.ExecuteInTransaction(ctx, uc.uow, uc.log, func(txCtx context.Context) (interface{}, error) {
		return nil, uc.enrollmentRepo.Update(txCtx, enrollment.ID, map[string]interface{}{
			"status":     entities.EnrollmentStatusEnrolled,
			"updated_at": time.Now(),
		})
	})
	if err != nil {
		return nil, err
	}

	return &ResumeStudentOutput{
		Message:             "Đã hoàn tác bảo lưu và cho học viên quay lại lớp",
		ClassID:             input.ClassID,
		StudentID:           input.StudentID,
		Status:              entities.EnrollmentStatusEnrolled,
		ImpactedLessonCount: impactedLessons,
		EffectiveAt:         effectiveAt,
	}, nil
}

func (uc *transferStudentUseCase) Execute(ctx context.Context, input TransferStudentInput) (*TransferStudentOutput, error) {
	if input.SourceClassID == "" || input.StudentID == "" {
		return nil, ErrEnrollmentNotFound
	}
	if input.TargetClassID == "" {
		return nil, ErrTargetClassRequired
	}
	if input.SourceClassID == input.TargetClassID {
		return nil, ErrTransferSameClass
	}

	effectiveAt := resolveEffectiveAt(input.EffectiveAt)
	sourceEnrollment, err := findActiveEnrollment(ctx, uc.enrollmentRepo, input.SourceClassID, input.StudentID)
	if err != nil {
		return nil, err
	}

	sourceClass, err := loadClassWithRelations(ctx, uc.classRepo, input.SourceClassID)
	if err != nil {
		return nil, err
	}
	targetClass, err := loadClassWithRelations(ctx, uc.classRepo, input.TargetClassID)
	if err != nil {
		return nil, err
	}
	if targetClass.Status != "OPEN" {
		return nil, ErrTargetClassNotOpen
	}
	if !areClassesCompatible(sourceClass, targetClass) {
		return nil, ErrTargetClassIncompatible
	}
	if _, err := findActiveEnrollment(ctx, uc.enrollmentRepo, input.TargetClassID, input.StudentID); err == nil {
		return nil, ErrStudentAlreadyInTarget
	}

	targetEnrollments, err := uc.enrollmentRepo.ListByClassID(ctx, input.TargetClassID)
	if err != nil {
		return nil, err
	}
	targetStudentCount := countActiveEnrollments(targetEnrollments)
	capacityLimit := scheduling.CalculateCapacityLimit(targetClass.Room.Capacity, targetClass.MaxStudents)
	if !scheduling.ValidateMakeupCapacity(targetStudentCount, capacityLimit, 1) {
		return nil, ErrTargetClassFull
	}

	if err := uc.validateTransferSchedule(ctx, input.StudentID, input.SourceClassID, targetClass, effectiveAt); err != nil {
		return nil, err
	}

	impactedLessons, err := countFutureLessons(ctx, uc.lessonRepo, input.SourceClassID, effectiveAt)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	newEnrollment := &entities.Enrollment{
		ClassID:    input.TargetClassID,
		StudentID:  input.StudentID,
		Status:     entities.EnrollmentStatusEnrolled,
		ApprovedAt: &now,
	}

	_, err = repositories.ExecuteInTransaction(ctx, uc.uow, uc.log, func(txCtx context.Context) (interface{}, error) {
		if err := uc.enrollmentRepo.Update(txCtx, sourceEnrollment.ID, map[string]interface{}{
			"status":     entities.EnrollmentStatusTransferred,
			"updated_at": now,
		}); err != nil {
			return nil, err
		}
		if _, err := uc.enrollmentRepo.Create(txCtx, newEnrollment); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	remainingCapacity, utilization := scheduling.CalculateUtilization(targetStudentCount+1, capacityLimit)
	return &TransferStudentOutput{
		Message:               "Đã chuyển học viên sang lớp đích",
		StudentID:             input.StudentID,
		SourceClassID:         input.SourceClassID,
		TargetClassID:         input.TargetClassID,
		SourceEnrollmentState: entities.EnrollmentStatusTransferred,
		TargetEnrollmentState: entities.EnrollmentStatusEnrolled,
		ImpactedLessonCount:   impactedLessons,
		RemainingCapacity:     remainingCapacity,
		CapacityUtilization:   utilization,
		EffectiveAt:           effectiveAt,
	}, nil
}

func (uc *transferStudentUseCase) validateTransferSchedule(
	ctx context.Context,
	studentID string,
	sourceClassID string,
	targetClass *entities.Class,
	effectiveAt time.Time,
) error {
	activeClassIDs, err := loadActiveStudentClassIDs(ctx, uc.enrollmentRepo, studentID, sourceClassID, targetClass.ID)
	if err != nil {
		return err
	}

	windowEnd := resolveValidationWindowEnd(targetClass.EndDate, effectiveAt)
	targetLessons, err := loadFutureLessons(ctx, uc.lessonRepo, []string{targetClass.ID}, effectiveAt, windowEnd)
	if err != nil {
		return err
	}
	if len(targetLessons) == 0 {
		return nil
	}

	otherLessons, err := loadFutureLessons(ctx, uc.lessonRepo, activeClassIDs, effectiveAt, windowEnd)
	if err != nil {
		return err
	}

	travelMap := uc.loadTravelMap(ctx)

	for _, targetLesson := range targetLessons {
		for _, otherLesson := range otherLessons {
			if lessonsOverlap(targetLesson, otherLesson) {
				return ErrStudentScheduleConflict
			}
			if sameCalendarDay(targetLesson.DateStart, otherLesson.DateStart) && !hasTravelGapBetweenLessons(targetLesson, otherLesson, travelMap) {
				return ErrStudentTravelConflict
			}
		}
	}

	return nil
}

func resolveEffectiveAt(value *time.Time) time.Time {
	if value == nil || value.IsZero() {
		return time.Now()
	}
	return value.UTC()
}

func resolveValidationWindowEnd(classEndDate *time.Time, effectiveAt time.Time) time.Time {
	windowEnd := effectiveAt.AddDate(0, 0, defaultTransferValidationWindowDays)
	if classEndDate != nil && classEndDate.After(effectiveAt) && classEndDate.Before(windowEnd) {
		return classEndDate.UTC().Add(24 * time.Hour)
	}
	return windowEnd
}

func loadClassWithRelations(
	ctx context.Context,
	classRepo repositoryinterface.ClassRepository,
	classID string,
) (*entities.Class, error) {
	condition := repositories.NewCommonCondition()
	condition.AddCondition("id", classID, repositories.Equal)
	condition.SetPreload([]string{"Course", "Room", "Teacher"})
	condition.SetPaging(1, 1)

	result, err := classRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 || result.Data[0] == nil {
		return nil, fmt.Errorf("class %s not found", classID)
	}
	return result.Data[0], nil
}

func findActiveEnrollment(
	ctx context.Context,
	enrollmentRepo repositoryinterface.EnrollmentRepository,
	classID string,
	studentID string,
) (*entities.Enrollment, error) {
	return findEnrollmentByStatuses(ctx, enrollmentRepo, classID, studentID, []string{
		entities.EnrollmentStatusEnrolled,
		entities.EnrollmentStatusApproved,
	})
}

func findEnrollmentByStatuses(
	ctx context.Context,
	enrollmentRepo repositoryinterface.EnrollmentRepository,
	classID string,
	studentID string,
	statuses []string,
) (*entities.Enrollment, error) {
	condition := repositories.NewCommonCondition()
	condition.AddCondition("class_id", classID, repositories.Equal)
	condition.AddCondition("student_id", studentID, repositories.Equal)
	condition.AddCondition("status", statuses, repositories.In)
	condition.SetPaging(1, 1)

	result, err := enrollmentRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 || result.Data[0] == nil {
		if len(statuses) == 1 && statuses[0] == entities.EnrollmentStatusSuspended {
			return nil, ErrReservedEnrollmentNotFound
		}
		return nil, ErrEnrollmentNotFound
	}
	return result.Data[0], nil
}

func loadActiveStudentClassIDs(
	ctx context.Context,
	enrollmentRepo repositoryinterface.EnrollmentRepository,
	studentID string,
	excludedClassIDs ...string,
) ([]string, error) {
	condition := repositories.NewCommonCondition()
	condition.AddCondition("student_id", studentID, repositories.Equal)
	condition.AddCondition("status", []string{
		entities.EnrollmentStatusEnrolled,
		entities.EnrollmentStatusApproved,
	}, repositories.In)
	condition.SetPaging(1, 2000)

	result, err := enrollmentRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}

	excluded := make(map[string]struct{}, len(excludedClassIDs))
	for _, classID := range excludedClassIDs {
		if classID != "" {
			excluded[classID] = struct{}{}
		}
	}

	unique := make(map[string]struct{})
	classIDs := make([]string, 0)
	if result == nil {
		return classIDs, nil
	}
	for _, enrollment := range result.Data {
		if enrollment == nil {
			continue
		}
		if _, skip := excluded[enrollment.ClassID]; skip {
			continue
		}
		if _, exists := unique[enrollment.ClassID]; exists {
			continue
		}
		unique[enrollment.ClassID] = struct{}{}
		classIDs = append(classIDs, enrollment.ClassID)
	}
	return classIDs, nil
}

func loadFutureLessons(
	ctx context.Context,
	lessonRepo repositoryinterface.LessonRepository,
	classIDs []string,
	from time.Time,
	to time.Time,
) ([]entities.Lesson, error) {
	if len(classIDs) == 0 {
		return []entities.Lesson{}, nil
	}
	return lessonRepo.FindOverlappingLessons(ctx, from, to, classIDs, nil, nil, []string{
		entities.LessonStatusDraft,
		entities.LessonStatusPublished,
	})
}

func countFutureLessons(
	ctx context.Context,
	lessonRepo repositoryinterface.LessonRepository,
	classID string,
	from time.Time,
) (int, error) {
	windowEnd := from.AddDate(0, 0, defaultTransferValidationWindowDays)
	lessons, err := loadFutureLessons(ctx, lessonRepo, []string{classID}, from, windowEnd)
	if err != nil {
		return 0, err
	}
	return len(lessons), nil
}

func areClassesCompatible(sourceClass, targetClass *entities.Class) bool {
	if sourceClass == nil || targetClass == nil {
		return false
	}
	if sourceClass.CourseID != nil && targetClass.CourseID != nil && *sourceClass.CourseID == *targetClass.CourseID {
		return true
	}
	if sourceClass.Course.Subject == "" || targetClass.Course.Subject == "" {
		return false
	}
	if sourceClass.Course.GradeLevel == "" || targetClass.Course.GradeLevel == "" {
		return false
	}
	return sourceClass.Course.Subject == targetClass.Course.Subject &&
		sourceClass.Course.GradeLevel == targetClass.Course.GradeLevel
}

func countActiveEnrollments(items []entities.Enrollment) int {
	total := 0
	for _, item := range items {
		if isActiveEnrollmentStatus(item.Status) {
			total++
		}
	}
	return total
}

func isActiveEnrollmentStatus(status string) bool {
	return status == entities.EnrollmentStatusEnrolled || status == entities.EnrollmentStatusApproved
}

func lessonsOverlap(left, right entities.Lesson) bool {
	return left.DateStart.Before(right.DateEnd) && left.DateEnd.After(right.DateStart)
}

func sameCalendarDay(left, right time.Time) bool {
	ly, lm, ld := left.Date()
	ry, rm, rd := right.Date()
	return ly == ry && lm == rm && ld == rd
}

func hasTravelGapBetweenLessons(left, right entities.Lesson, travelMap map[string]int) bool {
	if left.DateStart.Before(right.DateStart) {
		return schedulingservice.HasSufficientTravelGap(left.DateEnd, right.DateStart, &left.Room, &right.Room, travelMap)
	}
	return schedulingservice.HasSufficientTravelGap(right.DateEnd, left.DateStart, &right.Room, &left.Room, travelMap)
}

func (uc *transferStudentUseCase) loadTravelMap(ctx context.Context) map[string]int {
	condition := repositories.NewCommonCondition()
	condition.AddCondition("is_active", true, repositories.Equal)
	condition.SetPaging(1, 2000)

	result, err := uc.campusTravelRepo.GetByCondition(ctx, condition)
	if err != nil || result == nil {
		return map[string]int{}
	}

	items := make([]entities.CampusTravelTime, 0, len(result.Data))
	for _, item := range result.Data {
		if item == nil {
			continue
		}
		items = append(items, *item)
	}
	return schedulingservice.BuildCampusTravelTimeMap(items)
}
