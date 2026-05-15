package teacher

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	skillservice "doan/internal/services/skills"
	"doan/pkg/logger"
	"strings"
)

// ListTeachersInput represents the input for listing teachers
type ListTeachersInput struct {
	Search         string `json:"search"`
	Status         string `json:"status"`
	EmploymentType string `json:"employment_type"`
	CourseID       string `json:"course_id"`
	ClassID        string `json:"class_id"`
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
	courseRepo  repointerface.CourseRepository
	classRepo   repointerface.ClassRepository
}

// NewListTeachersUseCase creates a new instance of ListTeachersUseCase
func NewListTeachersUseCase(
	teacherRepo repointerface.TeacherRepository,
	courseRepo repointerface.CourseRepository,
	classRepo repointerface.ClassRepository,
) ListTeachersUseCase {
	return &listTeachersUseCase{
		teacherRepo: teacherRepo,
		courseRepo:  courseRepo,
		classRepo:   classRepo,
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

	requiredSkills, err := uc.resolveRequiredSkills(ctx, input.CourseID, input.ClassID)
	if err != nil {
		ctxLogger.Errorf("Failed to resolve required skills: %v", err)
		return nil, err
	}

	if len(requiredSkills) > 0 {
		teachers, err := uc.teacherRepo.ListByRequiredSkillCodes(ctx, requiredSkills)
		if err != nil {
			ctxLogger.Errorf("Failed to list teachers by required skills: %v", err)
			return nil, err
		}

		filtered := uc.filterTeachersInMemory(teachers, input)
		start := (input.Page - 1) * input.Limit
		if start > len(filtered) {
			start = len(filtered)
		}
		end := start + input.Limit
		if end > len(filtered) {
			end = len(filtered)
		}

		meta := repositories.NewMeta(&repositories.Paging{
			Limit: uint64(input.Limit),
			Page:  uint64(input.Page),
		}, uint64(len(filtered)))

		return &ListTeachersOutput{
			Teachers:   filtered[start:end],
			Pagination: &meta,
		}, nil
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

func (uc *listTeachersUseCase) resolveRequiredSkills(ctx context.Context, courseID string, classID string) ([]string, error) {
	if strings.TrimSpace(classID) != "" {
		classEntity, err := uc.classRepo.GetByID(ctx, classID)
		if err != nil {
			return nil, err
		}
		if classEntity.CourseID != nil {
			courseID = *classEntity.CourseID
		}
	}

	if strings.TrimSpace(courseID) == "" {
		return []string{}, nil
	}

	courseEntity, err := uc.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return skillservice.NormalizeCodes([]string(courseEntity.RequiredSkills)), nil
}

func (uc *listTeachersUseCase) filterTeachersInMemory(teachers []*entities.Teacher, input ListTeachersInput) []*entities.Teacher {
	keyword := strings.TrimSpace(strings.ToLower(input.Search))
	filtered := make([]*entities.Teacher, 0, len(teachers))

	for _, teacherEntity := range teachers {
		if teacherEntity == nil {
			continue
		}
		if input.Status != "" && teacherEntity.Status != input.Status {
			continue
		}
		if input.EmploymentType != "" && teacherEntity.EmploymentType != input.EmploymentType {
			continue
		}
		if keyword != "" {
			content := strings.ToLower(strings.Join([]string{
				teacherEntity.FullName,
				teacherEntity.Email,
				teacherEntity.Phone,
				teacherEntity.Code,
			}, " "))
			if !strings.Contains(content, keyword) {
				continue
			}
		}
		filtered = append(filtered, teacherEntity)
	}

	return filtered
}
