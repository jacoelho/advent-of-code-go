package aoc2017

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func knotProduct(r io.Reader, size int) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	var lengths []int
	input = bytes.TrimSpace(input)
	if len(input) > 0 {
		for token := range bytes.SplitSeq(input, []byte{','}) {
			length, err := convert.ScanNumber[int](bytes.TrimSpace(token))
			if err != nil {
				return "", err
			}
			if length < 0 || length > size {
				return "", fmt.Errorf("invalid knot length %d for list size %d", length, size)
			}
			lengths = append(lengths, length)
		}
	}

	values := make([]byte, size)
	for i := range values {
		values[i] = byte(i)
	}

	position, skip := 0, 0
	knotRound(values, lengths, &position, &skip)

	return strconv.Itoa(int(values[0]) * int(values[1])), nil
}

func day10p01(r io.Reader) (string, error) {
	return knotProduct(r, 256)
}

func day10p02(r io.Reader) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	input = bytes.TrimSuffix(input, []byte{'\n'})
	input = bytes.TrimSuffix(input, []byte{'\r'})
	hash := knotHash(input)

	return hex.EncodeToString(hash[:]), nil
}
