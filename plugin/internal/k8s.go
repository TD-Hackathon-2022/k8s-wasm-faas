package internal

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func newK8sClientset() kubernetes.Clientset {
	config := newK8sConfig()

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	return *clientset
}

func newK8sConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))

	if err != nil {
		fmt.Println("load k8s config error: ", err)
		os.Exit(1)
	}

	return config
}
