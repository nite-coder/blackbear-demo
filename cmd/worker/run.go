package worker

import (
	"fmt"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	starterWorkflow "github.com/jasonsoft/starter/pkg/workflow"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/api/global"
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
		fn := initTracer(cfg)
		defer fn()
		tr := global.Tracer("")
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
			log.Fatalf("Unable to create client", err)
		}
		defer c.Close()

		w := worker.New(c, "default", worker.Options{
			WorkerStopTimeout: 10, // 10 sec
		})

		w.RegisterWorkflow(starterWorkflow.PublishEventWorkflow)
		w.RegisterActivity(starterWorkflow.WithdrawActivity)
		w.RegisterActivity(starterWorkflow.PublishEventActivity)

		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Unable to start worker", err)
		}

		log.Info("worker has stopped")
	},
}
