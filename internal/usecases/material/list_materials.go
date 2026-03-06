package material

import (
	"context"

	repositoryinterface "doan/internal/repositories/interface"
)

type ListMaterialsInput struct {
	TeacherID string
	Status    string
	Queue     string
}

type ListMaterialsOutput struct {
	Materials []MaterialView `json:"materials"`
}

type ListMaterialsUseCase interface {
	Execute(ctx context.Context, input ListMaterialsInput) (*ListMaterialsOutput, error)
}

type listMaterialsUseCase struct {
	materialRepo repositoryinterface.MaterialRepository
}

func NewListMaterialsUseCase(materialRepo repositoryinterface.MaterialRepository) ListMaterialsUseCase {
	return &listMaterialsUseCase{materialRepo: materialRepo}
}

func (uc *listMaterialsUseCase) Execute(ctx context.Context, input ListMaterialsInput) (*ListMaterialsOutput, error) {
	materials, err := uc.materialRepo.ListDetailed(ctx, repositoryinterface.MaterialFilter{
		TeacherID: input.TeacherID,
		Status:    input.Status,
		Queue:     input.Queue,
	})
	if err != nil {
		return nil, err
	}

	views := make([]MaterialView, 0, len(materials))
	for _, item := range materials {
		materialCopy := item
		views = append(views, mapMaterial(&materialCopy))
	}

	return &ListMaterialsOutput{Materials: views}, nil
}
