package aoc2022

import (
	"fmt"
	"io"
	"math/bits"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type monkey struct {
	items       []uint64
	operator    byte
	squares     bool
	operand     uint64
	divisor     uint64
	targets     [2]int
	inspections uint64
}

func parseMonkeys(r io.Reader) ([]monkey, error) {
	blocks, err := scanner.ReadBlocks(r)
	if err != nil {
		return nil, err
	}
	if len(blocks) < 2 {
		return nil, fmt.Errorf("expected at least two monkeys")
	}
	monkeys := make([]monkey, len(blocks))
	defined := make([]bool, len(blocks))
	number := func(text string) (uint64, error) {
		value, err := convert.ScanNumber[int64]([]byte(text))
		if err != nil {
			return 0, err
		}
		if value < 0 {
			return 0, fmt.Errorf("negative monkey value %d", value)
		}
		return uint64(value), nil
	}
	for _, lines := range blocks {
		if len(lines) != 6 {
			return nil, fmt.Errorf("monkey requires six records")
		}
		for i := range lines {
			lines[i] = strings.TrimSpace(lines[i])
		}
		header, ok := strings.CutPrefix(lines[0], "Monkey ")
		if !ok || !strings.HasSuffix(header, ":") {
			return nil, fmt.Errorf("invalid monkey header")
		}
		id, err := number(strings.TrimSuffix(header, ":"))
		if err != nil {
			return nil, err
		}
		if id >= uint64(len(blocks)) || defined[id] {
			return nil, fmt.Errorf("monkey ID %d is out of range or duplicated", id)
		}
		defined[id] = true
		m := &monkeys[id]
		items, ok := strings.CutPrefix(lines[1], "Starting items:")
		if !ok {
			return nil, fmt.Errorf("invalid monkey items")
		}
		if items = strings.TrimSpace(items); items != "" {
			for item := range strings.SplitSeq(items, ",") {
				value, err := number(strings.TrimSpace(item))
				if err != nil {
					return nil, err
				}
				m.items = append(m.items, value)
			}
		}
		operation, ok := strings.CutPrefix(lines[2], "Operation: new = old ")
		fields := strings.Fields(operation)
		if !ok || len(fields) != 2 || (fields[0] != "+" && fields[0] != "*") {
			return nil, fmt.Errorf("invalid monkey operation")
		}
		m.operator = fields[0][0]
		if fields[1] == "old" {
			m.squares = true
		} else {
			m.operand, err = number(fields[1])
			if err != nil {
				return nil, err
			}
		}
		for i, prefix := range []string{"Test: divisible by ", "If true: throw to monkey ", "If false: throw to monkey "} {
			text, ok := strings.CutPrefix(lines[i+3], prefix)
			if !ok {
				return nil, fmt.Errorf("invalid monkey test or target")
			}
			value, err := number(text)
			if err != nil {
				return nil, err
			}
			if i == 0 {
				if value == 0 {
					return nil, fmt.Errorf("zero monkey divisor")
				}
				m.divisor = value
			} else {
				if value >= uint64(len(blocks)) || value == id {
					return nil, fmt.Errorf("unknown or unsupported self-target monkey")
				}
				m.targets[i-1] = int(value)
			}
		}
	}
	return monkeys, nil
}

func monkeyBusiness(r io.Reader, rounds int, relief bool) (string, error) {
	monkeys, err := parseMonkeys(r)
	if err != nil {
		return "", err
	}
	// Without relief, reducing modulo the divisors' LCM preserves every
	// divisibility test.
	modulus := uint64(1)
	if !relief {
		for _, m := range monkeys {
			var ok bool
			if modulus, ok = xmath.CheckedMul(modulus/xmath.GCD(modulus, m.divisor), m.divisor); !ok {
				return "", fmt.Errorf("monkey divisor LCM overflows uint64")
			}
		}
	}
	for range rounds {
		for i := range monkeys {
			m := &monkeys[i]
			m.inspections += uint64(len(m.items))
			for _, item := range m.items {
				operand := m.operand
				if m.squares {
					operand = item
				}
				var hi, lo uint64
				if m.operator == '+' {
					lo, hi = bits.Add64(item, operand, 0)
				} else {
					hi, lo = bits.Mul64(item, operand)
				}
				if relief {
					if hi != 0 {
						return "", fmt.Errorf("worry level overflows uint64")
					}
					item = lo / 3
				} else {
					item = bits.Rem64(hi, lo, modulus)
				}
				target := m.targets[1]
				if item%m.divisor == 0 {
					target = m.targets[0]
				}
				monkeys[target].items = append(monkeys[target].items, item)
			}
			m.items = m.items[:0]
		}
	}
	var first, second uint64
	for _, m := range monkeys {
		if m.inspections > first {
			first, second = m.inspections, first
		} else {
			second = max(second, m.inspections)
		}
	}
	return strconv.FormatUint(first*second, 10), nil
}

func day11p01(r io.Reader) (string, error) {
	return monkeyBusiness(r, 20, true)
}

func day11p02(r io.Reader) (string, error) {
	return monkeyBusiness(r, 10_000, false)
}
