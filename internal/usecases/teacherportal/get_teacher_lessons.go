package teacherportal

import (
	"context"
	"sort"
	"strings"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetTeacherLessonsInput struct {
	Actor    Actor
	ClassID  string
	DateFrom *time.Time
	DateTo   *time.Time
}

type TeacherLessonShift struct {
	ID              string `json:"id"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	DurationMinutes int    `json:"duration_minutes"`
	SessionType     string `json:"session_type"`
}

type TeacherLessonItem struct {
	ID        string              `json:"id"`
	ClassID   string              `json:"class_id"`
	ClassName string              `json:"class_name"`
	ClassCode string              `json:"class_code"`
	RoomID    *string             `json:"room_id,omitempty"`
	RoomName  *string             `json:"room_name,omitempty"`
	DateStart time.Time           `json:"date_start"`
	DateEnd   time.Time           `json:"date_end"`
	Notes     string              `json:"notes"`
	Shift     *TeacherLessonShift `json:"shift,omitempty"`
}

type GetTeacherLessonsOutput struct {
	TeacherID string              `json:"teacher_id"`
	Lessons   []TeacherLessonItem `json:"lessons"`
}

type GetTeacherLessonsUseCase interface {
	Execute(ctx context.Context, input GetTeacherLessonsInput) (*GetTeacherLessonsOutput, error)
}

type getTeacherLessonsUseCase struct {
	teacherRepo       repointerface.TeacherRepository
	classScheduleRepo repointerface.ClassScheduleRepository
}

func NewGetTeacherLessonsUseCase(
	teacherRepo repointerface.TeacherRepository,
	classScheduleRepo repointerface.ClassScheduleRepository,
) GetTeacherLessonsUseCase {
	return &getTeacherLessonsUseCase{
		teacherRepo:       teacherRepo,
		classScheduleRepo: classScheduleRepo,
	}
}

func (uc *getTeacherLessonsUseCase) Execute(ctx context.Context, input GetTeacherLessonsInput) (*GetTeacherLessonsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if strings.TrimSpace(input.Actor.Role) != "TEACHER" {
		return nil, ErrTeacherAccessDenied
	}

	teacher, err := resolveTeacherByEmail(ctx, uc.teacherRepo, input.Actor.Email)
	if err != nil {
		ctxLogger.Errorf("Failed to resolve teacher from actor email %s: %v", input.Actor.Email, err)
		return nil, err
	}

	var from time.Time
	if input.DateFrom != nil {
		from = input.DateFrom.UTC()
	}

	var to time.Time
	if input.DateTo != nil {
		to = input.DateTo.UTC()
	}

	lessons, err := uc.teacherRepo.GetTeacherLessons(ctx, teacher.ID, from, to)
	if err != nil {
		ctxLogger.Errorf("Failed to get teacher lessons %s: %v", teacher.ID, err)
		return nil, err
	}

	classScheduleMap := make(map[string][]entities.ClassSchedule)
	items := make([]TeacherLessonItem, 0, len(lessons))

	for _, lesson := range lessons {
		if input.ClassID != "" && lesson.ClassID != input.ClassID {
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

		item := TeacherLessonItem{
			ID:        lesson.ID,
			ClassID:   lesson.ClassID,
			ClassName: lesson.Class.Name,
			ClassCode: lesson.Class.Code,
			DateStart: lesson.DateStart,
			DateEnd:   lesson.DateEnd,
			Notes:     lesson.Notes,
			RoomID:    lesson.RoomID,
		}

		if lesson.Room.Name != "" {
			roomName := lesson.Room.Name
			item.RoomName = &roomName
		}

		if matched := matchShiftForLesson(lesson, schedules); matched != nil {
			item.Shift = &TeacherLessonShift{
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

	sort.Slice(items, func(i, j int) bool {
		return items[i].DateStart.Before(items[j].DateStart)
	})

	return &GetTeacherLessonsOutput{
		TeacherID: teacher.ID,
		Lessons:   items,
	}, nil
}
