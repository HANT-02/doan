package lesson

import (
	"context"
	"errors"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
)

type GetLessonInput struct {
	ID string
}

type GetLessonOutput struct {
	Lesson *entities.Lesson `json:"lesson"`
}

type GetLessonUseCase interface {
	Execute(ctx context.Context, input GetLessonInput) (*GetLessonOutput, error)
}

type getLessonUseCase struct {
	lessonRepo repointerface.LessonRepository
}

func NewGetLessonUseCase(lessonRepo repointerface.LessonRepository) GetLessonUseCase {
	return &getLessonUseCase{lessonRepo: lessonRepo}
}

func (uc *getLessonUseCase) Execute(ctx context.Context, input GetLessonInput) (*GetLessonOutput, error) {
	lesson, err := uc.lessonRepo.GetLessonWithRelations(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if lesson == nil {
		return nil, errors.New("lesson not found")
	}
	return &GetLessonOutput{Lesson: lesson}, nil
}
