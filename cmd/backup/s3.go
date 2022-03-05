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
	client *minio.Client
	config *Config
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
		client: mc,
		config: config,
	}, nil
}

func (s *s3Storage) id() storageID {
	return storageIDS3
}

func (s *s3Storage) copy(files []string) (msgs []string, errors []error) {
	for _, file := range files {
		_, name := path.Split(file)
		if _, err := s.client.FPutObject(context.Background(), s.config.AwsS3BucketName, filepath.Join(s.config.AwsS3Path, name), file, minio.PutObjectOptions{
			ContentType: "application/tar+gzip",
		}); err != nil {
			errors = append(errors, fmt.Errorf("copy: error uploading backup to remote storage: %w", err))
			continue
		}
		msgs = append(msgs, fmt.Sprintf("Uploaded a copy of backup `%s` to bucket `%s`.", file, s.config.AwsS3BucketName))
	}
	return
}

func (s *s3Storage) delete(files []string) ([]string, []error) {
	return nil, nil
}

func (s *s3Storage) list() ([]backupInfo, error) {
	return nil, nil
}
