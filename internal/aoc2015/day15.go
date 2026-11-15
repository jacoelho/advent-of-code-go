package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func cookieScore(r io.Reader, calories int64) (string, error) {
	pattern := regexp.MustCompile(`^(\w+): capacity (-?\d+), durability (-?\d+), flavor (-?\d+), texture (-?\d+), calories (\d+)$`)
	var ingredients [][5]int64
	names := make(map[string]bool)
	s := bufio.NewScanner(r)
	for s.Scan() {
		match := pattern.FindStringSubmatch(s.Text())
		if match == nil || names[match[1]] {
			return "", fmt.Errorf("invalid or duplicate ingredient %q", s.Text())
		}
		names[match[1]] = true
		var ingredient [5]int64
		for i := range ingredient {
			value, err := strconv.ParseInt(match[i+2], 10, 32)
			if err != nil {
				return "", fmt.Errorf("invalid ingredient property: %w", err)
			}
			ingredient[i] = value
		}
		ingredients = append(ingredients, ingredient)
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(ingredients) == 0 {
		return "", fmt.Errorf("empty ingredient list")
	}
	var best int64
	var overflow error
	var visit func(int, int64, [5]int64)
	visit = func(index int, remaining int64, totals [5]int64) {
		if index == len(ingredients)-1 {
			for property, value := range ingredients[index] {
				totals[property] += remaining * value
			}
			if calories >= 0 && totals[4] != calories {
				return
			}
			score := int64(1)
			for _, value := range totals[:4] {
				if value <= 0 {
					return
				}
				next, ok := xmath.CheckedMul(score, value)
				if !ok {
					overflow = fmt.Errorf("cookie score overflows int64")
					return
				}
				score = next
			}
			best = max(best, score)
			return
		}
		for quantity := int64(0); quantity <= remaining && overflow == nil; quantity++ {
			next := totals
			for property, value := range ingredients[index] {
				next[property] += quantity * value
			}
			visit(index+1, remaining-quantity, next)
		}
	}
	visit(0, 100, [5]int64{})
	if overflow != nil {
		return "", overflow
	}
	return strconv.FormatInt(best, 10), nil
}

func day15p01(r io.Reader) (string, error) {
	return cookieScore(r, -1)
}

func day15p02(r io.Reader) (string, error) {
	return cookieScore(r, 500)
}
