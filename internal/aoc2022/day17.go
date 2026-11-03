package aoc2022

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type rockChamber struct {
	rows []byte
	jets []byte
	jet  int
	base int64
	work int
}

func (c *rockChamber) blocked(shape []grid.Position2D[int], x, y int) bool {
	for _, point := range shape {
		px, py := x+point.X, y+point.Y
		if px < 0 || px >= 7 || py < 0 || py < len(c.rows) && c.rows[py]&(1<<px) != 0 {
			return true
		}
	}
	return false
}

func (c *rockChamber) drop(shape []grid.Position2D[int]) error {
	x, y := 2, len(c.rows)+3
	for {
		c.work++
		if c.work > 100_000_000 {
			return fmt.Errorf("rock simulation exceeds 100000000 steps and air visits")
		}
		dx := -1
		if c.jets[c.jet] == '>' {
			dx = 1
		}
		c.jet = (c.jet + 1) % len(c.jets)
		if !c.blocked(shape, x+dx, y) {
			x += dx
		}
		if c.blocked(shape, x, y-1) {
			break
		}
		y--
	}
	for _, point := range shape {
		px, py := x+point.X, y+point.Y
		for len(c.rows) <= py {
			c.rows = append(c.rows, 0)
		}
		c.rows[py] |= 1 << px
	}
	if len(c.rows) > 100_000 {
		return fmt.Errorf("rock surface exceeds 100000 retained rows")
	}
	return c.trim()
}

func (c *rockChamber) trim() error {
	visited := make([]byte, len(c.rows)+1)
	queue := collections.NewDeque[int](128)
	top, lowest := len(c.rows), len(c.rows)
	visited[top] = 127
	for x := range 7 {
		queue.PushBack(top*7 + x)
	}
	for queue.Size() > 0 {
		cell, _ := queue.PopFront()
		x, y := cell%7, cell/7
		lowest = min(lowest, y)
		c.work++
		if c.work > 100_000_000 {
			return fmt.Errorf("rock simulation exceeds 100000000 steps and air visits")
		}
		for _, next := range [4]grid.Position2D[int]{{X: x - 1, Y: y}, {X: x + 1, Y: y}, {X: x, Y: y - 1}, {X: x, Y: y + 1}} {
			if next.X < 0 || next.X >= 7 || next.Y < 0 || next.Y > top {
				continue
			}
			bit := byte(1 << next.X)
			if visited[next.Y]&bit != 0 || next.Y < top && c.rows[next.Y]&bit != 0 {
				continue
			}
			visited[next.Y] |= bit
			queue.PushBack(next.Y*7 + next.X)
		}
	}
	// Future rock cells cannot enter air disconnected from the sky; rocks only close paths.
	if lowest > 0 {
		c.rows = slices.Clone(c.rows[lowest:])
		c.base += int64(lowest)
	}
	return nil
}

func fallingRockHeight(r io.Reader, target int64) (string, error) {
	data, err := io.ReadAll(io.LimitReader(r, 1_000_001))
	if err != nil {
		return "", err
	}
	if len(data) > 1_000_000 {
		return "", fmt.Errorf("jet pattern exceeds 1000000 bytes")
	}
	jets := bytes.TrimSpace(data)
	if len(jets) == 0 || len(bytes.Trim(jets, "<>")) != 0 {
		return "", fmt.Errorf("jet pattern must contain only < and >")
	}
	if target < 0 || target > 1_000_000_000_000 {
		return "", fmt.Errorf("rock count must be in 0..1000000000000")
	}
	shapes := [][]grid.Position2D[int]{
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}},
		{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 1, Y: 2}},
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 1}, {X: 2, Y: 2}},
		{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}},
		{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}},
	}
	type state struct {
		surface string
		shape   int
		jet     int
	}
	type occurrence struct {
		count, height int64
	}
	chamber := rockChamber{jets: jets}
	seen := map[state]occurrence{{}: {}}
	var count, skippedHeight int64
	simulated, stored := 0, 0
	skipped := false
	for count < target {
		if simulated == 200_000 {
			return "", fmt.Errorf("rock simulation exceeded 200000 rocks without a sufficient cycle")
		}
		if err := chamber.drop(shapes[count%int64(len(shapes))]); err != nil {
			return "", err
		}
		count++
		simulated++
		if skipped {
			continue
		}
		key := state{string(chamber.rows), int(count % int64(len(shapes))), chamber.jet}
		height := chamber.base + int64(len(chamber.rows))
		if previous, exists := seen[key]; exists {
			period := count - previous.count
			cycles := (target - count) / period
			if cycles > 0 {
				count += cycles * period
				skippedHeight += cycles * (height - previous.height)
				skipped = true
			}
		} else {
			stored += len(key.surface)
			if stored > 64_000_000 {
				return "", fmt.Errorf("rock cycle history exceeds 64000000 surface bytes")
			}
			seen[key] = occurrence{count, height}
		}
	}
	return strconv.FormatInt(chamber.base+int64(len(chamber.rows))+skippedHeight, 10), nil
}

func day17p01(r io.Reader) (string, error) {
	return fallingRockHeight(r, 2022)
}

func day17p02(r io.Reader) (string, error) {
	return fallingRockHeight(r, 1_000_000_000_000)
}
