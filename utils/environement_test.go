package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvReturnsValueWhenPresent(t *testing.T) {
	t.Setenv("MESSAGE_SERVICE_TEST_ENV", "configured")

	value := GetEnv("MESSAGE_SERVICE_TEST_ENV", "fallback")

	assert.Equal(t, "configured", value)
}

func TestGetEnvReturnsDefaultWhenMissing(t *testing.T) {
	value := GetEnv("MESSAGE_SERVICE_TEST_ENV_MISSING", "fallback")

	assert.Equal(t, "fallback", value)
}
