package studentportal

import (
	"context"
	"strings"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
)

func resolveStudentByEmail(ctx context.Context, studentRepo repointerface.StudentRepository, email string) (*entities.Student, error) {
	condition := repositories.NewCommonCondition()
	condition.AddCondition("email", strings.TrimSpace(email), repositories.Equal)
	condition.SetPaging(1, 1)

	result, err := studentRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 || result.Data[0] == nil {
		return nil, ErrStudentNotFound
	}

	return result.Data[0], nil
}

func normalizeDayOfWeek(value time.Weekday) string {
	switch value {
	case time.Monday:
		return "MONDAY"
	case time.Tuesday:
		return "TUESDAY"
	case time.Wednesday:
		return "WEDNESDAY"
	case time.Thursday:
		return "THURSDAY"
	case time.Friday:
		return "FRIDAY"
	case time.Saturday:
		return "SATURDAY"
	case time.Sunday:
		return "SUNDAY"
	default:
		return ""
	}
}

func matchShiftForLesson(lesson entities.Lesson, schedules []entities.ClassSchedule) *entities.ClassSchedule {
	lessonDay := normalizeDayOfWeek(lesson.DateStart.Weekday())
	lessonStart := lesson.DateStart.Format("15:04")
	lessonEnd := lesson.DateEnd.Format("15:04")

	var fallback *entities.ClassSchedule
	for idx := range schedules {
		schedule := &schedules[idx]
		if strings.TrimSpace(strings.ToUpper(schedule.DayOfWeek)) != lessonDay {
			continue
		}

		if fallback == nil {
			fallback = schedule
		}

		if schedule.Shift.StartTime == lessonStart && schedule.Shift.EndTime == lessonEnd {
			if lesson.RoomID == nil || schedule.RoomID == nil || *schedule.RoomID == *lesson.RoomID {
				return schedule
			}
		}

		if lesson.RoomID != nil && schedule.RoomID != nil && *schedule.RoomID == *lesson.RoomID {
			return schedule
		}
	}

	return fallback
}
