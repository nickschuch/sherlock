package types

import (
	"time"
)

type Storage interface {
	Put(PutParams) (PutResponse, error)
	List(ListParams) (ListResponse, error)
	Inspect(InspectParams) (InspectResponse, error)
}

type PutParams struct {
	Incident Incident `yaml:"incident" json:"incident"`
}

type PutResponse struct {
	ID string `yaml:"id" json:"id"`
}

type InspectParams struct {
	ID string `yaml:"id" json:"id"`
}

type InspectResponse struct {
	Incident Incident `yaml:"incident" json:"incident"`
}

type ListParams struct {
	Namespace string `yaml:"namespace" json:"namespace"`
}

type ListResponse struct {
	Incidents []Incident `yaml:"incidents" json:"incidents"`
}

// Incident occurs when a pod is restarted.
type Incident struct {
	ID        string    `yaml:"id"        json:"id"`
	Created   time.Time `yaml:"created"   json:"created"`
	Cluster   string    `yaml:"cluster"   json:"cluster"`
	Namespace string    `yaml:"namespace" json:"namespace"`
	Pod       string    `yaml:"pod"       json:"pod"`
	Container string    `yaml:"container" json:"container"`
	Clues     []Clue    `yaml:"clues"     json:"clues"`
}

type Clue struct {
	Name string    `yaml:"name"    json:"name"`
	Content string `yaml:"content" json:"content"`
}
