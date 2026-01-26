package utils

func IsInterfaceNumber(value interface{}) bool {
	switch value.(type) {
	case int64, float64:
		return true
	default:
		return false
	}
}

func IsInterfaceArrayNumber(value interface{}) bool {
	valueArray, ok := value.([]interface{})
	if !ok {
		return false
	}
	for _, item := range valueArray {
		if !IsInterfaceNumber(item) {
			return false
		}
	}
	return true
}
