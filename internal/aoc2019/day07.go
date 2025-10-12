package aoc2019

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/xiter"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

// runAmplifierChain runs a series of 5 amplifiers with the given phase settings
func runAmplifierChain(program []int, phaseSettings []int) (int, error) {
	signal := 0

	for _, phase := range phaseSettings {
		amp := New(program)
		amp.SetInput(phase, signal)
		if err := amp.Run(); err != nil {
			return 0, err
		}
		output, err := amp.LastOutput()
		if err != nil {
			return 0, err
		}
		signal = output
	}

	return signal, nil
}

// findMaxSignal finds the maximum signal by trying all permutations of phase settings
func findMaxSignal(program []int, phaseSettings []int, runFunc func([]int, []int) (int, error)) (int, error) {
	maxSignal := 0
	for perm := range xiter.Permutations(phaseSettings) {
		signal, err := runFunc(program, perm)
		if err != nil {
			return 0, err
		}
		maxSignal = max(maxSignal, signal)
	}
	return maxSignal, nil
}

func day7p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	maxSignal, err := findMaxSignal(program, []int{0, 1, 2, 3, 4}, runAmplifierChain)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(maxSignal), nil
}

// runAmplifierFeedbackLoop runs 5 amplifiers in feedback loop mode
func runAmplifierFeedbackLoop(program []int, phaseSettings []int) (int, error) {
	amps := make([]*IntcodeComputer, 5)
	for i := range 5 {
		amps[i] = New(program)
		amps[i].SetInput(phaseSettings[i])
	}

	signal := 0
	ampIndex := 0

	for {
		amps[ampIndex].AddInput(signal)

		if err := amps[ampIndex].Run(); err != nil {
			return 0, err
		}

		if output, err := amps[ampIndex].LastOutput(); err == nil {
			signal = output
		}

		if xslices.Every(func(amp *IntcodeComputer) bool { return amp.IsHalted() }, amps) {
			break
		}

		ampIndex = (ampIndex + 1) % 5
	}

	return signal, nil
}

func day7p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	maxSignal, err := findMaxSignal(program, []int{5, 6, 7, 8, 9}, runAmplifierFeedbackLoop)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(maxSignal), nil
}
