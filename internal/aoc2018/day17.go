package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

const springX = 500

type clayVein struct {
	x0, x1, y0, y1 int
}

// reservoir is a ground scan whose columns span one beyond the clay and the
// spring on each side, which is as far as overflowing water can reach.
type reservoir struct {
	tiles       []byte
	left, width int
	top, bottom int
}

func parseReservoir(r io.Reader) (reservoir, error) {
	pattern := regexp.MustCompile(`^([xy])=([0-9]+), ([xy])=([0-9]+)\.\.([0-9]+)$`)
	var veins []clayVein
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
		if low > high {
			return reservoir{}, fmt.Errorf("invalid clay vein range %q", s.Text())
		}
		vein := clayVein{x0: fixed, x1: fixed, y0: low, y1: high}
		if fields[1][0] == 'y' {
			vein = clayVein{x0: low, x1: high, y0: fixed, y1: fixed}
		}
		veins = append(veins, vein)
	}
	if err := s.Err(); err != nil {
		return reservoir{}, err
	}
	if len(veins) == 0 {
		return reservoir{}, fmt.Errorf("empty clay scan")
	}

	left, right := springX, springX
	top, bottom := veins[0].y0, veins[0].y1
	for _, vein := range veins {
		left, right = min(left, vein.x0), max(right, vein.x1)
		top, bottom = min(top, vein.y0), max(bottom, vein.y1)
	}
	ground := reservoir{left: left - 1, width: right - left + 3, top: top, bottom: bottom}
	ground.tiles = make([]byte, ground.width*(bottom+1))
	for _, vein := range veins {
		for y := vein.y0; y <= vein.y1; y++ {
			for x := vein.x0; x <= vein.x1; x++ {
				ground.tiles[ground.index(x, y)] = '#'
			}
		}
	}
	return ground, nil
}

func (ground reservoir) index(x, y int) int {
	return y*ground.width + x - ground.left
}

// fill flows water from (x, y) and reports whether the water there is supported.
func (ground reservoir) fill(x, y int) bool {
	if y > ground.bottom {
		return false
	}
	switch ground.tiles[ground.index(x, y)] {
	case '#', '~':
		return true
	case '|':
		return false
	}
	ground.tiles[ground.index(x, y)] = '|'
	if !ground.fill(x, y+1) {
		return false
	}

	spread := func(direction int) (int, bool) {
		for edge := x; ; {
			next := ground.index(edge+direction, y)
			if ground.tiles[next] == '#' {
				return edge, true
			}
			edge += direction
			if ground.tiles[next] != '~' {
				ground.tiles[next] = '|'
			}
			if !ground.fill(edge, y+1) {
				return edge, false
			}
		}
	}
	left, leftWall := spread(-1)
	right, rightWall := spread(1)
	if !leftWall || !rightWall {
		return false
	}
	for edge := left; edge <= right; edge++ {
		ground.tiles[ground.index(edge, y)] = '~'
	}
	return true
}

func reservoirWater(r io.Reader) (int, int, error) {
	ground, err := parseReservoir(r)
	if err != nil {
		return 0, 0, err
	}
	ground.fill(springX, 0)

	reached, settled := 0, 0
	for _, tile := range ground.tiles[ground.top*ground.width:] {
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
