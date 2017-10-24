package storage

import (
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	metaKeyCluster      = "Cluster"
	metaKeyNamespace    = "Namespace"
	metaKeyPod          = "Pod"
	metaKeyContainer    = "Container"
	metaKeyIncident     = "Incident"
	metaKeyIncidentType = "Incidenttype"
)

// Client is used for interacting with the storage.
type Client struct {
	bucket   string
	uploader *s3manager.Uploader
	client   *s3.S3
}

// Incidents is a list of incidents.
type Incidents map[string]Incident

// Incident occurs when a pod is restarted.
type Incident struct {
	File      string
	Type      string
	Created   time.Time
	Cluster   string
	Namespace string
	Pod       string
	Container string
}
