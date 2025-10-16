package aoc2023

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmaps"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

// Pulse represents a high or low pulse
type Pulse bool

const (
	LowPulse  Pulse = false
	HighPulse Pulse = true
)

// PulseMessage represents a pulse being sent between modules
type PulseMessage struct {
	From string
	To   string
	Type Pulse
}

// createPulseMessages creates pulse messages for all destinations
func createPulseMessages(from string, destinations []string, pulseType Pulse) []PulseMessage {
	return xslices.Map(func(dest string) PulseMessage {
		return PulseMessage{From: from, To: dest, Type: pulseType}
	}, destinations)
}

// Module represents a communication module
type Module interface {
	Name() string
	ProcessPulse(from string, pulse Pulse) []PulseMessage
	Destinations() []string
}

// FlipFlopModule represents a flip-flop module (%)
type FlipFlopModule struct {
	name         string
	destinations []string
	on           bool
}

func (f *FlipFlopModule) Name() string {
	return f.name
}

func (f *FlipFlopModule) Destinations() []string {
	return f.destinations
}

func (f *FlipFlopModule) ProcessPulse(from string, pulse Pulse) []PulseMessage {
	if pulse == HighPulse {
		return nil // Ignore high pulses
	}

	f.on = !f.on
	outputPulse := LowPulse
	if f.on {
		outputPulse = HighPulse
	}

	return createPulseMessages(f.name, f.destinations, outputPulse)
}

// ConjunctionModule represents a conjunction module (&)
type ConjunctionModule struct {
	name         string
	destinations []string
	memory       map[string]Pulse
}

func (c *ConjunctionModule) Name() string {
	return c.name
}

func (c *ConjunctionModule) Destinations() []string {
	return c.destinations
}

func (c *ConjunctionModule) ProcessPulse(from string, pulse Pulse) []PulseMessage {
	c.memory[from] = pulse

	pulseValues := slices.Collect(maps.Values(c.memory))

	allHigh := xslices.Every(func(p Pulse) bool { return p == HighPulse }, pulseValues)

	outputPulse := LowPulse
	if !allHigh {
		outputPulse = HighPulse
	}

	return createPulseMessages(c.name, c.destinations, outputPulse)
}

// BroadcasterModule represents the broadcaster module
type BroadcasterModule struct {
	name         string
	destinations []string
}

func (b *BroadcasterModule) Name() string {
	return b.name
}

func (b *BroadcasterModule) Destinations() []string {
	return b.destinations
}

func (b *BroadcasterModule) ProcessPulse(from string, pulse Pulse) []PulseMessage {
	return createPulseMessages(b.name, b.destinations, pulse)
}

// parseModuleLine parses a single line into a module
func parseModuleLine(line []byte) (Module, error) {
	lineStr := strings.TrimSpace(string(line))
	if lineStr == "" {
		return nil, nil
	}

	parts := strings.Split(lineStr, " -> ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid line format: %s", lineStr)
	}

	moduleSpec := parts[0]
	destinations := strings.Split(strings.ReplaceAll(parts[1], " ", ""), ",")

	switch {
	case moduleSpec == "broadcaster":
		return &BroadcasterModule{
			name:         "broadcaster",
			destinations: destinations,
		}, nil
	case strings.HasPrefix(moduleSpec, "%"):
		return &FlipFlopModule{
			name:         moduleSpec[1:],
			destinations: destinations,
		}, nil
	case strings.HasPrefix(moduleSpec, "&"):
		return &ConjunctionModule{
			name:         moduleSpec[1:],
			destinations: destinations,
			memory:       make(map[string]Pulse),
		}, nil
	default:
		return nil, fmt.Errorf("unknown module type: %s", moduleSpec)
	}
}

// initConjunctionInputs initializes conjunction memory with low pulses for all inputs
func initConjunctionInputs(modules map[string]Module) {
	// identify conjunction modules
	conjunctions := make(map[string]*ConjunctionModule)
	for _, module := range modules {
		if conj, ok := module.(*ConjunctionModule); ok {
			conjunctions[conj.name] = conj
		}
	}

	// identify inputs for each conjunction
	for moduleName, module := range modules {
		for _, dest := range module.Destinations() {
			if conj, exists := conjunctions[dest]; exists {
				conj.memory[moduleName] = LowPulse
			}
		}
	}
}

// parseModules parses the input and builds a map of modules
func parseModules(r io.Reader) (map[string]Module, error) {
	modules := make(map[string]Module)

	sc := scanner.NewScanner(r, parseModuleLine)
	for module := range sc.Values() {
		if module == nil {
			continue
		}
		modules[module.Name()] = module
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	initConjunctionInputs(modules)
	return modules, nil
}

// findModuleFeeding finds which module feeds the target module and returns its name
func findModuleFeeding(modules map[string]Module, target string) string {
	result, found := xmaps.Find(modules, func(_ string, module Module) bool {
		return slices.Contains(module.Destinations(), target)
	})
	if found {
		return result.K
	}
	return ""
}

// findInputs finds all modules that feed the given target module
func findInputs(modules map[string]Module, target string) []string {
	matches := xmaps.Filter(func(_ string, module Module) bool {
		return slices.Contains(module.Destinations(), target)
	}, modules)

	return xslices.Map(func(p xmaps.Pair[string, Module]) string {
		return p.K
	}, matches)
}

// simulateButtonPress simulates one button press and calls onPulse for each pulse
func simulateButtonPress(modules map[string]Module, onPulse func(PulseMessage)) {
	queue := collections.NewDeque[PulseMessage](100)
	queue.PushBack(PulseMessage{From: "button", To: "broadcaster", Type: LowPulse})

	for pulse, ok := queue.PopFront(); ok; pulse, ok = queue.PopFront() {
		onPulse(pulse)
		if module, exists := modules[pulse.To]; exists {
			for _, newPulse := range module.ProcessPulse(pulse.From, pulse.Type) {
				queue.PushBack(newPulse)
			}
		}
	}
}

func day20p01(r io.Reader) (string, error) {
	modules, err := parseModules(r)
	if err != nil {
		return "", err
	}

	var lowPulses, highPulses int64
	for range 1000 {
		simulateButtonPress(modules, func(pulse PulseMessage) {
			if pulse.Type == LowPulse {
				lowPulses++
			} else {
				highPulses++
			}
		})
	}

	return strconv.FormatInt(lowPulses*highPulses, 10), nil
}

func day20p02(r io.Reader) (string, error) {
	modules, err := parseModules(r)
	if err != nil {
		return "", err
	}

	// identify the module that feeds rx
	rxFeeder := findModuleFeeding(modules, "rx")
	if rxFeeder == "" {
		return "", fmt.Errorf("no module found feeding rx")
	}

	// identify all inputs to the rx feeder
	rxFeederInputs := findInputs(modules, rxFeeder)
	if len(rxFeederInputs) == 0 {
		return "", fmt.Errorf("no inputs found for rx feeder %s", rxFeeder)
	}

	// track when each input sends HIGH to the rx feeder
	cycleLengths := make(map[string]int)
	buttonPress := 0

	// simulate until we find the first HIGH pulse from each input
	for len(cycleLengths) < len(rxFeederInputs) && buttonPress < 100000 {
		buttonPress++
		simulateButtonPress(modules, func(pulse PulseMessage) {
			// check if this is a HIGH pulse from one of our target inputs to the rx feeder
			if pulse.Type == HighPulse && pulse.To == rxFeeder {
				for _, input := range rxFeederInputs {
					if pulse.From == input && cycleLengths[input] == 0 {
						cycleLengths[input] = buttonPress
					}
				}
			}
		})
	}

	// check if we found all cycle lengths
	if len(cycleLengths) != len(rxFeederInputs) {
		return "", fmt.Errorf("could not find all cycle lengths after %d button presses", buttonPress)
	}

	// calculate LCM of all cycle lengths
	lengths := slices.Collect(maps.Values(cycleLengths))

	if len(lengths) == 1 {
		return strconv.Itoa(lengths[0]), nil
	}

	lcm := xmath.LCM(lengths[0], lengths[1:]...)
	return strconv.Itoa(lcm), nil
}
