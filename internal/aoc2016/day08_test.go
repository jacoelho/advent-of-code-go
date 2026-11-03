package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestScreenExample(t *testing.T) {
	pixels, err := swipeScreen(strings.NewReader("rect 3x2\nrotate column x=1 by 1\nrotate row y=0 by 4\nrotate column x=1 by 1"), 7, 3)
	if err != nil {
		t.Fatal(err)
	}
	want := " #  # #\n# #    \n #     "
	if got := strings.Join(screenRows(pixels), "\n"); got != want {
		t.Fatalf("got display:\n%s\nwant:\n%s", got, want)
	}
}

func TestScreenRotationWraps(t *testing.T) {
	pixels, err := swipeScreen(strings.NewReader("rect 1x1\nrotate row y=0 by 8\nrotate column x=1 by 4"), 7, 3)
	if err != nil {
		t.Fatal(err)
	}
	want := "       \n #     \n       "
	if got := strings.Join(screenRows(pixels), "\n"); got != want {
		t.Fatalf("got display:\n%s\nwant:\n%s", got, want)
	}
}

func Test_day08p01(t *testing.T) {
	aoc.AOCTest(t, day08p01, []aoc.TestInput{
		{Input: strings.NewReader("rect 3x2\nrotate column x=1 by 1\nrotate row y=0 by 4\nrotate column x=1 by 1"), Want: "6"},
		{Input: aoc.FileInput(t, 2016, 8), Want: "119"},
	})
}

func Test_day08p02(t *testing.T) {
	aoc.AOCTest(t, day08p02, []aoc.TestInput{
		{Input: aoc.FileInput(t, 2016, 8), Want: "ZFHFSFOGPO"},
	})
}

func TestDay08ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day08p01", day08p01},
		{"day08p02", day08p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay08MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"screen bounds", "rect 51x1", day08p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
