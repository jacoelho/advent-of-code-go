package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type bingoBoard struct {
	values [25]int64
	marked uint32
	won    bool
}

func bingoInput(r io.Reader) ([]int64, []bingoBoard, error) {
	s := bufio.NewScanner(r)
	if !s.Scan() {
		if err := s.Err(); err != nil {
			return nil, nil, err
		}
		return nil, nil, fmt.Errorf("empty bingo input")
	}
	var draws []int64
	for field := range strings.SplitSeq(s.Text(), ",") {
		value, err := convert.ScanNumber[int64]([]byte(strings.TrimSpace(field)))
		if err != nil {
			return nil, nil, err
		}
		draws = append(draws, value)
	}
	var boards []bingoBoard
	var board bingoBoard
	row := 0
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) == 0 {
			continue
		}
		if len(fields) != 5 {
			return nil, nil, fmt.Errorf("bingo row must contain five numbers")
		}
		for column, field := range fields {
			value, err := convert.ScanNumber[int64]([]byte(field))
			if err != nil {
				return nil, nil, err
			}
			board.values[row*5+column] = value
		}
		row++
		if row == 5 {
			boards = append(boards, board)
			board, row = bingoBoard{}, 0
		}
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	if row != 0 || len(boards) == 0 {
		return nil, nil, fmt.Errorf("missing or incomplete bingo board")
	}
	return draws, boards, nil
}

func bingoScores(r io.Reader) ([]int64, error) {
	draws, boards, err := bingoInput(r)
	if err != nil {
		return nil, err
	}
	var scores []int64
	for _, draw := range draws {
		for i := range boards {
			board := &boards[i]
			if board.won {
				continue
			}
			for cell, value := range board.values {
				if value == draw {
					board.marked |= 1 << cell
				}
			}
			for line := range 5 {
				row, column := uint32(31)<<(line*5), uint32(0x108421)<<line
				if board.marked&row == row || board.marked&column == column {
					board.won = true
				}
			}
			if board.won {
				var unmarked int64
				for cell, value := range board.values {
					if board.marked&(1<<cell) == 0 {
						unmarked += value
					}
				}
				scores = append(scores, unmarked*draw)
			}
		}
	}
	if len(scores) == 0 {
		return nil, fmt.Errorf("no bingo board wins")
	}
	return scores, nil
}

func day04p01(r io.Reader) (string, error) {
	scores, err := bingoScores(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(scores[0], 10), nil
}

func day04p02(r io.Reader) (string, error) {
	scores, err := bingoScores(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(scores[len(scores)-1], 10), nil
}
