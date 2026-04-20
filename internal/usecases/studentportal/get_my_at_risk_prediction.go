package studentportal

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	repointerface "doan/internal/repositories/interface"
	predictiveservice "doan/internal/services/predictive"
	"doan/pkg/logger"
)

type GetMyAtRiskPredictionInput struct {
	Actor   Actor
	Refresh bool
}

type StudentAtRiskTopFeature struct {
	Key          string  `json:"key"`
	Label        string  `json:"label"`
	Value        float64 `json:"value"`
	DisplayValue string  `json:"display_value"`
}

type StudentAtRiskPrediction struct {
	StudentID      string                                            `json:"student_id"`
	StudentCode    string                                            `json:"student_code"`
	StudentName    string                                            `json:"student_name"`
	GradeLevel     string                                            `json:"grade_level"`
	ClassID        string                                            `json:"class_id"`
	ClassCode      string                                            `json:"class_code"`
	ClassName      string                                            `json:"class_name"`
	SnapshotAt     string                                            `json:"snapshot_at"`
	RiskLabel      string                                            `json:"risk_label"`
	RiskScore      float64                                           `json:"risk_score"`
	RiskBand       string                                            `json:"risk_band"`
	ModelName      string                                            `json:"model_name"`
	ModelVersion   string                                            `json:"model_version"`
	PrimaryReason  string                                            `json:"primary_reason"`
	Reasons        []string                                          `json:"reasons"`
	TopFeatures    []StudentAtRiskTopFeature                         `json:"top_features"`
	FeatureSummary predictiveservice.AtRiskPredictionFeatureSnapshot `json:"feature_summary"`
}

type GetMyAtRiskPredictionOutput struct {
	StudentID  string                   `json:"student_id"`
	Prediction *StudentAtRiskPrediction `json:"prediction,omitempty"`
}

type GetMyAtRiskPredictionUseCase interface {
	Execute(ctx context.Context, input GetMyAtRiskPredictionInput) (*GetMyAtRiskPredictionOutput, error)
}

type getMyAtRiskPredictionUseCase struct {
	studentRepo   repointerface.StudentRepository
	atRiskService predictiveservice.AtRiskService
}

func NewGetMyAtRiskPredictionUseCase(
	studentRepo repointerface.StudentRepository,
	atRiskService predictiveservice.AtRiskService,
) GetMyAtRiskPredictionUseCase {
	return &getMyAtRiskPredictionUseCase{
		studentRepo:   studentRepo,
		atRiskService: atRiskService,
	}
}

func (uc *getMyAtRiskPredictionUseCase) Execute(ctx context.Context, input GetMyAtRiskPredictionInput) (*GetMyAtRiskPredictionOutput, error) {
	ctxLogger := logger.NewLogger(ctx)

	if strings.TrimSpace(input.Actor.Role) != "STUDENT" {
		return nil, ErrStudentAccessDenied
	}

	student, err := resolveStudentByEmail(ctx, uc.studentRepo, input.Actor.Email)
	if err != nil {
		ctxLogger.Errorf("Failed to resolve student from actor email %s: %v", input.Actor.Email, err)
		return nil, err
	}

	predictions, err := uc.atRiskService.ListStudentPredictions(ctx, predictiveservice.ListStudentPredictionsInput{
		Page:    1,
		Limit:   100000,
		Refresh: input.Refresh,
	})
	if err != nil {
		if isPredictionUnavailable(err) {
			return &GetMyAtRiskPredictionOutput{
				StudentID:  student.ID,
				Prediction: nil,
			}, nil
		}
		ctxLogger.Errorf("Failed to load at-risk predictions for student %s: %v", student.ID, err)
		return nil, err
	}

	matches := make([]predictiveservice.StudentRiskPrediction, 0, 4)
	for _, item := range predictions.Items {
		if item.StudentID == student.ID {
			matches = append(matches, item)
		}
	}

	if len(matches) == 0 {
		return &GetMyAtRiskPredictionOutput{
			StudentID:  student.ID,
			Prediction: nil,
		}, nil
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Label != matches[j].Label {
			return matches[i].Label == predictiveservice.LabelAtRisk
		}
		if matches[i].RiskScore != matches[j].RiskScore {
			return matches[i].RiskScore > matches[j].RiskScore
		}
		return matches[i].SnapshotAt.After(matches[j].SnapshotAt)
	})

	selected := matches[0]
	return &GetMyAtRiskPredictionOutput{
		StudentID: student.ID,
		Prediction: &StudentAtRiskPrediction{
			StudentID:      selected.StudentID,
			StudentCode:    selected.StudentCode,
			StudentName:    selected.StudentName,
			GradeLevel:     selected.GradeLevel,
			ClassID:        selected.ClassID,
			ClassCode:      selected.ClassCode,
			ClassName:      selected.ClassName,
			SnapshotAt:     selected.SnapshotAt.Format(time.RFC3339),
			RiskLabel:      selected.Label,
			RiskScore:      selected.RiskScore,
			RiskBand:       selected.RiskBand,
			ModelName:      selected.ModelName,
			ModelVersion:   selected.ModelVersion,
			PrimaryReason:  selected.PrimaryReason,
			Reasons:        append([]string(nil), selected.Reasons...),
			TopFeatures:    buildTopFeatures(selected.FeatureSummary),
			FeatureSummary: selected.FeatureSummary,
		},
	}, nil
}

func isPredictionUnavailable(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "khong co du lieu prediction") ||
		strings.Contains(message, "khong tim thay candidate prediction")
}

func buildTopFeatures(summary predictiveservice.AtRiskPredictionFeatureSnapshot) []StudentAtRiskTopFeature {
	candidates := []struct {
		key      string
		label    string
		value    float64
		display  string
		severity float64
	}{
		{
			key:      "attendance_rate_28d",
			label:    "Chuyên cần 28 ngày",
			value:    summary.AttendanceRate28d,
			display:  formatPercent(summary.AttendanceRate28d),
			severity: 1 - summary.AttendanceRate28d,
		},
		{
			key:      "average_total_score_28d",
			label:    "Điểm trung bình 28 ngày",
			value:    summary.AverageTotalScore28d,
			display:  formatDecimal(summary.AverageTotalScore28d),
			severity: maxSeverity((5-summary.AverageTotalScore28d)/5, (10-summary.AverageTotalScore28d)/10),
		},
		{
			key:      "homework_completion_rate_28d",
			label:    "Tỷ lệ hoàn thành bài tập",
			value:    summary.HomeworkCompletionRate28d,
			display:  formatPercent(summary.HomeworkCompletionRate28d),
			severity: 1 - summary.HomeworkCompletionRate28d,
		},
		{
			key:      "absence_count_28d",
			label:    "Số buổi vắng 28 ngày",
			value:    summary.AbsenceCount28d,
			display:  formatDecimal(summary.AbsenceCount28d),
			severity: clamp(summary.AbsenceCount28d / 4),
		},
		{
			key:      "approved_leave_count_28d",
			label:    "Số đơn nghỉ được duyệt",
			value:    summary.ApprovedLeaveCount28d,
			display:  formatDecimal(summary.ApprovedLeaveCount28d),
			severity: clamp(summary.ApprovedLeaveCount28d / 3),
		},
		{
			key:      "days_since_last_lesson",
			label:    "Số ngày từ buổi học gần nhất",
			value:    summary.DaysSinceLastLesson,
			display:  formatDecimal(summary.DaysSinceLastLesson),
			severity: clamp(summary.DaysSinceLastLesson / 14),
		},
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].severity != candidates[j].severity {
			return candidates[i].severity > candidates[j].severity
		}
		return candidates[i].label < candidates[j].label
	})

	top := make([]StudentAtRiskTopFeature, 0, 3)
	for _, item := range candidates {
		top = append(top, StudentAtRiskTopFeature{
			Key:          item.key,
			Label:        item.label,
			Value:        item.value,
			DisplayValue: item.display,
		})
		if len(top) == 3 {
			break
		}
	}

	return top
}

func formatPercent(value float64) string {
	return formatDecimal(value*100) + "%"
}

func formatDecimal(value float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", value), "0"), ".")
}

func clamp(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func maxSeverity(a, b float64) float64 {
	if a > b {
		return clamp(a)
	}
	return clamp(b)
}
