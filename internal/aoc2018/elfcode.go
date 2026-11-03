package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

const elfOperationNames = "addr addi mulr muli banr bani borr bori setr seti gtir gtri gtrr eqir eqri eqrr"

const elfStepLimit = 10_000_000

type elfInstruction struct {
	operation string
	a, b, c   int
}

func (instruction elfInstruction) validate(registerCount int) error {
	if instruction.c < 0 || instruction.c >= registerCount {
		return fmt.Errorf("invalid destination register %d", instruction.c)
	}
	var registerA, registerB bool
	switch instruction.operation {
	case "addr", "mulr", "banr", "borr", "gtrr", "eqrr":
		registerA, registerB = true, true
	case "addi", "muli", "bani", "bori", "setr", "gtri", "eqri":
		registerA = true
	case "gtir", "eqir":
		registerB = true
	case "seti":
	default:
		return fmt.Errorf("unknown operation %q", instruction.operation)
	}
	if registerA && (instruction.a < 0 || instruction.a >= registerCount) {
		return fmt.Errorf("invalid source register %d", instruction.a)
	}
	if registerB && (instruction.b < 0 || instruction.b >= registerCount) {
		return fmt.Errorf("invalid source register %d", instruction.b)
	}
	return nil
}

func elfAdd(a, b int64) (int64, error) {
	if b > 0 && a > math.MaxInt64-b || b < 0 && a < math.MinInt64-b {
		return 0, fmt.Errorf("register addition overflows int64")
	}
	return a + b, nil
}

func elfMultiply(a, b int64) (int64, error) {
	product := a * b
	if a == -1 && b == math.MinInt64 || b == -1 && a == math.MinInt64 || a != 0 && product/a != b {
		return 0, fmt.Errorf("register multiplication overflows int64")
	}
	return product, nil
}

func (instruction elfInstruction) execute(registers []int64) error {
	a, b := instruction.a, instruction.b
	var value int64
	var err error
	switch instruction.operation {
	case "addr":
		value, err = elfAdd(registers[a], registers[b])
	case "addi":
		value, err = elfAdd(registers[a], int64(b))
	case "mulr":
		value, err = elfMultiply(registers[a], registers[b])
	case "muli":
		value, err = elfMultiply(registers[a], int64(b))
	case "banr":
		value = registers[a] & registers[b]
	case "bani":
		value = registers[a] & int64(b)
	case "borr":
		value = registers[a] | registers[b]
	case "bori":
		value = registers[a] | int64(b)
	case "setr":
		value = registers[a]
	case "seti":
		value = int64(a)
	case "gtir":
		if int64(a) > registers[b] {
			value = 1
		}
	case "gtri":
		if registers[a] > int64(b) {
			value = 1
		}
	case "gtrr":
		if registers[a] > registers[b] {
			value = 1
		}
	case "eqir":
		if int64(a) == registers[b] {
			value = 1
		}
	case "eqri":
		if registers[a] == int64(b) {
			value = 1
		}
	case "eqrr":
		if registers[a] == registers[b] {
			value = 1
		}
	}
	if err != nil {
		return err
	}
	registers[instruction.c] = value
	return nil
}

type elfProgram struct {
	ipRegister int
	code       []elfInstruction
}

func parseElfProgram(r io.Reader) (elfProgram, error) {
	program := elfProgram{ipRegister: -1}
	s := bufio.NewScanner(r)
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if program.ipRegister == -1 {
			if len(fields) != 2 || fields[0] != "#ip" {
				return program, fmt.Errorf("missing instruction pointer binding")
			}
			binding, err := convert.ScanNumber[int]([]byte(fields[1]))
			if err != nil {
				return program, err
			}
			if binding < 0 || binding >= 6 {
				return program, fmt.Errorf("invalid instruction pointer register %d", binding)
			}
			program.ipRegister = binding
			continue
		}
		if len(fields) != 4 {
			return program, fmt.Errorf("invalid instruction %q", s.Text())
		}
		instruction := elfInstruction{operation: fields[0]}
		for i, target := range []*int{&instruction.a, &instruction.b, &instruction.c} {
			value, err := convert.ScanNumber[int]([]byte(fields[i+1]))
			if err != nil {
				return program, err
			}
			*target = value
		}
		if err := instruction.validate(6); err != nil {
			return program, fmt.Errorf("instruction %d: %w", len(program.code), err)
		}
		program.code = append(program.code, instruction)
	}
	if err := s.Err(); err != nil {
		return program, err
	}
	if program.ipRegister == -1 || len(program.code) == 0 {
		return program, fmt.Errorf("empty elf program")
	}
	return program, nil
}

type elfMachine struct {
	program   elfProgram
	registers [6]int64
	ip        int64
}

func (machine *elfMachine) running() bool {
	return machine.ip >= 0 && machine.ip < int64(len(machine.program.code))
}

func (machine *elfMachine) step() error {
	machine.registers[machine.program.ipRegister] = machine.ip
	if err := machine.program.code[machine.ip].execute(machine.registers[:]); err != nil {
		return fmt.Errorf("instruction %d: %w", machine.ip, err)
	}
	next := machine.registers[machine.program.ipRegister]
	if next == math.MaxInt64 {
		return fmt.Errorf("instruction pointer overflows int64")
	}
	machine.ip = next + 1
	return nil
}

func (machine *elfMachine) run(limit int) error {
	for range limit {
		if !machine.running() {
			return nil
		}
		if err := machine.step(); err != nil {
			return err
		}
	}
	if machine.running() {
		return fmt.Errorf("program exceeds %d literal instructions", limit)
	}
	return nil
}
