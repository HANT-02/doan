package teacherportal

import (
	"context"
	"time"

	lessonactivity "doan/internal/usecases/lessonactivity"
)

type TeacherLessonSummary struct {
	ID               string     `json:"id"`
	LessonID         string     `json:"lesson_id"`
	Topic            string     `json:"topic"`
	LessonContent    string     `json:"lesson_content"`
	ClassFeedback    string     `json:"class_feedback"`
	Homework         string     `json:"homework"`
	HomeworkDeadline *time.Time `json:"homework_deadline,omitempty"`
	TeacherNotes     string     `json:"teacher_notes"`
	CreatedByID      *string    `json:"created_by_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type GetLessonSummaryInput struct {
	Actor    Actor
	LessonID string
}

type GetLessonSummaryOutput struct {
	Lesson  TeacherLessonItem     `json:"lesson"`
	Summary *TeacherLessonSummary `json:"summary,omitempty"`
}

type GetLessonSummaryUseCase interface {
	Execute(ctx context.Context, input GetLessonSummaryInput) (*GetLessonSummaryOutput, error)
}

type getLessonSummaryUseCase struct {
	lessonSummaryUseCase lessonactivity.GetLessonSummaryUseCase
}

func NewGetLessonSummaryUseCase(
	baseGetLessonSummaryUseCase lessonactivity.GetLessonSummaryUseCase,
) GetLessonSummaryUseCase {
	return &getLessonSummaryUseCase{
		lessonSummaryUseCase: baseGetLessonSummaryUseCase,
	}
}

func (uc *getLessonSummaryUseCase) Execute(ctx context.Context, input GetLessonSummaryInput) (*GetLessonSummaryOutput, error) {
	output, err := uc.lessonSummaryUseCase.Execute(ctx, lessonactivity.GetLessonSummaryInput{
		LessonID: input.LessonID,
		Actor:    buildLessonActor(input.Actor),
	})
	if err != nil {
		return nil, err
	}

	result := &GetLessonSummaryOutput{
		Lesson: TeacherLessonItem{
			ID:        output.Lesson.ID,
			ClassID:   output.Lesson.ClassID,
			ClassName: output.Lesson.Class.Name,
			ClassCode: output.Lesson.Class.Code,
			DateStart: output.Lesson.DateStart,
			DateEnd:   output.Lesson.DateEnd,
			Notes:     output.Lesson.Notes,
		},
	}
	if output.Lesson.RoomID != nil {
		result.Lesson.RoomID = output.Lesson.RoomID
		roomName := output.Lesson.Room.Name
		result.Lesson.RoomName = &roomName
	}

	if output.Summary != nil {
		summary := &TeacherLessonSummary{
			ID:            output.Summary.ID,
			LessonID:      output.Summary.LessonID,
			Topic:         output.Summary.Topic,
			LessonContent: output.Summary.LessonContent,
			ClassFeedback: output.Summary.ClassFeedback,
			Homework:      output.Summary.Homework,
			TeacherNotes:  output.Summary.TeacherNotes,
			CreatedByID:   output.Summary.CreatedByID,
			CreatedAt:     output.Summary.CreatedAt,
			UpdatedAt:     output.Summary.UpdatedAt,
		}
		if !output.Summary.HomeworkDeadline.IsZero() {
			deadline := output.Summary.HomeworkDeadline
			summary.HomeworkDeadline = &deadline
		}
		result.Summary = summary
	}

	return result, nil
}
