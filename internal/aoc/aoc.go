package aoc

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
)

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func MustAtoi(s string) int {
	return Must(strconv.Atoi(s))
}

func Must2[T, V any](t T, v V, err error) (T, V) {
	if err != nil {
		panic(err)
	}
	return t, v
}

func FileInput(t *testing.T, year, day int) io.Reader {
	t.Helper()

	f, err := os.Open(fmt.Sprintf("../../inputs/%d/%02d.txt", year, day))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}

type TestInput struct {
	Input io.Reader
	Want  string
}

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
