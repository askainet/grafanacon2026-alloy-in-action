package telemetry

import (
	"log/slog"
	"os"
	"strings"
)

// InitLogging initializes structured logging with slog to stdout.
// Logs are written in JSON format for easy parsing by Alloy/Loki.
func InitLogging(logLevel string) {
	// Parse log level
	level := parseLogLevel(logLevel)

	// JSON handler for structured logs to stdout
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	slog.SetDefault(slog.New(handler))
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
