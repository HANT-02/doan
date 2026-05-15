package skills

import (
	"reflect"
	"testing"
)

func TestNormalizeCodes(t *testing.T) {
	t.Parallel()

	values := []string{" ielts_8.0 ", "TESOL", "", "tesol", "Ielts_8.0"}
	got := NormalizeCodes(values)
	want := []string{"IELTS_8.0", "TESOL"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected normalized codes %v, got %v", want, got)
	}
}

func TestMissingRequiredCodes(t *testing.T) {
	t.Parallel()

	got := MissingRequiredCodes(
		[]string{"tesol", "ielts_8.0"},
		[]string{"TESOL", "CELTA", "IELTS_8.0"},
	)
	want := []string{"CELTA"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected missing required codes %v, got %v", want, got)
	}
}
