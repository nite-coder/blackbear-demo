package middleware

import (
	"github.com/jasonsoft/napnap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

// TracerMW is tracer middleware struct
type TracerMW struct {
}

// NewTracerMW returns TracerMW middlware instance
func NewTracerMW() *TracerMW {
	return &TracerMW{}
}

// Invoke function is a middleware entry
func (m *TracerMW) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	tr := otel.Tracer("")

	ctx, span := tr.Start(c.StdContext(), "web") // TODO: we need to replace "web" to gql operation's name which is more meaningful
	c.SetStdContext(ctx)

	lblRequestID := label.KeyValue{
		Key:   label.Key("request_id"),
		Value: label.StringValue(RequestIDFromContext(ctx)),
	}

	span.SetAttributes(lblRequestID)
	defer span.End()

	next(c)
}
