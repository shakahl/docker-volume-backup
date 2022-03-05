package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const storageIDLocal = "Local"

func newLocalStorage(config *Config) (storage, error) {
	return &localStorage{
		config: config,
	}, nil
}

type localStorage struct {
	config *Config
}

func (s *localStorage) id() storageID {
	return storageIDLocal
}

func (s *localStorage) copy(files []string) (errors []error) {
	for _, file := range files {

		_, name := path.Split(file)
		if err := copyFile(file, path.Join(s.config.BackupArchive, name)); err != nil {
			errors = append(errors, fmt.Errorf("copy: error copying file to local archive: %w", err))
			continue
		}
		if s.config.BackupLatestSymlink != "" {
			symlink := path.Join(s.config.BackupArchive, s.config.BackupLatestSymlink)
			if _, err := os.Lstat(symlink); err == nil {
				os.Remove(symlink)
			}
			if err := os.Symlink(name, symlink); err != nil {
				errors = append(errors, fmt.Errorf("copy: error creating latest symlink: %w", err))
			}
		}
	}
	return
}

func (s *localStorage) list(prefix string) ([]backupInfo, error) {
	globPattern := path.Join(
		s.config.BackupArchive,
		fmt.Sprintf("%s*", prefix),
	)
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

func (s *localStorage) delete(files []string) (errors []error) {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			errors = append(errors, err)
		}
	}
	return
}
