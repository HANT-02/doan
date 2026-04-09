package shift

import (
	"context"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type ListShiftsInput struct {
	Search      string
	SessionType string
	IsActive    *bool
	Page        int
	Limit       int
	SortBy      string
	SortOrder   string
}

type ListShiftsOutput struct {
	Shifts     []entities.Shift
	Pagination struct {
		CurrentPage  int
		ItemsPerPage int
		TotalItems   int64
		TotalPages   int
	}
}

type ListShiftsUseCase interface {
	Execute(ctx context.Context, input ListShiftsInput) (*ListShiftsOutput, error)
}

type listShiftsUseCase struct {
	shiftRepo repointerface.ShiftRepository
}

func NewListShiftsUseCase(shiftRepo repointerface.ShiftRepository) ListShiftsUseCase {
	return &listShiftsUseCase{shiftRepo: shiftRepo}
}

func (uc *listShiftsUseCase) Execute(ctx context.Context, input ListShiftsInput) (*ListShiftsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	commonCond := repositories.NewCommonCondition()
	if input.Search != "" {
		commonCond.AddCondition("name ILIKE ? OR code ILIKE ?", "%"+input.Search+"%", repositories.Like)
	}
	if input.SessionType != "" {
		commonCond.AddCondition("session_type", input.SessionType, repositories.Equal)
	}
	if input.IsActive != nil {
		commonCond.AddCondition("is_active", *input.IsActive, repositories.Equal)
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
	commonCond.AddSorting(orderBy, "")

	result, err := uc.shiftRepo.GetByCondition(ctx, commonCond)
	if err != nil {
		ctxLogger.Errorf("Failed to list shifts: %v", err)
		return nil, err
	}

	var shifts []entities.Shift
	total := int64(0)
	totalPages := 0
	if result != nil {
		for _, ptr := range result.Data {
			shifts = append(shifts, *ptr)
		}
		total = int64(result.Meta.TotalItems)
		totalPages = int(result.Meta.TotalPages)
	}

	return &ListShiftsOutput{
		Shifts: shifts,
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
