package scheduling

import (
	"context"
	"errors"

	schedulingstore "doan/internal/services/scheduling"
)

type GetPreviewInput struct {
	RunID string
}

type GetPreviewUseCase interface {
	Execute(ctx context.Context, input GetPreviewInput) (*PreviewResult, error)
	GetLatest(ctx context.Context) (*PreviewResult, error)
}

type getPreviewUseCase struct {
	store schedulingstore.PreviewStore[PreviewResult]
}

func NewGetPreviewUseCase(store schedulingstore.PreviewStore[PreviewResult]) GetPreviewUseCase {
	return &getPreviewUseCase{store: store}
}

func (uc *getPreviewUseCase) Execute(_ context.Context, input GetPreviewInput) (*PreviewResult, error) {
	if input.RunID == "" {
		return nil, errors.New("run_id is required")
	}

	result, ok := uc.store.Get(input.RunID)
	if !ok {
		return nil, errors.New("preview run not found")
	}
	return &result, nil
}

func (uc *getPreviewUseCase) GetLatest(_ context.Context) (*PreviewResult, error) {
	result, ok := uc.store.GetLatest()
	if !ok {
		return nil, errors.New("preview run not found")
	}
	return &result, nil
}
