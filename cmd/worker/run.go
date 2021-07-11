package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/nite-coder/blackbear-demo/internal/pkg/config"
	starterWorkflow "github.com/nite-coder/blackbear-demo/pkg/workflow"
	"github.com/nite-coder/blackbear/pkg/log"
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
		tp, err := cfg.TracerProvider("worker")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Cleanly shutdown and flush telemetry when the application exits.
		defer func(ctx context.Context) {
			// Do not make the application hang when it is shutdown.
			ctx, cancel = context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				log.Err(err).Panic("tp shutdown failed")
			}
		}(ctx)

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

		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Err(err).Fatalf("Unable to start worker")
		}

		log.Info("worker has stopped")
	},
}
