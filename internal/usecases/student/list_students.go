package student

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type ListStudentsInput struct {
	Search    string
	Status    string
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

type ListStudentsOutput struct {
	Students   []entities.Student
	Pagination struct {
		CurrentPage  int
		ItemsPerPage int
		TotalItems   int64
		TotalPages   int
	}
}

type ListStudentsUseCase interface {
	Execute(ctx context.Context, input ListStudentsInput) (*ListStudentsOutput, error)
}

type listStudentsUseCase struct {
	studentRepo repointerface.StudentRepository
}

func NewListStudentsUseCase(studentRepo repointerface.StudentRepository) ListStudentsUseCase {
	return &listStudentsUseCase{
		studentRepo: studentRepo,
	}
}

func (uc *listStudentsUseCase) Execute(ctx context.Context, input ListStudentsInput) (*ListStudentsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	commonCond := repositories.NewCommonCondition()

	if input.Search != "" {
		commonCond.AddCondition("full_name ILIKE ? OR code ILIKE ?", "%"+input.Search+"%", repositories.Like)
	}

	if input.Status != "" {
		commonCond.AddCondition("status", input.Status, repositories.Equal)
	}

	if input.Page > 0 && input.Limit > 0 {
		commonCond.SetPaging(uint64(input.Limit), uint64(input.Page))
	}

	orderBy := "created_at"
	if input.SortBy != "" {
		order := repositories.Asc
		if input.SortOrder == "desc" || input.SortOrder == "DESC" {
			order = repositories.Desc
		}
		orderBy = input.SortBy + " " + order
	}
	commonCond.AddSorting(orderBy, repositories.Asc)

	result, err := uc.studentRepo.GetByCondition(ctx, commonCond)
	if err != nil {
		ctxLogger.Errorf("Failed to list students: %v", err)
		return nil, err
	}

	var students []entities.Student
	total := int64(0)
	totalPages := 0
	if result != nil {
		for _, ptr := range result.Data {
			students = append(students, *ptr)
		}
		total = int64(result.Meta.TotalItems)
		totalPages = int(result.Meta.TotalPages)
	}

	return &ListStudentsOutput{
		Students: students,
		Pagination: struct {
			CurrentPage  int
			ItemsPerPage int
			TotalItems   int64
			TotalPages   int
		}{
			CurrentPage:  input.Page,
			ItemsPerPage: input.Limit,
			TotalItems:   total,
			TotalPages:   totalPages,
		},
	}, nil
}
