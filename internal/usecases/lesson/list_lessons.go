package lesson

import (
	"context"
	"strings"
	"time"

	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type ListLessonsInput struct {
	ClassID   string
	TeacherID string
	Status    string
	DateFrom  *time.Time
	DateTo    *time.Time
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

type ListLessonsOutput struct {
	Lessons    []entities.Lesson
	Pagination struct {
		CurrentPage  int
		ItemsPerPage int
		TotalItems   int64
		TotalPages   int
	}
}

type ListLessonsUseCase interface {
	Execute(ctx context.Context, input ListLessonsInput) (*ListLessonsOutput, error)
}

type listLessonsUseCase struct {
	lessonRepo repointerface.LessonRepository
}

func NewListLessonsUseCase(lessonRepo repointerface.LessonRepository) ListLessonsUseCase {
	return &listLessonsUseCase{lessonRepo: lessonRepo}
}

func (uc *listLessonsUseCase) Execute(ctx context.Context, input ListLessonsInput) (*ListLessonsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	commonCond := repositories.NewCommonCondition()
	commonCond.SetPreload([]string{"Class", "Teacher", "Room"})

	if input.ClassID != "" {
		commonCond.AddCondition("class_id", input.ClassID, repositories.Equal)
	}
	if input.TeacherID != "" {
		commonCond.AddCondition("teacher_id", input.TeacherID, repositories.Equal)
	}
	if input.Status != "" {
		statuses := splitLessonStatuses(input.Status)
		if len(statuses) == 1 {
			commonCond.AddCondition("status", statuses[0], repositories.Equal)
		} else if len(statuses) > 1 {
			commonCond.AddCondition("status", statuses, repositories.In)
		}
	}
	if input.DateFrom != nil {
		commonCond.AddCondition("date_start", *input.DateFrom, repositories.GreaterThanOrEqual)
	}
	if input.DateTo != nil {
		commonCond.AddCondition("date_end", *input.DateTo, repositories.LessThanOrEqual)
	}

	if input.Page > 0 && input.Limit > 0 {
		commonCond.SetPaging(uint64(input.Limit), uint64(input.Page))
	}

	sortField := "date_start"
	switch strings.TrimSpace(input.SortBy) {
	case "date_start", "date_end", "created_at", "updated_at":
		sortField = input.SortBy
	}

	sortOrder := repositories.Asc
	if strings.EqualFold(strings.TrimSpace(input.SortOrder), repositories.Desc) {
		sortOrder = repositories.Desc
	}
	commonCond.AddSorting(sortField, sortOrder)

	result, err := uc.lessonRepo.GetByCondition(ctx, commonCond)
	if err != nil {
		ctxLogger.Errorf("Failed to list lessons: %v", err)
		return nil, err
	}

	var lessons []entities.Lesson
	total := int64(0)
	totalPages := 0
	if result != nil {
		for _, ptr := range result.Data {
			lessons = append(lessons, *ptr)
		}
		total = int64(result.Meta.TotalItems)
		totalPages = int(result.Meta.TotalPages)
	}

	return &ListLessonsOutput{
		Lessons: lessons,
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

func splitLessonStatuses(raw string) []string {
	parts := strings.Split(raw, ",")
	statuses := make([]string, 0, len(parts))
	for _, part := range parts {
		normalized := strings.ToUpper(strings.TrimSpace(part))
		if normalized == "" {
			continue
		}
		statuses = append(statuses, normalized)
	}
	return statuses
}
