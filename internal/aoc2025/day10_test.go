package aoc2025

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day10p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}`),
			Want: "7",
		},
		{
			Input: aoc.FileInput(t, 2025, 10),
			Want:  "545",
		},
	}
	aoc.AOCTest(t, day10p01, tests)
}

func Test_day10p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}`),
			Want: "33",
		},
		{
			Input: strings.NewReader(`[##] (0) (1) (0,1) {3,2}`),
			Want:  "3",
		},
		{
			Input: aoc.FileInput(t, 2025, 10),
			Want:  "22430",
		},
	}
	aoc.AOCTest(t, day10p02, tests)
}

func TestJoltageFailureStopsAggregation(t *testing.T) {
	for _, input := range []string{
		"[..] (0,1) {1,2}",
		"[.] (0) {2}\n[..] (0,1) {1,2}",
	} {
		if _, err := day10p02(strings.NewReader(input)); err == nil {
			t.Fatalf("expected infeasible-machine error for %q", input)
		}
	}
}

func TestMachineFeasibility(t *testing.T) {
	if got, err := day10p01(strings.NewReader("[#] () {0}")); err == nil || !strings.Contains(err.Error(), "lights have no reachable") {
		t.Fatalf("got %q, %v; want unreachable-lights error", got, err)
	}
	for _, input := range []string{
		"[...] (0,1) (0,2) (1,2) {1,1,1}",
		"[..] (0) (0,1) {0,1}",
	} {
		if got, err := day10p02(strings.NewReader(input)); err == nil || !strings.Contains(err.Error(), "no nonnegative integer solution") {
			t.Fatalf("input %q: got %q, %v; want infeasible-joltage error", input, got, err)
		}
	}
	for _, solve := range []func(io.Reader) (string, error){day10p01, day10p02} {
		got, err := solve(strings.NewReader("[.] {0}"))
		if err != nil || got != "0" {
			t.Fatalf("zero buttons: got %q, %v; want 0", got, err)
		}
	}
}

func TestMachineArithmeticLimits(t *testing.T) {
	const maximum = "[.] (0) {9223372036854775807}"
	got, err := day10p02(strings.NewReader(maximum))
	if err != nil || got != strconv.FormatInt(math.MaxInt64, 10) {
		t.Fatalf("got %q, %v; want MaxInt64", got, err)
	}
	if got, err := day10p02(strings.NewReader(maximum + "\n" + maximum)); err == nil || !strings.Contains(err.Error(), "machine 2: total presses") {
		t.Fatalf("got %q, %v; want aggregate-overflow error", got, err)
	}
	if got, err := day10p02(strings.NewReader("[.] (0) (0) {9223372036854775807}")); err == nil || !strings.Contains(err.Error(), "unsupported joltage range") {
		t.Fatalf("got %q, %v; want unsupported arithmetic-range error", got, err)
	}
}

func TestMachineWorkLimit(t *testing.T) {
	step, err := parseMachine([]byte("[.] (0) (0) {2}"))
	if err != nil {
		t.Fatal(err)
	}
	for _, solve := range []func(machineRequirements, *int) (int64, error){minButtonPresses, minButtonPressesForJoltage} {
		remaining := 0
		if got, err := solve(step, &remaining); err == nil || !strings.Contains(err.Error(), "work limit") {
			t.Fatalf("got %d, %v; want exhausted-work error", got, err)
		}
	}
}

func TestMachineParsing(t *testing.T) {
	wideLights := strings.Repeat(".", 63) + "#"
	wideJoltage := strings.Repeat("0,", 63) + "1"
	input := "[" + wideLights + "] (63) {" + wideJoltage + "}"
	for _, solve := range []func(io.Reader) (string, error){day10p01, day10p02} {
		got, err := solve(strings.NewReader(input))
		if err != nil || got != "1" {
			t.Fatalf("64 lights: got %q, %v; want 1", got, err)
		}
	}

	invalid := []string{
		"\n",
		"[.]",
		"[.] (-1) {0}",
		"[.] (0,0) {0}",
		"[.] (0x) {0}",
		"[.] (0) {-1}",
		"[.] (0) {0x}",
		"[.] (0) {9223372036854775808}",
		"[..] (0) {0}",
		"[" + wideLights + "] (64) {" + wideJoltage + "}",
		"[." + wideLights + "] (64) {0," + wideJoltage + "}",
		"[.] " + strings.Repeat("(0) ", 25) + "{0}",
	}
	for i, input := range invalid {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			for _, solve := range []func(io.Reader) (string, error){day10p01, day10p02} {
				if got, err := solve(strings.NewReader(input)); err == nil || !strings.Contains(err.Error(), "machine 1:") {
					t.Fatalf("input %q: got %q, %v; want contextual parse error", input, got, err)
				}
			}
		})
	}
}

func TestMachineSmallExhaustive(t *testing.T) {
	buttonText := [4]string{"()", "(0)", "(1)", "(0,1)"}
	for a := range 4 {
		for b := range 4 {
			for c := range 4 {
				buttons := [3]int{a, b, c}
				for lights := range 4 {
					pattern := []byte("..")
					for i := range 2 {
						if lights&(1<<i) != 0 {
							pattern[i] = '#'
						}
					}
					input := fmt.Sprintf("[%s] %s %s %s {0,0}", pattern, buttonText[a], buttonText[b], buttonText[c])
					want := exhaustiveMachinePresses(buttons, [2]int{lights & 1, lights >> 1}, true)
					checkSmallMachine(t, day10p01, input, want)
				}
				for x := range 3 {
					for y := range 3 {
						input := fmt.Sprintf("[..] %s %s %s {%d,%d}", buttonText[a], buttonText[b], buttonText[c], x, y)
						want := exhaustiveMachinePresses(buttons, [2]int{x, y}, false)
						checkSmallMachine(t, day10p02, input, want)
					}
				}
			}
		}
	}
}

func exhaustiveMachinePresses(buttons [3]int, target [2]int, lights bool) int64 {
	limit := 2
	if lights {
		limit = 1
	}
	best := int64(-1)
	for a := 0; a <= limit; a++ {
		for b := 0; b <= limit; b++ {
			for c := 0; c <= limit; c++ {
				presses := [3]int{a, b, c}
				var counters [2]int
				for i, button := range buttons {
					for counter := range 2 {
						if button&(1<<counter) != 0 {
							counters[counter] += presses[i]
						}
					}
				}
				if lights {
					counters[0] %= 2
					counters[1] %= 2
				}
				total := int64(a + b + c)
				if counters == target && (best < 0 || total < best) {
					best = total
				}
			}
		}
	}
	return best
}

func checkSmallMachine(t *testing.T, solve func(io.Reader) (string, error), input string, want int64) {
	t.Helper()
	got, err := solve(strings.NewReader(input))
	if want < 0 {
		if err == nil || (!strings.Contains(err.Error(), "lights have no reachable") && !strings.Contains(err.Error(), "joltage has")) {
			t.Fatalf("input %q: got %q, %v; want infeasible-machine error", input, got, err)
		}
		return
	}
	if err != nil || got != strconv.FormatInt(want, 10) {
		t.Fatalf("input %q: got %q, %v; want %d", input, got, err, want)
	}
}

func TestDay10ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day10p01, "[.] (0) {0}\n"},
		{"part2", day10p02, "[.] (0) {0}\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				input := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(input); !errors.Is(err, failure) {
					t.Fatalf("got %v; want reader failure", err)
				}
			}
		})
	}
}
