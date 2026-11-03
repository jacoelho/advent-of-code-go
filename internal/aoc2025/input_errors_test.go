package aoc2025

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func Test2025ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"day01p01", day01p01, "R1\n"},
		{"day01p02", day01p02, "R1\n"},
		{"day03p01", day03p01, "123456789012\n"},
		{"day03p02", day03p02, "123456789012\n"},
		{"day06p01", day06p01, "12\n3\n+\n"},
		{"day06p02", day06p02, "12\n3\n+\n"},
		{"day07p01", day07p01, "S\n"},
		{"day07p02", day07p02, "S\n"},
		{"day09p01", day09p01, "0,0\n2,0\n2,1\n0,1\n"},
		{"day09p02", day09p02, "0,0\n2,0\n2,1\n0,1\n"},
		{"day10p01", day10p01, "[.] (0) {0}\n"},
		{"day10p02", day10p02, "[.] (0) {0}\n"},
		{"day12p01", day12p01, "0:\n#\n\n1x1: 1\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				input := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(input); !errors.Is(err, failure) {
					t.Fatalf("got %v; want reader failure", err)
				}
			}
		})
	}
}
