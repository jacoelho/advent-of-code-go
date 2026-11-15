package xmath

import (
	"math"
	"math/big"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestCheckedArithmeticMatchesExactArithmetic(t *testing.T) {
	for a := math.MinInt8; a <= math.MaxInt8; a++ {
		for b := math.MinInt8; b <= math.MaxInt8; b++ {
			checkExact(t, "add", int8(a), int8(b), a+b, CheckedAdd[int8])
			checkExact(t, "sub", int8(a), int8(b), a-b, CheckedSub[int8])
			checkExact(t, "mul", int8(a), int8(b), a*b, CheckedMul[int8])
		}
	}
	for a := 0; a <= math.MaxUint8; a++ {
		for b := 0; b <= math.MaxUint8; b++ {
			checkExact(t, "add", uint8(a), uint8(b), a+b, CheckedAdd[uint8])
			checkExact(t, "sub", uint8(a), uint8(b), a-b, CheckedSub[uint8])
			checkExact(t, "mul", uint8(a), uint8(b), a*b, CheckedMul[uint8])
		}
	}
}

func checkExact[T constraints.Integer](t *testing.T, op string, a, b T, exact int, f func(T, T) (T, bool)) {
	t.Helper()
	got, ok := f(a, b)
	fits := exact == int(T(exact))
	if ok != fits || ok && int(got) != exact {
		t.Fatalf("%s(%d, %d) = %d, %v; exact %d", op, a, b, got, ok, exact)
	}
}

func TestCheckedArithmeticAtInt64Limits(t *testing.T) {
	if _, ok := CheckedMul(int64(math.MinInt64), -1); ok {
		t.Fatal("MinInt64 * -1 reported as fitting")
	}
	if got, ok := CheckedMul(int64(math.MinInt64), 1); !ok || got != math.MinInt64 {
		t.Fatalf("MinInt64 * 1 = %d, %v", got, ok)
	}
	if _, ok := CheckedAdd(int64(math.MaxInt64), 1); ok {
		t.Fatal("MaxInt64 + 1 reported as fitting")
	}
	if got, ok := CheckedSub(int64(-1), math.MaxInt64); !ok || got != math.MinInt64 {
		t.Fatalf("-1 - MaxInt64 = %d, %v", got, ok)
	}
}

func TestMulModMatchesExactArithmetic(t *testing.T) {
	for m := 1; m <= math.MaxInt8; m++ {
		for a := math.MinInt8; a <= math.MaxInt8; a++ {
			for b := math.MinInt8; b <= math.MaxInt8; b += 7 {
				want := ((a*b)%m + m) % m
				if got := MulMod(int8(a), int8(b), int8(m)); int(got) != want {
					t.Fatalf("MulMod(%d, %d, %d) = %d; want %d", a, b, m, got, want)
				}
			}
		}
	}
	a, b, m := uint64(math.MaxUint64-1), uint64(math.MaxUint64-2), uint64(math.MaxUint64)
	want := new(big.Int).Mod(new(big.Int).Mul(new(big.Int).SetUint64(a), new(big.Int).SetUint64(b)), new(big.Int).SetUint64(m))
	if got := MulMod(a, b, m); got != want.Uint64() {
		t.Fatalf("MulMod near MaxUint64 = %d; want %s", got, want)
	}
}

func TestModInverseMatchesExhaustiveSearch(t *testing.T) {
	for m := 1; m <= math.MaxInt8; m++ {
		for a := math.MinInt8; a <= math.MaxInt8; a++ {
			want, exists := -1, false
			for x := range m {
				if ((a*x)%m+m)%m == 1%m {
					want, exists = x, true
					break
				}
			}
			got, ok := ModInverse(int8(a), int8(m))
			if ok != exists || ok && int(got) != want {
				t.Fatalf("ModInverse(%d, %d) = %d, %v; want %d, %v", a, m, got, ok, want, exists)
			}
		}
	}
	m := uint64(math.MaxUint64 - 58) // 2^64 - 59 is prime.
	inverse, ok := ModInverse(uint64(3), m)
	if !ok || MulMod(inverse, 3, m) != 1 {
		t.Fatalf("ModInverse(3, 2^64-59) = %d, %v", inverse, ok)
	}
}

func TestModularOperationsRejectNonPositiveModulus(t *testing.T) {
	for name, f := range map[string]func(){
		"MulMod":     func() { MulMod(1, 1, 0) },
		"ModInverse": func() { ModInverse(1, -3) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("%s did not panic", name)
				}
			}()
			f()
		}()
	}
}
