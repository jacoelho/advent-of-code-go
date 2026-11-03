package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestJSONReadErrorAfterDocument(t *testing.T) {
	failure := errors.New("trailing read failed")
	r := io.MultiReader(strings.NewReader("[1]"), iotest.ErrReader(failure))
	if _, err := day12p01(r); !errors.Is(err, failure) {
		t.Fatalf("got error %v, want trailing read failure", err)
	}
}

func Test_day12p01(t *testing.T) {
	aoc.AOCTest(t, day12p01, []aoc.TestInput{
		{Input: strings.NewReader("[1,2,3]"), Want: "6"},
		{Input: strings.NewReader("{\"a\":2,\"b\":4}"), Want: "6"},
		{Input: strings.NewReader("[[[3]]]"), Want: "3"},
		{Input: strings.NewReader("{\"a\":{\"b\":4},\"c\":-1}"), Want: "3"},
		{Input: strings.NewReader("{\"a\":[-1,1]}"), Want: "0"},
		{Input: strings.NewReader("[-1,{\"a\":1}]"), Want: "0"},
		{Input: strings.NewReader("[]"), Want: "0"},
		{Input: strings.NewReader("{}"), Want: "0"},
		{Input: aoc.FileInput(t, 2015, 12), Want: "191164"},
	})
}

func Test_day12p02(t *testing.T) {
	aoc.AOCTest(t, day12p02, []aoc.TestInput{
		{Input: strings.NewReader("[1,2,3]"), Want: "6"},
		{Input: strings.NewReader("[1,{\"c\":\"red\",\"b\":2},3]"), Want: "4"},
		{Input: strings.NewReader("{\"d\":\"red\",\"e\":[1,2,3,4],\"f\":5}"), Want: "0"},
		{Input: strings.NewReader("[1,\"red\",5]"), Want: "6"},
		{Input: strings.NewReader("{\"red\":5,\"nested\":{\"value\":\"red\",\"n\":7}}"), Want: "5"},
		{Input: aoc.FileInput(t, 2015, 12), Want: "87842"},
	})
}
