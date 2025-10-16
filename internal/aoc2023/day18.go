package aoc2023

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type instruction struct {
	direction byte
	distance  int
	color     string
}

func parseInstruction(line []byte) (instruction, error) {
	parts := strings.Fields(string(line))
	if len(parts) != 3 {
		return instruction{}, fmt.Errorf("invalid instruction: %s", line)
	}

	distance := aoc.MustAtoi(parts[1])
	color := strings.Trim(parts[2], "()")

	return instruction{
		direction: parts[0][0],
		distance:  distance,
		color:     color,
	}, nil
}

func decodeHexInstruction(color string) (byte, int, error) {
	if len(color) != 7 || color[0] != '#' {
		return 0, 0, fmt.Errorf("invalid color format: %s", color)
	}

	distanceHex := color[1:6]
	distance, err := strconv.ParseInt(distanceHex, 16, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid distance hex: %s", distanceHex)
	}

	directionHex := color[6]
	var direction byte
	switch directionHex {
	case '0':
		direction = 'R'
	case '1':
		direction = 'D'
	case '2':
		direction = 'L'
	case '3':
		direction = 'U'
	default:
		return 0, 0, fmt.Errorf("invalid direction hex: %c", directionHex)
	}

	return direction, int(distance), nil
}

func movePosition(pos grid.Position2D[int], direction byte, distance int) grid.Position2D[int] {
	switch direction {
	case 'R':
		pos.X += distance
	case 'L':
		pos.X -= distance
	case 'U':
		pos.Y -= distance
	case 'D':
		pos.Y += distance
	}
	return pos
}

type instructionDecoder func(instruction) (direction byte, distance int, err error)

func calculateLagoonArea(instructions []instruction, decoder instructionDecoder) (int, error) {
	vertices := make([]grid.Position2D[int], 0, len(instructions)+1)
	current := grid.Position2D[int]{X: 0, Y: 0}
	vertices = append(vertices, current)
	perimeter := 0

	for _, inst := range instructions {
		direction, distance, err := decoder(inst)
		if err != nil {
			return 0, err
		}
		perimeter += distance
		current = movePosition(current, direction, distance)
		vertices = append(vertices, current)
	}

	area := xmath.PolygonArea(vertices)

	// Pick's theorem: A = i + b/2 - 1
	// where A is area, i is interior points, b is boundary points
	// Total area = i + b = A + b/2 + 1
	return area + perimeter/2 + 1, nil
}

func part1Decoder(inst instruction) (byte, int, error) {
	return inst.direction, inst.distance, nil
}

func part2Decoder(inst instruction) (byte, int, error) {
	return decodeHexInstruction(inst.color)
}

func day18p01(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseInstruction)
	instructions := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return "", err
	}
	result, err := calculateLagoonArea(instructions, part1Decoder)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(result), nil
}

func day18p02(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseInstruction)
	instructions := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return "", err
	}
	result, err := calculateLagoonArea(instructions, part2Decoder)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(result), nil
}
