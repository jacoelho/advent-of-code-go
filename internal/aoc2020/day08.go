package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type bootInstruction struct {
	operation string
	argument  int
}

func parseBootCode(r io.Reader) ([]bootInstruction, error) {
	s := scanner.NewScanner(r, func(line []byte) (bootInstruction, error) {
		operation, argument, ok := strings.Cut(string(line), " ")
		value, err := strconv.Atoi(argument)
		if !ok || err != nil {
			return bootInstruction{}, fmt.Errorf("invalid boot instruction %q", line)
		}
		switch operation {
		case "nop", "acc", "jmp":
		default:
			return bootInstruction{}, fmt.Errorf("unknown boot operation %q", operation)
		}
		return bootInstruction{operation, value}, nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func executeBoot(code []bootInstruction, swap int) (accumulator int, terminated bool, err error) {
	seen := make([]bool, len(code))
	for ip := 0; ; {
		if ip == len(code) {
			return accumulator, true, nil
		}
		if ip < 0 || ip > len(code) {
			return 0, false, fmt.Errorf("instruction pointer %d outside boot code", ip)
		}
		if seen[ip] {
			return accumulator, false, nil
		}
		seen[ip] = true
		instruction := code[ip]
		if ip == swap {
			switch instruction.operation {
			case "nop":
				instruction.operation = "jmp"
			case "jmp":
				instruction.operation = "nop"
			}
		}
		switch instruction.operation {
		case "acc":
			accumulator += instruction.argument
			ip++
		case "nop":
			ip++
		case "jmp":
			ip += instruction.argument
		}
	}
}

func day08p01(r io.Reader) (string, error) {
	code, err := parseBootCode(r)
	if err != nil {
		return "", err
	}
	accumulator, _, err := executeBoot(code, -1)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(accumulator), nil
}

func day08p02(r io.Reader) (string, error) {
	code, err := parseBootCode(r)
	if err != nil {
		return "", err
	}
	for i, instruction := range code {
		if instruction.operation == "acc" {
			continue
		}
		accumulator, terminated, err := executeBoot(code, i)
		if err == nil && terminated {
			return strconv.Itoa(accumulator), nil
		}
	}
	return "", fmt.Errorf("no single nop/jmp swap repairs boot code")
}
