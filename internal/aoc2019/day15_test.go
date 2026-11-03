package aoc2019

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day15p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(twoRoomDroidProgram), Want: "1"},
		{
			Input: aoc.FileInput(t, 2019, 15),
			Want:  "216",
		},
	}
	aoc.AOCTest(t, day15p01, tests)
}

func Test_day15p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(twoRoomDroidProgram), Want: "1"},
		{
			Input: aoc.FileInput(t, 2019, 15),
			Want:  "326",
		},
	}
	aoc.AOCTest(t, day15p02, tests)
}

const twoRoomDroidProgram = "1101,99,0,0,109,100,203,1,21208,0,0,2,1205,2,31,21208,1,2,2,1206,2,47,21101,0,0,0,104,1,1105,1,6,21208,1,1,2,1206,2,47,21101,1,0,0,104,2,1105,1,6,104,0,1105,1,6"

func TestDroidProtocolErrors(t *testing.T) {
	programs := []string{
		"42",
		"3,100,99",
		"3,100,104,0,104,1,99",
		"3,100,104,3,1105,1,0",
		"3,100,104,0,1105,1,0",
	}
	for _, solve := range []func(io.Reader) (string, error){day15p01, day15p02} {
		for _, program := range programs {
			if _, err := solve(strings.NewReader(program)); err == nil {
				t.Errorf("expected droid error for %s", program)
			}
		}
	}
}

func TestDay15ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name, prefix string
		solve        func(io.Reader) (string, error)
	}{
		{"day15p01", "99", day15p01},
		{"day15p02", "99", day15p02},
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
