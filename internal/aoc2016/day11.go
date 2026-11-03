package aoc2016

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

type elevatorState struct {
	floor uint8
	pairs [15]uint8
}

func safeElevatorState(state elevatorState, count int) bool {
	var generators [4]bool
	for _, pair := range state.pairs[:count] {
		generators[pair>>2] = true
	}
	for _, pair := range state.pairs[:count] {
		chip, generator := pair&3, pair>>2
		if chip != generator && generators[chip] {
			return false
		}
	}
	return true
}

func parseElevator(r io.Reader, extraPairs int) (elevatorState, int, error) {
	pattern := regexp.MustCompile(`^a ([a-z]+)( generator|-compatible microchip)$`)
	floors := [...]string{"first", "second", "third", "fourth"}
	pairs := make(map[string][2]uint8)
	floor := 0
	s := bufio.NewScanner(r)
	for s.Scan() {
		if floor == len(floors) {
			return elevatorState{}, 0, fmt.Errorf("too many elevator floors")
		}
		contents, ok := strings.CutPrefix(s.Text(), "The "+floors[floor]+" floor contains ")
		contents, suffixOK := strings.CutSuffix(contents, ".")
		if !ok || !suffixOK {
			return elevatorState{}, 0, fmt.Errorf("invalid floor description %q", s.Text())
		}
		if contents != "nothing relevant" {
			contents = strings.ReplaceAll(contents, ", and ", ", ")
			contents = strings.ReplaceAll(contents, " and ", ", ")
			for item := range strings.SplitSeq(contents, ", ") {
				match := pattern.FindStringSubmatch(item)
				if match == nil {
					return elevatorState{}, 0, fmt.Errorf("invalid elevator item %q", item)
				}
				pair, exists := pairs[match[1]]
				if !exists {
					pair = [2]uint8{4, 4}
				}
				kind := 0
				if match[2] == " generator" {
					kind = 1
				}
				if pair[kind] != 4 {
					return elevatorState{}, 0, fmt.Errorf("duplicate elevator item %q", item)
				}
				pair[kind] = uint8(floor)
				pairs[match[1]] = pair
			}
		}
		floor++
	}
	if err := s.Err(); err != nil {
		return elevatorState{}, 0, err
	}
	if floor != 4 || len(pairs) == 0 || len(pairs)+extraPairs > 15 {
		return elevatorState{}, 0, fmt.Errorf("expected four floors with 1 to 15 complete element pairs")
	}
	state := elevatorState{}
	count := 0
	for name, pair := range pairs {
		if pair[0] == 4 || pair[1] == 4 {
			return elevatorState{}, 0, fmt.Errorf("incomplete element pair %s", name)
		}
		state.pairs[count] = pair[0] | pair[1]<<2
		count++
	}
	count += extraPairs
	if !safeElevatorState(state, count) {
		return elevatorState{}, 0, fmt.Errorf("unsafe initial elevator layout")
	}
	// Element names are interchangeable; retain each chip's generator pairing.
	slices.Sort(state.pairs[:count])
	return state, count, nil
}

func elevatorMoves(r io.Reader, extraPairs int) (string, error) {
	initial, count, err := parseElevator(r, extraPairs)
	if err != nil {
		return "", err
	}
	queue := collections.NewDeque[elevatorState](256)
	queue.PushBack(initial)
	distances := map[elevatorState]int{initial: 0}
	for queue.Size() > 0 {
		state, _ := queue.PopFront()
		distance := distances[state]
		finished := true
		var items []int
		for i, pair := range state.pairs[:count] {
			finished = finished && pair == 15
			if pair&3 == state.floor {
				items = append(items, 2*i)
			}
			if pair>>2 == state.floor {
				items = append(items, 2*i+1)
			}
		}
		if finished {
			return strconv.Itoa(distance), nil
		}
		for _, delta := range []int{-1, 1} {
			floor := int(state.floor) + delta
			if floor < 0 || floor > 3 {
				continue
			}
			for i, first := range items {
				for j := i; j < len(items); j++ {
					next := state
					next.floor = uint8(floor)
					for _, item := range []int{first, items[j]} {
						shift := 2 * (item % 2)
						next.pairs[item/2] = next.pairs[item/2]&^(3<<shift) | uint8(floor)<<shift
					}
					if !safeElevatorState(next, count) {
						continue
					}
					slices.Sort(next.pairs[:count])
					if _, seen := distances[next]; seen {
						continue
					}
					if len(distances) == 1_000_000 {
						return "", fmt.Errorf("elevator search exceeds 1000000 states")
					}
					distances[next] = distance + 1
					queue.PushBack(next)
				}
			}
		}
	}
	return "", fmt.Errorf("elevator layout has no solution")
}

func day11p01(r io.Reader) (string, error) {
	return elevatorMoves(r, 0)
}

func day11p02(r io.Reader) (string, error) {
	return elevatorMoves(r, 2)
}
