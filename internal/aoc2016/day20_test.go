package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	aoc.AOCTest(t, day20p01, []aoc.TestInput{{Input: strings.NewReader("5-8\n0-2\n4-7"), Want: "3"}, {Input: aoc.FileInput(t, 2016, 20), Want: "23923783"}})
}

func Test_day20p02(t *testing.T) {
	aoc.AOCTest(t, day20p02, []aoc.TestInput{{Input: strings.NewReader(""), Want: "4294967296"}, {Input: strings.NewReader("0-4294967294"), Want: "1"}, {Input: strings.NewReader("0-4294967295"), Want: "0"}, {Input: aoc.FileInput(t, 2016, 20), Want: "125"}})
}

func TestAddressGaps(t *testing.T) {
	first, count, err := allowedAddresses(strings.NewReader("5-8\n0-2\n4-7"), 9)
	if err != nil || first != 3 || count != 2 {
		t.Fatalf("gaps=%d,%d,%v; want3,2", first, count, err)
	}
	if _, err := day20p01(strings.NewReader("0-4294967295")); err == nil {
		t.Fatal("all addresses blocked")
	}
	for _, input := range []string{"5-2", "-1-2", "0-4294967296", "x-2"} {
		if _, err := day20p02(strings.NewReader(input)); err == nil {
			t.Errorf("accepted%q", input)
		}
	}
}

func TestDay20ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day20p01,
		day20p02,
	} {
		if _, err := solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Errorf("read error%v lost input cause", err)
		}
	}
}
