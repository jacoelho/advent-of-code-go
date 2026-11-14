package grid

import (
	"testing"

	"golang.org/x/exp/constraints"
)

func TestGrid2D_TurnRight(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]rune
		expected [][]rune
	}{
		{
			name: "2x2 diagonal pattern",
			input: [][]rune{
				{'#', '.'},
				{'.', '#'},
			},
			expected: [][]rune{
				{'.', '#'},
				{'#', '.'},
			},
		},
		{
			name: "2x2 top row",
			input: [][]rune{
				{'#', '#'},
				{'.', '.'},
			},
			expected: [][]rune{
				{'.', '#'},
				{'.', '#'},
			},
		},
		{
			name: "3x3 L shape",
			input: [][]rune{
				{'#', '.', '.'},
				{'#', '.', '.'},
				{'#', '#', '#'},
			},
			expected: [][]rune{
				{'#', '#', '#'},
				{'#', '.', '.'},
				{'#', '.', '.'},
			},
		},
		{
			name: "1x3 horizontal",
			input: [][]rune{
				{'#', '#', '#'},
			},
			expected: [][]rune{
				{'#'},
				{'#'},
				{'#'},
			},
		},
		{
			name: "3x1 vertical",
			input: [][]rune{
				{'#'},
				{'#'},
				{'#'},
			},
			expected: [][]rune{
				{'#', '#', '#'},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGrid2D[int](tt.input)
			rotated := grid.TurnRight()
			expected := NewGrid2D[int](tt.expected)

			if !gridsEqual(rotated, expected) {
				t.Errorf("TurnRight() failed\nGot:\n%sExpected:\n%s",
					gridToString(rotated), gridToString(expected))
			}
		})
	}
}

func TestGrid2D_TurnLeft(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]rune
		expected [][]rune
	}{
		{
			name: "2x2 diagonal pattern",
			input: [][]rune{
				{'#', '.'},
				{'.', '#'},
			},
			expected: [][]rune{
				{'.', '#'},
				{'#', '.'},
			},
		},
		{
			name: "2x2 top row",
			input: [][]rune{
				{'#', '#'},
				{'.', '.'},
			},
			expected: [][]rune{
				{'#', '.'},
				{'#', '.'},
			},
		},
		{
			name: "3x3 L shape",
			input: [][]rune{
				{'#', '.', '.'},
				{'#', '.', '.'},
				{'#', '#', '#'},
			},
			expected: [][]rune{
				{'.', '.', '#'},
				{'.', '.', '#'},
				{'#', '#', '#'},
			},
		},
		{
			name: "1x3 horizontal",
			input: [][]rune{
				{'#', '#', '#'},
			},
			expected: [][]rune{
				{'#'},
				{'#'},
				{'#'},
			},
		},
		{
			name: "3x1 vertical",
			input: [][]rune{
				{'#'},
				{'#'},
				{'#'},
			},
			expected: [][]rune{
				{'#', '#', '#'},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := NewGrid2D[int](tt.input)
			rotated := grid.TurnLeft()
			expected := NewGrid2D[int](tt.expected)

			if !gridsEqual(rotated, expected) {
				t.Errorf("TurnLeft() failed\nGot:\n%sExpected:\n%s",
					gridToString(rotated), gridToString(expected))
			}
		})
	}
}

func TestGrid2D_TurnRightFourTimes(t *testing.T) {
	input := [][]rune{
		{'#', '.', '.'},
		{'#', '#', '.'},
		{'.', '.', '#'},
	}

	grid := NewGrid2D[int](input)
	original := NewGrid2D[int](input)

	rotated := (&grid).TurnRight()
	rotated = (&rotated).TurnRight()
	rotated = (&rotated).TurnRight()
	rotated = (&rotated).TurnRight()

	if !gridsEqual(rotated, original) {
		t.Errorf("Four right rotations should return to original\nGot:\n%sExpected:\n%s",
			gridToString(rotated), gridToString(original))
	}
}

func TestGrid2D_TurnLeftFourTimes(t *testing.T) {
	input := [][]rune{
		{'#', '.', '.'},
		{'#', '#', '.'},
		{'.', '.', '#'},
	}

	grid := NewGrid2D[int](input)
	original := NewGrid2D[int](input)

	rotated := (&grid).TurnLeft()
	rotated = (&rotated).TurnLeft()
	rotated = (&rotated).TurnLeft()
	rotated = (&rotated).TurnLeft()

	if !gridsEqual(rotated, original) {
		t.Errorf("Four left rotations should return to original\nGot:\n%sExpected:\n%s",
			gridToString(rotated), gridToString(original))
	}
}

func TestGrid2D_TurnRightThenLeft(t *testing.T) {
	input := [][]rune{
		{'#', '.', '.'},
		{'#', '#', '.'},
		{'.', '.', '#'},
	}

	grid := NewGrid2D[int](input)
	original := NewGrid2D[int](input)

	rotated := (&grid).TurnRight()
	rotated = (&rotated).TurnLeft()

	if !gridsEqual(rotated, original) {
		t.Errorf("Right then left rotation should return to original\nGot:\n%sExpected:\n%s",
			gridToString(rotated), gridToString(original))
	}
}

func TestGrid2D_Turn180(t *testing.T) {
	input := [][]rune{
		{'#', '.'},
		{'.', '#'},
	}

	expected := [][]rune{
		{'#', '.'},
		{'.', '#'},
	}

	grid := NewGrid2D[int](input)
	expectedGrid := NewGrid2D[int](expected)

	rotated := (&grid).TurnRight()
	rotated = (&rotated).TurnRight()

	if !gridsEqual(rotated, expectedGrid) {
		t.Errorf("180 degree rotation failed\nGot:\n%sExpected:\n%s",
			gridToString(rotated), gridToString(expectedGrid))
	}
}

func gridsEqual[T constraints.Signed, V comparable](a, b Grid2D[T, V]) bool {
	if len(a) != len(b) {
		return false
	}

	for pos, val := range a {
		if bVal, exists := b[pos]; !exists || bVal != val {
			return false
		}
	}

	return true
}

func gridToString[T constraints.Signed](g Grid2D[T, rune]) string {
	minX, maxX, minY, maxY := g.Dimensions()
	var result string

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if v, exists := g[Position2D[T]{X: x, Y: y}]; exists {
				result += string(v)
			} else {
				result += "?"
			}
		}
		result += "\n"
	}

	return result
}

func TestGrid2D_FlipsAndEdgesOfOffsetGrid(t *testing.T) {
	g := Grid2D[int, rune]{}
	for y, row := range []string{"ab", "cd", "ef"} {
		for x, ch := range row {
			g[Position2D[int]{X: x - 3, Y: y + 5}] = ch
		}
	}
	for _, tt := range []struct {
		name string
		got  Grid2D[int, rune]
		want [][]rune
	}{
		{"horizontal", g.FlipHorizontal(), [][]rune{[]rune("ba"), []rune("dc"), []rune("fe")}},
		{"vertical", g.FlipVertical(), [][]rune{[]rune("ef"), []rune("cd"), []rune("ab")}},
		{"right", g.TurnRight(), [][]rune{[]rune("eca"), []rune("fdb")}},
		{"left", g.TurnLeft(), [][]rune{[]rune("bdf"), []rune("ace")}},
	} {
		if want := NewGrid2D[int](tt.want); !gridsEqual(tt.got, want) {
			t.Errorf("%s:\ngot\n%swant\n%s", tt.name, gridToString(tt.got), gridToString(want))
		}
	}
	for _, tt := range []struct {
		name string
		got  []rune
		want string
	}{
		{"top", g.Top(), "ab"},
		{"bottom", g.Bottom(), "ef"},
		{"left", g.Left(), "ace"},
		{"right", g.Right(), "bdf"},
	} {
		if string(tt.got) != tt.want {
			t.Errorf("%s = %q; want %q", tt.name, string(tt.got), tt.want)
		}
	}
}

func TestNeighboursIncludeNegativeCoordinates(t *testing.T) {
	var got []Position2D[int]
	for n := range Neighbours4(Position2D[int]{}) {
		got = append(got, n)
	}
	want := []Position2D[int]{{X: 1}, {X: -1}, {Y: -1}, {Y: 1}}
	if len(got) != len(want) {
		t.Fatalf("Neighbours4 = %v; want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Neighbours4 = %v; want %v", got, want)
		}
	}
}

func TestPosition2DWithin(t *testing.T) {
	for _, tt := range []struct {
		p    Position2D[int]
		want bool
	}{
		{Position2D[int]{X: 0, Y: 0}, true},
		{Position2D[int]{X: 2, Y: 1}, true},
		{Position2D[int]{X: -1, Y: 0}, false},
		{Position2D[int]{X: 0, Y: -1}, false},
		{Position2D[int]{X: 3, Y: 0}, false},
		{Position2D[int]{X: 0, Y: 2}, false},
	} {
		if got := tt.p.Within(3, 2); got != tt.want {
			t.Errorf("%v.Within(3, 2) = %v; want %v", tt.p, got, tt.want)
		}
	}
}
