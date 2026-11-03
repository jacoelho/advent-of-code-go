package aoc2024

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func lists(r io.Reader) ([]int, []int, error) {
	s := scanner.NewScanner(r, func(b []byte) ([2]int, error) {
		numbers := bytes.Fields(b)

		if len(numbers) != 2 {
			return [2]int{}, fmt.Errorf("location row needs two numbers")
		}
		var pair [2]int
		for i, number := range numbers {
			value, err := convert.ScanNumber[int](number)
			if err != nil {
				return [2]int{}, fmt.Errorf("location column %d: %w", i+1, err)
			}
			pair[i] = value
		}
		return pair, nil
	})

	l := xiter.From(s.Values()).Reduce(make([][]int, 2), func(sum [][]int, v [2]int) [][]int {
		sum[0] = append(sum[0], v[0])
		sum[1] = append(sum[1], v[1])
		return sum
	})
	if s.Err() != nil {
		return nil, nil, s.Err()
	}
	return l[0], l[1], nil
}

func day01p01(r io.Reader) (string, error) {
	left, right, err := lists(r)
	if err != nil {
		return "", err
	}

	slices.Sort(left)
	slices.Sort(right)

	pairs := xiter.Zip(slices.Values(left), slices.Values(right))

	result := xiter.From(pairs).SumBy(func(zipped xiter.Zipped[int, int]) int {
		return xmath.Abs(zipped.V2 - zipped.V1)
	})

	return strconv.Itoa(result), nil
}

func day01p02(r io.Reader) (string, error) {
	left, right, err := lists(r)
	if err != nil {
		return "", err
	}

	frequencies := xiter.Frequencies(slices.Values(right))

	result := xiter.FromSlice(left).SumBy(func(v int) int {
		return v * frequencies[v]
	})

	return strconv.Itoa(result), nil
}
