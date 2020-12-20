package worker

import (
	"context"
	"fmt"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	starterWorkflow "github.com/jasonsoft/starter/pkg/workflow"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/bridge/opentracing"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// RunCmd 是 worker service 的進入口
var RunCmd = &cobra.Command{
	Use:   "worker",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		defer log.Flush()
		defer func() {
			if r := recover(); r != nil {
				// unknown error
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("unknown error: %v", r)
				}
				log.Err(err).Panic("unknown error")
			}
		}()

		config.EnvPrefix = "STARTER"
		cfg := config.New("app.yml")
		err := initialize(cfg)
		if err != nil {
			log.Fatalf("main: initialize failed: %v", err)
			return
		}

		// enable tracer
		fn := cfg.InitTracer("worker")
		defer fn()

		tr := otel.Tracer("")
		bridgeTracer, _ := opentracing.NewTracerPair(tr)

		// start worker, one worker per process mode
		c, err := client.NewClient(client.Options{
			HostPort: cfg.Temporal.Address,
			ContextPropagators: []workflow.ContextPropagator{
				starterWorkflow.NewContextPropagator(),
			},
			Tracer: bridgeTracer,
		})
		if err != nil {
			log.Err(err).Fatal("Unable to create client")
		}
		defer c.Close()

		w := worker.New(c, "default", worker.Options{
			WorkerStopTimeout: 10, // 10 sec
		})

		w.RegisterWorkflow(starterWorkflow.PublishEventWorkflow)
		w.RegisterActivity(starterWorkflow.WithdrawActivity)
		w.RegisterActivity(starterWorkflow.PublishEventActivity)

		w.RegisterWorkflow(starterWorkflow.CronWorkflow)
		w.RegisterActivity(starterWorkflow.CronActivity)

		workflowOptions := client.StartWorkflowOptions{
			ID:           "cron_workflowID",
			TaskQueue:    "default",
			CronSchedule: "* * * * *",
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, starterWorkflow.CronWorkflow)
		if err != nil {
			panic(err)
		}

		log.Infof("Started workflow. WorkflowID: %s, RunID: %s", we.GetID(), we.GetRunID())

		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Err(err).Fatalf("Unable to start worker")
		}

		log.Info("worker has stopped")
	},
}
