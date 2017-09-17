package storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func New(region, bucket string) Client {
	var sess = session.New(&aws.Config{Region: aws.String(region)})

	return Client{
		bucket:   bucket,
		uploader: s3manager.NewUploader(sess),
		client:   s3.New(sess),
	}
}

func (s Client) Write(namespace, pod, container, incident, name string, data []byte) error {
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filepath.Join(namespace, pod, container, incident, name)),
		ACL:    aws.String("private"),
		Metadata: map[string]*string{
			metaKeyNamespace:    aws.String(namespace),
			metaKeyPod:          aws.String(pod),
			metaKeyContainer:    aws.String(container),
			metaKeyIncident:     aws.String(incident),
			metaKeyIncidentType: aws.String(name),
		},
	})
	return err
}

func (s Client) Incidents() (Incidents, error) {
	incidents := make(Incidents)

	result, err := s.client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return incidents, err
	}

	for _, file := range result.Contents {
		id, incident, err := lookup(s.client, s.bucket, file.Key)
		if err != nil {
			return incidents, err
		}

		incidents[id] = incident
	}

	return incidents, nil
}

func (s Client) IncidentDetails(incidentID string) (map[string]string, error) {
	files := make(map[string]string)

	result, err := s.client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return files, err
	}

	for _, file := range result.Contents {
		id, incident, err := lookup(s.client, s.bucket, file.Key)
		if err != nil {
			return files, err
		}

		if id != incidentID {
			continue
		}

		detail, err := s.client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(incident.File),
		})
		if err != nil {
			return files, fmt.Errorf("failed to lookup incident detail: %s", err)
		}

		content, err := ioutil.ReadAll(detail.Body)
		if err != nil {
			return files, fmt.Errorf("failed to read incident detail: %s", err)
		}

		files[incident.Type] = string(content)
	}

	return files, nil
}
