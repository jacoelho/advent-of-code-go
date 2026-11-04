package matrix

// Nonintegral inputs, singular systems and native overflow return false so
// general elimination can handle them. The augmented rows are read-only.
func solveTwoIntegers(rows *[2][3]Rat) (Rat, Rat, bool) {
	a, b, c := rows[0][0], rows[0][1], rows[0][2]
	d, e, f := rows[1][0], rows[1][1], rows[1][2]
	if a.value != nil || a.denominator != 0 ||
		b.value != nil || b.denominator != 0 ||
		c.value != nil || c.denominator != 0 ||
		d.value != nil || d.denominator != 0 ||
		e.value != nil || e.denominator != 0 ||
		f.value != nil || f.denominator != 0 {
		return Rat{}, Rat{}, false
	}
	determinant, ok := subtractProducts(a.numerator, e.numerator, b.numerator, d.numerator)
	if !ok || determinant == 0 {
		return Rat{}, Rat{}, false
	}
	x, ok := subtractProducts(c.numerator, e.numerator, b.numerator, f.numerator)
	if !ok {
		return Rat{}, Rat{}, false
	}
	y, ok := subtractProducts(a.numerator, f.numerator, c.numerator, d.numerator)
	if !ok {
		return Rat{}, Rat{}, false
	}
	return terminalRat(x, determinant), terminalRat(y, determinant), true
}

func subtractProducts(a, b, c, d int64) (int64, bool) {
	left, ok := multiplyIntegers(a, b)
	if !ok {
		return 0, false
	}
	right, ok := multiplyIntegers(c, d)
	if !ok {
		return 0, false
	}
	return subtractIntegers(left, right)
}
