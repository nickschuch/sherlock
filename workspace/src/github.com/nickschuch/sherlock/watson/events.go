package main

import (
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

const fileEvents = "events.log"

func getEvents(kubeClient *kubernetes.Clientset, pod *v1.Pod) ([]byte, error) {
	var events []byte

	list, err := kubeClient.CoreV1().Events(pod.Namespace).Search(api.Scheme, pod)
	if err != nil {
		return events, err
	}

	var tmp []string

	for _, event := range list.Items {
		tmp = append(tmp, fmt.Sprintf("%s - %s - %s - %s", event.CreationTimestamp, event.Type, event.Reason, event.Message))
	}

	file := strings.Join(tmp, "\n")

	return []byte(file), nil
}
