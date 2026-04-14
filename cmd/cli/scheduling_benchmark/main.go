package main

import (
	"context"
	"fmt"
	"os"

	schedulingservice "doan/internal/services/scheduling"
)

func main() {
	catalog := schedulingservice.NewSchedulingSolverCatalog(
		schedulingservice.NewGraphColoringSolver(),
		schedulingservice.NewCPSATSolver(),
		schedulingservice.NewTabuSearchSolver(),
	)

	study, err := schedulingservice.RunBenchmarkStudy(context.Background(), catalog, schedulingservice.BuildDefaultBenchmarkScenarios())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run scheduling benchmark study: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(schedulingservice.RenderBenchmarkStudyMarkdown(study))
}
