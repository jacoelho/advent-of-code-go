package aoc2024

import (
	"bytes"
	"errors"
	"io"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`),
			Want: "3749",
		},
		{
			Input: aoc.FileInput(t, 2024, 7),
			Want:  "945512582195",
		},
	}

	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`),
			Want: "11387",
		},
		{
			Input: aoc.FileInput(t, 2024, 7),
			Want:  "271691107779347",
		},
	}

	aoc.AOCTest(t, day07p02, tests)
}

func TestCalibrationRowsAndZero(t *testing.T) {
	tests := []struct {
		input string
		want  [2]string
	}{
		{"10: 5 5\n10: 2 5", [2]string{"20", "20"}},
		{"10: 2 0 10", [2]string{"10", "10"}},
		{"11: 2 0 10", [2]string{"0", "0"}},
		{"120: 12 0", [2]string{"0", "120"}},
		{"10: 9223372036854775807 2 0 10", [2]string{"10", "10"}},
	}
	for _, tt := range tests {
		for part, solve := range []func(io.Reader) (string, error){day07p01, day07p02} {
			got, err := solve(strings.NewReader(tt.input))
			if err != nil || got != tt.want[part] {
				t.Errorf("part %d input %q: got %q, %v; want %q", part+1, tt.input, got, err, tt.want[part])
			}
		}
	}
}

func TestCalibrationAgainstForwardEvaluation(t *testing.T) {
	combinations := 1
	for length := 1; length <= 4; length++ {
		combinations *= 4
		for encoded := range combinations {
			values := make([]int64, length)
			n := encoded
			for i := range values {
				values[i] = int64(n % 4)
				n /= 4
			}
			for _, operators := range []calibrationOperators{calibrationBasic, calibrationWithConcat} {
				reachable := make(map[string]bool)
				var evaluate func(int, *big.Int)
				evaluate = func(index int, value *big.Int) {
					if index == len(values) {
						reachable[value.String()] = true
						return
					}
					next := big.NewInt(values[index])
					evaluate(index+1, new(big.Int).Add(value, next))
					evaluate(index+1, new(big.Int).Mul(value, next))
					if operators == calibrationWithConcat {
						joined, ok := new(big.Int).SetString(value.String()+next.String(), 10)
						if !ok {
							t.Fatal("invalid decimal concatenation")
						}
						evaluate(index+1, joined)
					}
				}
				evaluate(1, big.NewInt(values[0]))
				for target := int64(0); target <= 50; target++ {
					equation := calibrationEquation{target: target, values: values}
					got, err := calibrationReachable(equation, operators, new(10_000))
					want := reachable[strconv.FormatInt(target, 10)]
					if err != nil || got != want {
						t.Fatalf("values %v, target %d, operators %d: got %v, %v; want %v", values, target, operators, got, err, want)
					}
				}
			}
		}
	}
}

func TestCalibrationInputAndTotal(t *testing.T) {
	for _, solve := range []func(io.Reader) (string, error){day07p01, day07p02} {
		for _, input := range []string{"10 5 5", "10:", "10: -1 11", "-1: 1"} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("input %q: got %q, want error", input, got)
			}
		}
		if _, err := solve(strings.NewReader("10: 5 x")); !errors.Is(err, strconv.ErrSyntax) {
			t.Errorf("got %v, want numeric syntax error", err)
		}
		const input = "9223372036854775807: 9223372036854775807"
		if got, err := solve(strings.NewReader(input)); err != nil || got != "9223372036854775807" {
			t.Errorf("got %q, %v; want MaxInt64", got, err)
		}
		if got, err := solve(strings.NewReader(input + "\n" + input)); err == nil {
			t.Errorf("got %q, want total overflow error", got)
		}
	}
}

func Benchmark_day07(b *testing.B) {
	input, err := io.ReadAll(aoc.FileInput(b, 2024, 7))
	if err != nil {
		b.Fatalf("read calibration input: %v", err)
	}

	const want = "271691107779347"
	for b.Loop() {
		got, err := day07p02(bytes.NewReader(input))
		if err != nil || got != want {
			b.Fatalf("day07p02() = %q, %v; want %q", got, err, want)
		}
	}
}
