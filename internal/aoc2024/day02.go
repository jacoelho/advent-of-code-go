package aoc2024

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseReport(b []byte) ([]int, error) {
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
}

func isReportSafe(report []int) bool {
	if len(report) < 2 {
		return true
	}
	levels := xiter.From(xslices.Window(2, report))

	ascending := levels.Every(func(v []int) bool {
		diff := v[1] - v[0]
		return diff >= 1 && diff <= 3
	})

	descending := levels.Every(func(v []int) bool {
		diff := v[1] - v[0]
		return diff <= -1 && diff >= -3
	})

	return ascending || descending
}

func day02p01(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseReport)
	safeReportCount := xiter.From(s.Values()).Filter(isReportSafe).Count()
	if err := s.Err(); err != nil {
		return "", err
	}
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
	s := scanner.NewScanner(r, parseReport)
	safeReportCount := xiter.From(s.Values()).Filter(func(report []int) bool {
		return xiter.From(tolerateOneLevel(report)).Any(isReportSafe)
	}).Count()
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.Itoa(safeReportCount), nil
}
