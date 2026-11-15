package aoc2023

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func hash(s string) int {
	current := 0
	for i := 0; i < len(s); i++ {
		current += int(s[i])
		current *= 17
		current %= 256
	}
	return current
}

type lensStep struct {
	text        string
	label       string
	remove      bool
	focalLength int64
}

func parseLensStep(step string) (lensStep, error) {
	result := lensStep{text: step}
	label, value, assignment := strings.Cut(step, "=")
	if !assignment {
		var ok bool
		label, ok = strings.CutSuffix(step, "-")
		if !ok {
			return lensStep{}, fmt.Errorf("invalid lens step %q", step)
		}
		result.remove = true
	}
	if label == "" || strings.ContainsAny(label, "=-,") {
		return lensStep{}, fmt.Errorf("invalid lens label in step %q", step)
	}
	result.label = label
	if assignment {
		focalLength, err := convert.ScanNumber[int64]([]byte(value))
		if err != nil {
			return lensStep{}, fmt.Errorf("lens %q focal length: %w", label, err)
		}
		if focalLength < 0 {
			return lensStep{}, fmt.Errorf("negative focal length for lens %q", label)
		}
		result.focalLength = focalLength
	}
	return result, nil
}

func initializationSequence(r io.Reader) ([]lensStep, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	sequence := strings.TrimSuffix(string(data), "\n")
	sequence = strings.TrimSuffix(sequence, "\r")
	if strings.ContainsAny(sequence, "\r\n") {
		return nil, fmt.Errorf("initialization sequence must be one line")
	}
	if sequence == "" {
		return nil, fmt.Errorf("empty initialization sequence")
	}
	var steps []lensStep
	for text := range strings.SplitSeq(sequence, ",") {
		step, err := parseLensStep(text)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	return steps, nil
}

func day15p01(r io.Reader) (string, error) {
	steps, err := initializationSequence(r)
	if err != nil {
		return "", err
	}
	sum := 0
	for _, step := range steps {
		sum += hash(step.text)
	}
	return strconv.Itoa(sum), nil
}

type lens struct {
	label       string
	focalLength int64
}

func day15p02(r io.Reader) (string, error) {
	steps, err := initializationSequence(r)
	if err != nil {
		return "", err
	}
	boxes := make([][]lens, 256)
	for _, step := range steps {
		boxNum := hash(step.label)
		index := slices.IndexFunc(boxes[boxNum], func(l lens) bool {
			return l.label == step.label
		})
		switch {
		case step.remove:
			if index >= 0 {
				boxes[boxNum] = slices.Delete(boxes[boxNum], index, index+1)
			}
		case index >= 0:
			boxes[boxNum][index].focalLength = step.focalLength
		default:
			boxes[boxNum] = append(boxes[boxNum], lens{label: step.label, focalLength: step.focalLength})
		}
	}

	var totalPower int64
	for boxNum, box := range boxes {
		for slotNum, l := range box {
			factor := int64(boxNum+1) * int64(slotNum+1)
			power, ok := xmath.CheckedMul(factor, l.focalLength)
			if !ok {
				return "", fmt.Errorf("lens focusing power overflows int64")
			}
			if totalPower, ok = xmath.CheckedAdd(totalPower, power); !ok {
				return "", fmt.Errorf("total focusing power overflows int64")
			}
		}
	}
	return strconv.FormatInt(totalPower, 10), nil
}
