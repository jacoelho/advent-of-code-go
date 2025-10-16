package aoc2023

import (
	"bufio"
	"io"
	"slices"
	"strconv"
	"strings"

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

func day15p01(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	line := scanner.Text()

	steps := strings.Split(line, ",")
	sum := xiter.Sum(xiter.Map(hash, slices.Values(steps)))

	return strconv.Itoa(sum), nil
}

type lens struct {
	label       string
	focalLength int
}

func day15p02(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	line := scanner.Text()

	boxes := make([][]lens, 256)

	steps := strings.SplitSeq(line, ",")
	for step := range steps {
		if strings.Contains(step, "=") {
			parts := strings.Split(step, "=")
			label := parts[0]
			focalLength, _ := strconv.Atoi(parts[1])
			boxNum := hash(label)

			idx := slices.IndexFunc(boxes[boxNum], func(l lens) bool {
				return l.label == label
			})
			if idx != -1 {
				boxes[boxNum][idx].focalLength = focalLength
			} else {
				boxes[boxNum] = append(boxes[boxNum], lens{label, focalLength})
			}
		} else {
			label := strings.TrimSuffix(step, "-")
			boxNum := hash(label)

			boxes[boxNum] = slices.DeleteFunc(boxes[boxNum], func(l lens) bool {
				return l.label == label
			})
		}
	}

	totalPower := 0
	for boxNum, box := range boxes {
		for slotNum, l := range box {
			totalPower += (boxNum + 1) * (slotNum + 1) * l.focalLength
		}
	}

	return strconv.Itoa(totalPower), nil
}
