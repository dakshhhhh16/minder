// SPDX-FileCopyrightText: Copyright 2026 The Minder Authors
// SPDX-License-Identifier: Apache-2.0

package ptr_test

import (
	"testing"

	"github.com/mindersec/minder/internal/util/ptr"
)

func TestPtr_Int(t *testing.T) {
	t.Parallel()
	v := 42
	got := ptr.Ptr(v)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *got != v {
		t.Errorf("*Ptr(%d) = %d, want %d", v, *got, v)
	}
}

func TestPtr_String(t *testing.T) {
	t.Parallel()
	s := "hello"
	got := ptr.Ptr(s)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *got != s {
		t.Errorf("*Ptr(%q) = %q, want %q", s, *got, s)
	}
}

func TestPtr_Bool(t *testing.T) {
	t.Parallel()
	for _, b := range []bool{true, false} {
		b := b
		t.Run("bool", func(t *testing.T) {
			t.Parallel()
			got := ptr.Ptr(b)
			if got == nil {
				t.Fatal("expected non-nil pointer")
			}
			if *got != b {
				t.Errorf("*Ptr(%v) = %v, want %v", b, *got, b)
			}
		})
	}
}

func TestPtr_ZeroValue(t *testing.T) {
	t.Parallel()
	got := ptr.Ptr(0)
	if got == nil {
		t.Fatal("expected non-nil pointer for zero value")
	}
	if *got != 0 {
		t.Errorf("*Ptr(0) = %d, want 0", *got)
	}
}

func TestPtr_Float64(t *testing.T) {
	t.Parallel()
	v := 3.14
	got := ptr.Ptr(v)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *got != v {
		t.Errorf("*Ptr(%f) = %f, want %f", v, *got, v)
	}
}

func TestPtr_Struct(t *testing.T) {
	t.Parallel()
	type point struct{ X, Y int }
	v := point{X: 1, Y: 2}
	got := ptr.Ptr(v)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *got != v {
		t.Errorf("*Ptr(%v) = %v, want %v", v, *got, v)
	}
	// Mutating the original should not affect the pointed-to value
	v.X = 99
	if got.X == 99 {
		t.Error("mutating original unexpectedly changed pointed-to struct")
	}
}

func TestPtr_ReturnsDifferentAddress(t *testing.T) {
	t.Parallel()
	v := 7
	p1 := ptr.Ptr(v)
	p2 := ptr.Ptr(v)
	if p1 == p2 {
		t.Error("expected two distinct pointer addresses for two Ptr calls")
	}
}
