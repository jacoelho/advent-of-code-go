package aoc2020

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type passwordPolicy struct {
	first, second int
	letter        byte
	password      string
}

func parsePasswordPolicy(line []byte) (passwordPolicy, error) {
	policy, password, ok := strings.Cut(string(line), ": ")
	bounds, letter, letterOK := strings.Cut(policy, " ")
	low, high, boundsOK := strings.Cut(bounds, "-")
	first, firstErr := strconv.Atoi(low)
	second, secondErr := strconv.Atoi(high)
	if !ok || !letterOK || !boundsOK || firstErr != nil || secondErr != nil || first < 1 || second < first || len(letter) != 1 || letter[0] < 'a' || letter[0] > 'z' {
		return passwordPolicy{}, fmt.Errorf("invalid password policy: %q", line)
	}
	return passwordPolicy{first, second, letter[0], password}, nil
}

func day02p01(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parsePasswordPolicy)
	valid := 0
	for p := range s.Values() {
		count := strings.Count(p.password, string(p.letter))
		if count >= p.first && count <= p.second {
			valid++
		}
	}
	return strconv.Itoa(valid), s.Err()
}

func day02p02(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parsePasswordPolicy)
	valid := 0
	for p := range s.Values() {
		first := p.first <= len(p.password) && p.password[p.first-1] == p.letter
		second := p.second <= len(p.password) && p.password[p.second-1] == p.letter
		if first != second {
			valid++
		}
	}
	return strconv.Itoa(valid), s.Err()
}
