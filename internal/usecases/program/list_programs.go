package program

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type ListProgramsInput struct {
	Search string
	Track  string
	Page   int
	Limit  int
}

type PaginationMeta struct {
	ItemsPerPage uint64 `json:"items_per_page"`
	TotalItems   uint64 `json:"total_items"`
	CurrentPage  uint64 `json:"current_page"`
	TotalPages   uint64 `json:"total_pages"`
}

type ListProgramsOutput struct {
	Programs   []*entities.Program `json:"programs"`
	Pagination PaginationMeta      `json:"pagination"`
}

type ListProgramsUseCase interface {
	Execute(ctx context.Context, input ListProgramsInput) (*ListProgramsOutput, error)
}

type listProgramsUseCaseImpl struct {
	repo repointerface.ProgramRepository
}

func NewListProgramsUseCase(repo repointerface.ProgramRepository) ListProgramsUseCase {
	return &listProgramsUseCaseImpl{repo: repo}
}

func (uc *listProgramsUseCaseImpl) Execute(ctx context.Context, input ListProgramsInput) (*ListProgramsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	condition := repositories.NewCommonCondition().WithPaging(uint64(input.Limit), uint64(input.Page))
	if input.Track != "" {
		condition.WithCondition("track", input.Track, "eq")
	}

	res, err := uc.repo.GetByCondition(ctx, condition)
	if err != nil {
		ctxLogger.Errorf("Error listing programs: %v", err)
		return nil, err
	}

	out := &ListProgramsOutput{Programs: res.Data}
	out.Pagination.CurrentPage = uint64(res.Meta.CurrentPage)
	out.Pagination.ItemsPerPage = uint64(res.Meta.ItemsPerPage)
	out.Pagination.TotalItems = uint64(res.Meta.TotalItems)
	out.Pagination.TotalPages = uint64(res.Meta.TotalPages)

	return out, nil
}
