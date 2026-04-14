package predictive

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"doan/internal/entities"

	"gorm.io/gorm"
)

const (
	DefaultTrainingSeed = int64(42)
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

type CurrentPredictionCandidate struct {
	Student    entities.Student
	Enrollment entities.Enrollment
	SnapshotAt time.Time
	Row        TrainingRow
}

type trainedModelSnapshot struct {
	Metadata      AtRiskModelMetadata
	LogisticModel *LogisticRegressionModel
	Recommended   string
}

type AtRiskService interface {
	ListStudentPredictions(ctx context.Context, input ListStudentPredictionsInput) (*ListStudentPredictionsOutput, error)
	GetModelMetadata(ctx context.Context, refresh bool) (*AtRiskModelMetadata, error)
}

type atRiskService struct {
	db         *gorm.DB
	definition DatasetDefinition
	seed       int64

	mu     sync.RWMutex
	cached *trainedModelSnapshot
}

func NewAtRiskService(db *gorm.DB) AtRiskService {
	return &atRiskService{
		db:         db,
		definition: DefaultAtRiskDatasetDefinition(),
		seed:       DefaultTrainingSeed,
	}
}

func (s *atRiskService) GetModelMetadata(ctx context.Context, refresh bool) (*AtRiskModelMetadata, error) {
	snapshot, err := s.getOrTrainSnapshot(ctx, refresh)
	if err != nil {
		return nil, err
	}

	metadata := snapshot.Metadata
	return &metadata, nil
}

func (s *atRiskService) ListStudentPredictions(ctx context.Context, input ListStudentPredictionsInput) (*ListStudentPredictionsOutput, error) {
	snapshot, err := s.getOrTrainSnapshot(ctx, input.Refresh)
	if err != nil {
		return nil, err
	}

	data, err := LoadTrainingSourceDataFromDB(ctx, s.db)
	if err != nil {
		return nil, err
	}

	candidates, err := BuildCurrentPredictionCandidates(s.definition, data, resolveCurrentSnapshotAt(data.Lessons))
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, errors.New("khong co du lieu prediction hien tai; hay kiem tra enrollment active va lesson trong 28 ngay gan nhat")
	}

	predictions := make([]StudentRiskPrediction, 0, len(candidates))
	for _, candidate := range candidates {
		predictions = append(predictions, s.predictCandidate(candidate, snapshot))
	}

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

func (s *atRiskService) getOrTrainSnapshot(ctx context.Context, refresh bool) (*trainedModelSnapshot, error) {
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

	rows, err := LoadTrainingRowsFromDB(ctx, s.db, s.definition)
	if err != nil {
		return nil, err
	}

	datasetName := fmt.Sprintf("db_%s", s.definition.Name)
	artifacts, err := TrainMinimalModels(datasetName, rows, s.seed)
	if err != nil {
		return nil, err
	}

	trainedAt := time.Now()
	metadata := AtRiskModelMetadata{
		Version:        fmt.Sprintf("%s-%s-%s", artifacts.Report.Recommended, s.definition.Version, trainedAt.Format("20060102150405")),
		ModelName:      artifacts.Report.Recommended,
		DatasetName:    artifacts.Report.DatasetName,
		DatasetSource:  "database",
		DefinitionName: s.definition.Name,
		DefinitionVer:  s.definition.Version,
		TrainedAt:      trainedAt,
		TrainSize:      artifacts.Report.TrainSize,
		TestSize:       artifacts.Report.TestSize,
		FeatureNames:   append([]string(nil), artifacts.Report.FeatureNames...),
		ModelReports:   append([]MinimalModelReport(nil), artifacts.Report.ModelReports...),
		Recommendation: append([]string(nil), artifacts.Report.Recommendation...),
	}

	s.cached = &trainedModelSnapshot{
		Metadata:      metadata,
		LogisticModel: artifacts.LogisticModel,
		Recommended:   artifacts.Report.Recommended,
	}

	return s.cached, nil
}

func (s *atRiskService) predictCandidate(candidate CurrentPredictionCandidate, snapshot *trainedModelSnapshot) StudentRiskPrediction {
	label := LabelNotAtRisk
	score := 0.0
	if snapshot.Recommended == "logistic_regression" && snapshot.LogisticModel != nil {
		score = snapshot.LogisticModel.PredictProbability(candidate.Row)
		if score >= 0.5 {
			label = LabelAtRisk
		}
	} else {
		label = PredictRuleBasedLabel(candidate.Row)
		score = PredictRuleBasedScore(candidate.Row)
	}

	reasons := buildPredictionReasons(candidate.Row, snapshot.LogisticModel, snapshot.Recommended == "logistic_regression")
	primaryReason := "Nguy cơ hiện tại thấp, chưa có tín hiệu nổi bật."
	if len(reasons) > 0 {
		primaryReason = reasons[0]
	}

	return StudentRiskPrediction{
		StudentID:     candidate.Student.ID,
		StudentCode:   candidate.Student.Code,
		StudentName:   candidate.Student.FullName,
		GradeLevel:    candidate.Student.GradeLevel,
		ClassID:       candidate.Enrollment.ClassID,
		ClassCode:     candidate.Enrollment.Class.Code,
		ClassName:     candidate.Enrollment.Class.Name,
		SnapshotAt:    candidate.SnapshotAt,
		Label:         label,
		RiskScore:     score,
		RiskBand:      classifyRiskBand(score),
		ModelName:     snapshot.Recommended,
		ModelVersion:  snapshot.Metadata.Version,
		PrimaryReason: primaryReason,
		Reasons:       reasons,
		FeatureSummary: AtRiskPredictionFeatureSnapshot{
			AttendanceRate28d:         candidate.Row.Features["attendance_rate_28d"],
			AbsenceCount28d:           candidate.Row.Features["absence_count_28d"],
			AverageTotalScore28d:      candidate.Row.Features["average_total_score_28d"],
			HomeworkCompletionRate28d: candidate.Row.Features["homework_completion_rate_28d"],
			ActiveEnrollmentCount28d:  candidate.Row.Features["active_enrollment_count_28d"],
			WeeklyLessonLoad28d:       candidate.Row.Features["weekly_lesson_load_28d"],
			ApprovedLeaveCount28d:     candidate.Row.Features["approved_leave_count_28d"],
			DaysSinceLastLesson:       candidate.Row.Features["days_since_last_lesson"],
		},
	}
}

func BuildCurrentPredictionCandidates(definition DatasetDefinition, data TrainingSourceData, snapshotAt time.Time) ([]CurrentPredictionCandidate, error) {
	if snapshotAt.IsZero() {
		snapshotAt = time.Now()
	}

	studentsByID := make(map[string]entities.Student, len(data.Students))
	for _, student := range data.Students {
		studentsByID[student.ID] = student
	}

	enrollmentsByStudent := make(map[string][]entities.Enrollment)
	for _, enrollment := range data.Enrollments {
		if strings.EqualFold(enrollment.Status, "ENROLLED") {
			enrollmentsByStudent[enrollment.StudentID] = append(enrollmentsByStudent[enrollment.StudentID], enrollment)
		}
	}

	lessonsByClass := make(map[string][]entities.Lesson)
	lessonByID := make(map[string]entities.Lesson, len(data.Lessons))
	for _, lesson := range data.Lessons {
		lessonsByClass[lesson.ClassID] = append(lessonsByClass[lesson.ClassID], lesson)
		lessonByID[lesson.ID] = lesson
	}

	attendanceByStudentClass := make(map[string][]attendanceObservation)
	for _, attendance := range data.Attendances {
		lesson, ok := lessonByID[attendance.LessonID]
		if !ok {
			continue
		}
		key := studentClassKey(attendance.StudentID, lesson.ClassID)
		attendanceByStudentClass[key] = append(attendanceByStudentClass[key], attendanceObservation{
			At:     lesson.DateStart,
			Status: attendance.Status,
		})
	}

	summaryByID := make(map[string]entities.LessonSummary, len(data.LessonSummaries))
	for _, summary := range data.LessonSummaries {
		summaryByID[summary.ID] = summary
	}

	academicByStudentClass := make(map[string][]academicObservation)
	for _, record := range data.AcademicRecords {
		summary, ok := summaryByID[record.LessonSummaryID]
		if !ok {
			continue
		}
		lesson, ok := lessonByID[summary.LessonID]
		if !ok {
			continue
		}
		key := studentClassKey(record.StudentID, lesson.ClassID)
		academicByStudentClass[key] = append(academicByStudentClass[key], academicObservation{
			At:                lesson.DateStart,
			HomeworkCompleted: record.HomeworkCompleted,
			TotalScore:        record.TotalScore,
			IsCompleted:       record.IsCompleted,
		})
	}

	leaveByStudent := make(map[string][]leaveObservation)
	for _, leaveRequest := range data.LeaveRequests {
		if !strings.EqualFold(leaveRequest.Status, "APPROVED") {
			continue
		}
		classID := ""
		if leaveRequest.ClassID != nil {
			classID = *leaveRequest.ClassID
		}
		leaveByStudent[leaveRequest.StudentID] = append(leaveByStudent[leaveRequest.StudentID], leaveObservation{
			At:      leaveRequest.ApplyDate,
			ClassID: classID,
		})
	}

	observationStart := snapshotAt.AddDate(0, 0, -definition.ObservationWindowDays)
	candidates := make([]CurrentPredictionCandidate, 0)
	seen := make(map[string]struct{})
	for studentID, enrollments := range enrollmentsByStudent {
		student, ok := studentsByID[studentID]
		if !ok || !strings.EqualFold(student.Status, "ACTIVE") {
			continue
		}

		activeClassIDs := activeEnrollmentClassIDs(enrollments, snapshotAt)
		for _, enrollment := range enrollments {
			if !isEnrollmentActiveAt(enrollment, snapshotAt) {
				continue
			}
			if _, ok := seen[studentClassKey(studentID, enrollment.ClassID)]; ok {
				continue
			}

			classKey := studentClassKey(studentID, enrollment.ClassID)
			row := TrainingRow{
				ID: makeTrainingRowID(studentID, enrollment.ClassID, snapshotAt),
				Features: map[string]float64{
					"attendance_rate_28d":          observationAttendanceRate(filterAttendances(attendanceByStudentClass[classKey], observationStart, snapshotAt)),
					"absence_count_28d":            float64(countAbsent(filterAttendances(attendanceByStudentClass[classKey], observationStart, snapshotAt))),
					"average_total_score_28d":      averageCompletedScore(completedAcademicRecords(filterAcademicRecords(academicByStudentClass[classKey], observationStart, snapshotAt))),
					"homework_completion_rate_28d": homeworkCompletionRate(completedAcademicRecords(filterAcademicRecords(academicByStudentClass[classKey], observationStart, snapshotAt))),
					"active_enrollment_count_28d":  float64(len(activeClassIDs)),
					"weekly_lesson_load_28d":       weeklyLessonLoad(lessonsByClass, activeClassIDs, observationStart, snapshotAt, definition.ObservationWindowDays),
					"approved_leave_count_28d":     float64(countApprovedLeaves(leaveByStudent[studentID], observationStart, snapshotAt, enrollment.ClassID)),
					"days_since_last_lesson":       daysSinceLastLesson(lessonsByClass, activeClassIDs, snapshotAt, definition.ObservationWindowDays),
				},
			}
			candidates = append(candidates, CurrentPredictionCandidate{
				Student:    student,
				Enrollment: enrollment,
				SnapshotAt: snapshotAt,
				Row:        row,
			})
			seen[classKey] = struct{}{}
		}
	}

	if len(candidates) == 0 {
		return nil, errors.New("khong tim thay candidate prediction nao tu database")
	}

	return candidates, nil
}

func PredictRuleBasedScore(row TrainingRow) float64 {
	deficits := []float64{
		maxFloat(0, (0.80-row.Features["attendance_rate_28d"])/0.80),
		maxFloat(0, (5.0-row.Features["average_total_score_28d"])/5.0),
		maxFloat(0, (0.60-row.Features["homework_completion_rate_28d"])/0.60),
	}

	score := 0.15 * clamp(row.Features["approved_leave_count_28d"]/4.0, 0, 1)
	score += 0.10 * clamp(row.Features["days_since_last_lesson"]/14.0, 0, 1)
	for _, deficit := range deficits {
		score += 0.25 * clamp(deficit, 0, 1)
	}

	return clamp(score, 0, 1)
}

func buildPredictionReasons(row TrainingRow, logisticModel *LogisticRegressionModel, useLogistic bool) []string {
	reasons := make([]string, 0, 4)
	if row.Features["attendance_rate_28d"] < 0.80 {
		reasons = append(reasons, fmt.Sprintf("Tỷ lệ đi học 28 ngày gần nhất chỉ %.0f%%, thấp hơn ngưỡng an toàn 80%%.", row.Features["attendance_rate_28d"]*100))
	}
	if row.Features["average_total_score_28d"] < 5.0 {
		reasons = append(reasons, fmt.Sprintf("Điểm trung bình 28 ngày gần nhất là %.2f, dưới ngưỡng 5.00.", row.Features["average_total_score_28d"]))
	}
	if row.Features["homework_completion_rate_28d"] < 0.60 {
		reasons = append(reasons, fmt.Sprintf("Tỷ lệ hoàn thành bài tập chỉ %.0f%% trong 28 ngày gần nhất.", row.Features["homework_completion_rate_28d"]*100))
	}
	if row.Features["approved_leave_count_28d"] >= 2 {
		reasons = append(reasons, fmt.Sprintf("Có %.0f đơn nghỉ đã duyệt trong 28 ngày gần nhất.", row.Features["approved_leave_count_28d"]))
	}
	if row.Features["days_since_last_lesson"] >= 10 {
		reasons = append(reasons, fmt.Sprintf("Khoảng cách từ buổi học gần nhất đã %.0f ngày.", row.Features["days_since_last_lesson"]))
	}

	if useLogistic && logisticModel != nil {
		reasons = append(reasons, logisticContributionReasons(row, logisticModel)...)
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "Các chỉ số chuyên cần và học tập gần đây đang trong vùng ổn định.")
	}

	if len(reasons) > 3 {
		return reasons[:3]
	}
	return reasons
}

func logisticContributionReasons(row TrainingRow, model *LogisticRegressionModel) []string {
	type contribution struct {
		feature string
		value   float64
	}

	contributions := make([]contribution, 0, len(model.FeatureNames))
	for idx, featureName := range model.FeatureNames {
		scaled := (row.Features[featureName] - model.Means[idx]) / model.Scales[idx]
		value := model.Weights[idx] * scaled
		if value > 0 {
			contributions = append(contributions, contribution{feature: featureName, value: value})
		}
	}

	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].value > contributions[j].value
	})

	reasons := make([]string, 0, 2)
	for _, item := range contributions {
		switch item.feature {
		case "absence_count_28d":
			reasons = append(reasons, "Mô hình ghi nhận số buổi vắng gần đây là tín hiệu rủi ro tăng.")
		case "approved_leave_count_28d":
			reasons = append(reasons, "Mô hình coi tần suất xin nghỉ gần đây là tín hiệu phụ làm tăng nguy cơ.")
		case "days_since_last_lesson":
			reasons = append(reasons, "Mô hình coi khoảng cách dài từ buổi học gần nhất là dấu hiệu cần theo dõi.")
		}
		if len(reasons) == 2 {
			break
		}
	}
	return reasons
}

func classifyRiskBand(score float64) string {
	switch {
	case score >= 0.75:
		return "CAO"
	case score >= 0.50:
		return "TRUNG_BINH"
	default:
		return "THAP"
	}
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

func resolveCurrentSnapshotAt(lessons []entities.Lesson) time.Time {
	now := time.Now()
	if len(lessons) == 0 {
		return now
	}

	latest := lessons[0].DateStart
	for _, lesson := range lessons[1:] {
		if lesson.DateStart.After(latest) {
			latest = lesson.DateStart
		}
	}

	if latest.Before(now.AddDate(0, 0, -7)) {
		return latest.Add(24 * time.Hour)
	}
	return now
}

func maxFloat(left, right float64) float64 {
	if left > right {
		return left
	}
	return right
}
