package aoc2018

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

const elfComparisonLimit = 1_000_000

type haltingRecurrence struct {
	initial, seed, orMask, byteMask, valueMask, multiplier, base int64
}

func recognizeHaltingRecurrence(program elfProgram) (haltingRecurrence, error) {
	var recurrence haltingRecurrence
	if len(program.code) != 31 {
		return recurrence, fmt.Errorf("unsupported halting program shape")
	}
	code, ip := program.code, program.ipRegister
	value, source, scratch, product := code[6].a, code[6].c, code[8].c, code[18].c
	registers := map[int]bool{0: true, ip: true, value: true, source: true, scratch: true, product: true}
	if len(registers) != 6 {
		return recurrence, fmt.Errorf("aliased halting program registers")
	}
	want := []elfInstruction{
		{operation: "seti", a: code[0].a, b: code[0].b, c: value},
		{operation: "bani", a: value, b: code[1].b, c: value},
		{operation: "eqri", a: value, b: code[2].b, c: value},
		{operation: "addr", a: value, b: ip, c: ip},
		{operation: "seti", a: 0, b: code[4].b, c: ip},
		{operation: "seti", a: code[5].a, b: code[5].b, c: value},
		{operation: "bori", a: value, b: code[6].b, c: source},
		{operation: "seti", a: code[7].a, b: code[7].b, c: value},
		{operation: "bani", a: source, b: code[8].b, c: scratch},
		{operation: "addr", a: value, b: scratch, c: value},
		{operation: "bani", a: value, b: code[10].b, c: value},
		{operation: "muli", a: value, b: code[11].b, c: value},
		{operation: "bani", a: value, b: code[10].b, c: value},
		{operation: "gtir", a: code[19].b, b: source, c: scratch},
		{operation: "addr", a: scratch, b: ip, c: ip},
		{operation: "addi", a: ip, b: 1, c: ip},
		{operation: "seti", a: 27, b: code[16].b, c: ip},
		{operation: "seti", a: 0, b: code[17].b, c: scratch},
		{operation: "addi", a: scratch, b: 1, c: product},
		{operation: "muli", a: product, b: code[19].b, c: product},
		{operation: "gtrr", a: product, b: source, c: product},
		{operation: "addr", a: product, b: ip, c: ip},
		{operation: "addi", a: ip, b: 1, c: ip},
		{operation: "seti", a: 25, b: code[23].b, c: ip},
		{operation: "addi", a: scratch, b: 1, c: scratch},
		{operation: "seti", a: 17, b: code[25].b, c: ip},
		{operation: "setr", a: scratch, b: code[26].b, c: source},
		{operation: "seti", a: 7, b: code[27].b, c: ip},
		{operation: "eqrr", a: value, b: 0, c: scratch},
		{operation: "addr", a: scratch, b: ip, c: ip},
		{operation: "seti", a: 5, b: code[30].b, c: ip},
	}
	for i, instruction := range want {
		if code[i] != instruction {
			return recurrence, fmt.Errorf("unsupported halting instruction at %d", i)
		}
	}
	if code[0].a&code[1].b != code[2].b && code[2].b != 0 {
		return recurrence, fmt.Errorf("halting program initialization never completes")
	}
	recurrence = haltingRecurrence{
		initial: int64(code[5].a), seed: int64(code[7].a), orMask: int64(code[6].b),
		byteMask: int64(code[8].b), valueMask: int64(code[10].b),
		multiplier: int64(code[11].b), base: int64(code[19].b),
	}
	if recurrence.initial < 0 || recurrence.seed < 0 || recurrence.orMask < 0 || recurrence.byteMask < 0 ||
		recurrence.valueMask < 0 || recurrence.multiplier < 0 || recurrence.base < 2 {
		return recurrence, fmt.Errorf("unsupported halting recurrence constants")
	}
	if _, err := elfAdd(max(recurrence.seed, recurrence.valueMask), recurrence.byteMask); err != nil {
		return recurrence, err
	}
	if _, err := elfMultiply(recurrence.valueMask, recurrence.multiplier); err != nil {
		return recurrence, err
	}
	largestSource := max(recurrence.initial|recurrence.orMask, recurrence.valueMask|recurrence.orMask)
	if largestSource/recurrence.base >= math.MaxInt64/recurrence.base {
		return recurrence, fmt.Errorf("halting division loop can overflow int64")
	}
	return recurrence, nil
}

func (recurrence haltingRecurrence) next(previous int64) int64 {
	source := previous | recurrence.orMask
	value := recurrence.seed
	for {
		value = (value + (source & recurrence.byteMask)) & recurrence.valueMask
		value = (value * recurrence.multiplier) & recurrence.valueMask
		if source < recurrence.base {
			return value
		}
		source /= recurrence.base
	}
}

func readHaltingRecurrence(r io.Reader) (haltingRecurrence, error) {
	program, err := parseElfProgram(r)
	if err != nil {
		return haltingRecurrence{}, err
	}
	return recognizeHaltingRecurrence(program)
}

func day21p01(r io.Reader) (string, error) {
	recurrence, err := readHaltingRecurrence(r)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(recurrence.next(recurrence.initial), 10), nil
}

func day21p02(r io.Reader) (string, error) {
	recurrence, err := readHaltingRecurrence(r)
	if err != nil {
		return "", err
	}

	seen := make(map[int64]bool)
	previous := recurrence.initial
	for {
		// Follow the non-halting branch even when the compared value is zero.
		value := recurrence.next(previous)
		if seen[value] {
			return strconv.FormatInt(previous, 10), nil
		}
		if len(seen) == elfComparisonLimit {
			return "", fmt.Errorf("halting sequence exceeds %d distinct comparisons", elfComparisonLimit)
		}
		seen[value] = true
		previous = value
	}
}
