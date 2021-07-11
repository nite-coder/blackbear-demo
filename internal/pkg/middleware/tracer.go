package middleware

import (
	"github.com/nite-coder/blackbear/pkg/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// TracerMW is tracer middleware struct
type TracerMW struct {
}

// NewTracerMW returns TracerMW middlware instance
func NewTracerMW() *TracerMW {
	return &TracerMW{}
}

// Invoke function is a middleware entry
func (m *TracerMW) Invoke(c *web.Context, next web.HandlerFunc) {
	tr := otel.Tracer("")

	ctx, span := tr.Start(c.StdContext(), "web") // TODO: we need to replace "web" to gql operation's name which is more meaningful
	c.SetStdContext(ctx)

	lblRequestID := attribute.KeyValue{
		Key:   attribute.Key("request_id"),
		Value: attribute.StringValue(RequestIDFromContext(ctx)),
	}

	span.SetAttributes(lblRequestID)
	defer span.End()

	next(c)
}
