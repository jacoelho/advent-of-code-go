package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
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
		if _, ok := xmath.CheckedAdd(a, b); a < 0 || b < 0 || !ok {
			return bridgeComponent{}, fmt.Errorf("invalid bridge port strengths %q", line)
		}
		return bridgeComponent{a: a, b: b}, nil
	})
	var components []bridgeComponent
	byPort := make(map[int64][]int)
	var total int64
	for component := range s.Values() {
		strength := component.a + component.b
		sum, ok := xmath.CheckedAdd(total, strength)
		if !ok {
			return bridgeScore{}, fmt.Errorf("total bridge strength overflows int64")
		}
		total = sum
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

	if len(components) == 0 {
		return bridgeScore{}, fmt.Errorf("no bridge components")
	}

	used := make([]bool, len(components))
	var visit func(port int64) bridgeScore
	visit = func(port int64) bridgeScore {
		var best bridgeScore
		for _, id := range byPort[port] {
			if used[id] {
				continue
			}
			component := components[id]
			next := component.a
			if next == port {
				next = component.b
			}
			used[id] = true
			suffix := visit(next)
			used[id] = false
			strength := component.a + component.b
			best.strongest = max(best.strongest, strength+suffix.strongest)
			length := suffix.longestLength + 1
			if length > best.longestLength || length == best.longestLength && strength+suffix.longestStrength > best.longestStrength {
				best.longestLength = length
				best.longestStrength = strength + suffix.longestStrength
			}
		}
		return best
	}

	return visit(0), nil
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
