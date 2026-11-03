package aoc2023

import "math/big"

// combineCongruences solves x ≡ a (mod m) and x ≡ b (mod n).
// It requires positive moduli and does not mutate its inputs. Incompatible
// congruences return nil, nil, false; otherwise residue is normalized modulo
// modulus, the LCM of m and n.
func combineCongruences(a, m, b, n *big.Int) (residue, modulus *big.Int, compatible bool) {
	var inverse, remainder big.Int
	divisor := new(big.Int).GCD(&inverse, nil, m, n)
	difference := new(big.Int).Sub(b, a)
	multiple := new(big.Int)
	multiple.QuoRem(difference, divisor, &remainder)
	if remainder.Sign() != 0 {
		return nil, nil, false
	}
	factor := new(big.Int).Quo(n, divisor)
	multiple.Mul(multiple, &inverse).Mod(multiple, factor)
	modulus = new(big.Int).Mul(m, factor)
	residue = new(big.Int).Mul(m, multiple)
	residue.Add(residue, a).Mod(residue, modulus)
	return residue, modulus, true
}

func firstCongruentAtLeast(residue, modulus, minimum *big.Int) *big.Int {
	difference := new(big.Int).Sub(minimum, residue)
	var remainder big.Int
	cycles := new(big.Int)
	cycles.QuoRem(difference, modulus, &remainder)
	if remainder.Sign() > 0 {
		cycles.Add(cycles, big.NewInt(1))
	}
	return cycles.Mul(cycles, modulus).Add(cycles, residue)
}
