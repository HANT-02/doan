package scheduling

import (
	"fmt"
	"strings"
)

func AvailableBenchmarkScenarioNames() []string {
	scenarios := BuildDefaultBenchmarkScenarios()
	names := make([]string, 0, len(scenarios))
	for _, scenario := range scenarios {
		names = append(names, scenario.Name)
	}
	return names
}

func SelectBenchmarkScenarios(selectedNames []string, iterations int) ([]BenchmarkScenario, error) {
	scenarios := BuildDefaultBenchmarkScenarios()
	if iterations > 0 {
		for index := range scenarios {
			scenarios[index].Iterations = iterations
		}
	}

	if len(selectedNames) == 0 {
		return scenarios, nil
	}

	selectedSet := make(map[string]struct{}, len(selectedNames))
	for _, name := range selectedNames {
		normalized := strings.TrimSpace(strings.ToLower(name))
		if normalized == "" || normalized == "tat_ca" || normalized == "all" {
			return scenarios, nil
		}
		selectedSet[normalized] = struct{}{}
	}

	filtered := make([]BenchmarkScenario, 0, len(selectedSet))
	for _, scenario := range scenarios {
		if _, ok := selectedSet[strings.ToLower(scenario.Name)]; ok {
			filtered = append(filtered, scenario)
			delete(selectedSet, strings.ToLower(scenario.Name))
		}
	}

	if len(selectedSet) > 0 {
		unknown := make([]string, 0, len(selectedSet))
		for name := range selectedSet {
			unknown = append(unknown, name)
		}
		return nil, fmt.Errorf("ten kich ban khong hop le: %s", strings.Join(unknown, ", "))
	}

	return filtered, nil
}
