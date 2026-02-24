package class

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type ListClassesInput struct {
	Search    string
	Status    string
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

type ListClassesOutput struct {
	Classes    []entities.Class
	Pagination struct {
		CurrentPage  int
		ItemsPerPage int
		TotalItems   int64
		TotalPages   int
	}
}

type ListClassesUseCase interface {
	Execute(ctx context.Context, input ListClassesInput) (*ListClassesOutput, error)
}

type listClassesUseCase struct {
	classRepo repointerface.ClassRepository
}

func NewListClassesUseCase(classRepo repointerface.ClassRepository) ListClassesUseCase {
	return &listClassesUseCase{
		classRepo: classRepo,
	}
}

func (uc *listClassesUseCase) Execute(ctx context.Context, input ListClassesInput) (*ListClassesOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	commonCond := repositories.NewCommonCondition()

	if input.Search != "" {
		commonCond.AddCondition("name ILIKE ?", "%"+input.Search+"%", repositories.Like)
	}

	if input.Status != "" {
		commonCond.AddCondition("status", input.Status, repositories.Equal)
	}

	if input.Page > 0 && input.Limit > 0 {
		commonCond.SetPaging(uint64(input.Limit), uint64(input.Page))
	}

	orderBy := "created_at DESC"
	if input.SortBy != "" {
		order := repositories.Asc
		if input.SortOrder == "desc" || input.SortOrder == "DESC" {
			order = repositories.Desc
		}
		orderBy = input.SortBy + " " + order
	}
	commonCond.AddSorting(orderBy, "")

	result, err := uc.classRepo.GetByCondition(ctx, commonCond)
	if err != nil {
		ctxLogger.Errorf("Failed to list classes: %v", err)
		return nil, err
	}

	var classes []entities.Class
	total := int64(0)
	totalPages := 0
	if result != nil {
		for _, ptr := range result.Data {
			classes = append(classes, *ptr)
		}
		total = int64(result.Meta.TotalItems)
		totalPages = int(result.Meta.TotalPages)
	}

	return &ListClassesOutput{
		Classes: classes,
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
