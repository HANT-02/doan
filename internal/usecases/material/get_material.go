package material

import (
	"context"

	repositoryinterface "doan/internal/repositories/interface"
)

type GetMaterialInput struct {
	ID string
}

type GetMaterialUseCase interface {
	Execute(ctx context.Context, input GetMaterialInput) (*MaterialView, error)
}

type getMaterialUseCase struct {
	materialRepo repositoryinterface.MaterialRepository
}

func NewGetMaterialUseCase(materialRepo repositoryinterface.MaterialRepository) GetMaterialUseCase {
	return &getMaterialUseCase{materialRepo: materialRepo}
}

func (uc *getMaterialUseCase) Execute(ctx context.Context, input GetMaterialInput) (*MaterialView, error) {
	materialEntity, err := uc.materialRepo.GetDetailed(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	view := mapMaterial(materialEntity)
	return &view, nil
}
