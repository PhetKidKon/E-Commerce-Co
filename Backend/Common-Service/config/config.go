// Package config reads configuration from environment variables with defaults.
package config

import (
	"os"
	"strconv"
)

// Get returns the env value for key, or def if unset/empty.
func Get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// GetInt returns the env value parsed as int, or def if unset/invalid.
func GetInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
