package aoc2021

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func probeTrajectories(r io.Reader) (int64, int64, error) {
	data, err := io.ReadAll(io.LimitReader(r, 1025))
	if err != nil {
		return 0, 0, err
	}
	pattern := regexp.MustCompile(`^target area: x=(-?\d+)\.\.(-?\d+), y=(-?\d+)\.\.(-?\d+)$`)
	match := pattern.FindSubmatch(bytes.TrimSpace(data))
	if len(data) > 1024 || match == nil {
		return 0, 0, fmt.Errorf("invalid target area")
	}
	var bounds [4]int64
	for i := range bounds {
		value, err := convert.ScanNumber[int32](match[i+1])
		if err != nil {
			return 0, 0, fmt.Errorf("target coordinate: %w", err)
		}
		bounds[i] = int64(value)
	}
	xmin, xmax, ymin, ymax := bounds[0], bounds[1], bounds[2], bounds[3]
	if xmin > xmax || ymin > ymax {
		return 0, 0, fmt.Errorf("target bounds are reversed")
	}
	if ymax >= 0 {
		return 0, 0, fmt.Errorf("target must lie wholly below the launch point")
	}
	// An initial vertical velocity of -ymin or more returns through y=0,
	// then falls below ymin in its next step.
	vxmin, vxmax := min(int64(0), xmin), max(int64(0), xmax)
	if vxmax-vxmin+1 > 5_000_000/(-2*ymin) {
		return 0, 0, fmt.Errorf("probe search exceeds 5000000 velocities")
	}
	var highest, hits int64
	work := 0
	for vx := vxmin; vx <= vxmax; vx++ {
		for vy := ymin; vy < -ymin; vy++ {
			x, y, dx, dy := int64(0), int64(0), vx, vy
			for {
				work++
				if work > 50_000_000 {
					return 0, 0, fmt.Errorf("probe search exceeds 50000000 steps")
				}
				x, y = x+dx, y+dy
				if dx > 0 {
					dx--
				} else if dx < 0 {
					dx++
				}
				dy--
				if x >= xmin && x <= xmax && y >= ymin && y <= ymax {
					hits++
					if vy > 0 {
						highest = max(highest, vy*(vy+1)/2)
					}
					break
				}
				if y < ymin || x < xmin && dx <= 0 || x > xmax && dx >= 0 {
					break
				}
			}
		}
	}
	return highest, hits, nil
}

func day17p01(r io.Reader) (string, error) {
	height, _, err := probeTrajectories(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(height, 10), nil
}

func day17p02(r io.Reader) (string, error) {
	_, count, err := probeTrajectories(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(count, 10), nil
}
