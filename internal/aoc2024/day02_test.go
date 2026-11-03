package aoc2024

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestShortReports(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  [2]string
	}{
		{"1", [2]string{"1", "1"}},
		{"1 2", [2]string{"1", "1"}},
		{"1 9", [2]string{"0", "1"}},
	} {
		for part, solve := range []func(io.Reader) (string, error){day02p01, day02p02} {
			got, err := solve(strings.NewReader(tt.input))
			if err != nil || got != tt.want[part] {
				t.Fatalf("part %d input %q = %q, %v; want %s", part+1, tt.input, got, err, tt.want[part])
			}
		}
	}
}

func Test_day02p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9`),
			Want: "2",
		},
		{
			Input: aoc.FileInput(t, 2024, 2),
			Want:  "472",
		},
	}

	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9`),
			Want: "4",
		},
		{
			Input: aoc.FileInput(t, 2024, 2),
			Want:  "520",
		},
	}

	aoc.AOCTest(t, day02p02, tests)
}

func TestReportParseErrors(t *testing.T) {
	for _, solve := range []func(io.Reader) (string, error){day02p01, day02p02} {
		_, err := solve(strings.NewReader("1 2 3\n4 nope 6"))
		if _, ok := errors.AsType[*strconv.NumError](err); !ok {
			t.Fatalf("got %v; want numeric conversion error", err)
		}
		if got, err := solve(strings.NewReader("\n")); err == nil {
			t.Fatalf("got %q; want empty-report error", got)
		}
	}
}
