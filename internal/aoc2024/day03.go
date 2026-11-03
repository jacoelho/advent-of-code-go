package aoc2024

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

var memoryInstructionPattern = regexp.MustCompile(`do\(\)|don't\(\)|mul\(([0-9]{1,3}),([0-9]{1,3})\)`)

type memoryOperation uint8

const (
	memoryMultiply memoryOperation = iota
	memoryEnable
	memoryDisable
)

type memoryInstruction struct {
	operation memoryOperation
	a, b      int64
}

func parseCorruptedMemory(r io.Reader) ([]memoryInstruction, error) {
	const inputLimit = 4 << 20
	input, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return nil, err
	}
	if len(input) > inputLimit {
		return nil, fmt.Errorf("corrupted memory exceeds 4 MiB")
	}

	var instructions []memoryInstruction
	for _, match := range memoryInstructionPattern.FindAllStringSubmatch(string(input), -1) {
		var instruction memoryInstruction
		switch {
		case match[0] == "do()":
			instruction.operation = memoryEnable
		case match[0] == "don't()":
			instruction.operation = memoryDisable
		case strings.HasPrefix(match[0], "mul("):
			instruction.a, err = convert.ScanNumber[int64]([]byte(match[1]))
			if err != nil {
				return nil, err
			}
			instruction.b, err = convert.ScanNumber[int64]([]byte(match[2]))
			if err != nil {
				return nil, err
			}
		}
		instructions = append(instructions, instruction)
	}
	return instructions, nil
}

func day03p01(r io.Reader) (string, error) {
	instructions, err := parseCorruptedMemory(r)
	if err != nil {
		return "", err
	}
	total := xiter.FromSlice(instructions).Filter(func(instruction memoryInstruction) bool {
		return instruction.operation == memoryMultiply
	}).SumBy(func(instruction memoryInstruction) int64 {
		return instruction.a * instruction.b
	})
	return strconv.FormatInt(total, 10), nil
}

func day03p02(r io.Reader) (string, error) {
	instructions, err := parseCorruptedMemory(r)
	if err != nil {
		return "", err
	}
	var total int64
	enabled := true
	for _, instruction := range instructions {
		switch instruction.operation {
		case memoryMultiply:
			if enabled {
				total += instruction.a * instruction.b
			}
		case memoryEnable:
			enabled = true
		case memoryDisable:
			enabled = false
		}
	}
	return strconv.FormatInt(total, 10), nil
}
