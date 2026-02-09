package teacher

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

// ListTeachersInput represents the input for listing teachers
type ListTeachersInput struct {
	Search         string `json:"search"`
	Status         string `json:"status"`
	EmploymentType string `json:"employment_type"`
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
	SortBy         string `json:"sort_by"`
	SortOrder      string `json:"sort_order"`
}

// ListTeachersOutput represents the output after listing teachers
type ListTeachersOutput struct {
	Teachers   []*entities.Teacher `json:"teachers"`
	Pagination *repositories.Meta  `json:"pagination"`
}

// ListTeachersUseCase defines the interface for listing teachers
type ListTeachersUseCase interface {
	Execute(ctx context.Context, input ListTeachersInput) (*ListTeachersOutput, error)
}

type listTeachersUseCase struct {
	teacherRepo repointerface.TeacherRepository
}

// NewListTeachersUseCase creates a new instance of ListTeachersUseCase
func NewListTeachersUseCase(teacherRepo repointerface.TeacherRepository) ListTeachersUseCase {
	return &listTeachersUseCase{
		teacherRepo: teacherRepo,
	}
}

func (uc *listTeachersUseCase) Execute(ctx context.Context, input ListTeachersInput) (*ListTeachersOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	// Set default pagination
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100 // Max limit
	}

	// Build condition
	condition := repositories.NewCommonCondition()
	condition.SetPaging(uint64(input.Limit), uint64(input.Page))

	// Add search filter (search in multiple fields using OR)
	if input.Search != "" {
		orConditions := []repositories.Condition{
			{Field: "full_name", Value: "%" + input.Search + "%", Op: repositories.Like},
			{Field: "email", Value: "%" + input.Search + "%", Op: repositories.Like},
			{Field: "phone", Value: "%" + input.Search + "%", Op: repositories.Like},
			{Field: "code", Value: "%" + input.Search + "%", Op: repositories.Like},
		}
		condition.AddOrCondition(orConditions)
	}

	// Add status filter
	if input.Status != "" {
		condition.AddCondition("status", input.Status, repositories.Equal)
	}

	// Add employment type filter
	if input.EmploymentType != "" {
		condition.AddCondition("employment_type", input.EmploymentType, repositories.Equal)
	}

	// Add sorting
	if input.SortBy != "" {
		order := repositories.Asc
		if input.SortOrder == repositories.Desc {
			order = repositories.Desc
		}
		condition.AddSorting(input.SortBy, order)
	} else {
		// Default sort by created_at DESC
		condition.AddSorting("created_at", repositories.Desc)
	}

	// Get teachers
	pagination, err := uc.teacherRepo.GetByCondition(ctx, condition)
	if err != nil {
		ctxLogger.Errorf("Failed to list teachers: %v", err)
		return nil, err
	}

	return &ListTeachersOutput{
		Teachers:   pagination.Data,
		Pagination: &pagination.Meta,
	}, nil
}
