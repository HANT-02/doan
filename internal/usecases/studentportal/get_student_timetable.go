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

type GetStudentTimetableInput struct {
	Actor    Actor
	ClassID  string
	DateFrom *time.Time
	DateTo   *time.Time
}

type StudentTimetableShift struct {
	ID              string `json:"id"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	DurationMinutes int    `json:"duration_minutes"`
	SessionType     string `json:"session_type"`
}

type StudentTimetableTeacher struct {
	ID       *string `json:"id,omitempty"`
	Code     *string `json:"code,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

type StudentTimetableLessonItem struct {
	ID        string                  `json:"id"`
	ClassID   string                  `json:"class_id"`
	ClassName string                  `json:"class_name"`
	ClassCode string                  `json:"class_code"`
	Teacher   StudentTimetableTeacher `json:"teacher"`
	RoomID    *string                 `json:"room_id,omitempty"`
	RoomName  *string                 `json:"room_name,omitempty"`
	DateStart time.Time               `json:"date_start"`
	DateEnd   time.Time               `json:"date_end"`
	Notes     string                  `json:"notes"`
	Shift     *StudentTimetableShift  `json:"shift,omitempty"`
}

type GetStudentTimetableOutput struct {
	StudentID string                       `json:"student_id"`
	Lessons   []StudentTimetableLessonItem `json:"lessons"`
}

type GetStudentTimetableUseCase interface {
	Execute(ctx context.Context, input GetStudentTimetableInput) (*GetStudentTimetableOutput, error)
}

type getStudentTimetableUseCase struct {
	studentRepo       repointerface.StudentRepository
	enrollmentRepo    repointerface.EnrollmentRepository
	lessonRepo        repointerface.LessonRepository
	classScheduleRepo repointerface.ClassScheduleRepository
}

func NewGetStudentTimetableUseCase(
	studentRepo repointerface.StudentRepository,
	enrollmentRepo repointerface.EnrollmentRepository,
	lessonRepo repointerface.LessonRepository,
	classScheduleRepo repointerface.ClassScheduleRepository,
) GetStudentTimetableUseCase {
	return &getStudentTimetableUseCase{
		studentRepo:       studentRepo,
		enrollmentRepo:    enrollmentRepo,
		lessonRepo:        lessonRepo,
		classScheduleRepo: classScheduleRepo,
	}
}

func (uc *getStudentTimetableUseCase) Execute(ctx context.Context, input GetStudentTimetableInput) (*GetStudentTimetableOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	role := strings.TrimSpace(input.Actor.Role)
	if role != "STUDENT" {
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
		return &GetStudentTimetableOutput{
			StudentID: student.ID,
			Lessons:   []StudentTimetableLessonItem{},
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
		ctxLogger.Errorf("Failed to get lessons for student %s: %v", student.ID, err)
		return nil, err
	}

	classScheduleMap := make(map[string][]entities.ClassSchedule)
	items := make([]StudentTimetableLessonItem, 0)

	if lessonsResult != nil {
		for _, lesson := range lessonsResult.Data {
			if lesson == nil {
				continue
			}

			schedules, exists := classScheduleMap[lesson.ClassID]
			if !exists {
				schedules, err = uc.classScheduleRepo.GetSchedulesByClassID(ctx, lesson.ClassID)
				if err != nil {
					ctxLogger.Errorf("Failed to load class schedules for class %s: %v", lesson.ClassID, err)
					return nil, err
				}
				classScheduleMap[lesson.ClassID] = schedules
			}

			item := StudentTimetableLessonItem{
				ID:        lesson.ID,
				ClassID:   lesson.ClassID,
				ClassName: lesson.Class.Name,
				ClassCode: lesson.Class.Code,
				DateStart: lesson.DateStart,
				DateEnd:   lesson.DateEnd,
				Notes:     lesson.Notes,
				RoomID:    lesson.RoomID,
				Teacher:   StudentTimetableTeacher{},
			}

			if lesson.TeacherID != nil {
				item.Teacher.ID = lesson.TeacherID
			}
			if lesson.Teacher.Code != "" {
				code := lesson.Teacher.Code
				item.Teacher.Code = &code
			}
			if lesson.Teacher.FullName != "" {
				name := lesson.Teacher.FullName
				item.Teacher.FullName = &name
			}
			if lesson.Room.Name != "" {
				roomName := lesson.Room.Name
				item.RoomName = &roomName
			}

			if matched := matchShiftForLesson(*lesson, schedules); matched != nil {
				item.Shift = &StudentTimetableShift{
					ID:              matched.Shift.ID,
					Code:            matched.Shift.Code,
					Name:            matched.Shift.Name,
					StartTime:       matched.Shift.StartTime,
					EndTime:         matched.Shift.EndTime,
					DurationMinutes: matched.Shift.DurationMinutes,
					SessionType:     matched.Shift.SessionType,
				}
				if item.RoomID == nil && matched.RoomID != nil {
					item.RoomID = matched.RoomID
				}
				if item.RoomName == nil && matched.Room.Name != "" {
					roomName := matched.Room.Name
					item.RoomName = &roomName
				}
			}

			items = append(items, item)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].DateStart.Before(items[j].DateStart)
	})

	return &GetStudentTimetableOutput{
		StudentID: student.ID,
		Lessons:   items,
	}, nil
}
