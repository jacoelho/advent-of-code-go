package aoc2023

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type cubeSet struct {
	red, green, blue int
}

type cubeGame struct {
	id   int
	sets []cubeSet
}

func parseCubeSet(setStr string) (cubeSet, error) {
	set := cubeSet{}
	cubes := strings.SplitSeq(setStr, ", ")

	for cube := range cubes {
		parts := strings.SplitN(cube, " ", 2)
		if len(parts) != 2 {
			return cubeSet{}, fmt.Errorf("invalid cube format: %s", cube)
		}

		count, err := strconv.Atoi(parts[0])
		if err != nil {
			return cubeSet{}, fmt.Errorf("invalid count: %w", err)
		}

		switch parts[1] {
		case "red":
			set.red = count
		case "green":
			set.green = count
		case "blue":
			set.blue = count
		default:
			return cubeSet{}, fmt.Errorf("invalid color: %s", parts[1])
		}
	}

	return set, nil
}

func parseCubeGame(line []byte) (cubeGame, error) {
	parts := strings.SplitN(string(line), ": ", 2)
	if len(parts) != 2 {
		return cubeGame{}, fmt.Errorf("invalid game format: %s", line)
	}

	digits := convert.ExtractDigits[int](parts[0])
	if len(digits) == 0 {
		return cubeGame{}, fmt.Errorf("invalid game ID: no digits found in %s", parts[0])
	}
	id := digits[0]

	setStrs := strings.Split(parts[1], "; ")
	sets := make([]cubeSet, 0, len(setStrs))

	for _, setStr := range setStrs {
		set, err := parseCubeSet(setStr)
		if err != nil {
			return cubeGame{}, err
		}
		sets = append(sets, set)
	}

	return cubeGame{id: id, sets: sets}, nil
}

func (g cubeGame) isValid(maxRed, maxGreen, maxBlue int) bool {
	return xiter.FromSlice(g.sets).Every(func(set cubeSet) bool {
		return set.red <= maxRed && set.green <= maxGreen && set.blue <= maxBlue
	})
}

func (g cubeGame) minCubes() cubeSet {
	return xiter.FromSlice(g.sets[1:]).Reduce(g.sets[0], func(set cubeSet, acc cubeSet) cubeSet {
		return cubeSet{
			red:   max(set.red, acc.red),
			green: max(set.green, acc.green),
			blue:  max(set.blue, acc.blue),
		}
	})
}

func (c cubeSet) power() int {
	return c.red * c.green * c.blue
}

func day02p01(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseCubeGame)

	sum := xiter.From(s.Values()).Filter(func(g cubeGame) bool {
		return g.isValid(12, 13, 14)
	}).SumBy(func(g cubeGame) int { return g.id })

	return strconv.Itoa(sum), s.Err()
}

func day02p02(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseCubeGame)

	sum := xiter.From(s.Values()).SumBy(func(g cubeGame) int {
		return g.minCubes().power()
	})

	return strconv.Itoa(sum), s.Err()
}
