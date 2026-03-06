// Package middleware provides Gin middleware used across the API server.
package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	log "cc.io/arena/pkg/logging"
)

// defaultSkipPaths contains routes that are excluded from access logging.
var defaultSkipPaths = map[string]bool{
	"/healthz": true,
}

// AccessLog returns a Gin middleware that writes a structured access-log entry
// for every request that is not in skipPaths.
//
// Logged fields:
//   - method, path, status, latency_ms, client_ip, user_agent, bytes_out
//   - request_id and correlation_id (when present in context)
//   - errors collected by c.Errors
//
// The default skip list (/healthz and /swagger/*) can be extended via the
// extraSkip parameter.
func AccessLog(logger *log.Logger, extraSkip ...string) gin.HandlerFunc {
	skip := make(map[string]bool, len(defaultSkipPaths)+len(extraSkip))
	for k, v := range defaultSkipPaths {
		skip[k] = v
	}
	for _, p := range extraSkip {
		skip[p] = true
	}

	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Skip swagger paths and anything the caller opted out of.
		if skip[path] || isSwaggerPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		latencyMs := time.Since(start).Milliseconds()

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("latency_ms", latencyMs),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("bytes_out", c.Writer.Size()),
		}

		// Attach errors gathered during request handling.
		if errs := c.Errors; len(errs) > 0 {
			fields = append(fields, zap.String("errors", errs.String()))
		}

		// Emit the log entry enriched with context IDs.
		logger.With(c.Request.Context()).Info("access", fields...)
	}
}

// isSwaggerPath returns true for any path under /swagger/.
func isSwaggerPath(path string) bool {
	return strings.HasPrefix(path, "/swagger")
}
