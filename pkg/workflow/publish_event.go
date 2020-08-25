package workflow

import (
	"context"
	"time"

	"github.com/jasonsoft/log/v2"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func PublishEventWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    1, // run once.  it means no retry
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	log.Info("workflow: publish event workflow started")

	err := workflow.ExecuteActivity(ctx, WithdrawActivity).Get(ctx, nil)
	if err != nil {
		log.Err(err).Warn("workflow: withdraw activity failed.")
		return err
	}

	err = workflow.ExecuteActivity(ctx, PublishEventActivity).Get(ctx, nil)
	if err != nil {
		log.Err(err).Warn("workflow: publish event activity failed.")
		return err
	}

	log.Info("workflow: publish event workflow completed")

	return nil
}

func WithdrawActivity(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Debug("workflow: begin WithdrawActivity fn")

	req := walletProto.WithdrawRequest{
		TransId: "ab",
		Amount:  100,
	}

	_, err := _manager.WalletClient.Withdraw(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}

func PublishEventActivity(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Debug("workflow: begin PublishEventActivity fn")

	req := eventProto.UpdatePublishStatusRequest{
		EventId:         1,
		TransId:         "ab",
		PublishedStatus: eventProto.PublishedStatus_Published,
	}

	_, err := _manager.EventClient.UpdatePublishStatus(ctx, &req)
	if err != nil {
		return err
	}
	return nil
}
