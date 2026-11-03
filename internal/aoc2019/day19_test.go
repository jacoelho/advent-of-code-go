package aoc2019

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day19p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("104,1,99"), Want: "2500"},
		{Input: strings.NewReader("104,0,99"), Want: "0"},
		{
			Input: aoc.FileInput(t, 2019, 19),
			Want:  "181",
		},
	}
	aoc.AOCTest(t, day19p01, tests)
}

func Test_day19p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("104,1,99"), Want: "0"},
		{
			Input: aoc.FileInput(t, 2019, 19),
			Want:  "4240964",
		},
	}
	aoc.AOCTest(t, day19p02, tests)
}

func TestBeamProtocolErrors(t *testing.T) {
	for _, solve := range []func(io.Reader) (string, error){day19p01, day19p02} {
		for _, program := range []string{"42", "99", "104,2,99", "104,0,104,1,99", "3,100,3,101,3,102,99"} {
			if _, err := solve(strings.NewReader(program)); err == nil {
				t.Errorf("expected beam error for %s", program)
			}
		}
	}
	if _, err := day19p02(strings.NewReader("104,0,99")); err == nil {
		t.Fatal("expected error when no square fits")
	}
}

func TestDay19ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name, prefix string
		solve        func(io.Reader) (string, error)
	}{
		{"day19p01", "99", day19p01},
		{"day19p02", "99", day19p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				reader := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(reader); !errors.Is(err, failure) {
					t.Fatalf("prefix %q: got %v, want input failure", prefix, err)
				}
			}
		})
	}
}
