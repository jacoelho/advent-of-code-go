package aoc2016

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

func isClockSignal(program []bunnyInstruction, initial int64, remaining *int) (bool, error) {
	machine := newBunnyMachine(program, [4]int64{initial, 0, 0, 0})
	seen := make(map[bunnyState]int)
	next, outputs := int64(0), 0
	for machine.running() {
		state := machine.state(next)
		if previous, exists := seen[state]; exists {
			// A repeated machine state proves an infinite clock only if the cycle emitted output.
			return outputs > previous, nil
		}
		if len(seen) == bunnyStateLimit {
			return false, fmt.Errorf("clock candidate %d exceeded %d states without a proof", initial, bunnyStateLimit)
		}
		seen[state] = outputs
		if *remaining <= 0 {
			return false, fmt.Errorf("clock work budget exhausted at candidate %d", initial)
		}
		*remaining--
		value, emitted, err := machine.step()
		if err != nil {
			return false, err
		}
		if emitted {
			if value != next {
				return false, nil
			}
			next ^= 1
			outputs++
		}
	}
	return false, nil
}

func day25p01(r io.Reader) (string, error) {
	program, err := parseAssembunny(r)
	if err != nil {
		return "", err
	}
	seed, err := clockSeed(program, 10_000_000)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(seed, 10), nil
}

func clockSeed(program []bunnyInstruction, remaining int) (int64, error) {
	for candidate := int64(1); candidate < math.MaxInt64; candidate++ {
		if remaining <= 0 {
			return 0, fmt.Errorf("clock work budget exhausted at candidate %d", candidate)
		}
		remaining--
		valid, err := isClockSignal(program, candidate, &remaining)
		if err != nil {
			return 0, err
		}
		if valid {
			return candidate, nil
		}
	}
	return 0, fmt.Errorf("no positive clock seed fits int64")
}
