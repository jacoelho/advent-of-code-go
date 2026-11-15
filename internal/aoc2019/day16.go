package aoc2019

import (
	"bytes"
	"fmt"
	"io"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

const (
	fftPhases       = 100
	fftMessageLen   = 8
	fftOffsetDigits = 7
	fftRepeats      = 10_000
)

func parseSignal(r io.Reader) ([]int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimSpace(data)
	if len(data) < fftMessageLen {
		return nil, fmt.Errorf("signal has %d digits, need at least %d", len(data), fftMessageLen)
	}
	digits := make([]int, len(data))
	for i, c := range data {
		if c < '0' || c > '9' {
			return nil, fmt.Errorf("invalid signal digit %q at %d", c, i)
		}
		digits[i] = int(c - '0')
	}
	return digits, nil
}

func digitsToString[T int | uint8](digits []T) string {
	text := make([]byte, len(digits))
	for i, d := range digits {
		text[i] = byte(d) + '0'
	}
	return string(text)
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
	signal, err := parseSignal(r)
	if err != nil {
		return "", err
	}

	for range fftPhases {
		signal = applyFFTPhase(signal)
	}

	return digitsToString(signal[:fftMessageLen]), nil
}

func day16p02(r io.Reader) (string, error) {
	signal, err := parseSignal(r)
	if err != nil {
		return "", err
	}

	offset := 0
	for _, digit := range signal[:fftOffsetDigits] {
		offset = offset*10 + digit
	}
	totalLen := len(signal) * fftRepeats
	// Past the midpoint every pattern coefficient is one, so each output digit
	// is the suffix sum of the input from its position.
	if offset < totalLen/2 || offset > totalLen-fftMessageLen {
		return "", fmt.Errorf("message offset %d is outside the second half of a %d-digit signal", offset, totalLen)
	}

	tail := make([]uint8, totalLen-offset)
	for i := range tail {
		tail[i] = uint8(signal[(offset+i)%len(signal)])
	}
	for range fftPhases {
		var sum uint8
		for i := len(tail) - 1; i >= 0; i-- {
			sum += tail[i]
			if sum >= 10 {
				sum -= 10
			}
			tail[i] = sum
		}
	}

	return digitsToString(tail[:fftMessageLen]), nil
}
