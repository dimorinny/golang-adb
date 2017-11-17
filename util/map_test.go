package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoinMapWithData(t *testing.T) {
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	assert.Equal(
		t,
		[]string{"-e key1 value1", "-e key2 value2", "-e key3 value3"},
		Join("-e %s %s", data),
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
