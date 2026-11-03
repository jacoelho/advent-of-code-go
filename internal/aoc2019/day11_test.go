package aoc2019

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("99"), Want: "0"},
		{Input: strings.NewReader(paintingIProgram()), Want: "18"},
		{
			Input: aoc.FileInput(t, 2019, 11),
			Want:  "2041",
		},
	}
	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(paintingIProgram()), Want: "I"},
		{
			Input: aoc.FileInput(t, 2019, 11),
			Want:  "ZRZPKEZR",
		},
	}
	aoc.AOCTest(t, day11p02, tests)
}

func paintingIProgram() string {
	const colors = "0100110110011011101"
	const turns = "1011011001101000110"
	var program strings.Builder
	for i := range len(colors) {
		fmt.Fprintf(&program, "3,10000,104,%c,104,%c,", colors[i], turns[i])
	}
	program.WriteString("99")
	return program.String()
}

func TestPaintingProtocolErrors(t *testing.T) {
	programs := []string{
		"42",
		"3,100,104,1,99",
		"3,100,104,2,104,0,99",
		"3,100,104,1,104,2,99",
		"3,100,104,0,104,0,104,1,104,1,99",
		"3,100,3,101,99",
	}
	for _, solve := range []func(io.Reader) (string, error){day11p01, day11p02} {
		for _, program := range programs {
			if _, err := solve(strings.NewReader(program)); err == nil {
				t.Errorf("expected painting error for %s", program)
			}
		}
	}
}
