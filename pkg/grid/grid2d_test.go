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

	// Rotate 4 times should return to original
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

	// Rotate 4 times should return to original
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

	// Right then left should return to original
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

	// 180 degree rotation (two rights or two lefts)
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
