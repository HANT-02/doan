package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"doan/internal/infrastructure/database/postgres"
	"doan/internal/services/predictive"
	"doan/pkg/config"
	"doan/pkg/constants"
	"doan/pkg/logger"
)

func main() {
	var (
		source         = flag.String("source", "demo", "Dataset source: demo | csv | db")
		datasetPath    = flag.String("dataset", "", "CSV dataset path when source=csv")
		datasetName    = flag.String("name", "", "Dataset name for the report")
		seed           = flag.Int64("seed", 42, "Random seed for split/demo generation")
		size           = flag.Int("size", 120, "Demo dataset size when source=demo")
		configFilePath = flag.String("config-file-path", "./configs", "Config file path when source=db")
		configFile     = flag.String("config-file", "config", "Config file name when source=db")
	)
	flag.Parse()

	var (
		rows []predictive.TrainingRow
		err  error
	)

	sourceValue := strings.ToLower(strings.TrimSpace(*source))
	switch sourceValue {
	case "csv":
		if strings.TrimSpace(*datasetPath) == "" {
			fmt.Fprintln(os.Stderr, "source=csv requires --dataset")
			os.Exit(1)
		}
		rows, err = predictive.LoadTrainingRowsFromCSV(*datasetPath, predictive.MinimalFeatureNames)
	case "db":
		rows, err = loadTrainingRowsFromDB(*configFilePath, *configFile)
	case "demo":
		rows = predictive.BuildDemoAtRiskDataset(*size, *seed)
	default:
		fmt.Fprintf(os.Stderr, "unsupported source %q, expected demo | csv | db\n", *source)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load predictive dataset: %v\n", err)
		os.Exit(1)
	}

	report, err := predictive.RunMinimalTrainingComparison(resolveDatasetName(sourceValue, *datasetName, *datasetPath), rows, *seed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run minimal predictive training: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(predictive.RenderMinimalTrainingReportMarkdown(report))
}

func resolveDatasetName(source, datasetName, datasetPath string) string {
	if strings.TrimSpace(datasetName) != "" {
		return datasetName
	}

	switch source {
	case "csv":
		return strings.TrimSuffix(filepath.Base(datasetPath), filepath.Ext(datasetPath))
	case "db":
		return "db_at_risk_dataset"
	default:
		return "demo_at_risk_dataset"
	}
}

func loadTrainingRowsFromDB(configFilePath, configFile string) ([]predictive.TrainingRow, error) {
	configSource := &config.Viper{
		ConfigType: constants.ConfigTypeFile,
		FilePath:   configFilePath,
		ConfigFile: configFile,
	}
	if err := configSource.InitConfig(); err != nil {
		return nil, err
	}

	log := logger.NewZapLogger(logger.Config{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		ServiceName: "predictive-train",
		Environment: "development",
	})
	ctx := context.Background()
	db, err := postgres.GetDBContext(ctx, log, config.GetManager())
	if err != nil {
		return nil, err
	}

	return predictive.LoadTrainingRowsFromDB(ctx, db, predictive.DefaultAtRiskDatasetDefinition())
}
