package predictive

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAtRiskService_LoadsPythonArtifacts(t *testing.T) {
	artifactsDir := filepath.Join("..", "..", "..", "ml", "at_risk_prediction", "artifacts")
	absoluteDir, err := filepath.Abs(artifactsDir)
	if err != nil {
		t.Fatalf("resolve artifacts dir: %v", err)
	}

	t.Setenv(predictiveArtifactsEnvKey, absoluteDir)

	service := NewAtRiskService(nil)

	metadata, err := service.GetModelMetadata(context.Background(), true)
	if err != nil {
		t.Fatalf("expected metadata, got error: %v", err)
	}
	if metadata.ModelName != "logistic_regression" {
		t.Fatalf("expected selected model logistic_regression, got %s", metadata.ModelName)
	}
	if metadata.TrainSize == 0 || metadata.TestSize == 0 {
		t.Fatalf("expected train/test size from artifact, got %+v", metadata)
	}

	output, err := service.ListStudentPredictions(context.Background(), ListStudentPredictionsInput{
		Page:  1,
		Limit: 100,
	})
	if err != nil {
		t.Fatalf("expected predictions, got error: %v", err)
	}
	if len(output.Items) == 0 {
		t.Fatalf("expected non-empty prediction items")
	}
	if output.Summary.AtRiskCount == 0 {
		t.Fatalf("expected at least one AT_RISK prediction")
	}
	if output.ModelMetadata.ModelName != "logistic_regression" {
		t.Fatalf("expected metadata model logistic_regression in list output, got %s", output.ModelMetadata.ModelName)
	}
}

func TestAtRiskService_MissingArtifactsReturnsClearError(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv(predictiveArtifactsEnvKey, tempDir)

	service := NewAtRiskService(nil)
	_, err := service.ListStudentPredictions(context.Background(), ListStudentPredictionsInput{
		Page:  1,
		Limit: 10,
	})
	if err == nil {
		t.Fatalf("expected error when artifacts are missing")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "khong co du lieu prediction") {
		t.Fatalf("expected clear missing-artifact message, got %v", err)
	}

	if _, statErr := os.Stat(filepath.Join(tempDir, "reports", "latest_predictions.json")); !os.IsNotExist(statErr) && statErr != nil {
		t.Fatalf("unexpected stat error: %v", statErr)
	}
}
