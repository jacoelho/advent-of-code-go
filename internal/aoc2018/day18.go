package aoc2018

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func lumberResource(r io.Reader, minutes int) (int, error) {
	area, err := grid.ReadBytes(r, ".|#")
	if err != nil {
		return 0, err
	}
	width, height, acres := area.Width(), area.Height(), area.Cells()

	seen := make(map[string]int)
	offsets := grid.OffsetsNeighbours8[int]()
	next := make([]byte, len(acres))
	for minute := 0; minute < minutes; {
		key := string(acres)
		if previous, ok := seen[key]; ok {
			period := minute - previous
			minute += (minutes - minute) / period * period
			if minute == minutes {
				break
			}
		} else {
			seen[key] = minute
		}

		for i, acre := range acres {
			x, y := i%width, i/width
			trees, yards := 0, 0
			for _, offset := range offsets {
				nx, ny := x+offset.X, y+offset.Y
				if nx < 0 || ny < 0 || nx >= width || ny >= height {
					continue
				}
				switch acres[ny*width+nx] {
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
		acres, next = next, acres
		minute++
	}

	trees, yards := 0, 0
	for _, acre := range acres {
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
