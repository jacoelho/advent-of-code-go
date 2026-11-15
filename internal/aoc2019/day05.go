package aoc2019

import (
	"fmt"
	"io"
	"strconv"
)

// runDiagnostic runs a self-testing program with one input and returns its
// final output; every earlier output reports a failed check unless it is zero.
func runDiagnostic(r io.Reader, input int) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	computer := New(program)
	computer.SetInput(input)
	if err := computer.RunToHalt(); err != nil {
		return "", err
	}

	output := computer.GetOutput()
	if len(output) == 0 {
		return "", fmt.Errorf("no output produced")
	}
	for i, value := range output[:len(output)-1] {
		if value != 0 {
			return "", fmt.Errorf("diagnostic check %d failed with %d", i, value)
		}
	}
	return strconv.Itoa(output[len(output)-1]), nil
}

func day05p01(r io.Reader) (string, error) {
	return runDiagnostic(r, 1)
}

func day05p02(r io.Reader) (string, error) {
	return runDiagnostic(r, 5)
}
