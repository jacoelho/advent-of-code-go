package aoc2025

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const forkliftLimit = 4

// paperRolls records, for every roll, how many of its eight neighbours hold a
// roll. Positions without a roll are absent.
type paperRolls map[int]int

func parseRollsPaper(r io.Reader) (paperRolls, int, error) {
	s := scanner.NewScanner(r, func(line []byte) (string, error) {
		return string(line), nil
	})
	var rows []string
	width := 0
	for row := range s.Values() {
		rows = append(rows, row)
		width = max(width, len(row))
	}
	if err := s.Err(); err != nil {
		return nil, 0, err
	}
	stride := width + 2
	rolls := make(paperRolls)
	for y, row := range rows {
		for x := range len(row) {
			if row[x] == '@' {
				rolls[(y+1)*stride+x+1] = 0
			}
		}
	}
	for position := range rolls {
		for _, offset := range rollOffsets(stride) {
			if _, roll := rolls[position+offset]; roll {
				rolls[position]++
			}
		}
	}
	return rolls, stride, nil
}

// rollOffsets returns the index offsets of the eight neighbours in a grid with
// a one-cell border, so neighbours of edge cells never wrap to another row.
func rollOffsets(stride int) [8]int {
	return [8]int{-stride - 1, -stride, -stride + 1, -1, 1, stride - 1, stride, stride + 1}
}

func day04p01(r io.Reader) (string, error) {
	rolls, _, err := parseRollsPaper(r)
	if err != nil {
		return "", err
	}
	accessible := 0
	for _, neighbours := range rolls {
		if neighbours < forkliftLimit {
			accessible++
		}
	}
	return strconv.Itoa(accessible), nil
}

// day04p02 removes accessible rolls until none remain. Removal only lowers
// neighbour counts, so the removed set does not depend on order, and each roll
// is queued once: initially, or when its count drops to the limit minus one.
func day04p02(r io.Reader) (string, error) {
	rolls, stride, err := parseRollsPaper(r)
	if err != nil {
		return "", err
	}
	var queue []int
	for position, neighbours := range rolls {
		if neighbours < forkliftLimit {
			queue = append(queue, position)
		}
	}
	removed := 0
	for len(queue) > 0 {
		position := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		delete(rolls, position)
		removed++
		for _, offset := range rollOffsets(stride) {
			neighbour := position + offset
			if count, roll := rolls[neighbour]; roll {
				rolls[neighbour] = count - 1
				if count == forkliftLimit {
					queue = append(queue, neighbour)
				}
			}
		}
	}
	return strconv.Itoa(removed), nil
}
