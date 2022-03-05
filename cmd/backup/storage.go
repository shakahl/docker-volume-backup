package main

import (
	"time"
)

type storage interface {
	id() storageID
	list(prefix string) ([]backupInfo, error)
	copy(files []string) ([]string, []error)
	delete(files []string) ([]string, []error)
}

type storageID string

type backupInfo struct {
	filename string
	mtime    time.Time
}
