package aoc2019

import (
	"io"
	"slices"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func parseInput(r io.Reader) ([]int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	input := string(data)
	digits := make([]int, 0, len(input))
	for _, ch := range input {
		if ch >= '0' && ch <= '9' {
			digits = append(digits, int(ch-'0'))
		}
	}
	return digits, nil
}

func digitsToString(digits []int, n int) string {
	var result strings.Builder
	for i := 0; i < min(n, len(digits)); i++ {
		result.WriteRune(rune(digits[i] + '0'))
	}
	return result.String()
}

func generatePattern(position, length int) []int {
	basePattern := []int{0, 1, 0, -1}
	pattern := make([]int, 0, length+1)

	for len(pattern) <= length {
		for _, val := range basePattern {
			for range position {
				pattern = append(pattern, val)
				if len(pattern) > length {
					break
				}
			}
			if len(pattern) > length {
				break
			}
		}
	}

	return pattern[1 : length+1]
}

func applyFFTPhase(input []int) []int {
	output := make([]int, len(input))
	n := len(input)
	halfPoint := n / 2

	for i := range halfPoint {
		pattern := generatePattern(i+1, n)
		sum := xiter.DotProduct(slices.Values(input), slices.Values(pattern))
		output[i] = xmath.Abs(sum) % 10
	}

	// backward summation
	output[n-1] = input[n-1]
	for i := n - 2; i >= halfPoint; i-- {
		output[i] = (output[i+1] + input[i]) % 10
	}

	return output
}

func day16p01(r io.Reader) (string, error) {
	signal, err := parseInput(r)
	if err != nil {
		return "", err
	}

	for range 100 {
		signal = applyFFTPhase(signal)
	}

	return digitsToString(signal, 8), nil
}

func day16p02(r io.Reader) (string, error) {
	signal, err := parseInput(r)
	if err != nil {
		return "", err
	}

	offset := convert.FromDigits(signal[:min(7, len(signal))])

	totalLen := len(signal) * 10000

	// optimization only works if offset is in second half
	if offset < totalLen/2 {
		return "", nil
	}

	neededLen := totalLen - offset

	current := make([]int, neededLen)
	for i := range neededLen {
		pos := offset + i
		current[i] = signal[pos%len(signal)]
	}

	for range 100 {
		for i := len(current) - 2; i >= 0; i-- {
			current[i] = (current[i+1] + current[i]) % 10
		}
	}

	return digitsToString(current, 8), nil
}
