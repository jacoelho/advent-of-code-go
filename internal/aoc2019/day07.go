package aoc2019

import (
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

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

func findMaxSignal(program []int, phaseSettings []int, runFunc func([]int, []int) (int, error)) (int, error) {
	maxSignal := math.MinInt
	for perm := range xiter.Permutations(phaseSettings) {
		signal, err := runFunc(program, perm)
		if err != nil {
			return 0, err
		}
		maxSignal = max(maxSignal, signal)
	}
	return maxSignal, nil
}

func day07p01(r io.Reader) (string, error) {
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

		if xiter.FromSlice(amps).Every(func(amp *IntcodeComputer) bool { return amp.IsHalted() }) {
			break
		}

		ampIndex = (ampIndex + 1) % 5
	}

	return signal, nil
}

func day07p02(r io.Reader) (string, error) {
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
