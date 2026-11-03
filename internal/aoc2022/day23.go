package aoc2022

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func parseElves(r io.Reader) (collections.Set[grid.Position2D[int]], error) {
	elves := collections.NewSet[grid.Position2D[int]]()
	s := bufio.NewScanner(r)
	width, height := 0, 0
	for s.Scan() {
		row := s.Bytes()
		if height == 0 {
			width = len(row)
		}
		if len(row) == 0 || len(row) != width {
			return nil, fmt.Errorf("elf grove must be a nonempty rectangle")
		}
		if height >= 1_000_000/width {
			return nil, fmt.Errorf("elf grove exceeds one million cells")
		}
		for x, c := range row {
			switch c {
			case '#':
				elves.Add(grid.Position2D[int]{X: x, Y: height})
			case '.':
			default:
				return nil, fmt.Errorf("invalid elf grove cell %q", c)
			}
		}
		height++
		if len(elves) > 100_000 {
			return nil, fmt.Errorf("elf grove exceeds one hundred thousand elves")
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if height == 0 {
		return nil, fmt.Errorf("empty elf grove")
	}
	return elves, nil
}

func elfRound(elves collections.Set[grid.Position2D[int]], round int) bool {
	type proposal struct {
		from  grid.Position2D[int]
		count int
	}
	proposals := make(map[grid.Position2D[int]]proposal)
	directions := [4][3]grid.Position2D[int]{
		{{X: 0, Y: -1}, {X: -1, Y: -1}, {X: 1, Y: -1}},
		{{X: 0, Y: 1}, {X: -1, Y: 1}, {X: 1, Y: 1}},
		{{X: -1, Y: 0}, {X: -1, Y: -1}, {X: -1, Y: 1}},
		{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 1, Y: 1}},
	}
	for elf := range elves {
		isolated := true
		for _, offset := range grid.OffsetsNeighbours8[int]() {
			if elves.Contains(elf.Add(offset)) {
				isolated = false
				break
			}
		}
		if isolated {
			continue
		}
		for i := range 4 {
			direction := directions[(round+i)%4]
			clear := true
			for _, offset := range direction {
				if elves.Contains(elf.Add(offset)) {
					clear = false
					break
				}
			}
			if clear {
				destination := elf.Add(direction[0])
				candidate := proposals[destination]
				candidate.from = elf
				candidate.count++
				proposals[destination] = candidate
				break
			}
		}
	}
	moved := false
	for destination, candidate := range proposals {
		if candidate.count == 1 {
			// Every destination was empty in the old grove, so these updates cannot block each other.
			elves.Remove(candidate.from)
			elves.Add(destination)
			moved = true
		}
	}
	return moved
}

func day23p01(r io.Reader) (string, error) {
	elves, err := parseElves(r)
	if err != nil {
		return "", err
	}
	for round := range 10 {
		elfRound(elves, round)
	}
	if len(elves) == 0 {
		return "0", nil
	}
	var left, right, top, bottom int
	first := true
	for elf := range elves {
		if first {
			left, right, top, bottom = elf.X, elf.X, elf.Y, elf.Y
			first = false
		}
		left, right = min(left, elf.X), max(right, elf.X)
		top, bottom = min(top, elf.Y), max(bottom, elf.Y)
	}
	empty := int64(right-left+1)*int64(bottom-top+1) - int64(len(elves))
	return strconv.FormatInt(empty, 10), nil
}

func day23p02(r io.Reader) (string, error) {
	elves, err := parseElves(r)
	if err != nil {
		return "", err
	}
	work := 100_000_000
	for round := range 100_000 {
		work -= len(elves)
		if work < 0 {
			return "", fmt.Errorf("elf diffusion exceeds one hundred million elf visits")
		}
		if !elfRound(elves, round) {
			return strconv.Itoa(round + 1), nil
		}
	}
	return "", fmt.Errorf("elf diffusion exceeds one hundred thousand rounds")
}
