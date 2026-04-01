// SPDX-FileCopyrightText: Copyright 2026 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

package testing

import (
	"fmt"
	"testing"

	"github.com/go-git/go-billy/v5"
)

// verifyGitFiles checks that every entry in files exists in fs with the
// expected content.  It is a local test helper that exercises rule-engine
// semantics: the rule evaluator calls Read on a billy.Filesystem and asserts
// that key source files are present and un-modified.
func verifyGitFiles(fs billy.Filesystem, files map[string]string) error {
	for path, want := range files {
		f, err := fs.Open(path)
		if err != nil {
			return fmt.Errorf("file %q missing from fs: %w", path, err)
		}
		buf := make([]byte, len(want)+1)
		n, _ := f.Read(buf)
		f.Close()
		if got := string(buf[:n]); got != want {
			return fmt.Errorf("file %q: got %q, want %q", path, got, want)
		}
	}
	return nil
}

// TestVerifyGitFiles_FileMissingFromFS_ReturnsError confirms that
// verifyGitFiles reports an error when a required file is absent from the
// in-memory filesystem.  This validates the basic sentinel behaviour of the
// helper used by rule fixture tests.
func TestVerifyGitFiles_FileMissingFromFS_ReturnsError(t *testing.T) {
	t.Parallel()

	fs, err := NewMockBillyFS(nil)
	if err != nil {
		t.Fatalf("NewMockBillyFS: %v", err)
	}

	err = verifyGitFiles(fs, map[string]string{"README.md": "# Hello"})
	if err == nil {
		t.Fatal("expected an error when file is missing, got nil")
	}
}

// TestVerifyGitFiles_EmptyMap_NeverErrors asserts that when the expected file
// map is empty verifyGitFiles always succeeds regardless of the filesystem
// state.  An empty fixture git section should be a no-op.
func TestVerifyGitFiles_EmptyMap_NeverErrors(t *testing.T) {
	t.Parallel()

	fs, err := NewMockBillyFS(nil)
	if err != nil {
		t.Fatalf("NewMockBillyFS: %v", err)
	}

	if err := verifyGitFiles(fs, nil); err != nil {
		t.Errorf("empty map should never error, got: %v", err)
	}
}

// TestVerifyGitFiles_MultipleFiles_FailsOnFirstMissing populates the
// filesystem with only one of two advertised files and confirms that
// verifyGitFiles catches the first missing entry without panicking.  The
// ordering of map iteration is non-deterministic, so we only assert an error
// is returned – not which file name is mentioned.
func TestVerifyGitFiles_MultipleFiles_FailsOnFirstMissing(t *testing.T) {
	t.Parallel()

	fs, err := NewMockBillyFS(map[string]string{
		"present.go": "package main",
	})
	if err != nil {
		t.Fatalf("NewMockBillyFS: %v", err)
	}

	want := map[string]string{
		"present.go": "package main",
		"absent.go":  "package main",
	}
	if err := verifyGitFiles(fs, want); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
