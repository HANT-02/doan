package audit

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"doan/pkg/utils"
)

type StorageService interface {
	Save(ctx context.Context, materialID string, fileName string, content []byte) (string, error)
	Resolve(ctx context.Context, storageKey string) (*StoredFile, error)
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

type StoredFile struct {
	StorageKey   string
	AbsolutePath string
	Size         int64
	OriginalName string
	ContentType  string
	ModifiedAt   time.Time
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

func (s *localStorageService) Save(_ context.Context, materialID string, fileName string, content []byte) (string, error) {
	now := time.Now()
	safeName := sanitizeFileName(fileName)
	relativePath := filepath.Join(
		fmt.Sprintf("%04d", now.Year()),
		fmt.Sprintf("%02d", int(now.Month())),
		materialID,
		safeName,
	)
	absolutePath := filepath.Join(s.basePath, relativePath)

	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		return "", err
	}

	if err := os.WriteFile(absolutePath, content, 0o644); err != nil {
		return "", err
	}

	return filepath.ToSlash(relativePath), nil
}

func (s *localStorageService) Resolve(_ context.Context, storageKey string) (*StoredFile, error) {
	cleanKey := filepath.Clean(filepath.FromSlash(storageKey))
	if cleanKey == "." || strings.HasPrefix(cleanKey, "..") || filepath.IsAbs(cleanKey) {
		return nil, fs.ErrPermission
	}

	basePathAbs, err := filepath.Abs(s.basePath)
	if err != nil {
		return nil, err
	}

	targetPath := filepath.Join(basePathAbs, cleanKey)
	targetPathAbs, err := filepath.Abs(targetPath)
	if err != nil {
		return nil, err
	}

	allowedPrefix := basePathAbs + string(os.PathSeparator)
	if targetPathAbs != basePathAbs && !strings.HasPrefix(targetPathAbs, allowedPrefix) {
		return nil, fs.ErrPermission
	}

	stat, err := os.Stat(targetPathAbs)
	if err != nil {
		return nil, err
	}

	return &StoredFile{
		StorageKey:   filepath.ToSlash(cleanKey),
		AbsolutePath: targetPathAbs,
		Size:         stat.Size(),
		OriginalName: filepath.Base(targetPathAbs),
		ModifiedAt:   stat.ModTime(),
	}, nil
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

func sanitizeFileName(fileName string) string {
	safeName := strings.ReplaceAll(strings.ToLower(fileName), " ", "_")
	safeName = strings.ReplaceAll(safeName, "..", "")
	safeName = strings.ReplaceAll(safeName, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	if safeName == "" {
		return utils.GenerateUUID()
	}

	return safeName
}
