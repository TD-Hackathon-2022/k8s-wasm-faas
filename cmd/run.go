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
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createdPod := internal.CreatePod(args[0])
			defer internal.DeletePod(createdPod.Name)

			internal.WaitPodStatus(createdPod.Name)

			stdout, stderr, err := internal.ExecInPod(createdPod.Name, []string{"name", "hello"})
			if err != nil {
				fmt.Println("Exec Error: ", err)
				return
			}

			if stdout != "" {
				fmt.Println(" STDIN: ", stdout)
			}

			if stderr != "" {
				fmt.Println("STDERR: ", stderr)
			}
		},
	}
)
