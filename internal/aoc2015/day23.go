package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type tinyInstruction struct {
	op               string
	register, offset int
}

func tinyComputer(r io.Reader, initialA uint64) ([2]uint64, error) {
	var program []tinyInstruction
	s := bufio.NewScanner(r)
	for s.Scan() {
		op, operands, ok := strings.Cut(s.Text(), " ")
		if !ok {
			return [2]uint64{}, fmt.Errorf("invalid instruction %q", s.Text())
		}
		args := strings.Split(operands, ", ")
		instruction := tinyInstruction{op: op}
		switch {
		case (op == "hlf" || op == "tpl" || op == "inc") && len(args) == 1,
			(op == "jie" || op == "jio") && len(args) == 2:
			if args[0] != "a" && args[0] != "b" {
				return [2]uint64{}, fmt.Errorf("unknown register in %q", s.Text())
			}
			instruction.register = int(args[0][0] - 'a')
		case op == "jmp" && len(args) == 1:
		default:
			return [2]uint64{}, fmt.Errorf("invalid instruction %q", s.Text())
		}
		if op == "jmp" || op == "jie" || op == "jio" {
			var err error
			instruction.offset, err = convert.ScanNumber[int]([]byte(args[len(args)-1]))
			if err != nil {
				return [2]uint64{}, fmt.Errorf("invalid jump offset in %q: %w", s.Text(), err)
			}
		}
		program = append(program, instruction)
	}
	if err := s.Err(); err != nil {
		return [2]uint64{}, err
	}
	registers := [2]uint64{initialA, 0}
	type state struct {
		ip        int
		registers [2]uint64
	}
	seen := make(map[state]struct{})
	for ip := 0; ip >= 0 && ip < len(program); {
		current := state{ip, registers}
		if _, ok := seen[current]; ok {
			return [2]uint64{}, fmt.Errorf("program loops at instruction %d", ip)
		}
		seen[current] = struct{}{}
		instruction := program[ip]
		value := &registers[instruction.register]
		offset := 1
		switch instruction.op {
		case "hlf":
			*value /= 2
		case "tpl":
			if *value > math.MaxUint64/3 {
				return [2]uint64{}, fmt.Errorf("register overflow at instruction %d", ip)
			}
			*value *= 3
		case "inc":
			if *value == math.MaxUint64 {
				return [2]uint64{}, fmt.Errorf("register overflow at instruction %d", ip)
			}
			*value++
		case "jmp":
			offset = instruction.offset
		case "jie":
			if *value%2 == 0 {
				offset = instruction.offset
			}
		case "jio":
			if *value == 1 {
				offset = instruction.offset
			}
		}
		if offset >= len(program)-ip || offset < -ip {
			break
		}
		ip += offset
	}
	return registers, nil
}

func day23p01(r io.Reader) (string, error) {
	registers, err := tinyComputer(r, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(registers[1], 10), nil
}

func day23p02(r io.Reader) (string, error) {
	registers, err := tinyComputer(r, 1)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(registers[1], 10), nil
}
