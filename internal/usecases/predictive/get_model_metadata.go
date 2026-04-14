package predictive

import (
	"context"

	predictiveservice "doan/internal/services/predictive"
	"doan/pkg/logger"
)

type GetModelMetadataOutput struct {
	ModelMetadata predictiveservice.AtRiskModelMetadata `json:"model_metadata"`
}

type GetModelMetadataUseCase interface {
	Execute(ctx context.Context, refresh bool) (*GetModelMetadataOutput, error)
}

type getModelMetadataUseCase struct {
	atRiskService predictiveservice.AtRiskService
}

func NewGetModelMetadataUseCase(atRiskService predictiveservice.AtRiskService) GetModelMetadataUseCase {
	return &getModelMetadataUseCase{atRiskService: atRiskService}
}

func (uc *getModelMetadataUseCase) Execute(ctx context.Context, refresh bool) (*GetModelMetadataOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	metadata, err := uc.atRiskService.GetModelMetadata(ctx, refresh)
	if err != nil {
		ctxLogger.Errorf("Failed to get predictive model metadata: %v", err)
		return nil, err
	}

	return &GetModelMetadataOutput{ModelMetadata: *metadata}, nil
}
