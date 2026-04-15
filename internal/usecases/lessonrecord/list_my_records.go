package lessonrecord

import (
	"context"
	"strings"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type ListMyAcademicRecordsInput struct {
	Actor   LessonActor
	ClassID string
}

type ListMyAcademicRecordsOutput struct {
	Records []entities.AcademicRecord `json:"records"`
}

type ListMyAcademicRecordsUseCase interface {
	Execute(ctx context.Context, input ListMyAcademicRecordsInput) (*ListMyAcademicRecordsOutput, error)
}

type listMyAcademicRecordsUseCase struct {
	studentRepo repointerface.StudentRepository
	recordRepo  repointerface.AcademicRecordRepository
}

func NewListMyAcademicRecordsUseCase(
	studentRepo repointerface.StudentRepository,
	recordRepo repointerface.AcademicRecordRepository,
) ListMyAcademicRecordsUseCase {
	return &listMyAcademicRecordsUseCase{
		studentRepo: studentRepo,
		recordRepo:  recordRepo,
	}
}

func (uc *listMyAcademicRecordsUseCase) Execute(ctx context.Context, input ListMyAcademicRecordsInput) (*ListMyAcademicRecordsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	condition := repositories.NewCommonCondition()
	condition.AddCondition("email", strings.TrimSpace(input.Actor.Email), repositories.Equal)
	condition.SetPaging(1, 1)
	result, err := uc.studentRepo.GetByCondition(ctx, condition)
	if err != nil {
		ctxLogger.Errorf("Failed to resolve student by email %s: %v", input.Actor.Email, err)
		return nil, err
	}
	if result == nil || len(result.Data) == 0 || result.Data[0] == nil {
		return nil, ErrStudentNotFound
	}

	records, err := uc.recordRepo.ListByStudentID(ctx, result.Data[0].ID)
	if err != nil {
		ctxLogger.Errorf("Failed to load academic records for student %s: %v", result.Data[0].ID, err)
		return nil, err
	}

	filtered := make([]entities.AcademicRecord, 0, len(records))
	for _, record := range records {
		if input.ClassID != "" {
			if record.LessonSummary.Lesson.ClassID != input.ClassID {
				continue
			}
		}
		filtered = append(filtered, record)
	}

	return &ListMyAcademicRecordsOutput{Records: filtered}, nil
}
