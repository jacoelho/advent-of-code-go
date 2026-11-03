package aoc2016

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const botExample = `value 5 goes to bot 2
bot 2 gives low to bot 1 and high to bot 0
value 3 goes to bot 1
bot 1 gives low to output 1 and high to bot 0
bot 0 gives low to output 2 and high to output 0
value 2 goes to bot 2`

func Test_day10p01(t *testing.T) {
	aoc.AOCTest(t, day10p01, []aoc.TestInput{
		{Input: strings.NewReader("value 17 goes to bot 7\nvalue 61 goes to bot 7\nbot 7 gives low to output 0 and high to output 1"), Want: "7"},
		{Input: aoc.FileInput(t, 2016, 10), Want: "116"},
	})
}

func Test_day10p02(t *testing.T) {
	aoc.AOCTest(t, day10p02, []aoc.TestInput{
		{Input: strings.NewReader(botExample), Want: "30"},
		{Input: aoc.FileInput(t, 2016, 10), Want: "23903"},
	})
}

func TestBotExampleIgnoresRuleOrder(t *testing.T) {
	lines := strings.Split(botExample, "\n")
	slices.Reverse(lines)
	for _, input := range []string{botExample, strings.Join(lines, "\n")} {
		bot, outputs, err := botFactory(strings.NewReader(input), [2]int{2, 5})
		if err != nil {
			t.Fatal(err)
		}
		if bot != 2 || !slices.Equal(outputs[0], []int{5}) || !slices.Equal(outputs[1], []int{2}) || !slices.Equal(outputs[2], []int{3}) {
			t.Errorf("bot=%d outputs=%v; want bot 2 and outputs 5,2,3", bot, outputs)
		}
	}
}

func TestBotRoutingRulePersists(t *testing.T) {
	input := `bot 0 gives low to output 0 and high to output 1
value 2 goes to bot 0
value 5 goes to bot 0
value 17 goes to bot 0
value 61 goes to bot 0`
	bot, outputs, err := botFactory(strings.NewReader(input), [2]int{17, 61})
	if err != nil {
		t.Fatal(err)
	}
	if bot != 0 || !slices.Equal(outputs[0], []int{2, 17}) || !slices.Equal(outputs[1], []int{5, 61}) {
		t.Errorf("bot=%d outputs=%v; want bot 0 and two deliveries per output", bot, outputs)
	}
}

func TestBotErrors(t *testing.T) {
	for _, input := range []string{"", "value x goes to bot 0", "value 1 goes to bot 0\nvalue 2 goes to bot 0"} {
		if _, err := day10p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day10p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day10p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
