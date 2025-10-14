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
	var foundHeader bool
	for i, line := range lines {
		if strings.HasPrefix(line, header) {
			foundHeader = true
			for j := i + 1; j < len(lines); j++ {
				listLine := lines[j]
				if after, ok := strings.CutPrefix(listLine, "- "); ok {
					result = append(result, after)
				} else if listLine == "" || !strings.HasPrefix(listLine, "-") {
					break
				}
			}
			break
		}
	}
	if !foundHeader {
		return nil
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

func isSafeToTake(computer *IntcodeComputer, item string) bool {
	if dangerousItems.Contains(item) {
		return false
	}
	output, err := tryCommand(computer, "take "+item)
	return err == nil && !isDeath(output)
}

func findCheckpointDirection(computer *IntcodeComputer, doors []string) string {
	for _, door := range doors {
		output, err := tryCommand(computer, door)
		if err == nil && isSecurityCheck(output) {
			return door
		}
	}
	return ""
}

func exploreDirection(computer *IntcodeComputer, door string) (*IntcodeComputer, string, bool) {
	output, err := tryCommand(computer, door)
	if err != nil || isDeath(output) || isSecurityCheck(output) {
		return nil, "", false
	}

	next := computer.Clone()
	next.SetInput(StringsToASCII(door)...)
	next.Run()
	return next, output, true
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
			if isSafeToTake(current.computer, item) {
				items.Add(item)
			}
		}

		if strings.Contains(state.roomName, "Security Checkpoint") && checkpointDir == "" {
			checkpointDir = findCheckpointDirection(current.computer, state.doors)
		}

		for _, door := range state.doors {
			if next, output, ok := exploreDirection(current.computer, door); ok {
				queue.PushBack(queueItem{
					computer:   next,
					inventory:  append([]string{}, current.inventory...),
					lastOutput: output,
				})
			}
		}
	}

	return items, checkpointDir
}

func collectItemsInRoom(
	computer *IntcodeComputer,
	inventory []string,
	items []string,
	safeItems collections.Set[string],
) (*IntcodeComputer, []string) {
	current := computer
	collected := inventory

	for _, item := range items {
		if !dangerousItems.Contains(item) && safeItems.Contains(item) {
			test := current.Clone()
			test.SetInput(StringsToASCII("take " + item)...)
			if err := test.Run(); err == nil {
				current = test
				collected = append(collected, item)
			}
		}
	}

	return current, collected
}

// collectItemsAndReachCheckpoint performs BFS to collect all items and navigate to checkpoint
func collectItemsAndReachCheckpoint(program []int, safeItems collections.Set[string]) (*IntcodeComputer, error) {
	collector := New(program)
	if err := collector.Run(); err != nil {
		return nil, err
	}
	collectorOutput := collector.ReadOutputString()

	visitedStates := make(map[string]bool)
	collectionQueue := collections.NewDeque[queueItem](16)
	collectionQueue.PushBack(queueItem{
		computer:   collector,
		inventory:  []string{},
		lastOutput: collectorOutput,
	})

	for current, ok := collectionQueue.PopFront(); ok; current, ok = collectionQueue.PopFront() {
		state := parseGameOutput(current.lastOutput)

		key := stateKey(state.roomName, current.inventory)
		if visitedStates[key] {
			continue
		}
		visitedStates[key] = true

		currentComputer, currentInventory := collectItemsInRoom(current.computer, current.inventory, state.items, safeItems)

		if strings.Contains(state.roomName, "Security Checkpoint") && len(currentInventory) == safeItems.Len() {
			return currentComputer, nil
		}

		for _, door := range state.doors {
			testOutput, err := tryCommand(currentComputer, door)
			if err != nil || isDeath(testOutput) {
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

	return nil, fmt.Errorf("never reached checkpoint with all items")
}

// exploreAndCollect performs BFS exploration using computer state cloning.
// It uses a two-pass strategy:
// 1. First pass: Map the entire ship, identify safe items, and locate the security checkpoint
// 2. Second pass: Collect all items and navigate to the checkpoint with full inventory
func exploreAndCollect(program []int) (*IntcodeComputer, collections.Set[string], string, error) {
	computer := New(program)
	if err := computer.Run(); err != nil {
		return nil, nil, "", err
	}
	initialOutput := computer.ReadOutputString()

	safeItems, checkpointDir := mapShip(computer, initialOutput)
	if checkpointDir == "" || safeItems.IsEmpty() {
		return nil, safeItems, checkpointDir, nil
	}

	checkpointComputer, err := collectItemsAndReachCheckpoint(program, safeItems)
	return checkpointComputer, safeItems, checkpointDir, err
}

type grayCodeChange struct {
	itemIdx int
	take    bool
}

// grayCodeChanges generates Gray code transitions for n items.
// Gray code ensures consecutive values differ by exactly one bit,
// minimizing the number of take/drop commands from O(n * 2^n) to O(2^n).
func grayCodeChanges(n int) func(yield func(grayCodeChange) bool) {
	return func(yield func(grayCodeChange) bool) {
		currentMask := 0
		numCombinations := 1 << n

		for i := range numCombinations {
			gray := i ^ (i >> 1)
			diff := currentMask ^ gray
			currentMask = gray

			if diff != 0 {
				itemIdx := bits.TrailingZeros(uint(diff))
				take := (gray & (1 << itemIdx)) != 0
				if !yield(grayCodeChange{itemIdx: itemIdx, take: take}) {
					return
				}
			}
		}
	}
}

func tryItemCombination(computer *IntcodeComputer, securityDir string) (string, bool) {
	computer.SetInput(StringsToASCII(securityDir)...)
	if err := computer.Run(); err != nil {
		return "", false
	}

	output := computer.ReadOutputString()
	if isSecurityCheck(output) {
		return "", false
	}

	if password := extractPassword(output); password != "" {
		return password, true
	}

	return "", false
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
	test := checkpointComputer.Clone()

	for _, item := range items {
		test.SetInput(StringsToASCII("drop " + item)...)
		test.Run()
		test.GetOutput()
	}

	tried := 0

	for change := range grayCodeChanges(len(items)) {
		item := items[change.itemIdx]
		action := "drop"
		if change.take {
			action = "take"
		}

		test.SetInput(StringsToASCII(action + " " + item)...)
		if err := test.Run(); err != nil {
			continue
		}
		test.GetOutput()

		if password, found := tryItemCombination(test, securityDir); found {
			return password, nil
		}
		tried++
	}

	return "", fmt.Errorf("tried all %d combinations, none worked", tried)
}
