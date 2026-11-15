package aoc2025

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/bits"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

const (
	noSolution        = -1
	maxMachineWidth   = 64
	maxMachineButtons = 16
)

type machineRequirements struct {
	lightsMask  uint64
	buttonMasks []uint64
	joltage     []int64
}

func parseLightsToken(token string) (uint64, error) {
	if len(token) < 3 || token[0] != '[' || token[len(token)-1] != ']' {
		return 0, fmt.Errorf("invalid lights token %q", token)
	}
	pattern := token[1 : len(token)-1]
	if len(pattern) > maxMachineWidth {
		return 0, fmt.Errorf("machine exceeds %d lights", maxMachineWidth)
	}

	var mask uint64
	for i, ch := range pattern {
		switch ch {
		case '#':
			mask |= uint64(1) << i
		case '.':
		default:
			return 0, fmt.Errorf("invalid light %q", ch)
		}
	}
	return mask, nil
}

func parseButtonTokens(tokens []string, width int) ([]uint64, error) {
	if len(tokens) > maxMachineButtons {
		return nil, fmt.Errorf("machine exceeds %d buttons", maxMachineButtons)
	}
	buttonMasks := make([]uint64, 0, len(tokens))
	for _, token := range tokens {
		if len(token) < 2 || token[0] != '(' || token[len(token)-1] != ')' {
			return nil, fmt.Errorf("invalid button token %q", token)
		}
		var mask uint64
		if token != "()" {
			for field := range strings.SplitSeq(token[1:len(token)-1], ",") {
				idx, err := convert.ScanNumber[int]([]byte(field))
				if err != nil {
					return nil, fmt.Errorf("button index: %w", err)
				}
				if idx < 0 || idx >= width {
					return nil, fmt.Errorf("button index %d outside %d lights", idx, width)
				}
				bit := uint64(1) << idx
				if mask&bit != 0 {
					return nil, fmt.Errorf("duplicate button index %d", idx)
				}
				mask |= bit
			}
		}
		buttonMasks = append(buttonMasks, mask)
	}
	return buttonMasks, nil
}

func parseJoltageToken(token string, width int) ([]int64, error) {
	if len(token) < 3 || token[0] != '{' || token[len(token)-1] != '}' {
		return nil, fmt.Errorf("invalid joltage token %q", token)
	}
	var joltage []int64
	for field := range strings.SplitSeq(token[1:len(token)-1], ",") {
		value, err := convert.ScanNumber[int64]([]byte(field))
		if err != nil {
			return nil, fmt.Errorf("joltage count: %w", err)
		}
		if value < 0 {
			return nil, fmt.Errorf("joltage count must be nonnegative")
		}
		joltage = append(joltage, value)
		if len(joltage) > width {
			return nil, fmt.Errorf("joltage width differs from lights")
		}
	}
	if len(joltage) != width {
		return nil, fmt.Errorf("joltage width differs from lights")
	}
	return joltage, nil
}

func parseMachine(line []byte) (machineRequirements, error) {
	tokens := strings.Fields(string(line))
	if len(tokens) < 2 {
		return machineRequirements{}, fmt.Errorf("machine needs lights and joltage")
	}
	lights, err := parseLightsToken(tokens[0])
	if err != nil {
		return machineRequirements{}, err
	}
	width := len(tokens[0]) - 2
	joltage, err := parseJoltageToken(tokens[len(tokens)-1], width)
	if err != nil {
		return machineRequirements{}, err
	}
	buttons, err := parseButtonTokens(tokens[1:len(tokens)-1], width)
	if err != nil {
		return machineRequirements{}, err
	}
	return machineRequirements{lightsMask: lights, buttonMasks: buttons, joltage: joltage}, nil
}

func parseInitializationProcedure(r io.Reader) ([]machineRequirements, error) {
	machineIndex := 0
	s := scanner.NewScanner(r, func(line []byte) (machineRequirements, error) {
		machineIndex++
		machine, err := parseMachine(line)
		if err != nil {
			return machineRequirements{}, fmt.Errorf("machine %d: %w", machineIndex, err)
		}
		return machine, nil
	})
	steps := slices.Collect(s.Values())
	return steps, s.Err()
}

// buttonSubset is one choice of buttons pressed once each and the amount it
// adds to every counter.
type buttonSubset struct {
	presses int64
	effect  []int64
}

// subsetsByParity groups every button subset by the counters it changes an odd
// number of times; that set is also the lights the subset toggles.
func subsetsByParity(step machineRequirements) map[uint64][]buttonSubset {
	n, width := len(step.buttonMasks), len(step.joltage)
	toggled := make([]uint64, 1<<n)
	effects := make([]int64, (1<<n)*width)
	groups := make(map[uint64][]buttonSubset)
	for subset := range 1 << n {
		effect := effects[subset*width : (subset+1)*width : (subset+1)*width]
		if subset != 0 {
			button := bits.TrailingZeros(uint(subset))
			rest := subset & (subset - 1)
			mask := step.buttonMasks[button]
			toggled[subset] = toggled[rest] ^ mask
			copy(effect, effects[rest*width:(rest+1)*width])
			for counter := range width {
				if mask&(uint64(1)<<counter) != 0 {
					effect[counter]++
				}
			}
		}
		groups[toggled[subset]] = append(groups[toggled[subset]], buttonSubset{
			presses: int64(bits.OnesCount(uint(subset))),
			effect:  effect,
		})
	}
	return groups
}

func minButtonPresses(step machineRequirements, groups map[uint64][]buttonSubset) (int64, error) {
	best := int64(noSolution)
	for _, subset := range groups[step.lightsMask] {
		if best == noSolution || subset.presses < best {
			best = subset.presses
		}
	}
	if best == noSolution {
		return 0, fmt.Errorf("lights have no reachable configuration")
	}
	return best, nil
}

// minJoltagePresses finds the fewest presses reaching the joltage targets. In
// any solution, the buttons pressed an odd number of times form a subset whose
// effect matches the targets' parity; removing it leaves every button pressed
// an even number of times, which is twice a solution for the halved remainder.
func minJoltagePresses(step machineRequirements, groups map[uint64][]buttonSubset) (int64, error) {
	width := len(step.joltage)
	memo := make(map[string]int64)
	var solve func(target []int64) (int64, error)
	solve = func(target []int64) (int64, error) {
		var parity uint64
		zero := true
		for counter, value := range target {
			parity |= uint64(value&1) << counter
			zero = zero && value == 0
		}
		if zero {
			return 0, nil
		}
		key := joltageKey(target)
		if cost, known := memo[key]; known {
			return cost, nil
		}
		best := int64(noSolution)
		next := make([]int64, width)
	subsets:
		for _, subset := range groups[parity] {
			for counter, value := range target {
				if subset.effect[counter] > value {
					continue subsets
				}
				next[counter] = (value - subset.effect[counter]) / 2
			}
			cost, err := solve(next)
			if err != nil {
				return 0, err
			}
			if cost == noSolution {
				continue
			}
			if cost > (math.MaxInt64-subset.presses)/2 {
				return 0, fmt.Errorf("joltage presses exceed int64")
			}
			if total := subset.presses + 2*cost; best == noSolution || total < best {
				best = total
			}
		}
		memo[key] = best
		return best, nil
	}
	best, err := solve(step.joltage)
	if err != nil {
		return 0, err
	}
	if best == noSolution {
		return 0, fmt.Errorf("joltage has no nonnegative integer solution")
	}
	return best, nil
}

func joltageKey(target []int64) string {
	key := make([]byte, 0, 8*len(target))
	for _, value := range target {
		key = binary.LittleEndian.AppendUint64(key, uint64(value))
	}
	return string(key)
}

func day10(r io.Reader, presses func(machineRequirements, map[uint64][]buttonSubset) (int64, error)) (string, error) {
	steps, err := parseInitializationProcedure(r)
	if err != nil {
		return "", err
	}
	var total int64
	for i, step := range steps {
		count, err := presses(step, subsetsByParity(step))
		if err != nil {
			return "", fmt.Errorf("machine %d: %w", i+1, err)
		}
		sum, ok := xmath.CheckedAdd(total, count)
		if !ok {
			return "", fmt.Errorf("machine %d: total presses exceed int64", i+1)
		}
		total = sum
	}
	return strconv.FormatInt(total, 10), nil
}

func day10p01(r io.Reader) (string, error) {
	return day10(r, minButtonPresses)
}

func day10p02(r io.Reader) (string, error) {
	return day10(r, minJoltagePresses)
}
