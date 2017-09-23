package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func lookup(client *s3.S3, bucket string, key *string) (string, Incident, error) {
	var (
		cluster      *string
		namespace    *string
		pod          *string
		container    *string
		incidentType *string
		id           *string
		ok           bool
		incident     Incident
	)

	head, err := client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    key,
	})
	if err != nil {
		return "", incident, err
	}

	// Check for Cluster metadata.
	if cluster, ok = head.Metadata[metaKeyCluster]; !ok {
		return "", incident, fmt.Errorf("cannot find object with metadata: %s", metaKeyCluster)
	}

	// Check for Namespace metadata.
	if namespace, ok = head.Metadata[metaKeyNamespace]; !ok {
		return "", incident, fmt.Errorf("cannot find object with metadata: %s", metaKeyNamespace)
	}

	// Check for Pod metadata.
	if pod, ok = head.Metadata[metaKeyPod]; !ok {
		return "", incident, fmt.Errorf("cannot find object with metadata: %s", metaKeyPod)
	}

	// Check for Container metadata.
	if container, ok = head.Metadata[metaKeyContainer]; !ok {
		return "", incident, fmt.Errorf("cannot find object with metadata: %s", metaKeyContainer)
	}

	// Check for Incident Type metadata.
	if incidentType, ok = head.Metadata[metaKeyIncidentType]; !ok {
		return "", incident, fmt.Errorf("cannot find object with metadata: %s", metaKeyIncidentType)
	}

	// Check for ID metadata.
	if id, ok = head.Metadata[metaKeyIncident]; !ok {
		return "", incident, fmt.Errorf("cannot find object with metadata: %s", metaKeyIncident)
	}

	incident.File = *key
	incident.Type = *incidentType
	incident.Created = *head.LastModified
	incident.Cluster = *cluster
	incident.Namespace = *namespace
	incident.Pod = *pod
	incident.Container = *container

	return *id, incident, nil
}
