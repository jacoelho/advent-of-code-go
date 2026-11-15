package aoc2023

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type digStep struct {
	direction byte
	distance  int
}

// digInstruction holds the plan step and the step encoded in its colour.
type digInstruction struct {
	plan, colour digStep
}

func parseDigInstruction(line []byte) (digInstruction, error) {
	parts := strings.Fields(string(line))
	if len(parts) != 3 || len(parts[0]) != 1 || !strings.Contains("UDLR", parts[0]) {
		return digInstruction{}, fmt.Errorf("invalid dig instruction %q", line)
	}
	distance, err := strconv.Atoi(parts[1])
	if err != nil || distance < 0 {
		return digInstruction{}, fmt.Errorf("invalid dig distance in %q", line)
	}
	colour, ok := strings.CutPrefix(parts[2], "(#")
	colour, closed := strings.CutSuffix(colour, ")")
	if !ok || !closed || len(colour) != 6 || strings.Trim(colour, "0123456789abcdef") != "" || colour[5] > '3' {
		return digInstruction{}, fmt.Errorf("invalid dig colour in %q", line)
	}
	colourDistance, _ := strconv.ParseInt(colour[:5], 16, 64)
	return digInstruction{
		plan:   digStep{direction: parts[0][0], distance: distance},
		colour: digStep{direction: "RDLU"[colour[5]-'0'], distance: int(colourDistance)},
	}, nil
}

func lagoonArea(r io.Reader, step func(digInstruction) digStep) (string, error) {
	s := scanner.NewScanner(r, parseDigInstruction)
	instructions := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(instructions) == 0 {
		return "", fmt.Errorf("empty dig plan")
	}

	vertices := make([]grid.Position2D[int], 0, len(instructions))
	var current grid.Position2D[int]
	perimeter := 0
	for _, instruction := range instructions {
		dig := step(instruction)
		switch dig.direction {
		case 'R':
			current.X += dig.distance
		case 'L':
			current.X -= dig.distance
		case 'U':
			current.Y -= dig.distance
		case 'D':
			current.Y += dig.distance
		}
		perimeter += dig.distance
		vertices = append(vertices, current)
	}
	if current != (grid.Position2D[int]{}) {
		return "", fmt.Errorf("dig plan ends at %d,%d instead of its start", current.X, current.Y)
	}

	// Pick's theorem gives i = area - perimeter/2 + 1 interior lattice points.
	// Including the boundary tiles gives i + perimeter = area + perimeter/2 + 1.
	return strconv.Itoa(xmath.PolygonArea(vertices) + perimeter/2 + 1), nil
}

func day18p01(r io.Reader) (string, error) {
	return lagoonArea(r, func(instruction digInstruction) digStep { return instruction.plan })
}

func day18p02(r io.Reader) (string, error) {
	return lagoonArea(r, func(instruction digInstruction) digStep { return instruction.colour })
}
