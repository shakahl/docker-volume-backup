// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/studio-b12/gowebdav"
)

const storageIDWebDAV = "WebDAV"

func newWebDAVStorage(url, username, password, path string) (storage, error) {
	if username == "" || password == "" {
		return nil, errors.New("newWebDAVStorage: WEBDAV_URL is defined, but no credentials were provided")
	}
	return &webDAVStorage{
		client: gowebdav.NewClient(url, username, password),
		path:   path,
	}, nil
}

type webDAVStorage struct {
	client *gowebdav.Client
	path   string
}

func (s *webDAVStorage) id() storageID {
	return storageIDWebDAV
}

func (s *webDAVStorage) copy(file string) error {
	_, name := path.Split(file)
	bytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("copy: error reading the file to be uploaded: %w", err)
	}
	if err := s.client.MkdirAll(s.path, 0644); err != nil {
		return fmt.Errorf("copy: error creating directory '%s' on WebDAV server: %w", s.path, err)
	}
	if err := s.client.Write(filepath.Join(s.path, name), bytes, 0644); err != nil {
		return fmt.Errorf("copy: error uploading the file to WebDAV server: %w", err)
	}
	return nil
}

func (s *webDAVStorage) list(prefix string) ([]backupInfo, error) {
	candidates, err := s.client.ReadDir(s.path)
	if err != nil {
		return nil, fmt.Errorf("list: error looking up candidates from remote storage: %w", err)
	}
	var matches []backupInfo
	for _, candidate := range candidates {
		if strings.HasPrefix(candidate.Name(), prefix) {
			matches = append(matches, backupInfo{
				filename: candidate.Name(),
				mtime:    candidate.ModTime(),
			})
		}
	}
	return matches, nil
}

func (s *webDAVStorage) delete(file string) error {
	if err := s.client.Remove(filepath.Join(s.path, file)); err != nil {
		return fmt.Errorf("delete: error removing file from WebDAV storage: %w", err)
	}
	return nil
}

func (s *webDAVStorage) symlink(string) error {
	return errNotSupported
}
