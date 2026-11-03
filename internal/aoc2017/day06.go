package aoc2017

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func memoryBankCycles(r io.Reader) (cycles, loop int, err error) {
	s := scanner.NewScannerWithSplit(r, bufio.ScanWords, convert.ScanNumber[int])
	banks := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return 0, 0, err
	}
	if len(banks) == 0 {
		return 0, 0, fmt.Errorf("no memory banks")
	}

	total := 0
	for _, blocks := range banks {
		if blocks < 0 || blocks > math.MaxInt-total {
			return 0, 0, fmt.Errorf("invalid memory-bank block count")
		}
		total += blocks
	}

	seen := make(map[string]int)
	for cycles := 0; ; cycles++ {
		var state []byte
		for _, blocks := range banks {
			state = strconv.AppendInt(state, int64(blocks), 10)
			state = append(state, ',')
		}
		if previous, ok := seen[string(state)]; ok {
			return cycles, cycles - previous, nil
		}
		seen[string(state)] = cycles
		largest := 0
		for i, blocks := range banks {
			if blocks > banks[largest] {
				largest = i
			}
		}
		blocks := banks[largest]
		banks[largest] = 0
		for i := range banks {
			banks[i] += blocks / len(banks)
		}
		for i := 1; i <= blocks%len(banks); i++ {
			banks[(largest+i)%len(banks)]++
		}
	}
}

func day06p01(r io.Reader) (string, error) {
	cycles, _, err := memoryBankCycles(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(cycles), nil
}

func day06p02(r io.Reader) (string, error) {
	_, loop, err := memoryBankCycles(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(loop), nil
}
