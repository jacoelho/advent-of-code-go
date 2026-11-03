package aoc2023

import (
	"math/big"
	"testing"
)

func TestCombineCongruences(t *testing.T) {
	for _, tt := range []struct {
		a, m, b, n, residue, modulus int64
	}{
		{2, 6, 5, 9, 14, 18},
		{5, 6, 2, 3, 5, 6},
	} {
		a, m, b, n := big.NewInt(tt.a), big.NewInt(tt.m), big.NewInt(tt.b), big.NewInt(tt.n)
		residue, modulus, compatible := combineCongruences(a, m, b, n)
		if !compatible || residue.Cmp(big.NewInt(tt.residue)) != 0 || modulus.Cmp(big.NewInt(tt.modulus)) != 0 {
			t.Fatalf("merge = %v modulo %v, %v; want %d modulo %d", residue, modulus, compatible, tt.residue, tt.modulus)
		}
		residue.SetInt64(99)
		modulus.SetInt64(100)
		for i, value := range []*big.Int{a, m, b, n} {
			want := [4]int64{tt.a, tt.m, tt.b, tt.n}[i]
			if value.Cmp(big.NewInt(want)) != 0 {
				t.Fatalf("input %d changed to %s; want %d", i, value, want)
			}
		}
	}
	if _, _, compatible := combineCongruences(big.NewInt(2), big.NewInt(6), big.NewInt(4), big.NewInt(9)); compatible {
		t.Fatal("incompatible congruences accepted")
	}
}

func TestFirstCongruentAtLeast(t *testing.T) {
	for _, tt := range []struct{ minimum, want int64 }{{2, 2}, {3, 5}, {5, 5}} {
		got := firstCongruentAtLeast(big.NewInt(2), big.NewInt(3), big.NewInt(tt.minimum))
		if got.Cmp(big.NewInt(tt.want)) != 0 {
			t.Fatalf("minimum %d = %s; want %d", tt.minimum, got, tt.want)
		}
	}
}
