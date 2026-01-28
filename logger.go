package gout

import (
	"log/slog"
	"os"
	"time"
)

type FormatterType int

const (
	TextFormatter FormatterType = iota
	JSONFormatter
)

type LoggerConfig struct {
	Formatter FormatterType
}

func Logger(configs ...LoggerConfig) HandlerFunc {
	cfg := LoggerConfig{
		Formatter: TextFormatter,
	}

	if len(configs) > 0 {
		cfg = configs[0]
	}

	var handler slog.Handler
	switch cfg.Formatter {
	case JSONFormatter:
		handler = slog.NewJSONHandler(os.Stdout, nil)
	default:
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	logger := slog.New(handler)

	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		latency := time.Since(t)

		logger.Info(
			"Request handled",
			slog.Int("status", c.StatusCode),
			slog.String(
				"method",
				c.Req.Method,
			),
			slog.String("path", c.Path),
			slog.Duration("latency", latency),
			slog.String("ip", c.ClientIP()),
		)
	}
}
