package aoc2018

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type skyPoint struct {
	x, y, dx, dy int64
}

func parseSky(r io.Reader) ([]skyPoint, error) {
	pattern := regexp.MustCompile(`^position=<\s*(-?\d+),\s*(-?\d+)> velocity=<\s*(-?\d+),\s*(-?\d+)>$`)
	var points []skyPoint
	s := bufio.NewScanner(r)
	for s.Scan() {
		match := pattern.FindStringSubmatch(s.Text())
		if match == nil {
			return nil, fmt.Errorf("invalid sky point %q", s.Text())
		}
		var values [4]int64
		for i := range values {
			value, err := strconv.ParseInt(match[i+1], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid sky coordinate or velocity: %w", err)
			}
			values[i] = value
		}
		points = append(points, skyPoint{values[0], values[1], values[2], values[3]})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(points) == 0 {
		return nil, fmt.Errorf("empty sky")
	}

	return points, nil
}

func skyBounds(points []skyPoint, time int64) (left, top, right, bottom int64) {
	left, top = points[0].x+time*points[0].dx, points[0].y+time*points[0].dy
	right, bottom = left, top
	for _, point := range points[1:] {
		x, y := point.x+time*point.dx, point.y+time*point.dy
		left, right = min(left, x), max(right, x)
		top, bottom = min(top, y), max(bottom, y)
	}

	return left, top, right, bottom
}

func skySpan(points []skyPoint, time int64) int64 {
	left, top, right, bottom := skyBounds(points, time)
	return right - left + bottom - top
}

// skyAlignment returns the first second at which width plus height of the
// bounding box is smallest. Each dimension is a maximum minus a minimum of
// linear functions of time, so the sum is convex: the first second that does
// not improve on the next one is a global minimum.
func skyAlignment(points []skyPoint) int64 {
	time := int64(0)
	for skySpan(points, time+1) < skySpan(points, time) {
		time++
	}
	return time
}

func day10p01(r io.Reader) (string, error) {
	points, err := parseSky(r)
	if err != nil {
		return "", err
	}
	time := skyAlignment(points)

	left, top, right, bottom := skyBounds(points, time)
	width, height := right-left+1, bottom-top+1
	area, ok := xmath.CheckedMul(width, height)
	if !ok {
		return "", fmt.Errorf("aligned sky raster %dx%d overflows int", width, height)
	}
	pixels := bytes.Repeat([]byte{'.'}, int(area))
	for _, point := range points {
		x, y := point.x+time*point.dx-left, point.y+time*point.dy-top
		pixels[y*width+x] = '#'
	}
	rows := make([]string, int(height))
	for y := range rows {
		rows[y] = string(pixels[int64(y)*width : int64(y+1)*width])
	}

	message, err := aoc.DecodeSeparatedDisplay(rows)
	if err != nil {
		return "", fmt.Errorf("decode sky at second %d: %w", time, err)
	}
	return message, nil
}

func day10p02(r io.Reader) (string, error) {
	points, err := parseSky(r)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(skyAlignment(points), 10), nil
}
