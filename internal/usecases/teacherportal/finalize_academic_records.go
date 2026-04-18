package teacherportal

import (
	"context"

	lessonrecord "doan/internal/usecases/lessonrecord"
)

type FinalizeAcademicRecordInput struct {
	Actor    Actor
	LessonID string
}

type FinalizeAcademicRecordOutput struct {
	FinalizedCount int `json:"finalized_count"`
}

type FinalizeAcademicRecordUseCase interface {
	Execute(ctx context.Context, input FinalizeAcademicRecordInput) (*FinalizeAcademicRecordOutput, error)
}

type finalizeAcademicRecordUseCase struct {
	finalizeLessonAcademicRecordsUseCase lessonrecord.FinalizeLessonAcademicRecordsUseCase
}

func NewFinalizeAcademicRecordUseCase(
	finalizeLessonAcademicRecordsUseCase lessonrecord.FinalizeLessonAcademicRecordsUseCase,
) FinalizeAcademicRecordUseCase {
	return &finalizeAcademicRecordUseCase{
		finalizeLessonAcademicRecordsUseCase: finalizeLessonAcademicRecordsUseCase,
	}
}

func (uc *finalizeAcademicRecordUseCase) Execute(ctx context.Context, input FinalizeAcademicRecordInput) (*FinalizeAcademicRecordOutput, error) {
	output, err := uc.finalizeLessonAcademicRecordsUseCase.Execute(ctx, lessonrecord.FinalizeLessonAcademicRecordsInput{
		LessonID: input.LessonID,
		Actor:    buildLessonRecordActor(input.Actor),
	})
	if err != nil {
		return nil, err
	}

	return &FinalizeAcademicRecordOutput{
		FinalizedCount: output.FinalizedCount,
	}, nil
}
