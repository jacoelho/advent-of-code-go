package aoc2022

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day19p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(blueprintExample), Want: "33"},
		{Input: strings.NewReader(strings.NewReplacer("Blueprint 1:", "Blueprint 7:", "Blueprint 2:", "Blueprint 11:").Replace(blueprintExample)), Want: "195"},
		{Input: aoc.FileInput(t, 2022, 19), Want: "1487"},
	}

	aoc.AOCTest(t, day19p01, tests)
}

func Test_day19p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(strings.NewReplacer("Blueprint 1:", "Blueprint 7:", "Blueprint 2:", "Blueprint 11:").Replace(blueprintExample)), Want: "3472"},
		{Input: aoc.FileInput(t, 2022, 19), Want: "13440"},
	}

	aoc.AOCTest(t, day19p02, tests)
}

func TestRobotBuildTiming(t *testing.T) {
	blueprint := robotBlueprint{costs: [4][3]int{{1, 0, 0}, {1, 0, 0}, {1, 1, 0}, {1, 0, 1}}}
	for _, test := range []struct{ minutes, want int }{{6, 0}, {7, 1}} {
		work := 0
		got, err := blueprintGeodes(blueprint, test.minutes, &work)
		if err != nil {
			t.Fatal(err)
		}
		if got != test.want {
			t.Errorf("%d minutes = %d geodes, want %d", test.minutes, got, test.want)
		}
	}
}

func TestBlueprintErrors(t *testing.T) {
	for _, input := range []string{"", "Blueprint 1: invalid", blueprintExample + "\n" + blueprintExample} {
		if _, err := day19p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected invalid blueprint error for %q", input)
		}
	}
}

const blueprintExample = `Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 2 ore. Each obsidian robot costs 3 ore and 14 clay. Each geode robot costs 2 ore and 7 obsidian.
Blueprint 2: Each ore robot costs 2 ore. Each clay robot costs 3 ore. Each obsidian robot costs 3 ore and 8 clay. Each geode robot costs 3 ore and 12 obsidian.`
