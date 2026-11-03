package aoc2018

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/graph"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type opcodeSample struct {
	before, after [4]int64
	instruction   [4]int
}

func parseOpcodeInstruction(line string) ([4]int, error) {
	var instruction [4]int
	fields := strings.Fields(line)
	if len(fields) != len(instruction) {
		return instruction, fmt.Errorf("invalid numeric instruction %q", line)
	}
	for i, field := range fields {
		value, err := convert.ScanNumber[int]([]byte(field))
		if err != nil {
			return instruction, err
		}
		instruction[i] = value
	}
	if instruction[0] < 0 || instruction[0] >= 16 || instruction[3] < 0 || instruction[3] >= 4 {
		return instruction, fmt.Errorf("invalid opcode or destination in %q", line)
	}
	return instruction, nil
}

func sampleRegisters(line, prefix string) ([4]int64, error) {
	var registers [4]int64
	text, ok := strings.CutPrefix(line, prefix)
	if !ok {
		return registers, fmt.Errorf("expected %s register sample", prefix)
	}
	text = strings.TrimSpace(text)
	if len(text) < 2 || text[0] != '[' || text[len(text)-1] != ']' {
		return registers, fmt.Errorf("invalid register sample %q", line)
	}
	fields := strings.Split(text[1:len(text)-1], ",")
	if len(fields) != len(registers) {
		return registers, fmt.Errorf("expected four registers in %q", line)
	}
	for i, field := range fields {
		value, err := convert.ScanNumber[int64]([]byte(strings.TrimSpace(field)))
		if err != nil {
			return registers, err
		}
		registers[i] = value
	}
	return registers, nil
}

func parseOpcodeInput(r io.Reader) ([]opcodeSample, [][4]int, error) {
	s := bufio.NewScanner(r)
	var lines []string
	for s.Scan() {
		if line := strings.TrimSpace(s.Text()); line != "" {
			lines = append(lines, line)
		}
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}

	var samples []opcodeSample
	var program [][4]int
	for i := 0; i < len(lines); {
		if !strings.HasPrefix(lines[i], "Before:") {
			instruction, err := parseOpcodeInstruction(lines[i])
			if err != nil {
				return nil, nil, err
			}
			program = append(program, instruction)
			i++
			continue
		}
		if len(program) != 0 || len(lines)-i < 3 {
			return nil, nil, fmt.Errorf("misplaced or incomplete opcode sample")
		}
		before, err := sampleRegisters(lines[i], "Before:")
		if err != nil {
			return nil, nil, err
		}
		instruction, err := parseOpcodeInstruction(lines[i+1])
		if err != nil {
			return nil, nil, err
		}
		after, err := sampleRegisters(lines[i+2], "After:")
		if err != nil {
			return nil, nil, err
		}
		samples = append(samples, opcodeSample{before: before, after: after, instruction: instruction})
		i += 3
	}
	return samples, program, nil
}

func matchingOperations(sample opcodeSample) uint16 {
	var matches uint16
	for i, operation := range strings.Fields(elfOperationNames) {
		instruction := elfInstruction{operation: operation, a: sample.instruction[1], b: sample.instruction[2], c: sample.instruction[3]}
		if instruction.validate(4) != nil {
			continue
		}
		registers := sample.before
		if instruction.execute(registers[:]) != nil {
			continue
		}
		if registers == sample.after {
			matches |= 1 << i
		}
	}
	return matches
}

func deduceOpcodes(samples []opcodeSample) ([16]int, error) {
	var candidates [16]uint16
	for opcode := range candidates {
		candidates[opcode] = ^uint16(0)
	}
	for _, sample := range samples {
		candidates[sample.instruction[0]] &= matchingOperations(sample)
	}
	options := make([][]int, len(candidates))
	for opcode, mask := range candidates {
		for operation := range 16 {
			if mask&(1<<operation) != 0 {
				options[opcode] = append(options[opcode], operation)
			}
		}
	}
	var assignment [16]int
	matched, err := graph.UniquePerfectMatching(options)
	if errors.Is(err, graph.ErrNoPerfectMatching) {
		return assignment, fmt.Errorf("inconsistent opcode samples")
	}
	if errors.Is(err, graph.ErrAmbiguousMatching) {
		return assignment, fmt.Errorf("ambiguous opcode samples")
	}
	if err != nil {
		return assignment, fmt.Errorf("opcode mapping: %w", err)
	}
	copy(assignment[:], matched)
	return assignment, nil
}

func day16p01(r io.Reader) (string, error) {
	samples, _, err := parseOpcodeInput(r)
	if err != nil {
		return "", err
	}

	count := xiter.FromSlice(samples).Filter(func(sample opcodeSample) bool {
		return bits.OnesCount16(matchingOperations(sample)) >= 3
	}).Count()
	return strconv.Itoa(count), nil
}

func day16p02(r io.Reader) (string, error) {
	samples, program, err := parseOpcodeInput(r)
	if err != nil {
		return "", err
	}
	assignment, err := deduceOpcodes(samples)
	if err != nil {
		return "", err
	}

	names := strings.Fields(elfOperationNames)
	var registers [4]int64
	for _, encoded := range program {
		instruction := elfInstruction{operation: names[assignment[encoded[0]]], a: encoded[1], b: encoded[2], c: encoded[3]}
		if err := instruction.validate(len(registers)); err != nil {
			return "", err
		}
		if err := instruction.execute(registers[:]); err != nil {
			return "", err
		}
	}
	return strconv.FormatInt(registers[0], 10), nil
}
