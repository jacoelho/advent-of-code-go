package aoc2024

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day18p01_example(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0`),
			Want: "22",
		},
	}
	aoc.AOCTest(t, day18p01(6, 12), tests)
}

func Test_day18p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2024, 18),
			Want:  "280",
		},
	}
	aoc.AOCTest(t, day18p01(70, 1024), tests)
}

func Test_day18p02_example(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0`),
			Want: "6,1",
		},
	}
	aoc.AOCTest(t, day18p02(6, 12), tests)
}

func Test_day18p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2024, 18),
			Want:  "28,56",
		},
	}
	aoc.AOCTest(t, day18p02(70, 1024), tests)
}

func TestMemoryRoutes(t *testing.T) {
	for _, tt := range []struct {
		name, input, want string
		solve             func(io.Reader) (string, error)
	}{
		{"blocked start", "0,0", "", day18p01(1, 1)},
		{"blocked end", "1,1", "", day18p01(1, 1)},
		{"first byte blocks", "0,0", "0,0", day18p02(1, 0)},
		{"last byte blocks", "1,0\n0,1", "0,1", day18p02(1, 0)},
		{"no blocker", "1,0", "", day18p02(1, 0)},
		{"already blocked", "0,0", "", day18p02(1, 1)},
		{"single cell route", "0,0", "0", day18p01(0, 0)},
		{"single cell blocked", "0,0", "0,0", day18p02(0, 0)},
		{"duplicate events", "1,0\n1,0\n0,1", "0,1", day18p02(1, 2)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.solve(strings.NewReader(tt.input))
			if tt.want == "" {
				if err == nil {
					t.Fatal("expected route error")
				}
			} else if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %s", got, err, tt.want)
			}
		})
	}
}

func TestMemoryInputErrors(t *testing.T) {
	for _, factory := range []func(int, int) func(io.Reader) (string, error){day18p01, day18p02} {
		for _, config := range [][2]int{{-1, 0}, {1025, 0}, {1, -1}, {1, 2}} {
			if _, err := factory(config[0], config[1])(strings.NewReader("0,0")); err == nil {
				t.Fatalf("configuration %v: expected error", config)
			}
		}
		for _, input := range []string{"1", "1,0,1", "1,0x", "9223372036854775808,0", "-1,0", "2,0"} {
			if _, err := factory(1, 0)(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: expected error", input)
			}
		}
	}
}
