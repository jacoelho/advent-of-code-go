package aoc2024

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseReports(r io.Reader) ([][]int, error) {
	s := scanner.NewScanner(r, func(b []byte) ([]int, error) {
		numbers := bytes.Fields(b)
		if len(numbers) == 0 {
			return nil, fmt.Errorf("report has no levels")
		}

		result := make([]int, 0, len(numbers))
		for _, number := range numbers {
			value, err := convert.ScanNumber[int](number)
			if err != nil {
				return nil, fmt.Errorf("report level: %w", err)
			}
			result = append(result, value)
		}
		return result, nil
	})

	reports := slices.Collect(s.Values())
	return reports, s.Err()
}

func isReportSafe(report []int) bool {
	if len(report) < 2 {
		return true
	}
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
	reports, err := parseReports(r)
	if err != nil {
		return "", err
	}

	safeReportCount := xiter.Reduce(func(sum int, report []int) int {
		if isReportSafe(report) {
			return sum + 1
		}
		return sum
	}, 0, slices.Values(reports))

	return strconv.Itoa(safeReportCount), nil
}

func tolerateOneLevel[T any](s []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		for i := range s {
			if !yield(append(s[:i:i], s[i+1:]...)) {
				return
			}
		}
	}
}

func day02p02(r io.Reader) (string, error) {
	reports, err := parseReports(r)
	if err != nil {
		return "", err
	}

	safeReportCount := xiter.Reduce(func(sum int, report []int) int {
		for variation := range tolerateOneLevel(report) {
			if isReportSafe(variation) {
				return sum + 1
			}
		}
		return sum
	}, 0, slices.Values(reports))

	return strconv.Itoa(safeReportCount), nil
}
