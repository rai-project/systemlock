// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd

package systemlock

import (
	"sync"
	"syscall"
)

const (
	mkdirPerm = 0750
)

// FileMutex is similar to sync.RWMutex, but also synchronizes across processes.
// This implementation is based on flock syscall.
type FileMutex struct {
	mu sync.RWMutex
	fd int
}

func New(filename string) (*FileMutex, error) {
	fd, err := syscall.Open(filename, syscall.O_CREAT|syscall.O_RDONLY, mkdirPerm)
	if err != nil {
		return nil, err
	}
	return &FileMutex{fd: fd}, nil
}

func (m *FileMutex) Lock() error {
	m.mu.Lock()
	return syscall.Flock(m.fd, syscall.LOCK_EX)
}

func (m *FileMutex) Unlock() error {
	if err := syscall.Flock(m.fd, syscall.LOCK_UN); err != nil {
		return err
	}
	m.mu.Unlock()
	return nil
}

func (m *FileMutex) RLock() error {
	m.mu.RLock()
	return syscall.Flock(m.fd, syscall.LOCK_SH)
}

func (m *FileMutex) RUnlock() error {
	if err := syscall.Flock(m.fd, syscall.LOCK_UN); err != nil {
		return err
	}
	m.mu.RUnlock()
	return nil
}
