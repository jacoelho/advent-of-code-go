package aoc2019

import "io"

func day09p01(r io.Reader) (string, error) {
	return runDiagnostic(r, 1)
}

func day09p02(r io.Reader) (string, error) {
	return runDiagnostic(r, 2)
}
