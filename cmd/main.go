package main

import (
	"fmt"
	"os"

	"github.com/jasonsoft/starter/cmd/bff"
	"github.com/jasonsoft/starter/cmd/event"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "choose instance to run: bff, event, wallet, worker",
	Long:  ``,
}

func main() {
	rootCmd.AddCommand(event.RunCmd)
	rootCmd.AddCommand(bff.RunCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
