package aoc2015

import (
	"fmt"
	"io"
	"math"
	"math/big"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func containerCombinations(r io.Reader, target int) (all, minimal *big.Int, err error) {
	ways := make([]big.Int, target+1)
	leastWays := make([]big.Int, target+1)
	leastCount := make([]int, target+1)
	ways[0].SetInt64(1)
	leastWays[0].SetInt64(1)
	for volume := 1; volume <= target; volume++ {
		leastCount[volume] = math.MaxInt
	}
	s := scanner.NewScanner(r, convert.ScanNumber[int])
	containers := 0
	for capacity := range s.Values() {
		if capacity <= 0 {
			return nil, nil, fmt.Errorf("container capacity must be positive")
		}
		containers++
		// Descending volumes use each physical container once, including equal-sized ones.
		for volume := target; volume >= capacity; volume-- {
			previous := volume - capacity
			ways[volume].Add(&ways[volume], &ways[previous])
			if leastCount[previous] == math.MaxInt {
				continue
			}
			count := leastCount[previous] + 1
			switch {
			case count < leastCount[volume]:
				leastCount[volume] = count
				leastWays[volume].Set(&leastWays[previous])
			case count == leastCount[volume]:
				leastWays[volume].Add(&leastWays[volume], &leastWays[previous])
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	if containers == 0 {
		return nil, nil, fmt.Errorf("empty container list")
	}
	return &ways[target], &leastWays[target], nil
}

func day17p01(r io.Reader) (string, error) {
	all, _, err := containerCombinations(r, 150)
	if err != nil {
		return "", err
	}
	return all.String(), nil
}

func day17p02(r io.Reader) (string, error) {
	_, minimal, err := containerCombinations(r, 150)
	if err != nil {
		return "", err
	}
	return minimal.String(), nil
}
