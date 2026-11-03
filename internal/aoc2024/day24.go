package aoc2024

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"maps"
	"math/big"
	"slices"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/graph"
)

type logicGate struct {
	left, operation, right, output string
}

type monitoringDevice struct {
	initial map[string]bool
	gates   []logicGate
}

func validLogicWire(name string) bool {
	if len(name) == 0 || len(name) > 64 {
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
	const inputLimit = 4 << 20
	data, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return monitoringDevice{}, err
	}
	if len(data) > inputLimit {
		return monitoringDevice{}, fmt.Errorf("circuit input exceeds 4 MiB")
	}

	input := strings.TrimSpace(strings.ReplaceAll(string(data), "\r\n", "\n"))
	initial, connections, ok := strings.Cut(input, "\n\n")
	if !ok {
		return monitoringDevice{}, fmt.Errorf("circuit needs initial wires and gates")
	}

	device := monitoringDevice{initial: make(map[string]bool)}
	producers := make(map[string]bool)
	for line := range strings.SplitSeq(initial, "\n") {
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
		if len(device.initial) == 4096 {
			return monitoringDevice{}, fmt.Errorf("circuit exceeds 4096 initial wires")
		}
		device.initial[name], producers[name] = bit == 1, true
	}

	for line := range strings.SplitSeq(connections, "\n") {
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
		if len(device.gates) == 4096 {
			return monitoringDevice{}, fmt.Errorf("circuit exceeds 4096 gates")
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

func deviceNumber(values map[string]bool, prefix byte) (*big.Int, error) {
	result := new(big.Int)
	seen := make(map[int]bool)
	for wire, value := range values {
		if wire[0] != prefix {
			continue
		}
		index, err := convert.ScanNumber[int]([]byte(wire[1:]))
		if err != nil {
			return nil, fmt.Errorf("output wire %s: %w", wire, err)
		}
		if index < 0 || index >= 4096 || seen[index] {
			return nil, fmt.Errorf("output bit indices must be unique and between 0 and 4095")
		}
		seen[index] = true
		if value {
			result.SetBit(result, index, 1)
		}
	}
	if len(seen) == 0 {
		return nil, fmt.Errorf("circuit has no %c output wires", prefix)
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
	return result.String(), nil
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

type adderPair struct {
	sum, half, left, right int
}

type adderCarry struct {
	gate, left, right int
}

type rippleAdder struct {
	device                          monitoringDevice
	width                           int
	labels                          []string
	outputs, propagate, generate, z []int
	pairs                           []adderPair
	carries                         []adderCarry
	forced                          []bool
}

func inspectRippleAdder(device monitoringDevice) (rippleAdder, error) {
	width := len(device.initial) / 2
	if width < 3 || width > 64 || len(device.initial) != 2*width || len(device.gates) != 5*width-3 {
		return rippleAdder{}, fmt.Errorf("repair requires a canonical 3 to 64 bit ripple adder")
	}
	for i := range width {
		for _, prefix := range []byte{'x', 'y'} {
			if _, ok := device.initial[fmt.Sprintf("%c%02d", prefix, i)]; !ok {
				return rippleAdder{}, fmt.Errorf("adder inputs must be contiguous x/y bits")
			}
		}
	}
	adder := rippleAdder{device: device, width: width}
	lookup := make(map[logicInputs]int)
	for i, gate := range device.gates {
		key := logicSignature(gate.operation, gate.left, gate.right)
		if _, exists := lookup[key]; exists {
			return rippleAdder{}, fmt.Errorf("adder has duplicate gate signatures")
		}
		lookup[key] = i
		adder.labels = append(adder.labels, gate.output)
	}
	slices.Sort(adder.labels)
	labelIndex := make(map[string]int, len(adder.labels))
	for i, label := range adder.labels {
		labelIndex[label] = i
	}
	adder.outputs = make([]int, len(device.gates))
	for i, gate := range device.gates {
		adder.outputs[i] = labelIndex[gate.output]
	}
	for i := range width + 1 {
		label, ok := labelIndex[fmt.Sprintf("z%02d", i)]
		if !ok {
			return rippleAdder{}, fmt.Errorf("adder outputs must be contiguous z bits")
		}
		adder.z = append(adder.z, label)
	}
	for _, label := range adder.labels {
		if strings.HasPrefix(label, "z") && !slices.Contains(adder.z, labelIndex[label]) {
			return rippleAdder{}, fmt.Errorf("adder has an unexpected z output")
		}
	}
	primary := make([]bool, len(device.gates))
	for i := range width {
		x, y := fmt.Sprintf("x%02d", i), fmt.Sprintf("y%02d", i)
		p, haveP := lookup[logicSignature("XOR", x, y)]
		g, haveG := lookup[logicSignature("AND", x, y)]
		if !haveP || !haveG {
			return rippleAdder{}, fmt.Errorf("adder lacks primary gates for bit %d", i)
		}
		adder.propagate = append(adder.propagate, p)
		adder.generate = append(adder.generate, g)
		primary[p], primary[g] = true, true
	}
	const (
		sumLabel   = 1
		pairLabel  = 2
		carryLabel = 3
	)
	categories := make([]int, len(adder.labels))
	for _, label := range adder.z {
		categories[label] = sumLabel
	}
	claimInputs := func(gate logicGate, category int) (int, int, bool) {
		left, haveLeft := labelIndex[gate.left]
		right, haveRight := labelIndex[gate.right]
		if !haveLeft || !haveRight || left == right || categories[left] != 0 || categories[right] != 0 {
			return 0, 0, false
		}
		categories[left], categories[right] = category, category
		return left, right, true
	}
	for i, gate := range device.gates {
		if primary[i] {
			continue
		}
		switch gate.operation {
		case "XOR":
			half, ok := lookup[logicSignature("AND", gate.left, gate.right)]
			if !ok || primary[half] {
				return rippleAdder{}, fmt.Errorf("adder lacks a matching carry AND gate")
			}
			left, right, ok := claimInputs(gate, pairLabel)
			if !ok {
				return rippleAdder{}, fmt.Errorf("unsupported adder sum input connections")
			}
			adder.pairs = append(adder.pairs, adderPair{sum: i, half: half, left: left, right: right})
		case "OR":
			left, right, ok := claimInputs(gate, carryLabel)
			if !ok {
				return rippleAdder{}, fmt.Errorf("unsupported adder carry input connections")
			}
			adder.carries = append(adder.carries, adderCarry{gate: i, left: left, right: right})
		}
	}
	if len(adder.pairs) != width-1 || len(adder.carries) != width-1 || slices.Contains(categories, 0) {
		return rippleAdder{}, fmt.Errorf("circuit does not have complete ripple-adder roles")
	}
	slices.SortFunc(adder.pairs, func(a, b adderPair) int {
		return strings.Compare(device.gates[a.sum].output, device.gates[b.sum].output)
	})
	slices.SortFunc(adder.carries, func(a, b adderCarry) int {
		return strings.Compare(device.gates[a.gate].output, device.gates[b.gate].output)
	})
	adder.forced = make([]bool, len(adder.labels))
	forceCategory := func(gate, category int) {
		label := adder.outputs[gate]
		adder.forced[label] = categories[label] != category
	}
	adder.forced[adder.outputs[adder.propagate[0]]] = adder.outputs[adder.propagate[0]] != adder.z[0]
	forceCategory(adder.generate[0], pairLabel)
	for i := 1; i < width; i++ {
		forceCategory(adder.propagate[i], pairLabel)
		forceCategory(adder.generate[i], carryLabel)
	}
	for _, pair := range adder.pairs {
		label := adder.outputs[pair.sum]
		adder.forced[label] = categories[label] != sumLabel || label == adder.z[0] || label == adder.z[width]
		forceCategory(pair.half, carryLabel)
	}
	for _, carry := range adder.carries {
		label := adder.outputs[carry.gate]
		adder.forced[label] = categories[label] != pairLabel && label != adder.z[width]
	}
	return adder, nil
}

func certifiesRippleAdder(device monitoringDevice, width int) bool {
	lookup := make(map[logicInputs]int)
	for i, gate := range device.gates {
		lookup[logicSignature(gate.operation, gate.left, gate.right)] = i
	}
	used := make([]bool, len(device.gates))
	find := func(operation, left, right string) (string, bool) {
		index, ok := lookup[logicSignature(operation, left, right)]
		if !ok || used[index] {
			return "", false
		}
		used[index] = true
		return device.gates[index].output, true
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

type adderRepair struct {
	device monitoringDevice
	wires  []string
}

var (
	errAdderRepairLimit     = errors.New("adder repair work limit exceeded")
	errAmbiguousAdderRepair = errors.New("adder repair has ambiguous swapped-wire sets")
)

type repairCheckpoint struct {
	changes, pairs, forced int
}

type outputRepair struct {
	permutation []int
	mustMove    []bool
	changes     []int
	pairs       int
	forced      int
}

func (r *outputRepair) checkpoint() repairCheckpoint {
	return repairCheckpoint{changes: len(r.changes), pairs: r.pairs, forced: r.forced}
}

func (r *outputRepair) restore(checkpoint repairCheckpoint) {
	for _, label := range r.changes[checkpoint.changes:] {
		r.permutation[label] = -1
	}
	r.changes = r.changes[:checkpoint.changes]
	r.pairs, r.forced = checkpoint.pairs, checkpoint.forced
}

func (r *outputRepair) assign(from, to int) bool {
	if r.permutation[from] >= 0 {
		return r.permutation[from] == to
	}
	if from == to {
		if r.mustMove[from] {
			return false
		}
		r.permutation[from] = from
		r.changes = append(r.changes, from)
		return true
	}
	if r.permutation[to] >= 0 || r.pairs == 4 {
		return false
	}
	r.permutation[from], r.permutation[to] = to, from
	r.changes = append(r.changes, from, to)
	r.pairs++
	if r.mustMove[from] {
		r.forced--
	}
	if r.mustMove[to] {
		r.forced--
	}
	return r.forced <= 2*(4-r.pairs)
}

func repairRippleAdder(device monitoringDevice, remaining int) (adderRepair, error) {
	adder, err := inspectRippleAdder(device)
	if err != nil {
		return adderRepair{}, err
	}
	state := outputRepair{permutation: make([]int, len(adder.labels)), mustMove: adder.forced}
	for i := range state.permutation {
		state.permutation[i] = -1
		if adder.forced[i] {
			state.forced++
		}
	}
	if state.forced > 8 {
		return adderRepair{}, fmt.Errorf("more than eight adder outputs must move")
	}
	// Eight necessary endpoints leave no room for any additional moved label.
	if state.forced == 8 {
		for i, forced := range adder.forced {
			if !forced {
				state.permutation[i] = i
			}
		}
	}
	if !state.assign(adder.outputs[adder.propagate[0]], adder.z[0]) {
		return adderRepair{}, fmt.Errorf("adder has no four-swap repair")
	}
	used := make([]bool, len(device.gates))
	for _, gate := range adder.propagate {
		used[gate] = true
	}
	for _, gate := range adder.generate {
		used[gate] = true
	}
	charge := func(work int) error {
		if work > remaining {
			return errAdderRepairLimit
		}
		remaining -= work
		return nil
	}
	var result adderRepair
	var answer string
	completed := make(map[string]bool)
	var visit func(int, int) error
	visit = func(bit, carry int) error {
		if err := charge(1); err != nil {
			return err
		}
		if bit == adder.width {
			checkpoint := state.checkpoint()
			defer state.restore(checkpoint)
			if !state.assign(adder.outputs[carry], adder.z[adder.width]) || state.pairs != 4 || slices.Contains(state.permutation, -1) || slices.Contains(used, false) {
				return nil
			}
			var names, pairs []string
			for from, to := range state.permutation {
				if from != to {
					names = append(names, adder.labels[from])
				}
				if from < to {
					pairs = append(pairs, adder.labels[from]+":"+adder.labels[to])
				}
			}
			key := strings.Join(pairs, ",")
			if completed[key] {
				return nil
			}
			if err := charge(len(device.gates)); err != nil {
				return err
			}
			repaired := monitoringDevice{initial: maps.Clone(device.initial), gates: slices.Clone(device.gates)}
			for i := range repaired.gates {
				repaired.gates[i].output = adder.labels[state.permutation[adder.outputs[i]]]
			}
			if !certifiesRippleAdder(repaired, adder.width) {
				return nil
			}
			candidate := strings.Join(names, ",")
			if answer != "" && answer != candidate {
				return errAmbiguousAdderRepair
			}
			completed[key] = true
			if answer == "" {
				answer, result = candidate, adderRepair{device: repaired, wires: names}
			}
			return nil
		}
		for _, pair := range adder.pairs {
			for orientation := range 2 {
				if err := charge(1); err != nil {
					return err
				}
				if used[pair.sum] || used[pair.half] {
					continue
				}
				left, right := pair.left, pair.right
				if orientation == 1 {
					left, right = right, left
				}
				checkpoint := state.checkpoint()
				if state.assign(adder.outputs[adder.propagate[bit]], left) && state.assign(adder.outputs[carry], right) && state.assign(adder.outputs[pair.sum], adder.z[bit]) {
					used[pair.sum], used[pair.half] = true, true
					for _, next := range adder.carries {
						for orientation := range 2 {
							if err := charge(1); err != nil {
								return err
							}
							if used[next.gate] {
								continue
							}
							left, right := next.left, next.right
							if orientation == 1 {
								left, right = right, left
							}
							carryCheckpoint := state.checkpoint()
							if state.assign(adder.outputs[adder.generate[bit]], left) && state.assign(adder.outputs[pair.half], right) {
								used[next.gate] = true
								err := visit(bit+1, next.gate)
								used[next.gate] = false
								if err != nil {
									return err
								}
							}
							state.restore(carryCheckpoint)
						}
					}
					used[pair.sum], used[pair.half] = false, false
				}
				state.restore(checkpoint)
			}
		}
		return nil
	}
	if err := visit(1, adder.generate[0]); err != nil {
		return adderRepair{}, err
	}
	if answer == "" {
		return adderRepair{}, fmt.Errorf("adder has no four-swap repair")
	}
	return result, nil
}

func day24p02(r io.Reader) (string, error) {
	device, err := parseMonitoringDevice(r)
	if err != nil {
		return "", err
	}
	repair, err := repairRippleAdder(device, 2_000_000)
	if err != nil {
		return "", err
	}
	return strings.Join(repair.wires, ","), nil
}
