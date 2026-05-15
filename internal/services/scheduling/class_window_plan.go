package scheduling

import (
	"time"

	"doan/internal/entities"
)

type ClassSchedulingWindow struct {
	DateFrom     time.Time `json:"date_from"`
	DateTo       time.Time `json:"date_to"`
	SessionTotal int       `json:"session_total"`
}

func BuildClassSchedulingWindows(
	globalFrom time.Time,
	classes []entities.Class,
	shifts []entities.Shift,
	now time.Time,
) map[string]ClassSchedulingWindow {
	windows := make(map[string]ClassSchedulingWindow, len(classes))
	today := startOfDay(now)
	baseFrom := startOfDay(globalFrom)
	if baseFrom.IsZero() || baseFrom.Before(today) {
		baseFrom = today
	}

	for _, classEntity := range classes {
		if classEntity.ID == "" || classEntity.Course.SessionCount <= 0 {
			continue
		}

		windowStart := baseFrom
		if !classEntity.StartDate.IsZero() {
			classStart := startOfDay(classEntity.StartDate)
			if classStart.After(windowStart) {
				windowStart = classStart
			}
		}

		windowEnd := estimateClassWindowEnd(
			windowStart,
			classEntity.Course.SessionCount,
			classEntity.Course.SessionDurationMinutes,
			classEntity.ClassSchedules,
			shifts,
		)
		if windowEnd.IsZero() {
			continue
		}

		windows[classEntity.ID] = ClassSchedulingWindow{
			DateFrom:     windowStart,
			DateTo:       windowEnd,
			SessionTotal: classEntity.Course.SessionCount,
		}
	}

	return windows
}

func AggregateClassSchedulingWindow(
	fallbackFrom, fallbackTo time.Time,
	windows map[string]ClassSchedulingWindow,
) (time.Time, time.Time) {
	from := startOfDay(fallbackFrom)
	to := startOfDay(fallbackTo)

	for _, window := range windows {
		if window.DateFrom.IsZero() || window.DateTo.IsZero() {
			continue
		}
		if from.IsZero() || window.DateFrom.Before(from) {
			from = window.DateFrom
		}
		if to.IsZero() || window.DateTo.After(to) {
			to = window.DateTo
		}
	}

	if to.IsZero() || to.Before(from) {
		to = from
	}
	return from, to
}

func estimateClassWindowEnd(
	windowStart time.Time,
	sessionTotal int,
	durationMinutes int,
	classSchedules []entities.ClassSchedule,
	shifts []entities.Shift,
) time.Time {
	if windowStart.IsZero() || sessionTotal <= 0 || durationMinutes <= 0 || len(classSchedules) == 0 {
		return time.Time{}
	}

	weeklySlots := countActiveClassScheduleSlots(classSchedules)
	if weeklySlots <= 0 {
		return time.Time{}
	}

	weeksNeeded := (sessionTotal + weeklySlots - 1) / weeklySlots
	horizonEnd := windowStart.AddDate(0, 0, weeksNeeded*7+14)
	slots := generateTimeSlotsForVariable(windowStart, horizonEnd, durationMinutes, classSchedules, shifts)
	if len(slots) < sessionTotal {
		extendedEnd := windowStart.AddDate(2, 0, 0)
		if extendedEnd.After(horizonEnd) {
			slots = generateTimeSlotsForVariable(windowStart, extendedEnd, durationMinutes, classSchedules, shifts)
		}
	}
	if len(slots) < sessionTotal {
		return time.Time{}
	}

	return startOfDay(slots[sessionTotal-1].Start)
}

func countActiveClassScheduleSlots(classSchedules []entities.ClassSchedule) int {
	if len(classSchedules) == 0 {
		return 0
	}

	seen := make(map[string]struct{}, len(classSchedules))
	for _, schedule := range classSchedules {
		if schedule.ShiftID == "" || schedule.Shift.ID == "" || !schedule.Shift.IsActive {
			continue
		}
		key := schedule.DayOfWeek + "|" + schedule.Shift.ID
		seen[key] = struct{}{}
	}

	return len(seen)
}

func resolveClassSchedulingWindow(input SolverInput, classEntity entities.Class) ClassSchedulingWindow {
	if input.ClassWindows != nil {
		if window, ok := input.ClassWindows[classEntity.ID]; ok {
			return window
		}
	}

	sessionTotal := classEntity.Course.SessionCount
	if sessionTotal <= 0 {
		sessionTotal = countUniqueTimeSlots(generateTimeSlotsForVariable(
			input.DateFrom,
			input.DateTo,
			classEntity.Course.SessionDurationMinutes,
			classEntity.ClassSchedules,
			input.Shifts,
		))
	}

	return ClassSchedulingWindow{
		DateFrom:     input.DateFrom,
		DateTo:       input.DateTo,
		SessionTotal: sessionTotal,
	}
}
