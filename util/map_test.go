package util

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestJoinMapWithData(t *testing.T) {
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	sortedExpected := []string{"-e key1 value1", "-e key2 value2", "-e key3 value3"}
	sort.Strings(sortedExpected)

	sortedResult := Join("-e %s %s", data)
	sort.Strings(sortedResult)

	assert.Equal(
		t,
		sortedExpected,
		sortedResult,
	)
}

func TestJoinEmptyMap(t *testing.T) {
	data := map[string]string{}

	assert.Equal(
		t,
		[]string{},
		Join("-e %s %s", data),
	)
}
