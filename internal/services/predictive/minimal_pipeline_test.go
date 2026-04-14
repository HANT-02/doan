package predictive

import "testing"

func TestRunMinimalTrainingComparison_LogisticIsUsableOnDemoDataset(t *testing.T) {
	t.Parallel()

	rows := BuildDemoAtRiskDataset(120, 42)
	report, err := RunMinimalTrainingComparison("demo", rows, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report == nil {
		t.Fatalf("expected report")
	}

	if report.TrainSize == 0 || report.TestSize == 0 {
		t.Fatalf("expected non-empty split: %+v", report)
	}

	var logisticF1 float64
	foundLogistic := false
	for _, item := range report.ModelReports {
		if item.Name == "logistic_regression" {
			logisticF1 = item.Metrics.F1Score
			foundLogistic = true
		}
	}

	if !foundLogistic {
		t.Fatalf("expected logistic regression report")
	}

	if logisticF1 < 0.75 {
		t.Fatalf("expected usable logistic F1 on demo dataset, got %.3f", logisticF1)
	}
}

func TestLoadTrainingRowsFromCSV_RejectsMissingColumns(t *testing.T) {
	t.Parallel()

	_, err := LoadTrainingRowsFromCSV("/path/that/does/not/exist.csv", MinimalFeatureNames)
	if err == nil {
		t.Fatalf("expected error for missing csv file")
	}
}
