package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	schedulingservice "doan/internal/services/scheduling"
)

func main() {
	var (
		scenarioFlag   string
		iterationsFlag int
		outputDirFlag  string
		listOnlyFlag   bool
	)

	flag.StringVar(&scenarioFlag, "kich-ban", "tat_ca", "Danh sach kich ban can chay, cach nhau boi dau phay: nho,trung_binh,lon")
	flag.IntVar(&iterationsFlag, "so-lan-chay", 0, "Neu > 0 thi ghi de so lan chay mac dinh cua moi kich ban")
	flag.StringVar(&outputDirFlag, "thu-muc-dau-ra", "", "Thu muc goc de ghi tap tin minh chung")
	flag.BoolVar(&listOnlyFlag, "liet-ke-kich-ban", false, "Chi in danh sach kich ban co san roi thoat")
	flag.Parse()

	if listOnlyFlag {
		fmt.Println(strings.Join(schedulingservice.AvailableBenchmarkScenarioNames(), "\n"))
		return
	}

	selectedNames := strings.Split(scenarioFlag, ",")
	scenarios, err := schedulingservice.SelectBenchmarkScenarios(selectedNames, iterationsFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "khong the chon kich ban benchmark: %v\n", err)
		os.Exit(1)
	}

	catalog := schedulingservice.NewSchedulingSolverCatalog(
		schedulingservice.NewGraphColoringSolver(),
		schedulingservice.NewCPSATSolver(),
		schedulingservice.NewTabuSearchSolver(),
	)

	study, err := schedulingservice.RunBenchmarkStudy(context.Background(), catalog, scenarios)
	if err != nil {
		fmt.Fprintf(os.Stderr, "khong the chay bo do dac doi sanh xep lich: %v\n", err)
		os.Exit(1)
	}

	artifactDir, err := schedulingservice.WriteBenchmarkArtifacts(outputDirFlag, study, scenarios)
	if err != nil {
		fmt.Fprintf(os.Stderr, "khong the ghi tap tin minh chung benchmark: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "da ghi tap tin minh chung benchmark vao: %s\n", artifactDir)
	fmt.Print(schedulingservice.RenderBenchmarkStudyMarkdown(study))
}
