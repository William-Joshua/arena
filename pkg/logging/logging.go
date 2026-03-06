// Package log provides a structured, context-aware logger built on go.uber.org/zap,
// with Fx lifecycle integration.
//
// The package is intentionally named "log" (not "logging") to allow clean import
// aliases in consuming packages:
//
//	import log "cc.io/arena/pkg/logging"
package log

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ctxKeyRequestID is the context key for the request ID.
type ctxKeyRequestID struct{}

// ctxKeyCorrelationID is the context key for the correlation ID.
type ctxKeyCorrelationID struct{}

// WithRequestID returns a new context that carries the given request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID{}, id)
}

// WithCorrelationID returns a new context that carries the given correlation ID.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyCorrelationID{}, id)
}

// RequestID extracts the request ID from ctx, returning "" when absent.
func RequestID(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyRequestID{}).(string)
	return v
}

// CorrelationID extracts the correlation ID from ctx, returning "" when absent.
func CorrelationID(ctx context.Context) string {
	v, _ := ctx.Value(ctxKeyCorrelationID{}).(string)
	return v
}

// Logger wraps *zap.Logger and adds context-aware helpers.
type Logger struct {
	*zap.Logger
}

// With returns a *zap.Logger enriched with request_id and correlation_id
// fields pulled from ctx (fields are omitted when empty).
func (l *Logger) With(ctx context.Context) *zap.Logger {
	logger := l.Logger
	if id := RequestID(ctx); id != "" {
		logger = logger.With(zap.String("request_id", id))
	}
	if id := CorrelationID(ctx); id != "" {
		logger = logger.With(zap.String("correlation_id", id))
	}
	return logger
}

// New constructs a production *Logger.
func New() (*Logger, error) {
	z, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{z}, nil
}

// NewFx constructs a *Logger suitable for use inside an Fx application.
// It registers an OnStop hook that syncs the underlying zap logger on shutdown.
func NewFx(lc fx.Lifecycle) (*Logger, error) {
	l, err := New()
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			// Sync flushes any buffered log entries.  Errors are intentionally
			// swallowed here since they are usually benign (e.g. /dev/stderr).
			_ = l.Sync()
			return nil
		},
	})
	return l, nil
}

// Module is an Fx option that provides *Logger.
var Module = fx.Provide(NewFx)
