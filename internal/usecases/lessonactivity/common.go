package lessonactivity

import (
	"context"
	"strings"

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
		return nil, ErrTeacherProfileMissing
	}

	teacher := result.Data[0]
	if lesson.TeacherID == nil || *lesson.TeacherID != teacher.ID {
		return nil, ErrLessonAccessDenied
	}

	return lesson, nil
}
