package aoc2019

import (
	"errors"
	"io"
	"math"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
)

func TestParseIntcodeInput(t *testing.T) {
	for _, input := range []string{"99\n", "99\r\n", " 99 "} {
		got, err := ParseIntcodeInput(strings.NewReader(input))
		if err != nil || !slices.Equal(got, []int{99}) {
			t.Errorf("ParseIntcodeInput(%q) = %v, %v; want [99]", input, got, err)
		}
	}
	got, err := ParseIntcodeInput(strings.NewReader(" 104 , +7 , 99\r\n"))
	if err != nil || !slices.Equal(got, []int{104, 7, 99}) {
		t.Fatalf("spaced input = %v, %v; want [104 7 99]", got, err)
	}
	for _, input := range []string{"", " ", ",99", "99,", "99,,1", "1 2", "1,x", "9223372036854775808"} {
		if _, err := ParseIntcodeInput(strings.NewReader(input)); err == nil {
			t.Errorf("ParseIntcodeInput(%q) accepted invalid input", input)
		}
	}
}

func TestIntcodeMemory(t *testing.T) {
	tests := []struct {
		name    string
		program []int
		output  []int
	}{
		{
			name:    "unwritten instruction parameters",
			program: []int{1101, 99, 0, 0, 1101, 1005, 0, 20, 1105, 1, 20},
		},
		{
			name:    "unwritten data beyond allocation limit",
			program: []int{4, 1_000_000, 99},
			output:  []int{0},
		},
		{
			name:    "large output",
			program: []int{104, 1125899906842624, 99},
			output:  []int{1125899906842624},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			computer := New(tt.program)
			if err := computer.Run(); err != nil {
				t.Fatal(err)
			}
			if got := computer.GetOutput(); !computer.IsHalted() || !slices.Equal(got, tt.output) {
				t.Fatalf("halted = %v, output = %v; want %v", computer.IsHalted(), got, tt.output)
			}
		})
	}

	program := []int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99}
	computer := New(program)
	if err := computer.Run(); err != nil {
		t.Fatal(err)
	}
	if got := computer.GetOutput(); !slices.Equal(got, program) {
		t.Fatalf("quine output = %v; want %v", got, program)
	}
}

func TestIntcodeAddressErrors(t *testing.T) {
	programs := [][]int{
		{109, -1, 204, 0, 99},
		{3, -1, 99},
		{1105, 1, -1},
		{1101, 1, 0, 1_000_000, 99},
		{109, math.MaxInt, 204, 1, 99},
		{109, math.MinInt, 109, -1, 99},
	}
	for _, program := range programs {
		computer := New(program)
		computer.SetInput(7)
		if err := computer.Run(); err == nil {
			t.Errorf("expected address error for %v", program)
		}
	}
}

func TestIntcodeInputOwnership(t *testing.T) {
	computer := New([]int{3, 100, 4, 100, 99})
	input := []int{7, 8}
	computer.SetInput(input...)
	input[0] = 42
	if err := computer.Run(); err != nil {
		t.Fatal(err)
	}
	if got := computer.GetOutput(); !slices.Equal(got, []int{7}) {
		t.Fatalf("output = %v; want [7]", got)
	}
}

func TestIntcodeCloneContinuation(t *testing.T) {
	computer := New([]int{1101, 0, 0, 40, 109, 30, 203, 0, 204, 0, 203, 1, 22201, 0, 1, 2, 204, 2, 99})
	computer.SetInput(7)
	if err := computer.Run(); err != nil {
		t.Fatal(err)
	}
	if !computer.IsWaiting() || computer.IsHalted() {
		t.Fatal("expected an input request after the first output")
	}
	left, right := computer.Clone(), computer.Clone()
	cases := []struct {
		computer *IntcodeComputer
		input    int
		want     int
	}{
		{left, 3, 10},
		{right, 5, 12},
		{computer, 11, 18},
	}
	for _, tt := range cases {
		tt.computer.AddInput(tt.input)
		if err := tt.computer.Run(); err != nil {
			t.Fatal(err)
		}
		if got := tt.computer.GetOutput(); !tt.computer.IsHalted() || !slices.Equal(got, []int{7, tt.want}) {
			t.Fatalf("input %d: halted = %v, output = %v; want [7 %d]", tt.input, tt.computer.IsHalted(), got, tt.want)
		}
	}
	for _, tt := range cases {
		memory := tt.computer.Memory()
		if memory[31] != tt.input || memory[32] != tt.want {
			t.Fatalf("input %d: continuation memory = %v", tt.input, memory[31:33])
		}
	}
}

func TestIntcodeReaderError(t *testing.T) {
	failure := errors.New("input read failed")
	reader := io.MultiReader(strings.NewReader("99"), iotest.ErrReader(failure))
	if _, err := ParseIntcodeInput(reader); !errors.Is(err, failure) {
		t.Fatalf("got %v, want input failure", err)
	}
}
