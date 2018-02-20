package s3

import (
	"fmt"
	"encoding/json"
	"strings"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"

	"github.com/nickschuch/sherlock/storage/types"
)

const (
	StorageName = "s3"
	MetaKeyCluster      = "Cluster"
	MetaKeyNamespace    = "Namespace"
	MetaKeyPod          = "Pod"
	MetaKeyContainer    = "Container"
)

// New returns a new storage client.
func New(region, bucket string) (Storage, error) {
	sess := session.New(&aws.Config{Region: aws.String(region)})

	return Storage{
		bucket:   bucket,
		uploader: s3manager.NewUploader(sess),
		client:   s3.New(sess),
	}, nil
}

type Storage struct {
	bucket   string
	uploader *s3manager.Uploader
	client   *s3.S3
}

func (s Storage) Put(params types.PutParams) (types.PutResponse, error) {
	var resp types.PutResponse

	// @todo, Validate.

	// Convert to json.
	content, err := json.Marshal(params.Incident)
	if err != nil {
		return resp, errors.Wrap(err, "failed to marshal incident to json")
	}

	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Body:   strings.NewReader(string(content)),
		Bucket: aws.String(s.bucket),
		Key:    aws.String(params.Incident.ID),
		ACL:    aws.String("private"),
		Metadata: map[string]*string{
			MetaKeyCluster:      aws.String(params.Incident.Cluster),
			MetaKeyNamespace:    aws.String(params.Incident.Namespace),
			MetaKeyPod:          aws.String(params.Incident.Pod),
			MetaKeyContainer:    aws.String(params.Incident.Container),
		},
	})
	if err != nil {
		return resp, errors.Wrap(err, "failed to PUT clue")
	}

	return resp, nil
}

func (s Storage) Inspect(params types.InspectParams) (types.InspectResponse, error) {
	var resp types.InspectResponse

	detail, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(params.ID),
	})
	if err != nil {
		return resp, errors.Wrap(err, "failed to lookup incident")
	}

	content, err := ioutil.ReadAll(detail.Body)
	if err != nil {
		return resp, errors.Wrap(err, "failed to read incident")
	}

	var incident types.Incident

	err = json.Unmarshal(content, &incident)
	if err != nil {
		return resp, errors.Wrap(err, "failed to unmarshal incident")
	}

	resp.Incident = incident

	return resp, nil
}

func (s Storage) List(params types.ListParams) (types.ListResponse, error) {
	var resp types.ListResponse

	result, err := s.client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return resp, errors.Wrap(err, "failed to list objects")
	}

	for _, object := range result.Contents {
		created, cluster, namespace, pod, container, err := lookupMetadata(s.client, s.bucket, *object.Key)
		if err != nil {
			return resp, errors.Wrap(err, "failed to get incident metadata")
		}

		resp.Incidents = append(resp.Incidents, types.Incident{
			ID: *object.Key,
			Created: created,
			Cluster: cluster,
			Namespace: namespace,
			Pod: pod,
			Container: container,
		})
	}

	return resp, nil
}


func lookupMetadata(client *s3.S3, bucket, key string) (time.Time, string, string, string, string, error) {
	var (
		cluster      *string
		namespace    *string
		pod          *string
		container    *string
		ok           bool
	)

	head, err := client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return time.Now(), "", "", "", "", fmt.Errorf("failed to lookup object")
	}

	// Check for Cluster metadata.
	if cluster, ok = head.Metadata[MetaKeyCluster]; !ok {
		return time.Now(), "", "", "", "", fmt.Errorf("cannot find object with metadata: %s", MetaKeyCluster)
	}

	// Check for Namespace metadata.
	if namespace, ok = head.Metadata[MetaKeyNamespace]; !ok {
		return time.Now(), "", "", "", "", fmt.Errorf("cannot find object with metadata: %s", MetaKeyNamespace)
	}

	// Check for Pod metadata.
	if pod, ok = head.Metadata[MetaKeyPod]; !ok {
		return time.Now(), "", "", "", "", fmt.Errorf("cannot find object with metadata: %s", MetaKeyPod)
	}

	// Check for Container metadata.
	if container, ok = head.Metadata[MetaKeyContainer]; !ok {
		return time.Now(), "", "", "", "", fmt.Errorf("cannot find object with metadata: %s", MetaKeyContainer)
	}

	return *head.LastModified, *cluster, *namespace, *pod, *container, nil
}
