package aoc2023

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

// brick spans the cells from low to high inclusive; they differ on at most one
// axis and low.z is at least 1.
type brick struct {
	low, high [3]int
}

func parseBrickEnd(s string) ([3]int, error) {
	var end [3]int
	fields := strings.Split(s, ",")
	if len(fields) != 3 {
		return end, fmt.Errorf("brick end %q needs three coordinates", s)
	}
	for i, field := range fields {
		value, err := strconv.Atoi(field)
		if err != nil || value < 0 {
			return end, fmt.Errorf("invalid brick coordinate %q", field)
		}
		end[i] = value
	}
	return end, nil
}

func parseBrick(line []byte) (brick, error) {
	first, second, ok := strings.Cut(string(line), "~")
	if !ok {
		return brick{}, fmt.Errorf("invalid brick %q", line)
	}
	a, err := parseBrickEnd(first)
	if err != nil {
		return brick{}, fmt.Errorf("brick %q: %w", line, err)
	}
	b, err := parseBrickEnd(second)
	if err != nil {
		return brick{}, fmt.Errorf("brick %q: %w", line, err)
	}
	var result brick
	axes := 0
	for i := range 3 {
		result.low[i], result.high[i] = min(a[i], b[i]), max(a[i], b[i])
		if a[i] != b[i] {
			axes++
		}
	}
	if axes > 1 {
		return brick{}, fmt.Errorf("brick %q is not a straight line", line)
	}
	if result.low[2] < 1 {
		return brick{}, fmt.Errorf("brick %q is below the ground", line)
	}
	return result, nil
}

// settleBricks drops the bricks and returns, for each settled brick, the
// bricks resting on it and the bricks it rests on.
func settleBricks(r io.Reader) (supports, supportedBy [][]int, err error) {
	s := scanner.NewScanner(r, parseBrick)
	bricks := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	slices.SortFunc(bricks, func(a, b brick) int { return cmp.Compare(a.low[2], b.low[2]) })

	type top struct{ z, id int }
	tops := make(map[[2]int]top)
	supports = make([][]int, len(bricks))
	supportedBy = make([][]int, len(bricks))
	for id, b := range bricks {
		rest := 1
		for x := b.low[0]; x <= b.high[0]; x++ {
			for y := b.low[1]; y <= b.high[1]; y++ {
				if t, ok := tops[[2]int{x, y}]; ok {
					rest = max(rest, t.z+1)
				}
			}
		}
		for x := b.low[0]; x <= b.high[0]; x++ {
			for y := b.low[1]; y <= b.high[1]; y++ {
				cell := [2]int{x, y}
				if t, ok := tops[cell]; ok && t.z+1 == rest && !slices.Contains(supportedBy[id], t.id) {
					supportedBy[id] = append(supportedBy[id], t.id)
					supports[t.id] = append(supports[t.id], id)
				}
				tops[cell] = top{z: rest + b.high[2] - b.low[2], id: id}
			}
		}
	}
	return supports, supportedBy, nil
}

func day22p01(r io.Reader) (string, error) {
	supports, supportedBy, err := settleBricks(r)
	if err != nil {
		return "", err
	}
	safe := 0
	for _, above := range supports {
		if !slices.ContainsFunc(above, func(id int) bool { return len(supportedBy[id]) == 1 }) {
			safe++
		}
	}
	return strconv.Itoa(safe), nil
}

func day22p02(r io.Reader) (string, error) {
	supports, supportedBy, err := settleBricks(r)
	if err != nil {
		return "", err
	}
	total := 0
	standing := make([]int, len(supportedBy))
	var queue []int
	for removed := range supports {
		for id, below := range supportedBy {
			standing[id] = len(below)
		}
		queue = append(queue[:0], removed)
		for head := 0; head < len(queue); head++ {
			for _, above := range supports[queue[head]] {
				standing[above]--
				if standing[above] == 0 {
					queue = append(queue, above)
				}
			}
		}
		total += len(queue) - 1
	}
	return strconv.Itoa(total), nil
}
