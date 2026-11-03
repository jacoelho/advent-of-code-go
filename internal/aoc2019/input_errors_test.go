package aoc2019

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func TestIntcodeReaderError(t *testing.T) {
	failure := errors.New("input read failed")
	reader := io.MultiReader(strings.NewReader("99"), iotest.ErrReader(failure))
	if _, err := ParseIntcodeInput(reader); !errors.Is(err, failure) {
		t.Fatalf("got %v, want input failure", err)
	}
}

func TestPuzzleReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name, prefix string
		solve        func(io.Reader) (string, error)
	}{
		{"day11p01", "99", day11p01}, {"day11p02", "99", day11p02},
		{"day13p01", "99", day13p01}, {"day13p02", "99", day13p02},
		{"day15p01", "99", day15p01}, {"day15p02", "99", day15p02},
		{"day19p01", "99", day19p01}, {"day19p02", "99", day19p02},
		{"day20p01", "AA...ZZ\n", day20p01}, {"day20p02", "AA...ZZ\n", day20p02},
		{"day24p01", ".....\n.....\n.....\n.....\n.....\n", day24p01},
		{"day24p02", ".....\n.....\n.....\n.....\n.....\n", day24p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				reader := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(reader); !errors.Is(err, failure) {
					t.Fatalf("prefix %q: got %v, want input failure", prefix, err)
				}
			}
		})
	}
}
