package aoc2019

import (
	"math/big"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day22p01(t *testing.T) {
	tests := []aoc.TestInput{
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
	inverse, err := composeInverse(instructions, big.NewInt(10))
	if err != nil {
		t.Fatal(err)
	}
	repeated := inverse.pow(3)
	want := []int64{7, 4, 1, 8, 5, 2, 9, 6, 3, 0}
	for position, card := range want {
		if got := repeated.apply(big.NewInt(int64(position))); got.Cmp(big.NewInt(card)) != 0 {
			t.Errorf("position %d: got %s, want %d", position, got, card)
		}
	}
}

func TestNoninvertibleShuffle(t *testing.T) {
	if _, err := day22p02(strings.NewReader("deal with increment 0")); err == nil {
		t.Fatal("expected noninvertible increment error")
	}
}
