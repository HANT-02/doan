package material

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"

	"doan/internal/entities"
	repositoryinterface "doan/internal/repositories/interface"
	auditservice "doan/internal/services/audit"
	"doan/pkg/logger"
)

type UploadMaterialInput struct {
	TeacherID   string
	Title       string
	Description string
	FileName    string
	FileType    string
	Content     []byte
}

type UploadMaterialUseCase interface {
	Execute(ctx context.Context, input UploadMaterialInput) (*MaterialView, error)
}

type uploadMaterialUseCase struct {
	materialRepo repositoryinterface.MaterialRepository
	auditLogRepo repositoryinterface.AuditLogRepository
	labelRepo    repositoryinterface.LabelRepository
	storage      auditservice.StorageService
	ocrService   auditservice.OCRService
	gemini       auditservice.GeminiService
}

func NewUploadMaterialUseCase(
	materialRepo repositoryinterface.MaterialRepository,
	auditLogRepo repositoryinterface.AuditLogRepository,
	labelRepo repositoryinterface.LabelRepository,
	storage auditservice.StorageService,
	ocrService auditservice.OCRService,
	gemini auditservice.GeminiService,
) UploadMaterialUseCase {
	return &uploadMaterialUseCase{
		materialRepo: materialRepo,
		auditLogRepo: auditLogRepo,
		labelRepo:    labelRepo,
		storage:      storage,
		ocrService:   ocrService,
		gemini:       gemini,
	}
}

func (uc *uploadMaterialUseCase) Execute(ctx context.Context, input UploadMaterialInput) (*MaterialView, error) {
	ctxLogger := logger.NewLogger(ctx)

	if input.TeacherID == "" || input.FileName == "" || len(input.Content) == 0 {
		return nil, errors.New("teacher_id and file are required")
	}

	filePath, err := uc.storage.Save(ctx, input.FileName, input.Content)
	if err != nil {
		ctxLogger.Errorf("Failed to store material file: %v", err)
		return nil, err
	}

	if input.Title == "" {
		input.Title = filepath.Base(input.FileName)
	}

	materialEntity, err := uc.materialRepo.Create(ctx, &entities.Material{
		TeacherID:   input.TeacherID,
		Title:       input.Title,
		Description: input.Description,
		FileName:    input.FileName,
		FilePath:    filePath,
		FileType:    input.FileType,
		Status:      "SCANNING",
	})
	if err != nil {
		ctxLogger.Errorf("Failed to create material: %v", err)
		return nil, err
	}

	rawText, err := uc.ocrService.ExtractText(ctx, input.FileName, input.Content)
	if err != nil {
		ctxLogger.Errorf("Failed to extract OCR text: %v", err)
		return nil, err
	}

	inference, err := uc.gemini.Analyze(ctx, rawText)
	if err != nil {
		ctxLogger.Errorf("Failed to analyze material via Gemini stub: %v", err)
		return nil, err
	}

	labelEntity, err := uc.labelRepo.GetByCode(ctx, inference.LabelCode)
	if err != nil {
		ctxLogger.Errorf("Failed to resolve label %s: %v", inference.LabelCode, err)
		return nil, err
	}

	detectedIssuesJSON, _ := json.Marshal(inference.DetectedIssues)
	_, err = uc.auditLogRepo.Create(ctx, &entities.AuditLog{
		MaterialID:      materialEntity.ID,
		LabelID:         &labelEntity.ID,
		Status:          "COMPLETED",
		Provider:        "STUB_OCR_GEMINI",
		RawOCRText:      rawText,
		ConfidenceScore: inference.ConfidenceScore,
		Reasoning:       inference.Reasoning,
		DetectedIssues:  string(detectedIssuesJSON),
		CompletedAt:     auditservice.NowPtr(),
	})
	if err != nil {
		ctxLogger.Errorf("Failed to create audit log: %v", err)
		return nil, err
	}

	if err := uc.materialRepo.Update(ctx, materialEntity.ID, map[string]interface{}{
		"status":          "AI_REVIEWED",
		"latest_label_id": labelEntity.ID,
	}); err != nil {
		ctxLogger.Errorf("Failed to update material audit status: %v", err)
		return nil, err
	}

	detailedMaterial, err := uc.materialRepo.GetDetailed(ctx, materialEntity.ID)
	if err != nil {
		ctxLogger.Errorf("Failed to load material detail after upload: %v", err)
		return nil, err
	}

	view := mapMaterial(detailedMaterial)
	return &view, nil
}
