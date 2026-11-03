package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestCorporatePasswordPolicy(t *testing.T) {
	for _, password := range []string{"hijklmmn", "abbceffg", "abbcegjk", "abcaaaaa"} {
		if validCorporatePassword([]byte(password)) {
			t.Errorf("password %q should be invalid", password)
		}
	}
	for _, password := range []string{"abcdffaa", "ghjaabcc"} {
		if !validCorporatePassword([]byte(password)) {
			t.Errorf("password %q should be valid", password)
		}
	}
}

func Test_day11p01(t *testing.T) {
	aoc.AOCTest(t, day11p01, []aoc.TestInput{
		{Input: strings.NewReader("abcdefgh"), Want: "abcdffaa"},
		{Input: strings.NewReader("ghijklmn"), Want: "ghjaabcc"},
		{Input: aoc.FileInput(t, 2015, 11), Want: "cqjxxyzz"},
	})
}

func Test_day11p02(t *testing.T) {
	aoc.AOCTest(t, day11p02, []aoc.TestInput{
		{Input: strings.NewReader("abcdefgh"), Want: "abcdffbb"},
		{Input: aoc.FileInput(t, 2015, 11), Want: "cqkaabcc"},
	})
}
