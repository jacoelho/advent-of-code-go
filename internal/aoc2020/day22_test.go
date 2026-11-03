package aoc2020

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day22p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Player 1:
9
2
6
3
1

Player 2:
5
8
4
7
10`),
			Want: "306",
		},
		{
			Input: aoc.FileInput(t, 2020, 22),
			Want:  "31455",
		},
	}
	aoc.AOCTest(t, day22p01, tests)
}

func Test_day22p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Player 1:
9
2
6
3
1

Player 2:
5
8
4
7
10`),
			Want: "291",
		},
		{
			Input: aoc.FileInput(t, 2020, 22),
			Want:  "32528",
		},
	}
	aoc.AOCTest(t, day22p02, tests)
}

func TestDay22ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"22/1", day22p01},
		{"22/2", day22p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestCombatInputErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, solve := range []func(io.Reader) (string, error){
		day22p01,
		day22p02,
	} {
		reader := io.MultiReader(strings.NewReader("Player 1:\n1\n\n"), iotest.ErrReader(failure))
		if _, err := solve(reader); !errors.Is(err, failure) {
			t.Fatalf("got %v, want input failure", err)
		}
		_, err := solve(strings.NewReader("Player 1:\nnope\n\nPlayer 2:\n2"))
		if _, ok := errors.AsType[*strconv.NumError](err); !ok {
			t.Fatalf("got %v, want number conversion error", err)
		}
	}
}
