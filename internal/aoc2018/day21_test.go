package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const haltingExample = `#ip 2
seti 123 0 4
bani 4 456 4
eqri 4 72 4
addr 4 2 2
seti 0 0 2
seti 0 7 4
bori 4 4 3
seti 1 1 4
bani 3 1 1
addr 4 1 4
bani 4 7 4
muli 4 3 4
bani 4 7 4
gtir 2 3 1
addr 1 2 2
addi 2 1 2
seti 27 8 2
seti 0 1 1
addi 1 1 5
muli 5 2 5
gtrr 5 3 5
addr 5 2 2
addi 2 1 2
seti 25 3 2
addi 1 1 1
seti 17 0 2
setr 1 4 3
seti 7 6 2
eqrr 4 0 1
addr 1 2 2
seti 5 2 2`

func Test_day21p01(t *testing.T) {
	zeroComparison := strings.NewReplacer("seti 0 7 4", "seti 5 7 4", "seti 1 1 4", "seti 0 1 4", "bani 3 1 1", "bani 3 3 1").Replace(haltingExample)
	tests := []aoc.TestInput{
		{Input: strings.NewReader(haltingExample), Want: "6"},
		{Input: strings.NewReader(zeroComparison), Want: "0"},
		{Input: strings.NewReader(strings.Replace(haltingExample, "eqri 4 72 4", "eqri 4 0 4", 1)), Want: "6"},
		{Input: aoc.FileInput(t, 2018, 21), Want: "13443200"},
	}

	aoc.AOCTest(t, day21p01, tests)
}

func Test_day21p02(t *testing.T) {
	zeroComparison := strings.NewReplacer("seti 0 7 4", "seti 5 7 4", "seti 1 1 4", "seti 0 1 4", "bani 3 1 1", "bani 3 3 1").Replace(haltingExample)
	tests := []aoc.TestInput{
		{Input: strings.NewReader(haltingExample), Want: "2"},
		{Input: strings.NewReader(zeroComparison), Want: "5"},
		{Input: strings.NewReader(strings.Replace(haltingExample, "eqri 4 72 4", "eqri 4 0 4", 1)), Want: "2"},
		{Input: aoc.FileInput(t, 2018, 21), Want: "7717135"},
	}

	aoc.AOCTest(t, day21p02, tests)
}

func TestHaltingProgramErrors(t *testing.T) {
	tests := map[string]string{
		"changed division jump":   strings.Replace(haltingExample, "seti 17 0 2", "seti 18 0 2", 1),
		"multiplication overflow": strings.Replace(haltingExample, "muli 4 3 4", "muli 4 9223372036854775807 4", 1),
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := day21p02(strings.NewReader(input)); err == nil {
				t.Fatal("expected unsupported or overflowing halting program error")
			}
		})
	}
}

func TestDay21ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"21/1", day21p01},
		{"21/2", day21p02},
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
