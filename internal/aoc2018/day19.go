package aoc2018

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type divisorRegisters struct {
	outer, inner, scratch, target int
}

func recognizeDivisorLoop(program elfProgram) (divisorRegisters, bool) {
	var roles divisorRegisters
	if len(program.code) < 16 {
		return roles, false
	}
	code, ip := program.code, program.ipRegister
	roles = divisorRegisters{outer: code[1].c, inner: code[2].c, scratch: code[3].c, target: code[4].b}
	registers := map[int]bool{0: true, ip: true, roles.outer: true, roles.inner: true, roles.scratch: true, roles.target: true}
	if len(registers) != 6 {
		return roles, false
	}
	want := []elfInstruction{
		{operation: "seti", a: 1, b: code[1].b, c: roles.outer},
		{operation: "seti", a: 1, b: code[2].b, c: roles.inner},
		{operation: "mulr", a: roles.outer, b: roles.inner, c: roles.scratch},
		{operation: "eqrr", a: roles.scratch, b: roles.target, c: roles.scratch},
		{operation: "addr", a: roles.scratch, b: ip, c: ip},
		{operation: "addi", a: ip, b: 1, c: ip},
		{operation: "addr", a: roles.outer, b: 0, c: 0},
		{operation: "addi", a: roles.inner, b: 1, c: roles.inner},
		{operation: "gtrr", a: roles.inner, b: roles.target, c: roles.scratch},
		{operation: "addr", a: ip, b: roles.scratch, c: ip},
		{operation: "seti", a: 2, b: code[11].b, c: ip},
		{operation: "addi", a: roles.outer, b: 1, c: roles.outer},
		{operation: "gtrr", a: roles.outer, b: roles.target, c: roles.scratch},
		{operation: "addr", a: roles.scratch, b: ip, c: ip},
		{operation: "seti", a: 1, b: code[15].b, c: ip},
	}
	for i, instruction := range want {
		if code[i+1] != instruction {
			return roles, false
		}
	}
	return roles, true
}

func (machine *elfMachine) skipDivisorLoop(roles divisorRegisters) bool {
	target := machine.registers[roles.target]
	if _, ok := xmath.CheckedMul(target, target); target <= 0 || !ok {
		return false
	}
	var sum int64
	for divisor := int64(1); divisor <= target/divisor; divisor++ {
		if target%divisor != 0 {
			continue
		}
		sum += divisor
		if other := target / divisor; other != divisor {
			sum += other
		}
	}
	total, ok := xmath.CheckedAdd(machine.registers[0], sum)
	if !ok {
		return false
	}
	machine.registers[0] = total
	machine.registers[roles.outer], machine.registers[roles.inner] = target+1, target+1
	machine.registers[roles.scratch] = 1
	// Resume at the loop's exit so the supplied program still owns its final result.
	machine.registers[machine.program.ipRegister] = 15
	machine.ip = 16
	return true
}

func runDivisorProgram(program elfProgram, initial int64) (int64, error) {
	roles, recognized := recognizeDivisorLoop(program)
	machine := elfMachine{program: program, registers: [6]int64{initial}}
	for machine.running() {
		if recognized && machine.ip == 1 && machine.skipDivisorLoop(roles) {
			continue
		}
		if err := machine.step(); err != nil {
			return 0, err
		}
	}
	return machine.registers[0], nil
}

func divisorProgramResult(r io.Reader, initial int64) (string, error) {
	program, err := parseElfProgram(r)
	if err != nil {
		return "", err
	}
	result, err := runDivisorProgram(program, initial)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(result, 10), nil
}

func day19p01(r io.Reader) (string, error) {
	return divisorProgramResult(r, 0)
}

func day19p02(r io.Reader) (string, error) {
	return divisorProgramResult(r, 1)
}
