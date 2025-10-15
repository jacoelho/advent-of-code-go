package aoc2024

import (
	"bytes"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func lists(r io.Reader) ([]int, []int, error) {
	s := scanner.NewScanner(r, func(b []byte) ([2]int, error) {
		numbers := bytes.Fields(b)

		v1 := aoc.Must(strconv.Atoi(string(numbers[0])))
		v2 := aoc.Must(strconv.Atoi(string(numbers[1])))

		return [2]int{v1, v2}, nil
	})

	l := xiter.Reduce(func(sum [][]int, v [2]int) [][]int {
		sum[0] = append(sum[0], v[0])
		sum[1] = append(sum[1], v[1])
		return sum
	}, make([][]int, 2), s.Values())
	if s.Err() != nil {
		return nil, nil, s.Err()
	}
	return l[0], l[1], nil
}

func day01p01(r io.Reader) (string, error) {
	left, right := aoc.Must2(lists(r))

	slices.Sort(left)
	slices.Sort(right)

	pairs := xiter.Zip(slices.Values(left), slices.Values(right))

	result := xiter.Reduce(func(sum int, zipped xiter.Zipped[int, int]) int {
		return sum + xmath.Abs(zipped.V2-zipped.V1)
	}, 0, pairs)

	return strconv.Itoa(result), nil
}

func day01p02(r io.Reader) (string, error) {
	left, right := aoc.Must2(lists(r))

	frequencies := xslices.Frequencies(right)

	result := xiter.Reduce(func(sum int, v int) int {
		return sum + v*frequencies[v]
	}, 0, slices.Values(left))

	return strconv.Itoa(result), nil
}
