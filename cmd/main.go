package main

import (
	"fmt"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jasonsoft/starter/cmd/bff"
	"github.com/jasonsoft/starter/cmd/event"
	"github.com/jasonsoft/starter/cmd/wallet"
	"github.com/jasonsoft/starter/cmd/worker"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "choose instance to run: bff, event, wallet, worker",
	Long:  ``,
}

func main() {
	rootCmd.AddCommand(worker.RunCmd)
	rootCmd.AddCommand(wallet.RunCmd)
	rootCmd.AddCommand(event.RunCmd)
	rootCmd.AddCommand(bff.RunCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
