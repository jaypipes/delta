package cmp_test

import (
	"testing"

	gogocmp "github.com/google/go-cmp/cmp"
	"github.com/jaypipes/delta/cmp"
)

func BenchmarkDeltaEqual_SliceString_NoOptions(b *testing.B) {
	s1 := []string{
		"foo",
		"bar",
	}
	s2 := []string{
		"bar",
		"baz",
	}
	for b.Loop() {
		_ = cmp.Equal(s1, s2)
	}
}

func BenchmarkDeltaEqual_SliceString_NoCast(b *testing.B) {
	s1 := []string{
		"foo",
		"bar",
	}
	s2 := []string{
		"bar",
		"baz",
	}
	for b.Loop() {
		_ = cmp.Equal(s1, s2, cmp.NoCast)
	}
}

func BenchmarkGoGoCmpEqual_SliceString(b *testing.B) {
	s1 := []string{
		"foo",
		"bar",
	}
	s2 := []string{
		"bar",
		"baz",
	}
	for b.Loop() {
		_ = gogocmp.Equal(s1, s2)
	}
}
