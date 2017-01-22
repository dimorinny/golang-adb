package util

import (
	"strconv"
)

func GetStringWithDefault(data map[string]interface{}, key string, defaultValue string) string {
	if value, ok := data[key]; ok {
		if value, ok := value.(string); ok {
			return value
		}
	}
	return defaultValue
}

func GetIntWithDefault(data map[string]interface{}, key string, defaultValue int) int {
	if value, ok := data[key]; ok {
		if value, ok := value.(string); ok {
			if intValue, err := strconv.Atoi(value); err == nil {
				return intValue
			}
		}
	}
	return defaultValue
}
