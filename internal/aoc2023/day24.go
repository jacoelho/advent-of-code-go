package aoc2023

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/matrix"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type hailVector [3]matrix.Rat

func (v hailVector) sub(other hailVector) (result hailVector) {
	for i := range result {
		result[i] = v[i].Sub(other[i])
	}
	return result
}

func (v hailVector) cross(other hailVector) (result hailVector) {
	for i := range result {
		j, k := (i+1)%3, (i+2)%3
		result[i] = v[j].Mul(other[k]).Sub(v[k].Mul(other[j]))
	}
	return result
}

type hailstone struct {
	position, velocity hailVector
}

func parseHailVector(data []byte) (hailVector, error) {
	var result hailVector
	parts := bytes.Split(data, []byte(","))
	if len(parts) != len(result) {
		return result, fmt.Errorf("hail vector needs three coordinates")
	}
	for i, part := range parts {
		value, err := convert.ScanNumber[int64](bytes.TrimSpace(part))
		if err != nil {
			return result, fmt.Errorf("hail coordinate %d: %w", i+1, err)
		}
		result[i] = matrix.NewRat(value, 1)
	}
	return result, nil
}

func parseHailstone(line []byte) (hailstone, error) {
	position, velocity, ok := bytes.Cut(line, []byte("@"))
	if !ok {
		return hailstone{}, fmt.Errorf("hailstone needs position @ velocity")
	}
	p, err := parseHailVector(position)
	if err != nil {
		return hailstone{}, err
	}
	v, err := parseHailVector(velocity)
	if err != nil {
		return hailstone{}, err
	}
	return hailstone{position: p, velocity: v}, nil
}

func parseHailstones(r io.Reader) ([]hailstone, error) {
	const inputLimit = 4 << 20
	limited := &io.LimitedReader{R: r, N: inputLimit + 1}
	s := scanner.NewScanner(limited, parseHailstone)
	var stones []hailstone
	for stone := range s.Values() {
		if len(stones) == 2000 {
			return nil, fmt.Errorf("hail input exceeds 2000 stones")
		}
		stones = append(stones, stone)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if limited.N == 0 {
		return nil, fmt.Errorf("hail input exceeds 4 MiB")
	}
	if len(stones) < 2 {
		return nil, fmt.Errorf("need at least two hailstones")
	}
	return stones, nil
}

func hailXYCross(a, b hailVector) matrix.Rat {
	return a[0].Mul(b[1]).Sub(a[1].Mul(b[0]))
}

func hailHasFutureTime(delta, velocity hailVector, axes int) bool {
	var time matrix.Rat
	known := false
	for i := range axes {
		if velocity[i].IsZero() {
			if !delta[i].IsZero() {
				return false
			}
			continue
		}
		candidate := delta[i].Div(velocity[i])
		if candidate.Sign() < 0 || known && candidate.Cmp(time) != 0 {
			return false
		}
		time, known = candidate, true
	}
	return true
}

func hailPointInBounds(point hailVector, lower, upper matrix.Rat) bool {
	return point[0].Cmp(lower) >= 0 && point[0].Cmp(upper) <= 0 &&
		point[1].Cmp(lower) >= 0 && point[1].Cmp(upper) <= 0
}

func parallelHailPathsIntersect(a, b hailstone, lower, upper matrix.Rat) bool {
	if a.velocity[0].IsZero() && a.velocity[1].IsZero() {
		return hailPointInBounds(a.position, lower, upper) &&
			hailHasFutureTime(a.position.sub(b.position), b.velocity, 2)
	}
	if b.velocity[0].IsZero() && b.velocity[1].IsZero() {
		return hailPointInBounds(b.position, lower, upper) &&
			hailHasFutureTime(b.position.sub(a.position), a.velocity, 2)
	}
	delta := b.position.sub(a.position)
	if !hailXYCross(delta, a.velocity).IsZero() {
		return false
	}

	axis := 0
	if a.velocity[axis].IsZero() {
		axis = 1
	}
	start := lower.Sub(a.position[axis]).Div(a.velocity[axis])
	end := upper.Sub(a.position[axis]).Div(a.velocity[axis])
	if start.Cmp(end) > 0 {
		start, end = end, start
	}
	if start.Sign() < 0 {
		start = matrix.Rat{}
	}

	other := 1 - axis
	if a.velocity[other].IsZero() {
		if a.position[other].Cmp(lower) < 0 || a.position[other].Cmp(upper) > 0 {
			return false
		}
	} else {
		first := lower.Sub(a.position[other]).Div(a.velocity[other])
		last := upper.Sub(a.position[other]).Div(a.velocity[other])
		if first.Cmp(last) > 0 {
			first, last = last, first
		}
		if first.Cmp(start) > 0 {
			start = first
		}
		if last.Cmp(end) < 0 {
			end = last
		}
	}

	// Express the second future ray as a lower or upper bound on the first ray's time.
	origin := delta[axis].Div(a.velocity[axis])
	if a.velocity[axis].Sign() == b.velocity[axis].Sign() {
		if origin.Cmp(start) > 0 {
			start = origin
		}
	} else if origin.Cmp(end) < 0 {
		end = origin
	}
	return start.Cmp(end) <= 0
}

func hailPathsIntersect(a, b hailstone, lower, upper matrix.Rat) bool {
	determinant := hailXYCross(a.velocity, b.velocity)
	if determinant.IsZero() {
		return parallelHailPathsIntersect(a, b, lower, upper)
	}
	delta := b.position.sub(a.position)
	timeA := hailXYCross(delta, b.velocity).Div(determinant)
	timeB := hailXYCross(delta, a.velocity).Div(determinant)
	if timeA.Sign() < 0 || timeB.Sign() < 0 {
		return false
	}
	point := a.position
	for i := range 2 {
		point[i] = point[i].Add(a.velocity[i].Mul(timeA))
	}
	return hailPointInBounds(point, lower, upper)
}

func countHailIntersections(stones []hailstone, minBound, maxBound int64) int {
	lower, upper := matrix.NewRat(minBound, 1), matrix.NewRat(maxBound, 1)
	count := 0
	for i, a := range stones {
		for _, b := range stones[i+1:] {
			if hailPathsIntersect(a, b, lower, upper) {
				count++
			}
		}
	}
	return count
}

func day24p01(r io.Reader) (string, error) {
	stones, err := parseHailstones(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(countHailIntersections(stones, 200000000000000, 400000000000000)), nil
}

func hailRockEquations(stones []hailstone) *matrix.Matrix {
	result := matrix.New(3*(len(stones)-1), 7)
	anchor := stones[0]
	anchorCross := anchor.position.cross(anchor.velocity)
	for i, stone := range stones[1:] {
		velocity := stone.velocity.sub(anchor.velocity)
		position := stone.position.sub(anchor.position)
		rhs := stone.position.cross(stone.velocity).sub(anchorCross)
		// r×(vi−v0) + (pi−p0)×u = pi×vi − p0×v0 cancels the common r×u term.
		rows := [3][7]matrix.Rat{
			{{}, velocity[2], velocity[1].Neg(), {}, position[2].Neg(), position[1], rhs[0]},
			{velocity[2].Neg(), {}, velocity[0], position[2], {}, position[0].Neg(), rhs[1]},
			{velocity[1], velocity[0].Neg(), {}, position[1].Neg(), position[0], {}, rhs[2]},
		}
		for j, row := range rows {
			for col, value := range row {
				result.Set(3*i+j, col, value)
			}
		}
	}
	return result
}

func validHailRock(values []matrix.Rat, stones []hailstone) bool {
	for _, value := range values {
		if !value.IsInt() {
			return false
		}
	}
	position := hailVector{values[0], values[1], values[2]}
	velocity := hailVector{values[3], values[4], values[5]}
	for _, stone := range stones {
		if !hailHasFutureTime(stone.position.sub(position), velocity.sub(stone.velocity), 3) {
			return false
		}
	}
	return true
}

func hailIntegerRoots(coefficients [3]matrix.Rat) []*big.Int {
	scale := big.NewInt(1)
	for _, coefficient := range coefficients {
		scale.Mul(scale, coefficient.Denominator())
	}
	var polynomial [3]*big.Int
	for i, coefficient := range coefficients {
		factor := new(big.Int).Quo(scale, coefficient.Denominator())
		polynomial[i] = factor.Mul(factor, coefficient.Numerator())
	}
	a, b, c := polynomial[0], polynomial[1], polynomial[2]
	if a.Sign() == 0 && b.Sign() == 0 {
		return nil
	}
	var roots []*big.Int
	addRoot := func(numerator, denominator *big.Int) {
		root, remainder := new(big.Int), new(big.Int)
		root.QuoRem(numerator, denominator, remainder)
		if remainder.Sign() != 0 || len(roots) != 0 && roots[0].Cmp(root) == 0 {
			return
		}
		roots = append(roots, root)
	}
	if a.Sign() == 0 {
		addRoot(new(big.Int).Neg(c), b)
		return roots
	}
	discriminant := new(big.Int).Mul(b, b)
	product := new(big.Int).Mul(a, c)
	product.Mul(product, big.NewInt(4))
	discriminant.Sub(discriminant, product)
	if discriminant.Sign() < 0 {
		return nil
	}
	squareRoot := new(big.Int).Sqrt(discriminant)
	if new(big.Int).Mul(squareRoot, squareRoot).Cmp(discriminant) != 0 {
		return nil
	}
	denominator := new(big.Int).Mul(big.NewInt(2), a)
	negativeB := new(big.Int).Neg(b)
	addRoot(new(big.Int).Add(negativeB, squareRoot), denominator)
	addRoot(new(big.Int).Sub(negativeB, squareRoot), denominator)
	return roots
}

func hailRockParameters(solution matrix.ParametricSolution, anchor hailstone) ([]*big.Int, error) {
	base := solution.Base
	coefficients := solution.Coeffs[0]
	position := hailVector{base[0], base[1], base[2]}.sub(anchor.position)
	velocity := hailVector{base[3], base[4], base[5]}.sub(anchor.velocity)
	positionSlope := hailVector{coefficients[0], coefficients[1], coefficients[2]}
	velocitySlope := hailVector{coefficients[3], coefficients[4], coefficients[5]}
	quadratic := positionSlope.cross(velocitySlope)
	linear := positionSlope.cross(velocity)
	otherLinear := position.cross(velocitySlope)
	constant := position.cross(velocity)
	for axis := range 3 {
		linear[axis] = linear[axis].Add(otherLinear[axis])
		if quadratic[axis].IsZero() && linear[axis].IsZero() && constant[axis].IsZero() {
			continue
		}
		// The free variable is a rock component, so integer roots are necessary.
		return hailIntegerRoots([3]matrix.Rat{quadratic[axis], linear[axis], constant[axis]}), nil
	}
	return nil, fmt.Errorf("unsupported rock family: collision polynomial vanishes identically")
}

func findHailRock(stones []hailstone) ([]matrix.Rat, error) {
	if len(stones) < 3 {
		return nil, fmt.Errorf("need at least three hailstones to determine a rock")
	}
	equations := hailRockEquations(stones)
	reduced := matrix.RREF(equations)
	if reduced.Inconsistent {
		return nil, fmt.Errorf("hailstone collision equations are inconsistent")
	}
	solution := matrix.ExtractParametricSolution(equations, reduced.PivotCols, 6)
	if len(solution.FreeCols) == 0 {
		if !validHailRock(solution.Base, stones) {
			return nil, fmt.Errorf("rock candidate does not collide with every hailstone at a nonnegative time with integer position and velocity")
		}
		return solution.Base, nil
	}
	if len(solution.FreeCols) != 1 {
		return nil, fmt.Errorf("unsupported rock family with %d free components", len(solution.FreeCols))
	}
	parameters, err := hailRockParameters(solution, stones[0])
	if err != nil {
		return nil, err
	}
	var result []matrix.Rat
	for _, parameter := range parameters {
		value := matrix.NewInteger(parameter)
		candidate := make([]matrix.Rat, 6)
		for i := range candidate {
			candidate[i] = solution.Base[i].Add(solution.Coeffs[0][i].Mul(value))
		}
		if validHailRock(candidate, stones) {
			if result != nil {
				return nil, fmt.Errorf("multiple integer rocks collide with every hailstone in the future")
			}
			result = candidate
		}
	}
	if result == nil {
		return nil, fmt.Errorf("no integer rock collides with every hailstone in the future")
	}
	return result, nil
}

func day24p02(r io.Reader) (string, error) {
	stones, err := parseHailstones(r)
	if err != nil {
		return "", err
	}
	rock, err := findHailRock(stones)
	if err != nil {
		return "", err
	}
	return rock[0].Add(rock[1]).Add(rock[2]).Numerator().String(), nil
}
