package storage

import (
	"github.com/pkg/errors"

	"github.com/nickschuch/sherlock/storage/s3"
	"github.com/nickschuch/sherlock/storage/types"
)

// New returns a new storage backend.
func New(name, region, bucket string) (types.Storage, error) {
	if name == s3.StorageName {
		return s3.New(region, bucket)
	}

	return nil, errors.New("cannot find storage backend")
}