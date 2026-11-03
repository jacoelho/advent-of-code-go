package aoc2024

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type guardStep struct {
	next, ahead int
}

type guardMap struct {
	steps []guardStep
	start int
}

type guardObstacle struct {
	cell, before int
}

func parseGuardMap(r io.Reader) (guardMap, error) {
	s := scanner.NewScanner(r, func(b []byte) (string, error) { return string(b), nil })
	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return guardMap{}, err
	}

	offsets := make([]int, len(rows))
	cellCount, start := 0, -1
	multiple := false
	for y, row := range rows {
		if len(row) > int(^uint(0)>>1)/4-cellCount {
			return guardMap{}, fmt.Errorf("guard map has too many cells")
		}
		offsets[y] = cellCount
		for x := range len(row) {
			if row[x] == '^' {
				if start >= 0 {
					multiple = true
				} else {
					start = 4 * (cellCount + x)
				}
			}
		}
		cellCount += len(row)
	}
	if multiple {
		return guardMap{}, fmt.Errorf("guard map has multiple guards")
	}
	if start < 0 {
		return guardMap{}, fmt.Errorf("guard not found")
	}

	// Each cell has north, east, south, and west states, including turns in place.
	m := guardMap{steps: make([]guardStep, 4*cellCount), start: start}
	dx, dy := [4]int{0, 1, 0, -1}, [4]int{-1, 0, 1, 0}
	for y, row := range rows {
		for x := range len(row) {
			cell := offsets[y] + x
			for direction := range 4 {
				step := guardStep{next: -1, ahead: -1}
				nx, ny := x+dx[direction], y+dy[direction]
				if ny >= 0 && ny < len(rows) && nx >= 0 && nx < len(rows[ny]) && rows[ny][nx] != 0 {
					step.ahead = offsets[ny] + nx
					if rows[ny][nx] == '#' {
						step.next = 4*cell + (direction+1)%4
					} else {
						step.next = 4*step.ahead + direction
					}
				}
				m.steps[4*cell+direction] = step
			}
		}
	}
	return m, nil
}

func (m guardMap) obstacles(marks []int) []guardObstacle {
	visited := make([]bool, len(m.steps)/4)
	visited[m.start/4] = true
	var obstacles []guardObstacle
	for state := m.start; state >= 0 && marks[state] == 0; {
		marks[state] = 1
		next := m.steps[state].next
		if next >= 0 && !visited[next/4] {
			visited[next/4] = true
			obstacles = append(obstacles, guardObstacle{cell: next / 4, before: state})
		}
		state = next
	}
	return obstacles
}

func day06p01(r io.Reader) (string, error) {
	m, err := parseGuardMap(r)
	if err != nil {
		return "", err
	}
	marks := make([]int, len(m.steps))
	return strconv.Itoa(len(m.obstacles(marks)) + 1), nil
}

func day06p02(r io.Reader) (string, error) {
	m, err := parseGuardMap(r)
	if err != nil {
		return "", err
	}
	marks := make([]int, len(m.steps))
	obstacles := m.obstacles(marks)
	var count int
	for i, obstacle := range obstacles {
		generation := i + 2
		// The unchanged prefix reaches this state; only the suffix needs replaying.
		state := obstacle.before
		for state >= 0 && marks[state] != generation {
			marks[state] = generation
			step := m.steps[state]
			if step.ahead == obstacle.cell {
				state = state&^3 | (state+1)&3
			} else {
				state = step.next
			}
		}
		if state >= 0 {
			count++
		}
	}
	return strconv.Itoa(count), nil
}
