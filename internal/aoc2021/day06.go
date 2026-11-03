package aoc2021

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func commaSeparatedNumbers(r io.Reader) ([]int64, error) {
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte(",")), func(data []byte) (int64, error) {
		return convert.ScanNumber[int64](bytes.TrimSpace(data))
	})
	var numbers []int64
	for value := range s.Values() {
		numbers = append(numbers, value)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(numbers) == 0 {
		return nil, fmt.Errorf("empty number list")
	}
	return numbers, nil
}

func lanternfishPopulation(r io.Reader, days int) (uint64, error) {
	timers, err := commaSeparatedNumbers(r)
	if err != nil {
		return 0, err
	}
	var population [9]uint64
	for _, timer := range timers {
		if timer < 0 || timer > 8 {
			return 0, fmt.Errorf("invalid lanternfish timer %d", timer)
		}
		population[timer]++
	}
	for range days {
		births := population[0]
		copy(population[:8], population[1:])
		population[8] = births
		if population[6] > math.MaxUint64-births {
			return 0, fmt.Errorf("lanternfish population overflows uint64")
		}
		population[6] += births
	}
	var total uint64
	for _, count := range population {
		if total > math.MaxUint64-count {
			return 0, fmt.Errorf("lanternfish population overflows uint64")
		}
		total += count
	}
	return total, nil
}

func day06p01(r io.Reader) (string, error) {
	population, err := lanternfishPopulation(r, 80)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(population, 10), nil
}

func day06p02(r io.Reader) (string, error) {
	population, err := lanternfishPopulation(r, 256)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(population, 10), nil
}
