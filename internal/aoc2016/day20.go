package aoc2016

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func allowedAddresses(r io.Reader, maximum int64) (int64, int64, error) {
	type interval struct{ low, high int64 }
	var blocked []interval
	s := bufio.NewScanner(r)
	for s.Scan() {
		first, last, ok := strings.Cut(s.Text(), "-")
		if !ok {
			return 0, 0, fmt.Errorf("invalid blocked interval %q", s.Text())
		}
		low, err := convert.ScanNumber[int64]([]byte(first))
		if err != nil {
			return 0, 0, fmt.Errorf("lower address: %w", err)
		}
		high, err := convert.ScanNumber[int64]([]byte(last))
		if err != nil {
			return 0, 0, fmt.Errorf("upper address: %w", err)
		}
		if low < 0 || high < low || high > maximum {
			return 0, 0, fmt.Errorf("invalid address interval %d-%d", low, high)
		}
		blocked = append(blocked, interval{low, high})
	}
	if err := s.Err(); err != nil {
		return 0, 0, err
	}
	slices.SortFunc(blocked, func(a, b interval) int { return cmp.Compare(a.low, b.low) })
	first, count, current := int64(-1), int64(0), int64(0)
	for _, span := range blocked {
		if current < span.low {
			if first < 0 {
				first = current
			}
			count += span.low - current
		}
		current = max(current, span.high+1)
	}
	if current <= maximum {
		if first < 0 {
			first = current
		}
		count += maximum - current + 1
	}
	return first, count, nil
}

func day20p01(r io.Reader) (string, error) {
	first, _, err := allowedAddresses(r, math.MaxUint32)
	if err != nil {
		return "", err
	}
	if first < 0 {
		return "", fmt.Errorf("all addresses are blocked")
	}
	return strconv.FormatInt(first, 10), nil
}

func day20p02(r io.Reader) (string, error) {
	_, count, err := allowedAddresses(r, math.MaxUint32)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(count, 10), nil
}
