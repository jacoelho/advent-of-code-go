package aoc2016

import (
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
)

// isClockSignal reports whether the program, started with register a set to
// initial, emits 0, 1, 0, 1, ... forever: it must reach a repeated state after
// emitting only matching values, and the cycle must emit output.
func isClockSignal(program []bunnyInstruction, initial int64) (bool, error) {
	machine := newBunnyMachine(program, [4]int64{initial, 0, 0, 0})
	next := int64(0)
	outcome, cycleOutputs, err := machine.run(func(value int64) bool {
		if value != next {
			return false
		}
		next ^= 1
		return true
	})
	if err != nil {
		return false, fmt.Errorf("clock candidate %d: %w", initial, err)
	}
	return outcome == bunnyCycled && cycleOutputs > 0, nil
}

func day25p01(r io.Reader) (string, error) {
	program, err := parseAssembunny(r)
	if err != nil {
		return "", err
	}
	// tgl never produces out, so a program without one can never emit a clock signal.
	if !slices.ContainsFunc(program, func(instruction bunnyInstruction) bool { return instruction.opcode == bunnyOut }) {
		return "", fmt.Errorf("program has no out instruction")
	}
	for candidate := int64(1); candidate < math.MaxInt64; candidate++ {
		valid, err := isClockSignal(program, candidate)
		if err != nil {
			return "", err
		}
		if valid {
			return strconv.FormatInt(candidate, 10), nil
		}
	}
	return "", fmt.Errorf("no positive clock seed fits int64")
}
