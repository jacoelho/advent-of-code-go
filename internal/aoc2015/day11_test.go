package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

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

func TestDay11ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day11p01", day11p01},
		{"day11p02", day11p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay11MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"password length", "abc", day11p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
