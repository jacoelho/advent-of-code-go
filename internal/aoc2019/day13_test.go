package aoc2019

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day13p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("104,0,104,0,104,2,104,0,104,0,104,0,99"), Want: "0"},
		{Input: strings.NewReader("104,0,104,0,104,2,104,0,104,0,104,2,104,1,104,0,104,2,104,0,104,0,104,0,99"), Want: "1"},
		{
			Input: aoc.FileInput(t, 2019, 13),
			Want:  "268",
		},
	}
	aoc.AOCTest(t, day13p01, tests)
}

func Test_day13p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("2,20,21,22,104,-1,104,0,104,17,99"), Want: "17"},
		{Input: strings.NewReader(joystickProgram), Want: "1"},
		{
			Input: aoc.FileInput(t, 2019, 13),
			Want:  "13989",
		},
	}
	aoc.AOCTest(t, day13p02, tests)
}

const joystickProgram = "2,100,100,100,104,0,104,0,104,3,104,2,104,0,104,4,3,101,104,-1,104,0,4,101,99"

func TestArcadeProtocolErrors(t *testing.T) {
	programs := []string{
		"2,100,100,100,104,0,104,0,99",
		"2,100,100,100,104,0,104,0,104,5,99",
		"2,100,100,100,104,-2,104,0,104,2,99",
		"2,100,100,100,104,0,104,-1,104,2,99",
	}
	for _, solve := range []func(io.Reader) (string, error){day13p01, day13p02} {
		for _, program := range programs {
			if _, err := solve(strings.NewReader(program)); err == nil {
				t.Errorf("expected arcade error for %s", program)
			}
		}
	}
	if _, err := day13p01(strings.NewReader(joystickProgram)); err == nil {
		t.Fatal("expected joystick request error during scan")
	}
}

func TestDay13ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name, prefix string
		solve        func(io.Reader) (string, error)
	}{
		{"day13p01", "99", day13p01},
		{"day13p02", "99", day13p02},
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
