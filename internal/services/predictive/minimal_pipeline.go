package predictive

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

var MinimalFeatureNames = []string{
	"attendance_rate_28d",
	"absence_count_28d",
	"average_total_score_28d",
	"homework_completion_rate_28d",
	"active_enrollment_count_28d",
	"weekly_lesson_load_28d",
	"approved_leave_count_28d",
	"days_since_last_lesson",
}

type TrainingRow struct {
	ID       string
	Features map[string]float64
	Label    string
}

type BinaryClassificationMetrics struct {
	Accuracy  float64 `json:"accuracy"`
	Precision float64 `json:"precision"`
	Recall    float64 `json:"recall"`
	F1Score   float64 `json:"f1_score"`
	Support   int     `json:"support"`
}

type MinimalModelReport struct {
	Name          string                      `json:"name"`
	Metrics       BinaryClassificationMetrics `json:"metrics"`
	PositiveClass string                      `json:"positive_class"`
	Notes         []string                    `json:"notes,omitempty"`
}

type MinimalTrainingReport struct {
	DatasetName    string               `json:"dataset_name"`
	TrainSize      int                  `json:"train_size"`
	TestSize       int                  `json:"test_size"`
	FeatureNames   []string             `json:"feature_names"`
	Recommended    string               `json:"recommended"`
	ModelReports   []MinimalModelReport `json:"model_reports"`
	Recommendation []string             `json:"recommendation"`
}

type MinimalTrainingArtifacts struct {
	Report          *MinimalTrainingReport
	FeatureNames    []string
	TrainRows       []TrainingRow
	TestRows        []TrainingRow
	LogisticModel   *LogisticRegressionModel
	Recommended     string
	RecommendedF1   float64
	RuleMetrics     BinaryClassificationMetrics
	LogisticMetrics BinaryClassificationMetrics
}

type LogisticRegressionConfig struct {
	LearningRate float64
	Epochs       int
	L2Penalty    float64
}

type LogisticRegressionModel struct {
	FeatureNames []string  `json:"feature_names"`
	Weights      []float64 `json:"weights"`
	Bias         float64   `json:"bias"`
	Means        []float64 `json:"means"`
	Scales       []float64 `json:"scales"`
}

func DefaultLogisticRegressionConfig() LogisticRegressionConfig {
	return LogisticRegressionConfig{
		LearningRate: 0.15,
		Epochs:       500,
		L2Penalty:    0.001,
	}
}

func BuildDemoAtRiskDataset(size int, seed int64) []TrainingRow {
	if size <= 0 {
		size = 80
	}

	rng := rand.New(rand.NewSource(seed))
	rows := make([]TrainingRow, 0, size)
	for index := 0; index < size; index++ {
		baseRisk := rng.Float64()

		attendanceRate := clamp(0.55+((1-baseRisk)*0.40)+rng.NormFloat64()*0.04, 0.35, 1.00)
		absenceCount := clamp(math.Round((1-attendanceRate)*10+rng.Float64()*2), 0, 12)
		averageScore := clamp(4.0+((1-baseRisk)*4.5)+rng.NormFloat64()*0.5, 2.5, 9.5)
		homeworkRate := clamp(0.45+((1-baseRisk)*0.50)+rng.NormFloat64()*0.05, 0.20, 1.00)
		enrollmentCount := clamp(math.Round(1+rng.Float64()*2), 1, 3)
		weeklyLessonLoad := clamp(1.5+enrollmentCount*1.1+rng.Float64()*2, 2, 8)
		approvedLeaveCount := clamp(math.Round(baseRisk*4+rng.Float64()*2), 0, 6)
		daysSinceLastLesson := clamp(math.Round(baseRisk*8+rng.Float64()*5), 0, 14)

		label := LabelNotAtRisk
		if attendanceRate < 0.80 || averageScore < 5.0 || homeworkRate < 0.60 {
			label = LabelAtRisk
		}

		rows = append(rows, TrainingRow{
			ID: fmt.Sprintf("demo-%03d", index+1),
			Features: map[string]float64{
				"attendance_rate_28d":          attendanceRate,
				"absence_count_28d":            absenceCount,
				"average_total_score_28d":      averageScore,
				"homework_completion_rate_28d": homeworkRate,
				"active_enrollment_count_28d":  enrollmentCount,
				"weekly_lesson_load_28d":       weeklyLessonLoad,
				"approved_leave_count_28d":     approvedLeaveCount,
				"days_since_last_lesson":       daysSinceLastLesson,
			},
			Label: label,
		})
	}

	return rows
}

func LoadTrainingRowsFromCSV(path string, featureNames []string) ([]TrainingRow, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, errors.New("csv dataset must contain header and at least one data row")
	}

	header := records[0]
	indexed := make(map[string]int, len(header))
	for idx, name := range header {
		indexed[strings.TrimSpace(name)] = idx
	}

	for _, field := range append(append([]string{}, featureNames...), "label") {
		if _, ok := indexed[field]; !ok {
			return nil, fmt.Errorf("missing required column %s in dataset", field)
		}
	}

	rows := make([]TrainingRow, 0, len(records)-1)
	for rowIndex, record := range records[1:] {
		row := TrainingRow{
			ID:       fmt.Sprintf("row-%03d", rowIndex+1),
			Features: make(map[string]float64, len(featureNames)),
			Label:    strings.TrimSpace(record[indexed["label"]]),
		}

		if idIndex, ok := indexed["id"]; ok && idIndex < len(record) {
			row.ID = strings.TrimSpace(record[idIndex])
		}

		for _, featureName := range featureNames {
			value, parseErr := strconv.ParseFloat(strings.TrimSpace(record[indexed[featureName]]), 64)
			if parseErr != nil {
				return nil, fmt.Errorf("parse feature %s at row %d: %w", featureName, rowIndex+2, parseErr)
			}
			row.Features[featureName] = value
		}

		if row.Label != LabelAtRisk && row.Label != LabelNotAtRisk {
			return nil, fmt.Errorf("invalid label at row %d: %s", rowIndex+2, row.Label)
		}

		rows = append(rows, row)
	}

	return rows, nil
}

func SplitTrainingRows(rows []TrainingRow, testRatio float64, seed int64) ([]TrainingRow, []TrainingRow) {
	if len(rows) == 0 {
		return nil, nil
	}
	if testRatio <= 0 || testRatio >= 1 {
		testRatio = 0.2
	}

	cloned := append([]TrainingRow(nil), rows...)
	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(cloned), func(i, j int) {
		cloned[i], cloned[j] = cloned[j], cloned[i]
	})

	testSize := int(math.Round(float64(len(cloned)) * testRatio))
	if testSize <= 0 {
		testSize = 1
	}
	if testSize >= len(cloned) {
		testSize = len(cloned) / 2
	}

	trainRows := append([]TrainingRow(nil), cloned[:len(cloned)-testSize]...)
	testRows := append([]TrainingRow(nil), cloned[len(cloned)-testSize:]...)
	return trainRows, testRows
}

func TrainLogisticRegression(rows []TrainingRow, featureNames []string, cfg LogisticRegressionConfig) (*LogisticRegressionModel, error) {
	if len(rows) == 0 {
		return nil, errors.New("cannot train logistic regression with empty dataset")
	}
	if len(featureNames) == 0 {
		return nil, errors.New("feature names are required")
	}

	means := make([]float64, len(featureNames))
	scales := make([]float64, len(featureNames))

	matrix, labels := rowsToMatrix(rows, featureNames)
	for col := range featureNames {
		for row := range matrix {
			means[col] += matrix[row][col]
		}
		means[col] /= float64(len(matrix))

		var variance float64
		for row := range matrix {
			diff := matrix[row][col] - means[col]
			variance += diff * diff
		}
		scales[col] = math.Sqrt(variance / float64(len(matrix)))
		if scales[col] == 0 {
			scales[col] = 1
		}

		for row := range matrix {
			matrix[row][col] = (matrix[row][col] - means[col]) / scales[col]
		}
	}

	weights := make([]float64, len(featureNames))
	bias := 0.0

	for epoch := 0; epoch < cfg.Epochs; epoch++ {
		for rowIndex, values := range matrix {
			probability := sigmoid(dot(weights, values) + bias)
			errorTerm := probability - labels[rowIndex]

			for col := range weights {
				gradient := errorTerm*values[col] + cfg.L2Penalty*weights[col]
				weights[col] -= cfg.LearningRate * gradient
			}
			bias -= cfg.LearningRate * errorTerm
		}
	}

	return &LogisticRegressionModel{
		FeatureNames: append([]string(nil), featureNames...),
		Weights:      weights,
		Bias:         bias,
		Means:        means,
		Scales:       scales,
	}, nil
}

func (m *LogisticRegressionModel) PredictProbability(row TrainingRow) float64 {
	values := make([]float64, len(m.FeatureNames))
	for idx, featureName := range m.FeatureNames {
		values[idx] = (row.Features[featureName] - m.Means[idx]) / m.Scales[idx]
	}
	return sigmoid(dot(m.Weights, values) + m.Bias)
}

func (m *LogisticRegressionModel) PredictLabel(row TrainingRow) string {
	if m.PredictProbability(row) >= 0.5 {
		return LabelAtRisk
	}
	return LabelNotAtRisk
}

func PredictRuleBasedLabel(row TrainingRow) string {
	if row.Features["attendance_rate_28d"] < 0.80 {
		return LabelAtRisk
	}
	if row.Features["average_total_score_28d"] < 5.0 {
		return LabelAtRisk
	}
	if row.Features["homework_completion_rate_28d"] < 0.60 {
		return LabelAtRisk
	}
	return LabelNotAtRisk
}

func EvaluateModel(rows []TrainingRow, predictor func(TrainingRow) string) BinaryClassificationMetrics {
	if len(rows) == 0 {
		return BinaryClassificationMetrics{}
	}

	var truePositive, trueNegative, falsePositive, falseNegative int
	for _, row := range rows {
		predicted := predictor(row)
		switch {
		case row.Label == LabelAtRisk && predicted == LabelAtRisk:
			truePositive++
		case row.Label == LabelNotAtRisk && predicted == LabelNotAtRisk:
			trueNegative++
		case row.Label == LabelNotAtRisk && predicted == LabelAtRisk:
			falsePositive++
		case row.Label == LabelAtRisk && predicted == LabelNotAtRisk:
			falseNegative++
		}
	}

	accuracy := float64(truePositive+trueNegative) / float64(len(rows))
	precision := safeDivide(float64(truePositive), float64(truePositive+falsePositive))
	recall := safeDivide(float64(truePositive), float64(truePositive+falseNegative))
	f1 := safeDivide(2*precision*recall, precision+recall)

	return BinaryClassificationMetrics{
		Accuracy:  accuracy,
		Precision: precision,
		Recall:    recall,
		F1Score:   f1,
		Support:   len(rows),
	}
}

func RunMinimalTrainingComparison(datasetName string, rows []TrainingRow, seed int64) (*MinimalTrainingReport, error) {
	artifacts, err := TrainMinimalModels(datasetName, rows, seed)
	if err != nil {
		return nil, err
	}

	return artifacts.Report, nil
}

func TrainMinimalModels(datasetName string, rows []TrainingRow, seed int64) (*MinimalTrainingArtifacts, error) {
	if len(rows) < 10 {
		return nil, errors.New("dataset must contain at least 10 rows for minimal train/test split")
	}

	featureNames := append([]string(nil), MinimalFeatureNames...)
	sort.Strings(featureNames)
	trainRows, testRows := SplitTrainingRows(rows, 0.2, seed)

	logisticModel, err := TrainLogisticRegression(trainRows, featureNames, DefaultLogisticRegressionConfig())
	if err != nil {
		return nil, err
	}

	ruleMetrics := EvaluateModel(testRows, PredictRuleBasedLabel)
	logisticMetrics := EvaluateModel(testRows, logisticModel.PredictLabel)

	modelReports := []MinimalModelReport{
		{
			Name:          "rule_based_baseline",
			Metrics:       ruleMetrics,
			PositiveClass: LabelAtRisk,
			Notes: []string{
				"Khong can train, dung de lam baseline nhanh de so sanh.",
				"Rule duoc suy ra truc tiep tu label definition cua F1.",
			},
		},
		{
			Name:          "logistic_regression",
			Metrics:       logisticMetrics,
			PositiveClass: LabelAtRisk,
			Notes: []string{
				"Train bang gradient descent thuan Go, khong phu thuoc sklearn.",
				"Phu hop may yeu va dataset tabular nho-trung binh.",
			},
		},
	}

	recommended := "rule_based_baseline"
	reason := []string{
		"Mac dinh uu tien model co F1 cua lop AT_RISK cao hon.",
		"Neu F1 bang nhau thi uu tien Logistic Regression vi van la ML model thuc su nhung van rat nhe.",
	}
	if logisticMetrics.F1Score >= ruleMetrics.F1Score {
		recommended = "logistic_regression"
	}

	report := &MinimalTrainingReport{
		DatasetName:    datasetName,
		TrainSize:      len(trainRows),
		TestSize:       len(testRows),
		FeatureNames:   featureNames,
		Recommended:    recommended,
		ModelReports:   modelReports,
		Recommendation: reason,
	}

	recommendedF1 := ruleMetrics.F1Score
	if recommended == "logistic_regression" {
		recommendedF1 = logisticMetrics.F1Score
	}

	return &MinimalTrainingArtifacts{
		Report:          report,
		FeatureNames:    featureNames,
		TrainRows:       trainRows,
		TestRows:        testRows,
		LogisticModel:   logisticModel,
		Recommended:     recommended,
		RecommendedF1:   recommendedF1,
		RuleMetrics:     ruleMetrics,
		LogisticMetrics: logisticMetrics,
	}, nil
}

func RenderMinimalTrainingReportMarkdown(report *MinimalTrainingReport) string {
	var builder strings.Builder
	builder.WriteString("# Minimal AT_RISK Training Report\n\n")
	builder.WriteString(fmt.Sprintf("- Dataset: %s\n", report.DatasetName))
	builder.WriteString(fmt.Sprintf("- Train size: %d\n", report.TrainSize))
	builder.WriteString(fmt.Sprintf("- Test size: %d\n", report.TestSize))
	builder.WriteString(fmt.Sprintf("- Recommended model: `%s`\n\n", report.Recommended))
	builder.WriteString("## Models\n\n")
	builder.WriteString("| Model | Accuracy | Precision | Recall | F1 |\n")
	builder.WriteString("| --- | ---: | ---: | ---: | ---: |\n")
	for _, item := range report.ModelReports {
		builder.WriteString(fmt.Sprintf(
			"| %s | %.3f | %.3f | %.3f | %.3f |\n",
			item.Name,
			item.Metrics.Accuracy,
			item.Metrics.Precision,
			item.Metrics.Recall,
			item.Metrics.F1Score,
		))
	}
	builder.WriteString("\n## Notes\n\n")
	for _, note := range report.Recommendation {
		builder.WriteString(fmt.Sprintf("- %s\n", note))
	}
	return builder.String()
}

func rowsToMatrix(rows []TrainingRow, featureNames []string) ([][]float64, []float64) {
	matrix := make([][]float64, 0, len(rows))
	labels := make([]float64, 0, len(rows))
	for _, row := range rows {
		values := make([]float64, len(featureNames))
		for idx, featureName := range featureNames {
			values[idx] = row.Features[featureName]
		}
		matrix = append(matrix, values)
		if row.Label == LabelAtRisk {
			labels = append(labels, 1)
		} else {
			labels = append(labels, 0)
		}
	}
	return matrix, labels
}

func dot(weights, values []float64) float64 {
	total := 0.0
	for idx := range weights {
		total += weights[idx] * values[idx]
	}
	return total
}

func sigmoid(value float64) float64 {
	if value < -35 {
		return 0
	}
	if value > 35 {
		return 1
	}
	return 1 / (1 + math.Exp(-value))
}

func safeDivide(numerator, denominator float64) float64 {
	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

func clamp(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
