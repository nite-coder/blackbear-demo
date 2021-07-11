package gql

import (
	"context"

	internalMiddleware "github.com/nite-coder/blackbear-demo/internal/pkg/middleware"
	starterWorkflow "github.com/nite-coder/blackbear-demo/pkg/workflow"
	"github.com/nite-coder/blackbear/pkg/log"
	"go.opentelemetry.io/otel"
	"go.temporal.io/sdk/client"
)

func (r *mutationResolver) PublishEvent(ctx context.Context, input PublishEventInput) (*bool, error) {
	logger := log.FromContext(ctx)
	logger.Debug("gql: begin publish event fn")

	workflowOptions := client.StartWorkflowOptions{
		ID:        "publish_event_workflow",
		TaskQueue: "default",
	}

	ctx = context.WithValue(ctx, starterWorkflow.PropagateKey, &starterWorkflow.Values{Key: "request_id", Value: internalMiddleware.RequestIDFromContext(ctx)})

	tr := otel.Tracer("")
	_, span := tr.Start(ctx, "StartWorkflow-PublishEventWorkflow-me")
	defer span.End()

	we, err := r.temporalClient.ExecuteWorkflow(ctx, workflowOptions, "PublishEventWorkflow")
	if err != nil {
		return nil, err
	}

	logger.Debugf("Started workflow. WorkflowID %s, RunID: %s", we.GetID(), we.GetRunID())

	// Synchronously wait for the workflow completion.
	err = we.Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	logger.Debug("Workflow done")

	return nil, nil
}
