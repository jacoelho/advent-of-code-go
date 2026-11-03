package aoc2018

import (
	"strings"
	"testing"

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
	remaining := 10
	skipped, err := machine.skipDivisorLoop(roles, &remaining)
	if err != nil {
		t.Fatal(err)
	}
	if !skipped {
		t.Fatal("expected loop acceleration")
	}
	want := [6]int64{12, 1, 7, 15, 6, 7}
	if machine.registers != want || machine.ip != 16 {
		t.Errorf("exit = %v at %d, want %v at 16", machine.registers, machine.ip, want)
	}
}

func TestDivisorLoopBudget(t *testing.T) {
	program, err := parseElfProgram(strings.NewReader(divisorExample))
	if err != nil {
		t.Fatal(err)
	}
	roles, recognized := recognizeDivisorLoop(program)
	if !recognized {
		t.Fatal("expected divisor loop")
	}
	before := [6]int64{0, 0, 0, 0, 6, 0}
	machine := elfMachine{program: program, registers: before, ip: 1}
	remaining := 1
	if skipped, err := machine.skipDivisorLoop(roles, &remaining); err == nil || skipped {
		t.Fatalf("skip = %v, error = %v; expected budget failure", skipped, err)
	}
	if machine.registers != before || machine.ip != 1 {
		t.Errorf("budget failure changed state: %v at %d", machine.registers, machine.ip)
	}
}

func TestDivisorProgramBudget(t *testing.T) {
	program, err := parseElfProgram(strings.NewReader(divisorExample))
	if err != nil {
		t.Fatal(err)
	}
	// Three initializer steps, one skip dispatch, two divisors, and two tail steps cost eight units.
	if _, err := runDivisorProgram(program, 0, 7); err == nil {
		t.Fatal("expected shared work budget exhaustion")
	}
	got, err := runDivisorProgram(program, 0, 8)
	if err != nil {
		t.Fatal(err)
	}
	if got != 22 {
		t.Errorf("got %d, want 22", got)
	}
}
