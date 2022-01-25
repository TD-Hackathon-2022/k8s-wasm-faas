package internal

import (
	"context"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"os"
)

func CreateConfigMap(name, content string) {
	configMap := coreV1.ConfigMap{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: coreV1.NamespaceDefault,
			Labels:    faasLabels,
		},
		Data: map[string]string{
			name: content,
		},
	}

	clientset := newK8sClientset()
	configMapsOperator := clientset.CoreV1().ConfigMaps(coreV1.NamespaceDefault)

	_, err := configMapsOperator.Create(context.Background(), &configMap, metaV1.CreateOptions{})
	if err != nil {
		fmt.Println("create configmap error!!!")
		os.Exit(1)
	}
}

func ListConfigMapNames() map[string]string {
	clientset := newK8sClientset()
	configMapsOperator := clientset.CoreV1().ConfigMaps(coreV1.NamespaceDefault)

	configMapList, err := configMapsOperator.List(context.Background(), metaV1.ListOptions{
		LabelSelector: labels.FormatLabels(faasLabels),
	})
	if err != nil {
		fmt.Println("get configmaps error!!!")
		os.Exit(1)
	}

	result := make(map[string]string)
	for index := range configMapList.Items {
		configmap := configMapList.Items[index]
		result[string(configmap.UID)] = configmap.Name
	}

	return result
}
