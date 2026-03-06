package audit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"doan/pkg/utils"
)

type StorageService interface {
	Save(ctx context.Context, fileName string, content []byte) (string, error)
}

type OCRService interface {
	ExtractText(ctx context.Context, fileName string, content []byte) (string, error)
}

type GeminiService interface {
	Analyze(ctx context.Context, text string) (*AuditInference, error)
}

type AuditInference struct {
	LabelCode       string
	ConfidenceScore float64
	Reasoning       string
	DetectedIssues  []string
}

type localStorageService struct {
	basePath string
}

type stubOCRService struct{}

type stubGeminiService struct{}

func NewLocalStorageService() StorageService {
	return &localStorageService{
		basePath: filepath.Join("storage", "materials"),
	}
}

func NewStubOCRService() OCRService {
	return &stubOCRService{}
}

func NewStubGeminiService() GeminiService {
	return &stubGeminiService{}
}

func (s *localStorageService) Save(_ context.Context, fileName string, content []byte) (string, error) {
	if err := os.MkdirAll(s.basePath, 0o755); err != nil {
		return "", err
	}

	safeName := strings.ReplaceAll(strings.ToLower(fileName), " ", "_")
	filePath := filepath.Join(s.basePath, fmt.Sprintf("%s_%s", utils.GenerateUUID(), safeName))
	if err := os.WriteFile(filePath, content, 0o644); err != nil {
		return "", err
	}

	return filePath, nil
}

func (s *stubOCRService) ExtractText(_ context.Context, fileName string, content []byte) (string, error) {
	if len(content) == 0 {
		return fmt.Sprintf("stub-ocr:%s:no-content", fileName), nil
	}

	preview := string(content)
	if len(preview) > 500 {
		preview = preview[:500]
	}

	return fmt.Sprintf("stub-ocr:%s:%s", fileName, preview), nil
}

func (s *stubGeminiService) Analyze(_ context.Context, text string) (*AuditInference, error) {
	lowerText := strings.ToLower(text)
	inference := &AuditInference{
		LabelCode:       "SAFE",
		ConfidenceScore: 0.92,
		Reasoning:       "Stub Gemini xác định nội dung an toàn cho demo.",
		DetectedIssues:  []string{},
	}

	switch {
	case strings.Contains(lowerText, "violence") || strings.Contains(lowerText, "gambling") || strings.Contains(lowerText, "danger"):
		inference.LabelCode = "DANGER"
		inference.ConfidenceScore = 0.88
		inference.Reasoning = "Stub Gemini phát hiện tín hiệu nội dung rủi ro cao."
		inference.DetectedIssues = []string{"Phát hiện từ khóa nhạy cảm", "Cần compliance review"}
	case strings.Contains(lowerText, "exam") || strings.Contains(lowerText, "cheat") || strings.Contains(lowerText, "warning"):
		inference.LabelCode = "WARNING"
		inference.ConfidenceScore = 0.76
		inference.Reasoning = "Stub Gemini phát hiện nội dung cần rà soát thêm."
		inference.DetectedIssues = []string{"Có tín hiệu cần kiểm duyệt thủ công"}
	}

	return inference, nil
}

func NowPtr() *time.Time {
	now := time.Now()
	return &now
}
