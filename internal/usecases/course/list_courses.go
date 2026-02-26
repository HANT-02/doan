package course

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type ListCoursesInput struct {
	Search    string
	Status    string
	Subject   string
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

type PaginationMeta struct {
	ItemsPerPage uint64 `json:"items_per_page"`
	TotalItems   uint64 `json:"total_items"`
	CurrentPage  uint64 `json:"current_page"`
	TotalPages   uint64 `json:"total_pages"`
}

type ListCoursesOutput struct {
	Courses    []*entities.Course `json:"courses"`
	Pagination PaginationMeta     `json:"pagination"`
}

type ListCoursesUseCase interface {
	Execute(ctx context.Context, input ListCoursesInput) (*ListCoursesOutput, error)
}

type listCoursesUseCaseImpl struct {
	repo repointerface.CourseRepository
}

func NewListCoursesUseCase(repo repointerface.CourseRepository) ListCoursesUseCase {
	return &listCoursesUseCaseImpl{repo: repo}
}

func (uc *listCoursesUseCaseImpl) Execute(ctx context.Context, input ListCoursesInput) (*ListCoursesOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	condition := repositories.NewCommonCondition().WithPaging(uint64(input.Limit), uint64(input.Page))
	if input.Status != "" {
		condition.WithCondition("status", input.Status, "eq")
	}
	if input.Subject != "" {
		condition.WithCondition("subject", input.Subject, "eq")
	}

	res, err := uc.repo.GetByCondition(ctx, condition)
	if err != nil {
		ctxLogger.Errorf("Error listing courses: %v", err)
		return nil, err
	}

	out := &ListCoursesOutput{
		Courses: res.Data,
	}
	out.Pagination.CurrentPage = uint64(res.Meta.CurrentPage)
	out.Pagination.ItemsPerPage = uint64(res.Meta.ItemsPerPage)
	out.Pagination.TotalItems = uint64(res.Meta.TotalItems)
	out.Pagination.TotalPages = uint64(res.Meta.TotalPages)

	return out, nil
}
