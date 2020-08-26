package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/jasonsoft/napnap"
)

type requestIDKey string

func (r requestIDKey) String() string {
	return string(r)
}

const (
	key requestIDKey = "X-Request-Id"
)

// RequestIDMW is a request_id middleware struct
// TODO: move to napnap
type RequestIDMW struct {
}

// NewRequestIDMW returns NewRequestID middlware instance
func NewRequestIDMW() *RequestIDMW {
	return &RequestIDMW{}
}

// Invoke function is a middleware entry
func (m *RequestIDMW) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	ctx := c.StdContext()

	requestID := c.RequestHeader(key.String())
	if requestID == "" {
		requestID = uuid.New().String()
	}
	c.Request.Header.Set(key.String(), requestID)

	// save request id to request context
	ctx = context.WithValue(ctx, key, requestID)
	c.SetStdContext(ctx)

	// Set X-Request-Id header
	c.Writer.Header().Set(key.String(), requestID)
	next(c)
}

// RequestIDFromContext get requestID from context
func RequestIDFromContext(ctx context.Context) string {
	rid, ok := ctx.Value(key).(string)
	if !ok {
		rid = uuid.New().String()
		return rid
	}
	return rid
}

// SetRequestIDToContext save the requestID into context
func SetRequestIDToContext(ctx context.Context, requestID string) context.Context {
	// save request id to request context
	ctx = context.WithValue(ctx, key, requestID)
	return ctx
}
