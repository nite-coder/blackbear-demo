package frontend

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nite-coder/blackbear-demo/internal/pkg/config"
	"github.com/nite-coder/blackbear-demo/pkg/frontend/delivery/gql"
	"github.com/nite-coder/blackbear/pkg/log"
	"github.com/spf13/cobra"
)

// RunCmd 是 frontend service 的進入口
var RunCmd = &cobra.Command{
	Use:   "frontend",
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
			log.Panicf("main: frontend initialize failed: %v", err)
			return
		}

		// enable tracer
		tp, err := cfg.TracerProvider("frontend")
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

		// start http server
		webServ := gql.NewHTTPServer(_eventClient, _walletClient, _temporalClient)
		httpServer := &http.Server{
			Addr:    cfg.Frontend.HTTPBind,
			Handler: webServ,
		}

		go func() {
			// service connections
			log.Infof("frontend is serving HTTP on %s\n", httpServer.Addr)
			err := httpServer.ListenAndServe()
			if err != nil {
				log.Errorf("main: http server listen failed: %v\n", err)
			}
		}()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
		<-stopChan
		log.Info("main: shutting down server...")

		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Errorf("main: http server shutdown error: %v", err)
		} else {
			log.Info("main: gracefully stopped")
		}

	},
}
