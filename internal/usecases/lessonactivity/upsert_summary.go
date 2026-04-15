package lessonactivity

import (
	"context"
	"strings"
	"time"

	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type UpsertLessonSummaryInput struct {
	LessonID         string
	Actor            LessonActor
	Topic            string
	LessonContent    string
	ClassFeedback    string
	Homework         string
	HomeworkDeadline *time.Time
	TeacherNotes     string
}

type UpsertLessonSummaryOutput struct {
	Summary entities.LessonSummary `json:"summary"`
}

type UpsertLessonSummaryUseCase interface {
	Execute(ctx context.Context, input UpsertLessonSummaryInput) (*UpsertLessonSummaryOutput, error)
}

type upsertLessonSummaryUseCase struct {
	lessonRepo  repointerface.LessonRepository
	teacherRepo repointerface.TeacherRepository
	summaryRepo repointerface.LessonSummaryRepository
}

func NewUpsertLessonSummaryUseCase(
	lessonRepo repointerface.LessonRepository,
	teacherRepo repointerface.TeacherRepository,
	summaryRepo repointerface.LessonSummaryRepository,
) UpsertLessonSummaryUseCase {
	return &upsertLessonSummaryUseCase{
		lessonRepo:  lessonRepo,
		teacherRepo: teacherRepo,
		summaryRepo: summaryRepo,
	}
}

func (uc *upsertLessonSummaryUseCase) Execute(ctx context.Context, input UpsertLessonSummaryInput) (*UpsertLessonSummaryOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	_, err := authorizeLessonAccess(ctx, uc.lessonRepo, uc.teacherRepo, input.LessonID, input.Actor)
	if err != nil {
		ctxLogger.Errorf("Failed to authorize lesson summary update %s: %v", input.LessonID, err)
		return nil, err
	}

	existing, err := uc.summaryRepo.GetByLessonID(ctx, input.LessonID)
	if err != nil {
		ctxLogger.Errorf("Failed to read existing summary for lesson %s: %v", input.LessonID, err)
		return nil, err
	}

	homeworkDeadline := time.Time{}
	if input.HomeworkDeadline != nil {
		homeworkDeadline = *input.HomeworkDeadline
	}

	if existing == nil {
		newSummary := &entities.LessonSummary{
			LessonID:         input.LessonID,
			Topic:            strings.TrimSpace(input.Topic),
			LessonContent:    strings.TrimSpace(input.LessonContent),
			ClassFeedback:    strings.TrimSpace(input.ClassFeedback),
			Homework:         strings.TrimSpace(input.Homework),
			HomeworkDeadline: homeworkDeadline,
			TeacherNotes:     strings.TrimSpace(input.TeacherNotes),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if input.Actor.UserID != "" {
			newSummary.CreatedByID = &input.Actor.UserID
		}

		created, createErr := uc.summaryRepo.Create(ctx, newSummary)
		if createErr != nil {
			ctxLogger.Errorf("Failed to create lesson summary %s: %v", input.LessonID, createErr)
			return nil, createErr
		}
		return &UpsertLessonSummaryOutput{Summary: *created}, nil
	}

	updateData := map[string]interface{}{
		"topic":             strings.TrimSpace(input.Topic),
		"lesson_content":    strings.TrimSpace(input.LessonContent),
		"class_feedback":    strings.TrimSpace(input.ClassFeedback),
		"homework":          strings.TrimSpace(input.Homework),
		"homework_deadline": homeworkDeadline,
		"teacher_notes":     strings.TrimSpace(input.TeacherNotes),
	}
	if existing.CreatedByID == nil && input.Actor.UserID != "" {
		updateData["created_by_id"] = input.Actor.UserID
	}

	if err := uc.summaryRepo.Update(ctx, existing.ID, updateData); err != nil {
		ctxLogger.Errorf("Failed to update lesson summary %s: %v", input.LessonID, err)
		return nil, err
	}

	summary, err := uc.summaryRepo.GetByLessonID(ctx, input.LessonID)
	if err != nil {
		return nil, err
	}

	return &UpsertLessonSummaryOutput{Summary: *summary}, nil
}
