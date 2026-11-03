package aoc2016

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

type chipDestination struct {
	kind string
	id   int
}
type botRule struct{ low, high chipDestination }
type chipDelivery struct {
	destination chipDestination
	value       int
}

func botFactory(r io.Reader, comparison [2]int) (int, map[int][]int, error) {
	valuePattern := regexp.MustCompile(`^value (\d+) goes to bot (\d+)$`)
	rulePattern := regexp.MustCompile(`^bot (\d+) gives low to (bot|output) (\d+) and high to (bot|output) (\d+)$`)
	rules := make(map[int]botRule)
	queue := collections.NewDeque[chipDelivery](64)
	s := bufio.NewScanner(r)
	for s.Scan() {
		if match := valuePattern.FindStringSubmatch(s.Text()); match != nil {
			value, err := strconv.Atoi(match[1])
			if err != nil {
				return -1, nil, fmt.Errorf("invalid chip value: %w", err)
			}
			bot, err := strconv.Atoi(match[2])
			if err != nil {
				return -1, nil, fmt.Errorf("invalid bot ID: %w", err)
			}
			queue.PushBack(chipDelivery{chipDestination{"bot", bot}, value})
			continue
		}
		match := rulePattern.FindStringSubmatch(s.Text())
		if match == nil {
			return -1, nil, fmt.Errorf("invalid bot instruction %q", s.Text())
		}
		var ids [3]int
		for i, field := range []int{1, 3, 5} {
			value, err := strconv.Atoi(match[field])
			if err != nil {
				return -1, nil, fmt.Errorf("invalid destination ID: %w", err)
			}
			ids[i] = value
		}
		if _, exists := rules[ids[0]]; exists {
			return -1, nil, fmt.Errorf("duplicate rule for bot %d", ids[0])
		}
		rules[ids[0]] = botRule{chipDestination{match[2], ids[1]}, chipDestination{match[4], ids[2]}}
	}
	if err := s.Err(); err != nil {
		return -1, nil, err
	}
	if queue.Size() == 0 {
		return -1, nil, fmt.Errorf("no input chips")
	}
	chips := make(map[int]int)
	outputs := make(map[int][]int)
	comparer := -1
	for deliveries := 0; queue.Size() > 0; deliveries++ {
		if deliveries == 1_000_000 {
			return -1, nil, fmt.Errorf("bot factory exceeds 1000000 deliveries")
		}
		delivery, _ := queue.PopFront()
		id := delivery.destination.id
		if delivery.destination.kind == "output" {
			outputs[id] = append(outputs[id], delivery.value)
			continue
		}
		previous, holding := chips[id]
		if !holding {
			chips[id] = delivery.value
			continue
		}
		delete(chips, id)
		low, high := min(previous, delivery.value), max(previous, delivery.value)
		if comparison == [2]int{low, high} {
			if comparer >= 0 && comparer != id {
				return -1, nil, fmt.Errorf("multiple bots compare chips %v", comparison)
			}
			comparer = id
		}
		rule, exists := rules[id]
		if !exists {
			return -1, nil, fmt.Errorf("bot %d has two chips but no routing rule", id)
		}
		queue.PushBack(chipDelivery{rule.low, low})
		queue.PushBack(chipDelivery{rule.high, high})
	}
	return comparer, outputs, nil
}

func day10p01(r io.Reader) (string, error) {
	bot, _, err := botFactory(r, [2]int{17, 61})
	if err != nil {
		return "", err
	}
	if bot < 0 {
		return "", fmt.Errorf("no bot compares chips 17 and 61")
	}
	return strconv.Itoa(bot), nil
}

func day10p02(r io.Reader) (string, error) {
	_, outputs, err := botFactory(r, [2]int{-1, -1})
	if err != nil {
		return "", err
	}
	product := big.NewInt(1)
	for output := range 3 {
		if len(outputs[output]) != 1 {
			return "", fmt.Errorf("output %d must contain exactly one chip", output)
		}
		product.Mul(product, big.NewInt(int64(outputs[output][0])))
	}
	return product.String(), nil
}
