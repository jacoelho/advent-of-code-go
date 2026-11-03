package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type aluOperand struct {
	register byte
	value    int64
}

type aluInstruction struct {
	operation string
	target    byte
	operand   aluOperand
}

func aluRegister(s string) bool {
	return len(s) == 1 && s[0] >= 'w' && s[0] <= 'z'
}

func parseALU(r io.Reader) ([]aluInstruction, error) {
	var program []aluInstruction
	s := bufio.NewScanner(r)
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) < 2 || !aluRegister(fields[1]) {
			return nil, fmt.Errorf("invalid ALU instruction %q", s.Text())
		}
		instruction := aluInstruction{operation: fields[0], target: fields[1][0]}
		switch instruction.operation {
		case "inp":
			if len(fields) != 2 {
				return nil, fmt.Errorf("inp requires one register")
			}
		case "add", "mul", "div", "mod", "eql":
			if len(fields) != 3 {
				return nil, fmt.Errorf("%s requires two operands", fields[0])
			}
			if aluRegister(fields[2]) {
				instruction.operand.register = fields[2][0]
			} else {
				value, err := convert.ScanNumber[int64]([]byte(fields[2]))
				if err != nil {
					return nil, fmt.Errorf("ALU instruction %d: %w", len(program)+1, err)
				}
				instruction.operand.value = value
			}
		default:
			return nil, fmt.Errorf("unknown ALU operation %q", fields[0])
		}
		program = append(program, instruction)
		if len(program) > 18*18 {
			return nil, fmt.Errorf("ALU program exceeds 324 instructions")
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(program) == 0 {
		return nil, fmt.Errorf("empty ALU program")
	}
	return program, nil
}

func executeALU(program []aluInstruction, input []int64) ([4]int64, error) {
	var registers [4]int64
	position := 0
	for i, instruction := range program {
		target := &registers[instruction.target-'w']
		if instruction.operation == "inp" {
			if position == len(input) {
				return registers, fmt.Errorf("ALU instruction %d has no input", i+1)
			}
			*target = input[position]
			position++
			continue
		}
		a, b := *target, instruction.operand.value
		if instruction.operand.register != 0 {
			b = registers[instruction.operand.register-'w']
		}
		var value int64
		switch instruction.operation {
		case "add":
			if b > 0 && a > math.MaxInt64-b || b < 0 && a < math.MinInt64-b {
				return registers, fmt.Errorf("ALU instruction %d addition overflows int64", i+1)
			}
			value = a + b
		case "mul":
			value = a * b
			if a == -1 && b == math.MinInt64 || b == -1 && a == math.MinInt64 || a != 0 && value/a != b {
				return registers, fmt.Errorf("ALU instruction %d multiplication overflows int64", i+1)
			}
		case "div":
			if b == 0 || a == math.MinInt64 && b == -1 {
				return registers, fmt.Errorf("ALU instruction %d has invalid division", i+1)
			}
			value = a / b
		case "mod":
			if a < 0 || b <= 0 {
				return registers, fmt.Errorf("ALU instruction %d has invalid modulo", i+1)
			}
			value = a % b
		case "eql":
			if a == b {
				value = 1
			}
		}
		*target = value
	}
	if position != len(input) {
		return registers, fmt.Errorf("ALU has %d unused inputs", len(input)-position)
	}
	return registers, nil
}

type monadBlock struct {
	divisor, offset, addend int64
}

func parseMONAD(program []aluInstruction) ([]monadBlock, error) {
	if len(program)%18 != 0 || len(program) < 36 {
		return nil, fmt.Errorf("unsupported MONAD program: expected complete 18-instruction blocks")
	}
	literal := func(op string, target byte, value int64) aluInstruction {
		return aluInstruction{operation: op, target: target, operand: aluOperand{value: value}}
	}
	register := func(op string, target, source byte) aluInstruction {
		return aluInstruction{operation: op, target: target, operand: aluOperand{register: source}}
	}
	blocks := make([]monadBlock, 0, len(program)/18)
	for start := 0; start < len(program); start += 18 {
		block := monadBlock{
			divisor: program[start+4].operand.value,
			offset:  program[start+5].operand.value,
			addend:  program[start+15].operand.value,
		}
		template := [...]aluInstruction{
			{operation: "inp", target: 'w'},
			literal("mul", 'x', 0),
			register("add", 'x', 'z'),
			literal("mod", 'x', 26),
			literal("div", 'z', block.divisor),
			literal("add", 'x', block.offset),
			register("eql", 'x', 'w'),
			literal("eql", 'x', 0),
			literal("mul", 'y', 0),
			literal("add", 'y', 25),
			register("mul", 'y', 'x'),
			literal("add", 'y', 1),
			register("mul", 'z', 'y'),
			literal("mul", 'y', 0),
			register("add", 'y', 'w'),
			literal("add", 'y', block.addend),
			register("mul", 'y', 'x'),
			register("add", 'z', 'y'),
		}
		for i, expected := range template {
			if program[start+i] != expected {
				return nil, fmt.Errorf("unsupported MONAD instruction %d", start+i+1)
			}
		}
		if block.addend < 0 || block.addend > 16 ||
			block.divisor == 1 && (block.offset < 10 || block.offset > 25) ||
			block.divisor == 26 && (block.offset < -25 || block.offset > 9) ||
			block.divisor != 1 && block.divisor != 26 {
			return nil, fmt.Errorf("unsupported MONAD constants in block %d", len(blocks)+1)
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func monadModel(r io.Reader, largest bool) (string, error) {
	program, err := parseALU(r)
	if err != nil {
		return "", err
	}
	blocks, err := parseMONAD(program)
	if err != nil {
		return "", err
	}
	digits := make([]int64, len(blocks))
	var stack []int
	// Nonzero base-26 digits make every pop match necessary for a balanced program to end at zero.
	for i, block := range blocks {
		if block.divisor == 1 {
			stack = append(stack, i)
			continue
		}
		if len(stack) == 0 {
			return "", fmt.Errorf("MONAD block %d pops an empty stack", i+1)
		}
		push := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		difference := blocks[push].addend + block.offset
		low, high := max(int64(1), 1-difference), min(int64(9), 9-difference)
		if low > high {
			return "", fmt.Errorf("MONAD blocks %d and %d have no valid digit pair", push+1, i+1)
		}
		digits[push] = low
		if largest {
			digits[push] = high
		}
		digits[i] = digits[push] + difference
	}
	if len(stack) != 0 {
		return "", fmt.Errorf("MONAD leaves unmatched push blocks")
	}
	registers, err := executeALU(program, digits)
	if err != nil {
		return "", err
	}
	if registers[3] != 0 {
		return "", fmt.Errorf("derived MONAD model does not finish with z=0")
	}
	model := make([]byte, len(digits))
	for i, digit := range digits {
		model[i] = byte(digit) + '0'
	}
	return string(model), nil
}

func day24p01(r io.Reader) (string, error) {
	return monadModel(r, true)
}

func day24p02(r io.Reader) (string, error) {
	return monadModel(r, false)
}
