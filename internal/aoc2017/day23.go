package aoc2017

import (
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"
)

func parseCoprocessor(r io.Reader) ([]assemblyInstruction, error) {
	program, err := parseAssembly(r)
	if err != nil {
		return nil, err
	}
	for _, instruction := range program {
		switch instruction.operation {
		case "set", "sub", "mul", "jnz":
		default:
			return nil, fmt.Errorf("unsupported coprocessor operation %q", instruction.operation)
		}
	}

	return program, nil
}

func runCoprocessor(program []assemblyInstruction, registers [26]int64, limit int) ([26]int64, int, error) {
	pointer, multiplications := int64(0), 0
	for steps := 0; pointer >= 0 && pointer < int64(len(program)); steps++ {
		if steps == limit {
			return registers, 0, fmt.Errorf("coprocessor exceeded %d instructions", limit)
		}
		instruction := program[pointer]
		if instruction.operation == "mul" {
			multiplications++
		}
		next, err := advanceAssembly(instruction, &registers, pointer)
		if err != nil {
			return registers, 0, fmt.Errorf("coprocessor instruction %d: %w", pointer, err)
		}
		pointer = next
	}

	return registers, multiplications, nil
}

type compositeRange struct{ start, end, step int64 }

func coprocessorComposites(program []assemblyInstruction) (compositeRange, bool) {
	const template = `set b 0
set c b
jnz a 2
jnz 1 5
mul b 0
sub b 0
set c b
sub c 0
set f 1
set d 2
set e 2
set g d
mul g e
sub g b
jnz g 2
set f 0
sub e -1
set g e
sub g b
jnz g -8
sub d -1
set g d
sub g b
jnz g -13
jnz f 2
sub h -1
set g b
sub g c
jnz g 2
jnz 1 3
sub b 0
jnz 1 -23`
	if len(program) != 32 {
		return compositeRange{}, false
	}
	expected, err := parseAssembly(strings.NewReader(template))
	if err != nil {
		return compositeRange{}, false
	}
	for _, index := range []int{0, 4, 5, 7, 30} {
		if program[index].y.register != -1 {
			return compositeRange{}, false
		}
		expected[index].y = program[index].y
	}
	if !slices.Equal(program, expected) {
		return compositeRange{}, false
	}

	var registers [26]int64
	registers[0] = 1
	for pointer := int64(0); pointer < 8; {
		next, err := advanceAssembly(program[pointer], &registers, pointer)
		if err != nil {
			return compositeRange{}, false
		}
		pointer = next
	}
	start, end := registers['b'-'a'], registers['c'-'a']
	decrement := program[30].y.literal
	if start < 3 || end < start || decrement >= 0 || decrement == math.MinInt64 {
		return compositeRange{}, false
	}
	step := -decrement
	if (end-start)%step != 0 || end-1 > math.MaxInt64/(end-1) {
		return compositeRange{}, false
	}

	return compositeRange{start: start, end: end, step: step}, true
}

func countComposites(progression compositeRange) (int, error) {
	count, work := 0, 0
	for value := progression.start; ; value += progression.step {
		work++
		for divisor := int64(2); divisor <= value/divisor; divisor++ {
			work++
			if work > 10_000_000 {
				return 0, fmt.Errorf("composite search exceeded 10000000 operations")
			}
			if value%divisor == 0 {
				count++
				break
			}
		}
		if value == progression.end {
			return count, nil
		}
	}
}

func day23p01(r io.Reader) (string, error) {
	program, err := parseCoprocessor(r)
	if err != nil {
		return "", err
	}
	_, multiplications, err := runCoprocessor(program, [26]int64{}, 10_000_000)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(multiplications), nil
}

func day23p02(r io.Reader) (string, error) {
	program, err := parseCoprocessor(r)
	if err != nil {
		return "", err
	}
	// The full loop proves h counts composites; matching only its prefix is insufficient.
	if progression, ok := coprocessorComposites(program); ok {
		count, err := countComposites(progression)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(count), nil
	}

	var registers [26]int64
	registers[0] = 1
	registers, _, err = runCoprocessor(program, registers, 10_000_000)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(registers['h'-'a'], 10), nil
}
