package convert

import (
	"errors"
	"strconv"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestScanNumber(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  int
	}{
		{"0", 0}, {"-0", 0}, {"+42", 42}, {"-42", -42}, {"0012", 12},
	} {
		got, err := ScanNumber[int]([]byte(tc.input))
		if err != nil || got != tc.want {
			t.Errorf("ScanNumber(%q) = %d, %v; want %d", tc.input, got, err, tc.want)
		}
	}
	for _, input := range []string{"", "+", "-", "1x", "1 2", " 1", "1 ", "0x10", "1_000"} {
		if _, err := ScanNumber[int]([]byte(input)); !errors.Is(err, strconv.ErrSyntax) {
			t.Errorf("ScanNumber(%q) error = %v; want ErrSyntax", input, err)
		}
	}
}

func checkScanNumberBounds[T constraints.Signed](t *testing.T, low, high T, below, above string) {
	t.Helper()
	for _, want := range []T{low, high} {
		input := strconv.FormatInt(int64(want), 10)
		got, err := ScanNumber[T]([]byte(input))
		if err != nil || got != want {
			t.Errorf("ScanNumber[%T](%q) = %d, %v; want %d", want, input, got, err, want)
		}
	}
	for _, input := range []string{below, above} {
		if _, err := ScanNumber[T]([]byte(input)); !errors.Is(err, strconv.ErrRange) {
			t.Errorf("ScanNumber[%T](%q) error = %v; want ErrRange", low, input, err)
		}
	}
}

func TestScanNumberBounds(t *testing.T) {
	checkScanNumberBounds[int8](t, -128, 127, "-129", "128")
	checkScanNumberBounds[int16](t, -32768, 32767, "-32769", "32768")
	checkScanNumberBounds[int32](t, -2147483648, 2147483647, "-2147483649", "2147483648")
	checkScanNumberBounds[int64](t, -9223372036854775808, 9223372036854775807, "-9223372036854775809", "9223372036854775808")
	if strconv.IntSize == 64 {
		low, high := int64(-9223372036854775808), int64(9223372036854775807)
		checkScanNumberBounds[int](t, int(low), int(high), "-9223372036854775809", "9223372036854775808")
	} else {
		checkScanNumberBounds[int](t, -2147483648, 2147483647, "-2147483649", "2147483648")
	}
	type offset int16
	checkScanNumberBounds[offset](t, -32768, 32767, "-32769", "32768")
}
