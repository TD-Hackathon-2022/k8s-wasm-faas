package cmd

import (
	"fmt"
	"github.com/hackathon-2022/k8s-faas-plugin/internal"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list all faas function on k8s clusters",
		Run: func(cmd *cobra.Command, args []string) {
			configMapNames := internal.ListConfigMapNames()
			for uid, name := range configMapNames {
				fmt.Printf("%s\t%s\n", uid, name)
			}
		},
	}
)
