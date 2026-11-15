package aoc2018

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type rescueCave struct {
	depth   int
	target  grid.Position2D[int]
	erosion [][]uint16
}

func parseRescueCave(r io.Reader) (*rescueCave, error) {
	s := bufio.NewScanner(r)
	cave := &rescueCave{}
	line := 0
	for s.Scan() {
		switch line {
		case 0:
			text, ok := bytes.CutPrefix(s.Bytes(), []byte("depth: "))
			if !ok {
				return nil, fmt.Errorf("missing cave depth")
			}
			depth, err := convert.ScanNumber[int64](text)
			if err != nil {
				return nil, err
			}
			if depth < 0 {
				return nil, fmt.Errorf("negative cave depth")
			}
			cave.depth = int(depth % 20183)
		case 1:
			text, ok := bytes.CutPrefix(s.Bytes(), []byte("target: "))
			if !ok {
				return nil, fmt.Errorf("missing cave target")
			}
			target, err := grid.ParsePosition2D[int32](string(text))
			if err != nil {
				return nil, fmt.Errorf("cave target: %w", err)
			}
			if target.X < 0 || target.Y < 0 {
				return nil, fmt.Errorf("negative cave target %q", text)
			}
			cave.target = grid.Position2D[int]{X: int(target.X), Y: int(target.Y)}
		default:
			return nil, fmt.Errorf("unexpected cave input %q", s.Text())
		}
		line++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if line != 2 {
		return nil, fmt.Errorf("expected cave depth and target")
	}

	return cave, nil
}

func (cave *rescueCave) region(pos grid.Position2D[int]) int {
	if pos.Y < len(cave.erosion) && pos.X < len(cave.erosion[pos.Y]) {
		return int(cave.erosion[pos.Y][pos.X] % 3)
	}
	width, height := pos.X+1, max(pos.Y+1, len(cave.erosion))
	if len(cave.erosion) != 0 {
		width = max(width, len(cave.erosion[0]))
	}
	for len(cave.erosion) < height {
		cave.erosion = append(cave.erosion, nil)
	}
	for y := range height {
		for x := len(cave.erosion[y]); x < width; x++ {
			var index int64
			switch {
			case x == 0 && y == 0 || x == cave.target.X && y == cave.target.Y:
			case y == 0:
				index = int64(x) * 16807
			case x == 0:
				index = int64(y) * 48271
			default:
				index = int64(cave.erosion[y][x-1]) * int64(cave.erosion[y-1][x])
			}
			value := uint16((index + int64(cave.depth)) % 20183)
			cave.erosion[y] = append(cave.erosion[y], value)
		}
	}

	return int(cave.erosion[pos.Y][pos.X] % 3)
}

func day22p01(r io.Reader) (string, error) {
	cave, err := parseRescueCave(r)
	if err != nil {
		return "", err
	}
	cave.region(cave.target)

	risk := 0
	for _, row := range cave.erosion {
		for _, erosion := range row {
			risk += int(erosion % 3)
		}
	}

	return strconv.Itoa(risk), nil
}

func day22p02(r io.Reader) (string, error) {
	cave, err := parseRescueCave(r)
	if err != nil {
		return "", err
	}
	const torch = 1
	if cave.depth%3 == torch {
		return "", fmt.Errorf("torch is not permitted at the cave entrance and target")
	}

	type state struct {
		pos  grid.Position2D[int]
		tool int
	}
	edges := func(current state, yield func(state, int) bool) error {
		region := cave.region(current.pos)
		for tool := range 3 {
			if tool != current.tool && tool != region {
				if !yield(state{pos: current.pos, tool: tool}, 7) {
					return nil
				}
			}
		}
		for pos := range grid.Neighbours4(current.pos) {
			if pos.X < 0 || pos.Y < 0 {
				continue
			}
			if current.tool != cave.region(pos) {
				if !yield(state{pos: pos, tool: current.tool}, 1) {
					return nil
				}
			}
		}
		return nil
	}
	heuristic := func(current state) int {
		distance := current.pos.Distance(cave.target)
		if current.tool != torch {
			distance += 7
		}
		return distance
	}
	minutes, found, err := search.AStarDistance(state{tool: torch}, func(current state) bool {
		return current.pos == cave.target && current.tool == torch
	}, edges, heuristic, search.Limits{})
	if err != nil {
		return "", fmt.Errorf("cave search: %w", err)
	}
	if !found {
		return "", fmt.Errorf("cave target is unreachable with the torch")
	}

	return strconv.Itoa(minutes), nil
}
