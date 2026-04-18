package teacherportal

import (
	"context"

	lessonrecord "doan/internal/usecases/lessonrecord"
)

type UpsertAcademicRecordInput struct {
	Actor              Actor
	LessonID           string
	StudentID          string
	HomeworkCompleted  bool
	HomeworkScore      float64
	AttitudeRating     int
	ParticipationScore float64
	PersonalComment    string
}

type UpsertAcademicRecordOutput struct {
	SavedCount int `json:"saved_count"`
}

type UpsertAcademicRecordUseCase interface {
	Execute(ctx context.Context, input UpsertAcademicRecordInput) (*UpsertAcademicRecordOutput, error)
}

type upsertAcademicRecordUseCase struct {
	upsertLessonAcademicRecordsUseCase lessonrecord.UpsertLessonAcademicRecordsUseCase
}

func NewUpsertAcademicRecordUseCase(
	upsertLessonAcademicRecordsUseCase lessonrecord.UpsertLessonAcademicRecordsUseCase,
) UpsertAcademicRecordUseCase {
	return &upsertAcademicRecordUseCase{
		upsertLessonAcademicRecordsUseCase: upsertLessonAcademicRecordsUseCase,
	}
}

func (uc *upsertAcademicRecordUseCase) Execute(ctx context.Context, input UpsertAcademicRecordInput) (*UpsertAcademicRecordOutput, error) {
	output, err := uc.upsertLessonAcademicRecordsUseCase.Execute(ctx, lessonrecord.UpsertLessonAcademicRecordsInput{
		LessonID: input.LessonID,
		Actor:    buildLessonRecordActor(input.Actor),
		Records: []lessonrecord.UpsertLessonAcademicRecordRow{
			{
				StudentID:          input.StudentID,
				HomeworkCompleted:  input.HomeworkCompleted,
				HomeworkScore:      input.HomeworkScore,
				AttitudeRating:     input.AttitudeRating,
				ParticipationScore: input.ParticipationScore,
				PersonalComment:    input.PersonalComment,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &UpsertAcademicRecordOutput{SavedCount: output.SavedCount}, nil
}
