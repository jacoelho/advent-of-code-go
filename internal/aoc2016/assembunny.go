package aoc2016

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"slices"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type bunnyOpcode byte

const (
	bunnyCopy bunnyOpcode = iota
	bunnyIncrement
	bunnyDecrement
	bunnyJump
	bunnyToggle
	bunnyOut
)

type bunnyOperand struct {
	register int
	literal  int64
}

func (o bunnyOperand) value(registers [4]int64) int64 {
	if o.register >= 0 {
		return registers[o.register]
	}
	return o.literal
}

type bunnyInstruction struct {
	opcode bunnyOpcode
	x, y   bunnyOperand
}

func parseAssembunny(r io.Reader) ([]bunnyInstruction, error) {
	operand := func(text string) (bunnyOperand, error) {
		if len(text) == 1 && text[0] >= 'a' && text[0] <= 'd' {
			return bunnyOperand{register: int(text[0] - 'a')}, nil
		}
		value, err := convert.ScanNumber[int64]([]byte(text))
		return bunnyOperand{register: -1, literal: value}, err
	}
	var program []bunnyInstruction
	s := bufio.NewScanner(r)
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) < 2 || len(fields) > 3 {
			return nil, fmt.Errorf("invalid instruction %q", s.Text())
		}
		instruction := bunnyInstruction{y: bunnyOperand{register: -1}}
		arity := 1
		switch fields[0] {
		case "cpy":
			instruction.opcode, arity = bunnyCopy, 2
		case "inc":
			instruction.opcode = bunnyIncrement
		case "dec":
			instruction.opcode = bunnyDecrement
		case "jnz":
			instruction.opcode, arity = bunnyJump, 2
		case "tgl":
			instruction.opcode = bunnyToggle
		case "out":
			instruction.opcode = bunnyOut
		default:
			return nil, fmt.Errorf("unknown instruction %q", fields[0])
		}
		if len(fields) != arity+1 {
			return nil, fmt.Errorf("wrong operands in %q", s.Text())
		}
		var err error
		instruction.x, err = operand(fields[1])
		if err != nil {
			return nil, fmt.Errorf("instruction %d: %w", len(program), err)
		}
		if arity == 2 {
			instruction.y, err = operand(fields[2])
			if err != nil {
				return nil, fmt.Errorf("instruction %d: %w", len(program), err)
			}
		}
		program = append(program, instruction)
	}
	return program, s.Err()
}

type bunnyMachine struct {
	program   []bunnyInstruction
	registers [4]int64
	ip        int64
	image     string
}

func newBunnyMachine(program []bunnyInstruction, registers [4]int64) bunnyMachine {
	m := bunnyMachine{program: slices.Clone(program), registers: registers}
	m.refreshImage()
	return m
}

func (m *bunnyMachine) refreshImage() {
	code := make([]byte, len(m.program))
	for i, instruction := range m.program {
		code[i] = byte(instruction.opcode)
	}
	m.image = string(code)
}

func (m *bunnyMachine) running() bool { return m.ip >= 0 && m.ip < int64(len(m.program)) }

func bunnyImmediate(o bunnyOperand, value int64) bool { return o.register < 0 && o.literal == value }

func bunnyAddPattern(code []bunnyInstruction) (int, int, bool) {
	if len(code) < 3 {
		return 0, 0, false
	}
	a, counter := code[0].x.register, code[1].x.register
	ok := code[0].opcode == bunnyIncrement && code[1].opcode == bunnyDecrement &&
		code[2].opcode == bunnyJump && code[2].x.register == counter && bunnyImmediate(code[2].y, -2) &&
		a >= 0 && counter >= 0 && a != counter
	return a, counter, ok
}

func (m *bunnyMachine) accelerate() bool {
	code := m.program[m.ip:]
	if len(code) >= 6 && code[0].opcode == bunnyCopy {
		a, inner, ok := bunnyAddPattern(code[1:])
		outer := code[4].x.register
		source := code[0].x
		if ok && code[0].y.register == inner && code[4].opcode == bunnyDecrement && outer >= 0 &&
			outer != inner && outer != a && code[5].opcode == bunnyJump && code[5].x.register == outer &&
			bunnyImmediate(code[5].y, -5) && source.register != a && source.register != inner && source.register != outer {
			x, y := source.value(m.registers), m.registers[outer]
			if x > 0 && y > 0 && x <= math.MaxInt64/y && m.registers[a] <= math.MaxInt64-x*y {
				m.registers[a] += x * y
				m.registers[inner], m.registers[outer] = 0, 0
				m.ip += 6
				return true
			}
		}
	}
	if a, counter, ok := bunnyAddPattern(code); ok {
		count := m.registers[counter]
		if count > 0 && m.registers[a] <= math.MaxInt64-count {
			m.registers[a] += count
			m.registers[counter] = 0
			m.ip += 3
			return true
		}
	}
	return false
}

func (m *bunnyMachine) step() (int64, bool, error) {
	// Match current instructions on every visit: tgl can change a previously seen loop.
	if m.accelerate() {
		return 0, false, nil
	}
	instruction := m.program[m.ip]
	offset := int64(1)
	var output int64
	emitted := false
	switch instruction.opcode {
	case bunnyCopy:
		if instruction.y.register >= 0 {
			m.registers[instruction.y.register] = instruction.x.value(m.registers)
		}
	case bunnyIncrement, bunnyDecrement:
		if instruction.x.register >= 0 {
			value := &m.registers[instruction.x.register]
			if instruction.opcode == bunnyIncrement {
				if *value == math.MaxInt64 {
					return 0, false, fmt.Errorf("register overflow at %d", m.ip)
				}
				*value++
			} else {
				if *value == math.MinInt64 {
					return 0, false, fmt.Errorf("register overflow at %d", m.ip)
				}
				*value--
			}
		}
	case bunnyJump:
		if instruction.x.value(m.registers) != 0 {
			offset = instruction.y.value(m.registers)
		}
	case bunnyToggle:
		distance := instruction.x.value(m.registers)
		if distance >= -m.ip && distance < int64(len(m.program))-m.ip {
			target := &m.program[m.ip+distance]
			switch target.opcode {
			case bunnyIncrement:
				target.opcode = bunnyDecrement
			case bunnyJump:
				target.opcode = bunnyCopy
			case bunnyCopy:
				target.opcode = bunnyJump
			default:
				target.opcode = bunnyIncrement
			}
			m.refreshImage()
		}
	case bunnyOut:
		output, emitted = instruction.x.value(m.registers), true
	}
	if offset < -m.ip || offset >= int64(len(m.program))-m.ip {
		m.ip = int64(len(m.program))
	} else {
		m.ip += offset
	}
	return output, emitted, nil
}

type bunnyState struct {
	registers [4]int64
	ip        int64
	image     string
	next      int64
}

func (m *bunnyMachine) state(next int64) bunnyState {
	return bunnyState{m.registers, m.ip, m.image, next}
}

const bunnyStateLimit = 1_000_000

func runAssembunny(program []bunnyInstruction, initial [4]int64) ([4]int64, error) {
	return executeAssembunny(program, initial, bunnyStateLimit)
}

func executeAssembunny(program []bunnyInstruction, initial [4]int64, stateLimit int) ([4]int64, error) {
	machine := newBunnyMachine(program, initial)
	seen := make(map[bunnyState]struct{})
	for machine.running() {
		state := machine.state(0)
		if _, exists := seen[state]; exists {
			return [4]int64{}, fmt.Errorf("program repeats a state at %d", machine.ip)
		}
		if len(seen) >= stateLimit {
			return [4]int64{}, fmt.Errorf("execution exceeded %d states", stateLimit)
		}
		seen[state] = struct{}{}
		if _, _, err := machine.step(); err != nil {
			return [4]int64{}, err
		}
	}
	return machine.registers, nil
}
