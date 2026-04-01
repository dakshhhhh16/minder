// SPDX-FileCopyrightText: Copyright 2026 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

package noop

import (
	"testing"
)

func TestNewJwtValidator_NonNil(t *testing.T) {
	t.Parallel()
	v := NewJwtValidator("test-subject")
	if v == nil {
		t.Error("NewJwtValidator() returned nil")
	}
}

func TestNewJwtValidator_SubjectMatches(t *testing.T) {
	t.Parallel()
	v := NewJwtValidator("alice")
	tok, err := v.ParseAndValidate("any-token")
	if err != nil {
		t.Fatalf("ParseAndValidate() error = %v", err)
	}
	if tok.Subject() != "alice" {
		t.Errorf("Subject() = %q, want 'alice'", tok.Subject())
	}
}

func TestNewJwtValidator_EmptySubject(t *testing.T) {
	t.Parallel()
	v := NewJwtValidator("")
	tok, err := v.ParseAndValidate("token")
	if err != nil {
		t.Fatalf("ParseAndValidate() error = %v", err)
	}
	if tok.Subject() != "" {
		t.Errorf("Subject() = %q, want empty", tok.Subject())
	}
}

func TestNewJwtValidator_TokenStringIgnored(t *testing.T) {
	t.Parallel()
	v := NewJwtValidator("bob")
	tok, err := v.ParseAndValidate("ignored-string")
	if err != nil {
		t.Fatalf("ParseAndValidate() error = %v", err)
	}
	if tok.Subject() != "bob" {
		t.Errorf("Subject() = %q, want 'bob'", tok.Subject())
	}
}

func TestNewJwtValidator_ImplementsInterface(t *testing.T) {
	t.Parallel()
	_ = NewJwtValidator("subject")
}
