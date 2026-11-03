package aoc2024

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0`),
			Want: "4,6,3,5,6,3,5,2,1,0",
		},
		{
			Input: aoc.FileInput(t, 2024, 17),
			Want:  "2,7,4,7,2,1,7,5,1",
		},
	}
	aoc.AOCTest(t, day17p01, tests)
}

func Test_day17p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Register A: 2024
		Register B: 0
		Register C: 0
		
		Program: 0,3,5,4,3,0`),
			Want: "117440",
		},
		{
			Input: aoc.FileInput(t, 2024, 17),
			Want:  "37221274271220",
		},
	}
	aoc.AOCTest(t, day17p02, tests)
}

func computerInput(a, b, c int64, program string) string {
	return fmt.Sprintf("Register A: %d\nRegister B: %d\nRegister C: %d\n\nProgram: %s", a, b, c, program)
}

func TestComputerInstructions(t *testing.T) {
	for _, tt := range []struct {
		name, program, want string
		a, b, c             int64
	}{
		{"literal seven", "1,7,5,5", "2", 0, 29, 0},
		{"ignored seven", "4,7,5,5", "2", 0, 2024, 43690},
		{"odd jump", "3,3,0,5,1", "1", 1, 0, 0},
		{"jump beyond program", "3,7", "", 1, 0, 0},
		{"large shift", "0,5,5,4", "0", math.MaxInt64, 100, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := day17p01(strings.NewReader(computerInput(tt.a, tt.b, tt.c, tt.program)))
			if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %q", got, err, tt.want)
			}
		})
	}
	for _, program := range []string{"5,7", "3,3,0,5,1"} {
		if _, err := day17p01(strings.NewReader(computerInput(0, 0, 0, program))); err == nil {
			t.Fatalf("program %q: expected instruction error", program)
		}
	}
	computer := chronospatialComputer{a: 1, program: []uint8{3, 0}}
	work := computerWork{steps: 8}
	if _, err := computer.run(&work); err == nil || !strings.Contains(err.Error(), "budget") {
		t.Fatalf("loop error = %v; want instruction budget error", err)
	}
}

func TestComputerQuines(t *testing.T) {
	got, err := day17p02(strings.NewReader(computerInput(0, 0, 0, "2,0,0,3,5,4,3,0")))
	if err != nil || got != "7516176" {
		t.Fatalf("quine = %q, %v; want 7516176", got, err)
	}
	for _, tt := range []struct{ program, category string }{
		{"5,0", "unsupported"},
		{"1,1,0,3,5,5,3,0", "unsupported"},
		{"0,3,5,0,3,0", "no quine"},
	} {
		_, err := day17p02(strings.NewReader(computerInput(0, 0, 0, tt.program)))
		if err == nil || !strings.Contains(err.Error(), tt.category) {
			t.Fatalf("program %q: got %v; want %s error", tt.program, err, tt.category)
		}
	}
}

func TestComputerInputErrors(t *testing.T) {
	valid := computerInput(0, 0, 0, "0,3,5,4,3,0")
	for _, input := range []string{
		strings.Replace(valid, "Register B: 0\n", "", 1),
		"Register A: 0\n" + valid,
		computerInput(-1, 0, 0, "0,3,5,4,3,0"),
		strings.Replace(valid, "Register A: 0", "Register A: 9223372036854775808", 1),
		computerInput(0, 0, 0, "8"),
		computerInput(0, 0, 0, "0,,3"),
	} {
		for _, solve := range []func(io.Reader) (string, error){day17p01, day17p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: expected parse error", input)
			}
		}
	}
}

func TestDay17ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day17p01, "Register A: 0\nRegister B: 0\nRegister C: 0\n\nProgram: 0,3,5,4,3,0\n"},
		{"part2", day17p02, "Register A: 0\nRegister B: 0\nRegister C: 0\n\nProgram: 0,3,5,4,3,0\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				r := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(want))
				if _, err := tt.solve(r); !errors.Is(err, want) {
					t.Errorf("prefix %q: got %v, want reader error", prefix, err)
				}
			}
		})
	}
}
