package aoc2023

import (
	"fmt"
	"io"
	"maps"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

const (
	maxNavigationStates = 1_000_000
	maxAlignmentWork    = 1_000_000
	maxAlignedResidues  = 10_000
)

type network map[string][2]string

type networkEntry struct {
	node string
	dest [2]string
}

type navigationState struct {
	node        string
	instruction int
}

type navigationSchedule struct {
	cycleStart int64
	period     int64
	transient  []int64
	residues   []int64
}

func parseNetworkEntry(line []byte) (networkEntry, error) {
	node, body, ok := strings.Cut(string(line), "=")
	if !ok {
		return networkEntry{}, fmt.Errorf("invalid network entry %q", line)
	}
	node, body = strings.TrimSpace(node), strings.TrimSpace(body)
	body, left := strings.CutPrefix(body, "(")
	body, right := strings.CutSuffix(body, ")")
	destinations := strings.Split(body, ",")
	if !left || !right || len(destinations) != 2 {
		return networkEntry{}, fmt.Errorf("invalid network destinations %q", line)
	}
	entry := networkEntry{node: node, dest: [2]string{strings.TrimSpace(destinations[0]), strings.TrimSpace(destinations[1])}}
	for _, name := range []string{entry.node, entry.dest[0], entry.dest[1]} {
		if len(name) == 0 || len(name) > 128 {
			return networkEntry{}, fmt.Errorf("network name must have 1 to 128 characters")
		}
		for _, c := range name {
			if !(c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9') {
				return networkEntry{}, fmt.Errorf("invalid network name %q", name)
			}
		}
	}
	return entry, nil
}

func parseNetwork(r io.Reader) (string, network, error) {
	const inputLimit = 4 << 20
	data, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return "", nil, err
	}
	if len(data) > inputLimit {
		return "", nil, fmt.Errorf("network input exceeds 4 MiB")
	}
	input := strings.ReplaceAll(string(data), "\r\n", "\n")
	instructions, entries, ok := strings.Cut(strings.TrimSpace(input), "\n\n")
	if !ok || len(instructions) == 0 || len(instructions) > 100_000 {
		return "", nil, fmt.Errorf("expected 1 to 100000 instructions and a network")
	}
	for _, instruction := range instructions {
		if instruction != 'L' && instruction != 'R' {
			return "", nil, fmt.Errorf("invalid navigation instruction %q", instruction)
		}
	}
	s := scanner.NewScanner(strings.NewReader(entries), parseNetworkEntry)
	net := make(network)
	for entry := range s.Values() {
		if _, exists := net[entry.node]; exists {
			return "", nil, fmt.Errorf("duplicate network node %q", entry.node)
		}
		if len(net) == 100_000 {
			return "", nil, fmt.Errorf("network exceeds 100000 nodes")
		}
		net[entry.node] = entry.dest
	}
	if err := s.Err(); err != nil {
		return "", nil, err
	}
	if len(net) == 0 {
		return "", nil, fmt.Errorf("empty network")
	}
	for node, destinations := range net {
		for _, destination := range destinations {
			if _, exists := net[destination]; !exists {
				return "", nil, fmt.Errorf("node %s has unknown destination %s", node, destination)
			}
		}
	}
	return instructions, net, nil
}

func (state navigationState) next(instructions string, net network) navigationState {
	branch := 0
	if instructions[state.instruction] == 'R' {
		branch = 1
	}
	return navigationState{
		node:        net[state.node][branch],
		instruction: (state.instruction + 1) % len(instructions),
	}
}

func analyzeNavigation(instructions string, net network, start string, remaining *int) (navigationSchedule, error) {
	seen := make(map[navigationState]int64)
	var hits []int64
	state := navigationState{node: start}
	for step := int64(0); ; step++ {
		if previous, exists := seen[state]; exists {
			schedule := navigationSchedule{cycleStart: previous, period: step - previous}
			for _, hit := range hits {
				if hit < previous {
					schedule.transient = append(schedule.transient, hit)
				} else {
					schedule.residues = append(schedule.residues, hit%schedule.period)
				}
			}
			slices.Sort(schedule.residues)
			return schedule, nil
		}
		if *remaining == 0 {
			return navigationSchedule{}, fmt.Errorf("navigation exceeds one million states")
		}
		*remaining -= 1
		seen[state] = step
		if strings.HasSuffix(state.node, "Z") {
			hits = append(hits, step)
		}
		state = state.next(instructions, net)
	}
}

func (schedule navigationSchedule) contains(step int64) bool {
	if step < schedule.cycleStart {
		_, found := slices.BinarySearch(schedule.transient, step)
		return found
	}
	_, found := slices.BinarySearch(schedule.residues, step%schedule.period)
	return found
}

func alignNavigation(schedules []navigationSchedule) (*big.Int, error) {
	latest := 0
	for i := range schedules {
		if schedules[i].cycleStart > schedules[latest].cycleStart {
			latest = i
		}
	}
	work := maxAlignmentWork
	// Before every walk is periodic, a common hit must be transient in this walk.
	for _, step := range schedules[latest].transient {
		shared := true
		for _, schedule := range schedules {
			if work == 0 {
				return nil, fmt.Errorf("navigation alignment work limit exceeded")
			}
			work--
			if !schedule.contains(step) {
				shared = false
				break
			}
		}
		if shared {
			return big.NewInt(step), nil
		}
	}

	period := big.NewInt(1)
	residues := []*big.Int{new(big.Int)}
	for _, schedule := range schedules {
		if len(schedule.residues) == 0 {
			return nil, fmt.Errorf("navigation has no common target time")
		}
		modulus := big.NewInt(schedule.period)
		var nextPeriod *big.Int
		seen := make(map[string]bool)
		var next []*big.Int
		for _, residue := range residues {
			for _, target := range schedule.residues {
				if work == 0 {
					return nil, fmt.Errorf("navigation alignment work limit exceeded")
				}
				work--
				combined, combinedPeriod, compatible := combineCongruences(residue, period, big.NewInt(target), modulus)
				if !compatible {
					continue
				}
				nextPeriod = combinedPeriod
				key := combined.String()
				if seen[key] {
					continue
				}
				if len(next) == maxAlignedResidues {
					return nil, fmt.Errorf("navigation exceeds 10000 aligned residues")
				}
				seen[key] = true
				next = append(next, combined)
			}
		}
		if len(next) == 0 {
			return nil, fmt.Errorf("navigation target schedules are incompatible")
		}
		period, residues = nextPeriod, next
	}

	minimum := big.NewInt(schedules[latest].cycleStart)
	var answer *big.Int
	for _, residue := range residues {
		candidate := firstCongruentAtLeast(residue, period, minimum)
		if answer == nil || candidate.Cmp(answer) < 0 {
			answer = candidate
		}
	}
	return answer, nil
}

func day08p01(r io.Reader) (string, error) {
	instructions, net, err := parseNetwork(r)
	if err != nil {
		return "", err
	}
	for _, node := range []string{"AAA", "ZZZ"} {
		if _, exists := net[node]; !exists {
			return "", fmt.Errorf("missing navigation node %s", node)
		}
	}
	state := navigationState{node: "AAA"}
	seen := make(map[navigationState]bool)
	for step := 0; ; step++ {
		if state.node == "ZZZ" {
			return strconv.Itoa(step), nil
		}
		if seen[state] {
			return "", fmt.Errorf("ZZZ is unreachable")
		}
		if len(seen) == maxNavigationStates {
			return "", fmt.Errorf("navigation exceeds one million states")
		}
		seen[state] = true
		state = state.next(instructions, net)
	}
}

func day08p02(r io.Reader) (string, error) {
	instructions, net, err := parseNetwork(r)
	if err != nil {
		return "", err
	}
	starts := xiter.From(maps.Keys(net)).Filter(func(node string) bool {
		return strings.HasSuffix(node, "A")
	}).Collect()
	if len(starts) == 0 || len(starts) > 1024 {
		return "", fmt.Errorf("expected 1 to 1024 navigation starts")
	}
	slices.Sort(starts)
	remaining := maxNavigationStates
	schedules := make([]navigationSchedule, 0, len(starts))
	for _, start := range starts {
		schedule, err := analyzeNavigation(instructions, net, start, &remaining)
		if err != nil {
			return "", err
		}
		schedules = append(schedules, schedule)
	}
	answer, err := alignNavigation(schedules)
	if err != nil {
		return "", err
	}
	return answer.String(), nil
}
