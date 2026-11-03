package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/bits"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type combatStats struct{ hp, damage, armor int }

func readBoss(r io.Reader) (combatStats, error) {
	values := make(map[string]int)
	s := bufio.NewScanner(r)
	for s.Scan() {
		name, raw, ok := strings.Cut(s.Text(), ": ")
		if !ok || (name != "Hit Points" && name != "Damage" && name != "Armor") {
			return combatStats{}, fmt.Errorf("invalid boss stat %q", s.Text())
		}
		value, err := convert.ScanNumber[int]([]byte(raw))
		if err != nil {
			return combatStats{}, fmt.Errorf("boss %s: %w", name, err)
		}
		if value < 0 {
			return combatStats{}, fmt.Errorf("invalid %s value %q", name, raw)
		}
		if _, exists := values[name]; exists {
			return combatStats{}, fmt.Errorf("duplicate boss stat %s", name)
		}
		values[name] = value
	}
	if err := s.Err(); err != nil {
		return combatStats{}, err
	}
	if values["Hit Points"] <= 0 || values["Damage"] <= 0 {
		return combatStats{}, fmt.Errorf("boss needs positive hit points and damage")
	}
	return combatStats{values["Hit Points"], values["Damage"], values["Armor"]}, nil
}

func playerWins(player, boss combatStats) bool {
	attack := max(1, player.damage-boss.armor)
	defend := max(1, boss.damage-player.armor)
	return (boss.hp-1)/attack <= (player.hp-1)/defend
}

func equipmentCosts(r io.Reader) (int, int, error) {
	boss, err := readBoss(r)
	if err != nil {
		return 0, 0, err
	}
	type item struct{ cost, damage, armor int }
	weapons := []item{{8, 4, 0}, {10, 5, 0}, {25, 6, 0}, {40, 7, 0}, {74, 8, 0}}
	armor := []item{{}, {13, 0, 1}, {31, 0, 2}, {53, 0, 3}, {75, 0, 4}, {102, 0, 5}}
	rings := []item{{25, 1, 0}, {50, 2, 0}, {100, 3, 0}, {20, 0, 1}, {40, 0, 2}, {80, 0, 3}}
	cheapest, dearest := math.MaxInt, -1
	for _, weapon := range weapons {
		for _, suit := range armor {
			for mask := range 1 << len(rings) {
				if bits.OnesCount(uint(mask)) > 2 {
					continue
				}
				cost := weapon.cost + suit.cost
				player := combatStats{100, weapon.damage, suit.armor}
				for i, ring := range rings {
					if mask&(1<<i) != 0 {
						cost += ring.cost
						player.damage += ring.damage
						player.armor += ring.armor
					}
				}
				if playerWins(player, boss) {
					cheapest = min(cheapest, cost)
				} else {
					dearest = max(dearest, cost)
				}
			}
		}
	}
	return cheapest, dearest, nil
}

func day21p01(r io.Reader) (string, error) {
	cheapest, _, err := equipmentCosts(r)
	if err != nil {
		return "", err
	}
	if cheapest == math.MaxInt {
		return "", fmt.Errorf("no winning equipment")
	}
	return strconv.Itoa(cheapest), nil
}

func day21p02(r io.Reader) (string, error) {
	_, dearest, err := equipmentCosts(r)
	if err != nil {
		return "", err
	}
	if dearest < 0 {
		return "", fmt.Errorf("no losing equipment")
	}
	return strconv.Itoa(dearest), nil
}
