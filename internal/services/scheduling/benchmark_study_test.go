package scheduling

import (
	"context"
	"strings"
	"testing"
)

func TestBuildDefaultBenchmarkScenarios_ReturnsExpectedScales(t *testing.T) {
	t.Parallel()

	scenarios := BuildDefaultBenchmarkScenarios()
	if len(scenarios) != 3 {
		t.Fatalf("expected 3 scenarios, got %d", len(scenarios))
	}

	expectedNames := []string{"small", "medium", "large"}
	expectedClasses := []int{6, 10, 16}
	expectedSessions := []int{2, 3, 4}

	for index, scenario := range scenarios {
		if scenario.Name != expectedNames[index] {
			t.Fatalf("expected scenario %d to be %s, got %s", index, expectedNames[index], scenario.Name)
		}
		if scenario.Iterations != 7 {
			t.Fatalf("expected default iterations to be 7, got %d for %s", scenario.Iterations, scenario.Name)
		}
		if len(scenario.Input.Classes) != expectedClasses[index] {
			t.Fatalf("expected %d classes for %s, got %d", expectedClasses[index], scenario.Name, len(scenario.Input.Classes))
		}
		for _, classEntity := range scenario.Input.Classes {
			if classEntity.Course.SessionCount != expectedSessions[index] {
				t.Fatalf(
					"expected session count %d for scenario %s, got %d",
					expectedSessions[index],
					scenario.Name,
					classEntity.Course.SessionCount,
				)
			}
		}
	}
}

func TestRunBenchmarkStudy_DefaultScenariosProducesComparableReport(t *testing.T) {
	t.Parallel()

	scenarios := BuildDefaultBenchmarkScenarios()
	for index := range scenarios {
		scenarios[index].Iterations = 2
	}

	catalog := NewSchedulingSolverCatalog(
		NewGraphColoringSolver(),
		NewCPSATSolver(),
		NewTabuSearchSolver(),
	)

	study, err := RunBenchmarkStudy(context.Background(), catalog, scenarios)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if study == nil {
		t.Fatalf("expected benchmark study")
	}

	if len(study.ScenarioReports) != 3 {
		t.Fatalf("expected 3 scenario reports, got %d", len(study.ScenarioReports))
	}

	if study.Recommendation.SelectedSolverKey == "" {
		t.Fatalf("expected selected solver key")
	}
	if len(study.Recommendation.Rationale) == 0 {
		t.Fatalf("expected recommendation rationale")
	}

	for _, scenario := range study.ScenarioReports {
		if scenario.Iterations != 2 {
			t.Fatalf("expected trimmed iteration count 2 for %s, got %d", scenario.Name, scenario.Iterations)
		}
		if len(scenario.Solvers) != 3 {
			t.Fatalf("expected 3 solver reports for %s, got %d", scenario.Name, len(scenario.Solvers))
		}

		for _, solver := range scenario.Solvers {
			if solver.Runs != 2 {
				t.Fatalf("expected 2 runs for solver %s in %s, got %d", solver.Key, scenario.Name, solver.Runs)
			}
			if solver.RepresentativeStatus == "" {
				t.Fatalf("expected representative status for solver %s in %s", solver.Key, scenario.Name)
			}
			if solver.AvgFeasibilityRate <= 0 {
				t.Fatalf("expected feasibility > 0 for solver %s in %s, got %f", solver.Key, scenario.Name, solver.AvgFeasibilityRate)
			}
			if solver.MinRuntimeMs > solver.MaxRuntimeMs {
				t.Fatalf("expected runtime bounds ordered for solver %s in %s", solver.Key, scenario.Name)
			}
		}
	}

	report := RenderBenchmarkStudyMarkdown(study)
	for _, expectedSection := range []string{
		"# Scheduling Benchmark Study",
		"## Scenario `small`",
		"## Scenario `medium`",
		"## Scenario `large`",
		"## Recommendation",
	} {
		if !strings.Contains(report, expectedSection) {
			t.Fatalf("expected markdown report to contain %q", expectedSection)
		}
	}
}
