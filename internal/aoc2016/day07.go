package aoc2016

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type ipv7Address struct{ supernets, hypernets []string }

func parseIPv7(line []byte) (ipv7Address, error) {
	var address ipv7Address
	start, inside := 0, false
	for i, letter := range line {
		switch letter {
		case '[':
			if inside {
				return address, fmt.Errorf("nested IPv7 brackets")
			}
			address.supernets = append(address.supernets, string(line[start:i]))
			start, inside = i+1, true
		case ']':
			if !inside {
				return address, fmt.Errorf("unmatched IPv7 closing bracket")
			}
			address.hypernets = append(address.hypernets, string(line[start:i]))
			start, inside = i+1, false
		default:
			if letter < 'a' || letter > 'z' {
				return address, fmt.Errorf("invalid IPv7 letter %q", letter)
			}
		}
	}
	if inside {
		return address, fmt.Errorf("unclosed IPv7 bracket")
	}
	address.supernets = append(address.supernets, string(line[start:]))
	return address, nil
}

func containsABBA(text string) bool {
	for i := 0; i+3 < len(text); i++ {
		if text[i] != text[i+1] && text[i] == text[i+3] && text[i+1] == text[i+2] {
			return true
		}
	}
	return false
}

func supportsTLS(address ipv7Address) bool {
	return !xiter.FromSlice(address.hypernets).Any(containsABBA) &&
		xiter.FromSlice(address.supernets).Any(containsABBA)
}

func supportsSSL(address ipv7Address) bool {
	pairs := collections.NewSet[[2]byte]()
	for _, text := range address.supernets {
		for i := 0; i+2 < len(text); i++ {
			if text[i] != text[i+1] && text[i] == text[i+2] {
				pairs.Add([2]byte{text[i], text[i+1]})
			}
		}
	}
	for _, text := range address.hypernets {
		for i := 0; i+2 < len(text); i++ {
			if text[i] == text[i+2] && pairs.Contains([2]byte{text[i+1], text[i]}) {
				return true
			}
		}
	}
	return false
}

func countIPv7(r io.Reader, supports func(ipv7Address) bool) (string, error) {
	s := scanner.NewScanner(r, parseIPv7)
	count := xiter.From(s.Values()).Filter(supports).Count()
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.Itoa(count), nil
}

func day07p01(r io.Reader) (string, error) {
	return countIPv7(r, supportsTLS)
}

func day07p02(r io.Reader) (string, error) {
	return countIPv7(r, supportsSSL)
}
