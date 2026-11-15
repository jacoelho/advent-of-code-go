package aoc2019

import (
	"fmt"
	"io"
	"math/bits"
	"regexp"
	"slices"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
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

// dangerousItems end or trap the game when taken; "infinite loop" never
// returns control, so it cannot be probed.
var dangerousItems = collections.NewSet(
	"infinite loop",
	"escape pod",
	"molten lava",
	"photons",
	"giant electromagnet",
)

var oppositeDoor = map[string]string{
	"north": "south",
	"south": "north",
	"east":  "west",
	"west":  "east",
}

func parseList(lines []string, header string) []string {
	var result []string
	for i, line := range lines {
		if !strings.HasPrefix(line, header) {
			continue
		}
		for _, listLine := range lines[i+1:] {
			after, ok := strings.CutPrefix(listLine, "- ")
			if !ok {
				break
			}
			result = append(result, after)
		}
		break
	}
	return result
}

// parseGameOutput describes the last room the output shows.
func parseGameOutput(text string) adventureRoom {
	var room adventureRoom
	matches := roomRe.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		return room
	}
	room.roomName = matches[len(matches)-1][1]
	lines := strings.Split(text[strings.LastIndex(text, "== "):], "\n")
	room.doors = parseList(lines, "Doors here lead:")
	room.items = parseList(lines, "Items here:")
	return room
}

func isDeath(text string) bool {
	for _, phrase := range []string{"ejected", "are frozen", "crushed", "knocked"} {
		if strings.Contains(text, phrase) {
			return true
		}
	}
	return false
}

func isSecurityCheck(text string) bool {
	return strings.Contains(text, "Alert") ||
		strings.Contains(text, "lighter") ||
		strings.Contains(text, "heavier")
}

func extractPassword(text string) string {
	return passwordRe.FindString(text)
}

func sendCommand(computer *IntcodeComputer, command string) (string, error) {
	computer.AddInput(StringsToASCII(command)...)
	if err := computer.Run(); err != nil {
		return "", fmt.Errorf("command %q: %w", command, err)
	}
	return computer.ReadOutputString(), nil
}

// shipSurvey is what one walk of the ship learns: the safe items, all carried
// by the droid, the route from the start to the security checkpoint, and the
// checkpoint door that weighs the droid.
type shipSurvey struct {
	items           []string
	checkpointFound bool
	checkpointPath  []string
	checkpointDoor  string
	password        string
}

// explore walks every room reachable from the current one depth-first,
// taking safe items, and returns to the current room. It stops early once a
// password appears.
func explore(computer *IntcodeComputer, output string, path []string, visited map[string]bool, survey *shipSurvey) error {
	room := parseGameOutput(output)
	if room.roomName == "" {
		return fmt.Errorf("game output names no room: %q", output)
	}
	visited[room.roomName] = true

	for _, item := range room.items {
		if dangerousItems.Contains(item) {
			continue
		}
		probe, err := sendCommand(computer.Clone(), "take "+item)
		if err != nil {
			return err
		}
		if isDeath(probe) {
			continue
		}
		if _, err := sendCommand(computer, "take "+item); err != nil {
			return err
		}
		survey.items = append(survey.items, item)
	}

	if strings.Contains(room.roomName, "Security Checkpoint") {
		survey.checkpointFound = true
		survey.checkpointPath = slices.Clone(path)
	}

	for _, door := range room.doors {
		back, ok := oppositeDoor[door]
		if !ok {
			return fmt.Errorf("unknown door %q in %s", door, room.roomName)
		}
		if len(path) > 0 && back == path[len(path)-1] {
			continue
		}
		next, err := sendCommand(computer, door)
		if err != nil {
			return err
		}
		if survey.password = extractPassword(next); survey.password != "" {
			return nil
		}
		if isSecurityCheck(next) {
			// The pressure floor ejects the droid back into this room.
			survey.checkpointDoor = door
			continue
		}
		if isDeath(next) {
			return fmt.Errorf("moving %s from %s ended the game", door, room.roomName)
		}
		if !visited[parseGameOutput(next).roomName] {
			if err := explore(computer, next, append(path, door), visited, survey); err != nil || survey.password != "" {
				return err
			}
		}
		if _, err := sendCommand(computer, back); err != nil {
			return err
		}
	}
	return nil
}

// grayCodeChange is one step of a Gray-code walk over item subsets: take or
// drop a single item.
type grayCodeChange struct {
	itemIdx int
	take    bool
}

// grayCodeChanges yields the single-item changes that enumerate subsets of n
// items in Gray-code order, starting from an empty inventory.
func grayCodeChanges(n int) func(yield func(grayCodeChange) bool) {
	return func(yield func(grayCodeChange) bool) {
		for i := 1; i < 1<<n; i++ {
			gray := i ^ (i >> 1)
			itemIdx := bits.TrailingZeros(uint(gray ^ (i - 1) ^ ((i - 1) >> 1)))
			if !yield(grayCodeChange{itemIdx: itemIdx, take: gray&(1<<itemIdx) != 0}) {
				return
			}
		}
	}
}

func tryItemCombination(computer *IntcodeComputer, securityDir string) (string, error) {
	output, err := sendCommand(computer, securityDir)
	if err != nil {
		return "", fmt.Errorf("try security checkpoint: %w", err)
	}
	if password := extractPassword(output); password != "" {
		return password, nil
	}
	if isSecurityCheck(output) {
		return "", nil
	}
	return "", fmt.Errorf("checkpoint response has neither rejection nor password")
}

// findWinningInventory stands at the security checkpoint carrying items and
// tries every subset of them on the pressure floor.
func findWinningInventory(computer *IntcodeComputer, items []string, securityDir string) (string, error) {
	for _, item := range items {
		if _, err := sendCommand(computer, "drop "+item); err != nil {
			return "", err
		}
	}
	password, err := tryItemCombination(computer, securityDir)
	if err != nil || password != "" {
		return password, err
	}
	for change := range grayCodeChanges(len(items)) {
		action := "drop "
		if change.take {
			action = "take "
		}
		if _, err := sendCommand(computer, action+items[change.itemIdx]); err != nil {
			return "", err
		}
		password, err := tryItemCombination(computer, securityDir)
		if err != nil || password != "" {
			return password, err
		}
	}
	return "", fmt.Errorf("none of the %d inventory combinations passes the checkpoint", 1<<len(items))
}

func day25p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	computer := New(program)
	if err := computer.Run(); err != nil {
		return "", err
	}
	var survey shipSurvey
	if err := explore(computer, computer.ReadOutputString(), nil, map[string]bool{}, &survey); err != nil {
		return "", err
	}
	if survey.password != "" {
		return survey.password, nil
	}
	if !survey.checkpointFound || survey.checkpointDoor == "" {
		return "", fmt.Errorf("ship has no security checkpoint exit")
	}

	for _, door := range survey.checkpointPath {
		if _, err := sendCommand(computer, door); err != nil {
			return "", err
		}
	}
	slices.Sort(survey.items)
	return findWinningInventory(computer, survey.items, survey.checkpointDoor)
}
