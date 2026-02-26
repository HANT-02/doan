package program

import (
	"context"
	"doan/internal/entities"
	repointerface "doan/internal/repositories/interface"
	"doan/pkg/logger"
	"time"
)

type UpdateProgramInput struct {
	ID            string     `json:"id"`
	Code          *string    `json:"code"`
	Name          *string    `json:"name"`
	Track         *string    `json:"track"`
	EffectiveFrom *time.Time `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`
	ApprovalNote  *string    `json:"approval_note"`
}

type UpdateProgramOutput struct {
	Program *entities.Program `json:"program"`
}

type UpdateProgramUseCase interface {
	Execute(ctx context.Context, input UpdateProgramInput) (*UpdateProgramOutput, error)
}

type updateProgramUseCaseImpl struct {
	repo repointerface.ProgramRepository
}

func NewUpdateProgramUseCase(repo repointerface.ProgramRepository) UpdateProgramUseCase {
	return &updateProgramUseCaseImpl{repo: repo}
}

func (uc *updateProgramUseCaseImpl) Execute(ctx context.Context, input UpdateProgramInput) (*UpdateProgramOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	prog, err := uc.repo.GetByID(ctx, input.ID)
	if err != nil {
		ctxLogger.Errorf("Error getting program by ID for update: %v", err)
		return nil, err
	}
	updateData := map[string]interface{}{}

	if input.Code != nil {
		prog.Code = *input.Code
		updateData["code"] = *input.Code
	}
	if input.Name != nil {
		prog.Name = *input.Name
		updateData["name"] = *input.Name
	}
	if input.Track != nil {
		prog.Track = *input.Track
		updateData["track"] = *input.Track
	}
	if input.EffectiveFrom != nil {
		prog.EffectiveFrom = input.EffectiveFrom
		updateData["effective_from"] = input.EffectiveFrom
	}
	if input.EffectiveTo != nil {
		prog.EffectiveTo = input.EffectiveTo
		updateData["effective_to"] = input.EffectiveTo
	}
	if input.ApprovalNote != nil {
		prog.ApprovalNote = *input.ApprovalNote
		updateData["approval_note"] = *input.ApprovalNote
	}

	err = uc.repo.Update(ctx, prog.ID, updateData)
	if err != nil {
		ctxLogger.Errorf("Error updating program: %v", err)
		return nil, err
	}
	return &UpdateProgramOutput{Program: prog}, nil
}
