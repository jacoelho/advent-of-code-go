package aoc2019

import (
	"fmt"
	"io"
	"slices"

	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
)

// Opcode represents an Intcode operation
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

// Mode represents a parameter mode
type Mode uint8

const (
	ModePosition  Mode = 0 // position mode
	ModeImmediate Mode = 1 // immediate mode
	ModeRelative  Mode = 2 // relative mode
)

// IntcodeComputer represents the state of an Intcode computer
type IntcodeComputer struct {
	memory       []int
	ip           int
	input        []int
	output       []int
	halted       bool
	waiting      bool // waiting for input
	relativeBase int
}

// New creates a new IntcodeComputer with the given program
func New(program []int) *IntcodeComputer {
	return &IntcodeComputer{
		memory:  slices.Clone(program),
		ip:      0,
		halted:  false,
		waiting: false,
	}
}

// ParseIntcodeInput parses comma-separated integers from input
func ParseIntcodeInput(r io.Reader) ([]int, error) {
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte(",")), convert.ScanNumber[int])
	program := slices.Collect(s.Values())
	return program, s.Err()
}

// readMemory reads a value from memory at the given address
func (c *IntcodeComputer) readMemory(addr int) (int, error) {
	if addr < 0 {
		return 0, fmt.Errorf("memory read out of bounds: address %d", addr)
	}
	if addr >= len(c.memory) {
		return 0, nil
	}
	return c.memory[addr], nil
}

// writeMemory writes a value to memory at the given address
func (c *IntcodeComputer) writeMemory(addr, value int) error {
	if addr < 0 {
		return fmt.Errorf("memory write out of bounds: address %d", addr)
	}
	if addr >= len(c.memory) {
		newMemory := make([]int, addr+1)
		copy(newMemory, c.memory)
		c.memory = newMemory
	}
	c.memory[addr] = value
	return nil
}

// parseOpcode extracts the opcode from an instruction value
func parseOpcode(value int) Opcode {
	return Opcode(value % 100)
}

// parseMode extracts the parameter mode for a given parameter position
func parseMode(instruction int, paramPos int) Mode {
	divisor := 100
	for range paramPos {
		divisor *= 10
	}
	return Mode((instruction / divisor) % 10)
}

// getParameter reads a parameter value based on its mode
func (c *IntcodeComputer) getParameter(mode Mode, offset int) (int, error) {
	if c.ip+offset >= len(c.memory) {
		return 0, fmt.Errorf("parameter read out of bounds: ip=%d, offset=%d", c.ip, offset)
	}

	param := c.memory[c.ip+offset]

	switch mode {
	case ModePosition:
		return c.readMemory(param)
	case ModeImmediate:
		return param, nil
	case ModeRelative:
		return c.readMemory(param + c.relativeBase)
	default:
		return 0, fmt.Errorf("unknown parameter mode: %d", mode)
	}
}

// getWriteAddress computes the address for write operations based on parameter mode
func (c *IntcodeComputer) getWriteAddress(mode Mode, offset int) (int, error) {
	if c.ip+offset >= len(c.memory) {
		return 0, fmt.Errorf("write address out of bounds: ip=%d, offset=%d", c.ip, offset)
	}

	param := c.memory[c.ip+offset]

	switch mode {
	case ModePosition:
		return param, nil
	case ModeRelative:
		return param + c.relativeBase, nil
	case ModeImmediate:
		return 0, fmt.Errorf("immediate mode not supported for write operations")
	default:
		return 0, fmt.Errorf("unknown parameter mode: %d", mode)
	}
}

// SetInput sets the input buffer for the computer
func (c *IntcodeComputer) SetInput(values ...int) {
	c.input = values
	c.waiting = false
}

// AddInput appends values to the input buffer
func (c *IntcodeComputer) AddInput(values ...int) {
	c.input = append(c.input, values...)
	c.waiting = false
}

// GetOutput returns the output buffer
func (c *IntcodeComputer) GetOutput() []int {
	return c.output
}

// LastOutput returns the last output value, or an error if no output was produced
func (c *IntcodeComputer) LastOutput() (int, error) {
	if len(c.output) == 0 {
		return 0, fmt.Errorf("no output produced")
	}
	return c.output[len(c.output)-1], nil
}

// IsHalted returns true if the computer has halted
func (c *IntcodeComputer) IsHalted() bool {
	return c.halted
}

// IsWaiting returns true if the computer is waiting for input
func (c *IntcodeComputer) IsWaiting() bool {
	return c.waiting
}

// executeInstruction executes the instruction at the current IP
// Returns true if the program should halt, false otherwise
func (c *IntcodeComputer) executeInstruction() (bool, error) {
	if c.ip >= len(c.memory) {
		return false, fmt.Errorf("instruction pointer out of bounds: %d", c.ip)
	}

	opcode := parseOpcode(c.memory[c.ip])

	instruction := c.memory[c.ip]

	switch opcode {
	case OpHalt:
		c.halted = true
		return true, nil

	case OpAdd:
		if c.ip+3 >= len(c.memory) {
			return false, fmt.Errorf("incomplete instruction at position %d", c.ip)
		}

		val1, err := c.getParameter(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}

		val2, err := c.getParameter(parseMode(instruction, 1), 2)
		if err != nil {
			return false, err
		}

		pos3, err := c.getWriteAddress(parseMode(instruction, 2), 3)
		if err != nil {
			return false, err
		}
		if err := c.writeMemory(pos3, val1+val2); err != nil {
			return false, err
		}

		c.ip += 4
		return false, nil

	case OpMul:
		if c.ip+3 >= len(c.memory) {
			return false, fmt.Errorf("incomplete instruction at position %d", c.ip)
		}

		val1, err := c.getParameter(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}

		val2, err := c.getParameter(parseMode(instruction, 1), 2)
		if err != nil {
			return false, err
		}

		pos3, err := c.getWriteAddress(parseMode(instruction, 2), 3)
		if err != nil {
			return false, err
		}
		if err := c.writeMemory(pos3, val1*val2); err != nil {
			return false, err
		}

		c.ip += 4
		return false, nil

	case OpInput:
		if c.ip+1 >= len(c.memory) {
			return false, fmt.Errorf("incomplete instruction at position %d", c.ip)
		}

		if len(c.input) == 0 {
			c.waiting = true
			return true, nil // Pause execution, waiting for input
		}

		pos, err := c.getWriteAddress(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}
		if err := c.writeMemory(pos, c.input[0]); err != nil {
			return false, err
		}
		c.input = c.input[1:]

		c.ip += 2
		return false, nil

	case OpOutput:
		if c.ip+1 >= len(c.memory) {
			return false, fmt.Errorf("incomplete instruction at position %d", c.ip)
		}

		val, err := c.getParameter(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}

		c.output = append(c.output, val)

		c.ip += 2
		return false, nil

	case OpJumpIfTrue:
		if c.ip+2 >= len(c.memory) {
			return false, fmt.Errorf("incomplete instruction at position %d", c.ip)
		}

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
		if c.ip+2 >= len(c.memory) {
			return false, fmt.Errorf("incomplete instruction at position %d", c.ip)
		}

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
		if c.ip+3 >= len(c.memory) {
			return false, fmt.Errorf("incomplete instruction at position %d", c.ip)
		}

		val1, err := c.getParameter(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}

		val2, err := c.getParameter(parseMode(instruction, 1), 2)
		if err != nil {
			return false, err
		}

		pos3, err := c.getWriteAddress(parseMode(instruction, 2), 3)
		if err != nil {
			return false, err
		}
		if val1 < val2 {
			if err := c.writeMemory(pos3, 1); err != nil {
				return false, err
			}
		} else {
			if err := c.writeMemory(pos3, 0); err != nil {
				return false, err
			}
		}

		c.ip += 4
		return false, nil

	case OpEquals:
		if c.ip+3 >= len(c.memory) {
			return false, fmt.Errorf("incomplete instruction at position %d", c.ip)
		}

		val1, err := c.getParameter(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}

		val2, err := c.getParameter(parseMode(instruction, 1), 2)
		if err != nil {
			return false, err
		}

		pos3, err := c.getWriteAddress(parseMode(instruction, 2), 3)
		if err != nil {
			return false, err
		}
		if val1 == val2 {
			if err := c.writeMemory(pos3, 1); err != nil {
				return false, err
			}
		} else {
			if err := c.writeMemory(pos3, 0); err != nil {
				return false, err
			}
		}

		c.ip += 4
		return false, nil

	case OpRelativeBase:
		if c.ip+1 >= len(c.memory) {
			return false, fmt.Errorf("incomplete instruction at position %d", c.ip)
		}

		val, err := c.getParameter(parseMode(instruction, 0), 1)
		if err != nil {
			return false, err
		}

		c.relativeBase += val
		c.ip += 2
		return false, nil

	default:
		return false, fmt.Errorf("unknown opcode: %d at position %d", opcode, c.ip)
	}
}

// Run executes the Intcode program until it halts
func (c *IntcodeComputer) Run() error {
	for {
		halt, err := c.executeInstruction()
		if err != nil {
			return err
		}
		if halt {
			return nil
		}
	}
}

// Memory returns a copy of the current memory state
func (c *IntcodeComputer) Memory() []int {
	return slices.Clone(c.memory)
}

// SetMemory sets the value at a specific memory address
func (c *IntcodeComputer) SetMemory(addr, value int) error {
	return c.writeMemory(addr, value)
}
