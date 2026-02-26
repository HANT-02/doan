package program

import (
	"context"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type DeleteProgramInput struct {
	ID string
}

type DeleteProgramOutput struct {
	Message string `json:"message"`
}

type DeleteProgramUseCase interface {
	Execute(ctx context.Context, input DeleteProgramInput) (*DeleteProgramOutput, error)
}

type deleteProgramUseCaseImpl struct {
	repo repointerface.ProgramRepository
}

func NewDeleteProgramUseCase(repo repointerface.ProgramRepository) DeleteProgramUseCase {
	return &deleteProgramUseCaseImpl{repo: repo}
}

func (uc *deleteProgramUseCaseImpl) Execute(ctx context.Context, input DeleteProgramInput) (*DeleteProgramOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	err := uc.repo.SoftDelete(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Error soft deleting program: %v", err)
		return nil, err
	}
	return &DeleteProgramOutput{Message: "Program deleted successfully"}, nil
}
