// Package config reads configuration from environment variables with defaults,
// and can optionally load a .env file for local development.
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadDotenv loads a .env file into the process environment if one exists.
// Optional by design: in docker/prod the env vars are set directly and there
// is no .env. Existing env vars are NOT overwritten (godotenv.Load behavior),
// so container-provided values always win over the file.
func LoadDotenv() {
	// run from Backend/ -> ".env" ; run from a service dir -> "../.env"
	for _, p := range []string{".env", "../.env"} {
		if err := godotenv.Load(p); err == nil {
			return
		}
	}
}

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
