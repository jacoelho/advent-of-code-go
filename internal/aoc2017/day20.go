package aoc2017

import (
	"fmt"
	"io"
	"math/big"
	"regexp"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type particle struct {
	position, velocity, acceleration [3]int64
}

func parseParticles(r io.Reader) ([]particle, error) {
	pattern := regexp.MustCompile(`^p=<\s*([+-]?\d+),\s*([+-]?\d+),\s*([+-]?\d+)>,\s*v=<\s*([+-]?\d+),\s*([+-]?\d+),\s*([+-]?\d+)>,\s*a=<\s*([+-]?\d+),\s*([+-]?\d+),\s*([+-]?\d+)>$`)
	s := scanner.NewScanner(r, func(line []byte) (particle, error) {
		fields := pattern.FindSubmatch(line)
		if fields == nil {
			return particle{}, fmt.Errorf("invalid particle %q", line)
		}
		var result particle
		for group, vector := range []*[3]int64{&result.position, &result.velocity, &result.acceleration} {
			for axis := range 3 {
				value, err := convert.ScanNumber[int64](fields[1+group*3+axis])
				if err != nil {
					return particle{}, err
				}
				vector[axis] = value
			}
		}
		return result, nil
	})
	particles := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(particles) == 0 {
		return nil, fmt.Errorf("empty particle swarm")
	}

	return particles, nil
}

// particleTrajectory holds, per axis, the coefficients a, b, c of twice the coordinate at tick t: a*t*t + b*t + c.
type particleTrajectory [3][3]*big.Int

func (p particle) trajectory() particleTrajectory {
	var trajectory particleTrajectory
	for axis := range 3 {
		a := big.NewInt(p.acceleration[axis])
		b := new(big.Int).Lsh(big.NewInt(p.velocity[axis]), 1)
		b.Add(b, a)
		c := new(big.Int).Lsh(big.NewInt(p.position[axis]), 1)
		trajectory[axis] = [3]*big.Int{a, b, c}
	}
	return trajectory
}

func eventualDistance(p particle) [3]*big.Int {
	distance := [3]*big.Int{new(big.Int), new(big.Int), new(big.Int)}
	for _, coefficients := range p.trajectory() {
		sign := 0
		for _, coefficient := range coefficients {
			if sign = coefficient.Sign(); sign != 0 {
				break
			}
		}
		for i, coefficient := range coefficients {
			if sign < 0 {
				distance[i].Sub(distance[i], coefficient)
			} else {
				distance[i].Add(distance[i], coefficient)
			}
		}
	}

	return distance
}

func day20p01(r io.Reader) (string, error) {
	particles, err := parseParticles(r)
	if err != nil {
		return "", err
	}

	closest := 0
	best := eventualDistance(particles[0])
	for id, p := range particles[1:] {
		distance := eventualDistance(p)
		for i := range 3 {
			comparison := distance[i].Cmp(best[i])
			if comparison < 0 {
				closest, best = id+1, distance
			}
			if comparison != 0 {
				break
			}
		}
	}

	return strconv.Itoa(closest), nil
}

// collisionSolver finds the earliest tick at which two trajectories meet, reusing its scratch values across pairs.
type collisionSolver struct {
	difference   [3][3]big.Int
	discriminant big.Int
	root         big.Int
	numerator    big.Int
	denominator  big.Int
	remainder    big.Int
	value        big.Int
	roots        [2]big.Int
}

// addRoot stores numerator/denominator as roots[count] when it is a nonnegative integer and returns the root count.
func (s *collisionSolver) addRoot(count int) int {
	root := &s.roots[count]
	root.QuoRem(&s.numerator, &s.denominator, &s.remainder)
	if s.remainder.Sign() == 0 && root.Sign() >= 0 {
		return count + 1
	}
	return count
}

// integerRoots stores the nonnegative integer roots of a*t*t + b*t + c and returns how many there are.
func (s *collisionSolver) integerRoots(coefficients *[3]big.Int) int {
	a, b, c := &coefficients[0], &coefficients[1], &coefficients[2]
	if a.Sign() == 0 {
		if b.Sign() == 0 {
			return 0
		}
		s.numerator.Neg(c)
		s.denominator.Set(b)
		return s.addRoot(0)
	}

	s.discriminant.Mul(b, b)
	s.value.Mul(a, c)
	s.discriminant.Sub(&s.discriminant, s.value.Lsh(&s.value, 2))
	if s.discriminant.Sign() < 0 {
		return 0
	}
	s.root.Sqrt(&s.discriminant)
	if s.value.Mul(&s.root, &s.root).Cmp(&s.discriminant) != 0 {
		return 0
	}
	s.denominator.Lsh(a, 1)
	s.numerator.Sub(&s.root, b)
	count := s.addRoot(0)
	if s.root.Sign() != 0 {
		s.numerator.Neg(b)
		s.numerator.Sub(&s.numerator, &s.root)
		count = s.addRoot(count)
	}
	return count
}

// earliest returns the first tick at which both trajectories coincide, or nil if they never do.
func (s *collisionSolver) earliest(first, second *particleTrajectory) *big.Int {
	constrained := -1
	for axis := range 3 {
		for i := range 3 {
			if s.difference[axis][i].Sub(first[axis][i], second[axis][i]).Sign() != 0 {
				constrained = axis
			}
		}
	}
	if constrained < 0 {
		return new(big.Int)
	}

	var earliest *big.Int
	for i := range s.integerRoots(&s.difference[constrained]) {
		time := &s.roots[i]
		matches := true
		for axis := range s.difference {
			coefficients := &s.difference[axis]
			s.value.Mul(&coefficients[0], time)
			s.value.Add(&s.value, &coefficients[1]).Mul(&s.value, time).Add(&s.value, &coefficients[2])
			if s.value.Sign() != 0 {
				matches = false
				break
			}
		}
		if matches && (earliest == nil || time.Cmp(earliest) < 0) {
			earliest = time
		}
	}
	if earliest == nil {
		return nil
	}

	return new(big.Int).Set(earliest)
}

func day20p02(r io.Reader) (string, error) {
	particles, err := parseParticles(r)
	if err != nil {
		return "", err
	}

	trajectories := make([]particleTrajectory, len(particles))
	for i, p := range particles {
		trajectories[i] = p.trajectory()
	}

	type collision struct {
		time          *big.Int
		first, second int
	}
	var solver collisionSolver
	var events []collision
	for i := range trajectories {
		for j := i + 1; j < len(trajectories); j++ {
			// After a pair's first collision, at least one member is gone forever.
			if time := solver.earliest(&trajectories[i], &trajectories[j]); time != nil {
				events = append(events, collision{time: time, first: i, second: j})
			}
		}
	}
	slices.SortFunc(events, func(a, b collision) int { return a.time.Cmp(b.time) })

	removed := make([]bool, len(particles))
	remaining := len(particles)
	for start := 0; start < len(events); {
		end := start + 1
		for end < len(events) && events[end].time.Cmp(events[start].time) == 0 {
			end++
		}
		var colliding []int
		for _, event := range events[start:end] {
			if !removed[event.first] && !removed[event.second] {
				colliding = append(colliding, event.first, event.second)
			}
		}
		for _, id := range colliding {
			if !removed[id] {
				removed[id] = true
				remaining--
			}
		}
		start = end
	}

	return strconv.Itoa(remaining), nil
}
