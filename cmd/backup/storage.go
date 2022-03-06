// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"errors"
	"time"
)

type storage interface {
	id() storageID
	list(prefix string) ([]backupInfo, error)
	symlink(file string) error
	copy(file string) error
	delete(file string) error
}

type storageID string

type backupInfo struct {
	filename string
	mtime    time.Time
}

var errNotSupported = errors.New("method not supported")
