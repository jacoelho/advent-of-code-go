package aoc2019

import (
	"fmt"
	"io"
	"math/bits"
	"regexp"
	"slices"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

var (
	passwordRe = regexp.MustCompile(`\d{5,}`)
	roomRe     = regexp.MustCompile(`== (.+) ==`)
)

type adventureRoom struct {
	roomName string
	doors    []string
	items    []string
}

var dangerousItems = collections.NewSet(
	"infinite loop",
	"escape pod",
	"molten lava",
	"photons",
	"giant electromagnet",
)

type adventureState struct {
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

func parseGameOutput(text string) adventureRoom {
	state := adventureRoom{}

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

	return state
}

func isDeath(text string) bool {
	deathPhrases := []string{"ejected", "are frozen", "crushed", "knocked"}
	return xiter.FromSlice(deathPhrases).Any(func(phrase string) bool {
		return strings.Contains(text, phrase)
	})
}

func isSecurityCheck(text string) bool {
	return strings.Contains(text, "Alert") ||
		strings.Contains(text, "lighter") ||
		strings.Contains(text, "heavier")
}

func extractPassword(text string) string {
	return passwordRe.FindString(text)
}

func tryCommand(computer *IntcodeComputer, command string) (*IntcodeComputer, string, error) {
	next := computer.Clone()
	next.SetInput(StringsToASCII(command)...)
	if err := next.Run(); err != nil {
		return nil, "", fmt.Errorf("command %q: %w", command, err)
	}
	return next, next.ReadOutputString(), nil
}

func stateKey(roomName string, inventory []string) string {
	sortedInv := slices.Clone(inventory)
	slices.Sort(sortedInv)
	return fmt.Sprintf("%s:%v", roomName, sortedInv)
}

func isSafeToTake(computer *IntcodeComputer, item string) (bool, error) {
	if dangerousItems.Contains(item) {
		return false, nil
	}
	_, output, err := tryCommand(computer, "take "+item)
	if err != nil {
		return false, err
	}
	return !isDeath(output), nil
}

func findCheckpointDirection(computer *IntcodeComputer, doors []string) (string, error) {
	for _, door := range doors {
		_, output, err := tryCommand(computer, door)
		if err != nil {
			return "", err
		}
		if isSecurityCheck(output) || extractPassword(output) != "" {
			return door, nil
		}
	}
	return "", nil
}

func exploreDirection(computer *IntcodeComputer, door string) (*IntcodeComputer, string, error) {
	next, output, err := tryCommand(computer, door)
	if err != nil {
		return nil, "", err
	}
	if isDeath(output) || isSecurityCheck(output) || extractPassword(output) != "" {
		return nil, "", nil
	}
	return next, output, nil
}

func mapShip(computer *IntcodeComputer, initialOutput string) (collections.Set[string], string, error) {
	var checkpointDir string
	visited := make(map[string]bool)
	queue := collections.NewDeque[adventureState](16)
	queue.PushBack(adventureState{
		computer:   computer,
		lastOutput: initialOutput,
	})

	items := collections.NewSet[string]()

	for current, ok := queue.PopFront(); ok; current, ok = queue.PopFront() {
		state := parseGameOutput(current.lastOutput)

		if visited[state.roomName] {
			continue
		}
		if state.roomName == "" || len(visited) == 256 {
			return nil, "", fmt.Errorf("ship exploration needs named rooms and at most 256 rooms")
		}
		visited[state.roomName] = true

		for _, item := range state.items {
			safe, err := isSafeToTake(current.computer, item)
			if err != nil {
				return nil, "", err
			}
			if safe {
				items.Add(item)
				if items.Len() > 16 {
					return nil, "", fmt.Errorf("ship has more than 16 safe items")
				}
			}
		}

		if strings.Contains(state.roomName, "Security Checkpoint") && checkpointDir == "" {
			var err error
			checkpointDir, err = findCheckpointDirection(current.computer, state.doors)
			if err != nil {
				return nil, "", err
			}
		}

		for _, door := range state.doors {
			next, output, err := exploreDirection(current.computer, door)
			if err != nil {
				return nil, "", err
			}
			if next != nil {
				if queue.Size() == 100_000 {
					return nil, "", fmt.Errorf("ship exploration exceeds 100000 queued states")
				}
				queue.PushBack(adventureState{
					computer:   next,
					lastOutput: output,
				})
			}
		}
	}

	return items, checkpointDir, nil
}

func collectItemsInRoom(
	computer *IntcodeComputer,
	inventory []string,
	items []string,
	safeItems collections.Set[string],
) (*IntcodeComputer, []string, error) {
	current := computer
	collected := slices.Clone(inventory)
	for _, item := range items {
		if !safeItems.Contains(item) || slices.Contains(collected, item) {
			continue
		}
		next, _, err := tryCommand(current, "take "+item)
		if err != nil {
			return nil, nil, err
		}
		current = next
		collected = append(collected, item)
	}
	return current, collected, nil
}

func collectItemsAndReachCheckpoint(program []int, safeItems collections.Set[string]) (*IntcodeComputer, error) {
	collector := New(program)
	if err := collector.Run(); err != nil {
		return nil, err
	}
	collectorOutput := collector.ReadOutputString()

	visitedStates := make(map[string]bool)
	collectionQueue := collections.NewDeque[adventureState](16)
	collectionQueue.PushBack(adventureState{
		computer:   collector,
		lastOutput: collectorOutput,
	})

	for current, ok := collectionQueue.PopFront(); ok; current, ok = collectionQueue.PopFront() {
		state := parseGameOutput(current.lastOutput)

		key := stateKey(state.roomName, current.inventory)
		if visitedStates[key] {
			continue
		}
		if state.roomName == "" || len(visitedStates) == 100_000 {
			return nil, fmt.Errorf("item collection needs named rooms and at most 100000 states")
		}
		visitedStates[key] = true

		currentComputer, currentInventory, err := collectItemsInRoom(current.computer, current.inventory, state.items, safeItems)
		if err != nil {
			return nil, err
		}

		if strings.Contains(state.roomName, "Security Checkpoint") && len(currentInventory) == safeItems.Len() {
			return currentComputer, nil
		}

		for _, door := range state.doors {
			next, output, err := exploreDirection(currentComputer, door)
			if err != nil {
				return nil, err
			}
			if next == nil {
				continue
			}
			if collectionQueue.Size() == 100_000 {
				return nil, fmt.Errorf("item collection exceeds 100000 queued states")
			}
			collectionQueue.PushBack(adventureState{
				computer:   next,
				inventory:  slices.Clone(currentInventory),
				lastOutput: output,
			})
		}
	}

	return nil, fmt.Errorf("never reached checkpoint with all items")
}

func exploreAndCollect(program []int) (*IntcodeComputer, collections.Set[string], string, error) {
	computer := New(program)
	if err := computer.Run(); err != nil {
		return nil, nil, "", err
	}
	initialOutput := computer.ReadOutputString()

	safeItems, checkpointDir, err := mapShip(computer, initialOutput)
	if err != nil {
		return nil, nil, "", err
	}
	if checkpointDir == "" {
		return nil, nil, "", fmt.Errorf("ship has no security checkpoint exit")
	}

	checkpointComputer, err := collectItemsAndReachCheckpoint(program, safeItems)
	return checkpointComputer, safeItems, checkpointDir, err
}

type grayCodeChange struct {
	itemIdx int
	take    bool
}

// grayCodeChanges yields the single-item changes that enumerate subsets of n items in
// Gray-code order, starting from an empty inventory.
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

func tryItemCombination(computer *IntcodeComputer, securityDir string) (string, error) {
	computer.SetInput(StringsToASCII(securityDir)...)
	if err := computer.Run(); err != nil {
		return "", fmt.Errorf("try security checkpoint: %w", err)
	}
	output := computer.ReadOutputString()
	if isSecurityCheck(output) {
		return "", nil
	}
	if password := extractPassword(output); password != "" {
		return password, nil
	}
	return "", fmt.Errorf("checkpoint response has neither rejection nor password")
}

func findWinningInventory(computer *IntcodeComputer, items []string, securityDir string) (string, error) {
	if len(items) > 16 {
		return "", fmt.Errorf("inventory search exceeds 16 items")
	}
	current := computer.Clone()
	for _, item := range items {
		next, _, err := tryCommand(current, "drop "+item)
		if err != nil {
			return "", err
		}
		current = next
	}
	password, err := tryItemCombination(current, securityDir)
	if err != nil || password != "" {
		return password, err
	}
	tried := 1
	for change := range grayCodeChanges(len(items)) {
		action := "drop"
		if change.take {
			action = "take"
		}
		next, _, err := tryCommand(current, action+" "+items[change.itemIdx])
		if err != nil {
			return "", err
		}
		current = next
		password, err := tryItemCombination(current, securityDir)
		if err != nil || password != "" {
			return password, err
		}
		tried++
	}
	return "", fmt.Errorf("tried all %d inventory combinations, none worked", tried)
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

	items := slices.Sorted(itemsSet.Iter())
	return findWinningInventory(checkpointComputer, items, securityDir)
}
