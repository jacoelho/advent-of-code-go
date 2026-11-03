package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func TestTriangleAndRoomReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name, prefix string
		solve        func(io.Reader) (string, error)
	}{
		{"day03p01", "3 4 5\n3 4 5\n3 4 5\n", day03p01},
		{"day03p02", "3 4 5\n3 4 5\n3 4 5\n", day03p02},
		{"day04p01", "northpole-object-storage-0[oetra]\n", day04p01},
		{"day04p02", "northpole-object-storage-0[oetra]\n", day04p02},
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
