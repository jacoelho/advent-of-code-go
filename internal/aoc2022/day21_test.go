package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestMonkeyNoninvertibleDomain(t *testing.T) {
	input := `root: left + zero
left: zero * expr
expr: onee / half
half: plus / twoo
plus: quot + onee
quot: twoo / humn
zero: 0
onee: 1
twoo: 2
humn: 2`
	for _, operation := range []string{"*", "/"} {
		program := strings.Replace(input, "zero * expr", "zero "+operation+" expr", 1)
		_, err := day21p02(strings.NewReader(program))
		if err == nil || !strings.Contains(err.Error(), "unsupported") {
			t.Fatalf("zero %s expr: got %v; want unsupported equation error", operation, err)
		}
	}
}

func Test_day21p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(monkeyExample), Want: "152"},
		{Input: aoc.FileInput(t, 2022, 21), Want: "223971851179174"},
	}

	aoc.AOCTest(t, day21p01, tests)
}

func Test_day21p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(monkeyExample), Want: "301"},
		{Input: aoc.FileInput(t, 2022, 21), Want: "3379022190351"},
	}

	aoc.AOCTest(t, day21p02, tests)
}

func TestMonkeyGraphErrors(t *testing.T) {
	for _, input := range []string{
		"root: miss + numb\nnumb: 1",
		"root: loop + numb\nloop: root + numb\nnumb: 1",
		"root: 1\nroot: 2",
		"root: left / rght\nleft: 1\nrght: 0",
		"root: left / rght\nleft: 3\nrght: 2",
	} {
		if _, err := day21p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected monkey expression error for %q", input)
		}
	}
}

func TestMonkeyInversion(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("root: left + rght\nleft: cons - humn\ncons: 10\nrght: 3\nhumn: 0"), Want: "7"},
		{Input: strings.NewReader("root: left + rght\nleft: cons / humn\ncons: 12\nrght: 3\nhumn: 1"), Want: "4"},
		{Input: strings.NewReader("root: left + rght\nleft: humn / cons\ncons: 3\nrght: 4\nhumn: 0"), Want: "12"},
	}

	aoc.AOCTest(t, day21p02, tests)
}

func TestMonkeyEquationErrors(t *testing.T) {
	for _, input := range []string{
		"root: left + rght\nleft: humn * cons\ncons: 2\nrght: 3\nhumn: 0",
		"root: left + rght\nleft: cons / humn\ncons: 0\nrght: 1\nhumn: 1",
		"root: left + rght\nleft: cons * humn\ncons: 0\nrght: 0\nhumn: 1",
		"root: left + rght\nleft: cons * humn\ncons: 0\nrght: 1\nhumn: 1",
		"root: left + rght\nleft: humn + humn\nrght: 4\nhumn: 1",
	} {
		if _, err := day21p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected unsupported or unsatisfiable equation error for %q", input)
		}
	}
}

const monkeyExample = `root: pppw + sjmn
dbpl: 5
cczh: sllz + lgvd
zczc: 2
ptdq: humn - dvpt
dvpt: 3
lfqf: 4
humn: 5
ljgn: 2
sjmn: drzm * dbpl
sllz: 4
pppw: cczh / lfqf
lgvd: ljgn * ptdq
drzm: hmdt - zczc
hmdt: 32`

func TestDay21ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day21p01},
		{"part2", day21p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
