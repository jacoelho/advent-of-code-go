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

const (
	maxPulsePresses       = 100_000
	maxPulsesPerPress     = 1_000_000
	maxQueuedPulses       = 100_000
	maxPulseEvents        = 10_000_000
	maxPulseStorage       = 64 << 20
	maxPulseAlignmentWork = 10_000_000
	maxPulseResidues      = 10_000
	maxPulseCuts          = 10_000
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

func parsePulseNetwork(r io.Reader) (*pulseNetwork, error) {
	const inputLimit = 4 << 20
	data, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return nil, err
	}
	if len(data) > inputLimit {
		return nil, fmt.Errorf("pulse input exceeds 4 MiB")
	}
	network := &pulseNetwork{broadcaster: -1, rx: -1}
	indices := make(map[string]int)
	index := func(name string) (int, error) {
		if !validPulseName(name) {
			return 0, fmt.Errorf("invalid pulse module name %q", name)
		}
		if idx, exists := indices[name]; exists {
			return idx, nil
		}
		if len(network.modules) == 4096 {
			return 0, fmt.Errorf("pulse network exceeds 4096 modules and sinks")
		}
		idx := len(network.modules)
		indices[name] = idx
		network.modules = append(network.modules, pulseModule{})
		return idx, nil
	}

	s := scanner.NewScanner(strings.NewReader(string(data)), func(line []byte) (string, error) {
		return string(line), nil
	})
	edges := 0
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
			edges++
			if edges > 16384 {
				return nil, fmt.Errorf("pulse network exceeds 16384 edges")
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

func (n *pulseNetwork) press(remaining *int, observe func(pulseMessage) error) error {
	queue := collections.NewDeque[pulseMessage](100)
	queue.PushBack(pulseMessage{from: -1, to: n.broadcaster})
	processed := 0
	for message, ok := queue.PopFront(); ok; message, ok = queue.PopFront() {
		if *remaining == 0 || processed == maxPulsesPerPress {
			return fmt.Errorf("pulse simulation event limit exceeded")
		}
		*remaining--
		processed++
		if err := observe(message); err != nil {
			return err
		}
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
			if queue.Size() == maxQueuedPulses {
				return fmt.Errorf("pulse queue exceeds %d events", maxQueuedPulses)
			}
			queue.PushBack(pulseMessage{from: message.to, to: to, depth: message.depth + 1, high: high})
		}
	}
	return nil
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

func chargePulseStorage(remaining *int, size int) error {
	if size > *remaining {
		return fmt.Errorf("pulse state and trace storage limit exceeded")
	}
	*remaining -= size
	return nil
}

func firstLowBySimulation(n *pulseNetwork) (*big.Int, error) {
	modules := make([]int, len(n.modules))
	for i := range modules {
		modules[i] = i
	}
	seen := make(map[string]bool)
	remaining, storage := maxPulseEvents, maxPulseStorage
	for press := 1; press <= maxPulsePresses; press++ {
		key := string(n.state(modules))
		if seen[key] {
			return nil, fmt.Errorf("rx never receives a low pulse")
		}
		if err := chargePulseStorage(&storage, len(key)+64); err != nil {
			return nil, err
		}
		seen[key] = true
		low := false
		err := n.press(&remaining, func(message pulseMessage) error {
			low = low || message.to == n.rx && !message.high
			return nil
		})
		if err != nil {
			return nil, err
		}
		if low {
			return big.NewInt(int64(press)), nil
		}
	}
	return nil, fmt.Errorf("pulse simulation exceeds %d presses", maxPulsePresses)
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
	remaining, storage := maxPulseEvents, maxPulseStorage
	sources := make(map[int]int, len(branches))
	for i, branch := range branches {
		sources[n.modules[gate].inputs[branch.gateSlot]] = i
	}
	for boundary := 0; boundary <= maxPulsePresses; boundary++ {
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
			if err := chargePulseStorage(&storage, len(key)+64+32); err != nil {
				return nil, err
			}
			branch.seen[key] = boundary
			branch.traces = append(branch.traces, pulseTrace{startHigh: memory != 0})
			complete = false
		}
		if complete {
			return alignPulseBranches(branches, boundary)
		}
		if boundary == maxPulsePresses {
			break
		}
		low := false
		err := n.press(&remaining, func(message pulseMessage) error {
			low = low || message.to == n.rx && !message.high
			if message.to != gate {
				return nil
			}
			branch := &branches[sources[message.from]]
			if branch.period != 0 {
				return nil
			}
			if err := chargePulseStorage(&storage, 16); err != nil {
				return err
			}
			trace := &branch.traces[len(branch.traces)-1]
			trace.events = append(trace.events, gatePulse{depth: message.depth, high: message.high})
			return nil
		})
		if err != nil {
			return nil, err
		}
		if low {
			return big.NewInt(int64(boundary + 1)), nil
		}
	}
	return nil, fmt.Errorf("pulse branch cycles exceed %d presses", maxPulsePresses)
}

type pulseCut struct {
	branch int
	depth  int
}

func pulseCutResidues(branch pulseBranch, index int, cut pulseCut, work *int) ([]int, error) {
	var residues []int
	for phase, trace := range branch.traces[branch.cycleStart : branch.cycleStart+branch.period] {
		if *work == 0 {
			return nil, fmt.Errorf("pulse alignment work limit exceeded")
		}
		*work--
		high := trace.startHigh
		if index == cut.branch {
			high = false
		}
		for _, event := range trace.events {
			if *work == 0 {
				return nil, fmt.Errorf("pulse alignment work limit exceeded")
			}
			*work--
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
	return residues, nil
}

func alignPulseCut(branches []pulseBranch, cut pulseCut, minimum *big.Int, work *int) (*big.Int, error) {
	modulus := big.NewInt(1)
	residues := []*big.Int{new(big.Int)}
	for i, branch := range branches {
		allowed, err := pulseCutResidues(branch, i, cut, work)
		if err != nil {
			return nil, err
		}
		if len(allowed) == 0 {
			return nil, nil
		}
		period := big.NewInt(int64(branch.period))
		var nextModulus *big.Int
		var next []*big.Int
		seen := make(map[string]bool)
		for _, residue := range residues {
			for _, phase := range allowed {
				if *work == 0 {
					return nil, fmt.Errorf("pulse alignment work limit exceeded")
				}
				*work--
				combined, combinedModulus, compatible := combineCongruences(residue, modulus, big.NewInt(int64(phase)), period)
				if !compatible {
					continue
				}
				nextModulus = combinedModulus
				key := combined.String()
				if seen[key] {
					continue
				}
				if len(next) == maxPulseResidues {
					return nil, fmt.Errorf("pulse alignment exceeds %d residues", maxPulseResidues)
				}
				seen[key] = true
				next = append(next, combined)
			}
		}
		if len(next) == 0 {
			return nil, nil
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
	return answer, nil
}

func alignPulseBranches(branches []pulseBranch, simulated int) (*big.Int, error) {
	work := maxPulseAlignmentWork
	seen := make(map[pulseCut]bool)
	var cuts []pulseCut
	for i, branch := range branches {
		for _, trace := range branch.traces[branch.cycleStart : branch.cycleStart+branch.period] {
			for _, event := range trace.events {
				if work == 0 {
					return nil, fmt.Errorf("pulse alignment work limit exceeded")
				}
				work--
				cut := pulseCut{branch: i, depth: event.depth}
				if event.high && !seen[cut] {
					if len(cuts) == maxPulseCuts {
						return nil, fmt.Errorf("pulse alignment exceeds %d event cuts", maxPulseCuts)
					}
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
		candidate, err := alignPulseCut(branches, cut, minimum, &work)
		if err != nil {
			return nil, err
		}
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
	remaining := maxPulseEvents
	var low, high int64
	for range 1000 {
		err := network.press(&remaining, func(message pulseMessage) error {
			if message.high {
				high++
			} else {
				low++
			}
			return nil
		})
		if err != nil {
			return "", err
		}
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
