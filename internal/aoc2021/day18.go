package aoc2021

import (
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type snailfishLeaf struct {
	value int64
	depth int
}

func parseSnailfish(line []byte) ([]snailfishLeaf, error) {
	var leaves []snailfishLeaf
	pos := 0
	var parse func(int) error
	parse = func(depth int) error {
		if pos == len(line) || depth > 64 || len(leaves) >= 1024 {
			return fmt.Errorf("truncated or oversized snailfish number")
		}
		if line[pos] == '[' {
			pos++
			if err := parse(depth + 1); err != nil {
				return err
			}
			if pos == len(line) || line[pos] != ',' {
				return fmt.Errorf("expected snailfish comma at byte %d", pos)
			}
			pos++
			if err := parse(depth + 1); err != nil {
				return err
			}
			if pos == len(line) || line[pos] != ']' {
				return fmt.Errorf("expected snailfish closing bracket at byte %d", pos)
			}
			pos++
			return nil
		}
		start := pos
		for pos < len(line) && line[pos] >= '0' && line[pos] <= '9' {
			pos++
		}
		value, err := convert.ScanNumber[int64](line[start:pos])
		if err != nil {
			return fmt.Errorf("snailfish regular number: %w", err)
		}
		leaves = append(leaves, snailfishLeaf{value, depth})
		return nil
	}
	if len(line) == 0 || line[0] != '[' {
		return nil, fmt.Errorf("snailfish number must be a pair")
	}
	if err := parse(0); err != nil {
		return nil, err
	}
	if pos != len(line) {
		return nil, fmt.Errorf("unexpected data after snailfish number")
	}
	return leaves, nil
}

func readSnailfish(r io.Reader) ([][]snailfishLeaf, error) {
	s := scanner.NewScanner(r, parseSnailfish)
	var numbers [][]snailfishLeaf
	for number := range s.Values() {
		if len(numbers) == 200 {
			return nil, fmt.Errorf("snailfish input exceeds 200 numbers")
		}
		numbers = append(numbers, number)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(numbers) == 0 {
		return nil, fmt.Errorf("empty snailfish input")
	}
	return numbers, nil
}

func explodeSnailfish(number []snailfishLeaf) ([]snailfishLeaf, bool, error) {
	for i := 0; i+1 < len(number); i++ {
		left, right := number[i], number[i+1]
		if left.depth <= 4 || left.depth != right.depth {
			continue
		}
		if i > 0 {
			if number[i-1].value > math.MaxInt64-left.value {
				return nil, false, fmt.Errorf("snailfish explosion overflows int64")
			}
			number[i-1].value += left.value
		}
		if i+2 < len(number) {
			if number[i+2].value > math.MaxInt64-right.value {
				return nil, false, fmt.Errorf("snailfish explosion overflows int64")
			}
			number[i+2].value += right.value
		}
		number[i] = snailfishLeaf{depth: left.depth - 1}
		return slices.Delete(number, i+1, i+2), true, nil
	}
	return number, false, nil
}

func addSnailfish(left, right []snailfishLeaf, work *int) ([]snailfishLeaf, error) {
	number := make([]snailfishLeaf, 0, len(left)+len(right))
	for _, operand := range [][]snailfishLeaf{left, right} {
		for _, leaf := range operand {
			leaf.depth++
			number = append(number, leaf)
		}
	}
	for {
		*work += len(number)
		if *work > 100_000_000 {
			return nil, fmt.Errorf("snailfish reduction exceeds 100000000 leaf visits")
		}
		var exploded bool
		var err error
		number, exploded, err = explodeSnailfish(number)
		if err != nil {
			return nil, err
		}
		if exploded {
			continue
		}
		split := false
		for i, leaf := range number {
			if leaf.value < 10 {
				continue
			}
			number[i] = snailfishLeaf{leaf.value / 2, leaf.depth + 1}
			number = slices.Insert(number, i+1, snailfishLeaf{leaf.value - leaf.value/2, leaf.depth + 1})
			split = true
			break
		}
		if !split {
			return number, nil
		}
	}
}

func snailfishMagnitude(number []snailfishLeaf) (int64, error) {
	stack := make([]snailfishLeaf, 0, len(number))
	for _, leaf := range number {
		stack = append(stack, leaf)
		for len(stack) >= 2 && stack[len(stack)-1].depth == stack[len(stack)-2].depth {
			left, right := stack[len(stack)-2], stack[len(stack)-1]
			if left.value > math.MaxInt64/3 || right.value > (math.MaxInt64-3*left.value)/2 {
				return 0, fmt.Errorf("snailfish magnitude overflows int64")
			}
			stack = stack[:len(stack)-1]
			stack[len(stack)-1] = snailfishLeaf{3*left.value + 2*right.value, left.depth - 1}
		}
	}
	return stack[0].value, nil
}

func day18p01(r io.Reader) (string, error) {
	numbers, err := readSnailfish(r)
	if err != nil {
		return "", err
	}
	sum, work := numbers[0], 0
	for _, number := range numbers[1:] {
		sum, err = addSnailfish(sum, number, &work)
		if err != nil {
			return "", err
		}
	}
	value, err := snailfishMagnitude(sum)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(value, 10), nil
}

func day18p02(r io.Reader) (string, error) {
	numbers, err := readSnailfish(r)
	if err != nil {
		return "", err
	}
	if len(numbers) < 2 {
		return "", fmt.Errorf("maximum snailfish pair requires two input numbers")
	}
	var largest int64
	work := 0
	for i, left := range numbers {
		for j, right := range numbers {
			if i == j {
				continue
			}
			sum, err := addSnailfish(left, right, &work)
			if err != nil {
				return "", err
			}
			value, err := snailfishMagnitude(sum)
			if err != nil {
				return "", err
			}
			largest = max(largest, value)
		}
	}
	return strconv.FormatInt(largest, 10), nil
}
