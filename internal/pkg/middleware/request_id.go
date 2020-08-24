package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/napnap"
)

type RequestIDKey string

func (r RequestIDKey) String() string {
	return string(r)
}

const (
	Key RequestIDKey = "X-Request-Id"
)

// RequestID is request_id middleware struct
// TODO: move to napnap
type RequestID struct {
}

// NewRequestID returns NewRequestID middlware instance
func NewRequestID() *RequestID {
	return &RequestID{}
}

// Invoke function is a middleware entry
func (m *RequestID) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	requestID := c.RequestHeader(Key.String())
	if requestID == "" {
		requestID = uuid.New().String()
	}
	c.Request.Header.Set(Key.String(), requestID)

	// save request id to logger
	ctx := log.Str("request_id", requestID).WithContext(c.Request.Context())

	// save request id to request context
	ctx = context.WithValue(ctx, Key, requestID)
	c.Request = c.Request.WithContext(ctx)

	// Set X-Request-Id header
	c.Writer.Header().Set(Key.String(), requestID)
	next(c)
}

// RequestIDFromContext 從 ctx 中取得 request id, 如果沒有即時產生一個
func RequestIDFromContext(ctx context.Context) string {
	rid, ok := ctx.Value(Key).(string)
	if !ok {
		// 產生 requestID 並傳下去
		rid = uuid.New().String()
		return rid
	}
	return rid
}
