package util

import (
	"os"
)

// Return the value of an environment variable, or a default value if it is not set.
func GetEnv(key, usual string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return usual
	}
	return value
}
