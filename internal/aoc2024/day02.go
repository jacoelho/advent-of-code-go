package aoc2024

import (
	"bytes"
	"io"
	"iter"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
	"github.com/jacoelho/advent-of-code-go/internal/xmath"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

func parseReports(r io.Reader) (iter.Seq[[]int], error) {
	s := scanner.NewScanner(r, func(b []byte) ([]int, error) {
		numbers := bytes.Fields(b)

		result := make([]int, 0, len(numbers))
		for _, number := range numbers {
			result = append(result, aoc.Must(strconv.Atoi(string(number))))
		}
		return result, nil
	})

	return s.Values(), s.Err()
}

func assertReport(report []int, fn func(a, b int) bool) bool {
	w := xslices.Window(report, 2)
	for v := range w {
		if !fn(v[0], v[1]) {
			return false
		}
	}
	return true
}

func isReportSafe(report []int) bool {
	ascending := assertReport(report, func(a, b int) bool { return a < b })
	descending := assertReport(report, func(a, b int) bool { return a > b })
	differ := assertReport(report, func(a, b int) bool {
		difference := xmath.Abs(a - b)
		return difference >= 1 && difference <= 3
	})

	return differ && (ascending || descending)
}

func day02p01(r io.Reader) (string, error) {
	reports := aoc.Must(parseReports(r))

	safeReportCount := xiter.Reduce(func(sum int, report []int) int {
		if isReportSafe(report) {
			return sum + 1
		}
		return sum
	}, 0, reports)

	return strconv.Itoa(safeReportCount), nil
}

func tolerateOneLevel[T any](s []T) iter.Seq[[]T] {
	if len(s) < 1 {
		panic("need at least one element")
	}
	return func(yield func([]T) bool) {
		for i := 0; i < len(s); i++ {
			if !yield(append(s[:i:i], s[i+1:]...)) {
				return
			}
		}
	}
}

func day02p02(r io.Reader) (string, error) {
	reports := aoc.Must(parseReports(r))

	safeReportCount := xiter.Reduce(func(sum int, report []int) int {
		for variation := range tolerateOneLevel(report) {
			if isReportSafe(variation) {
				return sum + 1
			}
		}
		return sum
	}, 0, reports)

	return strconv.Itoa(safeReportCount), nil
}
