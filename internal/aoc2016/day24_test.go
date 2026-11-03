package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const cleaningExample = "###########\n#0.1.....2#\n#.#######.#\n#4.......3#\n###########\n"

func Test_day24p01(t *testing.T) {
	aoc.AOCTest(t, day24p01, []aoc.TestInput{{Input: strings.NewReader(cleaningExample), Want: "14"}, {Input: strings.NewReader("0.7"), Want: "2"}, {Input: strings.NewReader("0"), Want: "0"}, {Input: aoc.FileInput(t, 2016, 24), Want: "430"}})
}

func Test_day24p02(t *testing.T) {
	aoc.AOCTest(t, day24p02, []aoc.TestInput{{Input: strings.NewReader(cleaningExample), Want: "20"}, {Input: strings.NewReader("0.7"), Want: "4"}, {Input: strings.NewReader("0"), Want: "0"}, {Input: aoc.FileInput(t, 2016, 24), Want: "700"}})
}

func TestCleaningUnreachableLocation(t *testing.T) {
	if _, err := day24p01(strings.NewReader("0#7")); err == nil {
		t.Fatal("accepted unreachable location")
	}
}

func TestDay24ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day24p01,
		day24p02,
	} {
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("read error%v lost input cause", err)
		}
	}
}
