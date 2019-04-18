package util

import (
	"fmt"
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

func Join(format string, data map[string]string) []string {
	var result []string

	for key, value := range data {
		result = append(
			result,
			fmt.Sprintf(
				format,
				key,
				value,
			),
		)
	}

	return result
}
