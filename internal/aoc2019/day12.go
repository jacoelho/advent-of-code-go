package aoc2019

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

// moonAxis holds every moon's position and velocity along one axis. Gravity
// on an axis depends only on positions along that axis, so axes evolve
// independently.
type moonAxis struct {
	positions  []int
	velocities []int
}

func parseMoon(line []byte) ([3]int, error) {
	var moon [3]int
	text, hasOpen := strings.CutPrefix(string(line), "<")
	text, hasClose := strings.CutSuffix(text, ">")
	if !hasOpen || !hasClose {
		return moon, fmt.Errorf("invalid moon %q", line)
	}
	fields := strings.Split(text, ", ")
	if len(fields) != 3 {
		return moon, fmt.Errorf("invalid moon %q", line)
	}
	for axis, name := range []string{"x=", "y=", "z="} {
		value, ok := strings.CutPrefix(fields[axis], name)
		if !ok {
			return moon, fmt.Errorf("invalid moon %q", line)
		}
		position, err := convert.ScanNumber[int]([]byte(value))
		if err != nil {
			return moon, fmt.Errorf("moon %q: %w", line, err)
		}
		moon[axis] = position
	}
	return moon, nil
}

func parseMoons(r io.Reader) ([3]moonAxis, error) {
	var axes [3]moonAxis
	s := scanner.NewScanner(r, parseMoon)
	for moon := range s.Values() {
		for axis, position := range moon {
			axes[axis].positions = append(axes[axis].positions, position)
			axes[axis].velocities = append(axes[axis].velocities, 0)
		}
	}
	if err := s.Err(); err != nil {
		return axes, err
	}
	if len(axes[0].positions) == 0 {
		return axes, fmt.Errorf("no moons")
	}
	return axes, nil
}

func (a *moonAxis) step() {
	for i, pi := range a.positions {
		for j := i + 1; j < len(a.positions); j++ {
			pull := 0
			switch pj := a.positions[j]; {
			case pi < pj:
				pull = 1
			case pi > pj:
				pull = -1
			}
			a.velocities[i] += pull
			a.velocities[j] -= pull
		}
	}
	for i, v := range a.velocities {
		a.positions[i] += v
	}
}

// period returns the number of steps until the axis first returns to its
// starting state. The simulation is reversible, so the starting state is the
// first to repeat.
func (a moonAxis) period() int {
	current := moonAxis{positions: slices.Clone(a.positions), velocities: slices.Clone(a.velocities)}
	for steps := 1; ; steps++ {
		current.step()
		if slices.Equal(current.positions, a.positions) && slices.Equal(current.velocities, a.velocities) {
			return steps
		}
	}
}

func day12p01(r io.Reader) (string, error) {
	axes, err := parseMoons(r)
	if err != nil {
		return "", err
	}
	for axis := range axes {
		for range 1000 {
			axes[axis].step()
		}
	}

	energy := 0
	for moon := range axes[0].positions {
		potential, kinetic := 0, 0
		for _, a := range axes {
			potential += xmath.Abs(a.positions[moon])
			kinetic += xmath.Abs(a.velocities[moon])
		}
		energy += potential * kinetic
	}
	return strconv.Itoa(energy), nil
}

func day12p02(r io.Reader) (string, error) {
	axes, err := parseMoons(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(xmath.LCM(axes[0].period(), axes[1].period(), axes[2].period())), nil
}
