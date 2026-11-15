package aoc2023

import (
	"cmp"
	"fmt"
	"io"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type pulseModuleKind uint8

const (
	pulseSink pulseModuleKind = iota
	pulseBroadcaster
	pulseFlipFlop
	pulseConjunction
)

type pulseModule struct {
	kind         pulseModuleKind
	destinations []int
	inputs       []int
	memory       []bool
	on           bool
}

type pulseNetwork struct {
	modules     []pulseModule
	broadcaster int
	rx          int
}

type pulseMessage struct {
	from, to int
	depth    int
	high     bool
}

func validPulseName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for _, c := range name {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func parsePulseNetwork(r io.Reader) (*pulseNetwork, error) {
	network := &pulseNetwork{broadcaster: -1, rx: -1}
	indices := make(map[string]int)
	index := func(name string) (int, error) {
		if !validPulseName(name) {
			return 0, fmt.Errorf("invalid pulse module name %q", name)
		}
		if idx, exists := indices[name]; exists {
			return idx, nil
		}
		idx := len(network.modules)
		indices[name] = idx
		network.modules = append(network.modules, pulseModule{})
		return idx, nil
	}

	s := scanner.NewScanner(r, func(line []byte) (string, error) {
		return string(line), nil
	})
	for line := range s.Values() {
		spec, destinations, ok := strings.Cut(line, " -> ")
		if !ok || len(spec) == 0 {
			return nil, fmt.Errorf("invalid pulse module %q", line)
		}
		name, kind := spec, pulseBroadcaster
		if spec != "broadcaster" {
			name = spec[1:]
			switch spec[0] {
			case '%':
				kind = pulseFlipFlop
			case '&':
				kind = pulseConjunction
			default:
				return nil, fmt.Errorf("unknown pulse module type %q", spec)
			}
			if name == "broadcaster" {
				return nil, fmt.Errorf("broadcaster cannot have a module prefix")
			}
		}
		from, err := index(name)
		if err != nil {
			return nil, err
		}
		if network.modules[from].kind != pulseSink {
			return nil, fmt.Errorf("duplicate pulse module %s", name)
		}
		network.modules[from].kind = kind
		for destination := range strings.SplitSeq(destinations, ",") {
			to, err := index(strings.TrimSpace(destination))
			if err != nil {
				return nil, err
			}
			network.modules[from].destinations = append(network.modules[from].destinations, to)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	if idx, exists := indices["broadcaster"]; exists && network.modules[idx].kind == pulseBroadcaster {
		network.broadcaster = idx
	} else {
		return nil, fmt.Errorf("missing broadcaster")
	}
	if idx, exists := indices["rx"]; exists {
		network.rx = idx
	}

	for from, module := range network.modules {
		for _, to := range module.destinations {
			destination := &network.modules[to]
			if destination.kind == pulseConjunction && !slices.Contains(destination.inputs, from) {
				destination.inputs = append(destination.inputs, from)
				destination.memory = append(destination.memory, false)
			}
		}
	}
	return network, nil
}

// press sends one button pulse and observes every pulse until the queue drains.
// A network whose pulses never die out makes press run forever.
func (n *pulseNetwork) press(observe func(pulseMessage)) {
	queue := collections.NewDeque[pulseMessage](100)
	queue.PushBack(pulseMessage{from: -1, to: n.broadcaster})
	for message, ok := queue.PopFront(); ok; message, ok = queue.PopFront() {
		observe(message)
		module := &n.modules[message.to]
		high := message.high
		switch module.kind {
		case pulseSink:
			continue
		case pulseFlipFlop:
			if high {
				continue
			}
			module.on = !module.on
			high = module.on
		case pulseConjunction:
			module.memory[slices.Index(module.inputs, message.from)] = high
			high = slices.Contains(module.memory, false)
		}
		for _, to := range module.destinations {
			queue.PushBack(pulseMessage{from: message.to, to: to, depth: message.depth + 1, high: high})
		}
	}
}

func (n *pulseNetwork) state(modules []int) []byte {
	var state []byte
	appendBit := func(high bool) {
		value := byte(0)
		if high {
			value = 1
		}
		state = append(state, value)
	}
	for _, idx := range modules {
		module := &n.modules[idx]
		if module.kind == pulseFlipFlop {
			appendBit(module.on)
		}
		for _, high := range module.memory {
			appendBit(high)
		}
	}
	return state
}

func firstLowBySimulation(n *pulseNetwork) (*big.Int, error) {
	modules := make([]int, len(n.modules))
	for i := range modules {
		modules[i] = i
	}
	seen := make(map[string]bool)
	for press := 1; ; press++ {
		key := string(n.state(modules))
		if seen[key] {
			return nil, fmt.Errorf("rx never receives a low pulse")
		}
		seen[key] = true
		low := false
		n.press(func(message pulseMessage) {
			low = low || message.to == n.rx && !message.high
		})
		if low {
			return big.NewInt(int64(press)), nil
		}
	}
}

type gatePulse struct {
	depth int
	high  bool
}

type pulseTrace struct {
	startHigh bool
	events    []gatePulse
}

type pulseBranch struct {
	modules    []int
	gateSlot   int
	seen       map[string]int
	traces     []pulseTrace
	cycleStart int
	period     int
}

func independentPulseBranches(n *pulseNetwork) (int, []pulseBranch, bool) {
	if n.rx < 0 || n.modules[n.rx].kind != pulseSink {
		return 0, nil, false
	}
	gate := -1
	for idx, module := range n.modules {
		if slices.Contains(module.destinations, n.broadcaster) {
			return 0, nil, false
		}
		if slices.Contains(module.destinations, n.rx) {
			if gate >= 0 {
				return 0, nil, false
			}
			gate = idx
		}
	}
	if gate < 0 || n.modules[gate].kind != pulseConjunction {
		return 0, nil, false
	}
	for _, to := range n.modules[gate].destinations {
		if n.modules[to].kind != pulseSink {
			return 0, nil, false
		}
	}
	roots := n.modules[n.broadcaster].destinations
	if len(n.modules[gate].inputs) != len(roots) {
		return 0, nil, false
	}
	owner := make([]int, len(n.modules))
	for i := range owner {
		owner[i] = -1
	}
	branches := make([]pulseBranch, len(roots))
	for i, root := range roots {
		if root == gate || n.modules[root].kind == pulseSink || owner[root] != -1 {
			return 0, nil, false
		}
		branch := &branches[i]
		branch.gateSlot = -1
		branch.seen = make(map[string]int)
		pending := []int{root}
		for len(pending) > 0 {
			idx := pending[len(pending)-1]
			pending = pending[:len(pending)-1]
			if idx == gate || n.modules[idx].kind == pulseSink {
				continue
			}
			if owner[idx] == i {
				continue
			}
			if owner[idx] != -1 || idx == n.broadcaster {
				return 0, nil, false
			}
			owner[idx] = i
			branch.modules = append(branch.modules, idx)
			pending = append(pending, n.modules[idx].destinations...)
		}
		slices.Sort(branch.modules)
	}
	for idx, module := range n.modules {
		if module.kind != pulseSink && idx != gate && idx != n.broadcaster && owner[idx] < 0 {
			return 0, nil, false
		}
	}
	for slot, source := range n.modules[gate].inputs {
		idx := owner[source]
		if idx < 0 || branches[idx].gateSlot >= 0 {
			return 0, nil, false
		}
		branches[idx].gateSlot = slot
	}
	return gate, branches, true
}

func firstLowByBranches(n *pulseNetwork, gate int, branches []pulseBranch) (*big.Int, error) {
	sources := make(map[int]int, len(branches))
	for i, branch := range branches {
		sources[n.modules[gate].inputs[branch.gateSlot]] = i
	}
	for boundary := 0; ; boundary++ {
		complete := true
		for i := range branches {
			branch := &branches[i]
			if branch.period != 0 {
				continue
			}
			state := n.state(branch.modules)
			memory := byte(0)
			if n.modules[gate].memory[branch.gateSlot] {
				memory = 1
			}
			key := string(append(state, memory))
			if previous, exists := branch.seen[key]; exists {
				branch.cycleStart, branch.period = previous, boundary-previous
				branch.seen = nil
				continue
			}
			branch.seen[key] = boundary
			branch.traces = append(branch.traces, pulseTrace{startHigh: memory != 0})
			complete = false
		}
		if complete {
			return alignPulseBranches(branches, boundary)
		}
		low := false
		n.press(func(message pulseMessage) {
			low = low || message.to == n.rx && !message.high
			if message.to != gate {
				return
			}
			branch := &branches[sources[message.from]]
			if branch.period != 0 {
				return
			}
			trace := &branch.traces[len(branch.traces)-1]
			trace.events = append(trace.events, gatePulse{depth: message.depth, high: message.high})
		})
		if low {
			return big.NewInt(int64(boundary + 1)), nil
		}
	}
}

type pulseCut struct {
	branch int
	depth  int
}

func pulseCutResidues(branch pulseBranch, index int, cut pulseCut) []int {
	var residues []int
	for phase, trace := range branch.traces[branch.cycleStart : branch.cycleStart+branch.period] {
		high := trace.startHigh
		if index == cut.branch {
			high = false
		}
		for _, event := range trace.events {
			if index == cut.branch {
				high = high || event.depth == cut.depth && event.high
			} else if event.depth < cut.depth || event.depth == cut.depth && index < cut.branch {
				high = event.high
			} else {
				break
			}
		}
		if high {
			residues = append(residues, (branch.cycleStart+phase+1)%branch.period)
		}
	}
	return residues
}

func alignPulseCut(branches []pulseBranch, cut pulseCut, minimum *big.Int) *big.Int {
	modulus := big.NewInt(1)
	residues := []*big.Int{new(big.Int)}
	for i, branch := range branches {
		allowed := pulseCutResidues(branch, i, cut)
		if len(allowed) == 0 {
			return nil
		}
		period := big.NewInt(int64(branch.period))
		var nextModulus *big.Int
		var next []*big.Int
		seen := make(map[string]bool)
		for _, residue := range residues {
			for _, phase := range allowed {
				combined, combinedModulus, compatible := combineCongruences(residue, modulus, big.NewInt(int64(phase)), period)
				if !compatible {
					continue
				}
				nextModulus = combinedModulus
				key := combined.String()
				if seen[key] {
					continue
				}
				seen[key] = true
				next = append(next, combined)
			}
		}
		if len(next) == 0 {
			return nil
		}
		modulus, residues = nextModulus, next
	}
	var answer *big.Int
	for _, residue := range residues {
		candidate := firstCongruentAtLeast(residue, modulus, minimum)
		if answer == nil || candidate.Cmp(answer) < 0 {
			answer = candidate
		}
	}
	return answer
}

func alignPulseBranches(branches []pulseBranch, simulated int) (*big.Int, error) {
	seen := make(map[pulseCut]bool)
	var cuts []pulseCut
	for i, branch := range branches {
		for _, trace := range branch.traces[branch.cycleStart : branch.cycleStart+branch.period] {
			for _, event := range trace.events {
				cut := pulseCut{branch: i, depth: event.depth}
				if event.high && !seen[cut] {
					seen[cut] = true
					cuts = append(cuts, cut)
				}
			}
		}
	}
	slices.SortFunc(cuts, func(a, b pulseCut) int {
		return cmp.Or(cmp.Compare(a.branch, b.branch), cmp.Compare(a.depth, b.depth))
	})
	minimum := big.NewInt(int64(simulated + 1))
	var answer *big.Int
	// FIFO orders independent branches by depth, then broadcaster-root order.
	// Each cut tests the conjunction memory immediately after a triggering HIGH.
	for _, cut := range cuts {
		candidate := alignPulseCut(branches, cut, minimum)
		if candidate != nil && (answer == nil || candidate.Cmp(answer) < 0) {
			answer = candidate
		}
	}
	if answer == nil {
		return nil, fmt.Errorf("rx never receives a low pulse")
	}
	return answer, nil
}

func day20p01(r io.Reader) (string, error) {
	network, err := parsePulseNetwork(r)
	if err != nil {
		return "", err
	}
	var low, high int64
	for range 1000 {
		network.press(func(message pulseMessage) {
			if message.high {
				high++
			} else {
				low++
			}
		})
	}
	return strconv.FormatInt(low*high, 10), nil
}

func day20p02(r io.Reader) (string, error) {
	network, err := parsePulseNetwork(r)
	if err != nil {
		return "", err
	}
	if network.rx < 0 {
		return "", fmt.Errorf("no module feeds rx")
	}
	var answer *big.Int
	if gate, branches, supported := independentPulseBranches(network); supported {
		answer, err = firstLowByBranches(network, gate, branches)
	} else {
		answer, err = firstLowBySimulation(network)
	}
	if err != nil {
		return "", err
	}
	return answer.String(), nil
}
