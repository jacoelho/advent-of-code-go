package aoc2020

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day10p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`16
10
15
5
1
11
7
19
6
12
4`), Want: "35"},
		{Input: strings.NewReader(`28
33
18
42
31
14
46
20
48
47
24
23
49
45
19
38
39
11
1
32
25
35
8
17
7
9
4
2
34
10
3`), Want: "220"},
		{Input: strings.NewReader(`1
3
4`), Want: "2"},
		{Input: aoc.FileInput(t, 2020, 10), Want: "2059"},
	}
	aoc.AOCTest(t, day10p01, tests)
}

func Test_day10p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`16
10
15
5
1
11
7
19
6
12
4`), Want: "8"},
		{Input: strings.NewReader(`28
33
18
42
31
14
46
20
48
47
24
23
49
45
19
38
39
11
1
32
25
35
8
17
7
9
4
2
34
10
3`), Want: "19208"},
		{Input: strings.NewReader(`1
3
4`), Want: "3"},
		{Input: aoc.FileInput(t, 2020, 10), Want: "86812553324672"},
	}
	aoc.AOCTest(t, day10p02, tests)
}

func TestAdaptersRejectBrokenChains(t *testing.T) {
	for _, input := range []string{"1\n5", "1\n1"} {
		if _, err := day10p02(strings.NewReader(input)); err == nil {
			t.Fatalf("expected chain error for %q", input)
		}
	}
}

func TestDay10ReadErrors(t *testing.T) {
	failure := errors.New("read failure")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"10/1", day10p01},
		{"10/2", day10p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestDay10MalformedInput(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"adapters", "-1", day10p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected malformed-input error")
			}
		})
	}
}
