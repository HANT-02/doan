package repositoryinterface

import (
	"context"
	"doan/internal/entities"
	"doan/internal/repositories"
)

type SkillRepository interface {
	repositories.BaseRepository[entities.Skill]
	ListCatalog(ctx context.Context, search string, limit int) ([]*entities.Skill, error)
	EnsureByCodes(ctx context.Context, codes []string) ([]*entities.Skill, error)
	SyncTeacherSkills(ctx context.Context, teacherID string, codes []string) error
	SyncCourseRequiredSkills(ctx context.Context, courseID string, codes []string) error
}
