package aoc2024

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func lists(r io.Reader) ([]int, []int, error) {
	s := scanner.NewScanner(r, func(b []byte) ([2]int, error) {
		numbers := bytes.Fields(b)
		if len(numbers) != 2 {
			return [2]int{}, fmt.Errorf("location row %q needs two numbers", b)
		}
		var pair [2]int
		for i, number := range numbers {
			value, err := convert.ScanNumber[int](number)
			if err != nil {
				return [2]int{}, fmt.Errorf("location column %d: %w", i+1, err)
			}
			if value < 0 {
				return [2]int{}, fmt.Errorf("location column %d: negative location ID %d", i+1, value)
			}
			pair[i] = value
		}
		return pair, nil
	})

	var left, right []int
	for pair := range s.Values() {
		left = append(left, pair[0])
		right = append(right, pair[1])
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	return left, right, nil
}

func day01p01(r io.Reader) (string, error) {
	left, right, err := lists(r)
	if err != nil {
		return "", err
	}

	slices.Sort(left)
	slices.Sort(right)

	var total int
	for i := range left {
		distance := xmath.Abs(right[i] - left[i])
		sum, ok := xmath.CheckedAdd(total, distance)
		if !ok {
			return "", fmt.Errorf("total location distance overflows int")
		}
		total = sum
	}
	return strconv.Itoa(total), nil
}

func day01p02(r io.Reader) (string, error) {
	left, right, err := lists(r)
	if err != nil {
		return "", err
	}

	frequencies := make(map[int]int, len(right))
	for _, id := range right {
		frequencies[id]++
	}

	var total int
	for _, id := range left {
		score, ok := xmath.CheckedMul(id, frequencies[id])
		if ok {
			total, ok = xmath.CheckedAdd(total, score)
		}
		if !ok {
			return "", fmt.Errorf("location similarity score overflows int")
		}
	}
	return strconv.Itoa(total), nil
}
