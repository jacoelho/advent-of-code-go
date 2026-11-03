package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type assemblyOperand struct {
	register int
	literal  int64
}

func (operand assemblyOperand) value(registers [26]int64) int64 {
	if operand.register >= 0 {
		return registers[operand.register]
	}

	return operand.literal
}

type assemblyInstruction struct {
	operation string
	x, y      assemblyOperand
}

func parseAssemblyOperand(text []byte) (assemblyOperand, error) {
	if len(text) == 1 && text[0] >= 'a' && text[0] <= 'z' {
		return assemblyOperand{register: int(text[0] - 'a')}, nil
	}

	value, err := convert.ScanNumber[int64](text)
	if err != nil {
		return assemblyOperand{}, fmt.Errorf("invalid assembly operand %q: %w", text, err)
	}

	return assemblyOperand{register: -1, literal: value}, nil
}

func parseAssembly(r io.Reader) ([]assemblyInstruction, error) {
	s := scanner.NewScanner(r, func(line []byte) (assemblyInstruction, error) {
		fields := bytes.Fields(line)
		if len(fields) == 0 {
			return assemblyInstruction{}, fmt.Errorf("empty assembly instruction")
		}

		instruction := assemblyInstruction{
			operation: string(fields[0]),
			x:         assemblyOperand{register: -1},
			y:         assemblyOperand{register: -1},
		}
		arity := 3
		switch instruction.operation {
		case "snd", "rcv":
			arity = 2
		case "set", "add", "sub", "mul", "mod", "jgz", "jnz":
		default:
			return assemblyInstruction{}, fmt.Errorf("unsupported assembly operation %q", instruction.operation)
		}
		if len(fields) != arity {
			return assemblyInstruction{}, fmt.Errorf("invalid arity for %q", line)
		}

		var err error
		instruction.x, err = parseAssemblyOperand(fields[1])
		if err != nil {
			return assemblyInstruction{}, err
		}
		if arity == 3 {
			instruction.y, err = parseAssemblyOperand(fields[2])
			if err != nil {
				return assemblyInstruction{}, err
			}
		}

		switch instruction.operation {
		case "set", "add", "sub", "mul", "mod":
			if instruction.x.register < 0 {
				return assemblyInstruction{}, fmt.Errorf("%s destination must be a register", instruction.operation)
			}
		}

		return instruction, nil
	})

	var program []assemblyInstruction
	for instruction := range s.Values() {
		program = append(program, instruction)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(program) == 0 {
		return nil, fmt.Errorf("empty assembly program")
	}

	return program, nil
}

func assemblySum(a, b int64) (int64, error) {
	if b > 0 && a > math.MaxInt64-b || b < 0 && a < math.MinInt64-b {
		return 0, fmt.Errorf("assembly integer overflow")
	}

	return a + b, nil
}

func advanceAssembly(instruction assemblyInstruction, registers *[26]int64, pointer int64) (int64, error) {
	x, y := instruction.x.value(*registers), instruction.y.value(*registers)
	var value int64
	var err error

	switch instruction.operation {
	case "set":
		value = y

	case "add":
		value, err = assemblySum(x, y)

	case "sub":
		if y > 0 && x < math.MinInt64+y || y < 0 && x > math.MaxInt64+y {
			err = fmt.Errorf("assembly integer overflow")
		} else {
			value = x - y
		}

	case "mul":
		value = x * y
		if x != 0 && (x == -1 && y == math.MinInt64 || value/x != y) {
			err = fmt.Errorf("assembly integer overflow")
		}

	case "mod":
		if y == 0 {
			return 0, fmt.Errorf("assembly modulo by zero")
		}
		value = x % y

	case "jgz", "jnz":
		if instruction.operation == "jgz" && x > 0 || instruction.operation == "jnz" && x != 0 {
			return assemblySum(pointer, y)
		}
		return pointer + 1, nil

	default:
		return 0, fmt.Errorf("unsupported assembly arithmetic operation %q", instruction.operation)
	}
	if err != nil {
		return 0, err
	}

	registers[instruction.x.register] = value
	return pointer + 1, nil
}
