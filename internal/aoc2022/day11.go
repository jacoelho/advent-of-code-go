package aoc2022

import (
	"fmt"
	"io"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type monkey struct {
	id          int
	items       []*big.Int
	operator    byte
	operand     *big.Int
	divisor     *big.Int
	targets     [2]int
	inspections uint64
}

func parseMonkeys(r io.Reader) ([]monkey, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return nil, err
	}
	blocks := strings.Split(input, "\n\n")
	if len(blocks) < 2 || len(blocks) > 128 {
		return nil, fmt.Errorf("expected 2 to 128 monkeys")
	}
	monkeys := make([]monkey, 0, len(blocks))
	itemCount := 0
	number := func(text string) (int64, error) {
		value, err := convert.ScanNumber[int64]([]byte(text))
		if err != nil {
			return 0, err
		}
		if value < 0 {
			return 0, fmt.Errorf("negative monkey value")
		}
		return value, nil
	}
	for _, block := range blocks {
		lines := strings.Split(block, "\n")
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
		if id >= int64(len(blocks)) {
			return nil, fmt.Errorf("monkey ID out of range")
		}
		m := monkey{id: int(id)}
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
				itemCount++
				if itemCount > 10_000 {
					return nil, fmt.Errorf("too many monkey items")
				}
				m.items = append(m.items, big.NewInt(value))
			}
		}
		operation, ok := strings.CutPrefix(lines[2], "Operation: new = old ")
		fields := strings.Fields(operation)
		if !ok || len(fields) != 2 || (fields[0] != "+" && fields[0] != "*") {
			return nil, fmt.Errorf("invalid monkey operation")
		}
		m.operator = fields[0][0]
		if fields[1] != "old" {
			value, err := number(fields[1])
			if err != nil {
				return nil, err
			}
			m.operand = big.NewInt(value)
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
				m.divisor = big.NewInt(value)
			} else {
				if value >= int64(len(blocks)) || value == id {
					return nil, fmt.Errorf("unknown or unsupported self-target monkey")
				}
				m.targets[i-1] = int(value)
			}
		}
		monkeys = append(monkeys, m)
	}
	slices.SortFunc(monkeys, func(a, b monkey) int { return a.id - b.id })
	for i, m := range monkeys {
		if m.id != i {
			return nil, fmt.Errorf("duplicate or missing monkey ID")
		}
	}
	return monkeys, nil
}

func monkeyBusiness(r io.Reader, rounds int, relief bool) (string, error) {
	monkeys, err := parseMonkeys(r)
	if err != nil {
		return "", err
	}
	modulus := big.NewInt(1)
	var gcd big.Int
	for _, m := range monkeys {
		gcd.GCD(nil, nil, modulus, m.divisor)
		modulus.Quo(modulus, &gcd).Mul(modulus, m.divisor)
	}
	var remainder big.Int
	three := big.NewInt(3)
	work := 0
	for range rounds {
		for i := range monkeys {
			m := &monkeys[i]
			items := m.items
			m.items = nil
			m.inspections += uint64(len(items))
			for _, item := range items {
				work++
				if work > 20_000_000 {
					return "", fmt.Errorf("monkey inspection limit exceeded")
				}
				operand := m.operand
				if operand == nil {
					operand = item
				}
				if m.operator == '+' {
					item.Add(item, operand)
				} else {
					item.Mul(item, operand)
				}
				if relief {
					if item.BitLen() > 4096 {
						return "", fmt.Errorf("worry value exceeds 4096 bits")
					}
					item.Quo(item, three)
				} else {
					item.Mod(item, modulus)
				}
				target := m.targets[1]
				if remainder.Mod(item, m.divisor).Sign() == 0 {
					target = m.targets[0]
				}
				monkeys[target].items = append(monkeys[target].items, item)
			}
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
