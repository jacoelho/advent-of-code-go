package aoc2019

import (
	"io"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
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

func applyFFTPhase(input []int) []int {
	output := make([]int, len(input))
	n := len(input)
	halfPoint := n / 2
	prefix := make([]int, n+1)
	for i, digit := range input {
		prefix[i+1] = prefix[i] + digit
	}

	for i := range halfPoint {
		width := i + 1
		sum, sign := 0, 1
		// Nonzero pattern blocks alternate signs, separated by equally wide zero blocks.
		for start := i; start < n; start += 2 * width {
			sum += sign * (prefix[min(start+width, n)] - prefix[start])
			sign = -sign
		}
		output[i] = xmath.Abs(sum) % 10
	}

	// From the midpoint onward, each output is the suffix sum modulo 10.
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

	// Suffix summation is valid only when the offset is in the second half of the
	// repeated signal.
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
