package predictive

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	predictiveservice "doan/internal/services/predictive"
	"doan/pkg/logger"
)

const (
	defaultPredictiveMLRelativePath = "ml/at_risk_prediction"
	predictiveMLDirEnvKey           = "PREDICTIVE_ML_DIR"
	predictivePythonBinEnvKey       = "PREDICTIVE_PYTHON_BIN"
	predictiveTrainTimeoutEnvKey    = "PREDICTIVE_TRAIN_TIMEOUT_SECONDS"
	defaultPredictivePythonBin      = "python3"
	defaultPredictiveTrainTimeout   = 10 * time.Minute
	commandOutputTailLimit          = 6000
)

type TrainAtRiskFromDBStep struct {
	Name       string `json:"name"`
	Command    string `json:"command"`
	Status     string `json:"status"`
	DurationMs int64  `json:"duration_ms"`
	OutputTail string `json:"output_tail,omitempty"`
}

type TrainAtRiskFromDBOutput struct {
	Message       string                                `json:"message"`
	DatasetName   string                                `json:"dataset_name"`
	StartedAt     time.Time                             `json:"started_at"`
	FinishedAt    time.Time                             `json:"finished_at"`
	DurationMs    int64                                 `json:"duration_ms"`
	MLDir         string                                `json:"ml_dir"`
	ModelMetadata predictiveservice.AtRiskModelMetadata `json:"model_metadata"`
	Steps         []TrainAtRiskFromDBStep               `json:"steps"`
}

type TrainAtRiskFromDBUseCase interface {
	Execute(ctx context.Context) (*TrainAtRiskFromDBOutput, error)
}

type trainAtRiskFromDBUseCase struct {
	atRiskService predictiveservice.AtRiskService
}

func NewTrainAtRiskFromDBUseCase(atRiskService predictiveservice.AtRiskService) TrainAtRiskFromDBUseCase {
	return &trainAtRiskFromDBUseCase{atRiskService: atRiskService}
}

func (uc *trainAtRiskFromDBUseCase) Execute(ctx context.Context) (*TrainAtRiskFromDBOutput, error) {
	ctxLogger := logger.NewLogger(ctx)
	startedAt := time.Now()

	mlDir, err := resolvePredictiveMLDir()
	if err != nil {
		ctxLogger.Errorf("Failed to resolve predictive ML directory: %v", err)
		return nil, err
	}

	timeout := predictiveTrainTimeout()
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pythonBin := predictivePythonBin()
	steps := make([]TrainAtRiskFromDBStep, 0, 2)
	trainWarning := ""

	trainStep, err := runPredictivePythonStep(runCtx, mlDir, pythonBin, "Lấy dữ liệu DB và huấn luyện", []string{
		"scripts/train_from_db.py",
		"--dataset-name",
		"at_risk_dataset_db",
	})
	steps = append(steps, trainStep)
	if err != nil {
		if !isInsufficientTrainingDatasetError(err, trainStep.OutputTail) {
			ctxLogger.Errorf("Failed to train at-risk model from DB: %v", err)
			return nil, buildTrainFromDBError("train model từ DB thất bại", err, steps)
		}
		trainWarning = "DB hiện chưa đủ dữ liệu có nhãn để huấn luyện lại, hệ thống giữ model hiện tại và chỉ chạy inference mới từ DB."
		steps[len(steps)-1].Status = "warning"
		ctxLogger.Warnf("Skip at-risk training because labelled dataset is insufficient: %v", err)
	}

	predictStep, err := runPredictivePythonStep(runCtx, mlDir, pythonBin, "Suy luận và ghi artifact dự báo", []string{
		"scripts/predict_from_db.py",
	})
	steps = append(steps, predictStep)
	if err != nil {
		ctxLogger.Errorf("Failed to generate at-risk predictions from DB: %v", err)
		return nil, buildTrainFromDBError("inference từ DB thất bại", err, steps)
	}

	predictions, err := uc.atRiskService.ListStudentPredictions(runCtx, predictiveservice.ListStudentPredictionsInput{
		Page:    1,
		Limit:   1,
		Refresh: true,
	})
	if err != nil {
		ctxLogger.Errorf("Failed to reload predictive artifacts after train: %v", err)
		return nil, buildTrainFromDBError("đã train xong nhưng không đọc được artifact mới", err, steps)
	}

	finishedAt := time.Now()
	message := "Đã train lại mô hình từ DB, chạy inference và tải lại artifact mới."
	if trainWarning != "" {
		message = trainWarning
	}
	return &TrainAtRiskFromDBOutput{
		Message:       message,
		DatasetName:   predictions.ModelMetadata.DatasetName,
		StartedAt:     startedAt,
		FinishedAt:    finishedAt,
		DurationMs:    finishedAt.Sub(startedAt).Milliseconds(),
		MLDir:         mlDir,
		ModelMetadata: predictions.ModelMetadata,
		Steps:         steps,
	}, nil
}

func isInsufficientTrainingDatasetError(err error, outputTail string) bool {
	text := strings.ToLower(err.Error() + "\n" + outputTail)
	markers := []string{
		"khong tao duoc dong dataset nao",
		"dataset rong",
		"dataset chi co mot nhan",
		"tap train chi co mot nhan",
	}
	for _, marker := range markers {
		if strings.Contains(text, marker) {
			return true
		}
	}
	return false
}

func runPredictivePythonStep(ctx context.Context, mlDir, pythonBin, name string, args []string) (TrainAtRiskFromDBStep, error) {
	startedAt := time.Now()
	command := strings.Join(append([]string{pythonBin}, args...), " ")
	step := TrainAtRiskFromDBStep{
		Name:    name,
		Command: command,
		Status:  "running",
	}

	cmd := exec.CommandContext(ctx, pythonBin, args...)
	cmd.Dir = mlDir
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	step.DurationMs = time.Since(startedAt).Milliseconds()
	step.OutputTail = tailUTF8(string(output), commandOutputTailLimit)
	if err != nil {
		step.Status = "failed"
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return step, fmt.Errorf("%s quá thời gian chờ sau %s", name, predictiveTrainTimeout())
		}
		return step, fmt.Errorf("%s lỗi: %w. Log cuối: %s", name, err, step.OutputTail)
	}

	step.Status = "success"
	return step, nil
}

func buildTrainFromDBError(message string, err error, steps []TrainAtRiskFromDBStep) error {
	var details []string
	for _, step := range steps {
		if step.Status != "failed" {
			continue
		}
		if strings.TrimSpace(step.OutputTail) != "" {
			details = append(details, fmt.Sprintf("%s: %s", step.Name, step.OutputTail))
		}
	}
	if len(details) == 0 {
		return fmt.Errorf("%s: %w", message, err)
	}
	return fmt.Errorf("%s: %w\n%s", message, err, strings.Join(details, "\n"))
}

func resolvePredictiveMLDir() (string, error) {
	if envPath := strings.TrimSpace(os.Getenv(predictiveMLDirEnvKey)); envPath != "" {
		if err := assertPredictiveMLDir(envPath); err != nil {
			return "", err
		}
		return envPath, nil
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := workingDir
	for {
		candidate := filepath.Join(current, defaultPredictiveMLRelativePath)
		if err := assertPredictiveMLDir(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("khong tim thay thu muc ML %s tu working dir %s", defaultPredictiveMLRelativePath, workingDir)
}

func assertPredictiveMLDir(path string) error {
	trainScript := filepath.Join(path, "scripts", "train_from_db.py")
	predictScript := filepath.Join(path, "scripts", "predict_from_db.py")
	if _, err := os.Stat(trainScript); err != nil {
		return fmt.Errorf("khong tim thay script train DB: %s", trainScript)
	}
	if _, err := os.Stat(predictScript); err != nil {
		return fmt.Errorf("khong tim thay script inference DB: %s", predictScript)
	}
	return nil
}

func predictivePythonBin() string {
	if value := strings.TrimSpace(os.Getenv(predictivePythonBinEnvKey)); value != "" {
		return value
	}
	return defaultPredictivePythonBin
}

func predictiveTrainTimeout() time.Duration {
	value := strings.TrimSpace(os.Getenv(predictiveTrainTimeoutEnvKey))
	if value == "" {
		return defaultPredictiveTrainTimeout
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return defaultPredictiveTrainTimeout
	}
	return time.Duration(seconds) * time.Second
}

func tailUTF8(value string, maxBytes int) string {
	if maxBytes <= 0 || len(value) <= maxBytes {
		return value
	}

	tail := value[len(value)-maxBytes:]
	for len(tail) > 0 && !utf8.ValidString(tail) {
		tail = tail[1:]
	}
	return "... " + tail
}
