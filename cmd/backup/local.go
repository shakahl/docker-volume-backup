package main

import (
	"fmt"
	"os"
	"path"
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

func (s *localStorage) copy(files []string) (messages []string, errors []error) {
	for _, file := range files {

		_, name := path.Split(file)
		if err := copyFile(file, path.Join(s.config.BackupArchive, name)); err != nil {
			errors = append(errors, fmt.Errorf("copy: error copying file to local archive: %w", err))
			continue
		}
		messages = append(messages, fmt.Sprintf("Stored copy of backup `%s` in local archive `%s`.", file, s.config.BackupArchive))
		if s.config.BackupLatestSymlink != "" {
			symlink := path.Join(s.config.BackupArchive, s.config.BackupLatestSymlink)
			if _, err := os.Lstat(symlink); err == nil {
				os.Remove(symlink)
			}
			if err := os.Symlink(name, symlink); err != nil {
				errors = append(errors, fmt.Errorf("copy: error creating latest symlink: %w", err))
			}
			messages = append(messages, fmt.Sprintf("Created/Updated symlink `%s` for latest backup.", s.config.BackupLatestSymlink))
		}
	}
	return
}

func (s *localStorage) delete(files []string) ([]string, []error) {
	return nil, nil
}

func (s *localStorage) list() ([]backupInfo, error) {
	return nil, nil
}
