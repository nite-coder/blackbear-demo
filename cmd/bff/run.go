package bff

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/starter/internal/pkg/config"
	"github.com/jasonsoft/starter/pkg/bff/delivery/gql"
	"github.com/spf13/cobra"
)

// RunCmd 是 bff service 的進入口
var RunCmd = &cobra.Command{
	Use:   "bff",
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
			log.Panicf("main: bff initialize failed: %v", err)
			return
		}

		// enable tracer
		fn := initTracer(cfg)
		defer fn()

		// start http server
		nap := gql.NewHTTPServer(_eventClient, _walletClient, _temporalClient)
		httpServer := &http.Server{
			Addr:    cfg.BFF.HTTPBind,
			Handler: nap,
		}

		go func() {
			// service connections
			log.Infof("bff is serving HTTP on %s\n", httpServer.Addr)
			err := httpServer.ListenAndServe()
			if err != nil {
				log.Errorf("main: http server listen failed: %v\n", err)
			}
		}()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
		<-stopChan
		log.Info("main: shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Errorf("main: http server shutdown error: %v", err)
		} else {
			log.Info("main: gracefully stopped")
		}

	},
}
