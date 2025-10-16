package aoc2023

import (
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parsePattern(data []byte) ([]string, error) {
	lines := strings.Split(string(data), "\n")
	return xslices.Filter(func(line string) bool { return len(line) > 0 }, lines), nil
}

func countDifferences(a, b string) int {
	count := 0
	for i := range a {
		if a[i] != b[i] {
			count++
		}
	}
	return count
}

func findVerticalReflection(pattern []string, allowedSmudges int) int {
	width := len(pattern[0])

	for col := 1; col < width; col++ {
		totalDifferences := 0

		for row := range pattern {
			left := col - 1
			right := col

			for left >= 0 && right < width {
				if pattern[row][left] != pattern[row][right] {
					totalDifferences++
				}
				left--
				right++
			}
		}

		if totalDifferences == allowedSmudges {
			return col
		}
	}

	return 0
}

func findHorizontalReflection(pattern []string, allowedSmudges int) int {
	for row := 1; row < len(pattern); row++ {
		totalDifferences := 0

		top := row - 1
		bottom := row

		for top >= 0 && bottom < len(pattern) {
			totalDifferences += countDifferences(pattern[top], pattern[bottom])
			top--
			bottom++
		}

		if totalDifferences == allowedSmudges {
			return row
		}
	}

	return 0
}

func day13(r io.Reader, allowedSmudges int) (string, error) {
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte("\n\n")), parsePattern)

	result := xiter.Sum(xiter.Map(func(pattern []string) int {
		if v := findVerticalReflection(pattern, allowedSmudges); v > 0 {
			return v
		}
		if h := findHorizontalReflection(pattern, allowedSmudges); h > 0 {
			return h * 100
		}
		return 0
	}, s.Values()))

	return strconv.Itoa(result), s.Err()
}

func day13p01(r io.Reader) (string, error) {
	return day13(r, 0)
}

func day13p02(r io.Reader) (string, error) {
	return day13(r, 1)
}
