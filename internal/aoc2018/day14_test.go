package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day14p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("9"),
			Want:  "5158916779",
		},
		{
			Input: strings.NewReader("5"),
			Want:  "0124515891",
		},
		{
			Input: strings.NewReader("18"),
			Want:  "9251071085",
		},
		{
			Input: strings.NewReader("2018"),
			Want:  "5941429882",
		},
		{
			Input: aoc.FileInput(t, 2018, 14),
			Want:  "1631191756",
		},
	}

	aoc.AOCTest(t, day14p01, tests)
}

func Test_day14p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("51589"),
			Want:  "9",
		},
		{
			Input: strings.NewReader("01245"),
			Want:  "5",
		},
		{
			Input: strings.NewReader("92510"),
			Want:  "18",
		},
		{
			Input: strings.NewReader("59414"),
			Want:  "2018",
		},
		{
			Input: strings.NewReader("37"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("1"),
			Want:  "2",
		},
		{
			Input: strings.NewReader("0"),
			Want:  "3",
		},
	}

	aoc.AOCTest(t, day14p02, tests)
}

func TestRecipeFullInput(t *testing.T) {
	if testing.Short() {
		t.Skip("full 20-million-score sequence search")
	}

	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2018, 14),
			Want:  "20219475",
		},
	}

	aoc.AOCTest(t, day14p02, tests)
}

func TestRecipeErrors(t *testing.T) {
	for _, input := range []string{"", "-1", "12x"} {
		if _, err := day14p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}
