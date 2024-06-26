package util

import (
	"fmt"
	"os"
	"strings"
)

// Return the value of an environment variable,
// or a default value if it is not set.
func DefaultEnv(key, usual string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return usual
	}
	return value
}

// initSecret returns secret either from env variable or from a file
func initSecret(prefix string) (string, error) {
	key, ok := os.LookupEnv(prefix)
	if ok {
		return key, nil
	}
	file, ok := os.LookupEnv(fmt.Sprintf("%s_FILE", prefix))
	if !ok {
		return "", fmt.Errorf("%s or %s_FILE must be set", prefix, prefix)
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("read %s_FILE: %w", prefix, err)
	}
	return strings.TrimSpace(string(raw)), nil
}
