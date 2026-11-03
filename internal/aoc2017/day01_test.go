package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("1122"), Want: "3"},
		{Input: strings.NewReader("1111"), Want: "4"},
		{Input: strings.NewReader("1234"), Want: "0"},
		{Input: strings.NewReader("91212129"), Want: "9"},
		{Input: aoc.FileInput(t, 2017, 1), Want: "1144"},
	}

	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("1212"), Want: "6"},
		{Input: strings.NewReader("1221"), Want: "0"},
		{Input: strings.NewReader("123425"), Want: "4"},
		{Input: strings.NewReader("123123"), Want: "12"},
		{Input: strings.NewReader("12131415"), Want: "4"},
		{Input: aoc.FileInput(t, 2017, 1), Want: "1194"},
	}

	aoc.AOCTest(t, day01p02, tests)
}
