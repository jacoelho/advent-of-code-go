package aoc2017

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func followTubes(r io.Reader) (string, int, error) {
	tubes := make(map[grid.Position2D[int]]byte)
	var start grid.Position2D[int]
	starts, y := 0, 0
	s := bufio.NewScanner(r)
	for s.Scan() {
		for x, tile := range s.Bytes() {
			if tile == ' ' {
				continue
			}
			if tile != '|' && tile != '-' && tile != '+' && (tile < 'A' || tile > 'Z') {
				return "", 0, fmt.Errorf("invalid tube tile %q", tile)
			}
			pos := grid.Position2D[int]{X: x, Y: y}
			tubes[pos] = tile
			if y == 0 && tile == '|' {
				start = pos
				starts++
			}
		}
		y++
	}
	if err := s.Err(); err != nil {
		return "", 0, err
	}
	if starts != 1 {
		return "", 0, fmt.Errorf("expected one tube entrance, got %d", starts)
	}

	type state struct{ pos, dir grid.Position2D[int] }
	visited := collections.NewSet[state]()
	pos, dir := start, grid.Position2D[int]{Y: 1}
	var letters strings.Builder
	steps := 0
	for tubes[pos] != 0 {
		current := state{pos: pos, dir: dir}
		if visited.Contains(current) {
			return "", 0, fmt.Errorf("tube path loops at %v", pos)
		}
		visited.Add(current)
		steps++

		tile := tubes[pos]
		if tile >= 'A' && tile <= 'Z' {
			letters.WriteByte(tile)
		}
		if tile == '+' {
			turns := 0
			var next grid.Position2D[int]
			for _, candidate := range []grid.Position2D[int]{dir.TurnLeft(), dir.TurnRight()} {
				if tubes[pos.Add(candidate)] != 0 {
					next = candidate
					turns++
				}
			}
			if turns != 1 {
				return "", 0, fmt.Errorf("expected one turn at %v, got %d", pos, turns)
			}
			dir = next
		}
		pos = pos.Add(dir)
	}

	return letters.String(), steps, nil
}

func day19p01(r io.Reader) (string, error) {
	letters, _, err := followTubes(r)
	return letters, err
}

func day19p02(r io.Reader) (string, error) {
	_, steps, err := followTubes(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(steps), nil
}
