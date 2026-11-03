package aoc2021

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func monadFixture(divisor, offset, addend int) string {
	return fmt.Sprintf(`inp w
mul x 0
add x z
mod x 26
div z %d
add x %d
eql x w
eql x 0
mul y 0
add y 25
mul y x
add y 1
mul z y
mul y 0
add y w
add y %d
mul y x
add z y
`, divisor, offset, addend)
}

func Test_day24p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(monadFixture(1, 12, 2) + monadFixture(26, -1, 0)), Want: "89"},
		{Input: strings.NewReader(monadFixture(1, 12, 2) + monadFixture(1, 11, 4) + monadFixture(26, -6, 5) + monadFixture(26, -1, 0)), Want: "8979"},
		{Input: aoc.FileInput(t, 2021, 24), Want: "99196997985942"},
	}
	aoc.AOCTest(t, day24p01, tests)
}

func Test_day24p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(monadFixture(1, 12, 2) + monadFixture(26, -1, 0)), Want: "12"},
		{Input: strings.NewReader(monadFixture(1, 12, 2) + monadFixture(1, 11, 4) + monadFixture(26, -6, 5) + monadFixture(26, -1, 0)), Want: "1312"},
		{Input: aoc.FileInput(t, 2021, 24), Want: "84191521311611"},
	}
	aoc.AOCTest(t, day24p02, tests)
}

func TestALUSemantics(t *testing.T) {
	tests := []struct {
		program string
		input   []int64
		want    [4]int64
	}{
		{program: "inp x\nmul x -1", input: []int64{7}, want: [4]int64{0, -7, 0, 0}},
		{program: "inp z\ninp x\nmul z 3\neql z x", input: []int64{3, 9}, want: [4]int64{0, 9, 0, 1}},
		{program: "inp x\ndiv x 3", input: []int64{-7}, want: [4]int64{0, -2, 0, 0}},
		{program: "inp x\ndiv x -3", input: []int64{7}, want: [4]int64{0, -2, 0, 0}},
		{program: "inp x\ndiv x -3", input: []int64{-7}, want: [4]int64{0, 2, 0, 0}},
		{program: "inp x\nmod x 3", input: []int64{8}, want: [4]int64{0, 2, 0, 0}},
	}
	for _, tt := range tests {
		program, err := parseALU(strings.NewReader(tt.program))
		if err != nil {
			t.Fatal(err)
		}
		got, err := executeALU(program, tt.input)
		if err != nil || got != tt.want {
			t.Fatalf("%q with %v: got %v, %v; want %v", tt.program, tt.input, got, err, tt.want)
		}
	}
}

func TestALUFailures(t *testing.T) {
	tests := []struct {
		program string
		input   []int64
	}{
		{program: "inp x"},
		{program: "inp x", input: []int64{1, 2}},
		{program: "add x 1\ndiv x 0"},
		{program: "add x -1\nmod x 3"},
		{program: "mod x 0"},
		{program: "mod x -1"},
		{program: "add x 9223372036854775807\nadd x 1"},
		{program: "add x -9223372036854775808\nadd x -1"},
		{program: "add x 9223372036854775807\nmul x 2"},
		{program: "add x -9223372036854775808\nmul x -1"},
		{program: "add x -9223372036854775808\ndiv x -1"},
	}
	for _, tt := range tests {
		program, err := parseALU(strings.NewReader(tt.program))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := executeALU(program, tt.input); err == nil {
			t.Fatalf("expected ALU error for %q with %v", tt.program, tt.input)
		}
	}
}

func TestMONADRejectsUnsupportedPrograms(t *testing.T) {
	valid := monadFixture(1, 12, 2) + monadFixture(26, -1, 0)
	tests := []string{
		"",
		"inp a",
		"inp w 1",
		"add w",
		"add w x extra",
		"xyz w 1",
		"add w +",
		valid + "add z 0\n",
		strings.Replace(valid, "mul y x", "mul y w", 1),
		strings.Replace(valid, "div z 1", "div z x", 1),
		monadFixture(1, 9, 2) + monadFixture(26, -1, 0),
		monadFixture(1, 12, 17) + monadFixture(26, -1, 0),
		monadFixture(1, 12, 0) + monadFixture(26, 9, 0),
		monadFixture(1, 12, 0) + monadFixture(26, -9, 0),
		monadFixture(26, -1, 0) + monadFixture(1, 12, 2),
		monadFixture(1, 12, 2) + monadFixture(1, 12, 2),
	}
	for _, input := range tests {
		if _, err := day24p01(strings.NewReader(input)); err == nil {
			t.Fatalf("expected MONAD error for %q", input)
		}
		if _, err := day24p02(strings.NewReader(input)); err == nil {
			t.Fatalf("expected MONAD error for %q", input)
		}
	}
}
