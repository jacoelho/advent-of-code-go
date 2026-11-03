package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func parseDiracPlayers(r io.Reader) ([2]int, error) {
	s := bufio.NewScanner(r)
	var positions [2]int
	count := 0
	for s.Scan() {
		if count == len(positions) {
			return positions, fmt.Errorf("expected exactly two Dirac players")
		}
		text, ok := strings.CutPrefix(s.Text(), fmt.Sprintf("Player %d starting position: ", count+1))
		if !ok {
			return positions, fmt.Errorf("invalid Dirac player %d", count+1)
		}
		position, err := convert.ScanNumber[int]([]byte(text))
		if err != nil || position < 1 || position > 10 {
			return positions, fmt.Errorf("dirac player %d position must be in 1..10", count+1)
		}
		positions[count] = position
		count++
	}
	if err := s.Err(); err != nil {
		return positions, err
	}
	if count != len(positions) {
		return positions, fmt.Errorf("expected exactly two Dirac players")
	}
	return positions, nil
}

func day21p01(r io.Reader) (string, error) {
	positions, err := parseDiracPlayers(r)
	if err != nil {
		return "", err
	}
	var scores [2]int
	die, rolls, turn := 1, 0, 0
	for {
		for range 3 {
			positions[turn] = (positions[turn]+die-1)%10 + 1
			die = die%100 + 1
			rolls++
		}
		scores[turn] += positions[turn]
		if scores[turn] >= 1000 {
			return strconv.Itoa(scores[1-turn] * rolls), nil
		}
		turn = 1 - turn
	}
}

func day21p02(r io.Reader) (string, error) {
	positions, err := parseDiracPlayers(r)
	if err != nil {
		return "", err
	}
	memo := new([10][10][21][21][2]uint64)
	var wins func(int, int, int, int) ([2]uint64, error)
	wins = func(position, otherPosition, score, otherScore int) ([2]uint64, error) {
		cached := &memo[position-1][otherPosition-1][score][otherScore]
		if *cached != [2]uint64{} {
			return *cached, nil
		}
		var result [2]uint64
		for i, frequency := range [...]uint64{1, 3, 6, 7, 6, 3, 1} {
			nextPosition := (position+i+2)%10 + 1
			nextScore := score + nextPosition
			if nextScore >= 21 {
				result[0] += frequency
				continue
			}
			next, err := wins(otherPosition, nextPosition, otherScore, nextScore)
			if err != nil {
				return result, err
			}
			for player := range 2 {
				if next[1-player] > (math.MaxUint64-result[player])/frequency {
					return result, fmt.Errorf("dirac universe count exceeds uint64")
				}
				result[player] += frequency * next[1-player]
			}
		}
		*cached = result
		return result, nil
	}
	result, err := wins(positions[0], positions[1], 0, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(max(result[0], result[1]), 10), nil
}
