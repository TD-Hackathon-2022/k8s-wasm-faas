package internal

import (
	"bytes"
	"context"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kubeInformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"strconv"
	"strings"
	"time"
)

func CreatePod(imageName string) *coreV1.Pod {
	clientset := newK8sClientset()
	podOperator := clientset.CoreV1().Pods(coreV1.NamespaceDefault)

	timestamp := strconv.Itoa(int(time.Now().Unix()))

	createdPod, err := podOperator.Create(
		context.TODO(),
		&coreV1.Pod{
			TypeMeta: metaV1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectMeta: metaV1.ObjectMeta{
				Name:   "wasm-lambda-function-" + timestamp,
				Labels: wasmPodLabel,
			},
			Spec: coreV1.PodSpec{
				Containers: []coreV1.Container{
					{
						Name:            "wasm-lambda-function-" + timestamp,
						Image:           imageName,
						ImagePullPolicy: coreV1.PullIfNotPresent,
					},
				},
				NodeSelector: map[string]string{
					"runtime": "wasm",
				},
			},
		},
		metaV1.CreateOptions{},
	)

	if err != nil {
		fmt.Println("create wasm function pod error:", err)
		os.Exit(1)
	}

	return createdPod
}

func DeletePod(podName string) {
	clientset := newK8sClientset()
	podOperator := clientset.CoreV1().Pods(coreV1.NamespaceDefault)

	err := podOperator.Delete(context.TODO(), podName, metaV1.DeleteOptions{})
	if err != nil {
		fmt.Println("delete faas lambda pod error: ", err)
		os.Exit(1)
	}
}

func ExecInPod(podName string, params []string) (string, string, error) {
	clientset := newK8sClientset()
	req := clientset.CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(coreV1.NamespaceDefault).
		SubResource("exec").
		VersionedParams(
			&coreV1.PodExecOptions{
				Command: []string{"-- " + strings.Join(params, "=")},
				Stdin:   false,
				Stdout:  true,
				Stderr:  true,
				TTY:     false,
			},
			scheme.ParameterCodec,
		)

	var stdout, stderr bytes.Buffer

	exec, err := remotecommand.NewSPDYExecutor(newK8sConfig(), "POST", req.URL())
	if err != nil {
		return "", "", err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return "", "", err
	}

	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}

func WaitPodStatus(podName string) {
	filterLabel := func(listOption *metaV1.ListOptions) {
		listOption.FieldSelector = "metadata.name=" + podName
		listOption.LabelSelector = labels.FormatLabels(wasmPodLabel)
	}

	clientset := newK8sClientset()

	kubeInformerFactory := kubeInformers.NewSharedInformerFactoryWithOptions(
		&clientset,
		time.Second*30,
		kubeInformers.WithNamespace(coreV1.NamespaceDefault),
		kubeInformers.WithTweakListOptions(filterLabel),
	)

	podInformer := kubeInformerFactory.Core().V1().Pods()

	stopChan := make(chan struct{})

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(_, newObj interface{}) {
			newPod := newObj.(*coreV1.Pod)
			if newPod.Status.Phase == coreV1.PodRunning {
				fmt.Println("build pod status run...")
				close(stopChan)
			}
		},
	})

	kubeInformerFactory.Start(stopChan)
	<-stopChan
}
