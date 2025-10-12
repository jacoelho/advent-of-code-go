package aoc2019

import (
	"fmt"
	"io"
	"strings"
)

func runSpringdroid(program []int, springscript []string) (string, error) {
	computer := New(program)
	computer.SetInput(StringsToASCII(springscript...)...)

	if err := computer.Run(); err != nil {
		return "", err
	}

	output := computer.GetOutput()

	if len(output) > 0 && output[len(output)-1] > 127 {
		return fmt.Sprintf("%d", output[len(output)-1]), nil
	}

	var sb strings.Builder
	for _, val := range output {
		if val <= 127 {
			sb.WriteByte(byte(val))
		} else {
			sb.WriteString(fmt.Sprintf("[%d]", val))
		}
	}
	return "", fmt.Errorf("springdroid failed: %s", sb.String())
}

func day21p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	// jump if there's a hole in A, B, or C AND there's ground at D
	// logic: (!A OR !B OR !C) AND D
	springscript := []string{
		"NOT A J",
		"NOT B T",
		"OR T J",
		"NOT C T",
		"OR T J",
		"AND D J",
		"WALK",
	}

	return runSpringdroid(program, springscript)
}

func day21p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	// jump if there's a hole in A, B, or C AND there's ground at D
	// AND we have a safe continuation (ground at E OR ground at H)
	// logic: (!A OR !B OR !C) AND D AND (E OR H)
	springscript := []string{
		"NOT A J",
		"NOT B T",
		"OR T J",
		"NOT C T",
		"OR T J",
		"AND D J",
		"NOT E T",
		"NOT T T",
		"OR H T",
		"AND T J",
		"RUN",
	}

	return runSpringdroid(program, springscript)
}
