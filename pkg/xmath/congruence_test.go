package xmath

import (
	"errors"
	"math"
	"testing"
)

func TestCombineCongruencesMatchesExhaustiveSearch(t *testing.T) {
	for m := 1; m <= 12; m++ {
		for n := 1; n <= 12; n++ {
			lcm := m / GCD(m, n) * n
			for a := -m; a < m; a++ {
				for b := -n; b < n; b++ {
					want := -1
					for x := range lcm {
						if Modulo(x-a, m) == 0 && Modulo(x-b, n) == 0 {
							want = x
							break
						}
					}
					residue, modulus, err := CombineCongruences(a, m, b, n)
					if want < 0 {
						if !errors.Is(err, ErrNoSolution) {
							t.Fatalf("x≡%d (mod %d), x≡%d (mod %d): got %d, %v; want no solution", a, m, b, n, residue, err)
						}
						continue
					}
					if err != nil || residue != want || modulus != lcm {
						t.Fatalf("x≡%d (mod %d), x≡%d (mod %d): got %d mod %d, %v; want %d mod %d", a, m, b, n, residue, modulus, err, want, lcm)
					}
				}
			}
		}
	}
}

func TestCombineCongruencesNearTypeLimit(t *testing.T) {
	residue, modulus, err := CombineCongruences[int8](126, 127, 0, 1)
	if err != nil || residue != 126 || modulus != 127 {
		t.Fatalf("got %d mod %d, %v", residue, modulus, err)
	}
	if _, _, err := CombineCongruences[int8](0, 16, 0, 9); !errors.Is(err, ErrOverflow) {
		t.Fatalf("lcm 144 in int8: got %v; want overflow", err)
	}
	wide, wideModulus, err := CombineCongruences[int64](math.MaxInt64/2-1, math.MaxInt64/2, 1, 2)
	if err != nil || wideModulus != math.MaxInt64-1 || wide%(math.MaxInt64/2) != math.MaxInt64/2-1 || wide%2 != 1 {
		t.Fatalf("got %d mod %d, %v", wide, wideModulus, err)
	}
	if got, _, err := CombineCongruences[uint8](3, 5, 0, 7); err != nil || got != 28 {
		t.Fatalf("unsigned: got %d, %v; want 28", got, err)
	}
}
