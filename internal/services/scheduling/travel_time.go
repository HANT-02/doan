package scheduling

import (
	"fmt"
	"time"

	"doan/internal/entities"
)

const (
	DefaultInterCampusTravelMinutes = 30
)

func BuildCampusTravelTimeMap(items []entities.CampusTravelTime) map[string]int {
	result := make(map[string]int, len(items))
	for _, item := range items {
		if item.FromCampusID == "" || item.ToCampusID == "" || !item.IsActive {
			continue
		}
		result[campusTravelPairKey(item.FromCampusID, item.ToCampusID)] = item.TravelMinutes
	}
	return result
}

func ResolveTravelMinutes(
	fromRoom *entities.Room,
	toRoom *entities.Room,
	travelMap map[string]int,
) int {
	if fromRoom == nil || toRoom == nil {
		return 0
	}
	if fromRoom.CampusID == nil || toRoom.CampusID == nil {
		return 0
	}
	if *fromRoom.CampusID == *toRoom.CampusID {
		return 0
	}

	key := campusTravelPairKey(*fromRoom.CampusID, *toRoom.CampusID)
	if minutes, ok := travelMap[key]; ok {
		return minutes
	}

	return DefaultInterCampusTravelMinutes
}

func HasSufficientTravelGap(
	previousEnd time.Time,
	nextStart time.Time,
	fromRoom *entities.Room,
	toRoom *entities.Room,
	travelMap map[string]int,
) bool {
	requiredMinutes := ResolveTravelMinutes(fromRoom, toRoom, travelMap)
	if requiredMinutes <= 0 {
		return true
	}

	gapMinutes := int(nextStart.Sub(previousEnd).Minutes())
	return gapMinutes >= requiredMinutes
}

func campusTravelPairKey(fromCampusID, toCampusID string) string {
	return fmt.Sprintf("%s->%s", fromCampusID, toCampusID)
}
