package aoc2023

import (
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
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

func initializationSequence(r io.Reader) (string, error) {
	const inputLimit = 4 << 20
	data, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return "", err
	}
	if len(data) > inputLimit {
		return "", fmt.Errorf("initialization sequence exceeds 4 MiB")
	}
	sequence := strings.TrimSuffix(string(data), "\n")
	sequence = strings.TrimSuffix(sequence, "\r")
	if strings.ContainsAny(sequence, "\r\n") {
		return "", fmt.Errorf("initialization sequence must be one line")
	}
	return sequence, nil
}

func day15p01(r io.Reader) (string, error) {
	sequence, err := initializationSequence(r)
	if err != nil {
		return "", err
	}
	sum := xiter.From(strings.SplitSeq(sequence, ",")).SumBy(hash)
	return strconv.Itoa(sum), nil
}

type lens struct {
	label       string
	focalLength int64
}

func day15p02(r io.Reader) (string, error) {
	sequence, err := initializationSequence(r)
	if err != nil {
		return "", err
	}
	if sequence == "" {
		return "0", nil
	}
	boxes := make([][]lens, 256)
	for step := range strings.SplitSeq(sequence, ",") {
		label, value, assignment := strings.Cut(step, "=")
		if !assignment {
			var ok bool
			label, ok = strings.CutSuffix(step, "-")
			if !ok {
				return "", fmt.Errorf("invalid lens step %q", step)
			}
		}
		if label == "" {
			return "", fmt.Errorf("empty lens label")
		}
		boxNum := hash(label)
		if !assignment {
			boxes[boxNum] = slices.DeleteFunc(boxes[boxNum], func(l lens) bool {
				return l.label == label
			})
			continue
		}
		focalLength, err := convert.ScanNumber[int64]([]byte(value))
		if err != nil {
			return "", fmt.Errorf("lens %q focal length: %w", label, err)
		}
		if focalLength < 0 {
			return "", fmt.Errorf("negative focal length for lens %q", label)
		}
		index := slices.IndexFunc(boxes[boxNum], func(l lens) bool {
			return l.label == label
		})
		if index != -1 {
			boxes[boxNum][index].focalLength = focalLength
		} else {
			boxes[boxNum] = append(boxes[boxNum], lens{label: label, focalLength: focalLength})
		}
	}

	var totalPower int64
	for boxNum, box := range boxes {
		for slotNum, l := range box {
			factor := int64(boxNum+1) * int64(slotNum+1)
			if l.focalLength > math.MaxInt64/factor {
				return "", fmt.Errorf("lens focusing power overflows int64")
			}
			power := factor * l.focalLength
			if totalPower > math.MaxInt64-power {
				return "", fmt.Errorf("total focusing power overflows int64")
			}
			totalPower += power
		}
	}
	return strconv.FormatInt(totalPower, 10), nil
}
