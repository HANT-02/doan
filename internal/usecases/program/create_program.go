package program

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"errors"
	"time"
)

type CreateProgramInput struct {
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Track         string     `json:"track"`
	EffectiveFrom *time.Time `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`
	ApprovalNote  string     `json:"approval_note"`
}

type CreateProgramOutput struct {
	Program *entities.Program `json:"program"`
}

type CreateProgramUseCase interface {
	Execute(ctx context.Context, input CreateProgramInput) (*CreateProgramOutput, error)
}

type createProgramUseCaseImpl struct {
	repo repointerface.ProgramRepository
}

func NewCreateProgramUseCase(repo repointerface.ProgramRepository) CreateProgramUseCase {
	return &createProgramUseCaseImpl{repo: repo}
}

func (uc *createProgramUseCaseImpl) Execute(ctx context.Context, input CreateProgramInput) (*CreateProgramOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	if input.Code == "" || input.Name == "" {
		return nil, errors.New("code and name are required")
	}
	prog := &entities.Program{
		Code:          input.Code,
		Name:          input.Name,
		Track:         input.Track,
		EffectiveFrom: input.EffectiveFrom,
		EffectiveTo:   input.EffectiveTo,
		ApprovalNote:  input.ApprovalNote,
	}
	created, err := uc.repo.Create(ctx, prog)
	if err != nil {
		ctxLogger.Errorf("Error creating program: %v", err)
		return nil, err
	}
	return &CreateProgramOutput{Program: created}, nil
}
