package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/studio-b12/gowebdav"
)

const storageIDWebDAV = "WebDAV"

func newWebDAVStorage(config *Config) (storage, error) {
	if config.WebdavUsername == "" || config.WebdavPassword == "" {
		return nil, errors.New("newWebDAVStorage: WEBDAV_URL is defined, but no credentials were provided")
	}
	return &webDAVStorage{
		config: config,
		client: gowebdav.NewClient(config.WebdavUrl, config.WebdavUsername, config.WebdavPassword),
	}, nil
}

type webDAVStorage struct {
	config *Config
	client *gowebdav.Client
}

func (s *webDAVStorage) id() storageID {
	return storageIDWebDAV
}

func (s *webDAVStorage) copy(files []string) (messages []string, errors []error) {
	for _, file := range files {
		_, name := path.Split(file)
		bytes, err := os.ReadFile(file)
		if err != nil {
			errors = append(errors, fmt.Errorf("copy: error reading the file to be uploaded: %w", err))
			continue
		}
		if err := s.client.MkdirAll(s.config.WebdavPath, 0644); err != nil {
			errors = append(errors, fmt.Errorf("copy: error creating directory '%s' on WebDAV server: %w", s.config.WebdavPath, err))
			continue
		}
		if err := s.client.Write(filepath.Join(s.config.WebdavPath, name), bytes, 0644); err != nil {
			errors = append(errors, fmt.Errorf("copy: error uploading the file to WebDAV server: %w", err))
			continue
		}
		messages = append(messages, fmt.Sprintf("Uploaded a copy of backup `%s` to WebDAV-URL '%s' at path '%s'.", file, s.config.WebdavUrl, s.config.WebdavPath))
	}
	return
}

func (s *webDAVStorage) delete(files []string) ([]string, []error) {
	return nil, nil
}

func (s *webDAVStorage) list() ([]backupInfo, error) {
	return nil, nil
}
