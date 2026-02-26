package course

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type UpdateCourseInput struct {
	ID                     string   `json:"id"`
	Code                   *string  `json:"code"`
	Name                   *string  `json:"name"`
	Description            *string  `json:"description"`
	GradeLevel             *string  `json:"grade_level"`
	Subject                *string  `json:"subject"`
	SessionCount           *int     `json:"session_count"`
	SessionDurationMinutes *int     `json:"session_duration_minutes"`
	TotalHours             *float64 `json:"total_hours"`
	Price                  *float64 `json:"price"`
	Status                 *string  `json:"status"`
}

type UpdateCourseOutput struct {
	Course *entities.Course `json:"course"`
}

type UpdateCourseUseCase interface {
	Execute(ctx context.Context, input UpdateCourseInput) (*UpdateCourseOutput, error)
}

type updateCourseUseCaseImpl struct {
	repo repointerface.CourseRepository
}

func NewUpdateCourseUseCase(repo repointerface.CourseRepository) UpdateCourseUseCase {
	return &updateCourseUseCaseImpl{repo: repo}
}

func (uc *updateCourseUseCaseImpl) Execute(ctx context.Context, input UpdateCourseInput) (*UpdateCourseOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	course, err := uc.repo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Error getting course by ID for update: %v", err)
		return nil, err
	}

	updateData := map[string]interface{}{}

	if input.Code != nil {
		course.Code = *input.Code
		updateData["code"] = *input.Code
	}
	if input.Name != nil {
		course.Name = *input.Name
		updateData["name"] = *input.Name
	}
	if input.Description != nil {
		course.Description = *input.Description
		updateData["description"] = *input.Description
	}
	if input.GradeLevel != nil {
		course.GradeLevel = *input.GradeLevel
		updateData["grade_level"] = *input.GradeLevel
	}
	if input.Subject != nil {
		course.Subject = *input.Subject
		updateData["subject"] = *input.Subject
	}
	if input.SessionCount != nil {
		course.SessionCount = *input.SessionCount
		updateData["session_count"] = *input.SessionCount
	}
	if input.SessionDurationMinutes != nil {
		course.SessionDurationMinutes = *input.SessionDurationMinutes
		updateData["session_duration_minutes"] = *input.SessionDurationMinutes
	}
	if input.TotalHours != nil {
		course.TotalHours = *input.TotalHours
		updateData["total_hours"] = *input.TotalHours
	}
	if input.Price != nil {
		course.Price = *input.Price
		updateData["price"] = *input.Price
	}
	if input.Status != nil {
		course.Status = *input.Status
		updateData["status"] = *input.Status
	}

	err = uc.repo.Update(ctx, course.ID, updateData)
	if err != nil {
		ctxLogger.Errorf("Error updating course: %v", err)
		return nil, err
	}
	return &UpdateCourseOutput{Course: course}, nil
}
