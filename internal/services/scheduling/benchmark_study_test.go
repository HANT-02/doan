package scheduling

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDefaultBenchmarkScenarios_ReturnsExpectedScales(t *testing.T) {
	t.Parallel()

	scenarios := BuildDefaultBenchmarkScenarios()
	if len(scenarios) != 3 {
		t.Fatalf("expected 3 scenarios, got %d", len(scenarios))
	}

	expectedNames := []string{"nho", "trung_binh", "lon"}
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

func TestSelectBenchmarkScenarios_FiltersAndOverridesIterations(t *testing.T) {
	t.Parallel()

	scenarios, err := SelectBenchmarkScenarios([]string{"nho", "lon"}, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(scenarios) != 2 {
		t.Fatalf("expected 2 scenarios, got %d", len(scenarios))
	}

	if scenarios[0].Name != "nho" || scenarios[1].Name != "lon" {
		t.Fatalf("unexpected scenario order: %s, %s", scenarios[0].Name, scenarios[1].Name)
	}

	for _, scenario := range scenarios {
		if scenario.Iterations != 3 {
			t.Fatalf("expected iteration override 3 for %s, got %d", scenario.Name, scenario.Iterations)
		}
	}
}

func TestSelectBenchmarkScenarios_RejectsUnknownScenario(t *testing.T) {
	t.Parallel()

	_, err := SelectBenchmarkScenarios([]string{"khong_ton_tai"}, 0)
	if err == nil {
		t.Fatalf("expected error for unknown scenario")
	}
}

func TestWriteBenchmarkArtifacts_WritesExpectedFiles(t *testing.T) {
	t.Parallel()

	scenarios := BuildDefaultBenchmarkScenarios()
	for index := range scenarios {
		scenarios[index].Iterations = 1
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

	outputRoot := t.TempDir()
	artifactDir, err := WriteBenchmarkArtifacts(outputRoot, study, scenarios)
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	for _, relativePath := range []string{
		"tong_hop.json",
		"tong_hop.md",
		"bang_tong_hop_theo_thuat_toan.csv",
		"so_lieu_tho_tung_lan_chay.csv",
		"bang_ke_tap_tin.json",
		filepath.Join("kich_ban_nho", "du_lieu_dau_vao.json"),
		filepath.Join("kich_ban_nho", "thuat_toan_graph_coloring", "lan_chay_001", "so_lieu_tho.json"),
	} {
		fullPath := filepath.Join(artifactDir, relativePath)
		if _, err := os.Stat(fullPath); err != nil {
			t.Fatalf("expected artifact file %s: %v", fullPath, err)
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
			if len(solver.RunRecords) != 2 {
				t.Fatalf("expected 2 run records for solver %s in %s, got %d", solver.Key, scenario.Name, len(solver.RunRecords))
			}
			if solver.StdDevRuntimeMs < 0 {
				t.Fatalf("expected non-negative runtime stddev for solver %s in %s", solver.Key, scenario.Name)
			}
			for _, run := range solver.RunRecords {
				if run.StartedAt.IsZero() || run.FinishedAt.IsZero() {
					t.Fatalf("expected timestamps for solver %s in %s", solver.Key, scenario.Name)
				}
				if run.Telemetry == nil {
					t.Fatalf("expected telemetry for solver %s in %s", solver.Key, scenario.Name)
				}
			}
		}
	}

	report := RenderBenchmarkStudyMarkdown(study)
	for _, expectedSection := range []string{
		"# Báo Cáo Đo Đạc Đối Sánh Xếp Lịch",
		"## Kịch bản `nho`",
		"## Kịch bản `trung_binh`",
		"## Kịch bản `lon`",
		"## Khuyến Nghị",
	} {
		if !strings.Contains(report, expectedSection) {
			t.Fatalf("expected markdown report to contain %q", expectedSection)
		}
	}
}
