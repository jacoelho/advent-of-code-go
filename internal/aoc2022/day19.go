package aoc2022

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type robotBlueprint struct {
	id    int64
	costs [4][3]int
}

func robotBlueprints(r io.Reader) ([]robotBlueprint, error) {
	pattern := regexp.MustCompile(`^Blueprint (\d+): Each ore robot costs (\d+) ore\. Each clay robot costs (\d+) ore\. Each obsidian robot costs (\d+) ore and (\d+) clay\. Each geode robot costs (\d+) ore and (\d+) obsidian\.$`)
	s := scanner.NewScanner(r, func(line []byte) (robotBlueprint, error) {
		match := pattern.FindSubmatch(line)
		if match == nil {
			return robotBlueprint{}, fmt.Errorf("invalid robot blueprint %q", line)
		}
		var values [7]int
		for i := range values {
			value, err := convert.ScanNumber[int32](match[i+1])
			if err != nil {
				return robotBlueprint{}, fmt.Errorf("blueprint value: %w", err)
			}
			if i > 0 && value > 10_000 {
				return robotBlueprint{}, fmt.Errorf("robot cost exceeds 10000")
			}
			values[i] = int(value)
		}
		if values[0] == 0 {
			return robotBlueprint{}, fmt.Errorf("blueprint ID must be positive")
		}
		return robotBlueprint{
			id: int64(values[0]),
			costs: [4][3]int{
				{values[1], 0, 0},
				{values[2], 0, 0},
				{values[3], values[4], 0},
				{values[5], 0, values[6]},
			},
		}, nil
	})
	var blueprints []robotBlueprint
	ids := make(map[int64]bool)
	for blueprint := range s.Values() {
		if len(blueprints) == 1000 || ids[blueprint.id] {
			return nil, fmt.Errorf("duplicate blueprint ID or more than 1000 blueprints")
		}
		ids[blueprint.id] = true
		blueprints = append(blueprints, blueprint)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(blueprints) == 0 {
		return nil, fmt.Errorf("empty blueprint input")
	}
	return blueprints, nil
}

func blueprintGeodes(blueprint robotBlueprint, minutes int, work *int) (int, error) {
	if minutes < 0 || minutes > 32 {
		return 0, fmt.Errorf("robot planning supports 0..32 minutes")
	}
	var maxSpend [3]int
	for _, cost := range blueprint.costs {
		for resource := range maxSpend {
			maxSpend[resource] = max(maxSpend[resource], cost[resource])
		}
	}
	type state struct {
		stock [3]int
		bots  [4]int
		time  int
	}
	seen := make(map[state]int)
	best := 0
	var explore func(state, int) error
	explore = func(current state, geodes int) error {
		*work++
		if *work > 20_000_000 || len(seen) >= 2_000_000 {
			return fmt.Errorf("robot planning exceeded its state or work limit")
		}
		best = max(best, geodes+current.bots[3]*current.time)
		// Even a free new geode robot each minute cannot exceed this bound.
		if geodes+current.bots[3]*current.time+current.time*(current.time-1)/2 <= best {
			return nil
		}
		for resource := range current.stock {
			cap := max(0, current.time*maxSpend[resource]-current.bots[resource]*(current.time-1))
			current.stock[resource] = min(current.stock[resource], cap)
		}
		if previous, exists := seen[current]; exists && previous >= geodes {
			return nil
		}
		seen[current] = geodes
		for robot := 3; robot >= 0; robot-- {
			if robot < 3 && current.bots[robot] >= maxSpend[robot] {
				continue
			}
			wait := 0
			possible := true
			for resource, cost := range blueprint.costs[robot] {
				missing := cost - current.stock[resource]
				if missing <= 0 {
					continue
				}
				if current.bots[resource] == 0 {
					possible = false
					break
				}
				wait = max(wait, (missing+current.bots[resource]-1)/current.bots[resource])
			}
			elapsed := wait + 1
			if !possible || elapsed >= current.time {
				continue
			}
			child := current
			child.time -= elapsed
			for resource, cost := range blueprint.costs[robot] {
				child.stock[resource] += current.bots[resource]*elapsed - cost
			}
			child.bots[robot]++
			if err := explore(child, geodes+current.bots[3]*elapsed); err != nil {
				return err
			}
		}
		return nil
	}
	if err := explore(state{bots: [4]int{1}, time: minutes}, 0); err != nil {
		return 0, err
	}
	return best, nil
}

func day19p01(r io.Reader) (string, error) {
	blueprints, err := robotBlueprints(r)
	if err != nil {
		return "", err
	}
	var total int64
	work := 0
	for _, blueprint := range blueprints {
		geodes, err := blueprintGeodes(blueprint, 24, &work)
		if err != nil {
			return "", fmt.Errorf("blueprint %d: %w", blueprint.id, err)
		}
		total += blueprint.id * int64(geodes)
	}
	return strconv.FormatInt(total, 10), nil
}

func day19p02(r io.Reader) (string, error) {
	blueprints, err := robotBlueprints(r)
	if err != nil {
		return "", err
	}
	total, work := 1, 0
	for _, blueprint := range blueprints[:min(3, len(blueprints))] {
		geodes, err := blueprintGeodes(blueprint, 32, &work)
		if err != nil {
			return "", fmt.Errorf("blueprint %d: %w", blueprint.id, err)
		}
		total *= geodes
	}
	return strconv.Itoa(total), nil
}
