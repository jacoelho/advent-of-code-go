package aoc2016

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	aoc.AOCTest(t, day25p01, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 25), Want: "198"}, {Input: strings.NewReader("out 0\nout 1\njnz 1 -2"), Want: "1"}})
}

func TestClockNeedsInfiniteAlternation(t *testing.T) {
	for _, tc := range []struct {
		program string
		want    bool
	}{
		{"out 0\nout 1\njnz 1 -2", true},
		{"out 0\njnz 1 -1", false},
		{"out 1\nout 0\njnz 1 -2", false},
		{"out 0\nout 1", false},
		{"out 0\nout 1\njnz 1 0", false},
		{"out 0\nout 1\ntgl -2\njnz 1 -3", false},
		{"cpy 600 b\nout 0\nout 1\ndec b\njnz b -3", false},
	} {
		program, err := parseAssembunny(strings.NewReader(tc.program))
		if err != nil {
			t.Fatal(err)
		}
		remaining := 10_000
		got, err := isClockSignal(program, 1, &remaining)
		if err != nil || got != tc.want {
			t.Errorf("%q: clock=%v,%v; want%v", tc.program, got, err, tc.want)
		}
	}
}

func TestClockTotalBudget(t *testing.T) {
	for _, input := range []string{"", "out 1"} {
		program, err := parseAssembunny(strings.NewReader(input))
		if err != nil {
			t.Fatal(err)
		}
		if got, err := clockSeed(program, 20); err == nil || got != 0 {
			t.Errorf("unbounded candidate search returned%d,%v", got, err)
		}
	}
	program, err := parseAssembunny(strings.NewReader("cpy a b\ndec b\njnz b 3\ninc c\njnz 1 -1\nout 0\nout 1\njnz 1 -2"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := clockSeed(program, 20); err == nil || !strings.HasSuffix(err.Error(), "candidate 1") {
		t.Fatalf("skipped undecided first candidate: %v", err)
	}
}
