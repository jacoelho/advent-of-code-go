package aoc2019

import (
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type direction uint8

const (
	directionUp direction = iota
	directionDown
	directionLeft
	directionRight
)

var directionMap = map[byte]direction{
	'U': directionUp,
	'D': directionDown,
	'L': directionLeft,
	'R': directionRight,
}

var directionsOffsets = map[direction]grid.Position2D[int]{
	directionUp:    {X: 0, Y: -1},
	directionDown:  {X: 0, Y: 1},
	directionLeft:  {X: -1, Y: 0},
	directionRight: {X: 1, Y: 0},
}

type instruction struct {
	direction direction
	distance  int
}

func parseWirePath(line string) ([]instruction, error) {
	segments := strings.Split(strings.TrimSpace(line), ",")
	instructions := make([]instruction, 0, len(segments))

	for _, seg := range segments {
		if len(seg) < 2 {
			return nil, fmt.Errorf("invalid segment: %s", seg)
		}

		dir, ok := directionMap[seg[0]]
		if !ok {
			return nil, fmt.Errorf("invalid direction: %c", seg[0])
		}

		distance, err := strconv.Atoi(seg[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid distance in segment %s: %w", seg, err)
		}

		instructions = append(instructions, instruction{direction: dir, distance: distance})
	}

	return instructions, nil
}

func parseWiresPaths(r io.Reader) ([]instruction, []instruction, error) {
	s := scanner.NewScanner(r, func(line []byte) ([]instruction, error) {
		return parseWirePath(string(line))
	})

	wires := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, nil, err
	}

	if len(wires) != 2 {
		return nil, nil, fmt.Errorf("expected 2 wires, got %d", len(wires))
	}

	return wires[0], wires[1], nil
}

func traceWire(instructions []instruction) map[grid.Position2D[int]]int {
	steps := make(map[grid.Position2D[int]]int)
	current := grid.Position2D[int]{X: 0, Y: 0}
	stepCount := 0

	for _, inst := range instructions {
		offset := directionsOffsets[inst.direction]
		for range inst.distance {
			current = current.Add(offset)
			stepCount++
			if _, exists := steps[current]; !exists {
				steps[current] = stepCount
			}
		}
	}

	return steps
}

func day3p01(r io.Reader) (string, error) {
	wire1, wire2, err := parseWiresPaths(r)
	if err != nil {
		return "", err
	}

	wire1Steps := traceWire(wire1)
	wire2Steps := traceWire(wire2)

	origin := grid.Position2D[int]{}
	minDistance := math.MaxInt

	for pos := range wire1Steps {
		if _, exists := wire2Steps[pos]; exists {
			minDistance = min(minDistance, origin.Distance(pos))
		}
	}

	return strconv.Itoa(minDistance), nil
}

func day3p02(r io.Reader) (string, error) {
	wire1, wire2, err := parseWiresPaths(r)
	if err != nil {
		return "", err
	}

	wire1Steps := traceWire(wire1)
	wire2Steps := traceWire(wire2)

	minSteps := math.MaxInt
	for pos, steps1 := range wire1Steps {
		if steps2, exists := wire2Steps[pos]; exists {
			minSteps = min(minSteps, steps1+steps2)
		}
	}

	return strconv.Itoa(minSteps), nil
}
