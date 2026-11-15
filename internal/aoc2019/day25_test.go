package aoc2019

import (
	"strconv"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(adventureProgramText(adventureReplyProgram(
				"== Security Checkpoint ==\nDoors here lead:\n- north\n\nCommand?\n",
				"typing 12345 on the keypad\n",
			))),
			Want: "12345",
		},
		{
			Input: aoc.FileInput(t, 2019, 25),
			Want:  "2214608912",
		},
	}
	aoc.AOCTest(t, day25p01, tests)
}

func appendAdventureText(program []int, text string) []int {
	for _, c := range text {
		program = append(program, 104, int(c))
	}
	return program
}

func adventureReplyProgram(initial string, replies ...string) []int {
	program := appendAdventureText(nil, initial)
	for _, reply := range replies {
		start := len(program)
		program = append(program, 3, 1000, 1008, 1000, 10, 1001, 1006, 1001, start)
		program = appendAdventureText(program, reply)
	}
	return append(program, 99)
}

func adventureProgramText(program []int) string {
	words := make([]string, len(program))
	for i, word := range program {
		words[i] = strconv.Itoa(word)
	}
	return strings.Join(words, ",")
}

func TestAdventureCommandContinuation(t *testing.T) {
	computer := New(adventureReplyProgram("", "A\n", "B\n"))
	output, err := sendCommand(computer, "north")
	if err != nil {
		t.Fatal(err)
	}
	if output != "A\n" || !computer.IsWaiting() || len(computer.GetOutput()) != 0 {
		t.Fatalf("first command: output = %q, waiting = %v", output, computer.IsWaiting())
	}
	output, err = sendCommand(computer, "south")
	if err != nil {
		t.Fatal(err)
	}
	if output != "B\n" || !computer.IsHalted() {
		t.Fatalf("second command: output = %q, halted = %v", output, computer.IsHalted())
	}
}

func inventoryCheckpointProgram(winningInventory int) []int {
	// Cell 1002 tracks the held item; drop/take commands clear/set it.
	program := []int{
		1101, 1, 0, 1002,
		3, 1000,
		3, 1001, 1008, 1001, 10, 1003, 1006, 1003, 6,
		1008, 1000, 100, 1003, 1005, 1003, 51,
		1008, 1000, 116, 1003, 1005, 1003, 58,
		1008, 1002, winningInventory, 1003, 1005, 1003, 65,
	}
	program = appendAdventureText(program, "Alert\n")
	program = append(program,
		1105, 1, 4,
		1101, 0, 0, 1002, 1105, 1, 4,
		1101, 1, 0, 1002, 1105, 1, 4,
	)
	return append(appendAdventureText(program, "typing 12345 on the keypad\n"), 99)
}

func TestAdventureInventorySearch(t *testing.T) {
	for _, winningInventory := range []int{0, 1} {
		computer := New(inventoryCheckpointProgram(winningInventory))
		got, err := findWinningInventory(computer, []string{"coin"}, "north")
		if err != nil || got != "12345" {
			t.Fatalf("winning inventory %d: got %q, %v; want 12345", winningInventory, got, err)
		}
	}
	computer := New(adventureReplyProgram("", "typing 12345 on the keypad\n"))
	if got, err := findWinningInventory(computer, nil, "north"); err != nil || got != "12345" {
		t.Fatalf("empty inventory: got %q, %v; want 12345", got, err)
	}
}

func TestAdventureCommandErrors(t *testing.T) {
	program := adventureReplyProgram("== Security Checkpoint ==\nDoors here lead:\n- north\n\nCommand?\n", "")
	program[len(program)-1] = 42
	if _, err := day25p01(strings.NewReader(adventureProgramText(program))); err == nil || !strings.Contains(err.Error(), "opcode: 42") {
		t.Fatalf("expected nested command error, got %v", err)
	}
	if _, err := sendCommand(New([]int{42}), "north"); err == nil {
		t.Fatal("expected command error")
	}
	if _, err := findWinningInventory(New([]int{42}), nil, "north"); err == nil {
		t.Fatal("expected checkpoint error")
	}
}

func TestAdventureExploration(t *testing.T) {
	const room = "== Hall ==\nDoors here lead:\n- north\n\nCommand?\n"
	exploreAfterReply := func(reply string) (shipSurvey, error) {
		var survey shipSurvey
		computer := New(adventureReplyProgram("", reply, room))
		err := explore(computer, room, nil, map[string]bool{}, &survey)
		return survey, err
	}

	survey, err := exploreAfterReply("Alert! Droids on this ship are heavier than the detected value!\n" + room)
	if err != nil || survey.checkpointDoor != "north" {
		t.Fatalf("pressure floor: door = %q, error = %v; want north", survey.checkpointDoor, err)
	}
	if _, err := exploreAfterReply("You are frozen.\n"); err == nil {
		t.Fatal("expected an error when moving ends the game")
	}
	if _, err := exploreAfterReply("nothing here\n"); err == nil {
		t.Fatal("expected an error for output without a room")
	}
	var unknownDoor shipSurvey
	err = explore(New([]int{99}), "== Hall ==\nDoors here lead:\n- up\n\nCommand?\n", nil, map[string]bool{}, &unknownDoor)
	if err == nil {
		t.Fatal("expected an error for an unknown door")
	}
}
