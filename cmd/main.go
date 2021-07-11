package main

import (
	"fmt"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/nite-coder/blackbear-demo/cmd/event"
	"github.com/nite-coder/blackbear-demo/cmd/frontend"
	"github.com/nite-coder/blackbear-demo/cmd/wallet"
	"github.com/nite-coder/blackbear-demo/cmd/worker"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "choose instance to run: frontend, event, wallet, worker",
	Long:  ``,
}

func main() {
	rootCmd.AddCommand(worker.RunCmd)
	rootCmd.AddCommand(wallet.RunCmd)
	rootCmd.AddCommand(event.RunCmd)
	rootCmd.AddCommand(frontend.RunCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
