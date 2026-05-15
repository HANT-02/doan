package scheduling

import (
	"testing"
	"time"

	"doan/internal/entities"
)

func TestResolveTravelMinutes(t *testing.T) {
	t.Parallel()

	campusA := "campus-a"
	campusB := "campus-b"

	roomA := &entities.Room{ID: "room-a", CampusID: &campusA}
	roomB := &entities.Room{ID: "room-b", CampusID: &campusB}

	travelMap := BuildCampusTravelTimeMap([]entities.CampusTravelTime{
		{FromCampusID: campusA, ToCampusID: campusB, TravelMinutes: 25, IsActive: true},
	})

	if minutes := ResolveTravelMinutes(roomA, roomB, travelMap); minutes != 25 {
		t.Fatalf("expected travel time 25, got %d", minutes)
	}
}

func TestResolveTravelMinutes_DefaultsForUnknownInterCampusTrip(t *testing.T) {
	t.Parallel()

	campusA := "campus-a"
	campusB := "campus-b"

	roomA := &entities.Room{ID: "room-a", CampusID: &campusA}
	roomB := &entities.Room{ID: "room-b", CampusID: &campusB}

	if minutes := ResolveTravelMinutes(roomA, roomB, map[string]int{}); minutes != DefaultInterCampusTravelMinutes {
		t.Fatalf("expected fallback travel time %d, got %d", DefaultInterCampusTravelMinutes, minutes)
	}
}

func TestHasSufficientTravelGap(t *testing.T) {
	t.Parallel()

	campusA := "campus-a"
	campusB := "campus-b"

	roomA := &entities.Room{ID: "room-a", CampusID: &campusA}
	roomB := &entities.Room{ID: "room-b", CampusID: &campusB}

	travelMap := BuildCampusTravelTimeMap([]entities.CampusTravelTime{
		{FromCampusID: campusA, ToCampusID: campusB, TravelMinutes: 20, IsActive: true},
	})

	previousEnd := time.Date(2026, 5, 14, 9, 0, 0, 0, time.UTC)
	nextStartTight := time.Date(2026, 5, 14, 9, 15, 0, 0, time.UTC)
	nextStartSafe := time.Date(2026, 5, 14, 9, 30, 0, 0, time.UTC)

	if HasSufficientTravelGap(previousEnd, nextStartTight, roomA, roomB, travelMap) {
		t.Fatalf("expected tight gap to be infeasible")
	}
	if !HasSufficientTravelGap(previousEnd, nextStartSafe, roomA, roomB, travelMap) {
		t.Fatalf("expected 30-minute gap to be feasible")
	}
}
