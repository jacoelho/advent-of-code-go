package aoc2018

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type nanobot struct {
	pos    [3]int64
	radius int64
}

func parseNanobots(r io.Reader) ([]nanobot, error) {
	pattern := regexp.MustCompile(`^pos=<([+-]?\d+),([+-]?\d+),([+-]?\d+)>, r=(\d+)$`)
	s := scanner.NewScanner(r, func(line []byte) (nanobot, error) {
		fields := pattern.FindSubmatch(line)
		if fields == nil {
			return nanobot{}, fmt.Errorf("invalid nanobot %q", line)
		}
		var bot nanobot
		for i, field := range fields[1:] {
			value, err := convert.ScanNumber[int32](field)
			if err != nil {
				return nanobot{}, err
			}
			if i == 3 {
				bot.radius = int64(value)
			} else {
				bot.pos[i] = int64(value)
			}
		}
		return bot, nil
	})
	var bots []nanobot
	for bot := range s.Values() {
		if len(bots) == 10_000 {
			return nil, fmt.Errorf("nanobot search supports at most 10000 bots")
		}
		bots = append(bots, bot)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(bots) == 0 {
		return nil, fmt.Errorf("empty nanobot scan")
	}

	return bots, nil
}

func day23p01(r io.Reader) (string, error) {
	bots, err := parseNanobots(r)
	if err != nil {
		return "", err
	}

	strongest := bots[0]
	for _, bot := range bots[1:] {
		if bot.radius > strongest.radius {
			strongest = bot
		}
	}
	count := 0
	for _, bot := range bots {
		var distance int64
		for axis := range 3 {
			distance += xmath.Abs(bot.pos[axis] - strongest.pos[axis])
		}
		if distance <= strongest.radius {
			count++
		}
	}

	return strconv.Itoa(count), nil
}

type botCube struct {
	min      [3]int64
	size     int64
	coverage int
	distance int64
}

func (cube botCube) distanceTo(point [3]int64) int64 {
	var distance int64
	for axis := range 3 {
		low, high := cube.min[axis], cube.min[axis]+cube.size-1
		distance += max(low-point[axis], point[axis]-high, 0)
	}
	return distance
}

func day23p02(r io.Reader) (string, error) {
	bots, err := parseNanobots(r)
	if err != nil {
		return "", err
	}

	low, high := bots[0].pos, bots[0].pos
	for _, bot := range bots {
		for axis := range 3 {
			low[axis] = min(low[axis], bot.pos[axis]-bot.radius)
			high[axis] = max(high[axis], bot.pos[axis]+bot.radius)
		}
	}
	size := int64(1)
	for axis := range 3 {
		for size < high[axis]-low[axis]+1 {
			size *= 2
		}
	}

	queue := collections.NewHeap(func(a, b botCube) bool {
		if a.coverage != b.coverage {
			return a.coverage > b.coverage
		}
		if a.distance != b.distance {
			return a.distance < b.distance
		}
		return a.size < b.size
	})
	work := 0
	add := func(cube botCube) error {
		for _, bot := range bots {
			work++
			if work > 50_000_000 {
				return fmt.Errorf("nanobot search exceeded 50000000 intersection checks")
			}
			if cube.distanceTo(bot.pos) <= bot.radius {
				cube.coverage++
			}
		}
		if cube.coverage == 0 {
			return nil
		}
		if queue.Len() == 1_000_000 {
			return fmt.Errorf("nanobot search exceeded 1000000 queued cubes")
		}
		cube.distance = cube.distanceTo([3]int64{})
		queue.Push(cube)
		return nil
	}
	if err := add(botCube{min: low, size: size}); err != nil {
		return "", err
	}

	for cube := range queue.PopSeq() {
		// Queue priorities bound coverage above and distance below for every point in a cube.
		if cube.size == 1 {
			return strconv.FormatInt(cube.distance, 10), nil
		}
		half := cube.size / 2
		for octant := range 8 {
			child := botCube{min: cube.min, size: half}
			for axis := range 3 {
				if octant&(1<<axis) != 0 {
					child.min[axis] += half
				}
			}
			if err := add(child); err != nil {
				return "", err
			}
		}
	}

	return "", fmt.Errorf("nanobot search found no covered point")
}
