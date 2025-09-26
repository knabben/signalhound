package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "stalker",
		Short: "stalker search for issues and flaky tests on Kubernetes",
		Long:  "stalker search for issues and flaky tests on Kubernetes",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
