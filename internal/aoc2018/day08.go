package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func licenseValues(r io.Reader) (int, int, error) {
	s := scanner.NewScannerWithSplit(r, bufio.ScanWords, convert.ScanNumber[int])
	var numbers []int
	for number := range s.Values() {
		if number < 0 {
			return 0, 0, fmt.Errorf("negative license value")
		}
		numbers = append(numbers, number)
	}
	if err := s.Err(); err != nil {
		return 0, 0, err
	}

	position := 0
	var node func() (int, int, error)
	node = func() (int, int, error) {
		if len(numbers)-position < 2 {
			return 0, 0, fmt.Errorf("incomplete node header")
		}
		children, metadata := numbers[position], numbers[position+1]
		position += 2
		if children > (len(numbers)-position)/2 {
			return 0, 0, fmt.Errorf("incomplete child nodes")
		}
		values := make([]int, children)
		var sum, value int
		for i := range children {
			childSum, childValue, err := node()
			if err != nil {
				return 0, 0, err
			}
			sum += childSum
			values[i] = childValue
		}
		if metadata > len(numbers)-position {
			return 0, 0, fmt.Errorf("incomplete metadata")
		}
		for _, entry := range numbers[position : position+metadata] {
			sum += entry
			if children == 0 {
				value += entry
			} else if entry > 0 && entry <= children {
				value += values[entry-1]
			}
		}
		position += metadata
		return sum, value, nil
	}

	sum, value, err := node()
	if err != nil {
		return 0, 0, err
	}
	if position != len(numbers) {
		return 0, 0, fmt.Errorf("trailing license data")
	}
	return sum, value, nil
}

func day08p01(r io.Reader) (string, error) {
	sum, _, err := licenseValues(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(sum), nil
}

func day08p02(r io.Reader) (string, error) {
	_, value, err := licenseValues(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(value), nil
}
