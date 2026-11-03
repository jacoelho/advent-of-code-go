package aoc2020

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"
)

func TestPortedDaysPropagateReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"01/1", day01p01}, {"01/2", day01p02}, {"02/1", day02p01}, {"02/2", day02p02},
		{"03/1", day03p01}, {"03/2", day03p02}, {"04/1", day04p01}, {"04/2", day04p02},
		{"05/1", day05p01}, {"05/2", day05p02}, {"06/1", day06p01}, {"06/2", day06p02},
		{"22/1", day22p01}, {"22/2", day22p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestCombatInputErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, solve := range []func(io.Reader) (string, error){day22p01, day22p02} {
		reader := io.MultiReader(strings.NewReader("Player 1:\n1\n\n"), iotest.ErrReader(failure))
		if _, err := solve(reader); !errors.Is(err, failure) {
			t.Fatalf("got %v, want input failure", err)
		}
		_, err := solve(strings.NewReader("Player 1:\nnope\n\nPlayer 2:\n2"))
		if _, ok := errors.AsType[*strconv.NumError](err); !ok {
			t.Fatalf("got %v, want number conversion error", err)
		}
	}
}

func TestTilePathReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, solve := range []func(io.Reader) (string, error){day24p01, day24p02} {
		for _, prefix := range []string{"", "e\n"} {
			reader := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
			if _, err := solve(reader); !errors.Is(err, failure) {
				t.Fatalf("prefix %q: got %v, want input failure", prefix, err)
			}
		}
	}
}

func TestPortedDaysRejectMalformedInput(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"expense", "not a number", day01p01},
		{"policy", "0-2 a: ab", day02p02},
		{"tree map", "..\n...", day03p01},
		{"passport", "byr", day04p01},
		{"boarding pass", "FBFBBFFRLF", day05p01},
		{"customs form", "aB", day06p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected invalid-input error")
			}
		})
	}
}
