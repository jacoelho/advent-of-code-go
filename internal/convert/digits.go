package convert

import (
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

func ExtractDigits[T constraints.Integer](line string) []T {
	var (
		digits   []T
		current  T
		inNumber bool
	)

	for _, ch := range line {
		if '0' <= ch && ch <= '9' {
			current = current*10 + T(ch-'0')
			inNumber = true
		} else if inNumber {
			digits = append(digits, current)
			current = 0
			inNumber = false
		}
	}

	if inNumber {
		digits = append(digits, current)
	}

	return digits
}
