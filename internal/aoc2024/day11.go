package aoc2024

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func parseStoneArrangement(r io.Reader) ([]int64, error) {
	s := scanner.NewScannerWithSplit(r, bufio.ScanWords, func(token []byte) (int64, error) {
		stone, err := convert.ScanNumber[int64](token)
		if err != nil {
			return 0, fmt.Errorf("stone number: %w", err)
		}
		if stone < 0 {
			return 0, fmt.Errorf("stone number must be nonnegative")
		}
		return stone, nil
	})
	stones := slices.Collect(s.Values())
	return stones, s.Err()
}

func memoizeCount() func(int64, int64) int64 {
	cache := make(map[[2]int64]int64)

	var countRec func(stone int64, steps int64) int64
	countRec = func(stone int64, steps int64) int64 {
		if steps == 0 {
			return 1
		}

		k := [2]int64{stone, steps}
		if v, ok := cache[k]; ok {
			return v
		}

		if stone == 0 {
			cache[k] = countRec(1, steps-1)
			return cache[k]
		}
		digits := convert.ToDigits(stone)
		if len(digits)%2 != 0 {
			cache[k] = countRec(stone*2024, steps-1)
			return cache[k]
		}

		mid := len(digits) / 2
		n1 := convert.FromDigits(digits[:mid])
		n2 := convert.FromDigits(digits[mid:])
		cache[k] = countRec(n1, steps-1) + countRec(n2, steps-1)
		return cache[k]
	}

	return countRec
}

func countBlinks(stones []int64, steps int64) int64 {
	countFunc := memoizeCount()

	return xiter.FromSlice(stones).SumBy(func(stone int64) int64 {
		return countFunc(stone, steps)
	})
}

func day11p01(r io.Reader) (string, error) {
	stones, err := parseStoneArrangement(r)
	if err != nil {
		return "", err
	}

	total := countBlinks(stones, 25)

	return strconv.FormatInt(total, 10), nil
}

func day11p02(r io.Reader) (string, error) {
	stones, err := parseStoneArrangement(r)
	if err != nil {
		return "", err
	}

	total := countBlinks(stones, 75)

	return strconv.FormatInt(total, 10), nil
}
