package aoc2021

import (
	"errors"
	"io"
	"slices"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const snailfishExample = `[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]
[[[5,[2,8]],4],[5,[[9,9],0]]]
[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]
[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]
[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]
[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]
[[[[5,4],[7,7]],8],[[8,3],8]]
[[9,3],[[9,9],[6,[4,9]]]]
[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]
[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]`

func Test_day18p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(snailfishExample), Want: "4140"},
		{Input: strings.NewReader("[1,1]\n[1,1]"), Want: "25"},
		{Input: strings.NewReader("[[[[4,3],4],4],[7,[[8,4],9]]]\n[1,1]"), Want: "1384"},
		{Input: aoc.FileInput(t, 2021, 18), Want: "4137"},
	}

	aoc.AOCTest(t, day18p01, tests)
}

func Test_day18p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(snailfishExample), Want: "3993"},
		{Input: strings.NewReader("[1,1]\n[1,1]"), Want: "25"},
		{Input: aoc.FileInput(t, 2021, 18), Want: "4573"},
	}

	aoc.AOCTest(t, day18p02, tests)
}

func TestSnailfishDeepExplosion(t *testing.T) {
	number, err := parseSnailfish([]byte("[[[[[[9,8],1],2],3],4],5]"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := parseSnailfish([]byte("[[[[[0,9],2],3],4],5]"))
	if err != nil {
		t.Fatal(err)
	}
	got, exploded, err := explodeSnailfish(number)
	if err != nil {
		t.Fatal(err)
	}
	if !exploded || !slices.Equal(got, want) {
		t.Fatalf("explosion = %v, %v; want %v, true", got, exploded, want)
	}
}

func TestSnailfishAdditionPreservesOperands(t *testing.T) {
	left, err := parseSnailfish([]byte("[[[[4,3],4],4],[7,[[8,4],9]]]"))
	if err != nil {
		t.Fatal(err)
	}
	right, err := parseSnailfish([]byte("[1,1]"))
	if err != nil {
		t.Fatal(err)
	}
	wantLeft, wantRight := slices.Clone(left), slices.Clone(right)
	work := 0
	if _, err := addSnailfish(left, right, &work); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(left, wantLeft) || !slices.Equal(right, wantRight) {
		t.Fatal("addition changed an input number")
	}
}

func TestSnailfishErrors(t *testing.T) {
	for _, input := range []string{"", "[1]", "[1,2,3]", "[1,-2]", "[1,2]x", "[1,2]"} {
		if _, err := day18p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

func TestDay18ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day18p01},
		{"part2", day18p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
