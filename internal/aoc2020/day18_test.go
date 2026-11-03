package aoc2020

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const homeworkExample = `1 + 2 * 3 + 4 * 5 + 6
1 + (2 * 3) + (4 * (5 + 6))
2 * 3 + (4 * 5)
5 + (8 * 3 + 9 + 3 * 4 * 3)
5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))
((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2`

func Test_day18p01(t *testing.T) {
	aoc.AOCTest(t, day18p01, []aoc.TestInput{
		{Input: strings.NewReader(homeworkExample), Want: "26457"},
		{Input: strings.NewReader("1 + 2 * 3 + 4 * 5 + 6"), Want: "71"},
		{Input: strings.NewReader("2 * 3 + (4 * 5)"), Want: "26"},
		{Input: strings.NewReader(" 12 * (3 + 4) \t"), Want: "84"},
		{Input: aoc.FileInput(t, 2020, 18), Want: "1408133923393"},
	})
}

func Test_day18p02(t *testing.T) {
	aoc.AOCTest(t, day18p02, []aoc.TestInput{
		{Input: strings.NewReader(homeworkExample), Want: "694173"},
		{Input: strings.NewReader("1 + 2 * 3 + 4 * 5 + 6"), Want: "231"},
		{Input: strings.NewReader("2 * 3 + (4 * 5)"), Want: "46"},
		{Input: strings.NewReader(" 12 * (3 + 4) \t"), Want: "84"},
		{Input: aoc.FileInput(t, 2020, 18), Want: "314455761823725"},
	})
}

func TestHomeworkErrors(t *testing.T) {
	for _, solve := range []func(io.Reader) (string, error){day18p01, day18p02} {
		for _, input := range []string{"\n", "1 +", "(2 * 3", "1 + 2)", "1 2", "1 + ()", "1 / 2", "9223372036854775807 + 1", "9223372036854775807 * 2", "9223372036854775808", "9223372036854775807\n1"} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("expected error for %q", input)
			}
		}
		failure := errors.New("read homework")
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("got %v, want %v", err, failure)
		}
	}
}
