package workflow

import (
	"context"
	"time"

	"github.com/nite-coder/blackbear/pkg/log"
	"go.temporal.io/sdk/workflow"
)

// CronWorkflow is a cron workflow definition.
func CronWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	err := workflow.ExecuteActivity(ctx, CronActivity).Get(ctx, nil)
	if err != nil {
		log.Err(err).Error("CronActivity failed.")
		return err
	}

	if workflow.HasLastCompletionResult(ctx) {
		log.Info("HasLastCompletionResult")
	} else {
		log.Info("no result from last task")
	}

	log.Info("CronWorkflow completed.")

	return nil
}

func CronActivity(ctx context.Context) error {
	logger := getLogger(ctx)
	logger.Infof("workflow: begin CronActivity fn at %s", time.Now().String())

	return nil
}
