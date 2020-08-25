package gql

import (
	"context"

	"github.com/jasonsoft/log/v2"
	"go.temporal.io/sdk/client"
)

func (r *mutationResolver) PublishEvent(ctx context.Context, input []*PublishEventInput) (*bool, error) {
	logger := log.FromContext(ctx)
	logger.Debug("gql: begin publish event fn")

	workflowOptions := client.StartWorkflowOptions{
		ID:        "publish_event_workflow",
		TaskQueue: "default",
	}

	we, err := r.temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, "PublishEventWorkflow")
	if err != nil {
		return nil, err
	}

	log.Debugf("Started workflow. WorkflowID %s, RunID: %s", we.GetID(), we.GetRunID())

	// Synchronously wait for the workflow completion.
	err = we.Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	log.Debug("Workflow done")

	return nil, nil
}
