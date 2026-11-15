package aoc2022

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func groveCoordinates(r io.Reader, key int64, rounds int) (string, error) {
	s := scanner.NewScanner(r, convert.ScanNumber[int64])
	var values []int64
	zero := -1
	for value := range s.Values() {
		decrypted, ok := xmath.CheckedMul(value, key)
		if !ok {
			return "", fmt.Errorf("decryption key multiplication overflows int64")
		}
		if value == 0 {
			if zero >= 0 {
				return "", fmt.Errorf("mixing input must contain exactly one zero")
			}
			zero = len(values)
		}
		values = append(values, decrypted)
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(values) < 2 || zero < 0 {
		return "", fmt.Errorf("mixing requires at least two values and exactly one zero")
	}

	previous, next := make([]int, len(values)), make([]int, len(values))
	for i := range values {
		previous[i] = (i + len(values) - 1) % len(values)
		next[i] = (i + 1) % len(values)
	}
	size := int64(len(values) - 1)
	for range rounds {
		for id, value := range values {
			steps := int(xmath.Modulo(value, size))
			if steps > int(size)/2 {
				steps -= int(size)
			}
			from := previous[id]
			next[previous[id]], previous[next[id]] = next[id], previous[id]
			if steps >= 0 {
				for range steps {
					from = next[from]
				}
			} else {
				for range -steps {
					from = previous[from]
				}
			}
			previous[id], next[id] = from, next[from]
			next[from], previous[next[id]] = id, id
		}
	}
	var total int64
	current := zero
	for range 3 {
		for range 1000 % len(values) {
			current = next[current]
		}
		sum, ok := xmath.CheckedAdd(total, values[current])
		if !ok {
			return "", fmt.Errorf("grove coordinate sum overflows int64")
		}
		total = sum
	}
	return strconv.FormatInt(total, 10), nil
}

func day20p01(r io.Reader) (string, error) {
	return groveCoordinates(r, 1, 1)
}

func day20p02(r io.Reader) (string, error) {
	return groveCoordinates(r, 811589153, 10)
}
