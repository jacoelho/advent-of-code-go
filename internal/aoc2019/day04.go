package aoc2019

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func parsePasswordRange(r io.Reader) (int, int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, 0, err
	}
	lowText, highText, ok := bytes.Cut(bytes.TrimSpace(data), []byte("-"))
	if !ok {
		return 0, 0, fmt.Errorf("password range %q needs the form low-high", data)
	}
	low, err := convert.ScanNumber[int](lowText)
	if err != nil {
		return 0, 0, fmt.Errorf("password range low: %w", err)
	}
	high, err := convert.ScanNumber[int](highText)
	if err != nil {
		return 0, 0, fmt.Errorf("password range high: %w", err)
	}
	if low < 0 || low > high {
		return 0, 0, fmt.Errorf("invalid password range %d-%d", low, high)
	}
	return low, high, nil
}

// countPasswords counts the six-digit numbers in the input range whose digits
// never decrease and whose run lengths satisfy valid. Run lengths are passed
// as a bit set: bit n is set when some run of equal digits has length n.
func countPasswords(r io.Reader, valid func(runLengths uint) bool) (string, error) {
	low, high, err := parsePasswordRange(r)
	if err != nil {
		return "", err
	}

	count := 0
	var walk func(value, length, last, run int, runLengths uint)
	walk = func(value, length, last, run int, runLengths uint) {
		if length == 6 {
			if value >= low && value <= high && valid(runLengths|1<<run) {
				count++
			}
			return
		}
		for digit := max(last, 1); digit <= 9; digit++ {
			if digit == last {
				walk(value*10+digit, length+1, digit, run+1, runLengths)
			} else {
				walk(value*10+digit, length+1, digit, 1, runLengths|1<<run)
			}
		}
	}
	walk(0, 0, 0, 0, 0)
	return strconv.Itoa(count), nil
}

func day04p01(r io.Reader) (string, error) {
	return countPasswords(r, func(runLengths uint) bool { return runLengths>>2 != 0 })
}

func day04p02(r io.Reader) (string, error) {
	return countPasswords(r, func(runLengths uint) bool { return runLengths&(1<<2) != 0 })
}
