package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const packageExample = "1\n2\n3\n4\n5\n7\n8\n9\n10\n11\n"

func Test_day24p01(t *testing.T) {
	aoc.AOCTest(t, day24p01, []aoc.TestInput{
		{Input: strings.NewReader(packageExample), Want: "99"},
		{Input: strings.NewReader("20\n17\n16\n15\n14\n11\n10\n6\n4\n1\n"), Want: "1120"},
		{Input: aoc.FileInput(t, 2015, 24), Want: "11266889531"},
	})
}

func Test_day24p02(t *testing.T) {
	aoc.AOCTest(t, day24p02, []aoc.TestInput{
		{Input: strings.NewReader(packageExample), Want: "44"},
		{Input: aoc.FileInput(t, 2015, 24), Want: "77387711"},
	})
}

func TestPackageWeightsRemainDistinctEntries(t *testing.T) {
	aoc.AOCTest(t, day24p01, []aoc.TestInput{{Input: strings.NewReader("1\n1\n1\n1\n1\n1\n"), Want: "1"}})
	for _, input := range []string{"1\n2\n4", "1\n1\n4", "0\n3\n3\n3", "3\n3\n3\n3\n3"} {
		if _, err := day24p01(strings.NewReader(input)); err == nil {
			t.Errorf("accepted impossible or invalid packages %q", input)
		}
	}
}

func TestDay24ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day24p01,
		day24p02,
	} {
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("read error = %v; want wrapped input failure", err)
		}
	}
}
