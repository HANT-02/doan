package account

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
	repositoryinterface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type ListUsersInput struct {
	Search    string
	Role      string
	IsActive  *bool
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

type ListUsersOutput struct {
	Users      []entities.User `json:"users"`
	Pagination struct {
		CurrentPage  int   `json:"current_page"`
		ItemsPerPage int   `json:"items_per_page"`
		TotalItems   int64 `json:"total_items"`
		TotalPages   int   `json:"total_pages"`
	} `json:"pagination"`
}

type ListUsersUseCase interface {
	Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error)
}

type listUsersUseCase struct {
	userRepo repositoryinterface.UserRepository
}

func NewListUsersUseCase(userRepo repositoryinterface.UserRepository) ListUsersUseCase {
	return &listUsersUseCase{userRepo: userRepo}
}

func (uc *listUsersUseCase) Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	commonCond := repositories.NewCommonCondition()
	if input.Search != "" {
		commonCond.AddOrCondition([]repositories.Condition{
			{Field: "full_name", Value: input.Search, Op: repositories.ILikeContains},
			{Field: "email", Value: input.Search, Op: repositories.ILikeContains},
			{Field: "code", Value: input.Search, Op: repositories.ILikeContains},
		})
	}
	if input.Role != "" {
		commonCond.AddCondition("role", input.Role, repositories.Equal)
	}
	if input.IsActive != nil {
		commonCond.AddCondition("is_active", *input.IsActive, repositories.Equal)
	}
	if input.Page > 0 && input.Limit > 0 {
		commonCond.SetPaging(uint64(input.Limit), uint64(input.Page))
	}

	sortField := "created_at"
	sortOrder := repositories.Desc
	if input.SortBy != "" {
		sortField = input.SortBy
		if input.SortOrder == "asc" || input.SortOrder == "ASC" {
			sortOrder = repositories.Asc
		}
	}
	commonCond.AddSorting(sortField, sortOrder)

	result, err := uc.userRepo.GetByCondition(ctx, commonCond)
	if err != nil {
		ctxLogger.Errorf("Failed to list users: %v", err)
		return nil, err
	}

	out := &ListUsersOutput{}
	if result != nil {
		for _, ptr := range result.Data {
			out.Users = append(out.Users, *ptr)
		}
		out.Pagination.TotalItems = int64(result.Meta.TotalItems)
		out.Pagination.TotalPages = int(result.Meta.TotalPages)
	}
	out.Pagination.CurrentPage = input.Page
	out.Pagination.ItemsPerPage = input.Limit

	return out, nil
}
