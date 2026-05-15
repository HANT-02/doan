package scheduling

import (
	"testing"
)

func TestCalculateCapacityLimit(t *testing.T) {
	tests := []struct {
		roomCapacity int
		maxStudents  int
		expected     int
	}{
		{0, 0, 0},
		{0, 20, 20},
		{30, 0, 30},
		{20, 30, 20},
		{30, 20, 20},
		{25, 25, 25},
	}

	for _, tc := range tests {
		actual := CalculateCapacityLimit(tc.roomCapacity, tc.maxStudents)
		if actual != tc.expected {
			t.Errorf("CalculateCapacityLimit(%d, %d) = %d, expected %d", tc.roomCapacity, tc.maxStudents, actual, tc.expected)
		}
	}
}

func TestCalculateUtilization(t *testing.T) {
	tests := []struct {
		studentCount  int
		capacityLimit int
		expectedSpots int
		expectedUtil  float64
	}{
		{0, 0, 0, 0.0},
		{10, 20, 10, 0.5},
		{20, 20, 0, 1.0},
		{25, 20, 0, 1.25},
		{0, 20, 20, 0.0},
	}

	for _, tc := range tests {
		spots, util := CalculateUtilization(tc.studentCount, tc.capacityLimit)
		if spots != tc.expectedSpots || util != tc.expectedUtil {
			t.Errorf("CalculateUtilization(%d, %d) = (%d, %f), expected (%d, %f)", tc.studentCount, tc.capacityLimit, spots, util, tc.expectedSpots, tc.expectedUtil)
		}
	}
}

func TestValidateMakeupCapacity(t *testing.T) {
	tests := []struct {
		studentCount       int
		capacityLimit      int
		additionalStudents int
		expected           bool
	}{
		{0, 0, 1, false},
		{10, 20, 5, true},
		{15, 20, 5, true},
		{16, 20, 5, false},
		{20, 20, 1, false},
	}

	for _, tc := range tests {
		actual := ValidateMakeupCapacity(tc.studentCount, tc.capacityLimit, tc.additionalStudents)
		if actual != tc.expected {
			t.Errorf("ValidateMakeupCapacity(%d, %d, %d) = %t, expected %t", tc.studentCount, tc.capacityLimit, tc.additionalStudents, actual, tc.expected)
		}
	}
}
