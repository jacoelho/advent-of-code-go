package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
)

func stepDependencies(r io.Reader) (map[byte]uint32, error) {
	s := bufio.NewScanner(r)
	dependencies := make(map[byte]uint32)
	for s.Scan() {
		var before, after byte
		if _, err := fmt.Sscanf(s.Text(), "Step %c must be finished before step %c can begin.", &before, &after); err != nil {
			return nil, fmt.Errorf("invalid dependency %q: %w", s.Text(), err)
		}
		if before < 'A' || before > 'Z' || after < 'A' || after > 'Z' {
			return nil, fmt.Errorf("invalid step in %q", s.Text())
		}
		if _, exists := dependencies[before]; !exists {
			dependencies[before] = 0
		}
		dependencies[after] |= 1 << (before - 'A')
	}
	return dependencies, s.Err()
}

func day07p01(r io.Reader) (string, error) {
	dependencies, err := stepDependencies(r)
	if err != nil {
		return "", err
	}

	var completed uint32
	order := make([]byte, 0, len(dependencies))
	for len(order) < len(dependencies) {
		next := byte(0)
		for step := byte('A'); step <= 'Z'; step++ {
			parents, exists := dependencies[step]
			if exists && completed&(1<<(step-'A')) == 0 && parents&^completed == 0 {
				next = step
				break
			}
		}
		if next == 0 {
			return "", fmt.Errorf("step dependencies contain a cycle")
		}
		order = append(order, next)
		completed |= 1 << (next - 'A')
	}
	return string(order), nil
}

func assemblyTime(r io.Reader, workers, baseDuration int) (int, error) {
	dependencies, err := stepDependencies(r)
	if err != nil {
		return 0, err
	}
	if workers <= 0 || baseDuration < 0 {
		return 0, fmt.Errorf("invalid worker count or base duration")
	}

	var started, completed uint32
	running := make(map[byte]int)
	elapsed, done := 0, 0
	for done < len(dependencies) {
		for step := byte('A'); step <= 'Z' && len(running) < workers; step++ {
			parents, exists := dependencies[step]
			if exists && started&(1<<(step-'A')) == 0 && parents&^completed == 0 {
				running[step] = elapsed + baseDuration + int(step-'A') + 1
				started |= 1 << (step - 'A')
			}
		}
		if len(running) == 0 {
			return 0, fmt.Errorf("step dependencies contain a cycle")
		}
		elapsed = math.MaxInt
		for _, finish := range running {
			elapsed = min(elapsed, finish)
		}
		for step, finish := range running {
			if finish == elapsed {
				completed |= 1 << (step - 'A')
				delete(running, step)
				done++
			}
		}
	}
	return elapsed, nil
}

func day07p02(r io.Reader) (string, error) {
	elapsed, err := assemblyTime(r, 5, 60)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(elapsed), nil
}
