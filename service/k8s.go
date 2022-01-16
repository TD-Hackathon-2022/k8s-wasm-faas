package service

import (
	"context"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
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

func CreateConfigMap(name, content string) {
	configMap := coreV1.ConfigMap{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: coreV1.NamespaceDefault,
			Labels: map[string]string{
				"runtime": "wasm",
			},
		},
		Data: map[string]string{
			name: content,
		},
	}

	clientset := newK8sClientset()
	configMaps := clientset.CoreV1().ConfigMaps(coreV1.NamespaceDefault)

	_, err := configMaps.Create(context.Background(), &configMap, metaV1.CreateOptions{})
	if err != nil {
		fmt.Println("create configmap error!!!")
		os.Exit(1)
	}
}
