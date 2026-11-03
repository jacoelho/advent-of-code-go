package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day10Example = `[({(<(())[]>[[{[]{<()<>>
[(()[<>])]({[<{<<[]>>(
{([(<{}[<>[]}>{[]{[(<()>
(((({<>}<{<{<>}{[]{[]{}
[[<[([]))<([[{}[[()]]]
[{[{({}]{}}([{[{{{}}([]
{<[[]]>}<{[{[{[]{()[[[]
[<(<(<(<{}))><([]([]()
<{([([[(<>()){}]>(<<{{
<{([{{}}[<[[[<>{}]]]>[]]`

func Test_day10p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day10Example), Want: "26397"},
		{Input: aoc.FileInput(t, 2021, 10), Want: "370407"},
	}

	aoc.AOCTest(t, day10p01, tests)
}

func Test_day10p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day10Example), Want: "288957"},
		{Input: strings.NewReader(`()
(
((
[`), Want: "2"},
		{Input: aoc.FileInput(t, 2021, 10), Want: "3249889609"},
	}

	aoc.AOCTest(t, day10p02, tests)
}

func TestDay10InvalidDelimiter(t *testing.T) {
	if _, err := day10p01(strings.NewReader(`x`)); err == nil {
		t.Fatal("expected invalid delimiter error")
	}
}

func TestDay10NoIncompleteLines(t *testing.T) {
	if _, err := day10p02(strings.NewReader(`()`)); err == nil {
		t.Fatal("expected no incomplete lines error")
	}
}

func TestDay10ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day10p01},
		{"part2", day10p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
