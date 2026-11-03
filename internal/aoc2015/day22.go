package aoc2015

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

type wizardMode int

const (
	wizardNormal wizardMode = iota
	wizardHard
)

type wizardState struct {
	hp, mana, boss           int
	shield, poison, recharge int
}

func (s *wizardState) tick() {
	if s.shield > 0 {
		s.shield--
	}
	if s.poison > 0 {
		s.boss -= 3
		s.poison--
	}
	if s.recharge > 0 {
		s.mana += 101
		s.recharge--
	}
}

func leastMana(initial wizardState, bossDamage int, mode wizardMode) (int, error) {
	type entry struct {
		state wizardState
		spent int
	}
	queue := collections.NewHeap(func(a, b entry) bool { return a.spent < b.spent })
	queue.Push(entry{state: initial})
	distance := map[wizardState]int{initial: 0}
	best := math.MaxInt
	for current := range queue.PopSeq() {
		if current.spent >= best {
			return best, nil
		}
		if distance[current.state] != current.spent {
			continue
		}
		state := current.state
		if mode == wizardHard {
			state.hp--
		}
		if state.hp <= 0 {
			continue
		}
		state.tick()
		if state.boss <= 0 {
			return current.spent, nil
		}
		for spell, cost := range [...]int{53, 73, 113, 173, 229} {
			spent := current.spent + cost
			if state.mana < cost || spent >= best {
				continue
			}
			next := state
			next.mana -= cost
			switch spell {
			case 0:
				next.boss -= 4
			case 1:
				next.boss -= 2
				next.hp += 2
			case 2:
				if next.shield != 0 {
					continue
				}
				next.shield = 6
			case 3:
				if next.poison != 0 {
					continue
				}
				next.poison = 6
			case 4:
				if next.recharge != 0 {
					continue
				}
				next.recharge = 5
			}
			if next.boss > 0 {
				next.tick()
			}
			if next.boss <= 0 {
				best = min(best, spent)
				continue
			}
			armor := 0
			if next.shield > 0 {
				armor = 7
			}
			next.hp -= max(1, bossDamage-armor)
			if next.hp <= 0 {
				continue
			}
			if previous, seen := distance[next]; seen && previous <= spent {
				continue
			}
			distance[next] = spent
			queue.Push(entry{next, spent})
		}
	}
	if best == math.MaxInt {
		return 0, fmt.Errorf("no winning spell sequence")
	}
	return best, nil
}

func wizardFight(r io.Reader, mode wizardMode) (string, error) {
	boss, err := readBoss(r)
	if err != nil {
		return "", err
	}
	spent, err := leastMana(wizardState{hp: 50, mana: 500, boss: boss.hp}, boss.damage, mode)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(spent), nil
}

func day22p01(r io.Reader) (string, error) {
	return wizardFight(r, wizardNormal)
}

func day22p02(r io.Reader) (string, error) {
	return wizardFight(r, wizardHard)
}
