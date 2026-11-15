package aoc2023

import (
	"fmt"
	"io"
	"strconv"
	"strings"

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
	seen := make(map[string]bool, 3)
	for cube := range strings.SplitSeq(setStr, ", ") {
		countText, colour, ok := strings.Cut(cube, " ")
		if !ok {
			return cubeSet{}, fmt.Errorf("invalid cube format %q", cube)
		}
		count, err := strconv.Atoi(countText)
		if err != nil || count < 0 {
			return cubeSet{}, fmt.Errorf("invalid cube count %q", countText)
		}
		var field *int
		switch colour {
		case "red":
			field = &set.red
		case "green":
			field = &set.green
		case "blue":
			field = &set.blue
		default:
			return cubeSet{}, fmt.Errorf("invalid colour %q", colour)
		}
		if seen[colour] {
			return cubeSet{}, fmt.Errorf("duplicate colour %q in %q", colour, setStr)
		}
		seen[colour] = true
		*field = count
	}
	return set, nil
}

func parseCubeGame(line []byte) (cubeGame, error) {
	header, body, ok := strings.Cut(string(line), ": ")
	idText, hasPrefix := strings.CutPrefix(header, "Game ")
	if !ok || !hasPrefix {
		return cubeGame{}, fmt.Errorf("invalid game format %q", line)
	}
	id, err := strconv.Atoi(idText)
	if err != nil || id < 0 {
		return cubeGame{}, fmt.Errorf("invalid game ID %q", idText)
	}

	setStrs := strings.Split(body, "; ")
	sets := make([]cubeSet, 0, len(setStrs))
	for _, setStr := range setStrs {
		set, err := parseCubeSet(setStr)
		if err != nil {
			return cubeGame{}, fmt.Errorf("game %d: %w", id, err)
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
