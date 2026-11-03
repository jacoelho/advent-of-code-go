package aoc2021

import (
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type reactorCuboid struct {
	low  [3]int64
	high [3]int64
}

type reactorStep struct {
	cuboid reactorCuboid
	on     bool
}

type reactorRegion uint8

const (
	initializationRegion reactorRegion = iota
	fullReactor
)

func intersectReactor(left, right reactorCuboid) (reactorCuboid, bool) {
	var result reactorCuboid
	for axis := range 3 {
		result.low[axis] = max(left.low[axis], right.low[axis])
		result.high[axis] = min(left.high[axis], right.high[axis])
		if result.low[axis] > result.high[axis] {
			return reactorCuboid{}, false
		}
	}
	return result, true
}

func subtractReactor(result []reactorCuboid, cuboid, removed reactorCuboid) []reactorCuboid {
	overlap, intersects := intersectReactor(cuboid, removed)
	if !intersects {
		return append(result, cuboid)
	}
	// Peel off disjoint slabs; the shrinking remainder ends at the overlap.
	for axis := range 3 {
		if cuboid.low[axis] < overlap.low[axis] {
			slab := cuboid
			slab.high[axis] = overlap.low[axis] - 1
			result = append(result, slab)
			cuboid.low[axis] = overlap.low[axis]
		}
		if cuboid.high[axis] > overlap.high[axis] {
			slab := cuboid
			slab.low[axis] = overlap.high[axis] + 1
			result = append(result, slab)
			cuboid.high[axis] = overlap.high[axis]
		}
	}
	return result
}

func rebootReactor(r io.Reader, region reactorRegion) (string, error) {
	pattern := regexp.MustCompile(`^(on|off) x=(-?\d+)\.\.(-?\d+),y=(-?\d+)\.\.(-?\d+),z=(-?\d+)\.\.(-?\d+)$`)
	s := scanner.NewScanner(r, func(line []byte) (reactorStep, error) {
		match := pattern.FindSubmatch(line)
		if match == nil {
			return reactorStep{}, fmt.Errorf("invalid reactor command %q", line)
		}
		step := reactorStep{on: string(match[1]) == "on"}
		for axis := range 3 {
			low, err := convert.ScanNumber[int32](match[2+2*axis])
			if err != nil {
				return reactorStep{}, fmt.Errorf("reactor lower coordinate: %w", err)
			}
			high, err := convert.ScanNumber[int32](match[3+2*axis])
			if err != nil {
				return reactorStep{}, fmt.Errorf("reactor upper coordinate: %w", err)
			}
			if low > high {
				return reactorStep{}, fmt.Errorf("reactor bounds are reversed")
			}
			step.cuboid.low[axis], step.cuboid.high[axis] = int64(low), int64(high)
		}
		return step, nil
	})

	var cuboids []reactorCuboid
	commands, work := 0, 0
	for step := range s.Values() {
		commands++
		if commands > 10_000 {
			return "", fmt.Errorf("reactor input exceeds 10000 commands")
		}
		if region == initializationRegion {
			var intersects bool
			step.cuboid, intersects = intersectReactor(step.cuboid, reactorCuboid{[3]int64{-50, -50, -50}, [3]int64{50, 50, 50}})
			if !intersects {
				continue
			}
		}
		var next []reactorCuboid
		for _, cuboid := range cuboids {
			work++
			if work > 20_000_000 {
				return "", fmt.Errorf("reactor reboot exceeds 20000000 intersections")
			}
			next = subtractReactor(next, cuboid, step.cuboid)
			if len(next) > 200_000 {
				return "", fmt.Errorf("reactor reboot exceeds 200000 disjoint cuboids")
			}
		}
		if step.on {
			next = append(next, step.cuboid)
		}
		if len(next) > 200_000 {
			return "", fmt.Errorf("reactor reboot exceeds 200000 disjoint cuboids")
		}
		cuboids = next
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if commands == 0 {
		return "", fmt.Errorf("empty reactor reboot")
	}
	var total int64
	for _, cuboid := range cuboids {
		volume := int64(1)
		for axis := range 3 {
			size := cuboid.high[axis] - cuboid.low[axis] + 1
			if volume > math.MaxInt64/size {
				return "", fmt.Errorf("reactor cuboid volume exceeds int64")
			}
			volume *= size
		}
		if total > math.MaxInt64-volume {
			return "", fmt.Errorf("reactor total volume exceeds int64")
		}
		total += volume
	}
	return strconv.FormatInt(total, 10), nil
}

func day22p01(r io.Reader) (string, error) {
	return rebootReactor(r, initializationRegion)
}

func day22p02(r io.Reader) (string, error) {
	return rebootReactor(r, fullReactor)
}
