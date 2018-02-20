package utils

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func HasRestarts(statuses []corev1.ContainerStatus, name string) (corev1.ContainerStatus, error) {
	for _, status := range statuses {
		if status.Name == name {
			return status, nil
		}
	}

	return corev1.ContainerStatus{}, fmt.Errorf("cannot find container status with name: %s", name)
}
