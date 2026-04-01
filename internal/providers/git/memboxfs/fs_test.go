// SPDX-FileCopyrightText: Copyright 2026 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

package memboxfs

import (
	"errors"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
)

func newLimitedFs(maxFiles, maxBytes int64) *LimitedFs {
	return &LimitedFs{
		Fs:            memfs.New(),
		MaxFiles:      maxFiles,
		TotalFileSize: maxBytes,
	}
}

func TestLimitedFs_CreateUnderLimit(t *testing.T) {
	t.Parallel()
	fs := newLimitedFs(10, 1024)
	f, err := fs.Create("file.txt")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if f == nil {
		t.Error("Create() returned nil file")
	}
}

func TestLimitedFs_CreateAtFileLimit_ReturnsError(t *testing.T) {
	t.Parallel()
	fs := newLimitedFs(2, 1024*1024)
	if _, err := fs.Create("first.txt"); err != nil {
		t.Fatalf("first Create() error = %v", err)
	}
	_, err := fs.Create("second.txt")
	if err == nil {
		t.Error("Create() expected error at file limit")
	}
	if !errors.Is(err, ErrTooManyFiles) {
		t.Errorf("Create() error = %v, want ErrTooManyFiles", err)
	}
}

func TestLimitedFs_WriteExceedsTotalSize(t *testing.T) {
	t.Parallel()
	fs := newLimitedFs(10, 4)
	f, err := fs.Create("small.txt")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	_, err = f.Write([]byte("hello world")) // > 4 bytes
	if err == nil {
		t.Error("Write() expected error exceeding total size")
	}
	if !errors.Is(err, ErrTooBig) {
		t.Errorf("Write() error = %v, want ErrTooBig", err)
	}
}

func TestLimitedFs_WriteWithinSize(t *testing.T) {
	t.Parallel()
	fs := newLimitedFs(10, 1024)
	f, err := fs.Create("file.txt")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	n, err := f.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if n != 5 {
		t.Errorf("Write() n = %d, want 5", n)
	}
}

func TestLimitedFs_Chroot_NotImplemented(t *testing.T) {
	t.Parallel()
	fs := newLimitedFs(10, 1024)
	_, err := fs.Chroot("/")
	if !errors.Is(err, ErrNotImplemented) {
		t.Errorf("Chroot() error = %v, want ErrNotImplemented", err)
	}
}

func TestLimitedFs_Join_Delegates(t *testing.T) {
	t.Parallel()
	fs := newLimitedFs(10, 1024)
	got := fs.Join("a", "b", "c")
	if got == "" {
		t.Error("Join() returned empty string")
	}
}

func TestLimitedFs_ErrorConstants_Distinct(t *testing.T) {
	t.Parallel()
	if ErrTooBig == ErrTooManyFiles || ErrTooBig == ErrNotImplemented {
		t.Error("error constants should be distinct")
	}
}

func TestLimitedFs_MkdirAll_AtFileLimit(t *testing.T) {
	t.Parallel()
	fs := newLimitedFs(2, 1024*1024)
	if err := fs.MkdirAll("firstdir", 0755); err != nil {
		t.Fatalf("first MkdirAll() error = %v", err)
	}
	err := fs.MkdirAll("seconddir", 0755)
	if err == nil {
		t.Error("MkdirAll() expected error at file limit")
	}
	if !errors.Is(err, ErrTooManyFiles) {
		t.Errorf("MkdirAll() error = %v, want ErrTooManyFiles", err)
	}
}
