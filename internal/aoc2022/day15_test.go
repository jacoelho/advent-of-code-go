package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day15Example = `Sensor at x=2, y=18: closest beacon is at x=-2, y=15
Sensor at x=9, y=16: closest beacon is at x=10, y=16
Sensor at x=13, y=2: closest beacon is at x=15, y=3
Sensor at x=12, y=14: closest beacon is at x=10, y=16
Sensor at x=10, y=20: closest beacon is at x=10, y=16
Sensor at x=14, y=17: closest beacon is at x=10, y=16
Sensor at x=8, y=7: closest beacon is at x=2, y=10
Sensor at x=2, y=0: closest beacon is at x=2, y=10
Sensor at x=0, y=11: closest beacon is at x=2, y=10
Sensor at x=20, y=14: closest beacon is at x=25, y=17
Sensor at x=17, y=20: closest beacon is at x=21, y=22
Sensor at x=16, y=7: closest beacon is at x=15, y=3
Sensor at x=14, y=3: closest beacon is at x=15, y=3
Sensor at x=20, y=1: closest beacon is at x=15, y=3`

func Test_day15p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: aoc.FileInput(t, 2022, 15), Want: "5525847"},
	}
	aoc.AOCTest(t, day15p01, tests)
}

func Test_day15p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: aoc.FileInput(t, 2022, 15), Want: "13340867187704"},
	}
	aoc.AOCTest(t, day15p02, tests)
}

func TestSensorExample(t *testing.T) {
	got, err := excludedBeacons(strings.NewReader(day15Example), 10)
	if err != nil {
		t.Fatal(err)
	}
	if got != "26" {
		t.Errorf("excluded %s, want 26", got)
	}
	got, err = distressBeacon(strings.NewReader(day15Example), 20)
	if err != nil {
		t.Fatal(err)
	}
	if got != "56000011" {
		t.Errorf("frequency %s, want 56000011", got)
	}
}

func TestBeaconSquareEdges(t *testing.T) {
	tests := []struct {
		input string
		limit int64
		want  string
	}{
		{"Sensor at x=0, y=0: closest beacon is at x=0, y=0\nSensor at x=1, y=0: closest beacon is at x=1, y=0\nSensor at x=0, y=1: closest beacon is at x=0, y=1", 1, "4000001"},
		{"Sensor at x=1, y=1: closest beacon is at x=1, y=1\nSensor at x=1, y=0: closest beacon is at x=1, y=0\nSensor at x=0, y=1: closest beacon is at x=0, y=1", 1, "0"},
		{"Sensor at x=1, y=1: closest beacon is at x=1, y=1", 0, "0"},
	}
	for _, tt := range tests {
		got, err := distressBeacon(strings.NewReader(tt.input), tt.limit)
		if err != nil {
			t.Fatal(err)
		}
		if got != tt.want {
			t.Errorf("got %s, want %s", got, tt.want)
		}
	}
	for _, limit := range []int64{0, 1} {
		if _, err := distressBeacon(strings.NewReader("Sensor at x=0, y=0: closest beacon is at x=0, y=0"), limit); err == nil {
			t.Errorf("expected uniqueness error at limit %d", limit)
		}
	}
}

func TestDay15ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day15p01},
		{"part2", day15p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
