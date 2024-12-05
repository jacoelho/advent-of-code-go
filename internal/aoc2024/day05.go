package aoc2024

import (
	"bufio"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

func parseSleighLaunchSafetyRules(r io.Reader) (map[[2]int]bool, [][]int, error) {
	var updatePages [][]int
	rules := make(map[[2]int]bool)

	s := bufio.NewScanner(r)
	for s.Scan() {
		text := s.Text()

		switch {
		case strings.Contains(text, "|"):
			v := xslices.Map(aoc.MustAtoi, strings.Split(text, "|"))
			rules[[2]int{v[0], v[1]}] = true
			rules[[2]int{v[1], v[0]}] = false

		case strings.Contains(text, ","):
			v := xslices.Map(aoc.MustAtoi, strings.Split(text, ","))
			updatePages = append(updatePages, v)
		}
	}
	return rules, updatePages, s.Err()
}

func isOrdered(rules map[[2]int]bool, update []int) bool {
	for i, u := range update {
		for _, v := range update[i+1:] {
			if ordered, found := rules[[2]int{u, v}]; found && !ordered {
				return false
			}
		}
	}
	return true
}

func day05p01(r io.Reader) (string, error) {
	rules, updatePages := aoc.Must2(parseSleighLaunchSafetyRules(r))

	var total int
	for _, updatePage := range updatePages {
		if isOrdered(rules, updatePage) {
			total += updatePage[len(updatePage)/2]
		}
	}
	return strconv.Itoa(total), nil
}

func day05p02(r io.Reader) (string, error) {
	rules, updatePages := aoc.Must2(parseSleighLaunchSafetyRules(r))

	var total int
	for _, updatePage := range updatePages {
		if isOrdered(rules, updatePage) {
			continue
		}

		slices.SortFunc(updatePage, func(i, j int) int {
			if ordered, found := rules[[2]int{i, j}]; found {
				if ordered {
					return -1
				}
				return 1
			}
			return 0
		})
		total += updatePage[len(updatePage)/2]
	}
	return strconv.Itoa(total), nil
}
