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
	var particles []particle
	for p := range s.Values() {
		if len(particles) == 2000 {
			return nil, fmt.Errorf("particle analysis supports at most 2000 particles")
		}
		particles = append(particles, p)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(particles) == 0 {
		return nil, fmt.Errorf("empty particle swarm")
	}

	return particles, nil
}

func particleAxis(p particle, axis int) [3]*big.Int {
	// Twice the coordinate at t is a*t*t + (a+2*v)*t + 2*p.
	a := big.NewInt(p.acceleration[axis])
	b := new(big.Int).Lsh(big.NewInt(p.velocity[axis]), 1)
	b.Add(b, a)
	c := new(big.Int).Lsh(big.NewInt(p.position[axis]), 1)
	return [3]*big.Int{a, b, c}
}

func eventualDistance(p particle) [3]*big.Int {
	distance := [3]*big.Int{new(big.Int), new(big.Int), new(big.Int)}
	for axis := range 3 {
		coefficients := particleAxis(p, axis)
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

func integerRoots(coefficients [3]*big.Int) []*big.Int {
	a, b, c := coefficients[0], coefficients[1], coefficients[2]
	var roots []*big.Int
	add := func(numerator, denominator *big.Int) {
		quotient, remainder := new(big.Int), new(big.Int)
		quotient.QuoRem(numerator, denominator, remainder)
		if remainder.Sign() == 0 && quotient.Sign() >= 0 {
			roots = append(roots, quotient)
		}
	}
	if a.Sign() == 0 {
		if b.Sign() != 0 {
			add(new(big.Int).Neg(c), b)
		}
		return roots
	}

	discriminant := new(big.Int).Mul(b, b)
	product := new(big.Int).Mul(a, c)
	discriminant.Sub(discriminant, product.Lsh(product, 2))
	if discriminant.Sign() < 0 {
		return nil
	}
	root := new(big.Int).Sqrt(discriminant)
	if new(big.Int).Mul(root, root).Cmp(discriminant) != 0 {
		return nil
	}
	denominator := new(big.Int).Lsh(a, 1)
	numerator := new(big.Int).Neg(b)
	add(new(big.Int).Add(numerator, root), denominator)
	if root.Sign() != 0 {
		add(new(big.Int).Sub(numerator, root), denominator)
	}

	return roots
}

func particleCollision(first, second particle) *big.Int {
	var differences [3][3]*big.Int
	constrained := -1
	for axis := range 3 {
		left, right := particleAxis(first, axis), particleAxis(second, axis)
		for i := range 3 {
			differences[axis][i] = left[i].Sub(left[i], right[i])
			if differences[axis][i].Sign() != 0 {
				constrained = axis
			}
		}
	}
	if constrained < 0 {
		return new(big.Int)
	}

	var earliest *big.Int
	for _, time := range integerRoots(differences[constrained]) {
		matches := true
		for _, coefficients := range differences {
			value := new(big.Int).Mul(coefficients[0], time)
			value.Add(value, coefficients[1]).Mul(value, time).Add(value, coefficients[2])
			if value.Sign() != 0 {
				matches = false
				break
			}
		}
		if matches && (earliest == nil || time.Cmp(earliest) < 0) {
			earliest = time
		}
	}

	return earliest
}

func day20p02(r io.Reader) (string, error) {
	particles, err := parseParticles(r)
	if err != nil {
		return "", err
	}

	type collision struct {
		time          *big.Int
		first, second int
	}
	var events []collision
	for i, first := range particles {
		for j := i + 1; j < len(particles); j++ {
			// After a pair's first collision, at least one member is gone forever.
			if time := particleCollision(first, particles[j]); time != nil {
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
