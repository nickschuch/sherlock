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

type Client struct {
	bucket   string
	uploader *s3manager.Uploader
	client   *s3.S3
}

type Incidents map[string]Incident

type Incident struct {
	File      string
	Type      string
	Created   time.Time
	Cluster   string
	Namespace string
	Pod       string
	Container string
}
