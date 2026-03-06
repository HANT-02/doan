package material

import (
	"context"
	"errors"

	"doan/internal/entities"
	repositoryinterface "doan/internal/repositories/interface"
)

type ReviewMaterialInput struct {
	MaterialID          string
	ComplianceOfficerID string
	Approved            bool
	RejectReason        string
	Notes               string
}

type ReviewMaterialUseCase interface {
	Execute(ctx context.Context, input ReviewMaterialInput) (*MaterialView, error)
}

type reviewMaterialUseCase struct {
	materialRepo         repositoryinterface.MaterialRepository
	approvalDecisionRepo repositoryinterface.ApprovalDecisionRepository
}

func NewReviewMaterialUseCase(
	materialRepo repositoryinterface.MaterialRepository,
	approvalDecisionRepo repositoryinterface.ApprovalDecisionRepository,
) ReviewMaterialUseCase {
	return &reviewMaterialUseCase{
		materialRepo:         materialRepo,
		approvalDecisionRepo: approvalDecisionRepo,
	}
}

func (uc *reviewMaterialUseCase) Execute(ctx context.Context, input ReviewMaterialInput) (*MaterialView, error) {
	if input.MaterialID == "" || input.ComplianceOfficerID == "" {
		return nil, errors.New("material_id and compliance_officer_id are required")
	}

	materialEntity, err := uc.materialRepo.GetDetailed(ctx, input.MaterialID)
	if err != nil {
		return nil, err
	}

	var auditLogID *string
	if len(materialEntity.AuditLogs) > 0 {
		auditLogID = &materialEntity.AuditLogs[0].ID
	}

	_, err = uc.approvalDecisionRepo.Create(ctx, &entities.ApprovalDecision{
		MaterialID:          input.MaterialID,
		AuditLogID:          auditLogID,
		ComplianceOfficerID: input.ComplianceOfficerID,
		Approved:            input.Approved,
		RejectReason:        input.RejectReason,
		Notes:               input.Notes,
	})
	if err != nil {
		return nil, err
	}

	status := "APPROVED"
	if !input.Approved {
		status = "REJECTED"
	}

	if err := uc.materialRepo.Update(ctx, input.MaterialID, map[string]interface{}{
		"status": status,
	}); err != nil {
		return nil, err
	}

	updatedMaterial, err := uc.materialRepo.GetDetailed(ctx, input.MaterialID)
	if err != nil {
		return nil, err
	}

	view := mapMaterial(updatedMaterial)
	return &view, nil
}
