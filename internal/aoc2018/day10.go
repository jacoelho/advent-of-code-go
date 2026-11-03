package aoc2018

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/bits"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
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
		if len(points) > 2000 {
			return nil, fmt.Errorf("sky contains more than 2000 points")
		}
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

func skyAlignment(points []skyPoint) (int64, error) {
	var horizon int64
	for i, first := range points {
		for _, second := range points[i+1:] {
			for _, axis := range [][2]int64{{second.x - first.x, first.dx - second.dx}, {second.y - first.y, first.dy - second.dy}} {
				distance, velocity := axis[0], axis[1]
				if velocity < 0 {
					distance, velocity = -distance, -velocity
				}
				if velocity > 0 && distance > 0 {
					horizon = max(horizon, (distance+velocity-1)/velocity)
				}
			}
		}
	}
	if horizon > 2_000_000 || (horizon+1)*int64(len(points)) > 100_000_000 {
		return 0, fmt.Errorf("sky alignment exceeds its time or point-evaluation limit")
	}

	// After the final pairwise crossing, both bounding dimensions are nondecreasing.
	bestTime := int64(0)
	bestHigh, bestLow := ^uint64(0), ^uint64(0)
	for time := int64(0); time <= horizon; time++ {
		left, top, right, bottom := skyBounds(points, time)
		high, low := bits.Mul64(uint64(right-left+1), uint64(bottom-top+1))
		if high < bestHigh || high == bestHigh && low < bestLow {
			bestTime, bestHigh, bestLow = time, high, low
		}
	}

	return bestTime, nil
}

func day10p01(r io.Reader) (string, error) {
	points, err := parseSky(r)
	if err != nil {
		return "", err
	}
	time, err := skyAlignment(points)
	if err != nil {
		return "", err
	}

	left, top, right, bottom := skyBounds(points, time)
	width, height := right-left+1, bottom-top+1
	if width > 1_000_000 || height > 1_000_000 || width > 1_000_000/height {
		return "", fmt.Errorf("aligned sky raster exceeds 1000000 pixels")
	}
	pixels := bytes.Repeat([]byte{'.'}, int(width*height))
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
	time, err := skyAlignment(points)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(time, 10), nil
}
