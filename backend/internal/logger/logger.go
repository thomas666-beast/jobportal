package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Init configures the global slog logger.
// LOG_LEVEL: debug, info, warn, error (default: info)
// LOG_FORMAT: json, text (default: json in production, text in dev)
func Init() {
	level := slog.LevelInfo
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	format := os.Getenv("LOG_FORMAT")
	if format == "" {
		// Default to text in dev (nicer terminal output), JSON elsewhere.
		if os.Getenv("APP_ENV") == "production" {
			format = "json"
		} else {
			format = "text"
		}
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level}
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
