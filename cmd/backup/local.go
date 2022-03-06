// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const storageIDLocal = "Local"

func newLocalStorage(archiveLocation, latestSymlink string) (storage, error) {
	return &localStorage{
		archiveLocation: archiveLocation,
		latestSymlink:   latestSymlink,
	}, nil
}

type localStorage struct {
	archiveLocation string
	latestSymlink   string
}

func (s *localStorage) id() storageID {
	return storageIDLocal
}

func (s *localStorage) copy(file string) error {
	_, name := path.Split(file)
	if err := copyFile(file, path.Join(s.archiveLocation, name)); err != nil {
		return fmt.Errorf("copy: error copying file to local archive: %w", err)
	}
	return nil
}

func (s *localStorage) symlink(file string) error {
	_, name := path.Split(file)
	symlink := path.Join(s.archiveLocation, s.latestSymlink)
	if _, err := os.Lstat(symlink); err == nil {
		os.Remove(symlink)
	}
	if err := os.Symlink(name, symlink); err != nil {
		return fmt.Errorf("symlink: error creating latest symlink: %w", err)
	}
	return nil
}

func (s *localStorage) list(prefix string) ([]backupInfo, error) {
	globPattern := path.Join(s.archiveLocation, fmt.Sprintf("%s*", prefix))
	globMatches, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf(
			"list: error looking up matching files using pattern %s: %w",
			globPattern,
			err,
		)
	}

	var candidates []backupInfo
	for _, candidate := range globMatches {
		fi, err := os.Lstat(candidate)
		if err != nil {
			return nil, fmt.Errorf(
				"list: error calling Lstat on file %s: %w",
				candidate,
				err,
			)
		}

		if fi.Mode()&os.ModeSymlink != os.ModeSymlink {
			candidates = append(candidates, backupInfo{
				filename: candidate,
				mtime:    fi.ModTime(),
			})
		}
	}
	return candidates, nil
}

func (s *localStorage) delete(file string) error {
	if err := os.Remove(file); err != nil {
		return fmt.Errorf("delete: error removing file: %w", err)
	}
	return nil
}
