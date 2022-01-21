package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "faas",
		Short: "wasm runtime is provided on k8s clusters, which may make it easier to run faas on k8s clusters",
	}
)

func init() {
	rootCmd.AddCommand(submitCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(runCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
