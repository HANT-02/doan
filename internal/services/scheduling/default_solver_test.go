package scheduling

import "testing"

func TestNewDefaultSchedulingSolver_SelectsCPSAT(t *testing.T) {
	t.Parallel()

	defaultSolver := NewDefaultSchedulingSolver(NewCPSATSolver())
	if defaultSolver == nil {
		t.Fatalf("expected default solver")
	}

	if defaultSolver.Key() != SolverKeyCPSAT {
		t.Fatalf("expected default solver key %s, got %s", SolverKeyCPSAT, defaultSolver.Key())
	}
}
