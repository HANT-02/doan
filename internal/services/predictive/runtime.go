package predictive

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

const (
	defaultPredictiveArtifactsRelativePath = "ml/at_risk_prediction/artifacts"
	predictiveArtifactsEnvKey              = "PREDICTIVE_ARTIFACTS_DIR"
	predictiveDefinitionVersion            = "python_ml_v1"
)

type AtRiskModelMetadata struct {
	Version        string               `json:"version"`
	ModelName      string               `json:"model_name"`
	DatasetName    string               `json:"dataset_name"`
	DatasetSource  string               `json:"dataset_source"`
	DefinitionName string               `json:"definition_name"`
	DefinitionVer  string               `json:"definition_version"`
	TrainedAt      time.Time            `json:"trained_at"`
	TrainSize      int                  `json:"train_size"`
	TestSize       int                  `json:"test_size"`
	FeatureNames   []string             `json:"feature_names"`
	ModelReports   []MinimalModelReport `json:"model_reports"`
	Recommendation []string             `json:"recommendation"`
}

type AtRiskPredictionFeatureSnapshot struct {
	AttendanceRate28d         float64 `json:"attendance_rate_28d"`
	AbsenceCount28d           float64 `json:"absence_count_28d"`
	AverageTotalScore28d      float64 `json:"average_total_score_28d"`
	HomeworkCompletionRate28d float64 `json:"homework_completion_rate_28d"`
	ActiveEnrollmentCount28d  float64 `json:"active_enrollment_count_28d"`
	WeeklyLessonLoad28d       float64 `json:"weekly_lesson_load_28d"`
	ApprovedLeaveCount28d     float64 `json:"approved_leave_count_28d"`
	DaysSinceLastLesson       float64 `json:"days_since_last_lesson"`
}

type StudentRiskPrediction struct {
	StudentID      string                          `json:"student_id"`
	StudentCode    string                          `json:"student_code"`
	StudentName    string                          `json:"student_name"`
	GradeLevel     string                          `json:"grade_level"`
	ClassID        string                          `json:"class_id"`
	ClassCode      string                          `json:"class_code"`
	ClassName      string                          `json:"class_name"`
	SnapshotAt     time.Time                       `json:"snapshot_at"`
	Label          string                          `json:"label"`
	RiskScore      float64                         `json:"risk_score"`
	RiskBand       string                          `json:"risk_band"`
	ModelName      string                          `json:"model_name"`
	ModelVersion   string                          `json:"model_version"`
	PrimaryReason  string                          `json:"primary_reason"`
	Reasons        []string                        `json:"reasons"`
	FeatureSummary AtRiskPredictionFeatureSnapshot `json:"feature_summary"`
}

type ListStudentPredictionsInput struct {
	Search     string
	OnlyAtRisk bool
	Page       int
	Limit      int
	Refresh    bool
}

type ListStudentPredictionsOutput struct {
	Items      []StudentRiskPrediction `json:"items"`
	Pagination struct {
		CurrentPage  int   `json:"current_page"`
		ItemsPerPage int   `json:"items_per_page"`
		TotalItems   int64 `json:"total_items"`
		TotalPages   int   `json:"total_pages"`
	} `json:"pagination"`
	Summary struct {
		TotalStudentsEvaluated int `json:"total_students_evaluated"`
		AtRiskCount            int `json:"at_risk_count"`
		NotAtRiskCount         int `json:"not_at_risk_count"`
	} `json:"summary"`
	ModelMetadata AtRiskModelMetadata `json:"model_metadata"`
}

type artifactSnapshot struct {
	Metadata    AtRiskModelMetadata
	Predictions []StudentRiskPrediction
	LoadedAt    time.Time
}

type AtRiskService interface {
	ListStudentPredictions(ctx context.Context, input ListStudentPredictionsInput) (*ListStudentPredictionsOutput, error)
	GetModelMetadata(ctx context.Context, refresh bool) (*AtRiskModelMetadata, error)
}

type atRiskService struct {
	mu           sync.RWMutex
	cached       *artifactSnapshot
	artifactsDir string
}

type pythonModelMetadataFile struct {
	GeneratedAt   string   `json:"generated_at"`
	DatasetName   string   `json:"dataset_name"`
	Source        string   `json:"source"`
	FeatureColumn []string `json:"feature_columns"`
	LabelMapping  struct {
		PositiveLabel string `json:"positive_label"`
		NegativeLabel string `json:"negative_label"`
	} `json:"label_mapping"`
	Selection pythonSelection `json:"selection"`
}

type pythonSelection struct {
	SelectedModel    string   `json:"selected_model"`
	TiedModels       []string `json:"tied_models"`
	SelectionCritera []string `json:"selection_criteria"`
	Rationale        string   `json:"rationale"`
	SelectedAt       string   `json:"selected_at"`
}

type pythonMetricsFile struct {
	GeneratedAt    string `json:"generated_at"`
	DatasetSummary struct {
		DatasetName string   `json:"dataset_name"`
		Source      string   `json:"source"`
		FeatureCols []string `json:"feature_columns"`
	} `json:"dataset_summary"`
	TrainingConfig struct {
		DatasetName  string `json:"dataset_name"`
		Source       string `json:"source"`
		TrainSize    int    `json:"train_size"`
		TestSizeRows int    `json:"test_size_rows"`
	} `json:"training_config"`
	Models map[string]pythonModelMetrics `json:"models"`
}

type pythonModelMetrics struct {
	ModelType string  `json:"model_type"`
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	F1        float64 `json:"f1"`
	Support   int     `json:"support"`
}

type pythonPredictionsFile struct {
	GeneratedAt     string                   `json:"generated_at"`
	Source          string                   `json:"source"`
	ModelVersion    string                   `json:"model_version"`
	PredictionCount int                      `json:"prediction_count"`
	Note            *string                  `json:"note"`
	Predictions     []pythonPredictionRecord `json:"predictions"`
}

type pythonPredictionRecord struct {
	SnapshotID     string               `json:"snapshot_id"`
	StudentID      string               `json:"student_id"`
	StudentCode    string               `json:"student_code"`
	StudentName    string               `json:"student_name"`
	ClassID        string               `json:"class_id"`
	ClassCode      string               `json:"class_code"`
	ClassName      string               `json:"class_name"`
	SnapshotAt     string               `json:"snapshot_at"`
	RiskLabel      string               `json:"risk_label"`
	RiskScore      float64              `json:"risk_score"`
	RiskBand       string               `json:"risk_band"`
	PrimaryReason  string               `json:"primary_reason"`
	TopFeatures    []pythonTopFeature   `json:"top_features"`
	FeatureSummary pythonFeatureSummary `json:"feature_summary"`
	ModelVersion   string               `json:"model_version"`
	PredictedAt    string               `json:"predicted_at"`
}

type pythonTopFeature struct {
	Feature  string  `json:"feature"`
	Value    float64 `json:"value"`
	Severity float64 `json:"severity"`
	Label    string  `json:"label"`
	Detail   string  `json:"detail"`
}

type pythonFeatureSummary struct {
	AttendanceRate28d         float64 `json:"attendance_rate_28d"`
	AbsenceCount28d           float64 `json:"absence_count_28d"`
	AverageTotalScore28d      float64 `json:"average_total_score_28d"`
	HomeworkCompletionRate28d float64 `json:"homework_completion_rate_28d"`
	ActiveEnrollmentCount28d  float64 `json:"active_enrollment_count_28d"`
	WeeklyLessonLoad28d       float64 `json:"weekly_lesson_load_28d"`
	ApprovedLeaveCount28d     float64 `json:"approved_leave_count_28d"`
	DaysSinceLastLesson       float64 `json:"days_since_last_lesson"`
}

func NewAtRiskService(_ *gorm.DB) AtRiskService {
	return &atRiskService{
		artifactsDir: resolvePredictiveArtifactsDir(),
	}
}

func (s *atRiskService) GetModelMetadata(ctx context.Context, refresh bool) (*AtRiskModelMetadata, error) {
	snapshot, err := s.getOrLoadSnapshot(ctx, refresh)
	if err != nil {
		return nil, err
	}

	metadata := snapshot.Metadata
	return &metadata, nil
}

func (s *atRiskService) ListStudentPredictions(ctx context.Context, input ListStudentPredictionsInput) (*ListStudentPredictionsOutput, error) {
	snapshot, err := s.getOrLoadSnapshot(ctx, input.Refresh)
	if err != nil {
		return nil, err
	}

	predictions := append([]StudentRiskPrediction(nil), snapshot.Predictions...)
	sort.Slice(predictions, func(i, j int) bool {
		if predictions[i].Label != predictions[j].Label {
			return predictions[i].Label == LabelAtRisk
		}
		if math.Abs(predictions[i].RiskScore-predictions[j].RiskScore) > 0.0001 {
			return predictions[i].RiskScore > predictions[j].RiskScore
		}
		if predictions[i].StudentName != predictions[j].StudentName {
			return predictions[i].StudentName < predictions[j].StudentName
		}
		return predictions[i].ClassCode < predictions[j].ClassCode
	})

	filtered := filterPredictions(predictions, input.Search, input.OnlyAtRisk)
	page, limit := normalizePaging(input.Page, input.Limit)
	start, end, totalPages := paginateRange(len(filtered), page, limit)

	output := &ListStudentPredictionsOutput{
		Items:         append([]StudentRiskPrediction(nil), filtered[start:end]...),
		ModelMetadata: snapshot.Metadata,
	}
	output.Pagination.CurrentPage = page
	output.Pagination.ItemsPerPage = limit
	output.Pagination.TotalItems = int64(len(filtered))
	output.Pagination.TotalPages = totalPages
	output.Summary.TotalStudentsEvaluated = len(predictions)
	for _, item := range predictions {
		if item.Label == LabelAtRisk {
			output.Summary.AtRiskCount++
		} else {
			output.Summary.NotAtRiskCount++
		}
	}

	return output, nil
}

func (s *atRiskService) getOrLoadSnapshot(_ context.Context, refresh bool) (*artifactSnapshot, error) {
	if !refresh {
		s.mu.RLock()
		if s.cached != nil {
			cached := s.cached
			s.mu.RUnlock()
			return cached, nil
		}
		s.mu.RUnlock()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if !refresh && s.cached != nil {
		return s.cached, nil
	}

	snapshot, err := loadArtifactSnapshot(s.artifactsDir)
	if err != nil {
		return nil, err
	}
	s.cached = snapshot
	return snapshot, nil
}

func loadArtifactSnapshot(artifactsDir string) (*artifactSnapshot, error) {
	modelMetadataPath := filepath.Join(artifactsDir, "models", "model_metadata.json")
	metricsPath := filepath.Join(artifactsDir, "reports", "metrics.json")
	predictionsPath := filepath.Join(artifactsDir, "reports", "latest_predictions.json")

	var modelMetadata pythonModelMetadataFile
	if err := readJSONFile(modelMetadataPath, &modelMetadata); err != nil {
		return nil, artifactReadError("model metadata", modelMetadataPath, err)
	}

	var metrics pythonMetricsFile
	if err := readJSONFile(metricsPath, &metrics); err != nil {
		return nil, artifactReadError("metrics", metricsPath, err)
	}

	var predictionsFile pythonPredictionsFile
	if err := readJSONFile(predictionsPath, &predictionsFile); err != nil {
		return nil, artifactReadError("prediction", predictionsPath, err)
	}

	metadata := mapModelMetadata(modelMetadata, metrics)
	predictions := mapPredictions(predictionsFile, metadata)
	if len(predictions) == 0 {
		note := ""
		if predictionsFile.Note != nil {
			note = strings.TrimSpace(*predictionsFile.Note)
		}
		if note == "" {
			note = "khong co du lieu prediction hien tai"
		}
		return nil, fmt.Errorf("khong co du lieu prediction hien tai: %s", note)
	}

	return &artifactSnapshot{
		Metadata:    metadata,
		Predictions: predictions,
		LoadedAt:    time.Now(),
	}, nil
}

func mapModelMetadata(modelMetadata pythonModelMetadataFile, metrics pythonMetricsFile) AtRiskModelMetadata {
	trainedAt := parseTimeOrNow(firstNonEmpty(modelMetadata.Selection.SelectedAt, modelMetadata.GeneratedAt, metrics.GeneratedAt))
	selectedModel := firstNonEmpty(modelMetadata.Selection.SelectedModel, "logistic_regression")
	version := fmt.Sprintf("%s-%s", selectedModel, trainedAt.Format("20060102150405"))

	featureNames := append([]string(nil), modelMetadata.FeatureColumn...)
	if len(featureNames) == 0 {
		featureNames = append([]string(nil), metrics.DatasetSummary.FeatureCols...)
	}

	modelReports := make([]MinimalModelReport, 0, len(metrics.Models))
	for _, modelName := range []string{"rule_based", "logistic_regression", "random_forest"} {
		metricsItem, ok := metrics.Models[modelName]
		if !ok {
			continue
		}
		modelReports = append(modelReports, MinimalModelReport{
			Name:          modelName,
			PositiveClass: LabelAtRisk,
			Metrics: BinaryClassificationMetrics{
				Accuracy:  metricsItem.Accuracy,
				Precision: metricsItem.Precision,
				Recall:    metricsItem.Recall,
				F1Score:   metricsItem.F1,
				Support:   metricsItem.Support,
			},
			Notes: []string{
				fmt.Sprintf("model_type=%s", metricsItem.ModelType),
			},
		})
	}

	recommendation := []string{
		fmt.Sprintf("Mô hình chính hiện dùng: %s.", selectedModel),
	}
	if rationale := strings.TrimSpace(modelMetadata.Selection.Rationale); rationale != "" {
		recommendation = append(recommendation, rationale)
	}
	if len(modelMetadata.Selection.SelectionCritera) > 0 {
		recommendation = append(
			recommendation,
			"Tiêu chí ưu tiên: "+strings.Join(modelMetadata.Selection.SelectionCritera, " -> ")+".",
		)
	}

	return AtRiskModelMetadata{
		Version:        version,
		ModelName:      selectedModel,
		DatasetName:    firstNonEmpty(modelMetadata.DatasetName, metrics.TrainingConfig.DatasetName, metrics.DatasetSummary.DatasetName),
		DatasetSource:  firstNonEmpty(modelMetadata.Source, metrics.TrainingConfig.Source, metrics.DatasetSummary.Source),
		DefinitionName: "student_at_risk_classification",
		DefinitionVer:  predictiveDefinitionVersion,
		TrainedAt:      trainedAt,
		TrainSize:      metrics.TrainingConfig.TrainSize,
		TestSize:       metrics.TrainingConfig.TestSizeRows,
		FeatureNames:   featureNames,
		ModelReports:   modelReports,
		Recommendation: recommendation,
	}
}

func mapPredictions(payload pythonPredictionsFile, metadata AtRiskModelMetadata) []StudentRiskPrediction {
	items := make([]StudentRiskPrediction, 0, len(payload.Predictions))
	for _, item := range payload.Predictions {
		reasons := make([]string, 0, len(item.TopFeatures))
		for _, top := range item.TopFeatures {
			if strings.TrimSpace(top.Detail) != "" {
				reasons = append(reasons, top.Detail)
				continue
			}
			if strings.TrimSpace(top.Label) != "" {
				reasons = append(reasons, top.Label)
			}
		}
		if len(reasons) == 0 && strings.TrimSpace(item.PrimaryReason) != "" {
			reasons = append(reasons, item.PrimaryReason)
		}

		items = append(items, StudentRiskPrediction{
			StudentID:     item.StudentID,
			StudentCode:   item.StudentCode,
			StudentName:   item.StudentName,
			GradeLevel:    "",
			ClassID:       item.ClassID,
			ClassCode:     item.ClassCode,
			ClassName:     item.ClassName,
			SnapshotAt:    parseTimeOrNow(item.SnapshotAt),
			Label:         item.RiskLabel,
			RiskScore:     item.RiskScore,
			RiskBand:      normalizeRiskBand(item.RiskBand),
			ModelName:     metadata.ModelName,
			ModelVersion:  metadata.Version,
			PrimaryReason: item.PrimaryReason,
			Reasons:       reasons,
			FeatureSummary: AtRiskPredictionFeatureSnapshot{
				AttendanceRate28d:         item.FeatureSummary.AttendanceRate28d,
				AbsenceCount28d:           item.FeatureSummary.AbsenceCount28d,
				AverageTotalScore28d:      item.FeatureSummary.AverageTotalScore28d,
				HomeworkCompletionRate28d: item.FeatureSummary.HomeworkCompletionRate28d,
				ActiveEnrollmentCount28d:  item.FeatureSummary.ActiveEnrollmentCount28d,
				WeeklyLessonLoad28d:       item.FeatureSummary.WeeklyLessonLoad28d,
				ApprovedLeaveCount28d:     item.FeatureSummary.ApprovedLeaveCount28d,
				DaysSinceLastLesson:       item.FeatureSummary.DaysSinceLastLesson,
			},
		})
	}
	return items
}

func readJSONFile(path string, target any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, target); err != nil {
		return err
	}
	return nil
}

func artifactReadError(kind, path string, err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("khong co du lieu prediction hien tai: thieu artifact %s (%s)", kind, path)
	}
	return fmt.Errorf("khong doc duoc artifact %s (%s): %w", kind, path, err)
}

func resolvePredictiveArtifactsDir() string {
	if envPath := strings.TrimSpace(os.Getenv(predictiveArtifactsEnvKey)); envPath != "" {
		return envPath
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return defaultPredictiveArtifactsRelativePath
	}

	current := workingDir
	for {
		candidate := filepath.Join(current, defaultPredictiveArtifactsRelativePath)
		if _, statErr := os.Stat(candidate); statErr == nil {
			return candidate
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return filepath.Join(workingDir, defaultPredictiveArtifactsRelativePath)
}

func normalizeRiskBand(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "HIGH", "CAO":
		return "CAO"
	case "MEDIUM", "TRUNG_BINH", "TRUNG BÌNH":
		return "TRUNG_BINH"
	default:
		return "THAP"
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func parseTimeOrNow(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Now()
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed
		}
	}
	return time.Now()
}

func filterPredictions(items []StudentRiskPrediction, search string, onlyAtRisk bool) []StudentRiskPrediction {
	filtered := make([]StudentRiskPrediction, 0, len(items))
	keyword := strings.ToLower(strings.TrimSpace(search))
	for _, item := range items {
		if onlyAtRisk && item.Label != LabelAtRisk {
			continue
		}
		if keyword != "" {
			searchable := strings.ToLower(strings.Join([]string{
				item.StudentCode,
				item.StudentName,
				item.ClassCode,
				item.ClassName,
				item.PrimaryReason,
			}, " "))
			if !strings.Contains(searchable, keyword) {
				continue
			}
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func normalizePaging(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return page, limit
}

func paginateRange(total, page, limit int) (int, int, int) {
	if total == 0 {
		return 0, 0, 0
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	return start, end, totalPages
}
