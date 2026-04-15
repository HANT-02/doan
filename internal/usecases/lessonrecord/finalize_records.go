package lessonrecord

import (
	"context"

	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type FinalizeLessonAcademicRecordsInput struct {
	LessonID string
	Actor    LessonActor
}

type FinalizeLessonAcademicRecordsOutput struct {
	FinalizedCount int `json:"finalized_count"`
}

type FinalizeLessonAcademicRecordsUseCase interface {
	Execute(ctx context.Context, input FinalizeLessonAcademicRecordsInput) (*FinalizeLessonAcademicRecordsOutput, error)
}

type finalizeLessonAcademicRecordsUseCase struct {
	lessonRepo  repointerface.LessonRepository
	teacherRepo repointerface.TeacherRepository
	summaryRepo repointerface.LessonSummaryRepository
	recordRepo  repointerface.AcademicRecordRepository
}

func NewFinalizeLessonAcademicRecordsUseCase(
	lessonRepo repointerface.LessonRepository,
	teacherRepo repointerface.TeacherRepository,
	summaryRepo repointerface.LessonSummaryRepository,
	recordRepo repointerface.AcademicRecordRepository,
) FinalizeLessonAcademicRecordsUseCase {
	return &finalizeLessonAcademicRecordsUseCase{
		lessonRepo:  lessonRepo,
		teacherRepo: teacherRepo,
		summaryRepo: summaryRepo,
		recordRepo:  recordRepo,
	}
}

func (uc *finalizeLessonAcademicRecordsUseCase) Execute(ctx context.Context, input FinalizeLessonAcademicRecordsInput) (*FinalizeLessonAcademicRecordsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if _, err := authorizeLessonAccess(ctx, uc.lessonRepo, uc.teacherRepo, input.LessonID, input.Actor); err != nil {
		ctxLogger.Errorf("Failed to authorize lesson academic finalize %s: %v", input.LessonID, err)
		return nil, err
	}

	summary, err := ensureLessonSummary(ctx, uc.summaryRepo, input.LessonID, input.Actor)
	if err != nil {
		return nil, err
	}

	records, err := uc.recordRepo.ListByLessonSummaryID(ctx, summary.ID)
	if err != nil {
		return nil, err
	}

	finalizedCount := 0
	for _, record := range records {
		if err := uc.recordRepo.Update(ctx, record.ID, map[string]interface{}{
			"is_completed": true,
			"total_score":  calculateTotalScore(record.HomeworkScore, record.ParticipationScore, record.AttitudeRating),
		}); err != nil {
			return nil, err
		}
		finalizedCount++
	}

	return &FinalizeLessonAcademicRecordsOutput{FinalizedCount: finalizedCount}, nil
}
