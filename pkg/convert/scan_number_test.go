package convert

import (
	"errors"
	"strconv"
	"strings"
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

func TestScanNumberDecimalLengths(t *testing.T) {
	long := strings.Repeat("0", 64) + "42"
	for _, tt := range []struct {
		input string
		want  int64
	}{
		{"999999999999999999", 999999999999999999},
		{long, 42},
	} {
		got, err := ScanNumber[int64]([]byte(tt.input))
		if err != nil || got != tt.want {
			t.Fatalf("ScanNumber[int64](%q) = %d, %v; want %d", tt.input, got, err, tt.want)
		}
	}
	if got, err := ScanNumber[int8]([]byte(long)); err != nil || got != 42 {
		t.Fatalf("ScanNumber[int8](%q) = %d, %v; want 42", long, got, err)
	}
}

func checkScanNumberError[T constraints.Signed](t *testing.T, input, wantFunc string, wantCause error) {
	t.Helper()
	got, err := ScanNumber[T]([]byte(input))
	var numberError *strconv.NumError
	if got != 0 || !errors.As(err, &numberError) {
		t.Fatalf("ScanNumber(%q) = %d, %v; want zero and NumError", input, got, err)
	}
	if numberError.Func != wantFunc || numberError.Num != input || numberError.Err != wantCause {
		t.Fatalf("ScanNumber(%q) error = %#v; want Func=%q, Num=%q, Err=%v", input, numberError, wantFunc, input, wantCause)
	}
}

func TestScanNumberErrorDetails(t *testing.T) {
	t.Run("narrow range", func(t *testing.T) {
		checkScanNumberError[int8](t, "128", "ScanNumber", strconv.ErrRange)
	})
	t.Run("syntax before narrow range", func(t *testing.T) {
		checkScanNumberError[int8](t, "128x", "ParseInt", strconv.ErrSyntax)
	})
	t.Run("empty", func(t *testing.T) {
		checkScanNumberError[int64](t, "", "ParseInt", strconv.ErrSyntax)
	})
	t.Run("integer range", func(t *testing.T) {
		checkScanNumberError[int64](t, "9223372036854775808", "ParseInt", strconv.ErrRange)
	})
	for _, input := range []string{"9223372036854775808x", "18446744073709551616x"} {
		t.Run(input, func(t *testing.T) {
			_, err := strconv.ParseInt(input, 10, 64)
			var numberError *strconv.NumError
			if !errors.As(err, &numberError) {
				t.Fatalf("ParseInt(%q) = %v; want NumError", input, err)
			}
			checkScanNumberError[int64](t, input, numberError.Func, numberError.Err)
		})
	}
}
