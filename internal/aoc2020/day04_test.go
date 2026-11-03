package aoc2020

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day04p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`ecl:gry pid:860033327 eyr:2020 hcl:#fffffd
byr:1937 iyr:2017 cid:147 hgt:183cm

iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884
hcl:#cfa07d byr:1929

hcl:#ae17e1 iyr:2013
eyr:2024
ecl:brn pid:760753108 byr:1931
hgt:179cm

hcl:#cfa07d eyr:2025 pid:166559648
iyr:2011 ecl:brn hgt:59in`), Want: "2"},
		{Input: aoc.FileInput(t, 2020, 4), Want: "264"},
	}
	aoc.AOCTest(t, day04p01, tests)
}

func Test_day04p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`ecl:gry pid:860033327 eyr:2020 hcl:#fffffd
byr:1937 iyr:2017 cid:147 hgt:183cm

iyr:2013 ecl:amb cid:350 eyr:2023 pid:028048884
hcl:#cfa07d byr:1929

hcl:#ae17e1 iyr:2013
eyr:2024
ecl:brn pid:760753108 byr:1931
hgt:179cm

hcl:#cfa07d eyr:2025 pid:166559648
iyr:2011 ecl:brn hgt:59in`), Want: "2"},
		{Input: aoc.FileInput(t, 2020, 4), Want: "224"},
	}
	aoc.AOCTest(t, day04p02, tests)
}

func TestPassportFieldRules(t *testing.T) {
	const valid = "byr:2002 iyr:2010 eyr:2030 hgt:150cm hcl:#123abc ecl:brn pid:000000001"
	for _, tt := range []struct{ old, replacement, want string }{
		{"byr:2002", "byr:1920", "1"},
		{"byr:2002", "byr:2003", "0"},
		{"byr:2002", "byr:02002", "0"},
		{"iyr:2010", "iyr:2009", "0"},
		{"eyr:2030", "eyr:2031", "0"},
		{"hgt:150cm", "hgt:193cm", "1"},
		{"hgt:150cm", "hgt:59in", "1"},
		{"hgt:150cm", "hgt:76in", "1"},
		{"hgt:150cm", "hgt:149cm", "0"},
		{"hgt:150cm", "hgt:77in", "0"},
		{"hcl:#123abc", "hcl:#123abz", "0"},
		{"hcl:#123abc", "hcl:#123ABC", "0"},
		{"ecl:brn", "ecl:wat", "0"},
		{"pid:000000001", "pid:0000000001", "0"},
		{"pid:000000001", "pid:00000000a", "0"},
	} {
		t.Run(tt.replacement, func(t *testing.T) {
			input := strings.Replace(valid, tt.old, tt.replacement, 1)
			aoc.AOCTest(t, day04p01, []aoc.TestInput{{Input: strings.NewReader(input), Want: "1"}})
			aoc.AOCTest(t, day04p02, []aoc.TestInput{{Input: strings.NewReader(input), Want: tt.want}})
		})
	}
	aoc.AOCTest(t, day04p02, []aoc.TestInput{{Input: strings.NewReader(valid), Want: "1"}})
}

func TestDay04ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"04/1", day04p01},
		{"04/2", day04p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestDay04MalformedInput(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"passport", "byr", day04p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected invalid-input error")
			}
		})
	}
}
