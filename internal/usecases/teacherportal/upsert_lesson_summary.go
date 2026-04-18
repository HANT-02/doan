package teacherportal

import (
	"context"
	"time"

	lessonactivity "doan/internal/usecases/lessonactivity"
)

type UpsertLessonSummaryInput struct {
	Actor            Actor
	LessonID         string
	Topic            string
	LessonContent    string
	ClassFeedback    string
	Homework         string
	HomeworkDeadline *time.Time
	TeacherNotes     string
}

type UpsertLessonSummaryOutput struct {
	Summary TeacherLessonSummary `json:"summary"`
}

type UpsertLessonSummaryUseCase interface {
	Execute(ctx context.Context, input UpsertLessonSummaryInput) (*UpsertLessonSummaryOutput, error)
}

type upsertLessonSummaryUseCase struct {
	lessonSummaryUseCase lessonactivity.UpsertLessonSummaryUseCase
}

func NewUpsertLessonSummaryUseCase(
	baseUpsertLessonSummaryUseCase lessonactivity.UpsertLessonSummaryUseCase,
) UpsertLessonSummaryUseCase {
	return &upsertLessonSummaryUseCase{
		lessonSummaryUseCase: baseUpsertLessonSummaryUseCase,
	}
}

func (uc *upsertLessonSummaryUseCase) Execute(ctx context.Context, input UpsertLessonSummaryInput) (*UpsertLessonSummaryOutput, error) {
	output, err := uc.lessonSummaryUseCase.Execute(ctx, lessonactivity.UpsertLessonSummaryInput{
		LessonID:         input.LessonID,
		Actor:            buildLessonActor(input.Actor),
		Topic:            input.Topic,
		LessonContent:    input.LessonContent,
		ClassFeedback:    input.ClassFeedback,
		Homework:         input.Homework,
		HomeworkDeadline: input.HomeworkDeadline,
		TeacherNotes:     input.TeacherNotes,
	})
	if err != nil {
		return nil, err
	}

	summary := TeacherLessonSummary{
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

	return &UpsertLessonSummaryOutput{
		Summary: summary,
	}, nil
}
