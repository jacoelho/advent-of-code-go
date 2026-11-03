package aoc2016

import (
	"bytes"
	"fmt"
	"io"
	"math/bits"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type officePosition struct{ x, y int }

func officeMaze(r io.Reader, destination officePosition, maxDistance int) (int, int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, 0, err
	}
	favorite, err := convert.ScanNumber[int32](bytes.TrimSpace(data))
	if err != nil || favorite < 0 {
		return 0, 0, fmt.Errorf("invalid favorite number %q", data)
	}
	isOpen := func(p officePosition) bool {
		if p.x < 0 || p.y < 0 {
			return false
		}
		x, y := uint64(p.x), uint64(p.y)
		value := x*x + 3*x + 2*x*y + y + y*y + uint64(favorite)
		return bits.OnesCount64(value)%2 == 0
	}
	start := officePosition{1, 1}
	if !isOpen(start) {
		return 0, 0, fmt.Errorf("maze starting position is a wall")
	}
	if maxDistance < 0 && !isOpen(destination) {
		return 0, 0, fmt.Errorf("maze destination is a wall")
	}
	queue := collections.NewDeque[officePosition](64)
	queue.PushBack(start)
	distances := map[officePosition]int{start: 0}
	for queue.Size() > 0 {
		position, _ := queue.PopFront()
		distance := distances[position]
		if maxDistance < 0 && position == destination {
			return distance, len(distances), nil
		}
		if maxDistance >= 0 && distance == maxDistance {
			continue
		}
		for _, offset := range [...]officePosition{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			next := officePosition{position.x + offset.x, position.y + offset.y}
			if _, seen := distances[next]; seen || !isOpen(next) {
				continue
			}
			if len(distances) == 1_000_000 {
				return 0, 0, fmt.Errorf("maze search exceeds 1000000 positions")
			}
			distances[next] = distance + 1
			queue.PushBack(next)
		}
	}
	if maxDistance < 0 {
		return 0, 0, fmt.Errorf("maze destination is unreachable")
	}
	return -1, len(distances), nil
}

func day13p01(r io.Reader) (string, error) {
	distance, _, err := officeMaze(r, officePosition{31, 39}, -1)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(distance), nil
}

func day13p02(r io.Reader) (string, error) {
	_, reachable, err := officeMaze(r, officePosition{}, 50)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(reachable), nil
}
