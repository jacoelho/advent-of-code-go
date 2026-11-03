package aoc2016

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const storageExample = `Filesystem Size Used Avail Use%
/dev/grid/node-x0-y0 10T 8T 2T 80%
/dev/grid/node-x0-y1 11T 6T 5T 54%
/dev/grid/node-x0-y2 32T 28T 4T 87%
/dev/grid/node-x1-y0 9T 7T 2T 77%
/dev/grid/node-x1-y1 8T 0T 8T 0%
/dev/grid/node-x1-y2 11T 7T 4T 63%
/dev/grid/node-x2-y0 10T 6T 4T 60%
/dev/grid/node-x2-y1 9T 8T 1T 88%
/dev/grid/node-x2-y2 9T 6T 3T 66%`

const mergeStorage = `/dev/grid/node-x0-y0 10T 0T 10T 0%
/dev/grid/node-x1-y0 10T 2T 8T 20%
/dev/grid/node-x2-y0 10T 3T 7T 30%`

const occupiedStorage = `/dev/grid/node-x0-y0 10T 3T 7T 30%
/dev/grid/node-x1-y0 10T 2T 8T 20%`

func Test_day22p01(t *testing.T) {
	aoc.AOCTest(t, day22p01, []aoc.TestInput{{Input: strings.NewReader(mergeStorage), Want: "4"}, {Input: strings.NewReader(occupiedStorage), Want: "2"}, {Input: aoc.FileInput(t, 2016, 22), Want: "1043"}})
}

func Test_day22p02(t *testing.T) {
	aoc.AOCTest(t, day22p02, []aoc.TestInput{{Input: strings.NewReader(storageExample), Want: "7"}, {Input: strings.NewReader(mergeStorage), Want: "2"}, {Input: strings.NewReader(occupiedStorage), Want: "1"}, {Input: aoc.FileInput(t, 2016, 22), Want: "185"}})
}

func TestStorageCannotSplitGoalData(t *testing.T) {
	input := "/dev/grid/node-x0-y0 3T 2T 1T 66%\n/dev/grid/node-x1-y0 10T 5T 5T 50%\n/dev/grid/node-x1-y1 10T 0T 10T 0%"
	if _, err := day22p02(strings.NewReader(input)); err == nil {
		t.Fatal("goal payload cannot fit access node")
	}
}
