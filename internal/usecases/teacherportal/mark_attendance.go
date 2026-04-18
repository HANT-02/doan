package teacherportal

import (
	"context"

	lessonactivity "doan/internal/usecases/lessonactivity"
)

type MarkAttendanceInput struct {
	Actor     Actor
	LessonID  string
	StudentID string
	Status    int
	Note      string
}

type MarkAttendanceOutput struct {
	SavedCount int `json:"saved_count"`
}

type MarkAttendanceUseCase interface {
	Execute(ctx context.Context, input MarkAttendanceInput) (*MarkAttendanceOutput, error)
}

type markAttendanceUseCase struct {
	upsertLessonAttendanceUseCase lessonactivity.UpsertLessonAttendanceUseCase
}

func NewMarkAttendanceUseCase(
	upsertLessonAttendanceUseCase lessonactivity.UpsertLessonAttendanceUseCase,
) MarkAttendanceUseCase {
	return &markAttendanceUseCase{
		upsertLessonAttendanceUseCase: upsertLessonAttendanceUseCase,
	}
}

func (uc *markAttendanceUseCase) Execute(ctx context.Context, input MarkAttendanceInput) (*MarkAttendanceOutput, error) {
	internalStatus, err := mapTeacherAttendanceStatusToInternal(input.Status)
	if err != nil {
		return nil, err
	}

	output, err := uc.upsertLessonAttendanceUseCase.Execute(ctx, lessonactivity.UpsertLessonAttendanceInput{
		LessonID: input.LessonID,
		Actor:    buildLessonActor(input.Actor),
		Records: []lessonactivity.UpsertLessonAttendanceRecord{
			{
				StudentID: input.StudentID,
				Status:    internalStatus,
				Note:      input.Note,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &MarkAttendanceOutput{SavedCount: output.SavedCount}, nil
}

type SubmitLessonAttendanceRecord struct {
	StudentID string `json:"student_id"`
	Status    int    `json:"status"`
	Note      string `json:"note"`
}

type SubmitLessonAttendanceInput struct {
	Actor    Actor
	LessonID string
	Records  []SubmitLessonAttendanceRecord
}

type SubmitLessonAttendanceOutput struct {
	SavedCount int `json:"saved_count"`
}

type SubmitLessonAttendanceUseCase interface {
	Execute(ctx context.Context, input SubmitLessonAttendanceInput) (*SubmitLessonAttendanceOutput, error)
}

type submitLessonAttendanceUseCase struct {
	upsertLessonAttendanceUseCase lessonactivity.UpsertLessonAttendanceUseCase
}

func NewSubmitLessonAttendanceUseCase(
	upsertLessonAttendanceUseCase lessonactivity.UpsertLessonAttendanceUseCase,
) SubmitLessonAttendanceUseCase {
	return &submitLessonAttendanceUseCase{
		upsertLessonAttendanceUseCase: upsertLessonAttendanceUseCase,
	}
}

func (uc *submitLessonAttendanceUseCase) Execute(ctx context.Context, input SubmitLessonAttendanceInput) (*SubmitLessonAttendanceOutput, error) {
	records := make([]lessonactivity.UpsertLessonAttendanceRecord, 0, len(input.Records))
	for _, record := range input.Records {
		internalStatus, err := mapTeacherAttendanceStatusToInternal(record.Status)
		if err != nil {
			return nil, err
		}
		records = append(records, lessonactivity.UpsertLessonAttendanceRecord{
			StudentID: record.StudentID,
			Status:    internalStatus,
			Note:      record.Note,
		})
	}

	output, err := uc.upsertLessonAttendanceUseCase.Execute(ctx, lessonactivity.UpsertLessonAttendanceInput{
		LessonID: input.LessonID,
		Actor:    buildLessonActor(input.Actor),
		Records:  records,
	})
	if err != nil {
		return nil, err
	}

	return &SubmitLessonAttendanceOutput{SavedCount: output.SavedCount}, nil
}
