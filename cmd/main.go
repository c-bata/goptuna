package main

import (
	"fmt"
	"os"

	"github.com/c-bata/goptuna/cmd/createstudy"
	"github.com/spf13/cobra"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	version  string
	revision string

	rootCmd = &cobra.Command{
		Use:   "goptuna",
		Short: "A command line interface for Goptuna",
	}
)

func main() {
	rootCmd.AddCommand(createstudy.GetCommand())
	if version != "" && revision != "" {
		rootCmd.Version = fmt.Sprintf("%s (rev: %s)", version, revision)
	}
	err := rootCmd.Execute()
	if err != nil {
		rootCmd.PrintErrln(err)
		os.Exit(1)
	}
}
