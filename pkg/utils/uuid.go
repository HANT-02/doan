package utils

import (
	"github.com/google/uuid"
)

func GenerateUUID() string {
	// UUID v4
	newUUID := uuid.New()
	return newUUID.String()
}

func GenerateUUIDWithPrefix(prefix string) string {
	newUUID := uuid.New()
	return prefix + newUUID.String()
}

func IsValidUUID(uuidStr string) bool {
	parsedStr, err := uuid.Parse(uuidStr)
	if err != nil {
		return false
	}
	return parsedStr.String() == uuidStr
}

func IsInterfaceUUID(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}
	return IsValidUUID(str)
}

func IsInterfaceArrayUUID(value interface{}) bool {
	valueSlice, ok := value.([]string)
	if !ok {
		return false
	}
	for _, v := range valueSlice {
		if !IsValidUUID(v) {
			return false
		}
	}
	return true
}
