package aoc2019

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day22p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("cut 20000"), Want: "2033"},
		{Input: strings.NewReader("cut -20000"), Want: "2005"},
		{Input: strings.NewReader("cut -9223372036854775808"), Want: "3918"},
		{Input: strings.NewReader("deal with increment 3"), Want: "6057"},
		{
			Input: aoc.FileInput(t, 2019, 22),
			Want:  "1234",
		},
	}
	aoc.AOCTest(t, day22p01, tests)
}

func Test_day22p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("cut 1"), Want: "101741582078681"},
		{Input: strings.NewReader("cut 0"), Want: "2020"},
		{Input: strings.NewReader("deal into new stack"), Want: "119315717512026"},
		{
			Input: aoc.FileInput(t, 2019, 22),
			Want:  "7757787935983",
		},
	}
	aoc.AOCTest(t, day22p02, tests)
}

func TestRepeatedShuffle(t *testing.T) {
	instructions, err := parseShuffleInstructions(strings.NewReader("deal with increment 3\ncut -2\ndeal into new stack"))
	if err != nil {
		t.Fatal(err)
	}
	f, err := shuffleMap(instructions, 10)
	if err != nil {
		t.Fatal(err)
	}
	cardAt, err := f.pow(3).inverse()
	if err != nil {
		t.Fatal(err)
	}
	want := []uint64{7, 4, 1, 8, 5, 2, 9, 6, 3, 0}
	for position, card := range want {
		if got := cardAt.apply(uint64(position)); got != card {
			t.Errorf("position %d: got %d, want %d", position, got, card)
		}
	}
}

func TestShuffleErrors(t *testing.T) {
	for _, input := range []string{
		"deal with increment 0",
		"deal with increment -3",
		"cut",
		"cut x",
		"xcut 5",
		"deal into new stacks",
	} {
		for _, solve := range []func(io.Reader) (string, error){day22p01, day22p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("expected shuffle error for %q", input)
			}
		}
	}
	if _, err := day22p01(strings.NewReader("deal with increment 10007")); err == nil {
		t.Error("expected error for an increment that does not permute the deck")
	}
}
