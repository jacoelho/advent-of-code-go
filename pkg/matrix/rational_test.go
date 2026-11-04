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

func wantRatValue(t *testing.T, got Rat, want string) {
	t.Helper()
	expected, ok := new(big.Rat).SetString(want)
	if !ok {
		t.Fatalf("invalid expected rational %q", want)
	}
	if got.Numerator().Cmp(expected.Num()) != 0 || got.Denominator().Cmp(expected.Denom()) != 0 {
		t.Fatalf("got %s/%s; want %s", got.Numerator(), got.Denominator(), want)
	}
}

func TestRatIntegerBoundaries(t *testing.T) {
	minimum, maximum := NewRat(math.MinInt64, 1), NewRat(math.MaxInt64, 1)
	one, negativeOne := NewRat(1, 1), NewRat(-1, 1)
	for _, tt := range []struct {
		name string
		got  Rat
		want string
	}{
		{"addition overflow", maximum.Add(one), "9223372036854775808"},
		{"addition underflow", minimum.Add(negativeOne), "-9223372036854775809"},
		{"subtraction overflow", maximum.Sub(negativeOne), "9223372036854775808"},
		{"subtraction underflow", minimum.Sub(one), "-9223372036854775809"},
		{"mixed subtraction", maximum.Add(one).Sub(maximum), "1"},
		{"minimum negation", minimum.Neg(), "9223372036854775808"},
		{"minimum multiplication", minimum.Mul(negativeOne), "9223372036854775808"},
		{"reversed minimum multiplication", negativeOne.Mul(minimum), "9223372036854775808"},
		{"product overflow", maximum.Mul(NewRat(2, 1)), "18446744073709551614"},
		{"minimum division", minimum.Div(negativeOne), "9223372036854775808"},
		{"fractional division", one.Div(NewRat(2, 1)), "1/2"},
		{"negative denominator", NewRat(4, -6), "-2/3"},
		{"minimum denominator", NewRat(1, math.MinInt64), "-1/9223372036854775808"},
		{"minimum numerator and denominator", NewRat(math.MinInt64, math.MinInt64), "1"},
	} {
		t.Run(tt.name, func(t *testing.T) { wantRatValue(t, tt.got, tt.want) })
	}
}

func TestRatMultiplicationBoundaries(t *testing.T) {
	const fits, overflows int64 = 3037000499, 3037000500
	for _, tt := range []struct {
		name        string
		left, right int64
		want        string
	}{
		{"positive fit", fits, fits, "9223372030926249001"},
		{"negative fit", -fits, fits, "-9223372030926249001"},
		{"positive overflow", overflows, overflows, "9223372037000250000"},
		{"two negative operands overflow", -overflows, -overflows, "9223372037000250000"},
		{"negative left overflow", -overflows, overflows, "-9223372037000250000"},
		{"negative right overflow", overflows, -overflows, "-9223372037000250000"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			wantRatValue(t, NewRat(tt.left, 1).Mul(NewRat(tt.right, 1)), tt.want)
		})
	}
}

func TestRatComparisonBoundaries(t *testing.T) {
	minimum, maximum := NewRat(math.MinInt64, 1), NewRat(math.MaxInt64, 1)
	aboveMaximum := maximum.Add(NewRat(1, 1))
	if minimum.Cmp(maximum) >= 0 || maximum.Cmp(minimum) <= 0 ||
		aboveMaximum.Cmp(maximum) <= 0 || maximum.Cmp(aboveMaximum) >= 0 {
		t.Fatal("comparison reversed integer bounds")
	}
	large, ok := new(big.Int).SetString("1208925819614629174706176", 10)
	if !ok {
		t.Fatal("invalid large integer fixture")
	}
	p := NewInteger(large)
	q := NewInteger(new(big.Int).Add(large, big.NewInt(1)))
	negativeP, negativeQ := p.Neg(), q.Neg()
	copiedP, copiedNegativeP := p, negativeP
	for _, tt := range []struct {
		name        string
		left, right Rat
		want        int
	}{
		{"copied positive", copiedP, p, 0},
		{"copied negative", copiedNegativeP, negativeP, 0},
		{"separately constructed equal", p, NewInteger(large), 0},
		{"large positive ordering", p, q, -1},
		{"large negative ordering", negativeP, negativeQ, 1},
		{"large opposite signs", negativeP, p, -1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.left.Cmp(tt.right); got != tt.want {
				t.Fatalf("Cmp = %d; want %d", got, tt.want)
			}
			if got := tt.right.Cmp(tt.left); got != -tt.want {
				t.Fatalf("reversed Cmp = %d; want %d", got, -tt.want)
			}
		})
	}
}

func TestRatCancellationAfterOverflow(t *testing.T) {
	promoted := NewRat(math.MaxInt64, 1).Add(NewRat(1, 1))
	wantRatValue(t, promoted, "9223372036854775808")
	for _, zero := range []Rat{promoted.Add(NewRat(math.MinInt64, 1)), promoted.Sub(promoted)} {
		wantRatValue(t, zero, "0")
		if !zero.IsZero() {
			t.Fatal("exact cancellation did not report zero")
		}
	}
	wantRatValue(t, promoted, "9223372036854775808")
}

func TestRatFractionArithmetic(t *testing.T) {
	type fractionValue struct {
		value Rat
		want  *big.Rat
	}
	var values []fractionValue
	for _, pair := range []struct{ numerator, denominator int64 }{
		{0, 1}, {1, 1}, {-1, 1}, {1, 2}, {-1, 2}, {2, 3}, {-2, 3}, {3, 5},
	} {
		values = append(values, fractionValue{NewRat(pair.numerator, pair.denominator), big.NewRat(pair.numerator, pair.denominator)})
	}
	half, third := solveDiagonalRats(t, 2, 1, 3, 1)
	values = append(values, fractionValue{half, big.NewRat(1, 2)}, fractionValue{third, big.NewRat(1, 3)})
	for _, a := range values {
		left, wantLeft := a.value, a.want
		if got, want := left.Sign(), wantLeft.Sign(); got != want {
			t.Fatalf("Sign(%s) = %d; want %d", wantLeft.RatString(), got, want)
		}
		wantRatValue(t, left.Neg(), new(big.Rat).Neg(wantLeft).RatString())
		for _, b := range values {
			right, wantRight := b.value, b.want
			t.Run(wantLeft.RatString()+","+wantRight.RatString(), func(t *testing.T) {
				wantRatValue(t, left.Add(right), new(big.Rat).Add(wantLeft, wantRight).RatString())
				wantRatValue(t, left.Sub(right), new(big.Rat).Sub(wantLeft, wantRight).RatString())
				wantRatValue(t, left.Mul(right), new(big.Rat).Mul(wantLeft, wantRight).RatString())
				if wantRight.Sign() != 0 {
					wantRatValue(t, left.Div(right), new(big.Rat).Quo(wantLeft, wantRight).RatString())
				}
				if got, want := left.Cmp(right), wantLeft.Cmp(wantRight); got != want {
					t.Fatalf("Cmp = %d; want %d", got, want)
				}
				wantRatValue(t, left, wantLeft.RatString())
				wantRatValue(t, right, wantRight.RatString())
			})
		}
	}
}

func solveDiagonalRats(t *testing.T, ax, cx, by, cy int64) (Rat, Rat) {
	t.Helper()
	x, y := Variable("x"), Variable("y")
	solution, err := NewSystem(x, y).Add(x.Scale(ax).Equal(cx)).Add(y.Scale(by).Equal(cy)).Solve()
	if err != nil {
		t.Fatal(err)
	}
	xExpression, err := solution.Expression(x)
	if err != nil {
		t.Fatal(err)
	}
	yExpression, err := solution.Expression(y)
	if err != nil {
		t.Fatal(err)
	}
	return xExpression.Constant(), yExpression.Constant()
}

func TestSolvedFractionsRemainReusable(t *testing.T) {
	half, third := solveDiagonalRats(t, 2, 1, 3, 1)
	if half.IsInt() || third.IsInt() || !half.Add(half).IsInt() || !half.Sub(half).IsZero() {
		t.Fatal("fraction arithmetic changed integer or zero classification")
	}
	a := Variable("a")
	solution, err := NewSystem(a).Add(a.ScaleRat(half).Add(a.ScaleRat(third)).Equal(5)).Solve()
	if err != nil {
		t.Fatal(err)
	}
	values, err := solution.Bind(nil)
	if err != nil {
		t.Fatal(err)
	}
	wantRatValue(t, values[a], "6")
	wantRatValue(t, half, "1/2")
	wantRatValue(t, third, "1/3")
}

func TestSolvedFractionSignBoundaries(t *testing.T) {
	for _, tt := range []struct {
		name           string
		ax, cx, by, cy int64
		wantX, wantY   string
	}{
		{"negative determinant", -2, 1, 3, 1, "-1/2", "1/3"},
		{"oversized positive numerator", -3, math.MinInt64, 1, 0, "9223372036854775808/3", "0"},
		{"oversized positive denominator", math.MinInt64, 1, 1, 0, "-1/9223372036854775808", "0"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			x, y := solveDiagonalRats(t, tt.ax, tt.cx, tt.by, tt.cy)
			wantRatValue(t, x, tt.wantX)
			wantRatValue(t, y, tt.wantY)
		})
	}
}

func TestRatFractionBoundaries(t *testing.T) {
	maximumHalf := NewRat(math.MaxInt64, 2)
	minimumThird := NewRat(math.MinInt64, 3)
	for _, tt := range []struct {
		name string
		got  Rat
		want string
	}{
		{"addition reduces overflow", maximumHalf.Add(maximumHalf), "9223372036854775807"},
		{"subtraction reduces underflow", minimumThird.Sub(NewRat(1, 3)), "-3074457345618258603"},
		{"multiplication cancels factors", maximumHalf.Mul(NewRat(2, math.MaxInt64)), "1"},
		{"denominator overflow", NewRat(1, math.MaxInt64).Mul(NewRat(1, 2)), "1/18446744073709551614"},
		{"minimum numerator negation", minimumThird.Neg(), "9223372036854775808/3"},
		{"reduced minimum denominator", NewRat(-2, math.MinInt64), "1/4611686018427387904"},
		{"negative reduced minimum denominator", NewRat(2, math.MinInt64), "-1/4611686018427387904"},
	} {
		t.Run(tt.name, func(t *testing.T) { wantRatValue(t, tt.got, tt.want) })
	}
	small := NewRat(1, math.MaxInt64)
	if maximumHalf.Cmp(small) <= 0 || small.Cmp(maximumHalf) >= 0 {
		t.Fatal("fraction comparison reversed values whose cross-products overflow")
	}
}

func TestRatDivisionByZero(t *testing.T) {
	for _, tt := range []struct {
		name string
		run  func()
	}{
		{"nonzero numerator", func() { NewRat(1, 0) }},
		{"zero numerator", func() { NewRat(0, 0) }},
		{"integer divisor", func() { NewRat(1, 1).Div(Rat{}) }},
		{"fraction divisor", func() { NewRat(1, 2).Div(Rat{}) }},
		{"zero divided by zero", func() { (Rat{}).Div(Rat{}) }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("division by zero did not panic")
				}
			}()
			tt.run()
		})
	}
}

func TestRatOwnership(t *testing.T) {
	for _, tt := range []struct {
		value    string
		fraction string
		sum      string
	}{
		{"7", "7/3", "3"},
		{"1208925819614629174706176", "1208925819614629174706176/3", "402975273204876391568726"},
	} {
		t.Run(tt.value, func(t *testing.T) {
			source, ok := new(big.Int).SetString(tt.value, 10)
			if !ok {
				t.Fatal("invalid integer fixture")
			}
			value := NewInteger(source)
			fraction := value.Div(NewRat(3, 1))
			source.SetInt64(42)
			value.Numerator().SetInt64(99)
			fraction.Denominator().SetInt64(1)
			wantRatValue(t, value, tt.value)
			wantRatValue(t, fraction, tt.fraction)
			wantRatValue(t, fraction.Add(NewRat(2, 3)), tt.sum)
			wantRatValue(t, fraction, tt.fraction)
		})
	}
}
