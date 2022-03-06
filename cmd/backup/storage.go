// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"time"
)

type storage interface {
	id() storageID
	list(prefix string) ([]backupInfo, error)
	copy(files []string) []error
	delete(files []string) []error
}

type storageID string

type backupInfo struct {
	filename string
	mtime    time.Time
}
