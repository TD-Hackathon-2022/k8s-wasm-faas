package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hackathon-2022/k8s-faas-plugin/service"
	"github.com/hackathon-2022/k8s-faas-plugin/utils"
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
			filename := utils.GetFileName(faasScriptPath)
			fileContent := utils.ReadAll(faasScriptPath)

			service.CreateConfigMap(filename, fileContent)
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
