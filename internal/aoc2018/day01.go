package aoc2018

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func calibrationFrequencies(r io.Reader) ([]int64, error) {
	s := scanner.NewScanner(r, convert.ScanNumber[int64])
	frequencies := []int64{0}
	for change := range s.Values() {
		previous := frequencies[len(frequencies)-1]
		if change > 0 && previous > math.MaxInt64-change || change < 0 && previous < math.MinInt64-change {
			return nil, fmt.Errorf("calibration frequency overflows int64")
		}
		frequencies = append(frequencies, previous+change)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(frequencies) == 1 {
		return nil, fmt.Errorf("no frequency changes")
	}
	return frequencies, nil
}

func day01p01(r io.Reader) (string, error) {
	frequencies, err := calibrationFrequencies(r)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(frequencies[len(frequencies)-1], 10), nil
}

func day01p02(r io.Reader) (string, error) {
	frequencies, err := calibrationFrequencies(r)
	if err != nil {
		return "", err
	}

	seen := make(map[int64]bool)
	for _, frequency := range frequencies {
		if seen[frequency] {
			return strconv.FormatInt(frequency, 10), nil
		}
		seen[frequency] = true
	}

	total := frequencies[len(frequencies)-1]
	prefixes := frequencies[:len(frequencies)-1]
	drift := uint64(total)
	if total < 0 {
		drift = -drift
	}
	var frequency int64
	var firstCycle uint64
	firstIndex := -1
	// Every visit is prefix[index] + cycle*total; compare cycle/index pairs without cycling forever.
	for _, earlier := range prefixes {
		for j, later := range prefixes {
			if total > 0 && earlier <= later || total < 0 && earlier >= later {
				continue
			}
			difference := uint64(earlier) - uint64(later)
			if total < 0 {
				difference = -difference
			}
			if difference%drift != 0 {
				continue
			}
			cycle := difference / drift
			if firstIndex == -1 || cycle < firstCycle || cycle == firstCycle && j < firstIndex {
				frequency, firstCycle, firstIndex = earlier, cycle, j
			}
		}
	}
	if firstIndex == -1 {
		return "", fmt.Errorf("frequency never repeats")
	}
	return strconv.FormatInt(frequency, 10), nil
}
