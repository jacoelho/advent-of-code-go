package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func spinlockStep(r io.Reader) (int64, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	step, err := convert.ScanNumber[int64](bytes.TrimSpace(data))
	if err != nil || step < 0 {
		return 0, fmt.Errorf("invalid spinlock step %q", data)
	}

	return step, nil
}

func spinlockAfterZero(step int64, insertions int) int {
	position, value := 0, 0
	for number := 1; number <= insertions; number++ {
		position = (position+int(step%int64(number)))%number + 1
		if position == 1 {
			value = number
		}
	}

	return value
}

func day17p01(r io.Reader) (string, error) {
	step, err := spinlockStep(r)
	if err != nil {
		return "", err
	}

	buffer := make([]int, 1, 2018)
	position := 0
	for number := 1; number <= 2017; number++ {
		position = (position+int(step%int64(number)))%number + 1
		buffer = slices.Insert(buffer, position, number)
	}

	return strconv.Itoa(buffer[(position+1)%len(buffer)]), nil
}

func day17p02(r io.Reader) (string, error) {
	step, err := spinlockStep(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(spinlockAfterZero(step, 50_000_000)), nil
}
