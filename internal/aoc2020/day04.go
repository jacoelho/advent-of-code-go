package aoc2020

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func parsePassports(r io.Reader) ([]map[string]string, error) {
	var passports []map[string]string
	fields := make(map[string]string)
	s := scanner.NewScanner(r, func(line []byte) (string, error) { return string(line), nil })
	for line := range s.Values() {
		if strings.TrimSpace(line) == "" {
			if len(fields) > 0 {
				passports = append(passports, fields)
				fields = make(map[string]string)
			}
			continue
		}
		for field := range strings.FieldsSeq(line) {
			key, value, ok := strings.Cut(field, ":")
			if !ok || key == "" || value == "" {
				return nil, fmt.Errorf("invalid passport field %q", field)
			}
			if _, exists := fields[key]; exists {
				return nil, fmt.Errorf("duplicate passport field %q", key)
			}
			fields[key] = value
		}
	}
	if len(fields) > 0 {
		passports = append(passports, fields)
	}
	return passports, s.Err()
}

func passportFieldsPresent(fields map[string]string) bool {
	for _, key := range []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"} {
		if _, ok := fields[key]; !ok {
			return false
		}
	}
	return true
}

func passportNumber(s string, low, high int) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	n, err := strconv.Atoi(s)
	return err == nil && n >= low && n <= high
}

func passportFieldsValid(fields map[string]string) bool {
	for _, year := range []struct {
		key       string
		low, high int
	}{{"byr", 1920, 2002}, {"iyr", 2010, 2020}, {"eyr", 2020, 2030}} {
		if value := fields[year.key]; len(value) != 4 || !passportNumber(value, year.low, year.high) {
			return false
		}
	}
	height := fields["hgt"]
	if cm, ok := strings.CutSuffix(height, "cm"); ok {
		if !passportNumber(cm, 150, 193) {
			return false
		}
	} else if inches, ok := strings.CutSuffix(height, "in"); ok {
		if !passportNumber(inches, 59, 76) {
			return false
		}
	} else {
		return false
	}
	hair := fields["hcl"]
	if len(hair) != 7 || hair[0] != '#' {
		return false
	}
	for _, c := range hair[1:] {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	switch fields["ecl"] {
	case "amb", "blu", "brn", "gry", "grn", "hzl", "oth":
	default:
		return false
	}
	return len(fields["pid"]) == 9 && passportNumber(fields["pid"], 0, 999999999)
}

func countPassports(r io.Reader, valid func(map[string]string) bool) (string, error) {
	passports, err := parsePassports(r)
	if err != nil {
		return "", err
	}
	count := xiter.FromSlice(passports).Filter(valid).Count()
	return strconv.Itoa(count), nil
}

func day04p01(r io.Reader) (string, error) {
	return countPassports(r, passportFieldsPresent)
}

func day04p02(r io.Reader) (string, error) {
	return countPassports(r, passportFieldsValid)
}
