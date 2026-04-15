package lessonrecord

import (
	"context"
	"math"
	"strings"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
)

func authorizeLessonAccess(
	ctx context.Context,
	lessonRepo repointerface.LessonRepository,
	teacherRepo repointerface.TeacherRepository,
	lessonID string,
	actor LessonActor,
) (*entities.Lesson, error) {
	lesson, err := lessonRepo.GetLessonWithRelations(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	if lesson == nil {
		return nil, ErrLessonNotFound
	}

	if actor.Role == "ADMIN" || actor.Role == "SUPER_ADMIN" {
		return lesson, nil
	}
	if actor.Role != "TEACHER" {
		return nil, ErrLessonAccessDenied
	}

	condition := repositories.NewCommonCondition()
	condition.AddCondition("email", strings.TrimSpace(actor.Email), repositories.Equal)
	condition.SetPaging(1, 1)
	result, err := teacherRepo.GetByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if result == nil || len(result.Data) == 0 || result.Data[0] == nil {
		return nil, ErrLessonAccessDenied
	}

	teacher := result.Data[0]
	if lesson.TeacherID == nil || *lesson.TeacherID != teacher.ID {
		return nil, ErrLessonAccessDenied
	}

	return lesson, nil
}

func ensureLessonSummary(
	ctx context.Context,
	summaryRepo repointerface.LessonSummaryRepository,
	lessonID string,
	actor LessonActor,
) (*entities.LessonSummary, error) {
	existing, err := summaryRepo.GetByLessonID(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	now := time.Now()
	entity := &entities.LessonSummary{
		LessonID:  lessonID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if actor.UserID != "" {
		entity.CreatedByID = &actor.UserID
	}
	return summaryRepo.Create(ctx, entity)
}

func calculateTotalScore(homeworkScore, participationScore float64, attitudeRating int) float64 {
	attitudeScore := math.Max(0, math.Min(float64(attitudeRating)*2, 10))
	total := (homeworkScore + participationScore + attitudeScore) / 3
	return math.Round(total*100) / 100
}
