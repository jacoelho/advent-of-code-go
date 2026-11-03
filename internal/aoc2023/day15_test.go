package aoc2023

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day15p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7`),
			Want:  "1320",
		},
		{
			Input: aoc.FileInput(t, 2023, 15),
			Want:  "516469",
		},
	}
	aoc.AOCTest(t, day15p01, tests)
}

func Test_day15p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7`),
			Want:  "145",
		},
		{
			Input: aoc.FileInput(t, 2023, 15),
			Want:  "221627",
		},
	}
	aoc.AOCTest(t, day15p02, tests)
}

func TestInitializationPayload(t *testing.T) {
	long := strings.TrimSuffix(strings.Repeat("rn=1,", 16384), ",")
	aoc.AOCTest(t, day15p01, []aoc.TestInput{
		{Input: strings.NewReader("rn=1\r\n"), Want: "30"},
		{Input: strings.NewReader(long), Want: "491520"},
	})
	aoc.AOCTest(t, day15p02, []aoc.TestInput{
		{Input: strings.NewReader("rn=1\r\n"), Want: "1"},
		{Input: strings.NewReader(long), Want: "1"},
		{Input: strings.NewReader("rn=9223372036854775807"), Want: "9223372036854775807"},
	})
	for _, solve := range []func(io.Reader) (string, error){day15p01, day15p02} {
		if _, err := solve(strings.NewReader("rn=1\nqp=2")); err == nil {
			t.Fatal("expected multiple payload lines error")
		}
	}
}

func TestLensFocalErrors(t *testing.T) {
	if _, err := day15p02(strings.NewReader("rn=nope")); !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("got %v, want numeric syntax cause", err)
	}
	for _, input := range []string{"rn=-1", "qp=9223372036854775807", "rn=9223372036854775807,qp=1"} {
		if _, err := day15p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected focal length or power error for %q", input)
		}
	}
}

func TestDay15ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day15p01, "rn=1\n"},
		{"part2", day15p02, "rn=1\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, mode := range []struct{ name, prefix string }{{"immediate", ""}, {"after prefix", tt.prefix}} {
				t.Run(mode.name, func(t *testing.T) {
					input := io.MultiReader(strings.NewReader(mode.prefix), iotest.ErrReader(failure))
					_, err := tt.solve(input)
					if !errors.Is(err, failure) {
						t.Fatalf("got %v, want reader failure", err)
					}
				})
			}
		})
	}
}
