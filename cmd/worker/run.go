package worker

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	starterWorkflow "github.com/jasonsoft/starter/pkg/workflow"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
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
				trace := make([]byte, 4096)
				runtime.Stack(trace, true)
				log.Str("stack_trace", string(trace)).Err(err).Panic("unknown error")
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

		// start worker
		c, err := client.NewClient(client.Options{
			HostPort: "temporal:7233",
		})
		if err != nil {
			log.Fatalf("Unable to create client", err)
		}
		defer c.Close()

		w := worker.New(c, "default", worker.Options{})

		w.RegisterWorkflow(starterWorkflow.PublishEventWorkflow)
		w.RegisterActivity(starterWorkflow.WithdrawActivity)
		w.RegisterActivity(starterWorkflow.PublishEventActivity)

		go func() {
			log.Info("worker started")
			err = w.Run(nil)
			if err != nil {
				log.Fatalf("Unable to start worker", err)
			}
		}()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
		<-stopChan
		log.Info("main: shutting down worker...")

		w.Stop()
		log.Info("main: worker was stopped")

	},
}
