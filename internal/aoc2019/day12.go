package aoc2019

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xmath"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

type Moon struct {
	x, y, z    int
	vx, vy, vz int
}

func parseMoon(line []byte) (*Moon, error) {
	coords := convert.ExtractDigits[int](string(line))
	if len(coords) != 3 {
		return nil, fmt.Errorf("expected 3 coordinates, got %d", len(coords))
	}
	return &Moon{x: coords[0], y: coords[1], z: coords[2]}, nil
}

func parseMoons(r io.Reader) ([]*Moon, error) {
	s := scanner.NewScanner(r, parseMoon)
	return slices.Collect(s.Values()), s.Err()
}

func applyGravityAxis(pos1, pos2 int) (delta1, delta2 int) {
	switch {
	case pos1 < pos2:
		return 1, -1
	case pos1 > pos2:
		return -1, 1
	default:
		return 0, 0
	}
}

func applyGravity(moons []*Moon) {
	for pair := range xslices.Pairwise(moons) {
		m1, m2 := pair.V1, pair.V2

		dx1, dx2 := applyGravityAxis(m1.x, m2.x)
		m1.vx += dx1
		m2.vx += dx2

		dy1, dy2 := applyGravityAxis(m1.y, m2.y)
		m1.vy += dy1
		m2.vy += dy2

		dz1, dz2 := applyGravityAxis(m1.z, m2.z)
		m1.vz += dz1
		m2.vz += dz2
	}
}

func applyVelocity(moons []*Moon) {
	for _, m := range moons {
		m.x += m.vx
		m.y += m.vy
		m.z += m.vz
	}
}

func simulate(moons []*Moon, steps int) {
	for range steps {
		applyGravity(moons)
		applyVelocity(moons)
	}
}

func totalEnergy(moons []*Moon) int {
	return xslices.Reduce(func(total int, m *Moon) int {
		potentialEnergy := xmath.Abs(m.x) + xmath.Abs(m.y) + xmath.Abs(m.z)
		kineticEnergy := xmath.Abs(m.vx) + xmath.Abs(m.vy) + xmath.Abs(m.vz)
		return total + potentialEnergy*kineticEnergy
	}, 0, moons)
}

func day12p01(r io.Reader) (string, error) {
	moons, err := parseMoons(r)
	if err != nil {
		return "", err
	}

	simulate(moons, 1000)
	energy := totalEnergy(moons)

	return strconv.Itoa(energy), nil
}

type axisState struct {
	positions  []int
	velocities []int
}

func (a axisState) equals(other axisState) bool {
	return slices.Equal(a.positions, other.positions) && slices.Equal(a.velocities, other.velocities)
}

func getAxisState(moons []*Moon, axis int) axisState {
	n := len(moons)
	state := axisState{
		positions:  make([]int, n),
		velocities: make([]int, n),
	}

	for i, m := range moons {
		switch axis {
		case 0: // x
			state.positions[i] = m.x
			state.velocities[i] = m.vx
		case 1: // y
			state.positions[i] = m.y
			state.velocities[i] = m.vy
		case 2: // z
			state.positions[i] = m.z
			state.velocities[i] = m.vz
		}
	}

	return state
}

func applyGravityAxis1D(state *axisState) {
	n := len(state.positions)
	for i := range n {
		for j := i + 1; j < n; j++ {
			dv1, dv2 := applyGravityAxis(state.positions[i], state.positions[j])
			state.velocities[i] += dv1
			state.velocities[j] += dv2
		}
	}
}

func applyVelocityAxis1D(state *axisState) {
	for i := range state.positions {
		state.positions[i] += state.velocities[i]
	}
}

func findAxisCycle(moons []*Moon, axis int) int {
	initial := getAxisState(moons, axis)
	current := axisState{
		positions:  slices.Clone(initial.positions),
		velocities: slices.Clone(initial.velocities),
	}

	steps := 0
	for {
		steps++
		applyGravityAxis1D(&current)
		applyVelocityAxis1D(&current)
		if current.equals(initial) {
			return steps
		}
	}
}

func day12p02(r io.Reader) (string, error) {
	moons, err := parseMoons(r)
	if err != nil {
		return "", err
	}

	cycleX := findAxisCycle(moons, 0)
	cycleY := findAxisCycle(moons, 1)
	cycleZ := findAxisCycle(moons, 2)

	result := xmath.LCM(cycleX, cycleY, cycleZ)

	return strconv.Itoa(result), nil
}
