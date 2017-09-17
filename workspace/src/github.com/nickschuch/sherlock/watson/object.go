package main

import (
	"github.com/ghodss/yaml"
	"k8s.io/client-go/pkg/api/v1"
)

const fileObject = "object.yaml"

func getObject(pod v1.Pod) ([]byte, error) {
	return yaml.Marshal(pod)
}
