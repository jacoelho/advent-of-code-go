package aoc2022

import (
	"fmt"
	"io"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

const snafuDigits = "=-012"

// day25p01 adds the numbers column by column in balanced base five, so the
// sum is exact whatever its magnitude.
func day25p01(r io.Reader) (string, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return "", err
	}
	var columns []int
	for line := range strings.SplitSeq(input, "\n") {
		if line == "" {
			return "", fmt.Errorf("empty SNAFU number")
		}
		for len(columns) < len(line) {
			columns = append(columns, 0)
		}
		for i := range len(line) {
			digit := strings.IndexByte(snafuDigits, line[len(line)-1-i])
			if digit < 0 {
				return "", fmt.Errorf("invalid SNAFU digit %q in %q", line[len(line)-1-i], line)
			}
			columns[i] += digit - 2
		}
	}
	var digits []byte
	carry := 0
	for i := 0; i < len(columns) || carry != 0; i++ {
		value := carry
		if i < len(columns) {
			value += columns[i]
		}
		digit := xmath.Modulo(value+2, 5) - 2
		carry = (value - digit) / 5
		digits = append(digits, snafuDigits[digit+2])
	}
	for len(digits) > 1 && digits[len(digits)-1] == '0' {
		digits = digits[:len(digits)-1]
	}
	var result strings.Builder
	for i := len(digits) - 1; i >= 0; i-- {
		result.WriteByte(digits[i])
	}
	return result.String(), nil
}
