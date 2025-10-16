package aoc2023

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

var calibrationDigits = map[string]int{
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
}

var calibrationDigitsAndWords = map[string]int{
	"1": 1, "one": 1,
	"2": 2, "two": 2,
	"3": 3, "three": 3,
	"4": 4, "four": 4,
	"5": 5, "five": 5,
	"6": 6, "six": 6,
	"7": 7, "seven": 7,
	"8": 8, "eight": 8,
	"9": 9, "nine": 9,
}

func calibration(table map[string]int, b []byte) (int, bool) {
	for k, v := range table {
		if bytes.HasPrefix(b, []byte(k)) {
			return v, true
		}
	}
	return 0, false
}

func day01(table map[string]int, r io.Reader) (string, error) {
	s := scanner.NewScanner(r, func(bytes []byte) (int, error) {
		var digits []int
		for i := range bytes {
			if v, ok := calibration(table, bytes[i:]); ok {
				digits = append(digits, v)
			}
		}

		if len(digits) == 0 {
			return 0, fmt.Errorf("expected digits")
		}

		return digits[0]*10 + digits[len(digits)-1], nil
	})

	result := xiter.Sum(s.Values())

	return strconv.Itoa(result), s.Err()
}

func day01p01(r io.Reader) (string, error) {
	return day01(calibrationDigits, r)
}

func day01p02(r io.Reader) (string, error) {
	return day01(calibrationDigitsAndWords, r)
}
