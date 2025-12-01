package utils

import (
	"os"
)

// GetEnv retrieves the value of the environment variable named by the key.
func GetEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}