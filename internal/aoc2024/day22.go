package aoc2024

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseSecretNumbers(r io.Reader) ([]int, error) {
	s := scanner.NewScanner(r, func(line []byte) (int, error) {
		n, err := convert.ScanNumber[int](line)
		if err != nil {
			return 0, fmt.Errorf("secret number: %w", err)
		}
		if n < 0 {
			return 0, fmt.Errorf("secret number %d is negative", n)
		}
		return n, nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func day22NextSecret(n int) int {
	n = (n ^ (n << 6)) & 16777215
	n = (n ^ (n >> 5)) & 16777215
	n = (n ^ (n << 11)) & 16777215
	return n
}

const day22Changes = 2000

func day22p01(r io.Reader) (string, error) {
	secretNumbers, err := parseSecretNumbers(r)
	if err != nil {
		return "", err
	}

	var total int
	for _, n := range secretNumbers {
		for range day22Changes {
			n = day22NextSecret(n)
		}
		total += n
	}
	return strconv.Itoa(total), nil
}

// day22Sequences is the number of four-change sequences; each change is a
// price difference in -9..9, so a sequence is a four-digit base-19 key.
const day22Sequences = 19 * 19 * 19 * 19

// day22BananaTotals returns, for each four-change sequence key, the bananas
// bought when every buyer sells at the first occurrence of that sequence.
func day22BananaTotals(secretNumbers []int) *[day22Sequences]int32 {
	totals := new([day22Sequences]int32)
	seen := new([day22Sequences]int32)
	for buyer, n := range secretNumbers {
		stamp := int32(buyer + 1)
		previous := n % 10
		key := 0
		for i := range day22Changes {
			n = day22NextSecret(n)
			price := n % 10
			key = key%(19*19*19)*19 + price - previous + 9
			previous = price
			if i >= 3 && seen[key] != stamp {
				seen[key] = stamp
				totals[key] += int32(price)
			}
		}
	}
	return totals
}

func day22p02(r io.Reader) (string, error) {
	secretNumbers, err := parseSecretNumbers(r)
	if err != nil {
		return "", err
	}
	if len(secretNumbers) == 0 {
		return "", fmt.Errorf("no secret numbers")
	}
	totals := day22BananaTotals(secretNumbers)
	return strconv.Itoa(int(slices.Max(totals[:]))), nil
}
