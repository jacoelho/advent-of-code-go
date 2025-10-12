package convert

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func ToDigits[T constraints.Signed](v T) []T {
	if v == 0 {
		return []T{0}
	}

	var digits []T
	for v > 0 {
		digits = append(digits, v%10)
		v /= 10
	}

	n := len(digits)
	for i := 0; i < n/2; i++ {
		digits[i], digits[n-1-i] = digits[n-1-i], digits[i]
	}

	return digits
}

func FromDigits[Slice ~[]T, T constraints.Signed](s Slice) T {
	var sum T
	for _, v := range s {
		sum = sum*10 + v
	}
	return sum
}

func ExtractDigits[T constraints.Signed](line string) []T {
	var (
		digits   []T
		current  T
		inNumber bool
		sign     T = 1
	)

	for _, ch := range line {
		switch {
		case ch == '-' && !inNumber:
			sign = -1
		case '0' <= ch && ch <= '9':
			current = current*10 + T(ch-'0')
			inNumber = true
		case inNumber:
			digits = append(digits, current*sign)
			current = 0
			inNumber = false
			sign = 1
		default:
			sign = 1
		}
	}

	if inNumber {
		digits = append(digits, current*sign)
	}

	return digits
}

func ScanNumber[T constraints.Signed](line []byte) (T, error) {
	if len(line) == 0 {
		return 0, fmt.Errorf("empty input")
	}

	var n T
	negative := false
	start := 0

	if line[0] == '-' {
		negative = true
		start = 1
	}

	for i := start; i < len(line); i++ {
		ch := line[i] - '0'
		if ch > 9 {
			return n, fmt.Errorf("invalid character '%c'", line[i])
		}
		n = n*10 + T(ch)
	}

	if negative {
		n = -n
	}
	return n, nil
}

func FromBinaryToBase10(digits []int) int {
	var result int
	for _, digit := range digits {
		result = (result << 1) | digit
	}
	return result
}
