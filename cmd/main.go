package main

import (
	"fmt"
	"os"

	"github.com/c-bata/goptuna/cmd/createstudy"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goptuna",
	Short: "A command line interface for Goptuna",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("I am root cmd.")
	},
}

func main() {
	rootCmd.AddCommand(createstudy.GetCommand())
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
