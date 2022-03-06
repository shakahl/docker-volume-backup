// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const storageIDS3 storageID = "S3"

type s3Storage struct {
	client     *minio.Client
	bucketName string
	path       string
}

func newS3Storage(config *Config) (storage, error) {
	var creds *credentials.Credentials
	if config.AwsAccessKeyID != "" && config.AwsSecretAccessKey != "" {
		creds = credentials.NewStaticV4(
			config.AwsAccessKeyID,
			config.AwsSecretAccessKey,
			"",
		)
	} else if config.AwsIamRoleEndpoint != "" {
		creds = credentials.NewIAM(config.AwsIamRoleEndpoint)
	} else {
		return nil, errors.New("newS3Storage: AWS_S3_BUCKET_NAME is defined, but no credentials were provided")
	}

	options := minio.Options{
		Creds:  creds,
		Secure: config.AwsEndpointProto == "https",
	}

	if config.AwsEndpointInsecure {
		if !options.Secure {
			return nil, errors.New("newS3Storage: AWS_ENDPOINT_INSECURE = true is only meaningful for https")
		}

		transport, err := minio.DefaultTransport(true)
		if err != nil {
			return nil, fmt.Errorf("newS3Storage: failed to create default minio transport")
		}
		transport.TLSClientConfig.InsecureSkipVerify = true
		options.Transport = transport
	}

	mc, err := minio.New(config.AwsEndpoint, &options)
	if err != nil {
		return nil, fmt.Errorf("newS3Storage: error setting up minio client: %w", err)
	}
	return &s3Storage{
		client:     mc,
		bucketName: config.AwsS3BucketName,
		path:       config.AwsS3Path,
	}, nil
}

func (s *s3Storage) id() storageID {
	return storageIDS3
}

func (s *s3Storage) copy(file string) error {
	_, name := path.Split(file)
	if _, err := s.client.FPutObject(context.Background(), s.bucketName, filepath.Join(s.path, name), file, minio.PutObjectOptions{
		ContentType: "application/tar+gzip",
	}); err != nil {
		return fmt.Errorf("copy: error uploading backup to remote storage: %w", err)
	}
	return nil
}

func (s *s3Storage) list(prefix string) ([]backupInfo, error) {
	candidates := s.client.ListObjects(context.Background(), s.bucketName, minio.ListObjectsOptions{
		WithMetadata: true,
		Prefix:       prefix,
	})
	var result []backupInfo
	for candidate := range candidates {
		result = append(result, backupInfo{
			filename: candidate.Key,
			mtime:    candidate.LastModified,
		})
	}
	return result, nil
}

func (s *s3Storage) delete(file string) error {
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		objectsCh <- minio.ObjectInfo{Key: file}
		close(objectsCh)
	}()
	errChan := s.client.RemoveObjects(context.Background(), s.bucketName, objectsCh, minio.RemoveObjectsOptions{})
	for result := range errChan {
		if result.Err != nil {
			return fmt.Errorf("delete: error deleting file: %w", result.Err)
		}
	}
	return nil
}

func (s *s3Storage) symlink(string) error {
	return errNotSupported
}
