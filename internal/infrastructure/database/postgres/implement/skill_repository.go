package implement

import (
	"context"
	"doan/internal/entities"
	"doan/internal/infrastructure/database/postgres"
	"doan/internal/repositories"
	repointerface "doan/internal/repositories/interface"
	skillservice "doan/internal/services/skills"
	"doan/pkg/base_struct"
	"doan/pkg/config"
	"doan/pkg/logger"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type skillRepository struct {
	base_struct.BaseDependency
	repositories.BaseRepository[entities.Skill]
	db *gorm.DB
}

func NewSkillRepository(
	db *gorm.DB,
	log logger.Logger,
	manager config.Manager,
) repointerface.SkillRepository {
	modelRepo := postgres.NewBaseRepository[entities.Skill](log, manager, db, "skills")
	return &skillRepository{
		BaseDependency: base_struct.BaseDependency{
			Log:           log,
			ConfigManager: manager,
		},
		BaseRepository: modelRepo,
		db:             db,
	}
}

func (r *skillRepository) ListCatalog(ctx context.Context, search string, limit int) ([]*entities.Skill, error) {
	if limit <= 0 {
		limit = 200
	}
	if limit > 500 {
		limit = 500
	}

	var items []*entities.Skill
	query := r.db.WithContext(ctx).
		Model(&entities.Skill{}).
		Where("status = ?", "ACTIVE")

	keyword := strings.TrimSpace(search)
	if keyword != "" {
		pattern := "%" + keyword + "%"
		query = query.Where("code ILIKE ? OR name ILIKE ?", pattern, pattern)
	}

	if err := query.Order("name ASC").Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *skillRepository) EnsureByCodes(ctx context.Context, codes []string) ([]*entities.Skill, error) {
	normalized := skillservice.NormalizeCodes(codes)
	if len(normalized) == 0 {
		return []*entities.Skill{}, nil
	}

	var existing []*entities.Skill
	if err := r.db.WithContext(ctx).
		Model(&entities.Skill{}).
		Where("code IN ?", normalized).
		Find(&existing).Error; err != nil {
		return nil, err
	}

	existingByCode := make(map[string]*entities.Skill, len(existing))
	for _, item := range existing {
		existingByCode[item.Code] = item
	}

	missing := make([]entities.Skill, 0)
	for _, code := range normalized {
		if _, ok := existingByCode[code]; ok {
			continue
		}
		missing = append(missing, entities.Skill{
			Code:   code,
			Name:   skillservice.HumanizeCode(code),
			Status: "ACTIVE",
		})
	}

	if len(missing) > 0 {
		for index := range missing {
			if err := r.db.WithContext(ctx).Create(&missing[index]).Error; err != nil {
				return nil, err
			}
			item := missing[index]
			existingByCode[item.Code] = &item
		}
	}

	result := make([]*entities.Skill, 0, len(normalized))
	for _, code := range normalized {
		if item, ok := existingByCode[code]; ok {
			result = append(result, item)
		}
	}

	return result, nil
}

func (r *skillRepository) SyncTeacherSkills(ctx context.Context, teacherID string, codes []string) error {
	skills, err := r.EnsureByCodes(ctx, codes)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM teacher_skills WHERE teacher_id = ?", teacherID).Error; err != nil {
			if isMissingRelationError(err) {
				r.Log.Warn(ctx, "teacher_skills table is missing; skipping teacher skill link sync until migrations run", "error", err)
				return nil
			}
			return err
		}

		now := time.Now()
		for _, item := range skills {
			if item == nil {
				continue
			}
			if err := tx.Exec(
				"INSERT INTO teacher_skills (teacher_id, skill_id, created_at) VALUES (?, ?, ?) ON CONFLICT DO NOTHING",
				teacherID,
				item.ID,
				now,
			).Error; err != nil {
				if isMissingRelationError(err) {
					r.Log.Warn(ctx, "teacher_skills table is missing; skipping teacher skill link sync until migrations run", "error", err)
					return nil
				}
				return err
			}
		}

		return nil
	})
}

func (r *skillRepository) SyncCourseRequiredSkills(ctx context.Context, courseID string, codes []string) error {
	skills, err := r.EnsureByCodes(ctx, codes)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM course_required_skills WHERE course_id = ?", courseID).Error; err != nil {
			if isMissingRelationError(err) {
				r.Log.Warn(ctx, "course_required_skills table is missing; skipping course skill link sync until migrations run", "error", err)
				return nil
			}
			return err
		}

		now := time.Now()
		for _, item := range skills {
			if item == nil {
				continue
			}
			if err := tx.Exec(
				"INSERT INTO course_required_skills (course_id, skill_id, created_at) VALUES (?, ?, ?) ON CONFLICT DO NOTHING",
				courseID,
				item.ID,
				now,
			).Error; err != nil {
				if isMissingRelationError(err) {
					r.Log.Warn(ctx, "course_required_skills table is missing; skipping course skill link sync until migrations run", "error", err)
					return nil
				}
				return err
			}
		}

		return nil
	})
}

func isMissingRelationError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(strings.ToLower(err.Error()), "does not exist") ||
		errors.Is(err, gorm.ErrInvalidDB)
}
