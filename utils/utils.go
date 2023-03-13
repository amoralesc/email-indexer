package utils

import "os"

// GetenvOrDefault returns the value of the environment variable named by the key.
// If the env variable is not present, it returns the defaultValue instead.
func GetenvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
