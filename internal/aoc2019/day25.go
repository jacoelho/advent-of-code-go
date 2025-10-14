// ABOUTME: Day 25 - Cryostasis: Interactive text adventure solver with automatic exploration
// ABOUTME: Uses BFS to explore ship, collect items, and brute force weight combinations

package aoc2019

import (
	"fmt"
	"io"
	"math/bits"
	"regexp"
	"slices"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

var (
	passwordRe = regexp.MustCompile(`\d{5,}`)
	roomRe     = regexp.MustCompile(`== (.+) ==`)
)

type GameState struct {
	roomName  string
	doors     []string
	items     []string
	inventory []string
}

var dangerousItems = collections.NewSet(
	"infinite loop",
	"escape pod",
	"molten lava",
	"photons",
	"giant electromagnet",
)

type queueItem struct {
	computer   *IntcodeComputer
	inventory  []string
	lastOutput string
}

func parseList(lines []string, header string) []string {
	var result []string
	for i, line := range lines {
		if !strings.HasPrefix(line, header) {
			continue
		}

		for j := i + 1; j < len(lines); j++ {
			listLine := lines[j]
			if after, ok := strings.CutPrefix(listLine, "- "); ok {
				result = append(result, after)
			} else if listLine == "" || !strings.HasPrefix(listLine, "-") {
				break
			}
		}
		return result
	}
	return result
}

func parseGameOutput(text string) GameState {
	state := GameState{}

	allMatches := roomRe.FindAllStringSubmatch(text, -1)
	if len(allMatches) > 0 {
		state.roomName = allMatches[len(allMatches)-1][1]
	}

	sections := strings.Split(text, "== ")
	if len(sections) < 2 {
		return state
	}
	lastSection := "== " + sections[len(sections)-1]

	lines := strings.Split(lastSection, "\n")

	state.doors = parseList(lines, "Doors here lead:")
	state.items = parseList(lines, "Items here:")
	state.inventory = parseList(lines, "Items in your inventory:")

	return state
}

func isDeath(text string) bool {
	deathPhrases := []string{"ejected", "are frozen", "crushed", "knocked"}
	return xslices.Any(func(phrase string) bool {
		return strings.Contains(text, phrase)
	}, deathPhrases)
}

func isSecurityCheck(text string) bool {
	return strings.Contains(text, "Alert") ||
		strings.Contains(text, "lighter") ||
		strings.Contains(text, "heavier")
}

func extractPassword(text string) string {
	return passwordRe.FindString(text)
}

func tryCommand(c *IntcodeComputer, command string) (string, error) {
	test := c.Clone()
	test.SetInput(StringsToASCII(command)...)
	if err := test.Run(); err != nil {
		return "", err
	}
	return test.ReadOutputString(), nil
}

func stateKey(roomName string, inventory []string) string {
	sortedInv := slices.Clone(inventory)
	slices.Sort(sortedInv)
	return fmt.Sprintf("%s:%v", roomName, sortedInv)
}

// mapShip performs BFS to explore all rooms, identify safe items, and locate the security checkpoint
func mapShip(computer *IntcodeComputer, initialOutput string) (safeItems collections.Set[string], checkpointDir string) {
	visited := make(map[string]bool)
	queue := collections.NewDeque[queueItem](16)
	queue.PushBack(queueItem{
		computer:   computer,
		inventory:  []string{},
		lastOutput: initialOutput,
	})

	items := collections.NewSet[string]()

	for current, ok := queue.PopFront(); ok; current, ok = queue.PopFront() {
		state := parseGameOutput(current.lastOutput)

		if visited[state.roomName] {
			continue
		}
		visited[state.roomName] = true

		for _, item := range state.items {
			if dangerousItems.Contains(item) {
				continue
			}

			testOutput, err := tryCommand(current.computer, "take "+item)
			if err != nil {
				continue
			}

			if isDeath(testOutput) {
				continue
			}

			items.Add(item)
		}

		if strings.Contains(state.roomName, "Security Checkpoint") && checkpointDir == "" {
			for _, door := range state.doors {
				testOutput, err := tryCommand(current.computer, door)
				if err != nil {
					continue
				}

				if isSecurityCheck(testOutput) {
					checkpointDir = door
					break
				}
			}

		}

		for _, door := range state.doors {
			testOutput, err := tryCommand(current.computer, door)
			if err != nil {
				continue
			}

			if isDeath(testOutput) || isSecurityCheck(testOutput) {
				continue
			}

			test := current.computer.Clone()
			test.SetInput(StringsToASCII(door)...)
			test.Run()

			queue.PushBack(queueItem{
				computer:   test,
				inventory:  append([]string{}, current.inventory...),
				lastOutput: testOutput,
			})
		}
	}

	return items, checkpointDir
}

// collectItemsAndReachCheckpoint performs BFS to collect all items and navigate to checkpoint
func collectItemsAndReachCheckpoint(program []int, safeItems collections.Set[string]) (*IntcodeComputer, error) {
	collector := New(program)
	if err := collector.Run(); err != nil {
		return nil, err
	}
	collectorOutput := collector.ReadAndResetOutputString()

	visitedStates := make(map[string]bool)
	collectionQueue := collections.NewDeque[queueItem](16)
	collectionQueue.PushBack(queueItem{
		computer:   collector,
		inventory:  []string{},
		lastOutput: collectorOutput,
	})

	statesExplored := 0
	for current, ok := collectionQueue.PopFront(); ok; current, ok = collectionQueue.PopFront() {
		state := parseGameOutput(current.lastOutput)

		stateKey := stateKey(state.roomName, current.inventory)
		if visitedStates[stateKey] {
			continue
		}
		visitedStates[stateKey] = true
		statesExplored++

		currentComputer := current.computer
		currentInventory := current.inventory

		for _, item := range state.items {
			if dangerousItems.Contains(item) || !safeItems.Contains(item) {
				continue
			}

			test := currentComputer.Clone()
			test.SetInput(StringsToASCII("take " + item)...)
			if err := test.Run(); err != nil {
				continue
			}

			currentComputer = test
			currentInventory = append(currentInventory, item)
		}

		if strings.Contains(state.roomName, "Security Checkpoint") && len(currentInventory) == safeItems.Len() {
			return currentComputer, nil
		}

		for _, door := range state.doors {
			testOutput, err := tryCommand(currentComputer, door)
			if err != nil {
				continue
			}

			if isDeath(testOutput) {
				continue
			}

			if isSecurityCheck(testOutput) && len(currentInventory) < safeItems.Len() {
				continue
			}

			test := currentComputer.Clone()
			test.SetInput(StringsToASCII(door)...)
			test.Run()

			collectionQueue.PushBack(queueItem{
				computer:   test,
				inventory:  slices.Clone(currentInventory),
				lastOutput: testOutput,
			})
		}
	}

	fmt.Printf("Explored %d states, never reached checkpoint with all items\n", statesExplored)
	return nil, nil
}

// exploreAndCollect performs BFS exploration using computer state cloning.
// It uses a two-pass strategy:
// 1. First pass: Map the entire ship, identify safe items, and locate the security checkpoint
// 2. Second pass: Collect all items and navigate to the checkpoint with full inventory
// Returns the computer state at checkpoint, set of safe items, and the direction to security
func exploreAndCollect(program []int) (*IntcodeComputer, collections.Set[string], string, error) {
	computer := New(program)
	if err := computer.Run(); err != nil {
		return nil, nil, "", err
	}
	initialOutput := computer.ReadAndResetOutputString()

	safeItems, checkpointDir := mapShip(computer, initialOutput)
	if checkpointDir == "" || safeItems.IsEmpty() {
		return nil, safeItems, checkpointDir, nil
	}

	checkpointComputer, err := collectItemsAndReachCheckpoint(program, safeItems)
	return checkpointComputer, safeItems, checkpointDir, err
}

func day25p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	checkpointComputer, itemsSet, securityDir, err := exploreAndCollect(program)
	if err != nil {
		return "", err
	}

	if checkpointComputer == nil || itemsSet.IsEmpty() || securityDir == "" {
		return "", fmt.Errorf("exploration failed: checkpoint=%v, items=%d, dir=%s", checkpointComputer != nil, itemsSet.Len(), securityDir)
	}

	items := slices.Collect(itemsSet.Iter())

	// Try all item combinations using Gray code to minimize state changes.
	// Gray code is a binary sequence where consecutive values differ by exactly one bit,
	// which means we only need to take/drop one item per iteration instead of 7+.
	// This reduces the number of commands from O(n * 2^n) to O(2^n).
	test := checkpointComputer.Clone()

	for _, item := range items {
		test.SetInput(StringsToASCII("drop " + item)...)
		test.Run()
		test.ClearOutput()
	}

	currentMask := 0
	numCombinations := 1 << len(items)

	for i := range numCombinations {
		gray := i ^ (i >> 1)
		diff := currentMask ^ gray
		currentMask = gray

		if diff != 0 {
			itemIdx := bits.TrailingZeros(uint(diff))
			item := items[itemIdx]
			action := "drop"
			if (gray & (1 << itemIdx)) != 0 {
				action = "take"
			}

			test.SetInput(StringsToASCII(action + " " + item)...)
			if err := test.Run(); err != nil {
				continue
			}
			test.ClearOutput()
		}

		test.SetInput(StringsToASCII(securityDir)...)
		if err := test.Run(); err != nil {
			continue
		}

		output := test.ReadOutputString()

		if isSecurityCheck(output) {
			continue
		}

		password := extractPassword(output)
		if password != "" {
			return password, nil
		}
	}

	return "", fmt.Errorf("tried all %d combinations, none worked", numCombinations)
}
