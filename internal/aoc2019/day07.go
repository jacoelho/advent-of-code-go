package aoc2019

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

// runAmplifiers feeds the signal through the amplifiers in order, looping back
// to the first until the last amplifier halts. Each amplifier must answer
// every signal with exactly one output.
func runAmplifiers(program []int, phases []int) (int, error) {
	amps := make([]*IntcodeComputer, len(phases))
	for i, phase := range phases {
		amps[i] = New(program)
		amps[i].SetInput(phase)
	}

	signal := 0
	for {
		for i, amp := range amps {
			amp.AddInput(signal)
			if err := amp.Run(); err != nil {
				return 0, fmt.Errorf("amplifier %d: %w", i, err)
			}
			output := amp.GetOutput()
			if len(output) != 1 {
				return 0, fmt.Errorf("amplifier %d produced %d outputs for one signal", i, len(output))
			}
			signal = output[0]
		}
		if amps[len(amps)-1].IsHalted() {
			return signal, nil
		}
	}
}

func maxThrusterSignal(r io.Reader, phases []int) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	best := math.MinInt
	for perm := range xiter.Permutations(phases) {
		signal, err := runAmplifiers(program, perm)
		if err != nil {
			return "", err
		}
		best = max(best, signal)
	}
	return strconv.Itoa(best), nil
}

func day07p01(r io.Reader) (string, error) {
	return maxThrusterSignal(r, []int{0, 1, 2, 3, 4})
}

func day07p02(r io.Reader) (string, error) {
	return maxThrusterSignal(r, []int{5, 6, 7, 8, 9})
}
