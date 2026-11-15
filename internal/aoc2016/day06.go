package aoc2016

import (
	"io"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type messageFrequency int

const (
	mostCommon messageFrequency = iota
	leastCommon
)

func recoverMessage(r io.Reader, frequency messageFrequency) (string, error) {
	g, err := grid.ReadBytes(r, "abcdefghijklmnopqrstuvwxyz")
	if err != nil {
		return "", err
	}
	columns := make([][26]int, g.Width())
	for i, letter := range g.Cells() {
		columns[i%g.Width()][letter-'a']++
	}
	message := make([]byte, len(columns))
	for i, counts := range columns {
		best := -1
		for letter, count := range counts {
			if count == 0 {
				continue
			}
			if best < 0 || (frequency == mostCommon && count > counts[best]) || (frequency == leastCommon && count < counts[best]) {
				best = letter
			}
		}
		message[i] = byte(best) + 'a'
	}
	return string(message), nil
}

func day06p01(r io.Reader) (string, error) {
	return recoverMessage(r, mostCommon)
}

func day06p02(r io.Reader) (string, error) {
	return recoverMessage(r, leastCommon)
}
