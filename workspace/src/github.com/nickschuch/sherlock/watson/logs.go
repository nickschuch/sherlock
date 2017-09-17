package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

const fileLogs = "output.log"

func getLogs(kubeClient *kubernetes.Clientset, namespace, pod, container string) ([]byte, error) {
	opts := &v1.PodLogOptions{
		Container: container,
		Previous:  true,
		TailLines: cliLogLines,
	}

	return kubeClient.CoreV1().Pods(namespace).GetLogs(pod, opts).DoRaw()
}
