package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day02Example = `forward 5
down 5
forward 8
up 3
down 8
forward 2`

func Test_day02p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day02Example), Want: "150"},
		{Input: strings.NewReader(`down 9223372036854775807
forward 2
up 9223372036854775807`), Want: "0"},
		{Input: aoc.FileInput(t, 2021, 2), Want: "1580000"},
	}

	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day02Example), Want: "900"},
		{Input: aoc.FileInput(t, 2021, 2), Want: "1251263225"},
	}

	aoc.AOCTest(t, day02p02, tests)
}

func TestDay02OverflowProductPartOne(t *testing.T) {
	if _, err := day02p01(strings.NewReader(`down 9223372036854775807
forward 2`)); err == nil {
		t.Fatal("expected overflow product part one error")
	}
}

func TestDay02OverflowProductPartTwo(t *testing.T) {
	if _, err := day02p02(strings.NewReader(`down 9223372036854775807
forward 2`)); err == nil {
		t.Fatal("expected overflow product part two error")
	}
}

func TestDay02OverflowPosition(t *testing.T) {
	if _, err := day02p01(strings.NewReader(`forward 9223372036854775807
forward 1`)); err == nil {
		t.Fatal("expected overflow position error")
	}
}

func TestDay02AimedOverflowBeforeCancellation(t *testing.T) {
	if _, err := day02p02(strings.NewReader(`down 9223372036854775807
forward 2
up 9223372036854775807`)); err == nil {
		t.Fatal("expected aimed overflow before cancellation error")
	}
}

func TestDay02UnknownCommand(t *testing.T) {
	if _, err := day02p01(strings.NewReader(`backward 1`)); err == nil {
		t.Fatal("expected unknown command error")
	}
}

func TestDay02ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day02p01},
		{"part2", day02p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
