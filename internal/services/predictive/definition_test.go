package predictive

import "testing"

func TestDefaultAtRiskDatasetDefinition_CoversTaskF1Scope(t *testing.T) {
	t.Parallel()

	definition := DefaultAtRiskDatasetDefinition()

	if definition.ProblemType != "classification" {
		t.Fatalf("expected classification problem type, got %s", definition.ProblemType)
	}

	if definition.Label.PositiveClass != LabelAtRisk || definition.Label.NegativeClass != LabelNotAtRisk {
		t.Fatalf("unexpected label classes: %+v", definition.Label)
	}

	requiredSources := map[string]bool{
		SourceStudent:        false,
		SourceAttendance:     false,
		SourceAcademicRecord: false,
		SourceEnrollment:     false,
	}

	for _, source := range definition.Sources {
		if _, ok := requiredSources[source.Key]; ok && source.Required {
			requiredSources[source.Key] = true
		}
	}

	for key, present := range requiredSources {
		if !present {
			t.Fatalf("expected required source %s in dataset definition", key)
		}
	}

	if len(definition.Features) < 10 {
		t.Fatalf("expected richer initial feature set, got %d features", len(definition.Features))
	}

	if definition.ObservationWindowDays != 28 || definition.PredictionHorizonDays != 28 {
		t.Fatalf("unexpected windows: observation=%d horizon=%d", definition.ObservationWindowDays, definition.PredictionHorizonDays)
	}
}
