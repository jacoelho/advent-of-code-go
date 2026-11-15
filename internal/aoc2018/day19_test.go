package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const boundIPExample = `#ip 0
seti 5 0 1
seti 6 0 2
addi 0 1 0
addr 1 2 3
setr 1 0 0
seti 8 0 4
seti 9 0 5`

const divisorExample = `#ip 3
seti 17 0 3
seti 1 0 5
seti 1 0 2
mulr 5 2 1
eqrr 1 4 1
addr 1 3 3
addi 3 1 3
addr 5 0 0
addi 2 1 2
gtrr 2 4 1
addr 3 1 3
seti 2 0 3
addi 5 1 5
gtrr 5 4 1
addr 1 3 3
seti 1 0 3
addi 0 10 0
seti 99 0 3
seti 6 0 4
seti 0 0 3`

func Test_day19p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(boundIPExample), Want: "6"},
		{Input: strings.NewReader(divisorExample), Want: "22"},
		{Input: strings.NewReader(strings.Replace(divisorExample, "addr 5 0 0", "addi 0 1 0", 1)), Want: "14"},
		{Input: aoc.FileInput(t, 2018, 19), Want: "1056"},
	}

	aoc.AOCTest(t, day19p01, tests)
}

func Test_day19p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(boundIPExample), Want: "6"},
		{Input: strings.NewReader(divisorExample), Want: "23"},
		{Input: strings.NewReader(strings.Replace(divisorExample, "addr 5 0 0", "addi 0 1 0", 1)), Want: "15"},
		{Input: aoc.FileInput(t, 2018, 19), Want: "10915260"},
	}

	aoc.AOCTest(t, day19p02, tests)
}

func TestDivisorLoopExit(t *testing.T) {
	program, err := parseElfProgram(strings.NewReader(divisorExample))
	if err != nil {
		t.Fatal(err)
	}
	roles, recognized := recognizeDivisorLoop(program)
	if !recognized {
		t.Fatal("expected divisor loop")
	}
	machine := elfMachine{program: program, registers: [6]int64{0, 0, 0, 0, 6, 0}, ip: 1}
	if !machine.skipDivisorLoop(roles) {
		t.Fatal("expected loop acceleration")
	}
	want := [6]int64{12, 1, 7, 15, 6, 7}
	if machine.registers != want || machine.ip != 16 {
		t.Errorf("exit = %v at %d, want %v at 16", machine.registers, machine.ip, want)
	}
}

func TestDay19ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"19/1", day19p01},
		{"19/2", day19p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
