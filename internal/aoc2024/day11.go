package aoc2024

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func parseStoneArrangement(r io.Reader) (map[int64]int64, error) {
	s := scanner.NewScannerWithSplit(r, bufio.ScanWords, func(token []byte) (int64, error) {
		stone, err := convert.ScanNumber[int64](token)
		if err != nil {
			return 0, fmt.Errorf("stone number: %w", err)
		}
		if stone < 0 {
			return 0, fmt.Errorf("stone number %d must be nonnegative", stone)
		}
		return stone, nil
	})
	stones := make(map[int64]int64)
	for stone := range s.Values() {
		stones[stone]++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(stones) == 0 {
		return nil, fmt.Errorf("stone arrangement is empty")
	}
	return stones, nil
}

// day11SplitStone returns the halves of a stone with an even number of digits.
func day11SplitStone(stone int64) (int64, int64, bool) {
	// Stones in [low, 10*low) have an even digit count and split at half.
	for low, half := int64(10), int64(10); ; low, half = low*100, half*10 {
		if stone < low {
			return 0, 0, false
		}
		high := low * 10
		if stone < high {
			return stone / half, stone % half, true
		}
		if high > math.MaxInt64/100 {
			return 0, 0, false
		}
	}
}

func day11AddCount(total, count int64) (int64, error) {
	sum, ok := xmath.CheckedAdd(total, count)
	if !ok {
		return 0, fmt.Errorf("stone count overflows int64")
	}
	return sum, nil
}

func day11CountStones(stones map[int64]int64, blinks int) (int64, error) {
	next := make(map[int64]int64, len(stones))
	add := func(stone, count int64) error {
		sum, err := day11AddCount(next[stone], count)
		next[stone] = sum
		return err
	}
	for range blinks {
		clear(next)
		for stone, count := range stones {
			var err error
			switch left, right, split := day11SplitStone(stone); {
			case stone == 0:
				err = add(1, count)
			case split:
				err = add(left, count)
				if err == nil {
					err = add(right, count)
				}
			default:
				if engraved, ok := xmath.CheckedMul(stone, 2024); ok {
					err = add(engraved, count)
				} else {
					err = fmt.Errorf("stone %d times 2024 overflows int64", stone)
				}
			}
			if err != nil {
				return 0, err
			}
		}
		stones, next = next, stones
	}
	var total int64
	for _, count := range stones {
		var err error
		if total, err = day11AddCount(total, count); err != nil {
			return 0, err
		}
	}
	return total, nil
}

func day11p01(r io.Reader) (string, error) {
	stones, err := parseStoneArrangement(r)
	if err != nil {
		return "", err
	}
	total, err := day11CountStones(stones, 25)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(total, 10), nil
}

func day11p02(r io.Reader) (string, error) {
	stones, err := parseStoneArrangement(r)
	if err != nil {
		return "", err
	}
	total, err := day11CountStones(stones, 75)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(total, 10), nil
}
