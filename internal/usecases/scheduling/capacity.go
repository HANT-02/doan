package scheduling

func CalculateCapacityLimit(roomCapacity, maxStudents int) int {
	if roomCapacity == 0 && maxStudents == 0 {
		return 0
	}
	if roomCapacity == 0 {
		return maxStudents
	}
	if maxStudents == 0 {
		return roomCapacity
	}
	if roomCapacity < maxStudents {
		return roomCapacity
	}
	return maxStudents
}

func CalculateUtilization(studentCount, capacityLimit int) (availableSpots int, utilization float64) {
	if capacityLimit <= 0 {
		return 0, 0.0
	}

	availableSpots = capacityLimit - studentCount
	if availableSpots < 0 {
		availableSpots = 0
	}

	utilization = float64(studentCount) / float64(capacityLimit)
	return availableSpots, utilization
}

func ValidateMakeupCapacity(studentCount, capacityLimit, additionalStudents int) bool {
	if capacityLimit <= 0 {
		return false
	}
	return studentCount+additionalStudents <= capacityLimit
}
