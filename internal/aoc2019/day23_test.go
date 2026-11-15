package aoc2019

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day23p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 23),
			Want:  "21897",
		},
	}
	aoc.AOCTest(t, day23p01, tests)
}

func Test_day23p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 23),
			Want:  "16424",
		},
	}
	aoc.AOCTest(t, day23p02, tests)
}

func TestNetworkProtocolErrors(t *testing.T) {
	for _, program := range []string{"99", "104,1,99", "104,7000,104,1,104,2,3,0,1105,1,6"} {
		for _, solve := range []func(io.Reader) (string, error){day23p01, day23p02} {
			if _, err := solve(strings.NewReader(program)); err == nil {
				t.Errorf("expected network error for %s", program)
			}
		}
	}
}
