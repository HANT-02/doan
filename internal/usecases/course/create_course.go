package course

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
)

type CreateCourseInput struct {
	Code                   string  `json:"code"`
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	GradeLevel             string  `json:"grade_level"`
	Subject                string  `json:"subject"`
	SessionCount           int     `json:"session_count"`
	SessionDurationMinutes int     `json:"session_duration_minutes"`
	TotalHours             float64 `json:"total_hours"`
	Price                  float64 `json:"price"`
	Status                 string  `json:"status"`
}

type CreateCourseOutput struct {
	Course *entities.Course `json:"course"`
}

type CreateCourseUseCase interface {
	Execute(ctx context.Context, input CreateCourseInput) (*CreateCourseOutput, error)
}

type createCourseUseCaseImpl struct {
	repo repointerface.CourseRepository
}

func NewCreateCourseUseCase(repo repointerface.CourseRepository) CreateCourseUseCase {
	return &createCourseUseCaseImpl{repo: repo}
}

func (uc *createCourseUseCaseImpl) Execute(ctx context.Context, input CreateCourseInput) (*CreateCourseOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	if input.Code == "" || input.Name == "" {
		return nil, errors.New("code and name are required")
	}
	if input.Status == "" {
		input.Status = "ACTIVE"
	}
	course := &entities.Course{
		Code:                   input.Code,
		Name:                   input.Name,
		Description:            input.Description,
		GradeLevel:             input.GradeLevel,
		Subject:                input.Subject,
		SessionCount:           input.SessionCount,
		SessionDurationMinutes: input.SessionDurationMinutes,
		TotalHours:             input.TotalHours,
		Price:                  input.Price,
		Status:                 input.Status,
	}
	created, err := uc.repo.Create(ctx, course)
	if err != nil {
		ctxLogger.Errorf("Error creating course: %v", err)
		return nil, err
	}
	return &CreateCourseOutput{Course: created}, nil
}
