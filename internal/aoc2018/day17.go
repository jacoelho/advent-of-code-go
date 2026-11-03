package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type reservoir struct {
	tiles       map[grid.Position2D[int]]byte
	top, bottom int
}

func parseReservoir(r io.Reader) (reservoir, error) {
	pattern := regexp.MustCompile(`^([xy])=([0-9]+), ([xy])=([0-9]+)\.\.([0-9]+)$`)
	ground := reservoir{tiles: make(map[grid.Position2D[int]]byte), top: math.MaxInt}
	s := bufio.NewScanner(r)
	for s.Scan() {
		fields := pattern.FindSubmatch(s.Bytes())
		if fields == nil || fields[1][0] == fields[3][0] {
			return reservoir{}, fmt.Errorf("invalid clay vein %q", s.Text())
		}
		var values [3]int
		for i, field := range [][]byte{fields[2], fields[4], fields[5]} {
			value, err := convert.ScanNumber[int32](field)
			if err != nil {
				return reservoir{}, err
			}
			values[i] = int(value)
		}
		fixed, low, high := values[0], values[1], values[2]
		if low > high || high-low+1 > 2_000_000-len(ground.tiles) {
			return reservoir{}, fmt.Errorf("invalid or excessive clay vein %q", s.Text())
		}
		for variable := low; variable <= high; variable++ {
			pos := grid.Position2D[int]{X: fixed, Y: variable}
			if fields[1][0] == 'y' {
				pos.X, pos.Y = variable, fixed
			}
			if pos.Y > 100_000 {
				return reservoir{}, fmt.Errorf("reservoir depth exceeds 100000")
			}
			ground.tiles[pos] = '#'
			ground.top = min(ground.top, pos.Y)
			ground.bottom = max(ground.bottom, pos.Y)
		}
	}
	if err := s.Err(); err != nil {
		return reservoir{}, err
	}
	if len(ground.tiles) == 0 {
		return reservoir{}, fmt.Errorf("empty clay scan")
	}

	return ground, nil
}

func reservoirWater(r io.Reader) (int, int, error) {
	ground, err := parseReservoir(r)
	if err != nil {
		return 0, 0, err
	}

	work := 0
	var fill func(grid.Position2D[int]) (bool, error)
	fill = func(pos grid.Position2D[int]) (bool, error) {
		work++
		if work > 10_000_000 || len(ground.tiles) > 2_000_000 {
			return false, fmt.Errorf("reservoir exceeded its 10000000-operation or 2000000-cell limit")
		}
		if pos.Y > ground.bottom {
			return false, nil
		}
		switch ground.tiles[pos] {
		case '#', '~':
			return true, nil
		case '|':
			return false, nil
		}
		ground.tiles[pos] = '|'
		supported, err := fill(grid.Position2D[int]{X: pos.X, Y: pos.Y + 1})
		if err != nil || !supported {
			return false, err
		}

		spread := func(direction int) (int, bool, error) {
			x := pos.X
			for {
				next := grid.Position2D[int]{X: x + direction, Y: pos.Y}
				if ground.tiles[next] == '#' {
					return x, true, nil
				}
				x = next.X
				if ground.tiles[next] != '~' {
					ground.tiles[next] = '|'
				}
				supported, err := fill(grid.Position2D[int]{X: x, Y: pos.Y + 1})
				if err != nil || !supported {
					return x, false, err
				}
			}
		}
		left, leftWall, err := spread(-1)
		if err != nil {
			return false, err
		}
		right, rightWall, err := spread(1)
		if err != nil {
			return false, err
		}
		if !leftWall || !rightWall {
			return false, nil
		}
		for x := left; x <= right; x++ {
			ground.tiles[grid.Position2D[int]{X: x, Y: pos.Y}] = '~'
		}

		return true, nil
	}
	if _, err := fill(grid.Position2D[int]{X: 500}); err != nil {
		return 0, 0, err
	}

	reached, settled := 0, 0
	for pos, tile := range ground.tiles {
		if pos.Y < ground.top || pos.Y > ground.bottom {
			continue
		}
		switch tile {
		case '|':
			reached++
		case '~':
			reached++
			settled++
		}
	}

	return reached, settled, nil
}

func day17p01(r io.Reader) (string, error) {
	reached, _, err := reservoirWater(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(reached), nil
}

func day17p02(r io.Reader) (string, error) {
	_, settled, err := reservoirWater(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(settled), nil
}
