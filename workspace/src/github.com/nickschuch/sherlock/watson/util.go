package main

import (
	"fmt"

	"k8s.io/client-go/pkg/api/v1"
)

func restarts(statuses []v1.ContainerStatus, name string) (v1.ContainerStatus, error) {
	for _, status := range statuses {
		if status.Name == name {
			return status, nil
		}
	}

	return v1.ContainerStatus{}, fmt.Errorf("cannot find container status with name: %s", name)
}
