package teacher

import (
	"context"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
	"time"
)

// GetTeacherTimetableInput represents the input for getting teacher timetable
type GetTeacherTimetableInput struct {
	TeacherID string    `json:"teacher_id"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

// TimetableLesson represents a lesson in the timetable
type TimetableLesson struct {
	ID        string    `json:"id"`
	ClassID   string    `json:"class_id"`
	ClassName string    `json:"class_name"`
	RoomID    *string   `json:"room_id"`
	RoomName  *string   `json:"room_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Notes     string    `json:"notes"`
}

// GetTeacherTimetableOutput represents the output after getting teacher timetable
type GetTeacherTimetableOutput struct {
	Lessons []TimetableLesson `json:"lessons"`
}

// GetTeacherTimetableUseCase defines the interface for getting teacher timetable
type GetTeacherTimetableUseCase interface {
	Execute(ctx context.Context, input GetTeacherTimetableInput) (*GetTeacherTimetableOutput, error)
}

type getTeacherTimetableUseCase struct {
	teacherRepo repointerface.TeacherRepository
}

// NewGetTeacherTimetableUseCase creates a new instance of GetTeacherTimetableUseCase
func NewGetTeacherTimetableUseCase(teacherRepo repointerface.TeacherRepository) GetTeacherTimetableUseCase {
	return &getTeacherTimetableUseCase{
		teacherRepo: teacherRepo,
	}
}

func (uc *getTeacherTimetableUseCase) Execute(ctx context.Context, input GetTeacherTimetableInput) (*GetTeacherTimetableOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.TeacherID == "" {
		return nil, errors.New("teacher ID is required")
	}

	// Verify teacher exists
	_, err := uc.teacherRepo.GetByID(ctx, input.TeacherID)
	if err != nil {
		ctxLogger.Errorf("Failed to get teacher: %v", err)
		return nil, err
	}

	// Get lessons
	lessons, err := uc.teacherRepo.GetTeacherLessons(ctx, input.TeacherID, input.From, input.To)
	if err != nil {
		ctxLogger.Errorf("Failed to get teacher lessons: %v", err)
		return nil, err
	}

	// Map to timetable lessons
	timetableLessons := make([]TimetableLesson, 0, len(lessons))
	for _, lesson := range lessons {
		tl := TimetableLesson{
			ID:        lesson.ID,
			ClassID:   lesson.ClassID,
			StartTime: lesson.DateStart,
			EndTime:   lesson.DateEnd,
			Notes:     lesson.Notes,
		}

		// Add class name if available
		if lesson.Class.Name != "" {
			tl.ClassName = lesson.Class.Name
		}

		// Add room info if available
		if lesson.RoomID != nil {
			tl.RoomID = lesson.RoomID
			if lesson.Room.Name != "" {
				roomName := lesson.Room.Name
				tl.RoomName = &roomName
			}
		}

		timetableLessons = append(timetableLessons, tl)
	}

	return &GetTeacherTimetableOutput{Lessons: timetableLessons}, nil
}
