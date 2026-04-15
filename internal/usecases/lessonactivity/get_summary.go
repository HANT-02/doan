package lessonactivity

import (
	"context"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetLessonSummaryInput struct {
	LessonID string
	Actor    LessonActor
}

type GetLessonSummaryOutput struct {
	Lesson  entities.Lesson         `json:"lesson"`
	Summary *entities.LessonSummary `json:"summary"`
}

type GetLessonSummaryUseCase interface {
	Execute(ctx context.Context, input GetLessonSummaryInput) (*GetLessonSummaryOutput, error)
}

type getLessonSummaryUseCase struct {
	lessonRepo  repointerface.LessonRepository
	teacherRepo repointerface.TeacherRepository
	summaryRepo repointerface.LessonSummaryRepository
}

func NewGetLessonSummaryUseCase(
	lessonRepo repointerface.LessonRepository,
	teacherRepo repointerface.TeacherRepository,
	summaryRepo repointerface.LessonSummaryRepository,
) GetLessonSummaryUseCase {
	return &getLessonSummaryUseCase{
		lessonRepo:  lessonRepo,
		teacherRepo: teacherRepo,
		summaryRepo: summaryRepo,
	}
}

func (uc *getLessonSummaryUseCase) Execute(ctx context.Context, input GetLessonSummaryInput) (*GetLessonSummaryOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	lesson, err := authorizeLessonAccess(ctx, uc.lessonRepo, uc.teacherRepo, input.LessonID, input.Actor)
	if err != nil {
		ctxLogger.Errorf("Failed to authorize summary access for lesson %s: %v", input.LessonID, err)
		return nil, err
	}

	summary, err := uc.summaryRepo.GetByLessonID(ctx, input.LessonID)
	if err != nil {
		ctxLogger.Errorf("Failed to load lesson summary %s: %v", input.LessonID, err)
		return nil, err
	}

	return &GetLessonSummaryOutput{
		Lesson:  *lesson,
		Summary: summary,
	}, nil
}
