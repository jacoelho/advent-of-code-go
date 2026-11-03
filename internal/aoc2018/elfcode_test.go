package aoc2018

import (
	"math"
	"strings"
	"testing"
)

func TestElfOperations(t *testing.T) {
	tests := []struct {
		operation string
		a, b      int
		want      int64
	}{
		{operation: "addr", a: 0, b: 1, want: 5},
		{operation: "addi", a: 0, b: 7, want: 10},
		{operation: "mulr", a: 0, b: 1, want: 6},
		{operation: "muli", a: 0, b: 4, want: 12},
		{operation: "banr", a: 0, b: 1, want: 2},
		{operation: "bani", a: 0, b: 1, want: 1},
		{operation: "borr", a: 0, b: 1, want: 3},
		{operation: "bori", a: 1, b: 4, want: 6},
		{operation: "setr", a: 0, b: 99, want: 3},
		{operation: "seti", a: 7, b: 99, want: 7},
		{operation: "gtir", a: 4, b: 0, want: 1},
		{operation: "gtri", a: 0, b: 2, want: 1},
		{operation: "gtrr", a: 0, b: 1, want: 1},
		{operation: "eqir", a: 3, b: 0, want: 1},
		{operation: "eqri", a: 0, b: 4, want: 0},
		{operation: "eqrr", a: 0, b: 1, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.operation, func(t *testing.T) {
			instruction := elfInstruction{operation: tt.operation, a: tt.a, b: tt.b, c: 3}
			if err := instruction.validate(4); err != nil {
				t.Fatal(err)
			}
			registers := [4]int64{3, 2, 1, 9}
			if err := instruction.execute(registers[:]); err != nil {
				t.Fatal(err)
			}
			want := [4]int64{3, 2, 1, tt.want}
			if registers != want {
				t.Errorf("got %v, want %v", registers, want)
			}
		})
	}
}

func TestElfArithmeticOverflow(t *testing.T) {
	tests := []struct {
		instruction elfInstruction
		before      [4]int64
	}{
		{instruction: elfInstruction{operation: "addi", a: 0, b: 1, c: 1}, before: [4]int64{math.MaxInt64}},
		{instruction: elfInstruction{operation: "muli", a: 0, b: -1, c: 1}, before: [4]int64{math.MinInt64}},
	}

	for _, tt := range tests {
		t.Run(tt.instruction.operation, func(t *testing.T) {
			registers := tt.before
			if err := tt.instruction.execute(registers[:]); err == nil {
				t.Fatal("expected arithmetic overflow")
			}
			if registers != tt.before {
				t.Errorf("rejected instruction changed registers: %v", registers)
			}
		})
	}
}

func TestElfExecutionLimit(t *testing.T) {
	program, err := parseElfProgram(strings.NewReader("#ip 5\nseti -1 0 5"))
	if err != nil {
		t.Fatal(err)
	}
	machine := elfMachine{program: program}
	if err := machine.run(3); err == nil {
		t.Fatal("expected execution limit error")
	}
}

func TestElfInstructionPointerOverflow(t *testing.T) {
	if _, err := day19p01(strings.NewReader("#ip 5\nseti 9223372036854775807 0 5")); err == nil {
		t.Fatal("expected instruction pointer overflow")
	}
}
