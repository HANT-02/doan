package scheduling

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	schedulingservice "doan/internal/services/scheduling"
)

const defaultMakeupWindowDays = 60

type MakeupSpot struct {
	LessonID             string    `json:"lesson_id"`
	ClassID              string    `json:"class_id"`
	ClassCode            string    `json:"class_code"`
	ClassName            string    `json:"class_name"`
	CourseID             string    `json:"course_id,omitempty"`
	CourseName           string    `json:"course_name,omitempty"`
	TeacherID            string    `json:"teacher_id,omitempty"`
	TeacherName          string    `json:"teacher_name,omitempty"`
	RoomID               string    `json:"room_id,omitempty"`
	RoomName             string    `json:"room_name,omitempty"`
	StartTime            time.Time `json:"start_time"`
	EndTime              time.Time `json:"end_time"`
	MatchType            string    `json:"match_type"`
	ExpectedStudentCount int       `json:"expected_student_count"`
	CapacityLimit        int       `json:"capacity_limit"`
	RemainingCapacity    int       `json:"remaining_capacity"`
	CapacityUtilization  float64   `json:"capacity_utilization"`
	Eligible             bool      `json:"eligible"`
	Reasons              []string  `json:"reasons"`
}

type FindMakeupSpotsInput struct {
	LessonID  string
	StudentID string
	Limit     int
}

type FindMakeupSpotsOutput struct {
	StudentID string       `json:"student_id"`
	LessonID  string       `json:"lesson_id"`
	Spots     []MakeupSpot `json:"spots"`
}

type MakeupUseCase interface {
	FindMakeupSpots(ctx context.Context, input FindMakeupSpotsInput) (*FindMakeupSpotsOutput, error)
}

type makeupUseCase struct {
	lessonRepo       repositoryinterface.LessonRepository
	classRepo        repositoryinterface.ClassRepository
	enrollmentRepo   repositoryinterface.EnrollmentRepository
	studentRepo      repositoryinterface.StudentRepository
	campusTravelRepo repositoryinterface.CampusTravelTimeRepository
}

func NewMakeupUseCase(
	lessonRepo repositoryinterface.LessonRepository,
	classRepo repositoryinterface.ClassRepository,
	enrollmentRepo repositoryinterface.EnrollmentRepository,
	studentRepo repositoryinterface.StudentRepository,
	campusTravelRepo repositoryinterface.CampusTravelTimeRepository,
) MakeupUseCase {
	return &makeupUseCase{
		lessonRepo:       lessonRepo,
		classRepo:        classRepo,
		enrollmentRepo:   enrollmentRepo,
		studentRepo:      studentRepo,
		campusTravelRepo: campusTravelRepo,
	}
}

func (uc *makeupUseCase) FindMakeupSpots(ctx context.Context, input FindMakeupSpotsInput) (*FindMakeupSpotsOutput, error) {
	if input.LessonID == "" || input.StudentID == "" {
		return nil, errors.New("lesson_id and student_id are required")
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}

	sourceLesson, err := uc.lessonRepo.GetLessonWithRelations(ctx, input.LessonID)
	if err != nil || sourceLesson == nil {
		return nil, err
	}

	student, err := uc.studentRepo.GetByID(ctx, input.StudentID)
	if err != nil || student == nil {
		return nil, err
	}

	sourceClass, err := uc.loadClassWithCourse(ctx, sourceLesson.ClassID)
	if err != nil || sourceClass == nil {
		return nil, err
	}

	enrolledClassIDs, err := uc.loadStudentClassIDs(ctx, input.StudentID)
	if err != nil {
		return nil, err
	}

	studentLessons, err := uc.loadStudentLessons(ctx, enrolledClassIDs, sourceLesson.DateStart.AddDate(0, 0, -1), sourceLesson.DateStart.AddDate(0, 0, defaultMakeupWindowDays))
	if err != nil {
		return nil, err
	}

	candidateLessons, err := uc.loadCandidateLessons(ctx, sourceClass, sourceLesson.DateStart)
	if err != nil {
		return nil, err
	}

	travelMap := uc.loadTravelMap(ctx)
	spots := make([]MakeupSpot, 0)

	for _, lesson := range candidateLessons {
		if lesson == nil || lesson.ID == sourceLesson.ID || lesson.ClassID == sourceLesson.ClassID {
			continue
		}

		matchType := determineMakeupMatchType(*sourceClass, lesson.Class)
		if matchType == "" {
			continue
		}

		classEnrollments, err := uc.enrollmentRepo.ListByClassID(ctx, lesson.ClassID)
		if err != nil {
			return nil, err
		}
		if classHasStudent(classEnrollments, input.StudentID) {
			continue
		}

		studentCount := countActiveEnrollments(classEnrollments)
		capacityLimit := CalculateCapacityLimit(lesson.Room.Capacity, lesson.Class.MaxStudents)
		remainingCapacity, utilization := CalculateUtilization(studentCount, capacityLimit)
		if !ValidateMakeupCapacity(studentCount, capacityLimit, 1) {
			continue
		}

		conflictReasons := buildStudentScheduleConflicts(studentLessons, *lesson, travelMap)
		if len(conflictReasons) > 0 {
			continue
		}

		teacherName := lesson.Teacher.FullName
		teacherID := ""
		if lesson.TeacherID != nil {
			teacherID = *lesson.TeacherID
		}
		roomID := ""
		if lesson.RoomID != nil {
			roomID = *lesson.RoomID
		}

		reasons := []string{matchType}
		if remainingCapacity > 0 {
			reasons = append(reasons, "Còn "+itoa(remainingCapacity)+" chỗ")
		}

		spots = append(spots, MakeupSpot{
			LessonID:             lesson.ID,
			ClassID:              lesson.ClassID,
			ClassCode:            lesson.Class.Code,
			ClassName:            lesson.Class.Name,
			CourseID:             derefString(lesson.Class.CourseID),
			CourseName:           lesson.Class.Course.Name,
			TeacherID:            teacherID,
			TeacherName:          teacherName,
			RoomID:               roomID,
			RoomName:             lesson.Room.Name,
			StartTime:            lesson.DateStart,
			EndTime:              lesson.DateEnd,
			MatchType:            matchType,
			ExpectedStudentCount: studentCount,
			CapacityLimit:        capacityLimit,
			RemainingCapacity:    remainingCapacity,
			CapacityUtilization:  utilization,
			Eligible:             true,
			Reasons:              reasons,
		})
	}

	sort.Slice(spots, func(i, j int) bool {
		if spots[i].MatchType != spots[j].MatchType {
			return spots[i].MatchType < spots[j].MatchType
		}
		if spots[i].RemainingCapacity != spots[j].RemainingCapacity {
			return spots[i].RemainingCapacity > spots[j].RemainingCapacity
		}
		return spots[i].StartTime.Before(spots[j].StartTime)
	})

	if len(spots) > input.Limit {
		spots = spots[:input.Limit]
	}

	return &FindMakeupSpotsOutput{
		StudentID: student.ID,
		LessonID:  sourceLesson.ID,
		Spots:     spots,
	}, nil
}

func (uc *makeupUseCase) loadClassWithCourse(ctx context.Context, classID string) (*entities.Class, error) {
	condition := repositories.NewCommonCondition()
	condition.AddCondition("id", classID, repositories.Equal)
	condition.SetPreload([]string{"Course", "Room", "Teacher"})
	condition.SetPaging(1, 1)
	page, err := uc.classRepo.GetByCondition(ctx, condition)
	if err != nil || page == nil || len(page.Data) == 0 || page.Data[0] == nil {
		return nil, err
	}
	return page.Data[0], nil
}

func (uc *makeupUseCase) loadStudentClassIDs(ctx context.Context, studentID string) ([]string, error) {
	condition := repositories.NewCommonCondition()
	condition.AddCondition("student_id", studentID, repositories.Equal)
	condition.SetPaging(500, 1)
	page, err := uc.enrollmentRepo.GetByCondition(ctx, condition)
	if err != nil || page == nil {
		return nil, err
	}

	ids := make([]string, 0)
	seen := map[string]struct{}{}
	for _, enrollment := range page.Data {
		if enrollment == nil || !isActiveEnrollmentStatus(enrollment.Status) {
			continue
		}
		if _, ok := seen[enrollment.ClassID]; ok {
			continue
		}
		seen[enrollment.ClassID] = struct{}{}
		ids = append(ids, enrollment.ClassID)
	}
	return ids, nil
}

func (uc *makeupUseCase) loadStudentLessons(ctx context.Context, classIDs []string, from, to time.Time) ([]entities.Lesson, error) {
	if len(classIDs) == 0 {
		return []entities.Lesson{}, nil
	}

	condition := repositories.NewCommonCondition()
	condition.AddCondition("class_id", classIDs, repositories.In)
	condition.AddCondition("status", []string{entities.LessonStatusPublished, entities.LessonStatusDraft}, repositories.In)
	condition.AddCondition("date_end", from.UTC(), repositories.GreaterThanOrEqual)
	condition.AddCondition("date_start", to.UTC(), repositories.LessThanOrEqual)
	condition.SetPreload([]string{"Class", "Class.Course", "Teacher", "Room"})
	condition.SetPaging(1000, 1)
	page, err := uc.lessonRepo.GetByCondition(ctx, condition)
	if err != nil || page == nil {
		return nil, err
	}

	items := make([]entities.Lesson, 0, len(page.Data))
	for _, lesson := range page.Data {
		if lesson != nil {
			items = append(items, *lesson)
		}
	}
	return items, nil
}

func (uc *makeupUseCase) loadCandidateLessons(ctx context.Context, sourceClass *entities.Class, sourceStart time.Time) ([]*entities.Lesson, error) {
	anchor := sourceStart
	now := time.Now().UTC()
	if anchor.Before(now) {
		anchor = now
	}

	condition := repositories.NewCommonCondition()
	condition.AddCondition("status", entities.LessonStatusPublished, repositories.Equal)
	condition.AddCondition("date_start", anchor, repositories.GreaterThanOrEqual)
	condition.AddCondition("date_start", anchor.AddDate(0, 0, defaultMakeupWindowDays), repositories.LessThanOrEqual)
	condition.SetPreload([]string{"Class", "Class.Course", "Teacher", "Room"})
	condition.SetPaging(1000, 1)
	condition.AddSorting("date_start", repositories.Asc)
	page, err := uc.lessonRepo.GetByCondition(ctx, condition)
	if err != nil || page == nil {
		return nil, err
	}
	return page.Data, nil
}

func (uc *makeupUseCase) loadTravelMap(ctx context.Context) map[string]int {
	travelMap := map[string]int{}
	page, err := uc.campusTravelRepo.GetByCondition(ctx, repositories.NewCommonCondition().WithPaging(1000, 1))
	if err != nil || page == nil {
		return travelMap
	}

	items := make([]entities.CampusTravelTime, 0, len(page.Data))
	for _, item := range page.Data {
		if item != nil {
			items = append(items, *item)
		}
	}
	return schedulingservice.BuildCampusTravelTimeMap(items)
}

func determineMakeupMatchType(sourceClass entities.Class, candidateClass entities.Class) string {
	switch {
	case sourceClass.CourseID != nil && candidateClass.CourseID != nil && *sourceClass.CourseID == *candidateClass.CourseID:
		return "same_course"
	case sourceClass.Course.Subject != "" && sourceClass.Course.Subject == candidateClass.Course.Subject &&
		sourceClass.Course.GradeLevel != "" && sourceClass.Course.GradeLevel == candidateClass.Course.GradeLevel:
		return "same_subject_grade"
	default:
		return ""
	}
}

func buildStudentScheduleConflicts(studentLessons []entities.Lesson, candidateLesson entities.Lesson, travelMap map[string]int) []string {
	conflicts := make([]string, 0)
	for _, lesson := range studentLessons {
		if lesson.ID == candidateLesson.ID || lesson.ClassID == candidateLesson.ClassID {
			continue
		}
		if overlaps(lesson.DateStart, lesson.DateEnd, candidateLesson.DateStart, candidateLesson.DateEnd) {
			conflicts = append(conflicts, "Trùng lịch với lesson khác của học sinh")
			continue
		}
		if sameCalendarDay(lesson.DateStart, candidateLesson.DateStart) && !hasTravelFeasibility(lesson, candidateLesson, travelMap) {
			conflicts = append(conflicts, "Không đủ thời gian di chuyển giữa 2 lesson")
		}
	}
	return conflicts
}

func countActiveEnrollments(items []entities.Enrollment) int {
	count := 0
	for _, item := range items {
		if isActiveEnrollmentStatus(item.Status) {
			count++
		}
	}
	return count
}

func isActiveEnrollmentStatus(status string) bool {
	return status == "ENROLLED" || status == "APPROVED"
}

func classHasStudent(items []entities.Enrollment, studentID string) bool {
	for _, item := range items {
		if item.StudentID == studentID && isActiveEnrollmentStatus(item.Status) {
			return true
		}
	}
	return false
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func itoa(value int) string {
	return strconv.Itoa(value)
}
