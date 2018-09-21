package utils

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

const (
	IgnorePodAnnotation = "sherlock.nickschuch.github.com/watson-ignore"
)

func HasRestarts(statuses []corev1.ContainerStatus, name string) (corev1.ContainerStatus, error) {
	for _, status := range statuses {
		if status.Name == name {
			return status, nil
		}
	}

	return corev1.ContainerStatus{}, fmt.Errorf("cannot find container status with name: %s", name)
}

// IsIgnored checks for the pod annotation which instructs watson to ignore it.
func IsIgnored(pod *corev1.Pod) bool {
	if _, ok := pod.ObjectMeta.Annotations[IgnorePodAnnotation]; ok {
		return true
	}
	return false
}