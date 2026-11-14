// Package convert converts between integers, their digits, and their decimal
// text.
package convert

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

// ToDigits returns the decimal digits of v, most significant first. It returns
// [0] for zero and nil for negative values.
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

// FromDigits returns the number whose decimal digits, most significant first,
// are s. Elements are not checked to be in [0, 9].
func FromDigits[Slice ~[]T, T constraints.Signed](s Slice) T {
	var sum T
	for _, v := range s {
		sum = sum*10 + v
	}
	return sum
}

// ExtractDigits returns the decimal integers in line, in order of appearance.
// A '-' makes a number negative only when it directly precedes the number and
// does not follow a digit, so "5 -3" yields [5 -3] but "5-3" yields [5 3].
// Values that overflow T wrap.
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

func scanSmallDecimal(line []byte) (int64, bool) {
	// Eighteen unsigned decimal digits always fit in int64.
	if len(line) == 0 || len(line) > 18 {
		return 0, false
	}
	var n int64
	for _, b := range line {
		digit := b - '0'
		if digit > 9 {
			return 0, false
		}
		n = n*10 + int64(digit)
	}
	return n, true
}

// ScanNumber parses line as a base-10 integer of type T. Invalid syntax or a
// value outside T's range returns a *strconv.NumError.
func ScanNumber[T constraints.Signed](line []byte) (T, error) {
	n, ok := scanSmallDecimal(line)
	if !ok {
		var err error
		n, err = strconv.ParseInt(string(line), 10, 64)
		if err != nil {
			return 0, err
		}
	}
	value := T(n)
	if int64(value) != n {
		return 0, &strconv.NumError{Func: "ScanNumber", Num: string(line), Err: strconv.ErrRange}
	}
	return value, nil
}

// FromBinaryToBase10 returns the number whose binary digits, most significant
// first, are digits. Elements are not checked to be 0 or 1.
func FromBinaryToBase10(digits []int) int {
	var result int
	for _, digit := range digits {
		result = (result << 1) | digit
	}
	return result
}
