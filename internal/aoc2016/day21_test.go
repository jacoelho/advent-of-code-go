package aoc2016

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const scrambleExample = "swap position 4 with position 0\nswap letter d with letter b\nreverse positions 0 through 4\nrotate left 1 step\nmove position 1 to position 4\nmove position 3 to position 0\nrotate based on position of letter b\nrotate based on position of letter d\n"

func Test_day21p01(t *testing.T) {
	aoc.AOCTest(t, day21p01, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 21), Want: "gbhafcde"}})
}

func Test_day21p02(t *testing.T) {
	aoc.AOCTest(t, day21p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 21), Want: "bcfaegdh"}})
}

func TestPasswordScrambling(t *testing.T) {
	instructions, err := parseScramble(strings.NewReader(scrambleExample))
	if err != nil {
		t.Fatal(err)
	}
	got, err := scramblePassword(instructions, "abcde")
	if err != nil || got != "decab" {
		t.Fatalf("scrambled=%q,%v; wantdecab", got, err)
	}
	for _, tc := range []struct{ operation, want string }{
		{"swap position 1 with position 5", "afcdebgh"},
		{"swap letter b with letter f", "afcdebgh"},
		{"rotate left 2 steps", "cdefghab"},
		{"rotate right 3 steps", "fghabcde"},
		{"reverse positions 2 through 6", "abgfedch"},
		{"move position 1 to position 6", "acdefgbh"},
		{"move position 6 to position 1", "agbcdefh"},
	} {
		program, err := parseScramble(strings.NewReader(tc.operation))
		if err != nil {
			t.Fatal(err)
		}
		got, err := scramblePassword(program, "abcdefgh")
		if err != nil || got != tc.want {
			t.Errorf("%s=%q,%v; want%q", tc.operation, got, err, tc.want)
		}
		previous, err := unscramblePassword(program, tc.want)
		if err != nil || previous != "abcdefgh" {
			t.Errorf("inverse%s=%q,%v", tc.operation, previous, err)
		}
	}
	for i, want := range []string{"habcdefg", "ghabcdef", "fghabcde", "efghabcd", "cdefghab", "bcdefgha", "abcdefgh", "habcdefg"} {
		program, err := parseScramble(strings.NewReader(fmt.Sprintf("rotate based on position of letter %c", 'a'+i)))
		if err != nil {
			t.Fatal(err)
		}
		got, err := scramblePassword(program, "abcdefgh")
		if err != nil || got != want {
			t.Errorf("rotate%c=%q,%v; want%q", 'a'+i, got, err, want)
		}
		previous, err := unscramblePassword(program, want)
		if err != nil || previous != "abcdefgh" {
			t.Errorf("inverse rotate%c=%q,%v", 'a'+i, previous, err)
		}
	}
}

func TestDay21ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day21p01,
		day21p02,
	} {
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("read error%v lost input cause", err)
		}
	}
}
