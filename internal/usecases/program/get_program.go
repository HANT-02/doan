package program

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
)

type GetProgramInput struct {
	ID string
}

type GetProgramOutput struct {
	Program *entities.Program `json:"program"`
}

type GetProgramUseCase interface {
	Execute(ctx context.Context, input GetProgramInput) (*GetProgramOutput, error)
}

type getProgramUseCaseImpl struct {
	repo repointerface.ProgramRepository
}

func NewGetProgramUseCase(repo repointerface.ProgramRepository) GetProgramUseCase {
	return &getProgramUseCaseImpl{repo: repo}
}

func (uc *getProgramUseCaseImpl) Execute(ctx context.Context, input GetProgramInput) (*GetProgramOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	prog, err := uc.repo.GetProgramWithCourses(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Error getting program by ID: %v", err)
		return nil, err
	}
	return &GetProgramOutput{Program: prog}, nil
}
