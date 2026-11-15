package aoc2019

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("3,0,3,1,104,-1,99"), Want: "-1"},
		{
			Input: aoc.FileInput(t, 2019, 7),
			Want:  "437860",
		},
	}
	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("3,0,3,1,104,-1,99"), Want: "-1"},
		{
			Input: strings.NewReader("3,26,1001,26,-4,26,3,27,1002,27,2,27,1,27,26,27,4,27,1001,28,-1,28,1005,28,6,99,0,0,5"),
			Want:  "139629729",
		},
		{
			Input: aoc.FileInput(t, 2019, 7),
			Want:  "49810599",
		},
	}
	aoc.AOCTest(t, day07p02, tests)
}

func TestAmplifierRejectsMissingOutput(t *testing.T) {
	for _, solve := range []func(io.Reader) (string, error){day07p01, day07p02} {
		if _, err := solve(strings.NewReader("3,0,3,1,99")); err == nil {
			t.Fatal("expected missing amplifier output error")
		}
	}
}
