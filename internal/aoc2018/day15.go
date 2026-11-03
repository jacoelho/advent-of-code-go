package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

type caveUnit struct {
	position, health int
	kind             byte
}

type caveBattle struct {
	width, height int
	open          []bool
	units         []caveUnit
}

type caveOutcome struct {
	score              int
	winner             byte
	elfDied, stalemate bool
}

type elfLossPolicy uint8

const (
	allowElfLoss elfLossPolicy = iota
	stopOnElfLoss
)

func parseCaveBattle(r io.Reader) (caveBattle, error) {
	var battle caveBattle
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if battle.height == 0 {
			battle.width = len(line)
		}
		if len(line) == 0 || len(line) != battle.width {
			return battle, fmt.Errorf("invalid combat map row width")
		}
		for _, cell := range []byte(line) {
			position := len(battle.open)
			switch cell {
			case '#':
				battle.open = append(battle.open, false)
			case '.':
				battle.open = append(battle.open, true)
			case 'E', 'G':
				battle.open = append(battle.open, true)
				battle.units = append(battle.units, caveUnit{position: position, health: 200, kind: cell})
			default:
				return battle, fmt.Errorf("invalid combat cell %q", cell)
			}
		}
		battle.height++
		if len(battle.open) > 100_000 {
			return battle, fmt.Errorf("combat map exceeds 100000 cells")
		}
	}
	if err := s.Err(); err != nil {
		return battle, err
	}
	if len(battle.units) == 0 {
		return battle, fmt.Errorf("combat map has no units")
	}

	return battle, nil
}

func (battle caveBattle) neighbours(position int) [4]int {
	x, y := position%battle.width, position/battle.width
	result := [4]int{-1, -1, -1, -1}
	if y > 0 {
		result[0] = position - battle.width
	}
	if x > 0 {
		result[1] = position - 1
	}
	if x+1 < battle.width {
		result[2] = position + 1
	}
	if y+1 < battle.height {
		result[3] = position + battle.width
	}

	return result
}

func (battle caveBattle) distances(start int, occupied []int) []int {
	distance := make([]int, len(battle.open))
	for i := range distance {
		distance[i] = -1
	}
	distance[start] = 0
	queue := collections.NewDeque[int](64)
	queue.PushBack(start)
	for position, ok := queue.PopFront(); ok; position, ok = queue.PopFront() {
		for _, next := range battle.neighbours(position) {
			if next >= 0 && battle.open[next] && occupied[next] < 0 && distance[next] < 0 {
				distance[next] = distance[position] + 1
				queue.PushBack(next)
			}
		}
	}

	return distance
}

func (battle caveBattle) attackTarget(unit caveUnit, units []caveUnit, occupied []int) int {
	target := -1
	for _, position := range battle.neighbours(unit.position) {
		if position < 0 {
			continue
		}
		id := occupied[position]
		if id >= 0 && units[id].kind != unit.kind && (target < 0 || units[id].health < units[target].health) {
			target = id
		}
	}

	return target
}

func (battle caveBattle) moveTarget(unit caveUnit, units []caveUnit, occupied []int) int {
	distance := battle.distances(unit.position, occupied)
	target, bestDistance := -1, len(battle.open)+1
	for _, enemy := range units {
		if enemy.health <= 0 || enemy.kind == unit.kind {
			continue
		}
		for _, position := range battle.neighbours(enemy.position) {
			if position < 0 || occupied[position] >= 0 || distance[position] < 0 {
				continue
			}
			if distance[position] < bestDistance || distance[position] == bestDistance && position < target {
				target, bestDistance = position, distance[position]
			}
		}
	}
	if target < 0 {
		return unit.position
	}

	// Resolve the destination first; reverse distances then select the first step in reading order.
	remaining := battle.distances(target, occupied)
	for _, position := range battle.neighbours(unit.position) {
		if position >= 0 && occupied[position] < 0 && remaining[position] == bestDistance-1 {
			return position
		}
	}

	return unit.position
}

func fightCave(battle caveBattle, elfPower int, policy elfLossPolicy) (caveOutcome, error) {
	units := slices.Clone(battle.units)
	occupied := make([]int, len(battle.open))
	for i := range occupied {
		occupied[i] = -1
	}
	counts := map[byte]int{'E': 0, 'G': 0}
	for id, unit := range units {
		occupied[unit.position] = id
		counts[unit.kind]++
	}
	elfDied := false
	for round := 0; round < 100_000; round++ {
		order := make([]int, 0, len(units))
		for id, unit := range units {
			if unit.health > 0 {
				order = append(order, id)
			}
		}
		slices.SortFunc(order, func(a, b int) int { return units[a].position - units[b].position })
		changed := false
		for _, id := range order {
			unit := &units[id]
			if unit.health <= 0 {
				continue
			}
			enemy := byte('E')
			if unit.kind == 'E' {
				enemy = 'G'
			}
			if counts[enemy] == 0 {
				health := 0
				for _, survivor := range units {
					health += max(0, survivor.health)
				}
				return caveOutcome{score: round * health, winner: unit.kind, elfDied: elfDied}, nil
			}

			target := battle.attackTarget(*unit, units, occupied)
			if target < 0 {
				position := battle.moveTarget(*unit, units, occupied)
				if position != unit.position {
					occupied[unit.position] = -1
					unit.position = position
					occupied[position] = id
					changed = true
				}
				target = battle.attackTarget(*unit, units, occupied)
			}
			if target < 0 {
				continue
			}
			power := 3
			if unit.kind == 'E' {
				power = elfPower
			}
			units[target].health -= power
			changed = true
			if units[target].health <= 0 {
				occupied[units[target].position] = -1
				counts[units[target].kind]--
				if units[target].kind == 'E' {
					elfDied = true
					if policy == stopOnElfLoss {
						return caveOutcome{elfDied: true}, nil
					}
				}
			}
		}
		if !changed {
			return caveOutcome{stalemate: true, elfDied: elfDied}, nil
		}
	}

	return caveOutcome{}, fmt.Errorf("combat exceeds 100000 rounds")
}

func day15p01(r io.Reader) (string, error) {
	battle, err := parseCaveBattle(r)
	if err != nil {
		return "", err
	}

	result, err := fightCave(battle, 3, allowElfLoss)
	if err != nil {
		return "", err
	}
	if result.stalemate {
		return "", fmt.Errorf("combat ended in stalemate")
	}

	return strconv.Itoa(result.score), nil
}

func day15p02(r io.Reader) (string, error) {
	battle, err := parseCaveBattle(r)
	if err != nil {
		return "", err
	}

	// At 200 power every goblin dies in one hit; larger powers cannot change a battle.
	for power := 3; power <= 200; power++ {
		result, err := fightCave(battle, power, stopOnElfLoss)
		if err != nil {
			return "", err
		}
		if result.winner == 'E' && !result.elfDied {
			return strconv.Itoa(result.score), nil
		}
	}

	return "", fmt.Errorf("no elf-survival victory exists, even with one-hit attacks")
}
