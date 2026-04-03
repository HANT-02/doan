package material

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"path"
	"path/filepath"
	"strings"

	"doan/internal/entities"
	repositoryinterface "doan/internal/repositories/interface"
	auditservice "doan/internal/services/audit"
	"doan/pkg/logger"
	"doan/pkg/utils"
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

const maxMaterialFileSize int64 = 10 * 1024 * 1024

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

	normalizedFileType, err := normalizeMaterialFileType(input.FileName, input.FileType)
	if err != nil {
		return nil, err
	}
	input.FileType = normalizedFileType

	if err := validateMaterialUpload(input.FileName, input.FileType, int64(len(input.Content))); err != nil {
		return nil, err
	}

	materialID := utils.GenerateUUID()
	filePath, err := uc.storage.Save(ctx, materialID, input.FileName, input.Content)
	if err != nil {
		ctxLogger.Errorf("Failed to store material file: %v", err)
		return nil, err
	}

	if input.Title == "" {
		input.Title = filepath.Base(input.FileName)
	}

	materialEntity, err := uc.materialRepo.Create(ctx, &entities.Material{
		ID:          materialID,
		TeacherID:   input.TeacherID,
		Title:       input.Title,
		Description: input.Description,
		FileName:    input.FileName,
		FilePath:    filePath,
		FileType:    input.FileType,
		FileSize:    int64(len(input.Content)),
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

func validateMaterialUpload(fileName string, fileType string, fileSize int64) error {
	if fileSize <= 0 {
		return errors.New("file is empty")
	}
	if fileSize > maxMaterialFileSize {
		return fmt.Errorf("file exceeds maximum size of %d MB", maxMaterialFileSize/(1024*1024))
	}

	extension := strings.ToLower(path.Ext(fileName))
	contentType := strings.ToLower(strings.TrimSpace(fileType))
	if contentType == "" {
		contentType = strings.ToLower(mime.TypeByExtension(extension))
	}

	allowedByExtension := map[string]struct{}{
		".pdf":  {},
		".doc":  {},
		".docx": {},
		".png":  {},
		".jpg":  {},
		".jpeg": {},
	}
	allowedByType := map[string]struct{}{
		"application/pdf":    {},
		"application/msword": {},
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {},
		"image/png":  {},
		"image/jpeg": {},
	}

	if _, ok := allowedByExtension[extension]; !ok {
		return fmt.Errorf("unsupported file extension: %s", extension)
	}
	if contentType != "" {
		if _, ok := allowedByType[contentType]; !ok {
			return fmt.Errorf("unsupported file type: %s", contentType)
		}
	}

	return nil
}

func normalizeMaterialFileType(fileName string, fileType string) (string, error) {
	extension := strings.ToLower(path.Ext(fileName))
	contentType := strings.ToLower(strings.TrimSpace(fileType))
	if contentType != "" {
		return contentType, nil
	}

	guessedType := strings.ToLower(mime.TypeByExtension(extension))
	if guessedType == "" {
		return "", fmt.Errorf("cannot determine file type for extension: %s", extension)
	}

	return guessedType, nil
}
