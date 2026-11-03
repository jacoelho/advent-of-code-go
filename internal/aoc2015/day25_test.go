package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	aoc.AOCTest(t, day25p01, []aoc.TestInput{
		{Input: strings.NewReader("row 1, column 1"), Want: "20151125"},
		{Input: strings.NewReader("row 2, column 1"), Want: "31916031"},
		{Input: strings.NewReader("row 1, column 2"), Want: "18749137"},
		{Input: strings.NewReader("row 6, column 6"), Want: "27995004"},
		{Input: aoc.FileInput(t, 2015, 25), Want: "9132360"},
	})
	for _, input := range []string{"", "row 0, column 1", "row 1, column 0"} {
		if _, err := day25p01(strings.NewReader(input)); err == nil {
			t.Errorf("accepted %q", input)
		}
	}
}

func TestDay25ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day25p01,
	} {
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("read error = %v; want wrapped input failure", err)
		}
	}
}
