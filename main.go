package main

import (
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	kubeInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	home := homedir.HomeDir()

	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	kubeInformerFactory := kubeInformers.NewSharedInformerFactory(clientset, time.Second*30)

	configMapsInformer := kubeInformerFactory.Core().V1().ConfigMaps()

	configMapsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(newConfigMap interface{}) {
			configmap := newConfigMap.(*coreV1.ConfigMap)
			fmt.Println(configmap)
		},
	})

	stopCh := SetupSignalHandler()
	kubeInformerFactory.Start(stopCh)
}

func SetupSignalHandler() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
