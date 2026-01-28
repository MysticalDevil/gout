package gout

import (
	"log/slog"
	"os"
	"time"
)

func Logger() HandlerFunc {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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
			slog.String("ip", c.Req.RemoteAddr),
		)
	}
}
