package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type lumberArea struct {
	width, height int
	acres         []byte
}

func parseLumberArea(r io.Reader) (lumberArea, error) {
	var area lumberArea
	s := bufio.NewScanner(r)
	for s.Scan() {
		row := s.Bytes()
		if area.height == 0 {
			area.width = len(row)
		}
		if len(row) == 0 || len(row) != area.width {
			return lumberArea{}, fmt.Errorf("lumber area must be a nonempty rectangle")
		}
		for _, acre := range row {
			if acre != '.' && acre != '|' && acre != '#' {
				return lumberArea{}, fmt.Errorf("invalid acre %q", acre)
			}
		}
		if len(area.acres)+len(row) > 1_000_000 {
			return lumberArea{}, fmt.Errorf("lumber area exceeds 1000000 acres")
		}
		area.acres = append(area.acres, row...)
		area.height++
	}
	if err := s.Err(); err != nil {
		return lumberArea{}, err
	}
	if area.height == 0 {
		return lumberArea{}, fmt.Errorf("empty lumber area")
	}

	return area, nil
}

func lumberResource(r io.Reader, minutes int) (int, error) {
	if minutes < 0 || minutes > 1_000_000_000 {
		return 0, fmt.Errorf("lumber minutes must be between 0 and 1000000000")
	}
	area, err := parseLumberArea(r)
	if err != nil {
		return 0, err
	}

	seen := make(map[string]int)
	stateLimit := (64 << 20) / (len(area.acres) + 64)
	offsets := grid.OffsetsNeighbours8[int]()
	next := make([]byte, len(area.acres))
	for minute := 0; minute < minutes; {
		key := string(area.acres)
		if previous, ok := seen[key]; ok {
			period := minute - previous
			minute += (minutes - minute) / period * period
			if minute == minutes {
				break
			}
		} else {
			if len(seen) == stateLimit {
				return 0, fmt.Errorf("lumber cycle search exceeded %d states", stateLimit)
			}
			seen[key] = minute
		}

		for i, acre := range area.acres {
			x, y := i%area.width, i/area.width
			trees, yards := 0, 0
			for _, offset := range offsets {
				nx, ny := x+offset.X, y+offset.Y
				if nx < 0 || ny < 0 || nx >= area.width || ny >= area.height {
					continue
				}
				switch area.acres[ny*area.width+nx] {
				case '|':
					trees++
				case '#':
					yards++
				}
			}
			next[i] = acre
			switch {
			case acre == '.' && trees >= 3:
				next[i] = '|'
			case acre == '|' && yards >= 3:
				next[i] = '#'
			case acre == '#' && (trees == 0 || yards == 0):
				next[i] = '.'
			}
		}
		area.acres, next = next, area.acres
		minute++
	}

	trees, yards := 0, 0
	for _, acre := range area.acres {
		switch acre {
		case '|':
			trees++
		case '#':
			yards++
		}
	}

	return trees * yards, nil
}

func day18p01(r io.Reader) (string, error) {
	value, err := lumberResource(r, 10)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(value), nil
}

func day18p02(r io.Reader) (string, error) {
	value, err := lumberResource(r, 1_000_000_000)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(value), nil
}
