package aoc2019

import (
	"bytes"
	"fmt"
	"io"
	"maps"
	"slices"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

// Opcode is the operation encoded in the two lowest decimal digits of an
// instruction.
type Opcode int

const (
	OpAdd          Opcode = 1
	OpMul          Opcode = 2
	OpInput        Opcode = 3
	OpOutput       Opcode = 4
	OpJumpIfTrue   Opcode = 5
	OpJumpIfFalse  Opcode = 6
	OpLessThan     Opcode = 7
	OpEquals       Opcode = 8
	OpRelativeBase Opcode = 9
	OpHalt         Opcode = 99
)

// Mode is how an instruction interprets a parameter.
type Mode uint8

// Parameter modes: a memory address, a literal value, or an address relative
// to the relative base.
const (
	ModePosition  Mode = 0
	ModeImmediate Mode = 1
	ModeRelative  Mode = 2
)

// denseGrowthSlack bounds how far past the dense memory a write may land and
// still extend it; farther writes go to a sparse map so a single large address
// cannot force a huge allocation.
const denseGrowthSlack = 1 << 16

var modeDivisors = [3]int{100, 1000, 10000}

// IntcodeComputer runs an Intcode program synchronously with buffered input
// and output. Run executes until the program halts or needs input it does not
// have; the caller owns the machine and resumes it explicitly.
type IntcodeComputer struct {
	memory       []int
	far          map[int]int
	ip           int
	relativeBase int
	input        []int
	output       []int
	halted       bool
	waiting      bool
}

// New snapshots program so each machine owns its memory.
func New(program []int) *IntcodeComputer {
	return &IntcodeComputer{memory: slices.Clone(program)}
}

// ParseIntcodeInput parses a comma-separated Intcode program, trimming space
// around each word. Empty input or an empty or invalid word returns an error.
func ParseIntcodeInput(r io.Reader) ([]int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var program []int
	for token := range bytes.SplitSeq(data, []byte(",")) {
		value, err := convert.ScanNumber[int](bytes.TrimSpace(token))
		if err != nil {
			return nil, fmt.Errorf("intcode word %d: %w", len(program), err)
		}
		program = append(program, value)
	}
	return program, nil
}

func (c *IntcodeComputer) load(addr int) int {
	if uint(addr) < uint(len(c.memory)) {
		return c.memory[addr]
	}
	return c.far[addr]
}

func (c *IntcodeComputer) store(addr, value int) {
	switch {
	case addr < len(c.memory):
		c.memory[addr] = value
	case addr-len(c.memory) < denseGrowthSlack:
		c.memory = slices.Grow(c.memory, max(addr+1, 2*len(c.memory))-len(c.memory))
		c.memory = c.memory[:cap(c.memory)]
		for at, v := range c.far {
			if at < len(c.memory) {
				c.memory[at] = v
				delete(c.far, at)
			}
		}
		c.memory[addr] = value
	default:
		if c.far == nil {
			c.far = make(map[int]int)
		}
		c.far[addr] = value
	}
}

func relativeAddress(base, offset int) (int, error) {
	address, ok := xmath.CheckedAdd(base, offset)
	if !ok {
		return 0, fmt.Errorf("relative address overflows: base %d, offset %d", base, offset)
	}
	return address, nil
}

// address resolves parameter n of the current instruction to a nonnegative
// memory address.
func (c *IntcodeComputer) address(instruction, n int) (int, error) {
	param := c.load(c.ip + 1 + n)
	addr := param
	switch mode := Mode(instruction / modeDivisors[n] % 10); mode {
	case ModePosition:
	case ModeRelative:
		var err error
		if addr, err = relativeAddress(c.relativeBase, param); err != nil {
			return 0, err
		}
	case ModeImmediate:
		return 0, fmt.Errorf("immediate mode for write parameter at position %d", c.ip)
	default:
		return 0, fmt.Errorf("unknown parameter mode %d at position %d", mode, c.ip)
	}
	if addr < 0 {
		return 0, fmt.Errorf("negative memory address %d at position %d", addr, c.ip)
	}
	return addr, nil
}

func (c *IntcodeComputer) value(instruction, n int) (int, error) {
	if Mode(instruction/modeDivisors[n]%10) == ModeImmediate {
		return c.load(c.ip + 1 + n), nil
	}
	addr, err := c.address(instruction, n)
	if err != nil {
		return 0, err
	}
	return c.load(addr), nil
}

func (c *IntcodeComputer) operands(instruction int) (int, int, error) {
	a, err := c.value(instruction, 0)
	if err != nil {
		return 0, 0, err
	}
	b, err := c.value(instruction, 1)
	return a, b, err
}

func (c *IntcodeComputer) jump(target int) error {
	if target < 0 {
		return fmt.Errorf("jump to negative address %d at position %d", target, c.ip)
	}
	c.ip = target
	return nil
}

// SetInput replaces the pending input with a copy of values.
func (c *IntcodeComputer) SetInput(values ...int) {
	c.input = slices.Clone(values)
	c.waiting = false
}

// AddInput appends values to the pending input.
func (c *IntcodeComputer) AddInput(values ...int) {
	c.input = append(c.input, values...)
	c.waiting = false
}

// GetOutput returns the buffered output and clears the buffer.
func (c *IntcodeComputer) GetOutput() []int {
	output := c.output
	c.output = nil
	return output
}

// ReadOutputString consumes the buffered output as bytes.
func (c *IntcodeComputer) ReadOutputString() string {
	output := c.GetOutput()
	result := make([]byte, len(output))
	for i, value := range output {
		result[i] = byte(value)
	}
	return string(result)
}

// IsHalted reports whether the program has executed its halt instruction.
func (c *IntcodeComputer) IsHalted() bool {
	return c.halted
}

// IsWaiting reports whether Run stopped at an input instruction with no
// pending input. SetInput and AddInput clear it.
func (c *IntcodeComputer) IsWaiting() bool {
	return c.waiting
}

// ReadMemory returns the word at addr; unwritten addresses hold zero.
func (c *IntcodeComputer) ReadMemory(addr int) (int, error) {
	if addr < 0 {
		return 0, fmt.Errorf("negative memory address %d", addr)
	}
	return c.load(addr), nil
}

// SetMemory writes value at addr.
func (c *IntcodeComputer) SetMemory(addr, value int) error {
	if addr < 0 {
		return fmt.Errorf("negative memory address %d", addr)
	}
	c.store(addr, value)
	return nil
}

// Clone returns an independent copy of the machine state.
func (c *IntcodeComputer) Clone() *IntcodeComputer {
	clone := *c
	clone.memory = slices.Clone(c.memory)
	clone.far = maps.Clone(c.far)
	clone.input = slices.Clone(c.input)
	clone.output = slices.Clone(c.output)
	return &clone
}

// RunToHalt runs the program and returns an error if it stops to wait for
// input instead of halting.
func (c *IntcodeComputer) RunToHalt() error {
	if err := c.Run(); err != nil {
		return err
	}
	if !c.halted {
		return fmt.Errorf("intcode program requested input at position %d after its input ran out", c.ip)
	}
	return nil
}

// Run executes until the program halts, needs input it does not have, or
// fails.
func (c *IntcodeComputer) Run() error {
	for !c.halted {
		instruction := c.load(c.ip)
		switch opcode := Opcode(instruction % 100); opcode {
		case OpAdd, OpMul, OpLessThan, OpEquals:
			a, b, err := c.operands(instruction)
			if err != nil {
				return err
			}
			target, err := c.address(instruction, 2)
			if err != nil {
				return err
			}
			var result int
			switch opcode {
			case OpAdd:
				result = a + b
			case OpMul:
				result = a * b
			case OpLessThan:
				result = boolToInt(a < b)
			case OpEquals:
				result = boolToInt(a == b)
			}
			c.store(target, result)
			c.ip += 4
		case OpInput:
			if len(c.input) == 0 {
				c.waiting = true
				return nil
			}
			target, err := c.address(instruction, 0)
			if err != nil {
				return err
			}
			c.store(target, c.input[0])
			c.input = c.input[1:]
			c.ip += 2
		case OpOutput:
			value, err := c.value(instruction, 0)
			if err != nil {
				return err
			}
			c.output = append(c.output, value)
			c.ip += 2
		case OpJumpIfTrue, OpJumpIfFalse:
			condition, target, err := c.operands(instruction)
			if err != nil {
				return err
			}
			if (condition != 0) == (opcode == OpJumpIfTrue) {
				if err := c.jump(target); err != nil {
					return err
				}
			} else {
				c.ip += 3
			}
		case OpRelativeBase:
			offset, err := c.value(instruction, 0)
			if err != nil {
				return err
			}
			base, err := relativeAddress(c.relativeBase, offset)
			if err != nil {
				return err
			}
			c.relativeBase = base
			c.ip += 2
		case OpHalt:
			c.halted = true
		default:
			return fmt.Errorf("unknown opcode: %d at position %d", opcode, c.ip)
		}
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// StringsToASCII returns the bytes of each line followed by a newline.
func StringsToASCII(lines ...string) []int {
	var result []int
	for _, line := range lines {
		for i := range len(line) {
			result = append(result, int(line[i]))
		}
		result = append(result, '\n')
	}
	return result
}
