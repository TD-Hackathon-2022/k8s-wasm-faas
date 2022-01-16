package service

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func newK8sClientset() kubernetes.Clientset {
	home := homedir.HomeDir()

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	return *clientset
}
