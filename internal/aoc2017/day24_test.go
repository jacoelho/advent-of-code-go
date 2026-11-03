package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const bridgeExample = "0/2\n2/2\n2/3\n3/4\n3/5\n0/1\n10/1\n9/10"

func Test_day24p01(t *testing.T) {
	aoc.AOCTest(t, day24p01, []aoc.TestInput{
		{Input: strings.NewReader(bridgeExample), Want: "31"},
		{Input: strings.NewReader("0/1\n0/1\n1/1"), Want: "4"},
		{Input: aoc.FileInput(t, 2017, 24), Want: "1695"},
	})
}

func Test_day24p02(t *testing.T) {
	aoc.AOCTest(t, day24p02, []aoc.TestInput{
		{Input: strings.NewReader(bridgeExample), Want: "19"},
		{Input: strings.NewReader("0/1\n0/1\n1/1"), Want: "4"},
		{Input: aoc.FileInput(t, 2017, 24), Want: "1673"},
	})
}

func TestDay24ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"24/1", day24p01},
		{"24/2", day24p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
