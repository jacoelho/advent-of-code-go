package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type registerInstruction struct {
	target, operation, conditionRegister, comparison string
	amount, conditionValue                           int
}

func registerMaxima(r io.Reader) (final, highest int, err error) {
	s := scanner.NewScanner(r, func(line []byte) (registerInstruction, error) {
		fields := bytes.Fields(line)
		if len(fields) != 7 || string(fields[3]) != "if" {
			return registerInstruction{}, fmt.Errorf("invalid register instruction %q", line)
		}
		operation, comparison := string(fields[1]), string(fields[5])
		if operation != "inc" && operation != "dec" {
			return registerInstruction{}, fmt.Errorf("unknown register operation %q", operation)
		}
		switch comparison {
		case ">", ">=", "<", "<=", "==", "!=":
		default:
			return registerInstruction{}, fmt.Errorf("invalid register comparison %q", comparison)
		}
		amount, err := convert.ScanNumber[int](fields[2])
		if err != nil {
			return registerInstruction{}, err
		}
		condition, err := convert.ScanNumber[int](fields[6])
		if err != nil {
			return registerInstruction{}, err
		}
		return registerInstruction{string(fields[0]), operation, string(fields[4]), comparison, amount, condition}, nil
	})
	registers := make(map[string]int)
	for instruction := range s.Values() {
		if _, exists := registers[instruction.target]; !exists {
			registers[instruction.target] = 0
		}
		if _, exists := registers[instruction.conditionRegister]; !exists {
			registers[instruction.conditionRegister] = 0
		}
		value := registers[instruction.conditionRegister]
		allowed := false
		switch instruction.comparison {
		case ">":
			allowed = value > instruction.conditionValue
		case ">=":
			allowed = value >= instruction.conditionValue
		case "<":
			allowed = value < instruction.conditionValue
		case "<=":
			allowed = value <= instruction.conditionValue
		case "==":
			allowed = value == instruction.conditionValue
		case "!=":
			allowed = value != instruction.conditionValue
		}
		if allowed {
			if instruction.operation == "inc" {
				registers[instruction.target] += instruction.amount
			} else {
				registers[instruction.target] -= instruction.amount
			}
			highest = max(highest, registers[instruction.target])
		}
	}
	if err := s.Err(); err != nil {
		return 0, 0, err
	}

	first := true
	for _, value := range registers {
		if first || value > final {
			final, first = value, false
		}
	}

	return final, highest, nil
}

func day08p01(r io.Reader) (string, error) {
	final, _, err := registerMaxima(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(final), nil
}

func day08p02(r io.Reader) (string, error) {
	_, highest, err := registerMaxima(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(highest), nil
}
