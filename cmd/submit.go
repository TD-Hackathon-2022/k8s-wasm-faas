package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hackathon-2022/k8s-faas-plugin/internal"
	"github.com/hackathon-2022/k8s-faas-plugin/tools"
	"github.com/spf13/cobra"
)

var (
	submitCmd = &cobra.Command{
		Use:   "submit",
		Short: "Store function in ConfigMap",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			validArgs(args)

			faasScriptPath := getAbsFaasScriptPath(args[0])
			filename := tools.GetFileName(faasScriptPath)
			fileContent := tools.ReadAll(faasScriptPath)

			internal.CreateConfigMap(filename, fileContent)
		},
	}
)

func validArgs(args []string) {
	if len(args) != 1 {
		fmt.Println("Need function file path")
		os.Exit(1)
	}
}

func getAbsFaasScriptPath(faasScriptPath string) string {
	if !filepath.IsAbs(faasScriptPath) {
		currentDir, _ := os.Getwd()
		return path.Join(currentDir, faasScriptPath)
	}

	return faasScriptPath
}
