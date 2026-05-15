package teacher

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
)

type ListSkillCatalogInput struct {
	Search string `json:"search"`
	Limit  int    `json:"limit"`
}

type ListSkillCatalogOutput struct {
	Skills []*entities.Skill `json:"skills"`
}

type ListSkillCatalogUseCase interface {
	Execute(ctx context.Context, input ListSkillCatalogInput) (*ListSkillCatalogOutput, error)
}

type listSkillCatalogUseCase struct {
	skillRepo repointerface.SkillRepository
}

func NewListSkillCatalogUseCase(skillRepo repointerface.SkillRepository) ListSkillCatalogUseCase {
	return &listSkillCatalogUseCase{skillRepo: skillRepo}
}

func (uc *listSkillCatalogUseCase) Execute(ctx context.Context, input ListSkillCatalogInput) (*ListSkillCatalogOutput, error) {
	skills, err := uc.skillRepo.ListCatalog(ctx, input.Search, input.Limit)
	if err != nil {
		return nil, err
	}

	return &ListSkillCatalogOutput{Skills: skills}, nil
}
