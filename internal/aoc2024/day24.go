package aoc2024

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/graph"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type logicGate struct {
	left, operation, right, output string
}

type monitoringDevice struct {
	initial map[string]bool
	gates   []logicGate
}

func validLogicWire(name string) bool {
	if name == "" {
		return false
	}
	for _, c := range name {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func parseMonitoringDevice(r io.Reader) (monitoringDevice, error) {
	blocks, err := scanner.ReadBlocks(r)
	if err != nil {
		return monitoringDevice{}, err
	}
	if len(blocks) != 2 || len(blocks[1]) == 0 {
		return monitoringDevice{}, fmt.Errorf("circuit needs initial wires and gates")
	}

	device := monitoringDevice{initial: make(map[string]bool)}
	producers := make(map[string]bool)
	for _, line := range blocks[0] {
		name, value, ok := strings.Cut(line, ":")
		name, value = strings.TrimSpace(name), strings.TrimSpace(value)
		if !ok || !validLogicWire(name) || producers[name] {
			return monitoringDevice{}, fmt.Errorf("invalid or duplicate initial wire %q", name)
		}
		bit, err := convert.ScanNumber[int]([]byte(value))
		if err != nil {
			return monitoringDevice{}, fmt.Errorf("initial wire %s: %w", name, err)
		}
		if bit != 0 && bit != 1 {
			return monitoringDevice{}, fmt.Errorf("initial wire %s is not a bit", name)
		}
		device.initial[name], producers[name] = bit == 1, true
	}

	for _, line := range blocks[1] {
		fields := strings.Fields(line)
		if len(fields) != 5 || fields[3] != "->" {
			return monitoringDevice{}, fmt.Errorf("invalid gate %q", line)
		}
		gate := logicGate{left: fields[0], operation: fields[1], right: fields[2], output: fields[4]}
		if !validLogicWire(gate.left) || !validLogicWire(gate.right) || !validLogicWire(gate.output) {
			return monitoringDevice{}, fmt.Errorf("invalid wire name in gate %q", line)
		}
		if gate.operation != "AND" && gate.operation != "OR" && gate.operation != "XOR" {
			return monitoringDevice{}, fmt.Errorf("unknown gate operation %q", gate.operation)
		}
		if producers[gate.output] {
			return monitoringDevice{}, fmt.Errorf("duplicate producer for %s", gate.output)
		}
		producers[gate.output] = true
		device.gates = append(device.gates, gate)
	}

	for _, gate := range device.gates {
		if !producers[gate.left] || !producers[gate.right] {
			return monitoringDevice{}, fmt.Errorf("gate %s references an unknown wire", gate.output)
		}
	}
	return device, nil
}

func evaluateDevice(device monitoringDevice) (map[string]bool, error) {
	outputs := make([]string, len(device.gates))
	byOutput := make(map[string]logicGate, len(device.gates))
	for i, gate := range device.gates {
		outputs[i], byOutput[gate.output] = gate.output, gate
	}
	order, err := graph.TopologicalOrder(outputs, func(output string) iter.Seq[string] {
		return func(yield func(string) bool) {
			gate := byOutput[output]
			for _, input := range [...]string{gate.left, gate.right} {
				if _, initial := device.initial[input]; !initial && !yield(input) {
					return
				}
			}
		}
	})
	if err != nil {
		return nil, fmt.Errorf("circuit dependencies: %w", err)
	}
	values := maps.Clone(device.initial)
	for i := len(order) - 1; i >= 0; i-- {
		gate := byOutput[order[i]]
		left, right := values[gate.left], values[gate.right]
		switch gate.operation {
		case "AND":
			values[gate.output] = left && right
		case "OR":
			values[gate.output] = left || right
		case "XOR":
			values[gate.output] = left != right
		}
	}
	return values, nil
}

func deviceNumber(values map[string]bool, prefix byte) (uint64, error) {
	var result, seen uint64
	for wire, value := range values {
		if wire[0] != prefix {
			continue
		}
		index, err := convert.ScanNumber[int]([]byte(wire[1:]))
		if err != nil {
			return 0, fmt.Errorf("output wire %s: %w", wire, err)
		}
		if index < 0 || index >= 64 {
			return 0, fmt.Errorf("output wire %s: bit index must be between 0 and 63", wire)
		}
		if seen&(1<<index) != 0 {
			return 0, fmt.Errorf("output wire %s: duplicate bit index %d", wire, index)
		}
		seen |= 1 << index
		if value {
			result |= 1 << index
		}
	}
	if seen == 0 {
		return 0, fmt.Errorf("circuit has no %c output wires", prefix)
	}
	return result, nil
}

func day24p01(r io.Reader) (string, error) {
	device, err := parseMonitoringDevice(r)
	if err != nil {
		return "", err
	}
	values, err := evaluateDevice(device)
	if err != nil {
		return "", err
	}
	result, err := deviceNumber(values, 'z')
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(result, 10), nil
}

type logicInputs struct {
	operation, left, right string
}

func logicSignature(operation, left, right string) logicInputs {
	if left > right {
		left, right = right, left
	}
	return logicInputs{operation: operation, left: left, right: right}
}

// certifiesRippleAdder reports whether the gates, with outputs[i] as the
// output of gates[i], form exactly a width-bit ripple-carry adder of x and y
// into z. lookup maps each gate's input signature to its index; output swaps
// do not change it.
func certifiesRippleAdder(gates []logicGate, lookup map[logicInputs]int, outputs []string, width int) bool {
	used := make([]bool, len(gates))
	find := func(operation, left, right string) (string, bool) {
		index, ok := lookup[logicSignature(operation, left, right)]
		if !ok || used[index] {
			return "", false
		}
		used[index] = true
		return outputs[index], true
	}
	sum, ok := find("XOR", "x00", "y00")
	if !ok || sum != "z00" {
		return false
	}
	carry, ok := find("AND", "x00", "y00")
	if !ok {
		return false
	}
	for i := 1; i < width; i++ {
		x, y := fmt.Sprintf("x%02d", i), fmt.Sprintf("y%02d", i)
		p, haveP := find("XOR", x, y)
		g, haveG := find("AND", x, y)
		sum, haveSum := find("XOR", p, carry)
		half, haveHalf := find("AND", p, carry)
		next, haveNext := find("OR", g, half)
		if !haveP || !haveG || !haveSum || !haveHalf || !haveNext || sum != fmt.Sprintf("z%02d", i) {
			return false
		}
		carry = next
	}
	return carry == fmt.Sprintf("z%02d", width) && !slices.Contains(used, false)
}

var errDay24RepairUncertified = errors.New("swapping the flagged wires does not yield a ripple-carry adder")

func day24AdderWidth(device monitoringDevice) (int, error) {
	width := len(device.initial) / 2
	for i := range width {
		for _, prefix := range []byte{'x', 'y'} {
			if _, ok := device.initial[fmt.Sprintf("%c%02d", prefix, i)]; !ok {
				return 0, fmt.Errorf("adder inputs must be contiguous x and y bits from 00")
			}
		}
	}
	if width == 0 || len(device.initial) != 2*width {
		return 0, fmt.Errorf("adder needs equally many x and y input bits")
	}
	return width, nil
}

// day24FlaggedWires returns the outputs that break the shape every gate has in
// a ripple-carry adder z = x + y, where for each bit i ≥ 1
// p = x XOR y, g = x AND y, z = p XOR c, h = p AND c, c' = g OR h,
// bit 0 is z00 = x00 XOR y00 with carry x00 AND y00, and the final carry is
// the most significant z.
func day24FlaggedWires(device monitoringDevice, mostSignificant string) []string {
	consumers := make(map[string][]string)
	for _, gate := range device.gates {
		consumers[gate.left] = append(consumers[gate.left], gate.operation)
		consumers[gate.right] = append(consumers[gate.right], gate.operation)
	}
	flagged := make(map[string]bool)
	for _, gate := range device.gates {
		_, leftInput := device.initial[gate.left]
		_, rightInput := device.initial[gate.right]
		fromInputs := leftInput && rightInput
		firstBit := min(gate.left, gate.right) == "x00" && max(gate.left, gate.right) == "y00"
		toZ := gate.output[0] == 'z'
		switch {
		// Every z below the top is a sum p XOR c; the top z is the final carry OR.
		case toZ && gate.output != mostSignificant && gate.operation != "XOR":
			flagged[gate.output] = true
		// The only XOR not reading x and y is the sum gate, which drives a z.
		case gate.operation == "XOR" && !fromInputs && !toZ:
			flagged[gate.output] = true
		// p feeds the sum XOR of its bit; bit 0 has no carry in, so x00 XOR y00 is z00.
		case gate.operation == "XOR" && fromInputs && !firstBit && !slices.Contains(consumers[gate.output], "XOR"):
			flagged[gate.output] = true
		// g and h feed the carry OR; bit 0's x00 AND y00 is itself the carry.
		case gate.operation == "AND" && !firstBit && !slices.Contains(consumers[gate.output], "OR"):
			flagged[gate.output] = true
		}
	}
	return slices.Sorted(maps.Keys(flagged))
}

// day24Pairings yields every way to split the indices 0..n-1 into disjoint
// pairs; n must be even. The yielded slice is reused between iterations.
func day24Pairings(n int) iter.Seq[[][2]int] {
	return func(yield func([][2]int) bool) {
		pairs := make([][2]int, 0, n/2)
		paired := make([]bool, n)
		var visit func() bool
		visit = func() bool {
			first := slices.Index(paired, false)
			if first < 0 {
				return yield(pairs)
			}
			paired[first] = true
			defer func() { paired[first] = false }()
			for second := first + 1; second < n; second++ {
				if paired[second] {
					continue
				}
				paired[second] = true
				pairs = append(pairs, [2]int{first, second})
				ok := visit()
				pairs = pairs[:len(pairs)-1]
				paired[second] = false
				if !ok {
					return false
				}
			}
			return true
		}
		visit()
	}
}

// repairRippleAdder returns four output swaps that turn device into a
// ripple-carry adder, choosing them only among the wires the adder shape flags.
func repairRippleAdder(device monitoringDevice) ([][2]string, error) {
	width, err := day24AdderWidth(device)
	if err != nil {
		return nil, err
	}
	flagged := day24FlaggedWires(device, fmt.Sprintf("z%02d", width))
	if len(flagged) != 8 {
		return nil, fmt.Errorf("adder shape flags %d wires, want 8: %v", len(flagged), flagged)
	}
	lookup := make(map[logicInputs]int, len(device.gates))
	outputs := make([]string, len(device.gates))
	producer := make(map[string]int, len(device.gates))
	for i, gate := range device.gates {
		lookup[logicSignature(gate.operation, gate.left, gate.right)] = i
		outputs[i] = gate.output
		producer[gate.output] = i
	}
	// Every pairing covers all flagged wires, so it overwrites every slot the
	// previous pairing changed.
	for pairs := range day24Pairings(len(flagged)) {
		for _, pair := range pairs {
			outputs[producer[flagged[pair[0]]]] = flagged[pair[1]]
			outputs[producer[flagged[pair[1]]]] = flagged[pair[0]]
		}
		if certifiesRippleAdder(device.gates, lookup, outputs, width) {
			swaps := make([][2]string, len(pairs))
			for i, pair := range pairs {
				swaps[i] = [2]string{flagged[pair[0]], flagged[pair[1]]}
			}
			return swaps, nil
		}
	}
	return nil, fmt.Errorf("flagged wires %v: %w", flagged, errDay24RepairUncertified)
}

func day24p02(r io.Reader) (string, error) {
	device, err := parseMonitoringDevice(r)
	if err != nil {
		return "", err
	}
	pairs, err := repairRippleAdder(device)
	if err != nil {
		return "", err
	}
	wires := make([]string, 0, 2*len(pairs))
	for _, pair := range pairs {
		wires = append(wires, pair[0], pair[1])
	}
	slices.Sort(wires)
	return strings.Join(wires, ","), nil
}
