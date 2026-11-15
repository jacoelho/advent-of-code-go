package aoc2018

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type recipeBoard struct {
	scores        []byte
	first, second int
}

func (board *recipeBoard) step() {
	a, b := board.scores[board.first], board.scores[board.second]
	if a+b >= 10 {
		board.scores = append(board.scores, 1)
	}
	board.scores = append(board.scores, (a+b)%10)
	board.first = (board.first + int(a) + 1) % len(board.scores)
	board.second = (board.second + int(b) + 1) % len(board.scores)
}

func recipeInput(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, fmt.Errorf("empty recipe input")
	}
	for _, digit := range data {
		if digit < '0' || digit > '9' {
			return nil, fmt.Errorf("invalid recipe digit %q", digit)
		}
	}
	return data, nil
}

func day14p01(r io.Reader) (string, error) {
	data, err := recipeInput(r)
	if err != nil {
		return "", err
	}
	target, err := convert.ScanNumber[int](data)
	if err != nil {
		return "", fmt.Errorf("invalid recipe count: %w", err)
	}
	end, ok := xmath.CheckedAdd(target, 10)
	if !ok {
		return "", fmt.Errorf("recipe count %d overflows the scoreboard index", target)
	}

	board := recipeBoard{scores: []byte{3, 7}, second: 1}
	for len(board.scores) < end {
		board.step()
	}
	result := make([]byte, 10)
	for i, score := range board.scores[target:end] {
		result[i] = score + '0'
	}
	return string(result), nil
}

func day14p02(r io.Reader) (string, error) {
	pattern, err := recipeInput(r)
	if err != nil {
		return "", err
	}
	for i := range pattern {
		pattern[i] -= '0'
	}
	prefix := make([]int, len(pattern))
	for i, matched := 1, 0; i < len(pattern); i++ {
		for matched > 0 && pattern[i] != pattern[matched] {
			matched = prefix[matched-1]
		}
		if pattern[i] == pattern[matched] {
			matched++
		}
		prefix[i] = matched
	}

	board := recipeBoard{scores: []byte{3, 7}, second: 1}
	matched, checked := 0, 0
	for {
		for checked < len(board.scores) {
			digit := board.scores[checked]
			for matched > 0 && digit != pattern[matched] {
				matched = prefix[matched-1]
			}
			if digit == pattern[matched] {
				matched++
			}
			checked++
			if matched == len(pattern) {
				return strconv.Itoa(checked - len(pattern)), nil
			}
		}
		board.step()
	}
}
