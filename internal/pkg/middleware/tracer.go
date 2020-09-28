package middleware

import (
	"github.com/jasonsoft/napnap"
	"go.opentelemetry.io/otel/api/global"
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
	tr := global.Tracer("")

	ctx, span := tr.Start(c.StdContext(), "web") // TODO: we need to replace "web" to gql operation's name which is more meaningful
	c.SetStdContext(ctx)
	span.SetAttribute("request_id", RequestIDFromContext(ctx))
	defer span.End()

	next(c)
}
