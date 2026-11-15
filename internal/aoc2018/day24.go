package aoc2018

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type armySide uint8

const (
	noArmy armySide = iota
	immuneArmy
	infectionArmy
)

type armyGroup struct {
	side                     armySide
	units, hitPoints, damage int64
	damageType               string
	initiative               int
	weaknesses, immunities   collections.Set[string]
}

func parseArmies(r io.Reader) ([]armyGroup, error) {
	pattern := regexp.MustCompile(`^(\d+) units each with (\d+) hit points(?: \(([^)]+)\))? with an attack that does (\d+) ([a-z]+) damage at initiative (\d+)$`)
	var groups []armyGroup
	var side armySide
	var headings [3]bool
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		if line == "Immune System:" || line == "Infection:" {
			side = immuneArmy
			if line == "Infection:" {
				side = infectionArmy
			}
			if headings[side] {
				return nil, fmt.Errorf("duplicate army heading %q", line)
			}
			headings[side] = true
			continue
		}
		fields := pattern.FindStringSubmatch(line)
		if fields == nil || side == noArmy {
			return nil, fmt.Errorf("invalid army group %q", line)
		}
		// int32 values keep every product of two values, and the boost bound, within int64.
		var values [4]int64
		for i, field := range []string{fields[1], fields[2], fields[4], fields[6]} {
			value, err := convert.ScanNumber[int32]([]byte(field))
			if err != nil {
				return nil, err
			}
			if value < 0 || value == 0 && i != 2 {
				return nil, fmt.Errorf("army group %q: values must be positive, except damage may be zero", line)
			}
			values[i] = int64(value)
		}
		group := armyGroup{
			side: side, units: values[0], hitPoints: values[1], damage: values[2], initiative: int(values[3]),
			damageType: fields[5], weaknesses: collections.NewSet[string](), immunities: collections.NewSet[string](),
		}
		if fields[3] != "" {
			for _, modifier := range strings.Split(fields[3], "; ") {
				kind, types, ok := strings.Cut(modifier, " to ")
				if !ok || kind != "weak" && kind != "immune" {
					return nil, fmt.Errorf("invalid army modifier %q", modifier)
				}
				for _, damageType := range strings.Split(types, ", ") {
					if damageType == "" || strings.IndexFunc(damageType, func(r rune) bool { return r < 'a' || r > 'z' }) >= 0 {
						return nil, fmt.Errorf("invalid damage type %q", damageType)
					}
					if kind == "weak" {
						group.weaknesses.Add(damageType)
					} else {
						group.immunities.Add(damageType)
					}
				}
			}
		}
		groups = append(groups, group)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	var counts [3]int
	for _, group := range groups {
		counts[group.side]++
	}
	if counts[immuneArmy] == 0 || counts[infectionArmy] == 0 {
		return nil, fmt.Errorf("expected at least one group in each army")
	}

	return groups, nil
}

func (group armyGroup) power() int64 {
	return group.units * group.damage
}

func (group armyGroup) damageTo(target armyGroup) int64 {
	if target.immunities.Contains(group.damageType) {
		return 0
	}
	if target.weaknesses.Contains(group.damageType) {
		return 2 * group.power()
	}
	return group.power()
}

func fightArmies(initial []armyGroup, boost int64) (armySide, int64) {
	groups := slices.Clone(initial)
	order := make([]int, len(groups))
	for i := range groups {
		order[i] = i
		if groups[i].side == immuneArmy {
			groups[i].damage += boost
		}
	}
	attackOrder := slices.Clone(order)
	initiativeOrder := func(a, b int) int {
		if groups[a].initiative != groups[b].initiative {
			return cmp.Compare(groups[b].initiative, groups[a].initiative)
		}
		return cmp.Compare(a, b)
	}
	slices.SortFunc(attackOrder, initiativeOrder)
	targets := make([]int, len(groups))
	selected := make([]bool, len(groups))

	for {
		var units [3]int64
		for _, group := range groups {
			units[group.side] += group.units
		}
		if units[immuneArmy] == 0 {
			return infectionArmy, units[infectionArmy]
		}
		if units[infectionArmy] == 0 {
			return immuneArmy, units[immuneArmy]
		}

		slices.SortFunc(order, func(a, b int) int {
			return cmp.Or(cmp.Compare(groups[b].power(), groups[a].power()), initiativeOrder(a, b))
		})
		clear(selected)
		for i := range targets {
			targets[i] = -1
		}
		for _, id := range order {
			attacker := groups[id]
			if attacker.units == 0 {
				continue
			}
			best := -1
			var bestDamage int64
			for targetID, target := range groups {
				if target.units == 0 || target.side == attacker.side || selected[targetID] {
					continue
				}
				damage := attacker.damageTo(target)
				if damage == 0 {
					continue
				}
				if best == -1 || damage > bestDamage || damage == bestDamage && (target.power() > groups[best].power() || target.power() == groups[best].power() && target.initiative > groups[best].initiative) {
					best, bestDamage = targetID, damage
				}
			}
			if best >= 0 {
				targets[id] = best
				selected[best] = true
			}
		}

		var killed int64
		for _, id := range attackOrder {
			attacker := groups[id]
			targetID := targets[id]
			if attacker.units == 0 || targetID < 0 || groups[targetID].units == 0 {
				continue
			}
			target := &groups[targetID]
			losses := min(target.units, attacker.damageTo(*target)/target.hitPoints)
			target.units -= losses
			killed += losses
		}
		if killed == 0 {
			return noArmy, units[immuneArmy] + units[infectionArmy]
		}
	}
}

func day24p01(r io.Reader) (string, error) {
	groups, err := parseArmies(r)
	if err != nil {
		return "", err
	}
	winner, units := fightArmies(groups, 0)
	if winner == noArmy {
		return "", fmt.Errorf("army battle is a stalemate")
	}

	return strconv.FormatInt(units, 10), nil
}

func day24p02(r io.Reader) (string, error) {
	groups, err := parseArmies(r)
	if err != nil {
		return "", err
	}
	// From this boost on, every landed immune attack destroys its target and the
	// ranking of immune groups by effective power follows units then damage, so
	// larger boosts replay the same battle.
	var finalBoost int64
	for _, group := range groups {
		if group.side == immuneArmy {
			finalBoost = max(finalBoost, group.units*group.damage)
		} else {
			finalBoost = max(finalBoost, group.units*group.hitPoints)
		}
	}
	for boost := range finalBoost + 1 {
		if winner, units := fightArmies(groups, boost); winner == immuneArmy {
			return strconv.FormatInt(units, 10), nil
		}
	}

	return "", fmt.Errorf("no boost lets the immune system win")
}
