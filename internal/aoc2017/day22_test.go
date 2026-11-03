package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const virusExample = "..#\n#..\n..."

func Test_day22p01(t *testing.T) {
	aoc.AOCTest(t, day22p01, []aoc.TestInput{
		{Input: strings.NewReader(virusExample), Want: "5587"},
		{Input: aoc.FileInput(t, 2017, 22), Want: "5570"},
	})
}

func Test_day22p02(t *testing.T) {
	aoc.AOCTest(t, day22p02, []aoc.TestInput{
		{Input: strings.NewReader(virusExample), Want: "2511944"},
		{Input: aoc.FileInput(t, 2017, 22), Want: "2512022"},
	})
}

func TestVirusBursts(t *testing.T) {
	tests := []struct {
		name   string
		mode   virusMode
		bursts int
		want   int
	}{
		{name: "simple", mode: simpleVirus, bursts: 7, want: 5},
		{name: "evolved", mode: evolvedVirus, bursts: 100, want: 26},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := virusBursts(strings.NewReader(virusExample), tt.mode, tt.bursts)
			if err != nil || got != tt.want {
				t.Fatalf("virusBursts() = %d, %v; want %d", got, err, tt.want)
			}
		})
	}
}
