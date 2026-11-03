package aoc2023

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day08p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`RL

AAA = (BBB, CCC)
BBB = (DDD, EEE)
CCC = (ZZZ, GGG)
DDD = (DDD, DDD)
EEE = (EEE, EEE)
GGG = (GGG, GGG)
ZZZ = (ZZZ, ZZZ)`),
			Want: "2",
		},
		{
			Input: strings.NewReader(`LLR

AAA = (BBB, BBB)
BBB = (AAA, ZZZ)
ZZZ = (ZZZ, ZZZ)`),
			Want: "6",
		},
		{
			Input: aoc.FileInput(t, 2023, 8),
			Want:  "13771",
		},
	}
	aoc.AOCTest(t, day08p01, tests)
}

func Test_day08p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`LR

11A = (11B, XXX)
11B = (XXX, 11Z)
11Z = (11B, XXX)
22A = (22B, XXX)
22B = (22C, 22C)
22C = (22Z, 22Z)
22Z = (22B, 22B)
XXX = (XXX, XXX)`),
			Want: "6",
		},
		{
			Input: aoc.FileInput(t, 2023, 8),
			Want:  "13129439557681",
		},
	}
	aoc.AOCTest(t, day08p02, tests)
}

func TestNavigationSchedules(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "offset cycles",
			input: `L

11A = (11B, 11B)
11B = (11Z, 11Z)
11Z = (11C, 11C)
11C = (11B, 11B)
22A = (22B, 22B)
22B = (22C, 22C)
22C = (22Z, 22Z)
22Z = (22D, 22D)
22D = (22B, 22B)`,
			want: "11",
		},
		{
			name: "transient only",
			input: `L

11A = (11Z, 11Z)
11Z = (11X, 11X)
11X = (11X, 11X)
22A = (22Z, 22Z)
22Z = (22X, 22X)
22X = (22X, 22X)`,
			want: "1",
		},
		{
			name: "multiple target residues",
			input: `L

11A = (11Z, 11Z)
11Z = (11B, 11B)
11B = (12Z, 12Z)
12Z = (11A, 11A)
22A = (22B, 22B)
22B = (22C, 22C)
22C = (22Z, 22Z)
22Z = (22D, 22D)
22D = (22E, 22E)
22E = (22A, 22A)`,
			want: "3",
		},
		{
			name: "cycle entry lower bound",
			input: `L

11A = (11B, 11B)
11B = (11C, 11C)
11C = (11D, 11D)
11D = (11Z, 11Z)
11Z = (11E, 11E)
11E = (11Z, 11Z)
22A = (22Z, 22Z)
22Z = (22Z, 22Z)`,
			want: "4",
		},
		{
			name: "instruction phase",
			input: `LR

AAA = (AAA, ZZZ)
ZZZ = (ZZZ, ZZZ)`,
			want: "2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := day08p02(strings.NewReader(tt.input))
			if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}

func TestNavigationLargeAlignment(t *testing.T) {
	var input strings.Builder
	input.WriteString("L\n\n")
	for cycle, length := range []int{97, 101, 103, 107, 109, 113, 127, 131, 137, 139} {
		names := make([]string, length)
		for i := range names {
			suffix := "N"
			if i == 0 {
				suffix = "A"
			} else if i == length-1 {
				suffix = "Z"
			}
			names[i] = fmt.Sprintf("C%dN%d%s", cycle, i, suffix)
		}
		for i, name := range names {
			next := names[(i+1)%length]
			fmt.Fprintf(&input, "%s = (%s, %s)\n", name, next, next)
		}
	}
	// Each prime-length cycle hits at -1 modulo its length.
	const want = "421336842070675358938"
	got, err := day08p02(strings.NewReader(input.String()))
	if err != nil || got != want {
		t.Fatalf("got %q, %v; want %q", got, err, want)
	}
}

func TestNavigationInvalidInput(t *testing.T) {
	for _, input := range []string{
		"",
		"X\n\nAAA = (AAA, AAA)",
		"L\n\nAAA = (ZZZ, ZZZ)",
		"L\n\nAAA = (AAA, AAA)\nAAA = (AAA, AAA)",
	} {
		for _, solve := range []func(io.Reader) (string, error){day08p01, day08p02} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("input %q: got %q, want error", input, got)
			}
		}
	}
}

func TestNavigationUnreachable(t *testing.T) {
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
		input string
	}{
		{
			name:  "unreachable ZZZ",
			solve: day08p01,
			input: "L\n\nAAA = (AAA, AAA)\nZZZ = (ZZZ, ZZZ)",
		},
		{
			name:  "no target",
			solve: day08p02,
			input: "L\n\nAAA = (AAA, AAA)",
		},
		{
			name:  "no starts",
			solve: day08p02,
			input: "L\n\nZZZ = (ZZZ, ZZZ)",
		},
		{
			name:  "incompatible noncoprime periods",
			solve: day08p02,
			input: `L

11A = (11Z, 11Z)
11Z = (11A, 11A)
22A = (22B, 22B)
22B = (22Z, 22Z)
22Z = (22C, 22C)
22C = (22A, 22A)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatalf("got %q, want error", got)
			}
		})
	}
}

func TestNavigationStateLimit(t *testing.T) {
	net := network{"AAA": {"ZZZ", "ZZZ"}, "ZZZ": {"ZZZ", "ZZZ"}}
	if _, err := analyzeNavigation("L", net, "AAA", new(1)); err == nil {
		t.Fatal("want state limit error")
	}
}

func TestNavigationReaderErrors(t *testing.T) {
	want := errors.New("navigation read failure")
	for _, solve := range []func(io.Reader) (string, error){day08p01, day08p02} {
		for _, prefix := range []string{"", "L\n\nAAA = (ZZZ, ZZZ)\nZZZ = (ZZZ, ZZZ)\n"} {
			r := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(want))
			if _, err := solve(r); !errors.Is(err, want) {
				t.Errorf("got %v, want reader error", err)
			}
		}
	}
}
