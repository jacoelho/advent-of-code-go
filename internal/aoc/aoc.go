// Package aoc provides input loading, test helpers and display decoding shared
// by the Advent of Code solutions.
package aoc

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
)

// Must returns t, panicking if err is non-nil.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// MustAtoi returns s parsed as a decimal int, panicking if it is invalid.
func MustAtoi(s string) int {
	return Must(strconv.Atoi(s))
}

// Must2 returns t and v, panicking if err is non-nil.
func Must2[T, V any](t T, v V, err error) (T, V) {
	if err != nil {
		panic(err)
	}
	return t, v
}

// TestHelper is the subset of testing.TB that FileInput uses.
type TestHelper interface {
	Helper()
	Cleanup(func())
	Fatal(...any)
}

// FileInput opens ../../inputs/<year>/<day>.txt, with a two-digit day,
// relative to the working directory. It fails the test if the file cannot be
// opened and closes it during test cleanup.
func FileInput(t TestHelper, year, day int) io.Reader {
	t.Helper()

	f, err := os.Open(fmt.Sprintf("../../inputs/%d/%02d.txt", year, day))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}

// TestInput is a puzzle input with its expected answer.
type TestInput struct {
	Input io.Reader
	Want  string
}

// AOCTest runs f on each input in its own subtest, reporting an error or an
// answer that differs from Want.
func AOCTest(t *testing.T, f func(io.Reader) (string, error), inputs []TestInput) {
	t.Helper()

	for i, tt := range inputs {
		t.Run(fmt.Sprintf("test %02d", i), func(t *testing.T) {
			got, err := f(tt.Input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.Want {
				t.Errorf("got = %v, want %v", got, tt.Want)
			}
		})
	}
}
