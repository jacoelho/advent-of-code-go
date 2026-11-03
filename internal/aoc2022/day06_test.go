package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day06Example = `mjqjpqmgbljsphdztnvjfqwrcgsmlb`

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day06Example), Want: "7"},
		{Input: aoc.FileInput(t, 2022, 6), Want: "1578"},
	}
	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day06Example), Want: "19"},
		{Input: aoc.FileInput(t, 2022, 6), Want: "2178"},
	}
	aoc.AOCTest(t, day06p02, tests)
}

func TestDatastreamMarkers(t *testing.T) {
	aoc.AOCTest(t, day06p01, []aoc.TestInput{{Input: strings.NewReader("aAbB\r\n"), Want: "4"}})
	inputs := []string{"bvwbjplbgvbhsrlpgdmjqwftvncz\n", "nppdvjthqldpwncqszvftbrmjlhg", "nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg", "zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw"}
	first := []string{"5", "6", "10", "11"}
	second := []string{"23", "23", "29", "26"}
	for i, input := range inputs {
		aoc.AOCTest(t, day06p01, []aoc.TestInput{{Input: strings.NewReader(input), Want: first[i]}})
		aoc.AOCTest(t, day06p02, []aoc.TestInput{{Input: strings.NewReader(input), Want: second[i]}})
	}
	for _, input := range []string{"aaaa\n", "abc\n", "abcd\nefgh"} {
		if _, err := day06p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected datastream error for %q", input)
		}
	}
}

func TestDay06ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day06p01},
		{"part2", day06p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
