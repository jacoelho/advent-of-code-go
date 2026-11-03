package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const firewallExample = "0: 3\n1: 2\n4: 4\n6: 4"

func Test_day13p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(firewallExample),
			Want:  "24",
		},
		{
			Input: strings.NewReader("0: 1\n3: 1"),
			Want:  "3",
		},
		{
			Input: aoc.FileInput(t, 2017, 13),
			Want:  "2164",
		},
	}

	aoc.AOCTest(t, day13p01, tests)
}

func Test_day13p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(firewallExample),
			Want:  "10",
		},
		{
			Input: strings.NewReader("0: 2"),
			Want:  "1",
		},
		{
			Input: aoc.FileInput(t, 2017, 13),
			Want:  "3861798",
		},
	}

	aoc.AOCTest(t, day13p02, tests)
}

func TestFirewallErrors(t *testing.T) {
	for _, input := range []string{"", "0: 0", "0: 1", "0: 2\n1: 2", "0: 3\n0: 4", "x: 3"} {
		if _, err := day13p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

func TestDay13ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"13/1", day13p01},
		{"13/2", day13p02},
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
