package middleware

import (
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/napnap"
)

// LoggerMW is a logger middleware struct
type LoggerMW struct {
}

// NewLoggerMW returns LoggerMW middlware instance
func NewLoggerMW() *LoggerMW {
	return &LoggerMW{}
}

// Invoke function is a middleware entry
func (m *LoggerMW) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	ctx := c.StdContext()

	// save request id to logger
	ctx = log.Str("request_id", RequestIDFromContext(ctx)).WithContext(ctx)
	c.SetStdContext(ctx)

	next(c)
}
