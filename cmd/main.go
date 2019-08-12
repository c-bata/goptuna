package main

import (
	"os"

	"github.com/c-bata/goptuna/cmd/createstudy"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goptuna",
	Short: "A command line interface for Goptuna",
}

func main() {
	rootCmd.AddCommand(createstudy.GetCommand())
	err := rootCmd.Execute()
	if err != nil {
		rootCmd.PrintErrln(err)
		os.Exit(1)
	}
}
