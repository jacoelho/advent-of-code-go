package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
	"strconv"
	"strings"
)

type displayEntry struct {
	signals [10]uint8
	output  [4]uint8
}

func displayEntries(r io.Reader) ([]displayEntry, error) {
	s := bufio.NewScanner(r)
	var entries []displayEntry
	for s.Scan() {
		left, right, ok := strings.Cut(s.Text(), "|")
		if !ok {
			return nil, fmt.Errorf("missing display delimiter")
		}
		signals, output := strings.Fields(left), strings.Fields(right)
		if len(signals) != 10 || len(output) != 4 {
			return nil, fmt.Errorf("display requires ten signals and four outputs")
		}
		var entry displayEntry
		var seen [128]bool
		for i, word := range append(signals, output...) {
			var mask uint8
			for _, segment := range word {
				if segment < 'a' || segment > 'g' || mask&(1<<(segment-'a')) != 0 {
					return nil, fmt.Errorf("invalid display segments %q", word)
				}
				mask |= 1 << (segment - 'a')
			}
			if i < 10 {
				if seen[mask] {
					return nil, fmt.Errorf("duplicate display signal %q", word)
				}
				seen[mask] = true
				entry.signals[i] = mask
			} else {
				entry.output[i-10] = mask
			}
		}
		entries = append(entries, entry)
	}
	return entries, s.Err()
}

func decodedDisplay(entry displayEntry) (int, error) {
	var digits [10]uint8
	for _, signal := range entry.signals {
		switch bits.OnesCount8(signal) {
		case 2:
			digits[1] = signal
		case 3:
			digits[7] = signal
		case 4:
			digits[4] = signal
		case 7:
			digits[8] = signal
		}
	}
	if digits[1] == 0 || digits[4] == 0 || digits[7] == 0 || digits[8] == 0 {
		return 0, fmt.Errorf("missing unique-length display digits")
	}
	for _, signal := range entry.signals {
		if bits.OnesCount8(signal) == 6 {
			switch {
			case signal&digits[4] == digits[4]:
				digits[9] = signal
			case signal&digits[1] == digits[1]:
				digits[0] = signal
			default:
				digits[6] = signal
			}
		}
	}
	for _, signal := range entry.signals {
		if bits.OnesCount8(signal) == 5 {
			switch {
			case signal&digits[1] == digits[1]:
				digits[3] = signal
			case signal&digits[6] == signal:
				digits[5] = signal
			default:
				digits[2] = signal
			}
		}
	}
	decoded := make(map[uint8]int, 10)
	for digit, signal := range digits {
		if signal == 0 {
			return 0, fmt.Errorf("inconsistent display signals")
		}
		decoded[signal] = digit
	}
	value := 0
	for _, signal := range entry.output {
		digit, ok := decoded[signal]
		if !ok {
			return 0, fmt.Errorf("unknown output signal")
		}
		value = value*10 + digit
	}
	return value, nil
}

func day08p01(r io.Reader) (string, error) {
	entries, err := displayEntries(r)
	if err != nil {
		return "", err
	}
	count := 0
	for _, entry := range entries {
		for _, signal := range entry.output {
			switch bits.OnesCount8(signal) {
			case 2, 3, 4, 7:
				count++
			}
		}
	}
	return strconv.Itoa(count), nil
}

func day08p02(r io.Reader) (string, error) {
	entries, err := displayEntries(r)
	if err != nil {
		return "", err
	}
	total := 0
	for _, entry := range entries {
		value, err := decodedDisplay(entry)
		if err != nil {
			return "", err
		}
		total += value
	}
	return strconv.Itoa(total), nil
}
