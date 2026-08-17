// Package logger configures the shared slog logger. Call Setup once at startup;
// afterwards use slog.Info / slog.Error anywhere (global logger).
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Setup installs a JSON logger writing to stdout, tagged with the service name.
func Setup(serviceName, level string) {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parseLevel(level)})
	slog.SetDefault(slog.New(h).With("service", serviceName))
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
