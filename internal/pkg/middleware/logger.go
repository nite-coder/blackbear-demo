package middleware

import (
	"github.com/nite-coder/blackbear/pkg/web"
	"github.com/nite-coder/blackbear/pkg/log"
)

// LoggerMW is a logger middleware struct
type LoggerMW struct {
}

// NewLoggerMW returns LoggerMW middlware instance
func NewLoggerMW() *LoggerMW {
	return &LoggerMW{}
}

// Invoke function is a middleware entry
func (m *LoggerMW) Invoke(c *web.Context, next web.HandlerFunc) {
	ctx := c.StdContext()

	// save request id to logger
	ctx = log.Str("request_id", RequestIDFromContext(ctx)).WithContext(ctx)
	c.SetStdContext(ctx)

	next(c)
}
