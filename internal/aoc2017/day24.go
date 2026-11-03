package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type bridgeComponent struct{ a, b int64 }

type bridgeScore struct {
	strongest       int64
	longestLength   int
	longestStrength int64
}

func buildBridges(r io.Reader) (bridgeScore, error) {
	s := scanner.NewScanner(r, func(line []byte) (bridgeComponent, error) {
		left, right, ok := bytes.Cut(line, []byte("/"))
		if !ok {
			return bridgeComponent{}, fmt.Errorf("invalid bridge component %q", line)
		}
		a, err := convert.ScanNumber[int64](left)
		if err != nil {
			return bridgeComponent{}, err
		}
		b, err := convert.ScanNumber[int64](right)
		if err != nil {
			return bridgeComponent{}, err
		}
		if a < 0 || b < 0 || a > math.MaxInt64-b {
			return bridgeComponent{}, fmt.Errorf("invalid bridge port strengths %q", line)
		}
		return bridgeComponent{a: a, b: b}, nil
	})
	var components []bridgeComponent
	byPort := make(map[int64][]int)
	var total int64
	for component := range s.Values() {
		if len(components) == 64 {
			return bridgeScore{}, fmt.Errorf("bridge search supports at most 64 components")
		}
		strength := component.a + component.b
		if total > math.MaxInt64-strength {
			return bridgeScore{}, fmt.Errorf("total bridge strength overflows int64")
		}
		total += strength
		id := len(components)
		components = append(components, component)
		byPort[component.a] = append(byPort[component.a], id)
		if component.b != component.a {
			byPort[component.b] = append(byPort[component.b], id)
		}
	}
	if err := s.Err(); err != nil {
		return bridgeScore{}, err
	}

	type state struct {
		port int64
		used uint64
	}
	memo := make(map[state]bridgeScore)
	states := 0
	var visit func(state) (bridgeScore, error)
	visit = func(current state) (bridgeScore, error) {
		if score, ok := memo[current]; ok {
			return score, nil
		}
		states++
		if states > 1_000_000 {
			return bridgeScore{}, fmt.Errorf("bridge search exceeded 1000000 states")
		}

		var best bridgeScore
		for _, id := range byPort[current.port] {
			bit := uint64(1) << id
			if current.used&bit != 0 {
				continue
			}
			component := components[id]
			port := component.a
			if port == current.port {
				port = component.b
			}
			suffix, err := visit(state{port: port, used: current.used | bit})
			if err != nil {
				return bridgeScore{}, err
			}
			strength := component.a + component.b
			best.strongest = max(best.strongest, strength+suffix.strongest)
			length := suffix.longestLength + 1
			if length > best.longestLength || length == best.longestLength && strength+suffix.longestStrength > best.longestStrength {
				best.longestLength = length
				best.longestStrength = strength + suffix.longestStrength
			}
		}
		memo[current] = best
		return best, nil
	}

	return visit(state{})
}

func day24p01(r io.Reader) (string, error) {
	score, err := buildBridges(r)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(score.strongest, 10), nil
}

func day24p02(r io.Reader) (string, error) {
	score, err := buildBridges(r)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(score.longestStrength, 10), nil
}
