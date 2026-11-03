package matrix

import (
	"math"
	"math/big"
	"testing"
)

func TestRatExactArithmetic(t *testing.T) {
	large := NewRat(math.MaxInt64, 1)
	got := large.Mul(NewRat(2, 1)).Div(NewRat(2, 1))
	if got.Cmp(large) != 0 {
		t.Fatalf("large product/quotient = %s/%s", got.Numerator(), got.Denominator())
	}
	got = NewRat(math.MinInt64, -1)
	if !got.IsInt() || got.Numerator().String() != "9223372036854775808" {
		t.Fatalf("normalized minimum integer = %s/%s", got.Numerator(), got.Denominator())
	}

	var zero Rat
	fraction := NewRat(2, 3)
	if zero.Add(fraction).Cmp(fraction) != 0 || !zero.Mul(fraction).IsZero() || !zero.Neg().IsZero() || !zero.Div(fraction).IsZero() {
		t.Fatal("zero-value arithmetic differs from exact zero")
	}
	if zero.Denominator().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("zero denominator = %s; want 1", zero.Denominator())
	}
}

func TestRatOwnership(t *testing.T) {
	source := big.NewInt(7)
	value := NewInteger(source)
	fraction := value.Div(NewRat(3, 1))
	source.SetInt64(42)
	value.Numerator().SetInt64(99)
	fraction.Denominator().SetInt64(1)
	if value.Cmp(NewRat(7, 1)) != 0 || fraction.Cmp(NewRat(7, 3)) != 0 {
		t.Fatal("caller mutation changed an owned rational")
	}
	if got := fraction.Add(NewRat(2, 3)); got.Cmp(NewRat(3, 1)) != 0 || fraction.Cmp(NewRat(7, 3)) != 0 {
		t.Fatal("arithmetic mutated its operand")
	}
}
