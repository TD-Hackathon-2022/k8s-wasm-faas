package main

import (
	"context"
	"fmt"
	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubeInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var faasLabels = map[string]string{
	"runtime": "wasm",
	"type":    "faas-wasm",
}

func main() {
	config, err := rest.InClusterConfig()

	if err != nil {
		fmt.Println("load k8s config error")
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		fmt.Println("create k8s client error")
		os.Exit(1)
	}

	kubeInformerFactory := kubeInformers.NewSharedInformerFactoryWithOptions(
		clientset,
		time.Second*30,
		kubeInformers.WithNamespace(coreV1.NamespaceDefault),
		kubeInformers.WithTweakListOptions(internalinterfaces.TweakListOptionsFunc(filterLabel)),
	)

	configMapsInformer := kubeInformerFactory.Core().V1().ConfigMaps()

	configMapsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(newConfigMap interface{}) {
			configmap := newConfigMap.(*coreV1.ConfigMap)
			fmt.Printf("%s\t%s\n", configmap.UID, configmap.Name)
			createFaasBuilderJob(clientset, configmap)
		},
	})

	stopCh := setupSignalHandler()
	kubeInformerFactory.Start(stopCh)
	<-stopCh
}

func setupSignalHandler() (stopCh <-chan struct{}) {
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

func filterLabel(listOption *v1.ListOptions) {
	listOption.LabelSelector = labels.FormatLabels(faasLabels)
}

func createFaasBuilderJob(clientset *kubernetes.Clientset, lambda *coreV1.ConfigMap) {
	var (
		ttLSecondsAfterFinished int32 = 100
		defaultMode             int32 = 0600
		targetHost                    = os.Getenv("TARGET_HOST")
		targetPort                    = os.Getenv("TARGET_PORT")
		targetUser                    = os.Getenv("TARGET_USER")
		targetPath                    = os.Getenv("TARGET_PATH")
		hostPathDir                   = coreV1.HostPathDirectory
	)

	_, err := clientset.BatchV1().Jobs(coreV1.NamespaceDefault).Create(
		context.TODO(),
		&batchV1.Job{
			TypeMeta: metaV1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Jobs",
			},
			ObjectMeta: metaV1.ObjectMeta{
				Name:      "faas-wasm-builder-" + string(lambda.UID),
				Namespace: coreV1.NamespaceDefault,
				Labels:    faasLabels,
			},
			Spec: batchV1.JobSpec{
				TTLSecondsAfterFinished: &ttLSecondsAfterFinished,
				Template: coreV1.PodTemplateSpec{
					Spec: coreV1.PodSpec{
						NodeSelector:  map[string]string{"faas-wasm-runtime": "wasm"},
						RestartPolicy: coreV1.RestartPolicyNever,
						Volumes: []coreV1.Volume{
							{
								Name: "faas-builder-ssh-private-key",
								VolumeSource: coreV1.VolumeSource{
									Secret: &coreV1.SecretVolumeSource{
										SecretName:  "faas-builder-ssh-private-key",
										DefaultMode: &defaultMode,
									},
								},
							},
							{
								Name: "lambda",
								VolumeSource: coreV1.VolumeSource{
									ConfigMap: &coreV1.ConfigMapVolumeSource{
										LocalObjectReference: coreV1.LocalObjectReference{Name: lambda.Name},
									},
								},
							},
							{
								Name: "cargo",
								VolumeSource: coreV1.VolumeSource{
									HostPath: &coreV1.HostPathVolumeSource{
										Type: &hostPathDir,
										Path: "/root/cargo-cache",
									},
								},
							},
						},
						Containers: []coreV1.Container{
							{
								Name:  "k8s-faas-builder",
								Image: "redxiiikk/k8s-faas-builder:latest",
								VolumeMounts: []coreV1.VolumeMount{
									{
										Name:      "faas-builder-ssh-private-key",
										ReadOnly:  false,
										MountPath: "/opt/k8s-faas-builder/config/",
									},
									{
										Name:      "lambda",
										ReadOnly:  true,
										MountPath: "/opt/k8s-faas-builder/lambda",
									},
									{
										Name:      "cargo",
										ReadOnly:  false,
										MountPath: "/root/.cargo/registry",
									},
								},
								Env: []coreV1.EnvVar{
									{Name: "FUNCTION_NAME", Value: lambda.Name},
									{Name: "TARGET_HOST", Value: targetHost},
									{Name: "TARGET_PORT", Value: targetPort},
									{Name: "TARGET_USER", Value: targetUser},
									{Name: "TARGET_PATH", Value: targetPath},
								},
							},
						},
					},
				},
			},
		},
		metaV1.CreateOptions{},
	)

	if err != nil {
		fmt.Println("create faas wasm builder job error", err)
		return
	}
}
