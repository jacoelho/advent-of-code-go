package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const bunnyToggleExample = "cpy 2 a\ntgl a\ntgl a\ntgl a\ncpy 1 a\ndec a\ndec a\n"

func Test_day23p01(t *testing.T) {
	aoc.AOCTest(t, day23p01, []aoc.TestInput{{Input: strings.NewReader(bunnyToggleExample), Want: "3"}, {Input: aoc.FileInput(t, 2016, 23), Want: "10880"}})
}

func Test_day23p02(t *testing.T) {
	aoc.AOCTest(t, day23p02, []aoc.TestInput{{Input: strings.NewReader(bunnyToggleExample), Want: "3"}, {Input: aoc.FileInput(t, 2016, 23), Want: "479007440"}})
}

func TestAssembunnySemantics(t *testing.T) {
	for _, tc := range []struct {
		program string
		want    [4]int64
	}{
		{"cpy 1 a\ntgl 1\njnz 1 2\ninc a", [4]int64{2}},
		{"tgl 1\ntgl 2\ninc a", [4]int64{1}},
		{"cpy 0 a\ntgl a\ninc a", [4]int64{1}},
		{"cpy 2 b\njnz 1 b\ninc a\ninc a", [4]int64{1, 2}},
		{"jnz 1 -1", [4]int64{}},
		{"cpy 2 a\ncpy 2 d\ncpy a c\ninc a\ndec c\njnz c -2\ndec d\njnz d -5", [4]int64{8, 0, 0, 0}},
		{"cpy 3 b\ncpy 2 d\ncpy b c\ninc a\ndec c\njnz c -2\ndec d\njnz d -5", [4]int64{6, 3, 0, 0}},
	} {
		program, err := parseAssembunny(strings.NewReader(tc.program))
		if err != nil {
			t.Fatal(err)
		}
		for range 2 {
			got, err := runAssembunny(program, [4]int64{})
			if err != nil || got != tc.want {
				t.Errorf("%q registers=%v,%v; want%v", tc.program, got, err, tc.want)
			}
		}
	}
	program, err := parseAssembunny(strings.NewReader("jnz 1 0"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runAssembunny(program, [4]int64{}); err == nil {
		t.Fatal("accepted nonterminating program")
	}
}

func TestAssembunnyNonpositiveCounters(t *testing.T) {
	for _, input := range []string{
		"cpy 0 b\ninc a\ndec b\njnz b -2",
		"cpy -1 b\ninc a\ndec b\njnz b -2",
		"cpy 2 d\ncpy 0 c\ninc a\ndec c\njnz c -2\ndec d\njnz d -5",
	} {
		program, err := parseAssembunny(strings.NewReader(input))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := executeAssembunny(program, [4]int64{}, 20); err == nil {
			t.Errorf("nonterminating loop incorrectly completed: %q", input)
		}
	}
}

func TestDay23ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day23p01,
		day23p02,
	} {
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("read error%v lost input cause", err)
		}
	}
}
