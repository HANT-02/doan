package predictive

import (
	"context"

	predictiveservice "doan/internal/services/predictive"
	"doan/pkg/logger"
)

type ListStudentPredictionsInput struct {
	Search     string
	OnlyAtRisk bool
	Page       int
	Limit      int
	Refresh    bool
}

type ListStudentPredictionsOutput = predictiveservice.ListStudentPredictionsOutput

type ListStudentPredictionsUseCase interface {
	Execute(ctx context.Context, input ListStudentPredictionsInput) (*ListStudentPredictionsOutput, error)
}

type listStudentPredictionsUseCase struct {
	atRiskService predictiveservice.AtRiskService
}

func NewListStudentPredictionsUseCase(atRiskService predictiveservice.AtRiskService) ListStudentPredictionsUseCase {
	return &listStudentPredictionsUseCase{atRiskService: atRiskService}
}

func (uc *listStudentPredictionsUseCase) Execute(ctx context.Context, input ListStudentPredictionsInput) (*ListStudentPredictionsOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	output, err := uc.atRiskService.ListStudentPredictions(ctx, predictiveservice.ListStudentPredictionsInput{
		Search:     input.Search,
		OnlyAtRisk: input.OnlyAtRisk,
		Page:       input.Page,
		Limit:      input.Limit,
		Refresh:    input.Refresh,
	})
	if err != nil {
		ctxLogger.Errorf("Failed to list at-risk predictions: %v", err)
		return nil, err
	}
	return output, nil
}
