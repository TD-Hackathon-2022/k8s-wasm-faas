package cmd

import (
	"fmt"
	"github.com/hackathon-2022/k8s-faas-plugin/internal"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "run wasm lambda function by name",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createdPod := internal.CreatePod(args[0])
			defer internal.DeletePod(createdPod.Name)

			internal.WaitPodStatus(createdPod.Name)

			var (
				stdout string
				stderr string
				err    error
			)

			if len(args) > 1 {
				stdout, stderr, err = internal.ExecInPod(createdPod.Name, args[1:])
			} else {
				stdout, stderr, err = internal.ExecInPod(createdPod.Name, []string{})
			}
			if err != nil {
				fmt.Println("Exec Error: ", err)
				return
			}

			if stdout != "" {
				fmt.Println("STDIN: ", stdout)
			}

			if stderr != "" {
				fmt.Println("STDERR: ", stderr)
			}
		},
	}
)
