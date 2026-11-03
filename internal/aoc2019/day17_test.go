package aoc2019

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestScaffoldAlignment(t *testing.T) {
	scaffold, err := parseScaffoldGrid(strings.NewReader(`..#..........
..#..........
#######...###
#.#...#...#.#
#############
..#...#...#..
..#####...^..`))
	if err != nil {
		t.Fatal(err)
	}
	if got := calculateAlignmentSum(scaffold); got != 76 {
		t.Fatalf("alignment sum = %d, want 76", got)
	}
}

func Test_day17p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 17),
			Want:  "6448",
		},
	}
	aoc.AOCTest(t, day17p01, tests)
}

func Test_day17p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 17),
			Want:  "914900",
		},
	}
	aoc.AOCTest(t, day17p02, tests)
}

func TestGenerateScaffoldPath(t *testing.T) {
	for _, tt := range []struct {
		name    string
		grid    string
		want    []string
		wantErr bool
	}{
		{
			name: "turns in place",
			grid: "^##\n..#",
			want: []string{"R", "2", "R", "1"},
		},
		{
			name:    "missing robot",
			grid:    "#",
			wantErr: true,
		},
		{
			name:    "multiple robots",
			grid:    "^#>",
			wantErr: true,
		},
		{
			name:    "closed loop",
			grid:    ">#\n##",
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			scaffold, err := parseScaffoldGrid(strings.NewReader(tt.grid))
			if err != nil {
				t.Fatal(err)
			}
			got, err := generatePath(scaffold)
			if (err != nil) != tt.wantErr {
				t.Fatalf("generatePath() error = %v, want error %v", err, tt.wantErr)
			}
			if err == nil && !slices.Equal(got, tt.want) {
				t.Fatalf("path = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("entrypoint preserves robot error", func(t *testing.T) {
		_, err := day17p02(strings.NewReader("104,35,104,10,99"))
		if err == nil || !strings.Contains(err.Error(), "robot") {
			t.Fatalf("expected robot error, got %v", err)
		}
	})
}

func assertMovementExpansion(t *testing.T, movement movementProgram, path []string) {
	t.Helper()

	if len(movement.main) > 20 || !regexp.MustCompile(`^[ABC](,[ABC])*$`).MatchString(movement.main) {
		t.Fatalf("invalid main routine %q", movement.main)
	}
	for _, function := range movement.functions {
		if len(function) > 20 || !regexp.MustCompile(`^[LR],[0-9]+(,[LR],[0-9]+)*$`).MatchString(function) {
			t.Fatalf("invalid movement function %q", function)
		}
		tokens := strings.Split(function, ",")
		for i := 1; i < len(tokens); i += 2 {
			steps, err := strconv.Atoi(tokens[i])
			if err != nil || steps <= 0 {
				t.Fatalf("invalid step count %q", tokens[i])
			}
		}
	}

	var expanded []string
	for call := range strings.SplitSeq(movement.main, ",") {
		expanded = append(expanded, strings.Split(movement.functions[call[0]-'A'], ",")...)
	}
	if !slices.Equal(expanded, path) {
		t.Fatalf("expanded path = %v, want %v", expanded, path)
	}
}

func TestCompressPath(t *testing.T) {
	const (
		a = "R,12,L,12,R,12,L,123"
		b = "L,34,R,34,L,34,R,345"
		c = "R,56,R,56,L,56,L,567"
	)
	boundaryRoute := strings.Join([]string{a, b, c, a, b, c, a, b, c, a}, ",")

	t.Run("multi-digit counts", func(t *testing.T) {
		path := []string{"R", "1", "R", "10"}
		movement, err := compressPath(path)
		if err != nil {
			t.Fatal(err)
		}
		assertMovementExpansion(t, movement, path)
	})

	t.Run("line limits", func(t *testing.T) {
		path := strings.Split(boundaryRoute, ",")
		movement, err := compressPath(path)
		if err != nil {
			t.Fatal(err)
		}
		assertMovementExpansion(t, movement, path)

		// This 209-byte route requires ten calls and 20-byte functions in any encoding.
		if len(movement.main) != 19 {
			t.Fatalf("main routine length = %d, want 19", len(movement.main))
		}
		for _, function := range movement.functions {
			if len(function) != 20 {
				t.Fatalf("function length = %d, want 20", len(function))
			}
		}
	})

	for _, tt := range []struct {
		name  string
		route string
	}{
		{name: "empty"},
		{name: "incomplete pair", route: "R"},
		{name: "invalid turn", route: "X,1"},
		{name: "zero steps", route: "L,0"},
		{name: "signed count", route: "L,+1"},
		{name: "function too long", route: "R,1234567890123456789"},
		{name: "eleven calls", route: boundaryRoute + "," + b},
		{name: "no three-function encoding", route: "R,10,L,11,R,12,L,13,R,14,L,15,R,16,L,17,R,18,L,19,R,20,L,21,R,22"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var path []string
			if tt.route != "" {
				path = strings.Split(tt.route, ",")
			}
			if _, err := compressPath(path); err == nil {
				t.Fatal("expected compression error")
			}
		})
	}
}

func TestRunVacuumRobot(t *testing.T) {
	for _, tt := range []struct {
		name    string
		program []int
		want    int
		wantErr bool
		cause   string
	}{
		{
			name:    "dust after transcript",
			program: []int{2, 100, 100, 100, 104, 65, 104, 10, 104, 300, 99},
			want:    300,
		},
		{
			name:    "empty response",
			program: []int{2, 100, 100, 100, 99},
			wantErr: true,
		},
		{
			name:    "ASCII-only response",
			program: []int{2, 100, 100, 100, 104, 65, 104, 10, 99},
			wantErr: true,
		},
		{
			name:    "dust before input request",
			program: []int{2, 100, 100, 100, 104, 300, 3, 101, 99},
			wantErr: true,
		},
		{
			name:    "multiple dust values",
			program: []int{2, 100, 100, 100, 104, 300, 104, 301, 99},
			wantErr: true,
		},
		{
			name:    "VM error",
			program: []int{2, 100, 100, 100, 42},
			wantErr: true,
			cause:   "opcode",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runVacuumRobot(tt.program, nil)
			if (err != nil) != tt.wantErr {
				t.Fatalf("runVacuumRobot() error = %v, want error %v", err, tt.wantErr)
			}
			if tt.cause != "" && !strings.Contains(err.Error(), tt.cause) {
				t.Fatalf("error %q does not retain cause %q", err, tt.cause)
			}
			if err == nil && got != tt.want {
				t.Fatalf("dust = %d, want %d", got, tt.want)
			}
		})
	}
}
