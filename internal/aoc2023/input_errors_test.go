package aoc2023

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func Test2023ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"day03p01", day03p01, "1*2\n"},
		{"day03p02", day03p02, "1*2\n"},
		{"day15p01", day15p01, "rn=1\n"},
		{"day15p02", day15p02, "rn=1\n"},
		{"day20p01", day20p01, "broadcaster -> rx\n"},
		{"day20p02", day20p02, "broadcaster -> rx\n"},
		{"day21p01", day21p01, "S\n"},
		{"day21p02", day21p02, "S\n"},
		{"day22p01", day22p01, "0,0,1~0,0,1\n"},
		{"day22p02", day22p02, "0,0,1~0,0,1\n"},
		{"day23p01", day23p01, "#.\n#.\n"},
		{"day23p02", day23p02, "#.\n#.\n"},
		{"day24p01", day24p01, hailExample + "\n"},
		{"day24p02", day24p02, hailExample + "\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, mode := range []struct{ name, prefix string }{{"immediate", ""}, {"after prefix", tt.prefix}} {
				t.Run(mode.name, func(t *testing.T) {
					input := io.MultiReader(strings.NewReader(mode.prefix), iotest.ErrReader(failure))
					_, err := tt.solve(input)
					if !errors.Is(err, failure) {
						t.Fatalf("got %v, want reader failure", err)
					}
				})
			}
		})
	}
}
