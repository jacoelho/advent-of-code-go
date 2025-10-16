package aoc2023

import (
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type race struct {
	time     int
	distance int
}

func parseRaces(input string) []race {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	timeFields := strings.Fields(lines[0])[1:]
	distanceFields := strings.Fields(lines[1])[1:]

	races := make([]race, len(timeFields))
	for i := range timeFields {
		races[i] = race{
			time:     aoc.MustAtoi(timeFields[i]),
			distance: aoc.MustAtoi(distanceFields[i]),
		}
	}
	return races
}

func parseSingleRace(input string) race {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	timeStr := strings.ReplaceAll(strings.TrimPrefix(lines[0], "Time:"), " ", "")
	distanceStr := strings.ReplaceAll(strings.TrimPrefix(lines[1], "Distance:"), " ", "")

	return race{
		time:     aoc.MustAtoi(timeStr),
		distance: aoc.MustAtoi(distanceStr),
	}
}

func waysToWin(r race) int {
	t := float64(r.time)
	d := float64(r.distance)

	discriminant := t*t - 4*d
	if discriminant < 0 {
		return 0
	}

	sqrtD := math.Sqrt(discriminant)
	minHold := (t - sqrtD) / 2
	maxHold := (t + sqrtD) / 2

	minInt := int(math.Floor(minHold + 1))
	maxInt := int(math.Ceil(maxHold - 1))

	if minInt > maxInt {
		return 0
	}

	return maxInt - minInt + 1
}

func day06p01(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	races := parseRaces(string(data))
	ways := xslices.Map(waysToWin, races)
	product := xslices.Product(ways)

	return strconv.Itoa(product), nil
}

func day06p02(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	race := parseSingleRace(string(data))
	ways := waysToWin(race)

	return strconv.Itoa(ways), nil
}
