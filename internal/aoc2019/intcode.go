package aoc2019

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"slices"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

const (
	maxIntcodeInput  = 16 << 20
	maxIntcodeMemory = 1_000_000
	maxIntcodeOutput = 1_000_000
	maxIntcodeSteps  = 10_000_000
)

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

type Mode uint8

const (
	ModePosition  Mode = 0
	ModeImmediate Mode = 1
	ModeRelative  Mode = 2
)

type IntcodeComputer struct {
	memory       []int
	ip           int
	input        []int
	output       []int
	halted       bool
	waiting      bool
	relativeBase int
}

// New snapshots program so each machine owns its memory.
func New(program []int) *IntcodeComputer {
	return &IntcodeComputer{
		memory: slices.Clone(program),
	}
}

func ParseIntcodeInput(r io.Reader) ([]int, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxIntcodeInput+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxIntcodeInput {
		return nil, fmt.Errorf("intcode input exceeds 16 MiB")
	}
	var program []int
	for token := range bytes.SplitSeq(data, []byte(",")) {
		if len(program) == maxIntcodeMemory {
			return nil, fmt.Errorf("intcode program exceeds %d words", maxIntcodeMemory)
		}
		value, err := convert.ScanNumber[int](bytes.TrimSpace(token))
		if err != nil {
			return nil, fmt.Errorf("intcode word %d: %w", len(program), err)
		}
		program = append(program, value)
	}
	return program, nil
}

// readMemory returns zero for unwritten nonnegative addresses without growing memory.
func (c *IntcodeComputer) readMemory(addr int) (int, error) {
	if addr < 0 {
		return 0, fmt.Errorf("memory read out of bounds: address %d", addr)
	}
	if addr >= len(c.memory) {
		return 0, nil
	}
	return c.memory[addr], nil
}

func (c *IntcodeComputer) writeMemory(addr, value int) error {
	if addr < 0 {
		return fmt.Errorf("memory write out of bounds: address %d", addr)
	}
	if addr >= maxIntcodeMemory {
		return fmt.Errorf("memory write exceeds %d words: address %d", maxIntcodeMemory, addr)
	}
	if addr >= len(c.memory) {
		newMemory := make([]int, addr+1)
		copy(newMemory, c.memory)
		c.memory = newMemory
	}
	c.memory[addr] = value
	return nil
}

func parseOpcode(value int) Opcode {
	return Opcode(value % 100)
}

func parseMode(instruction int, paramPos int) Mode {
	divisor := 100
	for range paramPos {
		divisor *= 10
	}
	return Mode((instruction / divisor) % 10)
}

func relativeAddress(base, offset int) (int, error) {
	if offset > 0 && base > math.MaxInt-offset || offset < 0 && base < math.MinInt-offset {
		return 0, fmt.Errorf("relative address overflows: base %d, offset %d", base, offset)
	}
	return base + offset, nil
}

func (c *IntcodeComputer) getParameter(mode Mode, offset int) (int, error) {
	param, err := c.readMemory(c.ip + offset)
	if err != nil {
		return 0, err
	}

	switch mode {
	case ModePosition:
		return c.readMemory(param)
	case ModeImmediate:
		return param, nil
	case ModeRelative:
		address, err := relativeAddress(c.relativeBase, param)
		if err != nil {
			return 0, err
		}
		return c.readMemory(address)
	default:
		return 0, fmt.Errorf("unknown parameter mode: %d", mode)
	}
}

func (c *IntcodeComputer) getWriteAddress(mode Mode, offset int) (int, error) {
	param, err := c.readMemory(c.ip + offset)
	if err != nil {
		return 0, err
	}

	switch mode {
	case ModePosition:
		return param, nil
	case ModeRelative:
		return relativeAddress(c.relativeBase, param)
	case ModeImmediate:
		return 0, fmt.Errorf("immediate mode not supported for write operations")
	default:
		return 0, fmt.Errorf("unknown parameter mode: %d", mode)
	}
}

// SetInput snapshots the replacement input; later caller mutations cannot change it.
func (c *IntcodeComputer) SetInput(values ...int) {
	c.input = slices.Clone(values)
	c.waiting = false
}

func (c *IntcodeComputer) AddInput(values ...int) {
	c.input = append(c.input, values...)
	c.waiting = false
}

// GetOutput returns a copy of the buffered output and clears the buffer.
func (c *IntcodeComputer) GetOutput() []int {
	output := slices.Clone(c.output)
	c.output = nil
	return output
}

// LastOutput consumes all buffered output and returns its last value.
func (c *IntcodeComputer) LastOutput() (int, error) {
	output := c.GetOutput()
	if len(output) == 0 {
		return 0, fmt.Errorf("no output produced")
	}
	return output[len(output)-1], nil
}

func (c *IntcodeComputer) IsHalted() bool {
	return c.halted
}

func (c *IntcodeComputer) IsWaiting() bool {
	return c.waiting
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (c *IntcodeComputer) execBinaryOp(instruction int, op func(int, int) int) error {
	val1, err := c.getParameter(parseMode(instruction, 0), 1)
	if err != nil {
		return err
	}

	val2, err := c.getParameter(parseMode(instruction, 1), 2)
	if err != nil {
		return err
	}

	pos3, err := c.getWriteAddress(parseMode(instruction, 2), 3)
	if err != nil {
		return err
	}

	if err := c.writeMemory(pos3, op(val1, val2)); err != nil {
		return err
	}

	c.ip += 4
	return nil
}

func (c *IntcodeComputer) execComparisonOp(instruction int, cmp func(int, int) bool) error {
	val1, err := c.getParameter(parseMode(instruction, 0), 1)
	if err != nil {
		return err
	}

	val2, err := c.getParameter(parseMode(instruction, 1), 2)
	if err != nil {
		return err
	}

	pos3, err := c.getWriteAddress(parseMode(instruction, 2), 3)
	if err != nil {
		return err
	}

	if err := c.writeMemory(pos3, boolToInt(cmp(val1, val2))); err != nil {
		return err
	}

	c.ip += 4
	return nil
}

func (c *IntcodeComputer) executeInstruction() (bool, error) {
	instruction, err := c.readMemory(c.ip)
	if err != nil {
		return false, err
	}
	opcode := parseOpcode(instruction)

	switch opcode {
	case OpHalt:
		c.halted = true
		return true, nil

	case OpAdd:
		return false, c.execBinaryOp(instruction, func(a, b int) int { return a + b })

	case OpMul:
		return false, c.execBinaryOp(instruction, func(a, b int) int { return a * b })

	case OpInput:
		if len(c.input) == 0 {
			c.waiting = true
			return true, nil
		}
		inputVal := c.input[0]

		pos, err := c.getWriteAddress(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}
		if err := c.writeMemory(pos, inputVal); err != nil {
			return false, err
		}
		c.input = c.input[1:]

		c.ip += 2
		return false, nil

	case OpOutput:
		val, err := c.getParameter(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}

		if len(c.output) == maxIntcodeOutput {
			return false, fmt.Errorf("intcode output exceeds %d words", maxIntcodeOutput)
		}
		c.output = append(c.output, val)

		c.ip += 2
		return false, nil

	case OpJumpIfTrue:
		val, err := c.getParameter(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}

		target, err := c.getParameter(parseMode(instruction, 1), 2)
		if err != nil {
			return false, err
		}

		if val != 0 {
			c.ip = target
		} else {
			c.ip += 3
		}
		return false, nil

	case OpJumpIfFalse:
		val, err := c.getParameter(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}

		target, err := c.getParameter(parseMode(instruction, 1), 2)
		if err != nil {
			return false, err
		}

		if val == 0 {
			c.ip = target
		} else {
			c.ip += 3
		}
		return false, nil

	case OpLessThan:
		return false, c.execComparisonOp(instruction, func(a, b int) bool { return a < b })

	case OpEquals:
		return false, c.execComparisonOp(instruction, func(a, b int) bool { return a == b })

	case OpRelativeBase:
		val, err := c.getParameter(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}

		base, err := relativeAddress(c.relativeBase, val)
		if err != nil {
			return false, err
		}
		c.relativeBase = base
		c.ip += 2
		return false, nil

	default:
		return false, fmt.Errorf("unknown opcode: %d at position %d", opcode, c.ip)
	}
}

// Run executes until halt, missing input, or an error.
// It returns an error if a resource limit or the per-call instruction budget is exceeded.
func (c *IntcodeComputer) Run() error {
	if len(c.memory) > maxIntcodeMemory || len(c.input) > maxIntcodeMemory {
		return fmt.Errorf("intcode memory or input exceeds %d words", maxIntcodeMemory)
	}
	for range maxIntcodeSteps {
		halt, err := c.executeInstruction()
		if err != nil {
			return err
		}
		if halt {
			return nil
		}
	}
	return fmt.Errorf("intcode execution exceeds %d instructions without pausing", maxIntcodeSteps)
}

// Memory returns a copy of the current memory state.
func (c *IntcodeComputer) Memory() []int {
	return slices.Clone(c.memory)
}

func (c *IntcodeComputer) SetMemory(addr, value int) error {
	return c.writeMemory(addr, value)
}

// ReadOutputString drains the output buffer, converting each value to a byte.
func (c *IntcodeComputer) ReadOutputString() string {
	output := c.GetOutput()
	var result []byte
	for _, val := range output {
		result = append(result, byte(val))
	}
	return string(result)
}

// Clone creates a deep copy of the computer state.
func (c *IntcodeComputer) Clone() *IntcodeComputer {
	return &IntcodeComputer{
		memory:       slices.Clone(c.memory),
		ip:           c.ip,
		input:        slices.Clone(c.input),
		output:       slices.Clone(c.output),
		halted:       c.halted,
		waiting:      c.waiting,
		relativeBase: c.relativeBase,
	}
}

// StringsToASCII appends the rune values in each line followed by a newline.
// It does not validate that the input is ASCII.
func StringsToASCII(lines ...string) []int {
	var result []int
	for _, line := range lines {
		for _, ch := range line {
			result = append(result, int(ch))
		}
		result = append(result, 10)
	}
	return result
}
