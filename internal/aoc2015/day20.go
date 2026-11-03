package aoc2015

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func firstHouse(r io.Reader, presents, visits int) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	target, err := convert.ScanNumber[int](bytes.TrimSpace(data))
	if err != nil {
		return "", fmt.Errorf("present target: %w", err)
	}
	if target <= 0 {
		return "", fmt.Errorf("invalid present target %q", data)
	}
	target = (target-1)/presents + 1
	// House target receives enough from its own elf, so the sieve needs no larger bound.
	houses := make([]int, target+1)
	for elf := 1; elf <= target; elf++ {
		last := target
		if visits > 0 && elf <= target/visits {
			last = elf * visits
		}
		for house := elf; house <= last; house += elf {
			houses[house] += elf
		}
	}
	for house := 1; house <= target; house++ {
		if houses[house] >= target {
			return strconv.Itoa(house), nil
		}
	}
	return "", fmt.Errorf("no house received enough presents")
}

func day20p01(r io.Reader) (string, error) {
	return firstHouse(r, 10, 0)
}

func day20p02(r io.Reader) (string, error) {
	return firstHouse(r, 11, 50)
}
