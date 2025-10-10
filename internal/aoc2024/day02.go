package aoc2024

import (
	"bytes"
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
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

func isReportSafe(report []int) bool {
	levels := slices.Collect(xslices.Window(2, report))

	ascending := xslices.Every(func(v []int) bool {
		diff := v[1] - v[0]
		return diff >= 1 && diff <= 3
	}, levels)

	descending := xslices.Every(func(v []int) bool {
		diff := v[1] - v[0]
		return diff <= -1 && diff >= -3
	}, levels)

	return ascending || descending
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
		for i := range s {
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
